package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Providers disagree about what to call their first page, and the
// disagreement is invisible. Algolia, Elasticsearch and everything shaped like
// them count from nought, so page 1 is the second page. Read as one-based, a
// client asking for page 1 is handed page 0 again: the same records twice, no
// error, and a loop that either never terminates or quietly returns
// duplicates.
//
// positionOf already carried a comment warning about losing or duplicating a
// record at every page break. It was hard-coded one-based, so it was doing
// exactly that for every provider that counts from nought.

func algoliaSearch(t *testing.T, body string) map[string]any {
	t.Helper()

	r, err := recipe.Open("algolia")
	if err != nil {
		t.Fatalf("open algolia: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/1/indexes/products/query", strings.NewReader(body))
	req.Header.Set("X-Algolia-API-Key", "cauldron-admin-key")
	req.Header.Set("X-Algolia-Application-Id", "CAULDRON01")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	return out
}

func firstHit(t *testing.T, out map[string]any) string {
	t.Helper()

	hits, ok := out["hits"].([]any)
	if !ok || len(hits) == 0 {
		t.Fatalf("no hits: %v", out)
	}

	hit, _ := hits[0].(map[string]any)
	id, _ := hit["objectID"].(string)

	return id
}

func TestAZeroBasedProviderCountsFromNought(t *testing.T) {
	first := algoliaSearch(t, `{"query":"","hitsPerPage":1,"page":0}`)
	second := algoliaSearch(t, `{"query":"","hitsPerPage":1,"page":1}`)

	if got := firstHit(t, first); got != "prod-0001" {
		t.Errorf("page 0 should be the first record, got %s", got)
	}

	// The whole point: page 1 is the second page, not the first one again.
	if got := firstHit(t, second); got != "prod-0002" {
		t.Errorf("page 1 should be the second record, got %s", got)
	}
}

func TestAPageBelowTheFirstIsTheFirstPage(t *testing.T) {
	// A number below the provider's own first page answers with the first
	// page, which is the reading that cannot silently skip records.
	out := algoliaSearch(t, `{"query":"","hitsPerPage":1,"page":-3}`)

	if got := firstHit(t, out); got != "prod-0001" {
		t.Errorf("a page below the first should be the first, got %s", got)
	}
}

func TestTheEnvelopeReportsThePageItServed(t *testing.T) {
	// These were constants -- 0 and 20 -- so a client that asked for page 1
	// was told it was looking at page 0, by the field whose entire purpose is
	// to say where you are.
	out := algoliaSearch(t, `{"query":"","hitsPerPage":1,"page":1}`)

	if page, _ := out["page"].(float64); page != 1 {
		t.Errorf("the response reported page %v for a request for page 1", out["page"])
	}

	if size, _ := out["hitsPerPage"].(float64); size != 1 {
		t.Errorf("the response reported hitsPerPage %v for a request for 1", out["hitsPerPage"])
	}
}

func TestTheEnvelopeFallsBackToTheProvidersDefaults(t *testing.T) {
	// A request that asks for no paging is told the provider's own defaults,
	// which is what Algolia does and what the constants used to hard-code.
	out := algoliaSearch(t, `{"query":"cauldron"}`)

	if page, _ := out["page"].(float64); page != 0 {
		t.Errorf("the default page should be nought here, got %v", out["page"])
	}

	if size, _ := out["hitsPerPage"].(float64); size != 20 {
		t.Errorf("the default page size should be 20, got %v", out["hitsPerPage"])
	}
}
