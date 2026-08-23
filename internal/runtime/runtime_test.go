package runtime

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

func stripe(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe recipe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	return s
}

// call makes an authenticated request against the sandbox.
func call(t *testing.T, s *Sandbox, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v\n%s", err, rec.Body.String())
	}

	return out
}

func TestCreateMintsAnIdentifierAndAppliesDefaults(t *testing.T) {
	s := stripe(t)

	rec := call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	body := decode(t, rec)

	id, _ := body["id"].(string)
	if !strings.HasPrefix(id, "cus_") {
		t.Errorf("id = %q, want a cus_ prefix", id)
	}

	if body["currency"] != "usd" {
		t.Errorf("currency = %v, want the declared default usd", body["currency"])
	}

	// created must come from the sandbox clock, not the wall clock.
	if got, want := body["created"], float64(s.Clock().Unix()); got != want {
		t.Errorf("created = %v, want %v", got, want)
	}
}

func TestCreateRejectsAMissingRequiredField(t *testing.T) {
	s := stripe(t)

	rec := call(t, s, http.MethodPost, "/v1/customers", `{}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400\n%s", rec.Code, rec.Body)
	}

	if !strings.Contains(rec.Body.String(), "email") {
		t.Errorf("the error should name the missing field\n%s", rec.Body)
	}
}

// Stripe's official SDKs post form encoding, so a fake that only spoke JSON
// would reject the clients it exists to serve.
func TestCreateAcceptsFormEncoding(t *testing.T) {
	s := stripe(t)

	req := httptest.NewRequest(http.MethodPost, "/v1/customers", strings.NewReader("email=grace@example.com&name=Grace+Hopper&metadata[order_id]=42"))
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	body := decode(t, rec)

	if body["email"] != "grace@example.com" {
		t.Errorf("email = %v", body["email"])
	}

	metadata, ok := body["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %v, want a nested object", body["metadata"])
	}

	if metadata["order_id"] != float64(42) {
		t.Errorf("metadata.order_id = %v, want 42 as a number", metadata["order_id"])
	}
}

func TestGetUpdateDelete(t *testing.T) {
	s := stripe(t)

	created := decode(t, call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com","name":"Ada"}`))
	id := created["id"].(string)

	fetched := decode(t, call(t, s, http.MethodGet, "/v1/customers/"+id, ""))
	if fetched["email"] != "ada@example.com" {
		t.Errorf("get returned %v", fetched)
	}

	updated := decode(t, call(t, s, http.MethodPost, "/v1/customers/"+id, `{"name":"Ada Lovelace"}`))
	if updated["name"] != "Ada Lovelace" {
		t.Errorf("update returned %v", updated["name"])
	}

	if updated["email"] != "ada@example.com" {
		t.Errorf("update dropped an unrelated field: %v", updated)
	}

	deleted := decode(t, call(t, s, http.MethodDelete, "/v1/customers/"+id, ""))
	if deleted["deleted"] != true {
		t.Errorf("delete returned %v", deleted)
	}

	if rec := call(t, s, http.MethodGet, "/v1/customers/"+id, ""); rec.Code != http.StatusNotFound {
		t.Errorf("get after delete = %d, want 404", rec.Code)
	}
}

func TestMissingResourceReturns404WithAUsefulMessage(t *testing.T) {
	s := stripe(t)

	rec := call(t, s, http.MethodGet, "/v1/customers/cus_nope", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "cus_nope") {
		t.Errorf("the error should name the identifier\n%s", rec.Body)
	}
}

