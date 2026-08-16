// Package runtime serves a Recipe over HTTP: routing, auth, state, pagination,
// failure injection and webhooks.
//
// One process hosts many sandboxes, routed by name, rather than a container per
// provider. Booting eight fakes should cost milliseconds and a few megabytes,
// because in CI it happens on every commit.
package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

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
	recipe *recipe.Recipe
	store  *store.Store
	clock  *clock.Clock
	router *router
	faults *faultSet
	log    *requestLog

	mu       sync.RWMutex
	fixture  string
	webhooks *webhookQueue
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
		log:      newRequestLog(opts.LogSize),
		webhooks: newWebhookQueue(r, c),
	}

	for name, resource := range r.Resources {
		s.store.DeclareStyle(name, resource.ID.Style, resource.ID.Prefix, resource.ID.Length)
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
// timestamps from the sandbox clock rather than the wall clock.
func (s *Sandbox) applyDefaults(resource string, record store.Record) {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return
	}

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

		if field.Type == "timestamp" {
			record[name] = s.clock.Unix()
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	_ = encoder.Encode(body)
}
