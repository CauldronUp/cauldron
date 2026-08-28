package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route's list override is read in full.
//
// The merge carried eighteen of the envelope's fields and silently dropped the
// rest, so a Recipe declaring collapse_single, cursor_url, final_field,
// complete_field or fields on one route had the line read, validated and thrown
// away. No conformance case could tell the difference: the emulator behaved
// exactly as though the line were not there, which is what makes this class of
// defect expensive rather than merely wrong.
//
// This is the fifth declared-and-ignored field closed here, after a list key
// the collection name already supplied, auth.header on the bearer scheme, a
// declared Content-Type overwritten by a default, and an empty object
// expectation that asserted nothing.
//
// Open Library found it. /api/books answers one book as the object under the
// caller's own query string and no books as {} -- collapse_single and
// omit_when_empty on one route -- and only the second reached the runtime.
func TestARoutesListOverrideIsReadInFull(t *testing.T) {
	r, err := recipe.Open("openlibrary")
	if err != nil {
		t.Fatalf("open openlibrary: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("tomsawyer"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet,
		"/api/books?bibkeys=ISBN:0451526538&format=json&jscmd=data", nil)
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	entry, ok := body["ISBN:0451526538"]
	if !ok {
		t.Fatalf("no entry under the caller's own key: %v", body)
	}

	// The whole point of collapse_single: one book is the object, not a list
	// holding it.
	if _, isList := entry.([]any); isList {
		t.Errorf("collapse_single was declared on the route and ignored: the entry is a list")
	}

	object, ok := entry.(map[string]any)
	if !ok {
		t.Fatalf("the entry is neither a list nor an object: %T", entry)
	}

	if got := object["key"]; got != "/books/OL1017798M" {
		t.Errorf("key = %v, want /books/OL1017798M", got)
	}
}
