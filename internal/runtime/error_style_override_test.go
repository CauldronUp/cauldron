package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

// One API, three failure shapes, and Shopify's GraphQL Admin API sends all
// three. A Recipe that could describe only one of them would be claiming the
// other two do not happen.
//
//	throttled     200 {"errors": [{"message": "Throttled", ...}]}
//	bad token     401 {"errors": "[API] Invalid API key or access token"}
//	userErrors    200 {"data": {"productCreate": {"userErrors": [...]}}}
//
// The second is the one that costs an afternoon: it is the same key as the
// first, holding a bare string rather than an array. errors[0].message reads
// the sentence on a throttle and reads the character "[" on a bad token,
// because indexing a string in JavaScript succeeds and .message on that is
// undefined. Nothing throws and nothing is logged.
//
// The third has no top-level errors at all, so a client that checks the status
// and then checks errors has checked neither of the places a business refusal
// appears.
//
// Style was already overridable per error. Key and MessageField are what the
// second and third shapes need, and without them the errors[0] path and the
// data.productCreate.userErrors path were both unreachable.
func TestOneRecipeMayDescribeThreeFailureShapes(t *testing.T) {
	r, err := recipe.Open("shopifygraphql")
	if err != nil {
		t.Fatalf("open shopifygraphql: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	const path = "/admin/api/2026-01/graphql.json"

	// Shape two: the same key, holding a string. Reached by presenting a
	// credential the Recipe does not accept.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"query":"{ products(first: 1) { edges { node { id } } } }"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", "shpat_not_the_token")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (body %s)", err, rec.Body.String())
	}

	// A string, not an array. This is the whole distinction: a client that
	// indexes it gets a character rather than an object, and does not fail.
	sentence, ok := body["errors"].(string)
	if !ok {
		t.Fatalf("errors = %#v, want a bare string", body["errors"])
	}

	if !strings.HasPrefix(sentence, "[API] Invalid API key") {
		t.Errorf("errors = %q, want Shopify's sentence", sentence)
	}

	if _, present := body["message"]; present {
		t.Error("the sentence belongs under errors, not under message")
	}
}
