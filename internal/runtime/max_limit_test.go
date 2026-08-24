package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some providers cap the page size and answer with less rather than refusing,
// which is the quieter of the two failures.
//
// Printify's own description says "default: 10, maximum: 10". So a client
// asking for a hundred orders is answered with ten and not told, and a paging
// loop that stops when it receives fewer records than it asked for stops on
// the first page -- a shop with four hundred orders reports ten, and nothing
// errored.
//
// Before max_limit the declared limit was only a default and the caller always
// won, so a Recipe could describe the cap in a comment and serve something
// else. Zero means the provider serves whatever is asked for, which is what
// every Recipe written before this assumed, so nothing that shipped changes.
func TestAProviderMayCapThePageSize(t *testing.T) {
	r, err := recipe.Open("printify")
	if err != nil {
		t.Fatalf("open printify: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-shop"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	perPage := func(query string) float64 {
		t.Helper()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/shops/9000/orders.json"+query, nil)
		req.Header.Set("Authorization", "Bearer cauldron.printify.token.0000")
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", query, rec.Code)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: decode: %v", query, err)
		}

		got, ok := body["per_page"].(float64)
		if !ok {
			t.Fatalf("%s: per_page = %#v, want a number", query, body["per_page"])
		}

		return got
	}

	// Asking for more than the cap is answered with the cap.
	if got := perPage("?limit=100"); got != 10 {
		t.Errorf("asking for 100 served %v, want the cap of 10", got)
	}

	// Asking for less than the cap is honoured, because a cap is a ceiling
	// rather than a fixed size -- a Recipe that ignored the request entirely
	// would be a different lie.
	if got := perPage("?limit=1"); got != 1 {
		t.Errorf("asking for 1 served %v, want 1", got)
	}

	// And asking for nothing gets the declared default.
	if got := perPage(""); got != 10 {
		t.Errorf("asking for nothing served %v, want the default of 10", got)
	}
}
