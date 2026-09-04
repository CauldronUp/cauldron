package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// One Recipe can hold two credentials, because one provider often has two
// surfaces.
//
// Auth was a single setting for a whole Recipe, and twenty-four Recipes wrote
// the same paragraph about it. Coinbase's two hosts want different credentials
// and no field recorded which route came from which. Healthchecks wanted to say
// "this route needs a key, and this one on a different host never does".
// Mezmo's two hosts resolve routing and the credential in opposite orders.
//
// Public already covered the never-needs-one case, and only that one.
func TestARouteCanCarryItsOwnCredential(t *testing.T) {
	r := twoSurfaces()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// The Recipe-wide key opens the Recipe-wide route and not the other one.
	if got := askKey(t, s, "/v1/things", "the-main-key"); got != http.StatusOK {
		t.Errorf("the main key answered %d on the main route, want 200", got)
	}

	if got := askKey(t, s, "/v1/reports", "the-main-key"); got != http.StatusUnauthorized {
		t.Errorf("the main key answered %d on the reporting route, want 401", got)
	}

	// And the route's own key opens the route's own surface and not the other.
	if got := askKey(t, s, "/v1/reports", "the-reporting-key"); got != http.StatusOK {
		t.Errorf("the reporting key answered %d on its own route, want 200", got)
	}

	if got := askKey(t, s, "/v1/things", "the-reporting-key"); got != http.StatusUnauthorized {
		t.Errorf("the reporting key answered %d on the main route, want 401", got)
	}
}

// A route may take a different header as well as a different key, which is
// what a second surface usually means in practice.
func TestARoutesCredentialCanArriveInItsOwnHeader(t *testing.T) {
	r := twoSurfaces()
	r.Routes[1].Auth.Scheme = "header"
	r.Routes[1].Auth.Header = "X-Report-Key"
	r.Routes[1].Auth.Prefix = ""

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/reports", nil)
	req.Header.Set("X-Report-Key", "the-reporting-key")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("the route's own header answered %d, want 200: %s", w.Code, w.Body.String())
	}

	// The Recipe-wide header is not what this route reads.
	if got := askKey(t, s, "/v1/reports", "the-reporting-key"); got != http.StatusUnauthorized {
		t.Errorf("the Recipe-wide header opened a route that declares another: %d", got)
	}
}

// A route's own ordering governs requests that reach it, and nothing else.
//
// This is the honest limit of the mechanism and it is worth stating precisely,
// because it is the half of the twenty-four that route-scoped credentials do
// not fix. Mezmo's two hosts resolve routing and the credential in opposite
// orders, and ordering is mostly a claim about what happens when routing
// FAILS -- an unrouted path, a wrong method. A request that matches no route
// has no route to take an ordering from, so the Recipe's own applies and
// Mezmo's split stays in prose.
//
// What a route's after_routing does govern is a request that matched it.
func TestARoutesOwnOrderingGovernsOnlyRequestsThatReachIt(t *testing.T) {
	r := twoSurfaces()
	r.Routes[1].Auth.AfterRouting = true

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// A wrong method matches no route, so the Recipe's ordering decides and
	// its credential is judged first.
	req := httptest.NewRequest(http.MethodDelete, "/v1/reports", nil)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("a wrong method answered %d; the Recipe's own ordering decides there, not the route's", w.Code)
	}

	// With the Recipe-wide credential accepted, routing is reached and the
	// method rejection is real -- which is what shows the request got past
	// the gate rather than never reaching it.
	req = httptest.NewRequest(http.MethodDelete, "/v1/reports", nil)
	req.Header.Set("Authorization", "Bearer the-main-key")

	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("with the Recipe's key accepted a wrong method answered %d, want 405", w.Code)
	}
}

// A path that matches no route has no route to take a credential from, so the
// Recipe's own applies. That is the only answer available and the right one: a
// path that does not exist belongs to no surface.
func TestAnUnroutedPathIsJudgedByTheRecipesOwnCredential(t *testing.T) {
	s, err := New(twoSurfaces(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got := askKey(t, s, "/v1/nothing-here", "the-main-key"); got != http.StatusNotFound {
		t.Errorf("an unrouted path with the main key answered %d, want 404", got)
	}

	if got := askKey(t, s, "/v1/nothing-here", "the-reporting-key"); got != http.StatusUnauthorized {
		t.Errorf("an unrouted path with a route's key answered %d, want the Recipe's own 401", got)
	}
}

// Naming a credential on a route that is also public is refused, because the
// exemption wins and the credential would be dead weight a reader takes for a
// claim.
func TestARouteCannotBeBothPublicAndCredentialled(t *testing.T) {
	r := twoSurfaces()
	r.Routes[1].Public = recipe.PublicMode{Always: true}

	if err := r.Validate(); err == nil {
		t.Fatal("a route that is public and names a scheme validated")
	}
}

func askKey(t *testing.T, s *Sandbox, path, key string) int {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", "Bearer "+key)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	return w.Code
}

// A Recipe with two surfaces behind one name, which is what most of the
// twenty-four actually have.
func twoSurfaces() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "twosurfaces",
		Capability: "observability",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme: "bearer",
			Prefix: "Bearer ",
			Keys:   []string{"the-main-key"},
		},
		Resources: map[string]recipe.Resource{
			"thing":  {ID: recipe.ID{Style: "opaque"}, Fields: map[string]recipe.Field{"id": {Type: "string"}}},
			"report": {ID: recipe.ID{Style: "opaque"}, Fields: map[string]recipe.Field{"id": {Type: "string"}}},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
			{
				Method: "GET", Path: "/v1/reports", Resource: "report", Operation: "list",
				Auth: &recipe.Auth{
					Scheme: "bearer",
					Prefix: "Bearer ",
					Keys:   []string{"the-reporting-key"},
				},
			},
		},
	}
}
