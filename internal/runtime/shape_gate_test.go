package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A credential can be the wrong shape or the wrong value, and several
// providers say so differently.
//
// Airtable answers "Invalid authentication token" to a wrong token of the
// right shape and the same thing it says to a missing header to one of any
// other length. Modern Treasury answers invalid_key to a password carrying
// neither its test- nor its live- prefix. Opsgenie answers "Key format is not
// valid!" to a key that is not UUID-shaped and "Could not authenticate" to a
// UUID nobody issued. Paddle answers "Authentication header included, but
// incorrectly formatted." to a token in a shape it does not recognise.
//
// All four wrote the same paragraph about why they could not model it. Pattern
// is the acceptance test, so declaring one authenticates every correctly
// shaped string -- a fake laxer than the provider, to fix a fake that was
// merely less specific. Shape sits in front of the keys instead.
func TestACredentialOfTheWrongShapeIsMalformedRatherThanRejected(t *testing.T) {
	s, err := New(shaped(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// The right shape and the right value.
	if got, _ := shapedAsk(t, s, "test-the-one-key"); got != http.StatusOK {
		t.Errorf("the key this Recipe publishes answered %d, want 200", got)
	}

	// The right shape, a value nobody issued.
	got, body := shapedAsk(t, s, "test-nobody-issued-this")
	if got != http.StatusUnauthorized {
		t.Errorf("a well-formed wrong key answered %d, want 401", got)
	}

	if !strings.Contains(body, "Invalid API key.") {
		t.Errorf("a well-formed wrong key got the wrong sentence: %s", body)
	}

	// The wrong shape, which the provider does not read as a key at all.
	got, body = shapedAsk(t, s, "nothing-like-a-key")
	if got != http.StatusBadRequest {
		t.Errorf("a wrongly shaped key answered %d, want the malformed sentence's 400", got)
	}

	if !strings.Contains(body, "Key format is not valid.") {
		t.Errorf("a wrongly shaped key got the wrong sentence: %s", body)
	}
}

// And a shape does not authenticate on its own, which is the whole difference
// from a pattern: matching it means the value is worth comparing, not that it
// was right.
func TestAShapeIsNotAnAcceptance(t *testing.T) {
	r := shaped()
	r.Auth.Keys = []string{"test-the-one-key"}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got, _ := shapedAsk(t, s, "test-some-other-string-entirely"); got == http.StatusOK {
		t.Error("a string matching the shape authenticated without matching any key")
	}
}

// A credential the carrier could not read and one the shape ruled out are two
// different failures, and Opsgenie answers them differently.
//
// A header that does not begin "GenieKey ", and a bare "GenieKey" with nothing
// after it, both answer 401 "Could not authenticate" -- the same as sending no
// header at all. A GenieKey-prefixed value that is not a UUID answers 422 "Key
// format is not valid!". All three struck live on 2026-09-05. Folding the
// shape failure into the malformed verdict would have made one of those two
// sentences unreachable.
func TestAShapeFailureAndAnUnreadableCarrierAreDifferentSentences(t *testing.T) {
	r := shaped()
	r.Auth.ShapeError = "bad_shape"
	r.Auth.MalformedError = "unreadable"
	r.Errors["unreadable"] = recipe.Error{Status: 401, Message: "Could not authenticate."}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// The prefix is wrong, so nothing was read.
	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Token test-the-one-key")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if got := w.Body.String(); !strings.Contains(got, "Could not authenticate.") {
		t.Errorf("an unreadable carrier got the shape sentence: %s", got)
	}

	// The prefix is right and the shape is not.
	if _, body := shapedAsk(t, s, "nothing-like-a-key"); !strings.Contains(body, "Key format is not valid.") {
		t.Errorf("a wrongly shaped credential got the unreadable sentence: %s", body)
	}
}

// A Recipe that declares a shape and no sentence for failing it falls back to
// the malformed one, which is the nearest thing it has said.
func TestAShapeWithNoSentenceOfItsOwnFallsBackToTheMalformedOne(t *testing.T) {
	r := shaped()
	r.Auth.ShapeError = ""
	r.Auth.MalformedError = "bad_shape"

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if _, body := shapedAsk(t, s, "nothing-like-a-key"); !strings.Contains(body, "Key format is not valid.") {
		t.Errorf("a shape failure with no sentence of its own got: %s", body)
	}
}

func shapedAsk(t *testing.T, s *Sandbox, key string) (int, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer "+key)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}

func shaped() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "shaped",
		Capability: "identity",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme:         "bearer",
			Prefix:         "Bearer ",
			Shape:          "^test-[a-z-]+$",
			Keys:           []string{"test-the-one-key"},
			MalformedError: "bad_shape",
			RejectedError:  "bad_key",
		},
		Errors: map[string]recipe.Error{
			"bad_shape": {Status: 400, Message: "Key format is not valid."},
			"bad_key":   {Status: 401, Message: "Invalid API key."},
		},
		Resources: map[string]recipe.Resource{
			"thing": {ID: recipe.ID{Style: "opaque"}, Fields: map[string]recipe.Field{"id": {Type: "string"}}},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
	}
}
