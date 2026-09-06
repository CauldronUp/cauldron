package runtime

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Django REST Framework's page shape is {count, next, previous, results}, and
// `previous` is not optional in it: null on the first page, a URL after that.
// Two Recipes wrote the same paragraph rather than serve it -- RAWG's said in
// so many words that its own "here is where you came from" was "real,
// documented, and not expressible in this format today".
//
// The absence is not cosmetic. A client deciding whether it is on the first
// page by testing whether the key is there takes the wrong branch every time
// against an emulator that omits it, and takes it silently.

func codecovPage(t *testing.T, target string) map[string]any {
	t.Helper()

	r, err := recipe.Open("codecov")
	if err != nil {
		t.Fatalf("open codecov: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("codecov-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("%s answered %d, want 200", target, rec.Code)
	}

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("reading the body: %v", err)
	}

	var page map[string]any
	if err := json.Unmarshal(body, &page); err != nil {
		t.Fatalf("decoding %s: %v", body, err)
	}

	return page
}

func TestThePreviousPageIsNullOnTheFirstPageRatherThanAbsent(t *testing.T) {
	t.Parallel()

	page := codecovPage(t, "/api/v2/github/codecov/repos/?page_size=1")

	value, present := page["previous"]
	if !present {
		t.Fatal("the first page has no previous key at all, which is the branch a client takes to decide it is at the start")
	}

	if value != nil {
		t.Errorf("previous on the first page is %v, want null", value)
	}

	// And next is a URL on the same response, so the two are being decided
	// independently rather than from one flag.
	next, _ := page["next"].(string)
	if !strings.Contains(next, "page=2") {
		t.Errorf("next on the first page is %q, want a URL carrying page=2", page["next"])
	}
}

func TestThePreviousPageIsAURLOnceThereIsOne(t *testing.T) {
	t.Parallel()

	page := codecovPage(t, "/api/v2/github/codecov/repos/?page_size=1&page=2")

	previous, ok := page["previous"].(string)
	if !ok {
		t.Fatalf("previous on page two is %v, want a URL", page["previous"])
	}

	if !strings.Contains(previous, "page=1") {
		t.Errorf("previous on page two is %q, want a URL carrying page=1", previous)
	}

	// The page size has to survive, or following it backwards lands on a
	// differently-sized page one and the caller sees records twice.
	if !strings.Contains(previous, "page_size=1") {
		t.Errorf("previous dropped the page size: %q", previous)
	}
}

// The last page still carries both keys, with the values the other way round.
// A client that stops when next is null and then walks back is the reason both
// have to be independent rather than two readings of one flag.
func TestTheLastPageCarriesAPreviousAndANullNext(t *testing.T) {
	t.Parallel()

	page := codecovPage(t, "/api/v2/github/codecov/repos/?page_size=1&page=3")

	if next, present := page["next"]; !present || next != nil {
		t.Errorf("next on the last page is %v (present=%v), want null and present", next, present)
	}

	if previous, ok := page["previous"].(string); !ok || !strings.Contains(previous, "page=2") {
		t.Errorf("previous on the last page is %v, want a URL carrying page=2", page["previous"])
	}
}

// A Recipe that has not named the field must not have one invented. Serving a
// previous key against a provider that sends none teaches a client to branch
// on something that will not be there.
func TestNoPreviousFieldUnlessTheRecipeNamesOne(t *testing.T) {
	t.Parallel()

	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	if r.Responses.List.PrevField != "" {
		t.Skip("github now names a previous field, so it cannot stand for the silent case")
	}

	body, err := io.ReadAll(issuesResponse(t, "/repos/octocat/hello-world/issues?state=all&per_page=1").Body)
	if err != nil {
		t.Fatalf("reading the body: %v", err)
	}

	if strings.Contains(string(body), `"previous"`) {
		t.Errorf("a Recipe naming no previous field sent one: %s", body)
	}
}