func TestListPaginates(t *testing.T) {
	s := stripe(t)

	for i := 0; i < 5; i++ {
		call(t, s, http.MethodPost, "/v1/customers", `{"email":"a@example.com"}`)
	}

	first := decode(t, call(t, s, http.MethodGet, "/v1/customers?limit=2", ""))

	if first["object"] != "list" {
		t.Errorf("object = %v, want list", first["object"])
	}

	data, _ := first["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("got %d records, want 2", len(data))
	}

	if first["has_more"] != true {
		t.Errorf("has_more = %v, want true", first["has_more"])
	}

	if _, present := first["next_cursor"]; present {
		t.Error("Stripe does not send a next_cursor, and neither should the emulator")
	}

	if first["url"] != "/v1/customers" {
		t.Errorf("url = %v, want the request path", first["url"])
	}

	// Page the way a Stripe client does: hand the last id back as starting_after.
	cursor, _ := data[len(data)-1].(map[string]any)["id"].(string)
	if cursor == "" {
		t.Fatal("records should carry an id to page from")
	}

	second := decode(t, call(t, s, http.MethodGet, "/v1/customers?limit=2&starting_after="+cursor, ""))
	secondData, _ := second["data"].([]any)

	if len(secondData) != 2 {
		t.Fatalf("second page has %d records, want 2", len(secondData))
	}

	firstID := data[0].(map[string]any)["id"]
	for _, record := range secondData {
		if record.(map[string]any)["id"] == firstID {
			t.Error("a record appeared on two pages")
		}
	}
}

func TestAuthIsEnforcedUsingTheRecipesOwnError(t *testing.T) {
	s := stripe(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	// The Recipe declares authentication_error as 401.
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401\n%s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	req.Header.Set("Authorization", "Bearer sk_test_wrong")
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("a wrong key returned %d, want 401", rec.Code)
	}
}

func TestUnknownRouteAndMethod(t *testing.T) {
	s := stripe(t)

	if rec := call(t, s, http.MethodGet, "/v1/nonsense", ""); rec.Code != http.StatusNotFound {
		t.Errorf("unknown path = %d, want 404", rec.Code)
	}

	rec := call(t, s, http.MethodPatch, "/v1/customers", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("wrong method = %d, want 405", rec.Code)
	}

	if allow := rec.Header().Get("Allow"); !strings.Contains(allow, "POST") {
		t.Errorf("Allow = %q, want it to include POST", allow)
	}
}

func TestFaultInjectionUsesTheRecipesStatusAndHeaders(t *testing.T) {
	s := stripe(t)

	if err := s.Arm(Fault{Error: "rate_limit"}); err != nil {
		t.Fatalf("Arm: %v", err)
	}

	rec := call(t, s, http.MethodGet, "/v1/customers", "")

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}

	if got := rec.Header().Get("Retry-After"); got != "1" {
		t.Errorf("Retry-After = %q, want 1 — headers come from the Recipe", got)
	}
}

func TestFaultsApplyBeforeAuth(t *testing.T) {
	s := stripe(t)
	_ = s.Arm(Fault{Error: "rate_limit"})

	// No credential at all: a real rate limiter does not check your key first.
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429", rec.Code)
	}
}

func TestFaultCanBeLimitedToACount(t *testing.T) {
	s := stripe(t)
	_ = s.Arm(Fault{Error: "rate_limit", Count: 2})

	for i := 1; i <= 2; i++ {
		if rec := call(t, s, http.MethodGet, "/v1/customers", ""); rec.Code != http.StatusTooManyRequests {
			t.Fatalf("request %d = %d, want 429", i, rec.Code)
		}
	}

	if rec := call(t, s, http.MethodGet, "/v1/customers", ""); rec.Code != http.StatusOK {
		t.Errorf("third request = %d, want 200 — the fault should be spent", rec.Code)
	}
}

func TestFaultExpiresOnTheSandboxClock(t *testing.T) {
	s := stripe(t)

	_ = s.Arm(Fault{Error: "rate_limit", Until: s.Clock().Now().Add(30 * time.Second)})

	if rec := call(t, s, http.MethodGet, "/v1/customers", ""); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("before expiry = %d, want 429", rec.Code)
	}

	s.Clock().Advance(31 * time.Second)

	if rec := call(t, s, http.MethodGet, "/v1/customers", ""); rec.Code != http.StatusOK {
		t.Errorf("after expiry = %d, want 200", rec.Code)
	}
}

func TestArmRejectsAnUndeclaredError(t *testing.T) {
	s := stripe(t)

	err := s.Arm(Fault{Error: "meteor_strike"})
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "rate_limit") {
		t.Errorf("the error should list what is available, got %q", err)
	}
}

