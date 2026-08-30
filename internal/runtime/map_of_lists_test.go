// A map whose values are lists.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// serveMapRecipe parses a Recipe from source and seeds one fixture, so the
// shape under test can be written out in full beside the test that reads it.
func serveMapRecipe(t *testing.T, body, fixture string) *Sandbox {
	t.Helper()

	r, err := recipe.Parse([]byte(body))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed(fixture); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func withoutLine(t *testing.T, body, line string) string {
	t.Helper()

	if strings.Count(body, line) != 1 {
		t.Fatalf("the line to remove appears %d times", strings.Count(body, line))
	}

	return strings.Replace(body, line, "", 1)
}

// PDBe answers every entry endpoint with {"4hhb": [ {...} ]} -- a map keyed by
// the identifier, whose value is an array holding one record. The map style
// already describes the key; what it could not describe is the array.
//
// The difference matters to a client and not a little. res["4hhb"].title is
// undefined and res["4hhb"][0].title is the answer, and an emulator that served
// the object directly would let a Recipe pass with the first, which is the
// reading that breaks against the real provider.
const mapOfLists = `
recipe: pdbe
capability: docs
version: 0.1.0
upstream:
  api: v1
auth:
  scheme: none
responses:
  list:
    style: map
    key: "-"
    entry_style: list
resources:
  entry:
    id:
      style: opaque
    fields:
      title:
        type: string
routes:
  - method: GET
    path: /entries
    resource: entry
    operation: list
fixtures:
  one:
    entry:
      - id: "4hhb"
        title: HAEMOGLOBIN
`

func TestAMapStyleCanHoldListsRatherThanObjects(t *testing.T) {
	server := serveMapRecipe(t, mapOfLists, "one")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/entries", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}

	held, ok := body["4hhb"]
	if !ok {
		t.Fatalf("the map is not keyed by the identifier: %s", rec.Body.String())
	}

	entries, ok := held.([]any)
	if !ok {
		t.Fatalf("the value is %T, want a list: %s", held, rec.Body.String())
	}

	if len(entries) != 1 {
		t.Fatalf("the list holds %d entries, want 1", len(entries))
	}

	record, ok := entries[0].(map[string]any)
	if !ok {
		t.Fatalf("the entry is %T, want an object", entries[0])
	}

	if record["title"] != "HAEMOGLOBIN" {
		t.Errorf("title is %v, want HAEMOGLOBIN", record["title"])
	}

	// The identifier is the key and nothing else, exactly as upstream: PDBe's
	// record carries no field holding it.
	if _, repeated := record["id"]; repeated {
		t.Errorf("the identifier is repeated inside the value: %v", record)
	}
}

// Without the declaration the map keeps holding objects, so no shipped Recipe
// changes shape underneath itself.
func TestAMapStyleStillHoldsObjectsByDefault(t *testing.T) {
	server := serveMapRecipe(t, withoutLine(t, mapOfLists, "    entry_style: list\n"), "one")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/entries", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}

	if _, isList := body["4hhb"].([]any); isList {
		t.Errorf("the map holds a list without being asked to: %s", rec.Body.String())
	}

	if _, isObject := body["4hhb"].(map[string]any); !isObject {
		t.Errorf("the map does not hold an object: %s", rec.Body.String())
	}
}
