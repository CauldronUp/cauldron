package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Three different quantities travel under the name "total", and Shopware sends
// all three from one endpoint depending on what the request asked for.
//
// The default is the length of the page in front of you, because counting rows
// costs a second query and the framework does not run one uninvited:
// Criteria.php starts at TOTAL_COUNT_MODE_NONE and EntitySearcher::getTotalCount
// returns count($data) for every mode except exact. So a shop with four hundred
// products answers a ten-record page with total: 10 -- a number that is not
// wrong about anything except the question it looks like it is answering.
//
// A fixture small enough to fit on one page makes all three modes agree, which
// is exactly why a Recipe could describe one and serve another and no
// conformance case would notice. Nine products at one per page is the smallest
// catalogue that tells them apart: the page reports 1, the exact count reports
// 9, and the six-page window reports 7.
func TestTheCountFieldNeedNotBeTheCount(t *testing.T) {
	r, err := recipe.Open("shopware")
	if err != nil {
		t.Fatalf("open shopware: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	total := func(query string) float64 {
		t.Helper()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/store-api/product"+query, nil)
		req.Header.Set("sw-access-key", "SWSCCAULDRONFIXTUREKEY00")
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200: %s", query, rec.Code, rec.Body.String())
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: decode: %v", query, err)
		}

		got, ok := body["total"].(float64)
		if !ok {
			t.Fatalf("%s: total = %#v, want a number", query, body["total"])
		}

		return got
	}

	// The default: one record on the page, so one is the total. Nine products
	// exist and nothing in the response says so.
	if got := total("?limit=1"); got != 1 {
		t.Errorf("the default mode reported %v, want the page length of 1", got)
	}

	// Opting in buys the second query and the real answer.
	if got := total("?limit=1&total-count-mode=exact"); got != 9 {
		t.Errorf("exact reported %v, want the whole catalogue of 9", got)
	}

	// And the window: limit * 6 + 1 rows are fetched and counted before the
	// slice, so seven of nine, with the seventh row being the sentinel that
	// says there is more beyond the window.
	if got := total("?limit=1&total-count-mode=next-pages"); got != 7 {
		t.Errorf("the lookahead reported %v, want the bounded 7", got)
	}

	// A lookahead is bounded rather than fixed: where the collection is
	// smaller than the window, the window is not what comes back. Three per
	// page reaches nineteen, which is past the whole catalogue, so the answer
	// is the catalogue.
	if got := total("?limit=3&total-count-mode=next-pages"); got != 9 {
		t.Errorf("the lookahead over a short collection reported %v, want 9", got)
	}
}

// A provider that caps the page size can answer with less or refuse outright,
// and Shopware does both -- on two listings of the same resource.
//
// /store-api/product runs the raw criteria builder, which throws
// QueryLimitExceededException: a 400 naming the ceiling it enforced.
// /store-api/product-listing/{categoryId} runs through PagingListingProcessor,
// which calls min($limit, $this->maxLimit) and says nothing. Same hundred, same
// provider, same resource, opposite answers -- and a client written against
// either one is broken against the other in the opposite direction.
func TestAProviderMayRefuseAnOversizedPageRatherThanTrimIt(t *testing.T) {
	r, err := recipe.Open("shopware")
	if err != nil {
		t.Fatalf("open shopware: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	ask := func(path string) (int, map[string]any) {
		t.Helper()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("sw-access-key", "SWSCCAULDRONFIXTUREKEY00")
		s.ServeHTTP(rec, req)

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: decode: %v", path, err)
		}

		return rec.Code, body
	}

	// The entity route refuses.
	status, body := ask("/store-api/product?limit=250")
	if status != http.StatusBadRequest {
		t.Errorf("an oversized page on the entity route = %d, want 400", status)
	}

	errors, _ := body["errors"].([]any)
	if len(errors) != 1 {
		t.Fatalf("errors = %#v, want one entry", body["errors"])
	}

	if entry, _ := errors[0].(map[string]any); entry["code"] != "FRAMEWORK__QUERY_LIMIT_EXCEEDED" {
		t.Errorf("code = %#v, want FRAMEWORK__QUERY_LIMIT_EXCEEDED", entry["code"])
	}

	// The category listing trims the same request to the same ceiling and
	// answers 200.
	status, body = ask("/store-api/product-listing/4a1c8e0b6d2f47a39c5e1b7d0f836a24?limit=250")
	if status != http.StatusOK {
		t.Errorf("an oversized page on the category listing = %d, want 200", status)
	}

	if got, _ := body["limit"].(float64); got != 100 {
		t.Errorf("the trimmed page size = %v, want the ceiling of 100", got)
	}

	// And a request under the ceiling is served as asked on both, because a
	// ceiling is not a fixed size.
	if status, _ = ask("/store-api/product?limit=2"); status != http.StatusOK {
		t.Errorf("a page under the ceiling = %d, want 200", status)
	}
}
