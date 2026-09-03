package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route can be open to a caller who presents nothing and closed to one who
// presents something wrong.
//
// football-data's competitions list does exactly that: it answers with no
// credential at all and refuses a wrong token. The exemption as first written
// was all-or-nothing, so modelling that route meant waving the wrong token
// through -- which teaches a caller their bad credential is fine on the one
// route they are most likely to try first.
//
// FireHydrant's route is the other reading and it was verified rather than
// assumed: byte-identical with a junk bearer attached and with no header at
// all. Both are real, so both are expressible.
func TestARouteCanBeOpenToNobodyAndClosedToAWrongCredential(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{WhenAbsent: true}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status, body := askFor(t, s, "/v1/runtimes"); status != http.StatusOK {
		t.Errorf("presenting nothing answered %d, want 200: %s", status, body)
	}

	// The half the blunt exemption gets wrong.
	if status := askWith(t, s, "/v1/runtimes", "Bearer not-a-key"); status != http.StatusUnauthorized {
		t.Errorf("a wrong credential answered %d, want the 401 the rest of the host gives", status)
	}

	// And a real key still works, which is what keeps this from being an
	// inverted gate rather than an exemption.
	// The Recipe declares no prefix, so the header value is the credential.
	if status := askWith(t, s, "/v1/runtimes", "a-key-the-recipe-holds"); status != http.StatusOK {
		t.Errorf("a key the Recipe holds answered %d, want 200", status)
	}
}

// The unconditional mode still ignores the credential entirely, which is the
// behaviour four shipped Recipes already rely on.
func TestTheUnconditionalModeStillIgnoresAWrongCredential(t *testing.T) {
	r := gatedAndNot()
	r.Routes[1].Public = recipe.PublicMode{Always: true}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if status := askWith(t, s, "/v1/runtimes", "Bearer not-a-key"); status != http.StatusOK {
		t.Errorf("the unconditional exemption refused a wrong credential: %d", status)
	}
}

func askWith(t *testing.T, s *Sandbox, path, credential string) int {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", credential)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	return w.Code
}
