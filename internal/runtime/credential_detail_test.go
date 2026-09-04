package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A provider that quotes the credential back can now say so.
//
// The error table holds one static sentence per named failure, and the
// auth-verdict call site passed no detail at all -- so a provider whose refusal
// echoes what was sent had to serve a fixed placeholder and explain in prose
// that the real one varies. NeverBounce answers "Invalid API key 'whatever you
// sent'", and New Relic echoes a wrong key into its body for the rejected
// verdict and not the absent one.
//
// The value is the one the credential check judged, after any prefix is
// stripped, because that is the part a provider quotes: the token, not the
// word Bearer in front of it.
func TestARefusalCanQuoteTheCredentialItRefused(t *testing.T) {
	r := quoting()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-key")
	r0, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	r0.ServeHTTP(w, req)

	if got := w.Body.String(); !strings.Contains(got, "'not-a-real-key'") {
		t.Errorf("the refusal did not quote the credential it refused: %s", got)
	}

	// The prefix is not part of what a provider quotes.
	if strings.Contains(w.Body.String(), "Bearer") {
		t.Errorf("the refusal quoted the scheme word as well as the token: %s", w.Body.String())
	}
}

// An absent credential quotes nothing, because there is nothing to quote -- and
// that is a different sentence, which is the whole reason the verdicts are
// separate.
func TestAnAbsentCredentialQuotesNothing(t *testing.T) {
	r := quoting()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/things", nil))

	if got := w.Body.String(); !strings.Contains(got, "No API key") {
		t.Errorf("an absent credential got the wrong sentence: %s", got)
	}
}

// A message that does not ask for the credential is unaffected, which is every
// Recipe that has not opted in.
func TestAMessageWithoutTheMarkerIsUnchanged(t *testing.T) {
	r := quoting()
	r.Errors["bad_key"] = recipe.Error{Status: 401, Message: "Invalid API key."}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-key")
	s.ServeHTTP(w, req)

	if got := w.Body.String(); strings.Contains(got, "not-a-real-key") {
		t.Errorf("a message with no marker leaked the credential anyway: %s", got)
	}
}

// A field may carry the marker too, and two requests must not see each other's.
//
// New Relic puts the rejected key back under error.api_key and leaves its
// sentence alone, so substituting only the message would have left that Recipe
// arming a fault to show a literal. Fields belong to the Recipe, which every
// request through one Sandbox shares -- so the substitution copies rather than
// writing through, and this is what proves it: the second request must not see
// the first one's credential.
func TestAFieldCarryingTheMarkerIsSubstitutedPerRequest(t *testing.T) {
	r := quoting()
	r.Errors["bad_key"] = recipe.Error{
		Status:  401,
		Message: "Invalid API key.",
		Fields:  map[string]any{"api_key": "{detail}"},
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	for _, sent := range []string{"first-wrong-key", "second-wrong-key"} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
		req.Header.Set("Authorization", "Bearer "+sent)
		s.ServeHTTP(w, req)

		if got := w.Body.String(); !strings.Contains(got, sent) {
			t.Errorf("sent %s and the field did not carry it: %s", sent, got)
		}
	}

	// And the Recipe's own map is untouched, which is the half a passing
	// response would not show.
	if got := r.Errors["bad_key"].Fields["api_key"]; got != "{detail}" {
		t.Errorf("the Recipe's shared field was overwritten with %v", got)
	}
}

// A Recipe whose rejection sentence quotes what it was given.
func quoting() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "quoting",
		Capability: "identity",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme:        "bearer",
			Prefix:        "Bearer ",
			Keys:          []string{"a-key-the-recipe-holds"},
			AbsentError:   "no_key",
			RejectedError: "bad_key",
		},
		Errors: map[string]recipe.Error{
			"no_key":  {Status: 401, Message: "No API key was provided."},
			"bad_key": {Status: 401, Message: "Invalid API key '{detail}'."},
		},
		Resources: map[string]recipe.Resource{
			"thing": {
				ID:     recipe.ID{Style: "opaque"},
				Fields: map[string]recipe.Field{"id": {Type: "string"}},
			},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
	}
}
