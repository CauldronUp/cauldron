package runtime

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Five providers modelled here advertise the next page in an RFC 5988 Link
// header, and for two of them it is the only mechanism there is. Without it
// the page size works and the next page is unreachable, which is the quietest
// way for a listing to be wrong: one page comes back, it is a correct page,
// and the loop that should have asked for the second one has nothing to
// follow.

func issuesResponse(t *testing.T, target string) *http.Response {
	t.Helper()

	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-repo"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer ghp_cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec.Result()
}

func TestALinkHeaderAdvertisesTheNextPage(t *testing.T) {
	got := issuesResponse(t, "/repos/octocat/hello-world/issues?state=all&per_page=1").Header.Get("Link")

	if got == "" {
		t.Fatal("no Link header on a truncated listing")
	}

	// The same query, one page further in, as a URL a client can request as
	// it stands.
	want := regexp.MustCompile(`^<http://[^>]*/repos/octocat/hello-world/issues\?[^>]*page=2[^>]*>; rel="next"$`)
	if !want.MatchString(got) {
		t.Errorf("Link header is not a next-page URL: %q", got)
	}

	// The rest of the query has to survive, or the next page is a different
	// question from the one that was asked. state=all is the filter here, and
	// dropping it would silently narrow page two to open issues.
	if !regexp.MustCompile(`state=all`).MatchString(got) {
		t.Errorf("the next page dropped the filter: %q", got)
	}
}

func TestTheLastPageAdvertisesNoNext(t *testing.T) {
	// How the loop ends, and this test had the same thing backwards that the
	// Recipe did: it asked for no Link header, and GitHub's last page carries
	// one. Checked against api.github.com on 2026-08-22, golang/example with
	// 74 issues: ?per_page=50&page=2 answers rel="prev" and nothing else.
	//
	// What has to be absent is the next, not the header. A client stopping on
	// a missing header never stops against GitHub.
	got := issuesResponse(t, "/repos/octocat/hello-world/issues?state=all&per_page=1&page=2").Header.Get("Link")

	if strings.Contains(got, `rel="next"`) {
		t.Errorf("the last page advertised a next page: %q", got)
	}

	if !strings.Contains(got, `rel="prev"`) {
		t.Errorf("the last page should still point back: %q", got)
	}
}

// A page reached from another page points back; a listing that never paged has
// nowhere to point. The two are different states and the header says so.
func TestAMiddlePageAdvertisesBoth(t *testing.T) {
	got := issuesResponse(t, "/repos/octocat/hello-world/issues?state=all&per_page=1&page=2").Header.Get("Link")

	if strings.Contains(got, `rel="next"`) && !strings.Contains(got, `rel="prev"`) {
		t.Errorf("a next without a prev on a page reached from another: %q", got)
	}
}

func TestAnUnpagedListingAdvertisesNothing(t *testing.T) {
	// Nothing was truncated, so there is no next page to name.
	if got := issuesResponse(t, "/repos/octocat/hello-world/issues?state=all").Header.Get("Link"); got != "" {
		t.Errorf("a complete listing advertised a next page: %q", got)
	}
}

func TestTheNextPageKeepsThePathTheCallerAsked(t *testing.T) {
	// The multi-provider server mounts each Recipe under its own name and
	// rewrites URL.Path, leaving RequestURI as it arrived. A header built
	// from the rewritten path sends the client to /repos/... when the server
	// serves /github/repos/... -- a next page that 404s, which is worse than
	// no next page at all.
	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-repo"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/github/repos/octocat/hello-world/issues?state=all&per_page=1", nil)
	req.Header.Set("Authorization", "Bearer ghp_cauldron")

	// What the server does before handing the request on: the sandbox sees
	// its own URL space, and RequestURI keeps what the caller sent.
	clone := req.Clone(req.Context())
	clone.URL.Path = "/repos/octocat/hello-world/issues"

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, clone)

	got := rec.Result().Header.Get("Link")
	if !regexp.MustCompile(`/github/repos/`).MatchString(got) {
		t.Errorf("the next page dropped the mounted prefix: %q", got)
	}
}
