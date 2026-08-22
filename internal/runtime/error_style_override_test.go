package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A provider can answer two shapes for the same status, and the npm registry
// does. Checked against registry.npmjs.org on 2026-08-22:
//
//	GET /cauldron-no-such-package-xyz  ->  {"error":"Not found"}
//	GET /left-pad/99.99.99             ->  "version not found: 99.99.99"
//
// Same registry, same 404, one object and one bare JSON string. Code reading
// body.error off the second finds undefined rather than failing, so a client
// reporting body.error as the reason reports "undefined" -- the quieter and
// worse of the two failures.
//
// The envelope was Recipe-wide, so a Recipe could describe one of those and
// not both. A single error may now override it.
func TestOneErrorMayOverrideTheEnvelope(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	clone := *r

	// A copy of the map, so the shared embedded Recipe is untouched.
	errors := make(map[string]recipe.Error, len(clone.Errors))
	for name, e := range clone.Errors {
		errors[name] = e
	}

	missing := errors["resource_missing"]
	missing.Style = "string"
	missing.Message = "version not found: 99.99.99"
	errors["resource_missing"] = missing
	clone.Errors = errors

	s, err := New(&clone, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/customers/cus_nosuchrecord", nil)
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	// Valid JSON, and not an object. That is the whole distinction: a client
	// calling .json() succeeds and is handed a string.
	var body any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("a bare string body must still be JSON: %v\n%s", err, rec.Body)
	}

	text, ok := body.(string)
	if !ok {
		t.Fatalf("body is %T, want a bare string: %s", body, rec.Body)
	}

	if text != "version not found: 99.99.99" {
		t.Errorf("body = %q", text)
	}
}

// The override is for one failure, not the Recipe. Every other error in the
// same Recipe keeps the shape it had, or a Recipe could not describe a
// provider that answers both.
func TestTheOtherErrorsKeepTheirShape(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	clone := *r

	errors := make(map[string]recipe.Error, len(clone.Errors))
	for name, e := range clone.Errors {
		errors[name] = e
	}

	missing := errors["resource_missing"]
	missing.Style = "string"
	errors["resource_missing"] = missing
	clone.Errors = errors

	s, err := New(&clone, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// An unauthenticated request raises a different error, which must still
	// arrive in Stripe's nested object.
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/customers", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("another error should still be an object: %v\n%s", err, rec.Body)
	}

	if _, ok := body["error"]; !ok {
		t.Errorf("the override leaked onto another failure: %s", rec.Body)
	}
}
