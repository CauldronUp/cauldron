package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Optimistic concurrency is the contract a test suite is structurally unable
// to test, which is exactly why an emulator has to enforce it.
//
// commercetools will not take a write unless the caller quotes the version it
// expects to replace. A suite is the one place where nothing else is writing
// to the record, so code that quotes a stale version -- or none at all --
// passes every test that exists. In production the same code does not error:
// the later write wins, the earlier one is gone, and nothing is logged
// anywhere.
//
// Three outcomes, because commercetools distinguishes three. The current
// version is accepted and moves the number on. A different one is a 409 that
// hands back the current number so the retry can be scripted. None at all is
// a 400 about a required field, which is a different thing entirely: a client
// that retries on the first should not retry on the second.
func TestAWriteMayHaveToQuoteTheVersionItReplaces(t *testing.T) {
	r, err := recipe.Open("commercetools")
	if err != nil {
		t.Fatalf("open commercetools: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	const product = "/cauldron-shop/products/e7ba4c75-b1bb-483d-94d8-2c4a10f78472"

	write := func(body string) (int, map[string]any) {
		t.Helper()

		req := httptest.NewRequest(http.MethodPost, product, strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer cauldron_commercetools_fixture_token_0")
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		var decoded map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
			t.Fatalf("%s: decode: %v", body, err)
		}

		return rec.Code, decoded
	}

	// The fixture is at version 2, so a write quoting 2 is the one the
	// provider is waiting for.
	status, body := write(`{"version":2,"actions":[{"action":"publish"}]}`)
	if status != http.StatusOK {
		t.Fatalf("a write quoting the current version = %d, want 200: %v", status, body)
	}

	if got, _ := body["version"].(float64); got != 3 {
		t.Errorf("the version after a write = %v, want 3", got)
	}

	// And replaying exactly that body is now stale, which is the whole reason
	// the number exists. A lock that let the same write through twice would
	// be decoration.
	status, body = write(`{"version":2,"actions":[{"action":"publish"}]}`)
	if status != http.StatusConflict {
		t.Errorf("replaying a write = %d, want 409", status)
	}

	errors, _ := body["errors"].([]any)
	if len(errors) != 1 {
		t.Fatalf("errors = %#v, want one entry", body["errors"])
	}

	if entry, _ := errors[0].(map[string]any); entry["code"] != "ConcurrentModification" {
		t.Errorf("code = %#v, want ConcurrentModification", entry["code"])
	}

	// A write carrying no version at all is a different refusal, with a
	// different status, because retrying it would never succeed.
	status, body = write(`{"actions":[{"action":"publish"}]}`)
	if status != http.StatusBadRequest {
		t.Errorf("a write with no version = %d, want 400", status)
	}

	errors, _ = body["errors"].([]any)
	if len(errors) != 1 {
		t.Fatalf("errors = %#v, want one entry", body["errors"])
	}

	if entry, _ := errors[0].(map[string]any); entry["code"] != "RequiredField" {
		t.Errorf("code = %#v, want RequiredField", entry["code"])
	}

	// The record moved once and only once. Two refused writes left it where
	// the accepted one put it.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, product, nil)
	req.Header.Set("Authorization", "Bearer cauldron_commercetools_fixture_token_0")
	s.ServeHTTP(rec, req)

	var current map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &current); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got, _ := current["version"].(float64); got != 3 {
		t.Errorf("the version after two refusals = %v, want 3", got)
	}
}

// A resource that declares no version field is unaffected, which is every
// Recipe written before this one. Printify's orders are read-only here, so
// Shopware's products stand in: a write to a resource with no lock is not
// asked for a number it does not have.
func TestAResourceWithNoVersionFieldIsUnaffected(t *testing.T) {
	r, err := recipe.Open("shopware")
	if err != nil {
		t.Fatalf("open shopware: %v", err)
	}

	if resource, ok := r.Resources["product"]; !ok || resource.VersionField != "" {
		t.Fatalf("shopware's product declares a version field, so it is the wrong control")
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// The listing answers as it always did, with no version demanded of
	// anybody.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/store-api/product?limit=1", nil)
	req.Header.Set("sw-access-key", "SWSCCAULDRONFIXTUREKEY00")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
