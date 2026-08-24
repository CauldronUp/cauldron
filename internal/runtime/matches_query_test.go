package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some APIs answer one path with two shapes, and which one you get depends on
// whether the request asked for more. Clover is the clearest case: an order
// carries no line items unless the query says ?expand=lineItems, and a request
// that forgets is not refused -- it gets an order with the right total and
// nothing in it, which is a perfectly ordinary thing for a shop to have.
//
// So the failure is a report saying a busy Saturday sold nothing in
// particular. The format could not express it: selects reads the request body
// and matches_header reads a header, and neither is where expand lives.
// Asana's Recipe says the same thing about opt_fields in as many words.
//
// matches_query is the third spelling of the one idea, and it matches on
// membership rather than equality because that is what these providers send:
// ?expand=lineItems,payments asks for two things in one parameter.
func TestARouteMayBePickedByAQueryParameter(t *testing.T) {
	r, err := recipe.Open("clover")
	if err != nil {
		t.Fatalf("open clover: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-merchant"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	const order = "/v3/merchants/CAULDRONMERCH1/orders/a1b2c3d4e5f60"

	get := func(path string) map[string]any {
		t.Helper()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer cauldron.clover.merchant.token")
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", path, rec.Code)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: decode: %v", path, err)
		}

		return body
	}

	// Forgetting is answered, not refused.
	compact := get(order)
	if compact["total"] == nil {
		t.Fatal("the compact order lost its total")
	}

	if _, present := compact["lineItems"]; present {
		t.Error("an order that did not ask to expand carried line items")
	}

	// Asking brings them.
	expanded := get(order + "?expand=lineItems")
	if _, present := expanded["lineItems"]; !present {
		t.Error("an order that asked to expand lineItems did not carry them")
	}

	// Asking for one thing brings one thing. This is the half that makes the
	// trap survive a partial fix: a client that expanded what it needed last
	// month gets the same silence for what it needs now.
	payments := get(order + "?expand=payments")
	if _, present := payments["payments"]; !present {
		t.Error("an order that asked to expand payments did not carry them")
	}

	if _, present := payments["lineItems"]; present {
		t.Error("expanding payments also expanded lineItems")
	}

	// Membership rather than equality: one parameter, two things asked for,
	// and a route selecting either member answers.
	both := get(order + "?expand=lineItems,payments")
	if _, present := both["lineItems"]; !present {
		t.Error("a comma list did not match a route selecting one of its members")
	}

	// And one request reaches one route, so only one of the two comes back.
	// Clover would send both. That is a stated gap in the Recipe rather than
	// a surprise here: matching composes, answering does not.
	if _, present := both["payments"]; present {
		t.Error("one request answered from two routes at once")
	}

	// And a parameter the Recipe does not model falls through to the compact
	// route rather than 404ing, because not being expanded is not an error.
	unknown := get(order + "?expand=customer")
	if _, present := unknown["lineItems"]; present {
		t.Error("an unmodelled expand matched a route that declares a different one")
	}

	if unknown["total"] == nil {
		t.Error("an unmodelled expand should fall through to the compact route")
	}
}
