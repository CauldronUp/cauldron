package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A webhook payload used to carry the store's internal field names, because
// the record went into it raw and never through the shaping every HTTP
// response goes through.
//
// So the same record, from the same sandbox, at the same instant, had two
// different shapes depending on how you looked at it: the response nested
// amount under amount_money as the Recipe declares, and the webhook did not.
// An application's handler written against the emulator reads
// event.data.object.amount and is a hundred per cent locally green and nought
// per cent correct, because real Square sends the nested one.
//
// It was live on 81 of the 85 webhook-emitting Recipes, and no conformance
// case could see it: the format had no way to assert a webhook at all.
func TestAWebhookCarriesTheSameShapeAsTheResponse(t *testing.T) {
	r, err := recipe.Open("square")
	if err != nil {
		t.Fatalf("open square: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v2/payments",
		strings.NewReader(`{"amount":2500,"currency":"USD","source_id":"cnon:card-nonce-ok"}`))
	req.Header.Set("Authorization", "Bearer EAAAcauldron")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Square-Version", "2026-01-22")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("create = %d\n%s", rec.Code, rec.Body)
	}

	deliveries := s.Webhooks().Deliveries()
	if len(deliveries) == 0 {
		t.Fatal("creating a payment emitted no webhook, so there is nothing to compare")
	}

	object := deliveryObject(t, deliveries[len(deliveries)-1])

	// The declared nesting, which the HTTP response has had all along.
	money, nested := object["amount_money"].(map[string]any)
	if !nested {
		t.Fatalf("the webhook object has no amount_money, so it is carrying the store's own field names: %v", keysOf(object))
	}

	if money["amount"] != float64(2500) {
		t.Errorf("amount_money.amount = %v, want 2500", money["amount"])
	}

	if money["currency"] != "USD" {
		t.Errorf("amount_money.currency = %v, want USD", money["currency"])
	}

	// And the flat names must not be there beside them, or a handler finds
	// both shapes and picks whichever it was written for.
	for _, internal := range []string{"amount", "currency"} {
		if _, present := object[internal]; present {
			t.Errorf("the webhook object still carries the internal name %q beside the declared one", internal)
		}
	}
}

// deliveryObject digs the changed record out of whatever envelope the Recipe
// declares, so the test is about the record's shape rather than the envelope's.
func deliveryObject(t *testing.T, d Delivery) map[string]any {
	t.Helper()

	data, ok := d.Payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("the payload has no data: %v", d.Payload)
	}

	object, ok := data["object"].(map[string]any)
	if !ok {
		t.Fatalf("the payload has no data.object: %v", data)
	}

	return object
}
