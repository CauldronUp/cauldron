package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The failures the runtime produces itself used to carry their own internal
// name as the code and the category, so every Stripe 404 said
// type "unknown_route" where the real one says "invalid_request_error".
//
// No provider has an error type called unknown_route, and no Recipe declared
// one, so this was live on all 97 of them, on the most commonly exercised
// error path in any integration. Retry and branching logic keyed on the
// category took one path locally and another in production.
func TestABuiltinFailureBorrowsAShapeTheRecipeDeclares(t *testing.T) {
	cases := []struct {
		recipe  string
		key     string
		method  string
		path    string
		status  int
		absent  []string
		expects map[string]string
	}{
		{
			// Stripe declares both of these itself now, having been struck
			// live: a 404 for an unrecognised URL, the identical sentence and
			// status for a method a real path does not take, and no
			// error.code on either -- where every other Stripe failure has
			// one. That is the same claim from the other side: what a Recipe
			// does not say the runtime fills in with a name no provider uses,
			// so a Recipe that has looked says what it saw.
			recipe: "stripe", key: "sk_test_cauldron",
			method: http.MethodGet, path: "/v1/nope", status: http.StatusNotFound,
			absent:  []string{"error.code"},
			expects: map[string]string{"error.type": "invalid_request_error"},
		},
		{
			recipe: "stripe", key: "sk_test_cauldron",
			method: http.MethodDelete, path: "/v1/customers", status: http.StatusNotFound,
			absent:  []string{"error.code"},
			expects: map[string]string{"error.type": "invalid_request_error"},
		},
		{
			// And a Recipe that has not looked still borrows a shape it does
			// declare rather than inventing one, which is what this test was
			// written for. Front declares no unknown_route, so the failure
			// takes the wording and the title of the resource_missing it does
			// declare -- and "Not found" is a phrase Front uses, where
			// "unknown_route" is a phrase nobody does.
			recipe: "front", key: "cauldron_front_token",
			method: http.MethodGet, path: "/v1/nope", status: http.StatusNotFound,
			expects: map[string]string{
				"_error.message": "Unknown resource ID: /v1/nope",
				"_error.title":   "Not found",
			},
		},
		{
			recipe: "gitlab", key: "glpat-cauldron",
			method: http.MethodGet, path: "/api/v4/nope", status: http.StatusNotFound,
			expects: map[string]string{"message": "404 /api/v4/nope Not Found"},
		},
	}

	for _, c := range cases {
		r, err := recipe.Open(c.recipe)
		if err != nil {
			t.Fatalf("open %s: %v", c.recipe, err)
		}

		s, err := New(r, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		req := httptest.NewRequest(c.method, c.path, nil)

		if c.recipe == "gitlab" {
			req.Header.Set("PRIVATE-TOKEN", c.key)
		} else {
			req.Header.Set("Authorization", "Bearer "+c.key)
		}

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != c.status {
			t.Errorf("%s %s %s = %d, want %d\n%s", c.recipe, c.method, c.path, rec.Code, c.status, rec.Body)
		}

		body := rec.Body.String()

		// The internal name must not appear anywhere in the response. It is
		// not a word any provider uses.
		for _, internal := range []string{"unknown_route", "method_not_allowed", "invalid_request\"", "unsupported_operation"} {
			if strings.Contains(body, internal) {
				t.Errorf("%s %s: the response carries the runtime's internal name %q\n%s", c.recipe, c.path, internal, body)
			}
		}

		for _, path := range c.absent {
			if got := lookupString(decode(t, rec), path); got != "" {
				t.Errorf("%s %s: %s = %q, and the provider sends no such field\n%s", c.recipe, c.path, path, got, body)
			}
		}

		for path, want := range c.expects {
			if got := lookupString(decode(t, rec), path); got != want {
				t.Errorf("%s %s: %s = %q, want %q\n%s", c.recipe, c.path, path, got, want, body)
			}
		}
	}
}

// A Recipe declaring one of these names itself still wins, so the fallback is
// a default rather than a rule.
func TestARecipeCanDeclareABuiltinFailureItself(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	clone := *r
	clone.Errors = map[string]recipe.Error{}

	for name, defined := range r.Errors {
		clone.Errors[name] = defined
	}

	clone.Errors["unknown_route"] = recipe.Error{
		Status:  http.StatusNotFound,
		Code:    "no_such_endpoint",
		Type:    "declared_by_the_recipe",
		Message: "This Recipe says so.",
	}

	s, err := New(&clone, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/nope", nil)
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := decode(t, rec)

	if got := lookupString(body, "error.type"); got != "declared_by_the_recipe" {
		t.Errorf("error.type = %q, want the Recipe's own", got)
	}

	if got := lookupString(body, "error.code"); got != "no_such_endpoint" {
		t.Errorf("error.code = %q, want the Recipe's own", got)
	}
}

// lookupString reads a dotted path out of a decoded body.
func lookupString(body map[string]any, path string) string {
	current := any(body)

	for _, segment := range strings.Split(path, ".") {
		holder, ok := current.(map[string]any)
		if !ok {
			return ""
		}

		current = holder[segment]
	}

	text, _ := current.(string)

	return text
}
