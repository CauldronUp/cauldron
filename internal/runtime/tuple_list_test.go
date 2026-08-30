// A listing whose envelope is a two-element array.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The World Bank's v2 API answers [{page, pages, per_page, total}, [...]]: the
// paging object first, the records second. No list style produced that -- bare
// gives the array alone, wrapped and map give an object, and the default gives
// Stripe's {object, data, has_more} -- so the shape had to be modelled as
// something it is not, or the provider left out.
//
// It is worth having as its own style rather than approximated, because of what
// the failure looks like. A bad country code answers [{message: [...]}] -- one
// element, not two -- with the same HTTP 200, so body[1] is the collection on
// success and undefined on failure, and the length of the outer array is the
// only structural signal that anything went wrong. A Recipe that served an
// object would let a client write body.data.map(...) and ship it.
func TestAListingCanBeATwoElementArray(t *testing.T) {
	r, err := recipe.Open("worldbank")
	if err != nil {
		t.Fatalf("open worldbank: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("countries"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v2/country/CAN?format=json", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body []any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("the body is not an array: %v -- %.200s", err, rec.Body.String())
	}

	if len(body) != 2 {
		t.Fatalf("want two elements, got %d: %.200s", len(body), rec.Body.String())
	}

	meta, ok := body[0].(map[string]any)
	if !ok {
		t.Fatalf("the first element is not the paging object: %#v", body[0])
	}

	// per_page is quoted and its three neighbours are not, which is the
	// difference a Recipe declaring this shape exists to keep.
	if meta["per_page"] != "50" {
		t.Errorf("per_page = %#v, want the string \"50\"", meta["per_page"])
	}

	if meta["total"] != float64(1) {
		t.Errorf("total = %#v, want the number 1", meta["total"])
	}

	records, ok := body[1].([]any)
	if !ok {
		t.Fatalf("the second element is not the collection: %#v", body[1])
	}

	if len(records) != 1 {
		t.Fatalf("want one country, got %d", len(records))
	}

	country, ok := records[0].(map[string]any)
	if !ok || country["name"] != "Canada" {
		t.Errorf("the collection does not hold the country: %#v", records[0])
	}
}
