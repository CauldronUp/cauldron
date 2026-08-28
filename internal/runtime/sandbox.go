// Package runtime serves a Recipe over HTTP: routing, auth, state, pagination,
// failure injection and webhooks.
//
// One process hosts many sandboxes, routed by name, rather than a container per
// provider. Booting eight fakes should cost milliseconds and a few megabytes,
// because in CI it happens on every commit.
package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CauldronUp/cauldron/internal/clock"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// Options configures a sandbox.
type Options struct {
	// Seed determines every generated identifier. The same seed and fixture
	// produce byte-identical state.
	Seed int64
	// Clock is shared when several sandboxes must agree on time. Optional.
	Clock *clock.Clock
	// LogSize caps how many exchanges are retained.
	LogSize int
}

// Sandbox is one running Recipe.
type Sandbox struct {
	recipe  *recipe.Recipe
	store   *store.Store
	clock   *clock.Clock
	router  *router
	faults  *faultSet
	network *networkSet
	log     *requestLog

	mu       sync.RWMutex
	fixture  string
	webhooks *webhookQueue

	// idShapes holds the compiled form of each resource's declared identifier
	// pattern, by resource name. Absent means the provider looks up whatever
	// it is given, which is most of them.
	idShapes map[string]*regexp.Regexp
}

// New builds a sandbox from a validated Recipe.
func New(r *recipe.Recipe, opts Options) (*Sandbox, error) {
	if r == nil {
		return nil, fmt.Errorf("runtime: recipe is nil")
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	c := opts.Clock
	if c == nil {
		c = clock.New()
	}

	s := &Sandbox{
		recipe:   r,
		store:    store.New(opts.Seed),
		clock:    c,
		router:   newRouter(r),
		faults:   newFaultSet(c),
		network:  newNetworkSet(c, opts.Seed),
		log:      newRequestLog(opts.LogSize),
		webhooks: newWebhookQueue(r, c),
	}

	// The identifier generator reads the sandbox clock, never the wall clock,
	// so a time-based identifier is as reproducible as any other.
	s.store.Clock(func() int64 { return s.clock.Unix() })

	for name, resource := range r.Resources {
		s.store.DeclareStyle(name, resource.ID.Style, resource.ID.Prefix, resource.ID.Length)

		if resource.ID.Pattern == "" {
			continue
		}

		// Validate has already refused a pattern that will not compile, so
		// reaching this is a Recipe loaded past it rather than a Recipe
		// mistake. Refusing here rather than ignoring it keeps the failure
		// loud: a pattern silently dropped would answer 404 to a malformed
		// id, which is the exact behaviour declaring one exists to change.
		shape, err := regexp.Compile(resource.ID.Pattern)
		if err != nil {
			return nil, fmt.Errorf("runtime: resource %s: id.pattern: %w", name, err)
		}

		if s.idShapes == nil {
			s.idShapes = map[string]*regexp.Regexp{}
		}

		s.idShapes[name] = shape
	}

	return s, nil
}

// Name returns the Recipe's name.
func (s *Sandbox) Name() string { return s.recipe.Name }

// Recipe returns the Recipe being served.
func (s *Sandbox) Recipe() *recipe.Recipe { return s.recipe }

// Clock returns the sandbox clock.
func (s *Sandbox) Clock() *clock.Clock { return s.clock }

// Store returns the sandbox state.
func (s *Sandbox) Store() *store.Store { return s.store }

// Webhooks returns the webhook queue.
func (s *Sandbox) Webhooks() *webhookQueue { return s.webhooks }

// Fixture returns the name of the fixture currently loaded.
func (s *Sandbox) Fixture() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.fixture
}

