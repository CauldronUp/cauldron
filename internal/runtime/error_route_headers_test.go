// The headers a route declares, on the route that answers a failure.

package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Route headers were written on the paths that answer with a record and
// nowhere else, so a route whose whole job is to answer a failure could
// declare headers and never send them -- silently, with nothing in the Recipe
// saying the declaration was inert.
//
// The iTunes Search API is what found it. Every JSON response there is
// announced as text/javascript, which is a real trap on its own, but a path the
// application never sees is answered by the web server instead: Apache's own
// page, as text/html; charset=iso-8859-1. Three content types on one host, and
// the only one that is not UTF-8 is the one a client is least likely to be
// ready for. A Recipe that could not declare a content type on an error route
// could not pin that at all.
func TestARouteAnsweringAFailureSendsTheHeadersItDeclares(t *testing.T) {
	r, err := recipe.Open("itunes")
	if err != nil {
		t.Fatalf("open itunes: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nosuchpath", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404: %s", rec.Code, rec.Body.String())
	}

	if got, want := rec.Header().Get("Content-Type"), "text/html; charset=iso-8859-1"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}

	// And the other failure keeps its own, which is the JavaScript one every
	// JSON response on this host carries.
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/search-bad-country", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400: %s", rec.Code, rec.Body.String())
	}

	if got, want := rec.Header().Get("Content-Type"), "text/javascript; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}
