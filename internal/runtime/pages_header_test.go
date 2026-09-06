package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Three quantities, and a provider can send any two of them without the third.
//
// Vikunja is the case. Its own description says every paginated endpoint
// answers `x-pagination-total-pages`, "the total number of available pages for
// this request", and `x-pagination-result-count`, "the number of items returned
// for this request". Pages, and this page. Nowhere does anything say how many
// items exist -- so a caller can compute an upper bound by multiplying and can
// never get the number.
//
// That is why pages_header is a separate key from count_header rather than a
// renaming of it, and why count_means applies to the header as it does to the
// body.

func vikunjaLabels(t *testing.T, target string) *http.Response {
	t.Helper()

	r, err := recipe.Open("vikunja")
	if err != nil {
		t.Fatalf("open vikunja: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("board"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer cauldron.vikunja.jwt")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec.Result()
}

func TestThePagesHeaderCountsPagesAndTheCountHeaderCountsThisPage(t *testing.T) {
	t.Parallel()

	// Two labels, one per page.
	got := vikunjaLabels(t, "/api/v1/labels?page=1&per_page=1")

	if got.StatusCode != http.StatusOK {
		t.Fatalf("listing answered %d, want 200", got.StatusCode)
	}

	if pages := got.Header.Get("x-pagination-total-pages"); pages != "2" {
		t.Errorf("x-pagination-total-pages is %q, want 2 -- two records at one a page", pages)
	}

	// And the count is the page, because Vikunja's is "the number of items
	// returned for this request". A Recipe that meant the whole set would say
	// so by leaving count_means out, and this one does not.
	if count := got.Header.Get("x-pagination-result-count"); count != "1" {
		t.Errorf("x-pagination-result-count is %q, want 1 -- it is this page, not the set", count)
	}
}

// The two move independently. On the last page there is still one record, and
// the page count has not changed -- so a client reading one of them cannot
// derive the other.
func TestThePageCountDoesNotMoveWithThePage(t *testing.T) {
	t.Parallel()

	got := vikunjaLabels(t, "/api/v1/labels?page=2&per_page=1")

	if pages := got.Header.Get("x-pagination-total-pages"); pages != "2" {
		t.Errorf("x-pagination-total-pages on page two is %q, want 2", pages)
	}

	if count := got.Header.Get("x-pagination-result-count"); count != "1" {
		t.Errorf("x-pagination-result-count on page two is %q, want 1", count)
	}
}

// And neither of them is the number of items, which is the finding. Two records
// exist and no header says two.
func TestNoHeaderSaysHowManyItemsExist(t *testing.T) {
	t.Parallel()

	got := vikunjaLabels(t, "/api/v1/labels?page=1&per_page=1")

	for _, name := range []string{"x-pagination-total-pages", "x-pagination-result-count", "X-Total-Count"} {
		if got.Header.Get(name) == "2" && name == "X-Total-Count" {
			t.Errorf("%s says 2, and Vikunja sends no such header -- the emulator invented a total", name)
		}
	}

	if got.Header.Get("X-Total-Count") != "" {
		t.Error("a total-count header was sent, and Vikunja's description names only the two pagination headers")
	}
}

// The header is opt-in, the same as the count. A Recipe that has not seen its
// provider send a page count must not have one invented, because a client
// taught to read it upstream finds an empty string and parses it as zero pages.
func TestNoPagesHeaderUnlessTheRecipeNamesOne(t *testing.T) {
	t.Parallel()

	r, err := recipe.Open("gitea")
	if err != nil {
		t.Fatalf("open gitea: %v", err)
	}

	if r.Responses.List.PagesHeader != "" {
		t.Skip("gitea now names a pages header, so it cannot stand for the silent case")
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("tea"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/repos/gitea/tea/issues?limit=1", nil)
	req.Header.Set("Authorization", "token cauldron0000000000000000000000000000000")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	for _, name := range []string{"x-pagination-total-pages", "X-Total-Pages", "X-Pages"} {
		if value := rec.Result().Header.Get(name); value != "" {
			t.Errorf("a Recipe naming no pages header sent %s: %q", name, value)
		}
	}

	// And the count header it does name is still the whole set, because gitea
	// declares no count_means.
	if got := rec.Result().Header.Get("X-Total-Count"); got != "3" {
		t.Errorf("gitea's X-Total-Count is %q, want 3 -- count_means is absent, so it is the set", got)
	}
}
