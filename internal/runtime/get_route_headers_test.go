// A route's declared response headers, on a route that fetches one record.

package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Route headers were written from the create path and the listing path and
// nowhere else, so a get declaring them had them dropped in silence. The Recipe
// said one thing and the response carried another, and no conformance case
// existed that could tell them apart -- every Recipe using route headers until
// now declared them on a create, which is the shape SendGrid needed them for.
//
// Advice Slip is what found it. Every one of its responses carries two
// Cache-Control field lines that contradict each other -- max-age=3600 beside
// max-age=600, private, must-revalidate -- and the route that shows it is a get.
// Before this the header simply was not there.
func TestAGetRouteSetsTheHeadersItDeclares(t *testing.T) {
	r, err := recipe.Open("adviceslip")
	if err != nil {
		t.Fatalf("open adviceslip: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("slips"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/advice/1", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	got := rec.Header().Get("Cache-Control")

	if got == "" {
		t.Fatal("the get route declared a Cache-Control header and the response carries none")
	}

	if want := "max-age=3600, max-age=600, private, must-revalidate"; got != want {
		t.Errorf("Cache-Control = %q, want %q", got, want)
	}
}
