package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Tradier's own words: if you have a single order, it will be returned as a
// JSON obj/dict whereas multiple orders will be returned as an array/list of
// JSON objects.
//
// This is what JSON generated from XML does, because a single child element
// and a repeated one are the same thing there. The format already described
// the safe half of this axis -- Xero sends one invoice as a list of one, which
// resource.array covers, and a client written for the list keeps working. This
// is the other half, where a client written against a fixture holding two
// records crashes the first time production holds one.

func tradierOrders(t *testing.T, fixture string) map[string]any {
	t.Helper()

	r, err := recipe.Open("tradier")
	if err != nil {
		t.Fatalf("open tradier: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed(fixture); err != nil {
		t.Fatalf("seed %s: %v", fixture, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/accounts/VA000001/orders", nil)
	req.Header.Set("Authorization", "Bearer cauldron000000000000000tradier")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("not JSON: %v %s", err, rec.Body.String())
	}

	return body
}

func orderNode(t *testing.T, body map[string]any) any {
	t.Helper()

	orders, ok := body["orders"].(map[string]any)
	if !ok {
		t.Fatalf("no orders envelope: %v", body)
	}

	return orders["order"]
}

func TestSeveralRecordsStayAList(t *testing.T) {
	if _, ok := orderNode(t, tradierOrders(t, "busy-account")).([]any); !ok {
		t.Errorf("two orders should be an array, got %T", orderNode(t, tradierOrders(t, "busy-account")))
	}
}

func TestOneRecordCollapsesToAnObject(t *testing.T) {
	node := orderNode(t, tradierOrders(t, "quiet-account"))

	object, ok := node.(map[string]any)
	if !ok {
		t.Fatalf("one order should be an object, got %T", node)
	}

	if object["symbol"] != "AAPL" {
		t.Errorf("the collapsed object is not the record: %v", object)
	}
}

func TestCollapsingIsOptIn(t *testing.T) {
	// The overwhelming majority of providers do not do this, and a listing
	// that quietly changed shape at one record would be a far worse bug than
	// the one it describes.
	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	if r.Responses.List.CollapseSingle {
		t.Error("github should not collapse")
	}
}
