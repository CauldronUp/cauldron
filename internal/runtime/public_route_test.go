package runtime

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A public route needs no credential on a Recipe whose other routes do.
//
// Auth is one setting for a whole Recipe, and nineteen Recipes have had to
// write down that their provider disagrees. Cognito serves a genuinely public
// key set beside SigV4-gated operations. Checkly's runtimes list answers 200
// with nothing presented. Ashby's job board is public and the rest of Ashby is
// not. Every one of them had two choices and both were wrong: declare the route
// and teach a caller to send a credential nowhere wanted, or leave it out and
// describe a provider as smaller than it is.
func TestAPublicRouteIsServedWithNoCredential(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{Always: true}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// The gated route still refuses, which is the half that keeps the
	// exemption from being a hole.
	if status, _ := askFor(t, s, "/v1/incidents"); status != http.StatusUnauthorized {
		t.Errorf("the gated route answered %d with no credential, want 401", status)
	}

	if status, body := askFor(t, s, "/v1/runtimes"); status != http.StatusOK {
		t.Errorf("the public route answered %d with no credential, want 200: %s", status, body)
	}
}

// Without the marking the same route is refused, which is what all nineteen of
// those Recipes have been living with.
func TestAnUnmarkedRouteIsStillGated(t *testing.T) {
	s, err := New(gatedAndNot(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status, _ := askFor(t, s, "/v1/runtimes"); status != http.StatusUnauthorized {
		t.Errorf("an unmarked route answered %d, want 401", status)
	}
}

// A path that is not a route is unaffected. The exemption is per-route, so
// there is nothing for an unrouted request to be exempt by, and it must still
// be refused before it is identified as missing -- which is the ordering this
// Recipe declares.
func TestAnUnroutedPathIsRefusedEvenBesideAPublicRoute(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{Always: true}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status, _ := askFor(t, s, "/v1/nothing-here"); status != http.StatusUnauthorized {
		t.Errorf("an unrouted path answered %d, want the 401 this ordering gives", status)
	}
}

// And the same holds for a Recipe that routes first: the public route is
// served, and the unrouted path gets the 404 that ordering promises rather than
// borrowing the exemption.
func TestAPublicRouteWorksOnARoutingFirstRecipe(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{Always: true}
	r.Auth.AfterRouting = true

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status, body := askFor(t, s, "/v1/runtimes"); status != http.StatusOK {
		t.Errorf("the public route answered %d, want 200: %s", status, body)
	}

	if status, _ := askFor(t, s, "/v1/incidents"); status != http.StatusUnauthorized {
		t.Errorf("the gated route answered %d, want 401", status)
	}

	if status, _ := askFor(t, s, "/v1/nothing-here"); status != http.StatusNotFound {
		t.Errorf("an unrouted path answered %d, want the 404 routing-first gives", status)
	}
}

// Public exempts the credential and nothing else.
//
// A provider can want its version header on a route that needs no key, and
// Cognito does. Folding the two together would have been the easy reading and
// the wrong one, because a caller would learn to omit a header the provider
// insists on.
func TestARequiredHeaderIsStillRequiredOnAPublicRoute(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{Always: true}
	r.Auth.AfterRouting = true
	r.RequiredHeaders = map[string]recipe.RequiredHeader{"X-Api-Version": {Error: "missing_version"}}
	r.Errors = map[string]recipe.Error{"missing_version": {Status: 400, Message: "The X-Api-Version header is required."}}
	// The validator insists a declared required header be shown enforced by a
	// case, which is exactly the claim this test is here to make.
	r.Conformance = []recipe.Case{{
		Name:    "the version header is required on the public route too",
		Source:  "constructed for this test; no provider is being described",
		Request: recipe.Request{Method: "GET", Path: "/v1/runtimes"},
		Expect:  recipe.Expectation{Status: 400},
	}}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status, _ := askFor(t, s, "/v1/runtimes"); status != http.StatusBadRequest {
		t.Errorf("a public route waved through a missing required header: %d, want 400", status)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/runtimes", nil)
	req.Header.Set("X-Api-Version", "1")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("the public route answered %d with the header present, want 200", w.Code)
	}
}

// A Recipe with a gated route and a public one beside it.
func gatedAndNot() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "gatedandnot",
		Capability: "observability",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "2026-09-02"},
		Auth: recipe.Auth{
			Scheme: "bearer",
			Keys:   []string{"a-key-the-recipe-holds"},
		},
		Resources: map[string]recipe.Resource{
			"incident": {
				ID:     recipe.ID{Prefix: "inc_"},
				Fields: map[string]recipe.Field{"id": {Type: "string"}},
			},
			"runtime": {
				ID:     recipe.ID{Prefix: "rt_"},
				Fields: map[string]recipe.Field{"id": {Type: "string"}},
			},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/incidents", Resource: "incident", Operation: "list"},
			{Method: "GET", Path: "/v1/runtimes", Resource: "runtime", Operation: "list"},
		},
	}
}

func askFor(t *testing.T, s *Sandbox, path string) (int, string) {
	t.Helper()

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

	body, _ := io.ReadAll(w.Body)

	return w.Code, strings.TrimSpace(string(body))
}