func TestCreateEmitsTheLifecycleWebhook(t *testing.T) {
	s := stripe(t)

	call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`)

	deliveries := s.Webhooks().Deliveries()

	if len(deliveries) != 1 {
		t.Fatalf("got %d deliveries, want 1", len(deliveries))
	}

	if deliveries[0].Event != "customer.created" {
		t.Errorf("event = %q", deliveries[0].Event)
	}

	if deliveries[0].Signature == "" {
		t.Error("expected a signature — applications verify with the provider SDK")
	}

	if !strings.HasPrefix(deliveries[0].Signature, "t=") {
		t.Errorf("signature = %q, want a Stripe-shaped t=…,v1=… value", deliveries[0].Signature)
	}
}

func TestWebhooksAreDeliveredToSubscribers(t *testing.T) {
	s := stripe(t)

	var (
		mu       sync.Mutex
		received []string
		gotSig   string
	)

	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		received = append(received, string(body))
		gotSig = r.Header.Get("Stripe-Signature")
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer endpoint.Close()

	s.Webhooks().Subscribe(endpoint.URL)

	call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`)

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 1 {
		t.Fatalf("endpoint received %d webhooks, want 1", len(received))
	}

	if !strings.Contains(received[0], "customer.created") {
		t.Errorf("payload = %s", received[0])
	}

	if gotSig == "" {
		t.Error("the signing header was not sent")
	}

	deliveries := s.Webhooks().Deliveries()
	if deliveries[0].Status != http.StatusOK {
		t.Errorf("recorded status = %d, want 200", deliveries[0].Status)
	}
}

func TestEmitRejectsAnUndeclaredEvent(t *testing.T) {
	s := stripe(t)

	_, _, err := s.Webhooks().Emit("customer.exploded", store.Record{})
	if err == nil {
		t.Fatal("expected an error for an event the provider never sends")
	}
}

// Every webhook used to leave in Stripe's envelope, whatever the provider. That
// was true while Stripe was the only Recipe and became a quiet lie as the
// collection grew: nobody else sends {id, type, created, data.object}.
//
// Adyen is the case that makes it matter. Its notification is an array of
// NotificationRequestItem, the payment's fields sit beside eventCode rather
// than under it, and success is the string "true" — so the obvious truth test
// on it is also true for "false", and every failed payment reads as
// successful. Emitting Stripe's shape here would hide the single most
// expensive thing about integrating with Adyen.
func TestADeclaredWebhookEnvelopeIsUsed(t *testing.T) {
	r, err := recipe.Open("adyen")
	if err != nil {
		t.Fatalf("open adyen: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	delivery, _, err := s.Webhooks().Emit("AUTHORISATION", store.Record{
		"pspReference":      "8836183744AAAAAA",
		"merchantReference": "order-1001",
	})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	// Stripe's envelope must not be there at all, or a client written against
	// the real Adyen finds two shapes and picks the wrong one.
	for _, absent := range []string{"id", "type", "created", "data"} {
		if _, present := delivery.Fields()[absent]; present {
			t.Errorf("Adyen sends no top-level %q", absent)
		}
	}

	items, _ := delivery.Fields()["notificationItems"].([]any)
	if len(items) != 1 {
		t.Fatalf("notificationItems = %v, want one entry", delivery.Fields()["notificationItems"])
	}

	wrapper, _ := items[0].(map[string]any)
	item, _ := wrapper["NotificationRequestItem"].(map[string]any)

	if item["eventCode"] != "AUTHORISATION" {
		t.Errorf("eventCode = %v", item["eventCode"])
	}

	// The string, not the boolean. A client writing if (item.success) is also
	// inside the branch for "false".
	if success, isString := item["success"].(string); !isString || success != "true" {
		t.Errorf("success = %#v, want the string \"true\"", item["success"])
	}

	// The record is spliced in beside eventCode, not nested under it.
	if item["merchantReference"] != "order-1001" {
		t.Errorf("merchantReference = %v, want the record merged in", item["merchantReference"])
	}

	if live, isString := delivery.Fields()["live"].(string); !isString || live != "false" {
		t.Errorf("live = %#v, want the string \"false\"", delivery.Fields()["live"])
	}
}

