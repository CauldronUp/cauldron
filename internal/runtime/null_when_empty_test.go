// An empty collection sent as the JSON null literal.

package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A Go handler that returns a nil slice marshals it as null, not as []. Pirsch
// does exactly that: a request carrying no credential at all answers 200 and
// the four bytes `null`, where a request carrying a good one answers an array.
//
// The difference is invisible in any test with data in it and certain in
// production the first time an account is empty: `response.ok` is true,
// `.json()` succeeds and hands back null, and the for-of on the next line
// throws. A client reaching `.length` instead gets undefined and reports no
// domains.
//
// omit_when_empty is the neighbouring case and not the same one -- a missing
// key reads as undefined, a null reads as null, and the two break different
// code -- so this is its own declaration rather than a mode of that one.
//
// A shipped Recipe on purpose: the shape is Pirsch's, not a fixture's.
func TestAnEmptyCollectionCanBeTheNullLiteral(t *testing.T) {
	r, err := recipe.Open("pirsch")
	if err != nil {
		t.Fatalf("open pirsch: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("empty"); err != nil {
		t.Fatalf("seed empty: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/domain", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	if got := strings.TrimSpace(rec.Body.String()); got != "null" {
		t.Errorf("empty body = %q, want null", got)
	}

	// And a seeded account still answers an array, because the declaration is
	// about emptiness rather than about the endpoint.
	if err := s.Seed("account"); err != nil {
		t.Fatalf("seed account: %v", err)
	}

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/domain", nil)
	req.Header.Set("Authorization", "Bearer cauldr0n.p1rsch.f1xturet0ken")
	s.ServeHTTP(rec, req)

	if body := strings.TrimSpace(rec.Body.String()); !strings.HasPrefix(body, "[") {
		t.Errorf("seeded body = %.60s, want an array", body)
	}
}
