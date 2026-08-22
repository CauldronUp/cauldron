package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A create used to store the decoded request body as the record, so anything
// sent came back. Posting a made-up field to Adyen's refund answered with the
// made-up field in it, and no provider does that.
//
// The stray key is not the cost. The cost is that a conformance case
// asserting a value it sent on a create cannot fail: the value returns
// whether or not the Recipe declares the field. Adyen's refund carried the
// payment's name for its reference -- merchantReference, where Adyen sends
// reference -- and a case asserting the name would have passed either way.
func TestACreateDoesNotEchoAnUndeclaredField(t *testing.T) {
	body := createRefund(t, `{"merchantAccount":"CauldronECOM","totallyMadeUpField":"xyzzy","amount":{"value":1000,"currency":"EUR"}}`)

	if _, ok := body["totallyMadeUpField"]; ok {
		t.Errorf("a field the resource does not declare came back: %v", body)
	}

	// And the declared ones still do, or this would pass by answering nothing.
	if body["status"] != "received" {
		t.Errorf("status = %v, want received", body["status"])
	}
}

// merchantAccount is sent by every Adyen request and is not a field on the
// refund. It used to be echoed for the same reason, which made the response
// look more like Adyen's than it was.
func TestAFieldTheProviderDoesNotAnswerWithIsNotEchoed(t *testing.T) {
	body := createRefund(t, `{"merchantAccount":"CauldronECOM","amount":{"value":1000,"currency":"EUR"}}`)

	if _, ok := body["merchantAccount"]; ok {
		t.Errorf("merchantAccount is not on Adyen's refund response and came back: %v", body)
	}
}

// Some providers really do accept arbitrary keys, and Stripe's metadata is
// the reason type: map exists. A Recipe says so rather than the runtime
// assuming it of every field.
func TestAMapFieldKeepsWhateverItWasSent(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/customers",
		strings.NewReader("email=grace@example.com&metadata[order_id]=42&metadata[tenant]=acme"))
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v\n%s", err, rec.Body)
	}

	metadata, ok := body["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %v, want the keys it was sent", body["metadata"])
	}

	if metadata["order_id"] != float64(42) || metadata["tenant"] != "acme" {
		t.Errorf("metadata = %v, want both keys through unchanged", metadata)
	}
}

func createRefund(t *testing.T, body string) map[string]any {
	t.Helper()

	r, err := recipe.Open("adyen")
	if err != nil {
		t.Fatalf("open adyen: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-merchant"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v71/payments/8836183744AAAAAA/refunds", strings.NewReader(body))
	req.Header.Set("X-API-Key", "AQEyhmfxK4cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v\n%s", err, rec.Body)
	}

	return decoded
}