// Recipes written before envelopes were declarable keep the shape they had.
func TestAnUndeclaredWebhookEnvelopeStaysTheDefault(t *testing.T) {
	s := stripe(t)

	delivery, _, err := s.Webhooks().Emit("customer.created", store.Record{"id": "cus_cauldron"})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	if delivery.Fields()["type"] != "customer.created" {
		t.Errorf("type = %v", delivery.Fields()["type"])
	}

	data, _ := delivery.Fields()["data"].(map[string]any)
	object, _ := data["object"].(map[string]any)

	if object["id"] != "cus_cauldron" {
		t.Errorf("data.object = %v, want the record", data["object"])
	}
}

// WooCommerce posts the resource itself, flat, with no envelope of any kind.
// Which event it was arrives in headers rather than the body, so a handler
// reading body.type or body.event finds nothing. That is the opposite extreme
// from Adyen's nested array, and both have to be expressible or the default
// envelope is a claim rather than a fallback.
func TestAWebhookEnvelopeCanBeNoEnvelopeAtAll(t *testing.T) {
	r, err := recipe.Open("woocommerce")
	if err != nil {
		t.Fatalf("open woocommerce: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	delivery, _, err := s.Webhooks().Emit("order.updated", store.Record{
		"id":     "1001",
		"status": "processing",
		"total":  "58.32",
	})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	if delivery.Fields()["status"] != "processing" {
		t.Errorf("status = %v, want the record at the top level", delivery.Fields()["status"])
	}

	// Nothing wrapping it, and nothing naming the event.
	for _, absent := range []string{"type", "event", "created", "data"} {
		if _, present := delivery.Fields()[absent]; present {
			t.Errorf("WooCommerce sends no top-level %q", absent)
		}
	}
}

func TestSeedLoadsAFixture(t *testing.T) {
	s := stripe(t)

	if err := s.Seed("small-shop"); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	body := decode(t, call(t, s, http.MethodGet, "/v1/customers", ""))

	data, _ := body["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("got %d seeded customers, want 2", len(data))
	}

	first := data[0].(map[string]any)
	if first["id"] != "cus_00000000000001" {
		t.Errorf("fixtures must keep their pinned identifiers; got %v", first["id"])
	}
}

func TestSeedRejectsAnUnknownFixture(t *testing.T) {
	s := stripe(t)

	err := s.Seed("enormous-shop")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "small-shop") {
		t.Errorf("the error should list available fixtures, got %q", err)
	}
}

func TestResetMakesASandboxIndistinguishableFromAFreshOne(t *testing.T) {
	s := stripe(t)
	_ = s.Seed("small-shop")

	before := decode(t, call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`))

	_ = s.Arm(Fault{Error: "rate_limit"})
	s.Clock().Advance(72 * time.Hour)

	if err := s.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}

	if len(s.ArmedFaults()) != 0 {
		t.Error("reset must disarm faults")
	}

	if len(s.Webhooks().Deliveries()) != 0 {
		t.Error("reset must clear webhook history")
	}

	after := decode(t, call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`))

	if after["id"] != before["id"] {
		t.Errorf("after reset id = %v, want %v — a reset sandbox must mint the same identifiers", after["id"], before["id"])
	}
}