// Seed loads a named fixture, replacing all current state.
func (s *Sandbox) Seed(name string) error {
	fixture, ok := s.recipe.Fixtures[name]
	if !ok {
		return fmt.Errorf("recipe %s has no fixture %q (available: %s)", s.recipe.Name, name, strings.Join(s.Fixtures(), ", "))
	}

	s.store.Reset()

	// Resource order is sorted so seeding is deterministic even though the
	// fixture is a map.
	resources := make([]string, 0, len(fixture))
	for resource := range fixture {
		resources = append(resources, resource)
	}

	sort.Strings(resources)

	for _, resource := range resources {
		for _, raw := range fixture[resource] {
			record := store.Record{}

			for key, value := range raw {
				record[key] = value
			}

			s.applyDefaults(resource, record)

			if _, err := s.store.Create(resource, record); err != nil {
				return fmt.Errorf("seeding %s: %w", resource, err)
			}
		}
	}

	s.mu.Lock()
	s.fixture = name
	s.mu.Unlock()

	return nil
}

// Fixtures returns the fixture names this Recipe ships.
func (s *Sandbox) Fixtures() []string {
	out := make([]string, 0, len(s.recipe.Fixtures))

	for name := range s.recipe.Fixtures {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

// Reset returns the sandbox to its seeded state, clears faults, empties the
// request log and rewinds the clock. A reset sandbox is indistinguishable from
// a freshly booted one — that is what makes it safe to reuse between tests.
func (s *Sandbox) Reset() error {
	fixture := s.Fixture()

	s.faults.Clear()
	s.log.reset()
	s.webhooks.reset()
	s.clock.Reset()
	s.store.Reset()

	if fixture == "" {
		return nil
	}

	return s.Seed(fixture)
}

// Arm installs a failure mode.
func (s *Sandbox) Arm(fault Fault) error {
	if _, ok := s.recipe.Errors[fault.Error]; !ok {
		return fmt.Errorf("recipe %s has no error %q (available: %s)", s.recipe.Name, fault.Error, strings.Join(s.Errors(), ", "))
	}

	s.faults.Arm(fault)

	return nil
}

// ClearFaults disarms every fault.
func (s *Sandbox) ClearFaults() { s.faults.Clear() }

// Degrade arms network conditions: latency, throttling, timeouts, resets.
//
// Unlike Arm, this needs nothing from the Recipe. A slow link is not something
// a provider declares; it happens to every provider equally.
func (s *Sandbox) Degrade(c Conditions) error {
	if !c.Degrades() {
		return fmt.Errorf("recipe %s: nothing to degrade (set at least one of latency, jitter, bandwidth, timeout, reset, limit, slice)", s.recipe.Name)
	}

	if c.Probability < 0 || c.Probability > 1 {
		return fmt.Errorf("recipe %s: probability must be between 0 and 1, got %v", s.recipe.Name, c.Probability)
	}

	s.network.Arm(c)

	return nil
}

// ClearNetwork restores an undegraded network.
func (s *Sandbox) ClearNetwork() { s.network.Clear() }

// ArmedNetwork returns the currently armed network conditions.
func (s *Sandbox) ArmedNetwork() []Conditions { return s.network.Armed() }

// ArmedFaults returns the currently armed faults.
func (s *Sandbox) ArmedFaults() []Fault { return s.faults.Armed() }

// Errors returns the injectable failure names.
func (s *Sandbox) Errors() []string {
	out := make([]string, 0, len(s.recipe.Errors))

	for name := range s.recipe.Errors {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

// Exchanges returns recent requests.
func (s *Sandbox) Exchanges(limit int) []Exchange { return s.log.Entries(limit) }

// applyDefaults fills unset fields declared with a default, and stamps created
// timestamps from the sandbox clock rather than the wall clock. Declared
// constants are stamped last and overwrite whatever was supplied, because the
// provider does not let a caller choose them either.
func (s *Sandbox) applyDefaults(resource string, record store.Record) {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return
	}

	defer func() {
		for name, value := range spec.Constants {
			record[name] = value
		}
	}()

	names := make([]string, 0, len(spec.Fields))
	for name := range spec.Fields {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		field := spec.Fields[name]

		if _, present := record[name]; present {
			continue
		}

		if field.Default != nil {
			record[name] = field.Default
			continue
		}

		// Present and null, for the providers that send the key regardless.
		// This has to come before the stamp, because a stamped value is not
		// an unset one and the whole point is that nothing has happened yet.
		if field.NullWhenUnset {
			record[name] = nil

			continue
		}

		// A field whose absence means something is left absent.
		if field.Stamped != nil && !*field.Stamped {
			continue
		}

		switch field.Type {
		case "timestamp":
			record[name] = s.clock.Unix()
		case "timestamp_ms":
			// Milliseconds. Code that divides by a thousand when it should not,
			// or fails to when it should, lands in 1970 or the year 55000, and
			// a seconds-only emulator makes both mistakes look fine.
			record[name] = s.clock.Unix() * 1000
		case "datetime":
			// RFC 3339, from the sandbox clock. GitHub, HubSpot and most newer
			// APIs send a string here, and a client parsing one does not
			// silently cope with the other.
			record[name] = s.clock.Now().UTC().Format(time.RFC3339)
		case "timestamp_ms_string":
			// Milliseconds as text, which is what ClickUp and Gmail send. It
			// looks like a number and is not one, so arithmetic on it
			// concatenates and new Date() of it is Invalid Date -- and a fake
			// sending the number makes both of those work, which is the wrong
			// direction for a difference to run in.
			record[name] = strconv.FormatInt(s.clock.Unix()*1000, 10)
		case "msdate":
			// Microsoft's date format, which looks like /Date(1552262400000+0000)/
			// and is not a date any ordinary parser accepts. Xero sends every
			// date this way and an ISO string beside it under a second name,
			// so new Date(invoice.Date) is Invalid Date and the value that
			// works is the one nobody reaches for first.
			record[name] = msdate(s.clock.Unix())
		}
	}
}

// missingRequired returns the fields a record is missing.
func (s *Sandbox) missingRequired(resource string, record store.Record) []string {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return nil
	}

	var missing []string

	names := make([]string, 0, len(spec.Fields))
	for name := range spec.Fields {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		if !spec.Fields[name].Required {
			continue
		}

		value, present := record[name]
		if !present || value == nil || value == "" {
			missing = append(missing, name)
		}
	}

	return missing
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, body any) {
	// Encoded into a buffer first, so a body that cannot be encoded becomes a
	// failure the caller can see rather than a success the caller cannot
	// parse. Writing the status line before asking the encoder meant an
	// unencodable value produced 200, application/json, and zero bytes, which
	// is the worst available answer: the client parses it before it fails, and
	// nothing anywhere says why.
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")
	// Off, because no provider here escapes for HTML. The standard library
	// turns <, > and & into <, > and & by default, on the
	// assumption the result may be pasted into a page, and that assumption is
	// wrong for every API in this collection. Discourse's fancy_title is an
	// HTML entity and Akeneo's paging links join two query parameters with an
	// ampersand; both were going out in a form their provider does not send.
	// Decoded values are identical either way, which is why no conformance
	// case could see it.
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintf(w, "{\"error\":{\"type\":\"cauldron_error\",\"message\":%s}}\n",
			mustQuote("The sandbox holds a value that cannot be encoded as JSON: "+err.Error()))

		return
	}

	// Only when nothing has said otherwise, which is the same rule the text
	// error style follows. A failure can declare its own content type, and
	// Packagist needs to: its 404 body is a bare JSON string -- valid JSON,
	// and not an object -- served as text/html. Both halves of that are the
	// finding, and forcing application/json here would have served the body
	// faithfully under a header the provider does not send.
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(status)

	_, _ = w.Write(buffer.Bytes())
}

// mustQuote renders a string as a JSON string literal.
//
// Used only by writeJSON's failure path, which has to produce valid JSON
// without going through the encoder that has just refused to produce any. A
// string always encodes, so the error is unreachable, and returning a quoted
// empty string rather than panicking keeps the failure path from having a
// failure path of its own.
func mustQuote(s string) string {
	encoded, err := json.Marshal(s)
	if err != nil {
		return `""`
	}

	return string(encoded)
}

// msdate renders Microsoft's date format, as Xero and other .NET-era APIs send
// it: the milliseconds since the epoch, wrapped, with an offset.
//
// Worth emulating exactly because it is not a date to anything that reads it.
// A client that passes it to a date parser gets Invalid Date, and one that
// stores it gets a string that looks like a timestamp and is not.
func msdate(seconds int64) string {
	return fmt.Sprintf("/Date(%d+0000)/", seconds*1000)
}
