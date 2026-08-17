package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The GitHub Recipe exists to prove the format is not secretly Stripe-shaped.
// These tests assert the two things GitHub does differently.

func github(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github recipe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	return s
}

func ghCall(t *testing.T, s *Sandbox, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	var req *http.Request
	if reader != nil {
		req = httptest.NewRequest(method, path, reader)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	req.Header.Set("Authorization", "Bearer ghp_cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

// A client doing json.Unmarshal into a slice must not receive an object.
func TestBareListStyleReturnsAnArray(t *testing.T) {
	s := github(t)

	if err := s.Seed("small-repo"); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	rec := ghCall(t, s, http.MethodGet, "/repos/octocat/hello-world/issues", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	var issues []map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &issues); err != nil {
		t.Fatalf("a bare list must unmarshal into a slice: %v\n%s", err, rec.Body)
	}

	if len(issues) != 2 {
		t.Fatalf("got %d issues, want 2", len(issues))
	}

	if issues[0]["title"] == nil {
		t.Errorf("issue is missing its title: %v", issues[0])
	}
}

func TestNumericIdentifiersAreSequential(t *testing.T) {
	s := github(t)

	first := ghCall(t, s, http.MethodPost, "/repos/octocat/hello-world/issues", `{"title":"First"}`)
	second := ghCall(t, s, http.MethodPost, "/repos/octocat/hello-world/issues", `{"title":"Second"}`)

	var a, b map[string]any

	_ = json.Unmarshal(first.Body.Bytes(), &a)
	_ = json.Unmarshal(second.Body.Bytes(), &b)

	if a["id"] != "1" || b["id"] != "2" {
		t.Errorf("numeric ids should be sequential; got %v then %v", a["id"], b["id"])
	}

	if strings.Contains(a["id"].(string), "_") {
		t.Errorf("a numeric id must not carry a prefix; got %v", a["id"])
	}
}

func TestNumericIdentifiersRewindOnReset(t *testing.T) {
	s := github(t)

	before := ghCall(t, s, http.MethodPost, "/repos/octocat/hello-world/issues", `{"title":"First"}`)

	if err := s.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}

	after := ghCall(t, s, http.MethodPost, "/repos/octocat/hello-world/issues", `{"title":"First"}`)

	var a, b map[string]any

	_ = json.Unmarshal(before.Body.Bytes(), &a)
	_ = json.Unmarshal(after.Body.Bytes(), &b)

	if a["id"] != b["id"] {
		t.Errorf("after reset id = %v, want %v", b["id"], a["id"])
	}
}

// GitHub rate limits with 403 and its own headers, not Stripe's 429.
func TestErrorsComeFromTheProvidersOwnTable(t *testing.T) {
	s := github(t)

	if err := s.Arm(Fault{Error: "rate_limit"}); err != nil {
		t.Fatalf("Arm: %v", err)
	}

	rec := ghCall(t, s, http.MethodGet, "/repos/octocat/hello-world/issues", "")

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403 — GitHub does not use 429 here", rec.Code)
	}

	if got := rec.Header().Get("x-ratelimit-remaining"); got != "0" {
		t.Errorf("x-ratelimit-remaining = %q, want 0", got)
	}
}

func TestPatchIsRoutedForProvidersThatUseIt(t *testing.T) {
	s := github(t)

	created := ghCall(t, s, http.MethodPost, "/repos/octocat/hello-world/issues", `{"title":"Needs triage"}`)

	var issue map[string]any
	_ = json.Unmarshal(created.Body.Bytes(), &issue)

	rec := ghCall(t, s, http.MethodPatch, "/repos/octocat/hello-world/issues/"+issue["id"].(string), `{"state":"closed"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	var updated map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)

	if updated["state"] != "closed" {
		t.Errorf("state = %v", updated["state"])
	}

	if updated["title"] != "Needs triage" {
		t.Errorf("patch dropped an unrelated field: %v", updated)
	}
}

// Two Recipes with different shapes must be able to run side by side without
// one imposing its conventions on the other.
func TestTwoRecipesCoexistWithDifferentShapes(t *testing.T) {
	gh := github(t)
	st := stripe(t)

	_ = gh.Seed("small-repo")
	_ = st.Seed("small-shop")

	ghBody := ghCall(t, gh, http.MethodGet, "/repos/octocat/hello-world/issues", "").Body.Bytes()
	stBody := call(t, st, http.MethodGet, "/v1/customers", "").Body.Bytes()

	var asSlice []map[string]any
	if err := json.Unmarshal(ghBody, &asSlice); err != nil {
		t.Errorf("github list should be an array: %v", err)
	}

	var asObject map[string]any
	if err := json.Unmarshal(stBody, &asObject); err != nil {
		t.Errorf("stripe list should be an object: %v", err)
	}

	if asObject["object"] != "list" {
		t.Errorf("stripe envelope changed: %v", asObject["object"])
	}
}

// Cloudflare wraps every single object under "result" regardless of what the
// object is, so the wrapper name has to be declarable rather than derived from
// the resource.
func TestASingleResourceCanBeWrappedUnderADeclaredKey(t *testing.T) {
	r, err := recipe.Open("cloudflare")
	if err != nil {
		t.Fatalf("open cloudflare: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/client/v4/zones/00000000000000000000000000000001/dns_records/0000000000000000000000000000000a", nil)
	req.Header.Set("Authorization", "Bearer cauldron-api-token")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := decode(t, rec)

	result, _ := body["result"].(map[string]any)
	if result == nil {
		t.Fatalf("no result in %v", body)
	}

	if result["name"] != "example.com" {
		t.Errorf("result.name = %v", result["name"])
	}

	// Not under the resource's own name, which is what an underived wrapper
	// would have produced.
	if _, wrong := body["dns_record"]; wrong {
		t.Error("the object should not also appear under its resource name")
	}

	if body["success"] != true {
		t.Errorf("success = %v, want the declared envelope flag", body["success"])
	}
}

// Xero answers a request for one invoice with a list of one, so client code
// reads Invoices[0]. An emulator returning the object directly lets code ship
// that breaks against the real API on its first call.
func TestASingleResourceCanComeBackAsAListOfOne(t *testing.T) {
	r, err := recipe.Open("xero")
	if err != nil {
		t.Fatalf("open xero: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-org"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/api.xro/2.0/Invoices/33333333-3333-4333-8333-333333333333", nil)
	req.Header.Set("Authorization", "Bearer cauldron-xero-access-token")
	req.Header.Set("Xero-tenant-id", "cauldron-tenant")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := decode(t, rec)

	// The plural collection name, because a list of one is still a collection.
	invoices, _ := body["Invoices"].([]any)
	if len(invoices) != 1 {
		t.Fatalf("Invoices = %v, want a list of one", body["Invoices"])
	}

	invoice, _ := invoices[0].(map[string]any)
	if invoice["InvoiceNumber"] != "INV-0001" {
		t.Errorf("InvoiceNumber = %v", invoice["InvoiceNumber"])
	}

	// Xero has no "id" anywhere: the identifier is InvoiceID.
	if _, wrong := invoice["id"]; wrong {
		t.Error("the identifier should have been renamed to InvoiceID")
	}

	if _, wrong := body["invoice"]; wrong {
		t.Error("the object should not also appear under its singular name")
	}
}
