package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mounted(t *testing.T) *Server {
	t.Helper()

	s := New()

	if err := s.Mount("stripe", 1, ""); err != nil {
		t.Fatalf("Mount: %v", err)
	}

	return s
}

func do(t *testing.T, s *Server, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, target, reader)
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func body(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, rec.Body.String())
	}

	return out
}

func TestPathPrefixAddressing(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost:4600/stripe/v1/customers", `{"email":"ada@example.com"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	if id, _ := body(t, rec)["id"].(string); !strings.HasPrefix(id, "cus_") {
		t.Errorf("id = %q", id)
	}
}

func TestHostPrefixAddressing(t *testing.T) {
	s := mounted(t)

	req := httptest.NewRequest(http.MethodPost, "http://stripe.cauldron.test/v1/customers", strings.NewReader(`{"email":"ada@example.com"}`))
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}
}

func TestUnmountedRecipeExplainsWhatIsRunning(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodGet, "http://localhost/shopify/admin/api/orders.json", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "stripe") {
		t.Errorf("the error should say what is mounted\n%s", rec.Body)
	}
}

func TestStatus(t *testing.T) {
	s := mounted(t)

	do(t, s, http.MethodPost, "http://localhost/stripe/v1/customers", `{"email":"ada@example.com"}`)

	rec := do(t, s, http.MethodGet, "http://localhost/_cauldron/status", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	out := body(t, rec)

	recipes, _ := out["recipes"].([]any)
	if len(recipes) != 1 {
		t.Fatalf("got %d recipes, want 1", len(recipes))
	}

	first := recipes[0].(map[string]any)

	if first["recipe"] != "stripe" {
		t.Errorf("recipe = %v", first["recipe"])
	}

	if first["requests"] != float64(1) {
		t.Errorf("requests = %v, want 1", first["requests"])
	}
}

func TestSeedAndReset(t *testing.T) {
	s := mounted(t)

	if rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/seed?fixture=small-shop", ""); rec.Code != http.StatusOK {
		t.Fatalf("seed status = %d\n%s", rec.Code, rec.Body)
	}

	list := body(t, do(t, s, http.MethodGet, "http://localhost/stripe/v1/customers", ""))
	if data, _ := list["data"].([]any); len(data) != 2 {
		t.Fatalf("got %d seeded customers, want 2", len(data))
	}

	if rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/reset", ""); rec.Code != http.StatusOK {
		t.Fatalf("reset status = %d", rec.Code)
	}

	after := body(t, do(t, s, http.MethodGet, "http://localhost/stripe/v1/customers", ""))
	if data, _ := after["data"].([]any); len(data) != 2 {
		t.Errorf("reset should restore the fixture, got %d records", len(data))
	}
}

func TestSeedRejectsAnUnknownFixture(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/seed?fixture=nope", "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "small-shop") {
		t.Errorf("should list available fixtures\n%s", rec.Body)
	}
}

func TestFaultViaControlAPI(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/fault", `{"error":"rate_limit","count":1}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("arm status = %d\n%s", rec.Code, rec.Body)
	}

	if got := do(t, s, http.MethodGet, "http://localhost/stripe/v1/customers", ""); got.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429", got.Code)
	}

	if got := do(t, s, http.MethodGet, "http://localhost/stripe/v1/customers", ""); got.Code != http.StatusOK {
		t.Errorf("the fault had count 1 and should be spent; got %d", got.Code)
	}
}

func TestFaultRejectsAnUnknownError(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/fault", `{"error":"meteor"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "rate_limit") {
		t.Errorf("should list available errors\n%s", rec.Body)
	}
}

func TestClockAdvanceMovesEverySandbox(t *testing.T) {
	s := mounted(t)

	before := s.Clock().Unix()

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/clock/advance", `{"duration":"30d"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	after := s.Clock().Unix()

	if after-before != 30*24*60*60 {
		t.Errorf("advanced by %d seconds, want 30 days", after-before)
	}

	created := body(t, do(t, s, http.MethodPost, "http://localhost/stripe/v1/customers", `{"email":"ada@example.com"}`))

	if created["created"] != float64(after) {
		t.Errorf("new records should carry the advanced time; got %v want %v", created["created"], after)
	}
}

func TestClockAdvanceRejectsNonsense(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/clock/advance", `{"duration":"next tuesday"}`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestEmitViaControlAPI(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/emit",
		`{"event":"payment_intent.payment_failed","data":{"id":"pi_test","amount":2500}}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	sandbox, _ := s.Sandbox("stripe")

	deliveries := sandbox.Webhooks().Deliveries()
	if len(deliveries) != 1 {
		t.Fatalf("got %d deliveries, want 1", len(deliveries))
	}

	if deliveries[0].Event != "payment_intent.payment_failed" {
		t.Errorf("event = %q", deliveries[0].Event)
	}
}

func TestSubscribeThenReceive(t *testing.T) {
	s := mounted(t)

	received := make(chan string, 1)

	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, _ := io.ReadAll(r.Body)
		received <- string(payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer endpoint.Close()

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/subscribe", `{"url":"`+endpoint.URL+`"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("subscribe status = %d\n%s", rec.Code, rec.Body)
	}

	do(t, s, http.MethodPost, "http://localhost/stripe/v1/customers", `{"email":"ada@example.com"}`)

	select {
	case payload := <-received:
		if !strings.Contains(payload, "customer.created") {
			t.Errorf("payload = %s", payload)
		}
	default:
		t.Fatal("the subscriber received nothing")
	}
}

func TestRequestsEndpoint(t *testing.T) {
	s := mounted(t)

	do(t, s, http.MethodPost, "http://localhost/stripe/v1/customers", `{"email":"ada@example.com"}`)

	out := body(t, do(t, s, http.MethodGet, "http://localhost/_cauldron/stripe/requests", ""))

	requests, _ := out["requests"].([]any)
	if len(requests) != 1 {
		t.Fatalf("got %d requests, want 1", len(requests))
	}

	first := requests[0].(map[string]any)
	if first["Op"] != "create" {
		t.Errorf("Op = %v", first["Op"])
	}
}

func TestUnknownControlEndpoint(t *testing.T) {
	s := mounted(t)

	if rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/summon", ""); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestControlPathIsNotConfusedForARecipe(t *testing.T) {
	s := mounted(t)

	// A recipe named "_cauldron" cannot shadow the control API.
	rec := do(t, s, http.MethodGet, "http://localhost/_cauldron/status", "")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d", rec.Code)
	}
}
