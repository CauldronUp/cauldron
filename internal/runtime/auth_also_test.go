package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// One secret, several carriers.
//
// Kickbox accepts its key as a query parameter or a bearer header, Watchmode
// takes three channels, Clearbit takes a bearer or Basic with a blank password.
// Nine Recipes recorded that and served one channel -- which is the wrong
// direction to be wrong in, because a client using the other was refused by the
// fake and accepted by the provider. Code written against a fake that is
// stricter than the thing it stands in for ships broken.
func TestOneCredentialCanArriveThroughSeveralCarriers(t *testing.T) {
	s, err := New(severalCarriers(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// The primary: a query parameter.
	if got := ask(t, s, "/v1/things?key=the-one-key", ""); got != http.StatusOK {
		t.Errorf("the query parameter answered %d, want 200", got)
	}

	// And the alternative: the same key in a header.
	if got := ask(t, s, "/v1/things", "the-one-key"); got != http.StatusOK {
		t.Errorf("the header carrier answered %d, want 200", got)
	}

	// A wrong key through either is still wrong.
	if got := ask(t, s, "/v1/things?key=not-the-key", ""); got != http.StatusUnauthorized {
		t.Errorf("a wrong key in the query answered %d, want 401", got)
	}

	if got := ask(t, s, "/v1/things", "not-the-key"); got != http.StatusUnauthorized {
		t.Errorf("a wrong key in the header answered %d, want 401", got)
	}

	// And nothing anywhere is still nothing.
	if got := ask(t, s, "/v1/things", ""); got != http.StatusUnauthorized {
		t.Errorf("no credential at all answered %d, want 401", got)
	}
}

// The primary's verdict wins when the primary was used.
//
// Getting the primary wrong is a different mistake from using another channel,
// so a wrong key in the primary must not be re-judged as absent because some
// other carrier held nothing.
func TestGettingThePrimaryWrongIsNotReadAsUsingAnotherChannel(t *testing.T) {
	r := severalCarriers()
	r.Auth.AbsentError = "nothing_sent"
	r.Auth.RejectedError = "wrong_key"
	r.Errors["nothing_sent"] = recipe.Error{Status: 401, Message: "No key."}
	r.Errors["wrong_key"] = recipe.Error{Status: 403, Message: "Bad key."}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got := ask(t, s, "/v1/things?key=not-the-key", ""); got != http.StatusForbidden {
		t.Errorf("a wrong key in the primary answered %d, want the rejected sentence's 403", got)
	}

	if got := ask(t, s, "/v1/things", ""); got != http.StatusUnauthorized {
		t.Errorf("nothing anywhere answered %d, want the absent sentence's 401", got)
	}
}

// A wrong credential in an alternative is reported as wrong, not as absent.
//
// The primary holds nothing, which says only that this caller used another
// channel. Answering "no credential" to somebody who sent one describes a
// request nobody made.
func TestAWrongCredentialInAnAlternativeIsNotReportedAsAbsent(t *testing.T) {
	r := severalCarriers()
	r.Auth.AbsentError = "nothing_sent"
	r.Auth.RejectedError = "wrong_key"
	r.Errors["nothing_sent"] = recipe.Error{Status: 401, Message: "No key."}
	r.Errors["wrong_key"] = recipe.Error{Status: 403, Message: "Bad key."}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got := ask(t, s, "/v1/things", "not-the-key"); got != http.StatusForbidden {
		t.Errorf("a wrong key in the alternative answered %d, want the rejected sentence's 403", got)
	}
}

// An alternative naming no carrier is the primary written twice, and a Recipe
// saying it is claiming a second channel it does not have.
func TestAnAlternativeMustNameACarrier(t *testing.T) {
	r := severalCarriers()
	r.Auth.Also = []recipe.Auth{{Keys: []string{"another"}}}

	if err := r.Validate(); err == nil {
		t.Fatal("an alternative naming no carrier validated")
	}
}

// And alternatives do not nest, because the inner ones would never be reached.
func TestAlternativesDoNotNest(t *testing.T) {
	r := severalCarriers()
	r.Auth.Also = []recipe.Auth{{
		Scheme: "header",
		Header: "X-Key",
		Also:   []recipe.Auth{{Scheme: "query", Param: "k"}},
	}}

	if err := r.Validate(); err == nil {
		t.Fatal("a nested alternative validated")
	}
}

// An alternative can take the bare secret where the primary takes a prefixed
// one, which needs a way to say "no prefix" rather than inherit.
//
// Doppler found this by being wired: its bearer token also works as a Basic
// username with a blank password, and the Basic alternative inherited the
// primary's "Bearer " and refused every key legitimately sent that way. A fake
// stricter than its provider is exactly what this field exists to prevent, so
// reintroducing it one level down was worse than not having the field.
func TestAnAlternativeCanClearAnInheritedPrefix(t *testing.T) {
	r := severalCarriers()
	r.Auth = recipe.Auth{
		Scheme: "bearer",
		Prefix: "Bearer ",
		Keys:   []string{"the-one-key"},
		Also: []recipe.Auth{{
			Scheme: "header",
			Header: "X-Api-Key",
			Prefix: "-",
		}},
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if got := ask(t, s, "/v1/things", "the-one-key"); got != http.StatusOK {
		t.Errorf("the bare secret through the alternative answered %d, want 200", got)
	}

	// And the primary still wants its prefix.
	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer the-one-key")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("the prefixed primary answered %d, want 200", w.Code)
	}
}

// Two carriers can share one header, and then the one that could read the
// value knows more about it than the one that could not.
//
// Clearbit takes a bearer or Basic with a blank password, both in
// Authorization. A bearer value is unreadable to the Basic check and merely
// wrong to the bearer one, and reporting "malformed" describes the carrier's
// confusion rather than the caller's mistake.
func TestTheCarrierThatCouldReadItDecides(t *testing.T) {
	r := severalCarriers()
	r.Auth = recipe.Auth{
		Scheme:        "basic",
		Keys:          []string{"the-one-key"},
		RejectedError: "wrong_key",
		Also: []recipe.Auth{{
			Scheme: "bearer",
			Header: "Authorization",
			Prefix: "Bearer ",
		}},
	}
	r.Errors["wrong_key"] = recipe.Error{Status: 403, Message: "Bad key."}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// A good key through the bearer channel is accepted.
	req := httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer the-one-key")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("a good key through the bearer channel answered %d, want 200", w.Code)
	}

	// A wrong one is rejected, not called malformed by the Basic check that
	// could not read it.
	req = httptest.NewRequest(http.MethodGet, "/v1/things", nil)
	req.Header.Set("Authorization", "Bearer not-the-key")

	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("a wrong bearer answered %d, want the rejected sentence's 403", w.Code)
	}
}

func ask(t *testing.T, s *Sandbox, target, header string) int {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	if header != "" {
		req.Header.Set("X-Api-Key", header)
	}

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	return w.Code
}

// A Recipe whose key travels in a query parameter or a header, as several do.
func severalCarriers() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "severalcarriers",
		Capability: "identity",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme: "query",
			Param:  "key",
			Keys:   []string{"the-one-key"},
			Also: []recipe.Auth{{
				Scheme: "header",
				Header: "X-Api-Key",
			}},
		},
		Errors: map[string]recipe.Error{},
		Resources: map[string]recipe.Resource{
			"thing": {ID: recipe.ID{Style: "opaque"}, Fields: map[string]recipe.Field{"id": {Type: "string"}}},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
	}
}
