package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A create that names its own identifier is stored under that identifier.
//
// Pinecone's POST /indexes sends {"name": "order-notes"} and the name is the
// index's identifier: there is no separate id anywhere in the API. The
// resource says so with id.field: name.
//
// That used to store nothing under id at all. The wire name was kept as an
// ordinary key, Store.Create found no identifier and minted one, and the
// response was then built with two writers for a single key -- the value the
// caller sent, under name, and the minted identifier rendered under the same
// name. Which one survived depended on map iteration order.
//
// So the same request answered "order-notes" or "rfBd56ti2SMtYv" from the same
// binary, and the conformance suite passed on Linux while failing on macOS and
// Windows in the same CI run.
//
// A non-deterministic fake is the worst failure this emulator can have. A flaky
// suite sends somebody looking through their own code first, and everything
// they find there will be innocent.
//
// The assertion is a fetch rather than a comparison of the create's body,
// because a fetch cannot be a coin flip: before the fix the record was filed
// under a minted identifier, so asking for it by the name the caller chose was
// a 404 every time.
func TestACreateIsFiledUnderTheIdentifierItWasGiven(t *testing.T) {
	r, err := recipe.Open("pinecone")
	if err != nil {
		t.Fatalf("open pinecone: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("project"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	const key = "pcsk_cauldronfixture_NotARealPineconeKey0000000"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/indexes", strings.NewReader(
		`{"name": "order-notes", "dimension": 1536, "metric": "cosine"}`))
	req.Header.Set("Api-Key", key)
	req.Header.Set("X-Pinecone-Api-Version", "2026-04")
	req.Header.Set("Content-Type", "application/json")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	if got := created["name"]; got != "order-notes" {
		t.Errorf("the create answered with name %v, not the name it was sent", got)
	}

	// The part that cannot be a coin flip.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/indexes/order-notes", nil)
	req.Header.Set("Api-Key", key)
	req.Header.Set("X-Pinecone-Api-Version", "2026-04")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("the created index is not filed under the name it was given: status = %d, body = %s",
			rec.Code, rec.Body.String())
	}

	var fetched map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode fetch: %v", err)
	}

	if got := fetched["name"]; got != "order-notes" {
		t.Errorf("fetched name = %v, want order-notes", got)
	}

	// And nothing else claims to be an identifier. The record carried the wire
	// name and a minted id at once before, which is what made the response
	// ambiguous in the first place.
	if _, present := fetched["id"]; present {
		t.Errorf("the response carries an id beside its wire name: %v", fetched)
	}
}

// The rename does not apply when the resource declares a field of that name
// too. Then the Recipe has said the identifier and the field are different
// things, and the field keeps the key: the identifier still has "id" to arrive
// under.
//
// Gemini is the case in the collection. Its model resource publishes the
// identifier as name -- a path rather than a word -- and declares displayName
// beside it, so nothing collides; what matters here is that a resource which
// does declare the wire name as a field is left alone.
func TestAFieldNamedLikeTheIdentifierIsStillAField(t *testing.T) {
	r, err := recipe.Open("gemini")
	if err != nil {
		t.Fatalf("open gemini: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	spec := r.Resources["model"]
	if spec.ID.Field != "name" {
		t.Fatalf("gemini's model no longer publishes its identifier as name; pick another resource")
	}

	kept := s.declaredOnly("model", map[string]any{"name": "models/gemini-2.5-flash"})

	if got := kept["id"]; got != "models/gemini-2.5-flash" {
		t.Errorf("the wire name did not become the identifier: %v", kept)
	}

	if _, present := kept["name"]; present {
		t.Errorf("the wire name was kept beside the identifier: %v", kept)
	}
}
