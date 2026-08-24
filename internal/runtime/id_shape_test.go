package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Two ways of asking for an order that is not there, and Squarespace refuses
// them differently.
//
// GET /1.0/commerce/orders/{id} answers 404 "The requested Order was not
// found" for an identifier that could exist and does not, and 400 "The id is
// not in the expected format" for one that could not exist at all. Every
// Recipe here answered 404 to both, because an identifier had a shape the
// emulator minted with and no shape it checked against.
//
// The distinction is not decorative. A 404 is a fact about the account -- the
// order was deleted, or belongs to somebody else -- and retrying will not
// help. A 400 is a fact about the caller: an id from the wrong provider, a
// truncated string, an empty variable interpolated into the path. Code that
// branches on the two takes the wrong branch for its own bugs when the fake
// collapses them, and the test that proves the handler works asks for
// "nonexistent" -- which is exactly the identifier that behaves differently.
func TestAMisshapenIdentifierIsRefusedBeforeItIsLookedUp(t *testing.T) {
	r, err := recipe.Open("squarespace")
	if err != nil {
		t.Fatalf("open squarespace: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("site"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	ask := func(id string) int {
		t.Helper()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/1.0/commerce/orders/"+id, nil)
		req.Header.Set("Authorization", "Bearer cauldron-squarespace-fixture-key")
		req.Header.Set("User-Agent", "cauldron-test")
		s.ServeHTTP(rec, req)

		return rec.Code
	}

	// The one the fixture holds, so the shape check cannot be refusing
	// everything.
	if got := ask("585d498fdee9f31a60284a37"); got != http.StatusOK {
		t.Errorf("the seeded order = %d, want 200", got)
	}

	for _, id := range []string{
		// A well-formed ObjectId nobody minted.
		"000000000000000000000000",
		// And one the store might plausibly have held.
		"ffffffffffffffffffffffff",
	} {
		if got := ask(id); got != http.StatusNotFound {
			t.Errorf("%s = %d, want 404: an id that could exist and does not is an absence", id, got)
		}
	}

	for _, id := range []string{
		// The string a test writes when it wants a miss, which is the whole
		// problem: against Squarespace this is not a miss.
		"nonexistent",
		// Twenty-three characters rather than twenty-four.
		"585d498fdee9f31a60284a3",
		// Hexadecimal until the last character.
		"585d498fdee9f31a60284a3z",
		// Somebody else's identifier, pasted in.
		"cus_NffrFeUfNV2Hib",
		// And the upper-case spelling of a real one, which a database that
		// stores them lower-case does not recognise.
		"585D498FDEE9F31A60284A37",
	} {
		if got := ask(id); got != http.StatusBadRequest {
			t.Errorf("%s = %d, want 400: an id that could not exist is the caller's mistake", id, got)
		}
	}
}
