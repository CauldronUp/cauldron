package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route's declared fields go in at a path, and a path may address a position
// in an array. The comparator that reads conformance cases had understood
// data.boards[0].name from the beginning; the code that writes the body had
// not, so it created a key literally spelled "boards[0]" and a case asserting
// the array form reported the field missing.
//
// The asymmetry is the whole bug: a Recipe could assert a shape it could not
// emit. A key like that has now been written four times by four different
// mechanisms, and every one of them was silent -- no provider sends it, so
// nothing downstream errors, it simply is not the field anyone asked for.

func mondayBody(t *testing.T, query string) map[string]any {
	t.Helper()

	r, err := recipe.Open("monday")
	if err != nil {
		t.Fatalf("open monday: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-board"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v2", strings.NewReader(string(body)))
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiJ9.cauldron.monday")
	req.Header.Set("API-Version", "2026-01")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v %s", err, rec.Body.String())
	}

	return out
}

func TestAnIndexedConstantPathBuildsAnArray(t *testing.T) {
	out := mondayBody(t, "query { boards { id name items_page { items { id } } } }")

	data, _ := out["data"].(map[string]any)
	if data == nil {
		t.Fatalf("no data object: %v", out)
	}

	if _, literal := data["boards[0]"]; literal {
		t.Fatal(`the index went in as part of the key: data["boards[0]"]`)
	}

	boards, ok := data["boards"].([]any)
	if !ok {
		t.Fatalf("boards is not an array: %T", data["boards"])
	}

	if len(boards) != 1 {
		t.Fatalf("expected one board, got %d", len(boards))
	}

	board, _ := boards[0].(map[string]any)
	if board == nil {
		t.Fatalf("the board is not an object: %v", boards[0])
	}

	if board["name"] != "Delivery" {
		t.Errorf("the constant did not reach the element: %v", board["name"])
	}
}

func TestAnIndexedConstantDoesNotDestroyItsSiblings(t *testing.T) {
	// Four constants share the data.boards[0] element and a fifth nests
	// below it. Whichever landed last used to be the only one there.
	out := mondayBody(t, "query { boards { id name items_count items_page { cursor items { id } } } }")

	data, _ := out["data"].(map[string]any)
	boards, _ := data["boards"].([]any)
	board, _ := boards[0].(map[string]any)

	for _, field := range []string{"id", "name", "items_count", "items_page"} {
		if _, present := board[field]; !present {
			t.Errorf("%s was erased by a sibling constant", field)
		}
	}

	page, _ := board["items_page"].(map[string]any)
	if page["cursor"] == nil {
		t.Error("the cursor was erased by the collection that shares its object")
	}

	if items, ok := page["items"].([]any); !ok || len(items) == 0 {
		t.Error("the collection was erased by the constants around it")
	}
}
