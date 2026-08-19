package runtime

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CauldronUp/cauldron/internal/clock"
)

// Conditions is a degraded network, armed against one sandbox. It is what
// `cauldron network stripe --latency 800 --jitter 200` installs.
//
// Cauldron's fault injection already covers what a provider does *deliberately*
// when something is wrong: rate limits, validation errors, 5xx. This covers
// what the network does to you regardless of what the provider intended, which
// is a different and usually worse class of failure. A rate limit arrives as a
// well-formed 429 your client library already understands. A connection that
// hangs for ninety seconds and then dies arrives as nothing at all, and that is
// the one that takes production down.
//
// The vocabulary is deliberately Toxiproxy's — latency, jitter, bandwidth,
// timeout, reset, slice, limit — because that is the vocabulary the ecosystem
// already has for this, and somebody who has used Toxiproxy against their
// Postgres should not have to learn a second set of words to use Cauldron
// against their payment provider. See docs/network.md for what each one means
// here, including the two places an HTTP emulator cannot be a TCP proxy.
type Conditions struct {
	// Latency delays the response by this long.
	Latency time.Duration
	// Jitter varies the delay by ± this much.
	Jitter time.Duration
	// Bandwidth caps the response body in kilobytes per second. Zero is
	// unlimited; a positive value throttles.
	Bandwidth int
	// Timeout accepts the request, answers nothing, and closes the connection
	// after this long. A client with no timeout of its own hangs for the whole
	// duration, which is the point.
	Timeout time.Duration
	// Reset closes the connection immediately without a response, the way a
	// peer that has just been restarted does.
	Reset bool
	// Limit closes the connection after this many bytes of body. Models a
	// response truncated mid-flight.
	Limit int64
	// Slice writes the body in chunks of roughly this many bytes, which
	// surfaces clients that assume a whole response arrives in one read.
	Slice int
	// Probability is the share of requests affected, from 0 to 1. Zero means
	// every request, so that the common case needs no flag.
	Probability float64
	// Path, when set, restricts the conditions to routes whose path contains it.
	Path string
	// Until is the sandbox time the conditions expire. Zero means they do not.
	Until time.Time
	// Count limits how many requests are affected. Zero means unlimited.
	Count int

	seen int
}

// Degrades reports whether these conditions would do anything at all. Arming
// an empty set is almost always a mistyped flag, and silently accepting it
// would leave somebody waiting for a failure that is never coming.
func (c Conditions) Degrades() bool {
	return c.Latency > 0 ||
		c.Jitter > 0 ||
		c.Bandwidth > 0 ||
		c.Timeout > 0 ||
		c.Reset ||
		c.Limit > 0 ||
		c.Slice > 0
}

// Fatal reports whether the conditions end the request without a response.
// Those are handled before the recipe runs, because there is no point building
// a body nobody will read.
func (c Conditions) Fatal() bool {
	return c.Reset || c.Timeout > 0
}

// networkSet holds the conditions armed against one sandbox.
type networkSet struct {
	mu         sync.Mutex
	conditions []*Conditions
	clock      *clock.Clock
	rand       *rand.Rand
}

func newNetworkSet(c *clock.Clock, seed int64) *networkSet {
	return &networkSet{
		clock: c,
		// Seeded from the sandbox seed, so a run with a probability below 1 is
		// as reproducible as the identifiers are. Chaos that cannot be replayed
		// is a bug report nobody can act on.
		rand: rand.New(rand.NewSource(seed)), // #nosec G404 -- not security-sensitive
	}
}

// Arm installs a set of conditions.
func (n *networkSet) Arm(c Conditions) {
	n.mu.Lock()
	defer n.mu.Unlock()

	copied := c
	n.conditions = append(n.conditions, &copied)
}

// Clear removes every armed condition.
func (n *networkSet) Clear() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.conditions = nil
}

// Armed returns the currently armed conditions.
func (n *networkSet) Armed() []Conditions {
	n.mu.Lock()
	defer n.mu.Unlock()

	out := make([]Conditions, 0, len(n.conditions))

	for _, c := range n.conditions {
		out = append(out, *c)
	}

	return out
}

// next returns the conditions to apply to a request, if any.
//
// Mirrors faultSet.next: expired and spent entries are dropped first so a dead
// one never consumes a slot from the entries behind it.
func (n *networkSet) next(path string) (Conditions, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := n.clock.Now()
	live := n.conditions[:0]

	for _, c := range n.conditions {
		if !c.Until.IsZero() && !now.Before(c.Until) {
			continue
		}

		if c.Count > 0 && c.seen >= c.Count {
			continue
		}

		live = append(live, c)
	}

	n.conditions = live

	for _, c := range n.conditions {
		if c.Path != "" && !contains(path, c.Path) {
			continue
		}

		c.seen++

		if c.Probability > 0 && c.Probability < 1 && n.rand.Float64() >= c.Probability {
			continue
		}

		return *c, true
	}

	return Conditions{}, false
}

// delay is the pause to apply before answering, jitter included.
func (n *networkSet) delay(c Conditions) time.Duration {
	if c.Latency <= 0 && c.Jitter <= 0 {
		return 0
	}

	d := c.Latency

	if c.Jitter > 0 {
		n.mu.Lock()
		offset := time.Duration(n.rand.Int63n(int64(c.Jitter)*2+1)) - c.Jitter
		n.mu.Unlock()

		d += offset
	}

	if d < 0 {
		return 0
	}

	return d
}

