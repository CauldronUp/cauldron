package runtime

import (
	"encoding/json"
	"net/http"
	"testing"
)

// Scoped routes are how a provider partitions a resource by something in the
// path: a repository, an account, a store. Getting this wrong leaks one
// tenant's records into another's, which is the single most damaging thing a
// fake can teach an application to tolerate.

func issues(t *testing.T, s *Sandbox, path string) []map[string]any {
	t.Helper()

	rec := ghCall(t, s, http.MethodGet, path, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s = %d\n%s", path, rec.Code, rec.Body)
	}

	var out []map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("not an array: %v\n%s", err, rec.Body)
	}

	return out
}

func TestScopedListOnlySeesItsOwnScope(t *testing.T) {
	s := github(t)

	if err := s.Seed("small-repo"); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	// state=all: the listing narrows itself to open issues, and this test is
	// about the scope rather than the filter.
	octocat := issues(t, s, "/repos/octocat/hello-world/issues?state=all")
	if len(octocat) != 2 {
		t.Fatalf("octocat/hello-world has %d issues, want 2", len(octocat))
	}

	for _, issue := range octocat {
		if issue["owner"] != "octocat" {
			t.Errorf("a foreign record leaked into the scope: %v", issue)
		}
	}

	acme := issues(t, s, "/repos/acme/internal-tools/issues")
	if len(acme) != 1 {
		t.Fatalf("acme/internal-tools has %d issues, want 1", len(acme))
	}

	empty := issues(t, s, "/repos/nobody/nothing/issues")
	if len(empty) != 0 {
		t.Errorf("an unknown scope returned %d records, want none", len(empty))
	}
}

func TestCreateStampsTheScopeFromThePath(t *testing.T) {
	s := github(t)

	rec := ghCall(t, s, http.MethodPost, "/repos/acme/widgets/issues", `{"title":"Scoped on create"}`)

	// GitHub answers a create with 201, which the Recipe declares.
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	var created map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &created)

	if created["owner"] != "acme" || created["repo"] != "widgets" {
		t.Errorf("scope was not stamped from the path: %v", created)
	}

	// And it must be invisible from a different repository.
	if got := issues(t, s, "/repos/octocat/hello-world/issues"); len(got) != 0 {
		t.Errorf("record created under acme/widgets is visible under octocat/hello-world: %v", got)
	}
}

// A record that exists but belongs elsewhere must look absent rather than
// forbidden: leaking existence across scopes is itself a bug.
func TestGetAcrossScopesIsANotFound(t *testing.T) {
	s := github(t)
	_ = s.Seed("small-repo")

	// Issue 3 belongs to acme/internal-tools.
	if rec := ghCall(t, s, http.MethodGet, "/repos/acme/internal-tools/issues/3", ""); rec.Code != http.StatusOK {
		t.Fatalf("owner should be able to read it: %d", rec.Code)
	}

	rec := ghCall(t, s, http.MethodGet, "/repos/octocat/hello-world/issues/3", "")

	if rec.Code != http.StatusNotFound {
		t.Errorf("reading another scope's record = %d, want 404", rec.Code)
	}
}

func TestUpdateAcrossScopesIsRefused(t *testing.T) {
	s := github(t)
	_ = s.Seed("small-repo")

	rec := ghCall(t, s, http.MethodPatch, "/repos/octocat/hello-world/issues/3", `{"state":"closed"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("cross-scope update = %d, want 404", rec.Code)
	}

	// And the record must be untouched.
	var issue map[string]any
	body := ghCall(t, s, http.MethodGet, "/repos/acme/internal-tools/issues/3", "").Body.Bytes()
	_ = json.Unmarshal(body, &issue)

	if issue["state"] != "open" {
		t.Errorf("a refused cross-scope update still mutated the record: %v", issue["state"])
	}
}

// Paging must happen after filtering, or a page can come back empty while
// claiming there is more.
func TestScopedPagingCountsOnlyTheScope(t *testing.T) {
	s := github(t)

	for i := 0; i < 5; i++ {
		ghCall(t, s, http.MethodPost, "/repos/acme/widgets/issues", `{"title":"mine"}`)
		ghCall(t, s, http.MethodPost, "/repos/other/repo/issues", `{"title":"theirs"}`)
	}

	rec := ghCall(t, s, http.MethodGet, "/repos/acme/widgets/issues?limit=3", "")

	var page []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &page)

	if len(page) != 3 {
		t.Fatalf("got %d records, want 3", len(page))
	}

	for _, issue := range page {
		if issue["owner"] != "acme" {
			t.Errorf("paging crossed a scope boundary: %v", issue)
		}
	}
}
