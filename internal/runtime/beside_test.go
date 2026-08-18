package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// One endpoint, several collections, and the fact that it is one endpoint is
// the point. GoCardless answers a request for transactions with a booked array
// and a pending array in one body: the same purchase appears in pending first
// and booked later, with a different identifier, so code that merges the two
// counts it twice and code that matches on the id across them loses it.
//
// Describing that as two endpoints would lose the thing worth describing, and
// describing only one of the arrays would answer with a shape no bank sends.

func gocardlessBody(t *testing.T, target string) map[string]any {
	t.Helper()

	r, err := recipe.Open("gocardlessbank")
	if err != nil {
		t.Fatalf("open gocardlessbank: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-customer"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer cauldron-gocardless-bank-access-token")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	return decode(t, rec)
}

func arrayAt(t *testing.T, body map[string]any, outer, inner string) []any {
	t.Helper()

	object, _ := body[outer].(map[string]any)
	if object == nil {
		t.Fatalf("%s is not an object: %v", outer, body[outer])
	}

	list, _ := object[inner].([]any)
	if list == nil {
		t.Fatalf("%s.%s is not a list: %v", outer, inner, object[inner])
	}

	return list
}

func TestBesideCollectionsShareOneResponse(t *testing.T) {
	body := gocardlessBody(t,
		"/api/v2/accounts/a7c31f95-2e84-4d60-9b17-5c02e8a41d73/transactions/")

	if got := len(arrayAt(t, body, "transactions", "booked")); got != 2 {
		t.Errorf("booked has %d records, want 2", got)
	}

	if got := len(arrayAt(t, body, "transactions", "pending")); got != 1 {
		t.Errorf("pending has %d records, want 1", got)
	}
}

func TestABesideCollectionIsScopedLikeTheRouteItRidesOn(t *testing.T) {
	// The scope is applied to every collection in the body, because they are
	// one request. A beside collection returned whole would hand one
	// customer's pending transactions to another.
	body := gocardlessBody(t,
		"/api/v2/accounts/00000000-0000-4000-8000-000000000000/transactions/")

	if got := len(arrayAt(t, body, "transactions", "pending")); got != 0 {
		t.Errorf("another account's pending transactions leaked: %d", got)
	}
}

func TestABesideCollectionDropsTheScopeFromItsRecords(t *testing.T) {
	// The partition is in the path and a provider that puts it there does not
	// repeat it in the body. The route's own returns names the route's own
	// resource's fields and means nothing for a different one.
	body := gocardlessBody(t,
		"/api/v2/accounts/a7c31f95-2e84-4d60-9b17-5c02e8a41d73/transactions/")

	for _, name := range []string{"booked", "pending"} {
		for i, record := range arrayAt(t, body, "transactions", name) {
			object, _ := record.(map[string]any)
			if object == nil {
				t.Fatalf("%s[%d] is not an object", name, i)
			}

			if _, present := object["account_id"]; present {
				t.Errorf("%s[%d] carries account_id, which is in the path: %v", name, i, object)
			}
		}
	}
}

func TestABesideCollectionIsPresentedLikeAnyOther(t *testing.T) {
	// Renaming and nesting run for it too, so a pending transaction gets its
	// transactionId and its transactionAmount rather than the store's names.
	pending := arrayAt(t, gocardlessBody(t,
		"/api/v2/accounts/a7c31f95-2e84-4d60-9b17-5c02e8a41d73/transactions/"),
		"transactions", "pending")

	first, _ := pending[0].(map[string]any)
	if first == nil {
		t.Fatalf("not an object: %v", pending[0])
	}

	if first["transactionId"] != "2026081501004466099000" {
		t.Errorf("transactionId = %v", first["transactionId"])
	}

	amount, _ := first["transactionAmount"].(map[string]any)
	if amount == nil {
		t.Fatalf("transactionAmount is not an object: %v", first)
	}

	if amount["amount"] != "-100.00" {
		t.Errorf("transactionAmount.amount = %v", amount["amount"])
	}
}
