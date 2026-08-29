// A single record answered as a bare array of one.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// envelope.array wraps one record in a list, and it was read only on the
// wrapped path. A route declaring it with a bare style got the object back and
// nothing said otherwise -- the Recipe said list, the runtime sent an object,
// and no conformance case existed that could tell them apart, because the only
// Recipes using it were Xero and Ghost and both wrap.
//
// PoetryDB is the provider that made the unwrapped form real. /title/Ozymandias
// answers [{...}]: no key, no envelope, and still a list, so client code reads
// body[0] exactly as it would read Xero's Invoices[0] with nothing to name it
// by. Its fetch of one poem is a get, and the get has to answer an array.
//
// This uses a shipped Recipe on purpose: the shape is the provider's, not a
// fixture's.
func TestASingleRecordCanBeAnsweredAsABareArray(t *testing.T) {
	r, err := recipe.Open("poetrydb")
	if err != nil {
		t.Fatalf("open poetrydb: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("shelley"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/title/Ozymandias", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	list, ok := body.([]any)
	if !ok {
		t.Fatalf("the body is not a list, so envelope.array did nothing: %.200s", rec.Body.String())
	}

	if len(list) != 1 {
		t.Fatalf("want one poem in the list, got %d", len(list))
	}

	poem, ok := list[0].(map[string]any)
	if !ok {
		t.Fatalf("the list holds something that is not the record: %#v", list[0])
	}

	if poem["title"] != "Ozymandias" {
		t.Errorf("the record is not the poem: %#v", poem)
	}

	// And nothing named it on the way out. A wrapped array would have put it
	// under a key, which is the shape this one is not.
	if _, wrapped := poem["poems"]; wrapped {
		t.Errorf("the record is wrapped inside the list: %#v", poem)
	}
}
