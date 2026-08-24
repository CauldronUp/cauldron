package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Absent and null are different on the wire, and for a paging loop the
// difference decides whether it terminates.
//
// Metronome's customer listing declares next_page required and nullable, and
// its own example sends "next_page": null on the last page. Every listing here
// omitted the cursor field when there was nothing to point at, so a loop
// written as `while ('next_page' in body)` ran for ever against the provider
// and stopped immediately against the emulator -- which is the failure this
// project exists to prevent, arrived at from the paging side.
//
// The key has to be present and its value has to be null. Asserting only that
// the value is falsy would pass against the omission this replaces, so the
// test reads the decoded map rather than the typed value.
func TestALastPageCanCarryANullCursorRatherThanNoCursor(t *testing.T) {
	r, err := recipe.Open("metronome")
	if err != nil {
		t.Fatalf("open metronome: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	req.Header.Set("Authorization", "Bearer cauldron-metronome-fixture-token")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	cursor, present := body["next_page"]
	if !present {
		t.Fatal("next_page is absent, and Metronome declares it required and sends null")
	}

	if cursor != nil {
		t.Errorf("next_page = %v, want null on the last page", cursor)
	}
}