func TestTwoSandboxesWithTheSameSeedAgree(t *testing.T) {
	var ids []string

	for run := 0; run < 3; run++ {
		s := stripe(t)
		_ = s.Seed("small-shop")

		body := decode(t, call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`))
		id := body["id"].(string)

		if run == 0 {
			ids = append(ids, id)
			continue
		}

		if id != ids[0] {
			t.Fatalf("run %d minted %s, run 0 minted %s", run, id, ids[0])
		}
	}
}

func TestRequestLogRecordsExchanges(t *testing.T) {
	s := stripe(t)

	call(t, s, http.MethodPost, "/v1/customers", `{"email":"ada@example.com"}`)
	call(t, s, http.MethodGet, "/v1/customers/cus_nope", "")

	entries := s.Exchanges(0)

	if len(entries) != 2 {
		t.Fatalf("got %d exchanges, want 2", len(entries))
	}

	if entries[0].Op != "create" || entries[0].Status != http.StatusOK {
		t.Errorf("first exchange = %+v", entries[0])
	}

	if entries[1].Status != http.StatusNotFound {
		t.Errorf("second exchange = %+v", entries[1])
	}

	if entries[0].Seq >= entries[1].Seq {
		t.Error("exchanges should be sequenced in order")
	}
}

func TestRequestLogRecordsTheFaultThatFired(t *testing.T) {
	s := stripe(t)
	_ = s.Arm(Fault{Error: "rate_limit"})

	call(t, s, http.MethodGet, "/v1/customers", "")

	entries := s.Exchanges(0)

	if entries[0].Fault != "rate_limit" {
		t.Errorf("Fault = %q, want rate_limit", entries[0].Fault)
	}
}

// The switch over auth schemes and the list of schemes the validator accepts
// are two separate pieces of code that have to agree, and nothing makes them.
//
// This test is the coupling. Adding a scheme to the valid list without adding
// a case to the switch is a one-line change that would silently authorise
// every request against every Recipe using it, and until now nothing would
// have failed.
func TestEveryValidAuthSchemeActuallyChecksSomething(t *testing.T) {
	schemes := map[string]recipe.Auth{
		"bearer": {Scheme: "bearer", Header: "Authorization", Prefix: "Bearer ", Keys: []string{"right"}},
		"basic":  {Scheme: "basic", Credential: "password", Keys: []string{"right"}},
		"header": {Scheme: "header", Header: "X-Api-Key", Keys: []string{"right"}},
		"query":  {Scheme: "query", Param: "api_key", Keys: []string{"right"}},
	}

	// none is the one scheme that is meant to accept anything, and it is
	// handled before the switch, so it is not in the table above. Anything
	// else appearing in the valid list without a case here is the bug.
	for _, scheme := range recipe.ValidAuthSchemes() {
		if scheme == "" || scheme == "none" {
			continue
		}

		if _, covered := schemes[scheme]; !covered {
			t.Errorf("auth scheme %q is accepted by the validator and this test does not exercise it, so nothing checks that the handler does either", scheme)
		}
	}

	for name, auth := range schemes {
		s := sandboxWithAuth(t, auth)

		if authorisedWith(s, auth, "right") != true {
			t.Errorf("%s: the correct credential was refused", name)
		}

		for _, wrong := range []string{"", "wrong", "righter", "righ"} {
			if authorisedWith(s, auth, wrong) {
				t.Errorf("%s: %q was accepted", name, wrong)
			}
		}
	}
}

// An unknown scheme is refused twice: once by the validator, which is the
// defence that works today, and once by the handler, which is the one that
// still works if somebody adds a scheme to the valid list and forgets the
// switch.
func TestAnUnknownAuthSchemeIsRefusedTwice(t *testing.T) {
	auth := recipe.Auth{Scheme: "some-scheme-nobody-implemented", Keys: []string{"right"}}

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	clone := *r
	clone.Auth = auth

	// The first defence. A Recipe naming a scheme nothing implements does not
	// load at all.
	if _, err := New(&clone, Options{Seed: 1}); err == nil {
		t.Fatal("a Recipe with an unknown auth scheme was accepted by the validator")
	}

	// The second. Reaching past the validator, as a future edit to the valid
	// list would, the handler must still refuse rather than wave everything
	// through, which is what it used to do.
	s := sandboxWithAuth(t, recipe.Auth{Scheme: "bearer", Prefix: "Bearer ", Keys: []string{"right"}})
	s.recipe.Auth = auth

	if authorisedWith(s, auth, "right") {
		t.Error("an unrecognised scheme accepted a request, which is the fail-open this branch exists to prevent")
	}
}

func sandboxWithAuth(t *testing.T, auth recipe.Auth) *Sandbox {
	t.Helper()

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	// A copy, so the shared embedded Recipe is not altered for anything else
	// running in this process.
	clone := *r
	clone.Auth = auth

	s, err := New(&clone, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	return s
}

func authorisedWith(s *Sandbox, auth recipe.Auth, credential string) bool {
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)

	switch auth.Scheme {
	case "bearer":
		req.Header.Set("Authorization", auth.Prefix+credential)
	case "basic":
		req.SetBasicAuth("user", credential)
	case "header":
		req.Header.Set(auth.Header, credential)
	case "query":
		q := req.URL.Query()
		q.Set(auth.Param, credential)
		req.URL.RawQuery = q.Encode()
	default:
		req.Header.Set("Authorization", credential)
	}

	return s.authorised(req)
}
