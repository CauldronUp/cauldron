package store

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// alphabet is the character set provider identifiers are drawn from. Digits and
// both cases, matching what Stripe, Shopify and most others actually emit —
// applications occasionally validate the shape, so the fakes have to match it.
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// shape describes how one resource's identifiers look.
type shape struct {
	// style is prefixed (Stripe object ids), numeric (GitHub issue numbers),
	// timestamp (Slack message ts) or opaque (SendGrid message ids, which
	// carry no prefix at all). Providers genuinely differ, applications parse
	// all of them, so the emulator has to match.
	style  string
	prefix string
	length int
	seq    int
	drawn  int
}

// Generator mints identifiers that look like the provider's own but are fully
// determined by the seed.
//
// Using math/rand here is deliberate and is not a security weakness: these are
// fixture identifiers inside a sandbox, and reproducibility is the entire
// point. Anything unpredictable would make CI failures impossible to replay.
type Generator struct {
	mu     sync.Mutex
	seed   int64
	rng    *rand.Rand
	shapes map[string]shape
	now    func() int64
}

// NewGenerator returns a generator for the given seed.
func NewGenerator(seed int64) *Generator {
	return &Generator{
		seed:   seed,
		rng:    rand.New(rand.NewSource(seed)),
		shapes: map[string]shape{},
	}
}

// Declare registers the identifier shape for a resource.
func (g *Generator) Declare(resource, prefix string, length int) {
	g.DeclareStyle(resource, "prefixed", prefix, length)
}

// DeclareStyle registers the identifier shape, including its style.
func (g *Generator) DeclareStyle(resource, style, prefix string, length int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if length <= 0 {
		length = 14
	}

	g.shapes[resource] = shape{style: style, prefix: prefix, length: length}
}

// Now is the sandbox clock, used by the timestamp style. It is injected rather
// than read from the wall clock so identifiers stay reproducible.
func (g *Generator) Now(now func() int64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.now = now
}

// Next mints the next identifier for a resource.
func (g *Generator) Next(resource string) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	s, ok := g.shapes[resource]
	if !ok {
		s = shape{prefix: resource + "_", length: 14}
	}

	s.drawn++

	switch s.style {
	case "numeric":
		s.seq++
		g.shapes[resource] = s

		return strconv.Itoa(s.seq)

	case "timestamp":
		// Slack identifies a message by the second it arrived plus a counter,
		// which is why thread_ts is a string that looks like a float. The
		// counter keeps two messages in the same second distinct, and comes
		// from the sandbox clock so the value is reproducible.
		s.seq++
		g.shapes[resource] = s

		seconds := int64(0)
		if g.now != nil {
			seconds = g.now()
		}

		return strconv.FormatInt(seconds, 10) + "." + fmt.Sprintf("%06d", s.seq)
	}

	g.shapes[resource] = s

	var b strings.Builder

	b.WriteString(s.prefix)

	for i := 0; i < s.length; i++ {
		b.WriteByte(alphabet[g.rng.Intn(len(alphabet))])
	}

	return b.String()
}

// Reset rewinds the generator to its initial state, so a reset sandbox
// reproduces the same identifiers as a fresh one.
func (g *Generator) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.rng = rand.New(rand.NewSource(g.seed))

	for resource, s := range g.shapes {
		s.seq = 0
		s.drawn = 0
		g.shapes[resource] = s
	}
}

// Drawn reports how many identifiers have been minted per resource.
func (g *Generator) Drawn() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()

	out := make(map[string]int, len(g.shapes))

	for resource, s := range g.shapes {
		out[resource] = s.drawn
	}

	return out
}

// Restore rewinds the generator and replays the recorded number of draws, so
// that identifiers minted after a restore continue from where the snapshot
// left off rather than repeating ones already in use.
func (g *Generator) Restore(drawn map[string]int) {
	g.Reset()

	g.mu.Lock()
	shapes := make([]string, 0, len(g.shapes))
	for resource := range g.shapes {
		shapes = append(shapes, resource)
	}
	g.mu.Unlock()

	// Sorted so the replay is identical on every machine: the shared random
	// source means the order draws happen in changes what they produce.
	sort.Strings(shapes)

	for _, resource := range shapes {
		for i := 0; i < drawn[resource]; i++ {
			g.Next(resource)
		}
	}
}
