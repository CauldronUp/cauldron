package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A trailing slash is part of the path for the providers that declare one, and
// the router used to trim it from both sides.
//
// Sixty-one routes across five Recipes -- Sentry, PostHog, GoCardless Bank
// Account Data, Uploadcare and Saleor -- write their paths with a trailing
// slash, and every one of them is a Django or Django-shaped API where the URL
// pattern is anchored and the slash is required. Trimming it meant all
// sixty-one accepted the path without it, so the claim in the Recipe was
// decoration.
//
// The failure that hides is the quiet kind. Django answers a slash-less path
// with a 301 to the slash, and a client that follows a redirect after a POST
// commonly drops the body -- so the request arrives, empty, at the right URL,
// and is refused for a reason that has nothing to do with the redirect. It
// worked locally, against an emulator that was being helpful.
func TestATrailingSlashIsPartOfThePath(t *testing.T) {
	r, err := recipe.Open("saleor")
	if err != nil {
		t.Fatalf("open saleor: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	ask := func(path string) int {
		t.Helper()

		body := `{"query":"{ orders(first: 1) { edges { node { id } } } }"}`
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer cauldron.saleor.fixture.jwt.000000")
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		return rec.Code
	}

	// urls.py anchors the pattern as ^graphql/$.
	if got := ask("/graphql/"); got != http.StatusOK {
		t.Errorf("the declared path = %d, want 200", got)
	}

	// And the same path without it is not that route.
	if got := ask("/graphql"); got != http.StatusNotFound {
		t.Errorf("the path without its slash = %d, want 404", got)
	}
}

// The converse, so the rule is symmetrical rather than a special case for one
// spelling: a route declared without a slash does not answer one either.
//
// Every Recipe that predates this declares its paths without a trailing slash,
// so getting this direction wrong would have broken all of them at once --
// which is the reason it is asserted rather than assumed.
func TestAPathWithoutASlashDoesNotAnswerOne(t *testing.T) {
	r, err := recipe.Open("shopware")
	if err != nil {
		t.Fatalf("open shopware: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("catalogue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	ask := func(path string) int {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("sw-access-key", "SWSCCAULDRONFIXTUREKEY00")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		return rec.Code
	}

	if got := ask("/store-api/product"); got != http.StatusOK {
		t.Errorf("the declared path = %d, want 200", got)
	}

	if got := ask("/store-api/product/"); got != http.StatusNotFound {
		t.Errorf("the path with a slash it does not declare = %d, want 404", got)
	}
}
