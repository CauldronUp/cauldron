package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A wrong method that is not answered with a 405 gets no Allow header.
//
// Allow belongs on a 405 and nowhere else, which is what RFC 9110 requires of
// it. Twenty-two Recipes say their provider has no 405 at all -- Turso answers
// the same route_not_found a missing path gets, Koyeb answers url_not_found,
// ClickHouse answers Express's bare "Cannot PUT /x", Workable answers the
// identical 404 an unrouted path does -- and computing an Allow to sit beside
// those invented a header no such response carries.
//
// That is the hardest kind of infidelity to notice, because nothing fails: the
// status is right, the body is right, and the client learns one extra thing
// that is true of nobody.
func TestAWrongMethodAnsweredWithoutA405GetsNoAllowHeader(t *testing.T) {
	r, err := recipe.Open("workable")
	if err != nil {
		t.Fatalf("open workable: %v", err)
	}

	declared, ok := r.Errors["method_not_allowed"]
	if !ok || declared.Status == http.StatusMethodNotAllowed {
		t.Skip("workable now answers a wrong method with a 405; pick another Recipe")
	}

	var path string

	for _, route := range r.Routes {
		if route.Method == http.MethodGet && !containsBrace(route.Path) {
			path = route.Path

			break
		}
	}

	if path == "" {
		t.Fatal("workable has no parameterless GET to send a wrong method to")
	}

	w := serveMethod(t, r, http.MethodPatch, path)

	// Workable has no method-not-allowed concept on the wire at all: a wrong
	// method answers the identical 404 an unrouted path does.
	if got := w.Header().Get("Allow"); got != "" {
		t.Errorf("Allow: %q was added to a failure the Recipe describes itself", got)
	}
}

// A Recipe whose provider does send a 405 keeps the computed header, and wants
// it.
//
// This is the half that makes the rule a rule rather than a removal. Airbrake,
// Api2Pdf and Papertrail all answer a real 405 with a real Allow, and all three
// assert its value in their own cases -- so the first version of this change,
// which suppressed the header for any Recipe that named the failure at all,
// broke every one of them.
func TestARecipeWhoseProviderReallySendsA405KeepsTheComputedAllow(t *testing.T) {
	r, err := recipe.Open("api2pdf")
	if err != nil {
		t.Fatalf("open api2pdf: %v", err)
	}

	declared, ok := r.Errors["method_not_allowed"]
	if ok && declared.Status != 0 && declared.Status != http.StatusMethodNotAllowed {
		t.Skip("api2pdf no longer answers a wrong method with a 405; pick another Recipe")
	}

	var path string

	for _, route := range r.Routes {
		if route.Method == http.MethodPost && !containsBrace(route.Path) {
			path = route.Path

			break
		}
	}

	if path == "" {
		t.Fatal("api2pdf has no parameterless POST to send a wrong method to")
	}

	w := serveMethod(t, r, http.MethodPatch, path)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH %s answered %d, not the 405 this test needs to observe", path, w.Code)
	}

	if w.Header().Get("Allow") == "" {
		t.Error("a provider that really sends a 405 lost the Allow header this runtime computes for it")
	}
}

func serveMethod(t *testing.T, r *recipe.Recipe, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(method, path, nil))

	return w
}
