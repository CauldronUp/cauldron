// A paging pointer that is a number rather than a token or an address.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Creem's listing envelope is {total_records, total_pages, current_page,
// next_page, prev_page} and every one of those is a JSON number: next_page is
// 2, null on the last page, and prev_page is 1, null on the first.
//
// Rendered as the string "2" it reads identically in a diff and is a different
// value on the wire. A client adding one to it -- which is the whole reason a
// page number is a number rather than an opaque token -- gets "21" instead of
// 3, and finds out in production rather than against the fake.
//
// A shipped Recipe on purpose: the shape is Creem's, not a fixture's.
func TestAPageNumberedListingSendsNumbersRatherThanTokens(t *testing.T) {
	body := creemPage(t, "?page_size=1&page_number=1")

	if got, ok := body["next_page"].(float64); !ok || got != 2 {
		t.Fatalf("next_page = %T %v, want the number 2", body["next_page"], body["next_page"])
	}

	if body["prev_page"] != nil {
		t.Fatalf("prev_page = %v, want null on the first page", body["prev_page"])
	}

	body = creemPage(t, "?page_size=1&page_number=2")

	if got, ok := body["prev_page"].(float64); !ok || got != 1 {
		t.Fatalf("prev_page = %T %v, want the number 1", body["prev_page"], body["prev_page"])
	}

	if body["next_page"] != nil {
		t.Fatalf("next_page = %v, want null on the last page", body["next_page"])
	}
}

// And a Recipe that says nothing keeps what it had. Emitting a number where a
// provider sends a token is the same mistake pointed the other way, so the
// declaration is opt in exactly as cursor_url beside it is.
func TestWithoutTheDeclarationThePagingPointerStaysAString(t *testing.T) {
	r, err := recipe.Open("creem")
	if err != nil {
		t.Fatalf("open creem: %v", err)
	}

	r.Responses.List.CursorNumber = false

	body := pageOf(t, r, "?page_size=1&page_number=2")

	if _, ok := body["next_page"].(float64); ok {
		t.Errorf("next_page = %v, want a string once the declaration is gone", body["next_page"])
	}

	if _, ok := body["prev_page"].(string); !ok {
		t.Errorf("prev_page = %T %v, want the address form", body["prev_page"], body["prev_page"])
	}
}

func creemPage(t *testing.T, query string) map[string]any {
	t.Helper()

	r, err := recipe.Open("creem")
	if err != nil {
		t.Fatalf("open creem: %v", err)
	}

	return pageOf(t, r, query)
}

func pageOf(t *testing.T, r *recipe.Recipe, query string) map[string]any {
	t.Helper()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed catalogue: %v", err)
	}

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/v1/products/search"+query, nil)
	req.Header.Set("x-api-key", r.Auth.Keys[0])

	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Pagination map[string]any `json:"pagination"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v: %s", err, rec.Body.String())
	}

	return body.Pagination
}
