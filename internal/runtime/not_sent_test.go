package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route that scopes by a path segment needs that segment as a field, because
// that is how the record is partitioned. Nothing then stopped it being
// emitted, and most providers do not repeat a partition they already put in
// the URL: Fly does not send app_name on a machine, Tradier does not say which
// account an order is in.
//
// The only way to say so was a route's returns naming every other field --
// twenty-three names to hide one on Fly, repeated per route. `in: "-"` says it
// once, on the field, where the fact belongs.
//
// An audit found 115 scope fields across 37 Recipes reaching the wire with no
// conformance case mentioning them. Some of those providers really do echo the
// partition, so this is a tool for the ones that have been read rather than a
// licence to change them all.

func fetch(t *testing.T, name, fixture, path, key string) map[string]any {
	t.Helper()

	r, err := recipe.Open(name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed(fixture); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", "Bearer "+key)

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("not an object: %v %s", err, rec.Body.String())
	}

	return body
}

func TestAFieldMarkedNotSentDoesNotReachTheWire(t *testing.T) {
	body := fetch(t, "fly", "one-app",
		"/v1/apps/cauldron-web/machines/148e2d1d7d1289", "FlyV1cauldron000000000000000000")

	if _, present := body["app_name"]; present {
		t.Errorf("the partition reached the wire: %v", body)
	}

	// And not under a key literally spelled "-" either. Without this the test
	// passes whether the field is withheld or nested under a dash, which is
	// how the first three mutations all survived.
	if _, present := body["-"]; present {
		t.Errorf(`the partition was nested under "-" instead of withheld: %v`, body)
	}

	// And the record is otherwise intact, so the field was withheld rather
	// than the resource being trimmed.
	if body["name"] != "winter-frost-4821" {
		t.Errorf("the rest of the machine is missing: %v", body)
	}
}

func TestTheRecordStillFindsItselfByThatField(t *testing.T) {
	// Withholding it from the wire must not stop it partitioning: asking the
	// wrong app for this machine still has to be a miss.
	r, err := recipe.Open("fly")
	if err != nil {
		t.Fatalf("open fly: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-app"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/apps/somebody-else/machines/148e2d1d7d1289", nil)
	req.Header.Set("Authorization", "Bearer FlyV1cauldron000000000000000000")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("another app's machine was served: %d %s", rec.Code, rec.Body.String())
	}
}

func TestAnOrdinaryNestedFieldStillNests(t *testing.T) {
	// "-" must not be read as a nesting path, or every in: field would break.
	body := fetch(t, "tradier", "busy-account",
		"/v1/accounts/VA000001/orders/228175", "cauldron000000000000000tradier")

	if _, present := body["account_id"]; present {
		t.Errorf("the partition reached the wire: %v", body)
	}

	if _, present := body["-"]; present {
		t.Errorf(`the partition was nested under "-": %v`, body)
	}

	if body["symbol"] != "AAPL" {
		t.Errorf("the order is missing its own fields: %v", body)
	}
}
