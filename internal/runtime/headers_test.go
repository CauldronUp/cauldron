package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Forgetting Notion-Version is the classic Notion integration bug: everything
// works against a hand-rolled fake and the first real call returns 400.

func notionSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("notion")
	if err != nil {
		t.Fatalf("open notion: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-workspace"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func notionRequest(t *testing.T, s *Sandbox, target string, version bool) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer secret_cauldron")

	if version {
		req.Header.Set("Notion-Version", "2022-06-28")
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func TestARequiredHeaderIsEnforced(t *testing.T) {
	s := notionSandbox(t)

	rec := notionRequest(t, s, "/v1/users", false)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400\n%s", rec.Code, rec.Body)
	}

	body := decode(t, rec)

	if body["code"] != "missing_version" {
		t.Errorf("code = %v, want the declared error", body["code"])
	}

	if body["object"] != "error" {
		t.Errorf("object = %v", body["object"])
	}
}

func TestTheSameRequestSucceedsWithTheHeader(t *testing.T) {
	s := notionSandbox(t)

	rec := notionRequest(t, s, "/v1/users", true)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200\n%s", rec.Code, rec.Body)
	}
}

// Auth is checked before the header, because a request with neither should
// report the problem a caller can actually act on first.
func TestAuthIsCheckedBeforeARequiredHeader(t *testing.T) {
	s := notionSandbox(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestListAndResourceCarryDifferentObjectTypes(t *testing.T) {
	s := notionSandbox(t)

	list := decode(t, notionRequest(t, s, "/v1/users", true))
	if list["object"] != "list" {
		t.Errorf("collection object = %v, want list", list["object"])
	}

	page := decode(t, notionRequest(t, s, "/v1/pages/11111111-1111-4111-8111-111111111111", true))
	if page["object"] != "page" {
		t.Errorf("page object = %v, want page", page["object"])
	}
}

// Providers disagree about which half of a basic credential carries the
// secret. Twilio puts the account SID in the username; Mailgun's username is
// the constant "api" and the key is the password. Checking the wrong half
// means a bad key is never rejected at all, so the failure path a test most
// wants to exercise silently returns 200.
func TestBasicAuthCanCheckThePasswordHalf(t *testing.T) {
	r, err := recipe.Open("mailgun")
	if err != nil {
		t.Fatalf("open mailgun: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-domain"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	call := func(user, password string) int {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, "/v3/mg.example.com/events", nil)
		req.SetBasicAuth(user, password)

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		return rec.Code
	}

	if code := call("api", "cauldron-mailgun-key"); code != http.StatusOK {
		t.Errorf("the right key should be accepted, got %d", code)
	}

	// The username is right and the key is not. Checking the username would
	// let this through.
	if code := call("api", "wrong-key"); code != http.StatusUnauthorized {
		t.Errorf("a bad key should be refused, got %d", code)
	}

	// And the halves are not interchangeable.
	if code := call("cauldron-mailgun-key", "api"); code != http.StatusUnauthorized {
		t.Errorf("the key in the username half should not be accepted, got %d", code)
	}
}
