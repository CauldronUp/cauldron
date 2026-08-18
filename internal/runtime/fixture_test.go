package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Reset promises a sandbox indistinguishable from a freshly booted one, and
// that promise is what makes it safe to reuse between tests.
//
// It did not hold for any fixture value that was a nested object or list, and
// the reason was two shallow copies in a row: Seed built each record one level
// deep from the Recipe's own embedded fixture, and the store cloned it one
// level deep again. So a stored record's nested object *was* the Recipe's
// object, the Recipe is process-global, and Reset re-seeded from the same map
// somebody had already written into.
//
// Jira is the Recipe to test it on because its fixtures nest three deep, and
// because it was the one the audit reproduced it with.
func TestAFixtureSurvivesBeingMutatedThroughTheSandbox(t *testing.T) {
	seeded := func(t *testing.T) *Sandbox {
		t.Helper()

		r, err := recipe.Open("jira")
		if err != nil {
			t.Fatalf("open jira: %v", err)
		}

		s, err := New(r, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		if err := s.Seed("platform-project"); err != nil {
			t.Fatalf("seed: %v", err)
		}

		return s
	}

	first := seeded(t)

	held, err := first.store.Get("issue", "10001")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	// A nested object, reached the way a handler reaches one.
	custom, ok := held["customfield_10041"].(map[string]any)
	if !ok {
		t.Fatalf("customfield_10041 is %T, expected a nested object to test with", held["customfield_10041"])
	}

	custom["value"] = "poisoned"
	custom["__probe__"] = true

	if err := first.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}

	// The same sandbox, re-seeded, must not carry it.
	after, err := first.store.Get("issue", "10001")
	if err != nil {
		t.Fatalf("get after reset: %v", err)
	}

	if value, _ := after["customfield_10041"].(map[string]any); value["value"] != "Critical" || value["__probe__"] != nil {
		t.Errorf("a reset sandbox carried the mutation: %v", value)
	}

	// And a brand new sandbox in the same process must not either, which is
	// the part that made this a whole-suite problem rather than a local one.
	second := seeded(t)

	fresh, err := second.store.Get("issue", "10001")
	if err != nil {
		t.Fatalf("get from a second sandbox: %v", err)
	}

	if value, _ := fresh["customfield_10041"].(map[string]any); value["value"] != "Critical" || value["__probe__"] != nil {
		t.Errorf("a freshly built sandbox carried a mutation made through another one: %v", value)
	}
}

// The same leak by the route a caller actually takes: a response body is
// shaped from a record the store handed out, so a handler holding one must not
// be holding the store's own nested maps.
func TestAResponseCannotBeUsedToReachIntoTheStore(t *testing.T) {
	r, err := recipe.Open("jira")
	if err != nil {
		t.Fatalf("open jira: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("platform-project"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	get := func() map[string]any {
		req := httptest.NewRequest(http.MethodGet, "/rest/api/3/issue/10001", nil)
		req.SetBasicAuth("work@cauldron.test", "cauldron-jira-api-token-0000")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
		}

		return decode(t, rec)
	}

	before := get()

	fields, _ := before["fields"].(map[string]any)
	fields["summary"] = "poisoned"

	after := get()
	if got, _ := after["fields"].(map[string]any); got["summary"] != "Rotate the signing keys" {
		t.Errorf("writing into one response changed the next one: %v", got["summary"])
	}
}

// The same collision through the API, on one of the eight recipes it was live
// on. A numeric counter starting at zero and a fixture pinning id "1" meant
// the first issue the API created replaced the first issue the fixture
// seeded: the collection did not grow, the seeded title was gone, and the
// response was a 201.
func TestCreatingThroughTheApiCannotDestroyAFixtureRecord(t *testing.T) {
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

	before, err := s.store.ListWhere("issue", nil, "", 100)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(before.Records) < 2 {
		t.Fatalf("expected the fixture to seed at least two issues, got %d", len(before.Records))
	}

	firstTitle := before.Records[0]["title"]

	req := httptest.NewRequest(http.MethodPost, "/repos/cauldron/example/issues",
		strings.NewReader(`{"title":"a brand new issue from the API"}`))
	req.Header.Set("Authorization", "Bearer ghp_cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated && rec.Code != http.StatusOK {
		t.Fatalf("create = %d\n%s", rec.Code, rec.Body)
	}

	after, err := s.store.ListWhere("issue", nil, "", 100)
	if err != nil {
		t.Fatalf("list after: %v", err)
	}

	if len(after.Records) != len(before.Records)+1 {
		t.Errorf("the collection went from %d to %d records after one create", len(before.Records), len(after.Records))
	}

	if after.Records[0]["title"] != firstTitle {
		t.Errorf("creating an issue overwrote the first seeded one: %v", after.Records[0]["title"])
	}
}
