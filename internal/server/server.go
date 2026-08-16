// Package server hosts many sandboxes in one process and exposes a small
// control API for driving them.
//
// One process rather than a container per provider: booting eight fakes should
// cost milliseconds and a few megabytes, because in CI it happens on every
// commit.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/CauldronUp/cauldron/internal/clock"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/runtime"
	"github.com/CauldronUp/cauldron/internal/store"
)

// controlPrefix namespaces the control API away from emulated provider routes.
// A leading underscore keeps it clear of anything a real provider would serve.
const controlPrefix = "/_cauldron"

// Server multiplexes sandboxes.
type Server struct {
	mu        sync.RWMutex
	sandboxes map[string]*runtime.Sandbox
	order     []string
	clock     *clock.Clock
}

// New returns an empty server. All sandboxes share one clock, so advancing
// time moves every provider together — which is what a developer means by
// "skip forward a month".
func New() *Server {
	return &Server{
		sandboxes: map[string]*runtime.Sandbox{},
		clock:     clock.New(),
	}
}

// Clock returns the shared clock.
func (s *Server) Clock() *clock.Clock { return s.clock }

// Mount boots a Recipe by name and adds it to the server.
func (s *Server) Mount(name string, seed int64, fixture string) error {
	r, err := recipe.Open(name)
	if err != nil {
		return err
	}

	sandbox, err := runtime.New(r, runtime.Options{Seed: seed, Clock: s.clock})
	if err != nil {
		return err
	}

	if fixture != "" {
		if err := sandbox.Seed(fixture); err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sandboxes[r.Name]; !exists {
		s.order = append(s.order, r.Name)
		sort.Strings(s.order)
	}

	s.sandboxes[r.Name] = sandbox

	return nil
}

// Names returns the mounted recipe names.
func (s *Server) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]string, len(s.order))
	copy(out, s.order)

	return out
}

// Sandbox returns a mounted sandbox.
func (s *Server) Sandbox(name string) (*runtime.Sandbox, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sandbox, ok := s.sandboxes[name]

	return sandbox, ok
}

// ServeHTTP routes a request to the right sandbox.
//
// Two addressing schemes are supported. The path prefix
// (http://localhost:4600/stripe/v1/customers) works everywhere with no DNS
// setup, which matters because it is what a developer pastes into a base URL.
// The host prefix (http://stripe.cauldron.test/v1/customers) is nicer once
// local domains exist.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, controlPrefix) {
		s.control(w, r)
		return
	}

	if name, ok := s.fromHost(r.Host); ok {
		s.dispatch(w, r, name, r.URL.Path)
		return
	}

	name, rest := splitFirstSegment(r.URL.Path)

	if name == "" {
		s.writeError(w, http.StatusNotFound, "No recipe addressed. Use /<recipe>/... or set a Host header like <recipe>.cauldron.test.")
		return
	}

	s.dispatch(w, r, name, rest)
}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request, name, path string) {
	sandbox, ok := s.Sandbox(name)
	if !ok {
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("No recipe %q is running. Mounted: %s.", name, strings.Join(s.Names(), ", ")))
		return
	}

	// Rewrite the path so the sandbox sees the provider's own URL space.
	clone := r.Clone(r.Context())
	clone.URL.Path = path

	if clone.URL.Path == "" {
		clone.URL.Path = "/"
	}

	sandbox.ServeHTTP(w, clone)
}

// fromHost matches a Host header's first label against a mounted recipe.
func (s *Server) fromHost(host string) (string, bool) {
	if host == "" {
		return "", false
	}

	if colon := strings.Index(host, ":"); colon >= 0 {
		host = host[:colon]
	}

	label, _, found := strings.Cut(host, ".")
	if !found {
		return "", false
	}

	if _, ok := s.Sandbox(label); ok {
		return label, true
	}

	return "", false
}

func splitFirstSegment(path string) (string, string) {
	trimmed := strings.TrimPrefix(path, "/")

	first, rest, found := strings.Cut(trimmed, "/")
	if !found {
		return first, "/"
	}

	return first, "/" + rest
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"type":    "cauldron_error",
			"message": message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	_ = encoder.Encode(body)
}

// recordFrom decodes a JSON body into a record, tolerating an empty body.
func recordFrom(r *http.Request) (store.Record, error) {
	if r.Body == nil {
		return store.Record{}, nil
	}

	var out store.Record

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&out); err != nil {
		if err.Error() == "EOF" {
			return store.Record{}, nil
		}

		return nil, err
	}

	return out, nil
}

// Snapshot is the state of every sandbox on this server, plus the shared
// clock.
//
// It is taken across all recipes rather than one, because a bug worth
// capturing usually spans providers: the payment succeeded, the order did not
// ship, and the state that proves it lives in two sandboxes at once.
type Snapshot struct {
	Version int                         `json:"version"`
	Clock   time.Time                   `json:"clock"`
	Recipes map[string]runtime.Snapshot `json:"recipes"`
}

// Snapshot captures every mounted sandbox.
func (s *Server) Snapshot() Snapshot {
	out := Snapshot{
		Version: runtime.SnapshotVersion,
		Clock:   s.clock.Now(),
		Recipes: map[string]runtime.Snapshot{},
	}

	for _, name := range s.Names() {
		if sandbox, ok := s.Sandbox(name); ok {
			out.Recipes[name] = sandbox.Snapshot()
		}
	}

	return out
}

// Restore applies a snapshot to every recipe it names.
//
// Recipes in the snapshot that are not mounted are reported rather than
// ignored, because a partially restored sandbox that claims success is worse
// than a refusal.
func (s *Server) Restore(snapshot Snapshot) error {
	var missing []string

	for name := range snapshot.Recipes {
		if _, ok := s.Sandbox(name); !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)

		return fmt.Errorf("this snapshot needs recipes that are not running: %s", strings.Join(missing, ", "))
	}

	for name, state := range snapshot.Recipes {
		sandbox, _ := s.Sandbox(name)

		if err := sandbox.Restore(state); err != nil {
			return fmt.Errorf("restoring %s: %w", name, err)
		}
	}

	if !snapshot.Clock.IsZero() {
		s.clock.Set(snapshot.Clock)
	}

	return nil
}

// HasFixture reports whether a mounted recipe ships a fixture.
func (s *Server) HasFixture(recipe, fixture string) bool {
	sandbox, ok := s.Sandbox(recipe)
	if !ok {
		return false
	}

	for _, candidate := range sandbox.Fixtures() {
		if candidate == fixture {
			return true
		}
	}

	return false
}

// SeedRecipe loads a fixture into one mounted recipe.
func (s *Server) SeedRecipe(recipe, fixture string) error {
	sandbox, ok := s.Sandbox(recipe)
	if !ok {
		return fmt.Errorf("no recipe %q is running", recipe)
	}

	return sandbox.Seed(fixture)
}
