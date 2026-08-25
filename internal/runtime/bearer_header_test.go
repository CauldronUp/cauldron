package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A bearer credential is read from whichever header the Recipe names.
//
// It used to be read from Authorization and nowhere else, which made auth.header
// inert on the bearer scheme -- 127 of the 130 Recipes using it declare that
// field, every one of them saying "Authorization", so nothing was ever wrong on
// the wire and nothing ever could be.
//
// That is the shape this project keeps finding and keeps closing: a declaration
// the runtime does not read. It is worse than a missing feature, because the
// file says one thing, the emulator does another, and no conformance case can
// tell them apart -- a mutation renaming the header survived every case in the
// Recipe that found this.
//
// The two schemes now share a branch, because they differ only in convention: a
// bearer token usually carries a prefix and a header credential usually does
// not, and the prefix is applied to both afterwards.
func TestABearerCredentialIsReadFromTheHeaderTheRecipeNames(t *testing.T) {
	r, err := recipe.Open("grafana")
	if err != nil {
		t.Fatalf("open grafana: %v", err)
	}

	if r.Auth.Scheme != "bearer" {
		t.Fatalf("grafana no longer uses the bearer scheme; pick another Recipe")
	}

	const token = "glsa_CauldronFixtureServiceAccountToken_000000"

	// The Recipe as written: Authorization, by default.
	assertAccepts(t, r, "Authorization", "Bearer "+token, http.StatusOK)
	assertAccepts(t, r, "X-Grafana-Token", "Bearer "+token, http.StatusUnauthorized)

	// The same Recipe, naming a different header. Before this was honoured,
	// both of these answered 200 and 401 the other way round -- the emulator
	// went on reading Authorization whatever the file said.
	moved := *r
	moved.Auth.Header = "X-Grafana-Token"

	assertAccepts(t, &moved, "X-Grafana-Token", "Bearer "+token, http.StatusOK)
	assertAccepts(t, &moved, "Authorization", "Bearer "+token, http.StatusUnauthorized)
}

func assertAccepts(t *testing.T, r *recipe.Recipe, header, value string, want int) {
	t.Helper()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("instance"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	req.Header.Set(header, value)
	s.ServeHTTP(rec, req)

	if rec.Code != want {
		t.Errorf("credential in %s: status = %d, want %d (body %s)",
			header, rec.Code, want, rec.Body.String())
	}
}
