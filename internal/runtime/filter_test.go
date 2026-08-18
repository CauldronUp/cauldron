package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// GitHub's issue listing answers with open issues unless you ask for
// otherwise. So does Alpaca's order listing, and so do enough others that a
// listing narrowing itself should be the assumption rather than the surprise.
//
// The failure it produces is the worst-shaped one there is: an issue is
// closed, the listing no longer contains it, and nothing errored. The list was
// correct. An emulator that returned everything would be helpful in exactly
// the direction that hides this, and Cauldron was: GitHub's own fixture holds
// a closed issue that the listing returned and the real API would not.

func githubList(t *testing.T, target string) []any {
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

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out []any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not a JSON array: %v %s", err, rec.Body.String())
	}

	return out
}

func statesOf(t *testing.T, records []any) []string {
	t.Helper()

	out := make([]string, 0, len(records))

	for _, record := range records {
		object, _ := record.(map[string]any)
		if object == nil {
			t.Fatalf("not an object: %v", record)
		}

		state, _ := object["state"].(string)
		out = append(out, state)
	}

	return out
}

func TestAListingAppliesItsDefaultFilter(t *testing.T) {
	states := statesOf(t, githubList(t, "/repos/octocat/hello-world/issues"))

	// The fixture holds one open and one closed issue in this repository.
	if len(states) != 1 || states[0] != "open" {
		t.Errorf("default listing should be the open issue alone, got %v", states)
	}
}

func TestASuppliedFilterOverridesTheDefault(t *testing.T) {
	states := statesOf(t, githubList(t, "/repos/octocat/hello-world/issues?state=closed"))

	if len(states) != 1 || states[0] != "closed" {
		t.Errorf("state=closed should be the closed issue alone, got %v", states)
	}
}

func TestTheEscapeValueTurnsTheFilterOff(t *testing.T) {
	states := statesOf(t, githubList(t, "/repos/octocat/hello-world/issues?state=all"))

	if len(states) != 2 {
		t.Errorf("state=all should be both issues, got %v", states)
	}
}

func TestAFilterDoesNotWidenTheScope(t *testing.T) {
	// The third fixture issue is open and belongs to another repository. A
	// filter joins the scope rather than replacing it, and getting that wrong
	// would leak one repository's issues into another's listing.
	states := statesOf(t, githubList(t, "/repos/octocat/hello-world/issues?state=all"))

	if len(states) != 2 {
		t.Errorf("another repository's issue must stay out of this listing, got %v", states)
	}
}
