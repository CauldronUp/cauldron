package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// pagedSandbox mounts a Recipe and seeds it.
func pagedSandbox(t *testing.T, name, fixture string) *Sandbox {
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
		t.Fatalf("seed: %v", err)
	}

	return s
}

// pagination.style was declared on 149 shipped routes and read by nothing, so
// every offset and page-numbered listing answered with the first page whatever
// was asked for.
//
// The consequence is not a slow test. A client looping until it has collected
// total_count records never finishes, and a client asking for page two
// processes page one a second time — which for a payments or messaging
// integration is repeated side effects rather than a wasted request.
func TestAnOffsetPagedListingHonoursTheOffset(t *testing.T) {
	// PagerDuty, whose own OpenAPI description declares limit and offset, so
	// the parameter names this exercises are verified rather than assumed.
	// Clerk was here first and its paging style has been withdrawn: it was
	// declared without naming its parameters, and nobody had checked what
	// Clerk calls them.
	s := pagedSandbox(t, "pagerduty", "small-account")

	get := func(query string) []any {
		req := httptest.NewRequest(http.MethodGet, "/incidents"+query, nil)
		req.Header.Set("Authorization", "Token token=cauldron_api_key")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /incidents%s = %d\n%s", query, rec.Code, rec.Body)
		}

		var body struct {
			Incidents []any `json:"incidents"`
		}

		decodeInto(t, rec, &body)

		return body.Incidents
	}

	all := get("?limit=100")
	if len(all) < 3 {
		t.Fatalf("expected the fixture to seed at least 3 incidents, got %d", len(all))
	}

	// Walking one at a time must visit every record exactly once, which is the
	// property that was broken: every offset returned the same first record.
	seen := map[string]bool{}

	for offset := 0; offset < len(all); offset++ {
		page := get(fmt.Sprintf("?limit=1&offset=%d", offset))

		if len(page) != 1 {
			t.Fatalf("offset=%d returned %d records, want 1", offset, len(page))
		}

		id := fmt.Sprint(page[0].(map[string]any)["id"])

		if seen[id] {
			t.Fatalf("offset=%d returned %s again, so a paging loop would repeat it forever", offset, id)
		}

		seen[id] = true
	}

	if len(seen) != len(all) {
		t.Errorf("walking by offset saw %d of %d records", len(seen), len(all))
	}

	// Past the end is an empty page, which is what stops a loop rather than
	// breaking it.
	if beyond := get(fmt.Sprintf("?limit=10&offset=%d", len(all)+5)); len(beyond) != 0 {
		t.Errorf("an offset past the end returned %d records", len(beyond))
	}
}

// A page number counts pages from one, so page two begins one page-length in.
// Getting that boundary wrong by a page loses or repeats a record at every
// page break, which is the bug an emulator exists to catch rather than commit.
func TestAPageNumberedListingHonoursThePage(t *testing.T) {
	s := pagedSandbox(t, "gitlab", "small-instance")

	get := func(query string) []any {
		req := httptest.NewRequest(http.MethodGet, "/api/v4/projects"+query, nil)
		req.Header.Set("PRIVATE-TOKEN", "glpat-cauldron")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v4/projects%s = %d\n%s", query, rec.Code, rec.Body)
		}

		var body []any

		decodeInto(t, rec, &body)

		return body
	}

	all := get("?per_page=100")
	if len(all) < 3 {
		t.Fatalf("expected at least 3 projects, got %d", len(all))
	}

	first := get("?per_page=1&page=1")
	second := get("?per_page=1&page=2")
	third := get("?per_page=1&page=3")

	for i, page := range [][]any{first, second, third} {
		if len(page) != 1 {
			t.Fatalf("page %d returned %d records, want 1", i+1, len(page))
		}
	}

	ids := []string{
		fmt.Sprint(first[0].(map[string]any)["id"]),
		fmt.Sprint(second[0].(map[string]any)["id"]),
		fmt.Sprint(third[0].(map[string]any)["id"]),
	}

	if ids[0] == ids[1] || ids[1] == ids[2] || ids[0] == ids[2] {
		t.Errorf("three consecutive pages returned overlapping records: %v", ids)
	}

	// Page one and no page at all are the same page. Nothing may be skipped by
	// treating an absent page number as page nought.
	if bare := get("?per_page=1"); fmt.Sprint(bare[0].(map[string]any)["id"]) != ids[0] {
		t.Errorf("no page parameter gave a different first page: %v against %v", bare[0], first[0])
	}
}

// decodeInto reads a JSON response body into any shape, for the listings that
// are a bare array rather than an object.
func decodeInto(t *testing.T, rec *httptest.ResponseRecorder, out any) {
	t.Helper()

	if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
		t.Fatalf("decoding the response: %v\n%s", err, rec.Body)
	}
}
