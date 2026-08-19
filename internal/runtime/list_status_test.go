package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// status and empty_body were once read on creates alone. They were extended to
// the routes that answer with a record, and listings were left out of both, so
// a Recipe could declare that its provider answers a page with 206 and a
// Next-Range header and be quietly ignored.
//
// Heroku is why it matters. It pages with the Range header, answers 206
// Partial Content while there is more to fetch, and puts the resume point in a
// header rather than in the body. A 200 there means you have everything, so an
// emulator that always answered 200 taught a client the one thing it must not
// believe -- and a client that compares against 200 rejects every page but the
// last.

func herokuApps(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	r, err := recipe.Open("heroku")
	if err != nil {
		t.Fatalf("open heroku: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/apps", nil)
	req.Header.Set("Authorization", "Bearer HRKU-cauldron00000000000000000000000000")
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func TestAListingAnswersWithItsDeclaredStatus(t *testing.T) {
	if got := herokuApps(t).Code; got != http.StatusPartialContent {
		t.Errorf("status %d, want 206", got)
	}
}

func TestAListingSetsItsDeclaredHeaders(t *testing.T) {
	rec := herokuApps(t)

	if got := rec.Header().Get("Next-Range"); got == "" {
		t.Error("Next-Range was not set, so nothing tells a client where to resume")
	}

	if got := rec.Header().Get("Accept-Ranges"); got != "id, name" {
		t.Errorf("Accept-Ranges is %q", got)
	}
}
