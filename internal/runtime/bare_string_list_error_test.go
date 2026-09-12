// A failure answered as a bare array of one string.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The string_list style always wrapped its array in a key. Datadog sends
// {"errors": ["..."]} and that was the only shape it could make, so a provider
// answering the array *as* the body had nowhere to declare it -- the closest a
// Recipe could get was an envelope the provider does not send, which is wrong
// in exactly the direction the finding is about.
//
// Storyblok is the provider that made the unwrapped form real. An unrouted
// path answers ["This record could not be found"] while a bad credential
// answers {"error":"Unauthorized"}, so one API answers an object for one
// failure and a bare array of strings for the other, and client code reading
// body.error finds the sentence on one and undefined on the other.
//
// The list style already spelled this with key: "-". This gives string_list
// the same spelling rather than a second one.
//
// A shipped Recipe on purpose: the shape is Storyblok's, not a fixture's.
func TestABareArrayOfStringsCanBeAFailureBody(t *testing.T) {
	r, err := recipe.Open("storyblok")
	if err != nil {
		t.Fatalf("open storyblok: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("space"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v2/cdn/cauldron-nope?token=notarealtoken", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404: %s", rec.Code, rec.Body.String())
	}

	var body any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	list, ok := body.([]any)
	if !ok {
		t.Fatalf("the failure is not a bare list, so key: \"-\" did nothing: %.200s", rec.Body.String())
	}

	if len(list) != 1 {
		t.Fatalf("want one sentence in the list, got %d", len(list))
	}

	if got, ok := list[0].(string); !ok || got != "This record could not be found" {
		t.Errorf("entry = %#v, want the sentence as a bare string", list[0])
	}

	// And the other failure on the same Recipe is still an object, because
	// that is the pair that makes this worth modelling at all.
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v2/cdn/stories", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401: %s", rec.Code, rec.Body.String())
	}

	var refused map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &refused); err != nil {
		t.Fatalf("the 401 is not an object, so the two failures no longer differ: %v", err)
	}

	if refused["error"] != "Unauthorized" {
		t.Errorf("error = %#v, want Unauthorized", refused["error"])
	}
}
