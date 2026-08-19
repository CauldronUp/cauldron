package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A declared constant at a dotted name nests, the same way a renamed
// identifier and a field's "in" already do. Attio's record identifier is an
// object of three UUIDs, and the two that are constant used to go out as one
// literal key spelled "id.workspace_id" -- a shape no provider sends, produced
// in silence, and the fifth mechanism in the runtime to have made that
// mistake.
//
// The exception is what makes it a rule worth stating: Dropbox names a field
// ".tag", where the leading dot is part of the name. A path is at least two
// segments and every one of them is a name, so ".tag" is a key and
// "id.workspace_id" is a path.

func attioRecord(t *testing.T) map[string]any {
	t.Helper()

	r, err := recipe.Open("attio")
	if err != nil {
		t.Fatalf("open attio: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-workspace"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v2/objects/people/records/query",
		strings.NewReader(`{"limit":1}`))
	req.Header.Set("Authorization", "Bearer attio_cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out struct {
		Data []map[string]any `json:"data"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if len(out.Data) == 0 {
		t.Fatal("no records")
	}

	return out.Data[0]
}

func TestADottedConstantNests(t *testing.T) {
	record := attioRecord(t)

	if _, literal := record["id.workspace_id"]; literal {
		t.Fatal(`the constant went in as one key: "id.workspace_id"`)
	}

	id, ok := record["id"].(map[string]any)
	if !ok {
		t.Fatalf("the identifier is not an object: %T", record["id"])
	}

	// All three, from three different places: the store's own key, a
	// declared constant, and a field with an "in".
	for _, part := range []string{"record_id", "workspace_id", "object_id"} {
		if id[part] == nil {
			t.Errorf("%s is missing from the identifier object", part)
		}
	}
}

func TestALeadingDotIsPartOfTheName(t *testing.T) {
	// Dropbox's union tag. Nesting this would put it under an empty key.
	r, err := recipe.Open("dropbox")
	if err != nil {
		t.Fatalf("open dropbox: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/2/files/list_folder",
		strings.NewReader(`{"path":"/documents"}`))
	req.Header.Set("Authorization", "Bearer sl.cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	var out struct {
		Entries []map[string]any `json:"entries"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if len(out.Entries) == 0 {
		t.Fatal("no entries")
	}

	if out.Entries[0][".tag"] == nil {
		t.Errorf(`.tag was treated as a path rather than a name: %v`, out.Entries[0])
	}

	if _, nested := out.Entries[0][""]; nested {
		t.Error(`.tag was split, and the empty first segment became a key`)
	}
}

func TestIsPathSeparatesAPathFromAName(t *testing.T) {
	for _, path := range []string{"id.workspace_id", "sys.id", "a.b.c", "to[0].name"} {
		if !recipe.IsPath(path) {
			t.Errorf("%q should be a path", path)
		}
	}

	for _, name := range []string{".tag", "tag.", "id", "", "a..b", "dist-tags"} {
		if recipe.IsPath(name) {
			t.Errorf("%q should be a name", name)
		}
	}
}