// hangUp ends the request without a response.
//
// Hijacking is the only way to do this: writing nothing and returning would
// still produce a 200 with an empty body, which is a completely different thing
// for the client under test to receive. With the connection taken over and
// closed, the client sees exactly what it sees when a peer goes away.
func hangUp(w http.ResponseWriter, after time.Duration) bool {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return false
	}

	conn, buffered, err := hijacker.Hijack()
	if err != nil {
		return false
	}

	if after > 0 {
		// The response is already unsalvageable; all this wait does is hold the
		// client's socket open, which is the failure being emulated.
		time.Sleep(after)
	}

	// SetLinger(0) makes the close a TCP RST rather than a FIN, so the client
	// sees "connection reset by peer" instead of a clean EOF. Clients treat the
	// two differently, and the difference is often the bug.
	if tcp, ok := conn.(*net.TCPConn); ok {
		_ = tcp.SetLinger(0)
	}

	_ = buffered.Flush()
	_ = conn.Close()

	return true
}

// degradedWriter applies bandwidth, byte limit and slicing to a response body.
//
// Headers are written normally: throttling those would model a slow link more
// completely, but it would also break the ResponseWriter contract, and the body
// is where a client's read loop actually lives.
type degradedWriter struct {
	http.ResponseWriter

	bytesPerSecond float64
	limit          int64
	slice          int

	written int64
	closed  bool
}

func newDegradedWriter(w http.ResponseWriter, c Conditions) *degradedWriter {
	d := &degradedWriter{ResponseWriter: w, limit: c.Limit, slice: c.Slice}

	if c.Bandwidth > 0 {
		d.bytesPerSecond = float64(c.Bandwidth) * 1024
	}

	return d
}

func (d *degradedWriter) Write(p []byte) (int, error) {
	if d.closed {
		// The byte limit was reached on an earlier write. Report success so the
		// handler finishes normally; the client has already lost the rest.
		return len(p), nil
	}

	if d.limit > 0 {
		remaining := d.limit - d.written

		if remaining <= 0 {
			d.closed = true

			return len(p), nil
		}

		if int64(len(p)) > remaining {
			written, err := d.writeThrottled(p[:remaining])
			d.closed = true

			if err != nil {
				return written, err
			}

			return len(p), nil
		}
	}

	return d.writeThrottled(p)
}

// writeThrottled writes in slice-sized pieces, sleeping between them to hold
// the average rate. Sleeping per chunk rather than once up front is what makes
// the response arrive gradually instead of all at once but late.
func (d *degradedWriter) writeThrottled(p []byte) (int, error) {
	chunk := d.slice

	if chunk <= 0 {
		chunk = len(p)
	}

	if chunk <= 0 {
		chunk = 1
	}

	total := 0

	for start := 0; start < len(p); start += chunk {
		end := start + chunk
		if end > len(p) {
			end = len(p)
		}

		written, err := d.ResponseWriter.Write(p[start:end])
		total += written
		d.written += int64(written)

		if err != nil {
			return total, err
		}

		// Flushing makes each piece a separate chunk on the wire. Without it Go
		// buffers the lot and the slicing would be invisible to the client.
		if flusher, ok := d.ResponseWriter.(http.Flusher); ok {
			flusher.Flush()
		}

		if d.bytesPerSecond > 0 && written > 0 {
			time.Sleep(time.Duration(float64(written) / d.bytesPerSecond * float64(time.Second)))
		}
	}

	return total, nil
}

// Flush satisfies http.Flusher so streaming handlers keep working.
func (d *degradedWriter) Flush() {
	if flusher, ok := d.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// truncated reports whether the byte limit cut the response short.
func (d *degradedWriter) truncated() bool { return d.closed }

// Describe renders armed conditions the way a person reads them, for the
// request log, `cauldron status` and the CLI's own confirmation line.
func Describe(c Conditions) string { return describeConditions(c) }

func describeConditions(c Conditions) string {
	var parts []string

	if c.Latency > 0 || c.Jitter > 0 {
		if c.Jitter > 0 {
			parts = append(parts, "latency "+c.Latency.String()+" ±"+c.Jitter.String())
		} else {
			parts = append(parts, "latency "+c.Latency.String())
		}
	}

	if c.Bandwidth > 0 {
		parts = append(parts, "bandwidth "+strconv.Itoa(c.Bandwidth)+"KB/s")
	}

	if c.Timeout > 0 {
		parts = append(parts, "timeout "+c.Timeout.String())
	}

	if c.Reset {
		parts = append(parts, "reset")
	}

	if c.Limit > 0 {
		parts = append(parts, "limit "+strconv.FormatInt(c.Limit, 10)+"B")
	}

	if c.Slice > 0 {
		parts = append(parts, "slice "+strconv.Itoa(c.Slice)+"B")
	}

	if len(parts) == 0 {
		return ""
	}

	out := strings.Join(parts, ", ")

	if c.Probability > 0 && c.Probability < 1 {
		out += fmt.Sprintf(" @ %.0f%% of requests", c.Probability*100)
	}

	if c.Path != "" {
		out += " on paths containing " + c.Path
	}

	return out
}
