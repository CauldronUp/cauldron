package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A key can authenticate and still be refused.
//
// The three credential verdicts are all about what was sent: nothing, something
// unreadable, or something readable and wrong. None of them describes a key the
// provider recognises and refuses anyway -- on the wrong plan, without the
// scope, asking for a region the tier does not include.
//
// Three Recipes document that refusal and none could serve it. Apollo.io's own
// description declares a 403 for a valid out-of-scope key. Electricity Maps'
// free tier is entitled to one zone. Fitbit publishes an insufficient-scope
// message whose blanks are never filled in. Until now the only way to show any
// of them was to arm a fault, which demonstrates a shape and says nothing about
// what triggers it.
func TestAnUnentitledKeyAuthenticatesAndIsStillRefused(t *testing.T) {
	r := entitled()

	for _, c := range []struct {
		name  string
		key   string
		want  recipe.Verdict
		serve int
	}{
		{"a key the Recipe holds", "a-good-key", recipe.Accepted, http.StatusOK},
		{"a key on the wrong plan", "a-key-without-the-scope", recipe.Unentitled, http.StatusForbidden},
		{"a key nobody holds", "not-a-key-at-all", recipe.Rejected, http.StatusUnauthorized},
	} {
		s, err := New(r, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/zones", nil)
		req.Header.Set("Authorization", "Bearer "+c.key)

		if got, _ := s.credential(req, r.Auth); got != c.want {
			t.Errorf("%s: verdict %d, want %d", c.name, got, c.want)
		}

		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		if w.Code != c.serve {
			t.Errorf("%s: served %d, want %d -- %s", c.name, w.Code, c.serve, w.Body.String())
		}
	}
}

// The refusal falls back to 403 rather than 401, because the caller is not
// unidentified. A Recipe naming its own error still gets that error's status.
func TestAnUnentitledRefusalCanNameItsOwnError(t *testing.T) {
	r := entitled()
	r.Errors["out_of_scope"] = recipe.Error{Status: 402, Message: "This endpoint is not on your plan."}
	r.Auth.UnentitledError = "out_of_scope"

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/zones", nil)
	req.Header.Set("Authorization", "Bearer a-key-without-the-scope")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != 402 {
		t.Errorf("served %d, want the declared error's own 402", w.Code)
	}

	if body := w.Body.String(); body == "" {
		t.Error("the declared error served no body")
	}
}

// A Recipe that says a key both works and does not is refused before it boots.
//
// The runtime resolves it by refusing, and that is the safer resolution, but a
// file claiming both is describing a mistake rather than a provider.
func TestAKeyCannotBeBothHeldAndUnentitled(t *testing.T) {
	r := entitled()
	r.Auth.Unentitled = append(r.Auth.Unentitled, "a-good-key")

	if err := r.Validate(); err == nil {
		t.Fatal("a key listed as held and unentitled validated")
	}
}

// And naming the refusal without naming a key that reaches it is refused too,
// because nothing could ever serve it.
func TestAnUnentitledErrorNeedsAKeyThatReachesIt(t *testing.T) {
	r := entitled()
	r.Auth.Unentitled = nil
	r.Auth.UnentitledError = "out_of_scope"
	r.Errors["out_of_scope"] = recipe.Error{Status: 403, Message: "no"}

	if err := r.Validate(); err == nil {
		t.Fatal("an unreachable unentitled_error validated")
	}
}

// A Recipe with one key that works and one that authenticates and is refused.
func entitled() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "entitled",
		Capability: "infrastructure",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme:     "bearer",
			Prefix:     "Bearer ",
			Keys:       []string{"a-good-key"},
			Unentitled: []string{"a-key-without-the-scope"},
		},
		Errors: map[string]recipe.Error{},
		Resources: map[string]recipe.Resource{
			"zone": {
				ID:     recipe.ID{Style: "opaque"},
				Fields: map[string]recipe.Field{"id": {Type: "string"}},
			},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/zones", Resource: "zone", Operation: "list"},
		},
	}
}
