package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Whether the credential is examined before or after the route is a fact about
// the provider, and providers disagree.
//
// Eighteen Recipes here had written a paragraph about it before this existed,
// because the runtime always checked first. Airbyte really does check first --
// it answered a byte-identical 401 to every path, real or invented. Fireworks
// really does route first: an unrouted path 404s and a wrong method 405s with
// no credential sent at all.
//
// Checking first stays the default, because it is the commoner arrangement and
// because 416 Recipes were written against it.
func TestWhetherTheCredentialIsCheckedBeforeTheRouteIsDeclarable(t *testing.T) {
	first, err := recipe.Open("airbyte")
	if err != nil {
		t.Fatalf("open airbyte: %v", err)
	}

	if first.Auth.AfterRouting {
		t.Fatalf("airbyte no longer checks the credential first; pick another Recipe")
	}

	// Credential first: a path this Recipe has never heard of is still refused
	// for its credential rather than for its path.
	assertStatus(t, first, http.MethodGet, "/v1/no-such-route-at-all", http.StatusUnauthorized)

	routed := *first
	routed.Auth.AfterRouting = true

	// Route first: the same request now gets an answer about the path, and the
	// credential never enters into it.
	assertStatus(t, &routed, http.MethodGet, "/v1/no-such-route-at-all", http.StatusNotFound)
}

// Routing first also means a wrong method is a wrong method.
//
// This is the half Fireworks demonstrates and Together contradicts: DELETE on a
// real path is a 405 naming the method on one host and a 401 on the other, and
// the difference is entirely this ordering.
func TestRoutingFirstAnswersAWrongMethodRatherThanTheCredential(t *testing.T) {
	r, err := recipe.Open("airbyte")
	if err != nil {
		t.Fatalf("open airbyte: %v", err)
	}

	var path string

	for _, route := range r.Routes {
		if route.Method == http.MethodGet {
			path = route.Path

			break
		}
	}

	if path == "" || containsBrace(path) {
		t.Skip("airbyte has no parameterless GET to send a wrong method to")
	}

	assertStatus(t, r, http.MethodDelete, path, http.StatusUnauthorized)

	routed := *r
	routed.Auth.AfterRouting = true

	assertStatus(t, &routed, http.MethodDelete, path, http.StatusMethodNotAllowed)
}

// A valid credential still has to be one, whichever order the checks run in.
//
// The risk in deferring a check is that it stops happening. This is the test
// that would catch it: a real route, a wrong credential, and a Recipe that
// routes first still answers 401 rather than serving the record.
func TestRoutingFirstStillRefusesAWrongCredentialOnARealRoute(t *testing.T) {
	r, err := recipe.Open("airbyte")
	if err != nil {
		t.Fatalf("open airbyte: %v", err)
	}

	var path string

	for _, route := range r.Routes {
		if route.Method == http.MethodGet && !containsBrace(route.Path) {
			path = route.Path

			break
		}
	}

	if path == "" {
		t.Skip("airbyte has no parameterless GET")
	}

	routed := *r
	routed.Auth.AfterRouting = true

	assertStatus(t, &routed, http.MethodGet, path, http.StatusUnauthorized)
}

func assertStatus(t *testing.T, r *recipe.Recipe, method, path string, want int) {
	t.Helper()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(method, path, nil))

	if w.Code != want {
		t.Errorf("%s %s answered %d, want %d", method, path, w.Code, want)
	}
}

func containsBrace(s string) bool {
	for i := range s {
		if s[i] == '{' {
			return true
		}
	}

	return false
}
