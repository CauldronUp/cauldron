package runtime

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A total that counts every record, not the page.
//
// The header exists because a Link header cannot say how far there is to go.
// Gitea's rel="last" is computed from the page size it actually served, which
// is not always the one the caller asked for -- send GitHub's per_page and
// Gitea ignores it, serves ten, and writes per_page back into the URLs it
// hands you. The count is then the only number in that response that survives
// having sent the wrong parameter name.

func countedListing(t *testing.T, name, fixture, target string) *http.Response {
	t.Helper()

	r, err := recipe.Open(name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed(fixture); err != nil {
		t.Fatalf("seed %s: %v", fixture, err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "token cauldron0000000000000000000000000000000")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec.Result()
}

func TestACountHeaderCountsTheWholeSetRatherThanThePage(t *testing.T) {
	t.Parallel()

	const path = "/api/v1/repos/gitea/tea/issues"

	// The whole set is three. A page of one still says three, which is the
	// only reason a client would read the header instead of len(body).
	page := countedListing(t, "gitea", "tea", path+"?limit=1")

	if page.StatusCode != http.StatusOK {
		t.Fatalf("listing answered %d, want 200", page.StatusCode)
	}

	if got := page.Header.Get("X-Total-Count"); got != "3" {
		t.Errorf("a page of one says %q records exist in total, want 3", got)
	}

	// And the last page says the same thing. A client that has walked to the
	// end still learns how many there were, without having counted.
	last := countedListing(t, "gitea", "tea", path+"?limit=1&page=3")

	if got := last.Header.Get("X-Total-Count"); got != "3" {
		t.Errorf("the last page says %q records exist in total, want 3", got)
	}

	// An integer, not a padded or quoted one: a client reads this with an
	// integer parse, and a stray space is a zero rather than an error.
	if _, err := strconv.Atoi(last.Header.Get("X-Total-Count")); err != nil {
		t.Errorf("the header does not parse as an integer: %v", err)
	}
}

// The count is of records, not of pages. Serving a page count under a name
// that says total is worse than sending nothing: the number is plausible and
// small, so a client sizing a progress bar by it finishes at a third.
func TestTheCountHeaderIsRecordsAndNotPages(t *testing.T) {
	t.Parallel()

	got := countedListing(t, "gitea", "tea",
		"/api/v1/repos/gitea/tea/issues?limit=1").Header.Get("X-Total-Count")

	if got == "3" {
		return
	}

	if got == "1" {
		t.Fatalf("X-Total-Count is %q, which is the page count, not the record count", got)
	}

	t.Fatalf("X-Total-Count is %q, want 3 records", got)
}

// The header is opt-in. A Recipe that has not seen its provider send one must
// not have this emulator invent it: a client taught to read a total here and
// finding none upstream parses the empty string as zero and decides the
// account is empty.
func TestNoCountHeaderUnlessTheRecipeNamesOne(t *testing.T) {
	t.Parallel()

	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	if r.Responses.List.CountHeader != "" {
		t.Skip("github now names a count header, so it cannot stand for the silent case")
	}

	got := issuesResponse(t, "/repos/octocat/hello-world/issues?state=all&per_page=1")

	for _, name := range []string{"X-Total-Count", "Total-Count", "X-Total"} {
		if value := got.Header.Get(name); value != "" {
			t.Errorf("a Recipe naming no count header sent %s: %q", name, value)
		}
	}
}
