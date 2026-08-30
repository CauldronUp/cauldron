// A listing that is a map with nothing wrapping it.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// CoinGecko answers /simple/price with {"bitcoin": {...}, "ethereum": {...}} at
// the top level. The map style already existed for Pusher, but it always wraps:
// the keyed object goes under the resource's collection name, or under whatever
// key the Recipe declares. There was no way to say "the map is the response",
// so a Recipe either invented a wrapper the provider does not send or served an
// array the provider does not send either.
//
// "-" is the convention id.field and message_field already use, where it means
// the provider does not send it, and it means the same thing here.
//
// What the wrapper hides is worth stating: with one, a client writes
// body.prices[id]; without one it writes body[id], and every key at the top
// level is a coin id rather than a property. A client that checks for
// body.error before reading is checking for a coin called "error".
func TestAListingCanBeAMapWithNoWrapper(t *testing.T) {
	r, err := recipe.Open("coingecko")
	if err != nil {
		t.Fatalf("open coingecko: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("market"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v3/simple/price?ids=bitcoin", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("the body is not an object: %v -- %.200s", err, rec.Body.String())
	}

	// The coin id is a top-level key, not a property of one.
	bitcoin, ok := body["bitcoin"].(map[string]any)
	if !ok {
		t.Fatalf("bitcoin is not at the top level: %.200s", rec.Body.String())
	}

	// Nothing wraps it, so there is no collection name to find.
	for _, wrapper := range []string{"prices", "price", "data", "result"} {
		if _, found := body[wrapper]; found {
			t.Errorf("the map is wrapped in %q and the provider sends no wrapper: %.200s", wrapper, rec.Body.String())
		}
	}

	// The key is the identifier, so it does not repeat inside the value.
	if _, repeated := bitcoin["id"]; repeated {
		t.Errorf("the identifier repeats inside the value: %#v", bitcoin)
	}

	// And the timestamp sits beside the currencies, at the same depth and
	// under the same kind of key, which is the trap the Recipe is about.
	if _, found := bitcoin["last_updated_at"]; !found {
		t.Errorf("last_updated_at is not beside the currencies: %#v", bitcoin)
	}
}
