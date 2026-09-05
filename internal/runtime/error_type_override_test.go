package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// One failure can put its category somewhere else, the way one failure can
// already move its message and its code.
//
// The envelope names three fields and two of them were overridable per error.
// Clio is what found the third missing: it answers one shape to a request that
// sent no credential and another to one that sent a wrong credential, and the
// second carries a property literally named "type" that the first does not.
// Its Recipe reached that by disabling type_field Recipe-wide and repurposing
// code_field on the one error -- the same bytes on the wire, produced by a knob
// named for something else, with a comment explaining the substitution to
// whoever reads it next.
func TestOneFailureCanMoveItsCategoryField(t *testing.T) {
	r := categorised()
	r.Errors["moved"] = recipe.Error{
		Status:    403,
		Type:      "forbidden",
		Message:   "Not for you.",
		TypeField: "detail.kind",
	}
	r.Routes[0] = recipe.Route{Method: "GET", Path: "/v1/things", Error: "moved"}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	body := categoryBody(t, s)

	if got := lookupString(body, "detail.kind"); got != "forbidden" {
		t.Errorf("the moved category read %q, want forbidden", got)
	}

	if got := lookupString(body, "type"); got != "" {
		t.Errorf("the category was also left at the Recipe-wide name: %q", got)
	}
}

// And "-" removes it for one failure, which is the other half of the same
// asymmetry: a provider that names its category everywhere except here.
func TestOneFailureCanDropItsCategoryField(t *testing.T) {
	r := categorised()
	r.Errors["bare"] = recipe.Error{
		Status:    404,
		Type:      "not_found",
		Message:   "No such thing.",
		TypeField: "-",
	}
	r.Routes[0] = recipe.Route{Method: "GET", Path: "/v1/things", Error: "bare"}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	body := categoryBody(t, s)

	if got := lookupString(body, "type"); got != "" {
		t.Errorf("the category survived a failure that removes it: %q", got)
	}

	if got := lookupString(body, "message"); got != "No such thing." {
		t.Errorf("the message went missing with it: %q", got)
	}
}

// A failure that says nothing about it keeps the Recipe's own name, which is
// every failure in every Recipe shipping today.
func TestAFailureThatSaysNothingKeepsTheRecipesCategoryField(t *testing.T) {
	r := categorised()
	r.Errors["ordinary"] = recipe.Error{Status: 400, Type: "bad_request", Message: "No."}
	r.Routes[0] = recipe.Route{Method: "GET", Path: "/v1/things", Error: "ordinary"}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got := lookupString(categoryBody(t, s), "type"); got != "bad_request" {
		t.Errorf("the Recipe-wide category read %q, want bad_request", got)
	}
}

func categoryBody(t *testing.T, s *Sandbox) map[string]any {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer the-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("not JSON: %s", rec.Body)
	}

	return decode(t, rec)
}

func categorised() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "categorised",
		Capability: "identity",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme: "bearer",
			Prefix: "Bearer ",
			Keys:   []string{"the-key"},
		},
		Responses: recipe.Responses{
			Error: recipe.ErrorResponse{
				Style:        "flat",
				MessageField: "message",
				CodeField:    "-",
				TypeField:    "type",
			},
		},
		Errors: map[string]recipe.Error{},
		Resources: map[string]recipe.Resource{
			"thing": {ID: recipe.ID{Style: "opaque"}, Fields: map[string]recipe.Field{"id": {Type: "string"}}},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
		// The envelope's own field names are claims, and the validator asks
		// for one case asserting each where the value exists.
		Conformance: []recipe.Case{{
			Name:    "the failure carries a message under the name the envelope declares",
			Source:  "https://example.invalid/errors",
			Request: recipe.Request{Method: "GET", Path: "/v1/nothing-here"},
			Expect: recipe.Expectation{
				Status: 404,
				Body:   map[string]any{"message": "Unrecognised request URL."},
			},
		}},
	}
}
