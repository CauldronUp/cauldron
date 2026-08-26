package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A failure's declared Content-Type is what gets sent.
//
// The text style exists because some providers answer a failure with prose
// rather than JSON, and it set text/plain on the way out. Headers declared on
// the failure are written a few lines earlier, so a Recipe saying
// Content-Type: text/html had it replaced -- declared, applied, and then
// overwritten by a default, with nothing in the response to notice it by.
//
// Homebrew is the case that found it. GET /api/formula/does-not-exist.json is a
// 404 whose body is a full HTML page, from a path ending in .json, and the
// content type is the difference between a client rendering an error page and a
// client throwing a parse error. Serving that as text/plain would have been the
// quiet kind of wrong: the body would still not be JSON, so the case asserting
// the body would still pass.
//
// This is the third declared-and-ignored field closed in this collection, after
// a list key the collection name already supplied and auth.header on the bearer
// scheme. The shape is always the same -- the file says something, the runtime
// does something else, and no conformance case can tell them apart.
func TestAFailureKeepsTheContentTypeItDeclares(t *testing.T) {
	r, err := recipe.Open("homebrew")
	if err != nil {
		t.Fatalf("open homebrew: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("tap"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/formula/not-a-real-formula.json", nil)
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	if got, want := rec.Header().Get("Content-Type"), "text/html; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}

	if body := rec.Body.String(); len(body) < 10 || body[:10] != "<!DOCTYPE " {
		t.Errorf("body does not begin with a doctype: %.40q", body)
	}
}

// And a failure that declares nothing still gets the default, because the
// default is what makes the text style worth having: Trello answers plain text
// and a client calling .json() on it throws, which is the thing being modelled.
func TestATextFailureWithoutADeclarationIsStillPlainText(t *testing.T) {
	r, err := recipe.Open("homebrew")
	if err != nil {
		t.Fatalf("open homebrew: %v", err)
	}

	// The same Recipe with the declaration taken away.
	stripped := *r
	stripped.Errors = map[string]recipe.Error{}

	for name, e := range r.Errors {
		if name == "resource_missing" {
			e.Headers = nil
		}

		stripped.Errors[name] = e
	}

	s, err := New(&stripped, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("tap"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/formula/not-a-real-formula.json", nil)
	s.ServeHTTP(rec, req)

	if got, want := rec.Header().Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}
