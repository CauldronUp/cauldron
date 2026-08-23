package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Two kinds of declared field on a list-shaped failure go to two different
// places, because they describe two different things.
//
// Fields on the error envelope belong to the response. Clerk sends a trace id
// that identifies the request, not any one of the failures inside it, so it
// sits beside the array.
//
// Fields on a named error belong to that failure. eBay sends
// {"errors":[{errorId, domain, category, message, longMessage}]} and a client
// reads domain and longMessage off the entry it is looping over, because one
// request can fail several ways at once and each one has its own domain.
//
// Both used to land beside the array, which made this branch disagree with the
// bare-array branch beside it: the same Recipe changing key to "-" moved its
// fields from the response into the entry. Nothing shipped depended on it, and
// a Recipe could not describe eBay's shape at all -- errors[0].longMessage was
// not expressible, so the field a client actually shows a user had nowhere to
// go.
func TestANamedErrorsFieldsGoInsideTheEntryAndTheEnvelopesBeside(t *testing.T) {
	r, err := recipe.Open("ebay")
	if err != nil {
		t.Fatalf("open ebay: %v", err)
	}

	clone := *r

	// A copy, so the shared embedded Recipe is untouched.
	responses := clone.Responses
	responses.Error.Fields = map[string]any{"traceId": "cauldron-trace-0001"}
	clone.Responses = responses

	s, err := New(&clone, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sell/fulfillment/v1/order", nil)
	req.Header.Set("Authorization", "Bearer not_the_token")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (body %s)", err, rec.Body.String())
	}

	// The envelope's field, beside the array.
	if body["traceId"] != "cauldron-trace-0001" {
		t.Errorf("traceId = %v, want it beside the array", body["traceId"])
	}

	entries, ok := body["errors"].([]any)
	if !ok || len(entries) == 0 {
		t.Fatalf("errors = %v, want a non-empty array", body["errors"])
	}

	entry, ok := entries[0].(map[string]any)
	if !ok {
		t.Fatalf("errors[0] = %v, want an object", entries[0])
	}

	// The named error's fields, inside the entry. These are the ones a client
	// reads off the failure it is looping over.
	for field, want := range map[string]string{
		"domain":      "OAuth",
		"category":    "REQUEST",
		"longMessage": "Invalid access token. Check the value of the Authorization header.",
		"message":     "Invalid access token",
	} {
		if entry[field] != want {
			t.Errorf("errors[0].%s = %v, want %q", field, entry[field], want)
		}
	}

	// And the two do not leak into each other. A client reading
	// response.longMessage finds nothing, which is what eBay does.
	for _, field := range []string{"domain", "category", "longMessage", "message"} {
		if _, present := body[field]; present {
			t.Errorf("%s is beside the array and belongs inside the entry", field)
		}
	}

	if _, present := entry["traceId"]; present {
		t.Error("traceId is inside the entry and belongs beside the array")
	}
}
