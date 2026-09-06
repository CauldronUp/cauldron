package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Each shipped Recipe is meant to prove a different part of the format is real
// rather than aspirational. These tests assert the parts that differ.

func sandboxFor(t *testing.T, name string) *Sandbox {
	t.Helper()

	r, err := recipe.Open(name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new %s sandbox: %v", name, err)
	}

	return s
}

func request(t *testing.T, s *Sandbox, method, path, body string, decorate func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request

	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if decorate != nil {
		decorate(req)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func shopify(t *testing.T) (*Sandbox, func(*http.Request)) {
	t.Helper()

	return sandboxFor(t, "shopify"), func(r *http.Request) {
		r.Header.Set("X-Shopify-Access-Token", "shpat_cauldron")
	}
}

func twilio(t *testing.T) (*Sandbox, func(*http.Request)) {
	t.Helper()

	return sandboxFor(t, "twilio"), func(r *http.Request) {
		r.SetBasicAuth("ACcauldron00000000000000000000000", "any-token")
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

// Shopify nests lists under a plural collection key, and the key differs per
// resource, so a single recipe-wide setting would not have been enough.
func TestShopifyWrapsListsUnderThePerResourceCollection(t *testing.T) {
	s, auth := shopify(t)

	if err := s.Seed("small-shop"); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	orders := request(t, s, http.MethodGet, "/admin/api/2026-01/orders.json", "", auth)

	if orders.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", orders.Code, orders.Body)
	}

	var wrapped map[string]any
	if err := json.Unmarshal(orders.Body.Bytes(), &wrapped); err != nil {
		t.Fatalf("not an object: %v", err)
	}

	list, ok := wrapped["orders"].([]any)
	if !ok {
		t.Fatalf("expected an orders key, got %v", keysOf(wrapped))
	}

	if len(list) != 2 {
		t.Errorf("got %d orders, want 2", len(list))
	}

	products := request(t, s, http.MethodGet, "/admin/api/2026-01/products.json", "", auth)

	var wrappedProducts map[string]any
	_ = json.Unmarshal(products.Body.Bytes(), &wrappedProducts)

	if _, ok := wrappedProducts["products"].([]any); !ok {
		t.Errorf("products should nest under their own key, got %v", keysOf(wrappedProducts))
	}
}

// The API version is a path parameter that partitions nothing, so any version
// must resolve to the same data.
func TestShopifyVersionSegmentIsNotAScope(t *testing.T) {
	s, auth := shopify(t)
	_ = s.Seed("small-shop")

	for _, version := range []string{"2026-01", "2025-10", "unstable"} {
		rec := request(t, s, http.MethodGet, "/admin/api/"+version+"/orders.json", "", auth)

		if rec.Code != http.StatusOK {
			t.Fatalf("version %s = %d", version, rec.Code)
		}

		var wrapped map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &wrapped)

		if list, _ := wrapped["orders"].([]any); len(list) != 2 {
			t.Errorf("version %s returned %d orders, want 2", version, len(list))
		}
	}
}

func TestShopifyUsesItsOwnAuthHeader(t *testing.T) {
	s, _ := shopify(t)

	// A bearer token is the wrong credential for Shopify.
	rec := request(t, s, http.MethodGet, "/admin/api/2026-01/orders.json", "", func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer shpat_cauldron")
	})

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestShopifyRateLimitCarriesItsOwnHeaders(t *testing.T) {
	s, auth := shopify(t)

	if err := s.Arm(Fault{Error: "rate_limit"}); err != nil {
		t.Fatalf("Arm: %v", err)
	}

	rec := request(t, s, http.MethodGet, "/admin/api/2026-01/orders.json", "", auth)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}

	if got := rec.Header().Get("X-Shopify-Shop-Api-Call-Limit"); got != "40/40" {
		t.Errorf("call limit header = %q", got)
	}
}

// Twilio authenticates with HTTP Basic, where the account SID is the username.
func TestTwilioBasicAuth(t *testing.T) {
	s, auth := twilio(t)

	path := "/2010-04-01/Accounts/ACcauldron00000000000000000000000/Messages.json"

	if rec := request(t, s, http.MethodGet, path, "", auth); rec.Code != http.StatusOK {
		t.Fatalf("valid credentials = %d\n%s", rec.Code, rec.Body)
	}

	wrong := request(t, s, http.MethodGet, path, "", func(r *http.Request) {
		r.SetBasicAuth("ACsomeoneelse", "token")
	})

	if wrong.Code != http.StatusUnauthorized {
		t.Errorf("wrong account = %d, want 401", wrong.Code)
	}

	none := request(t, s, http.MethodGet, path, "", nil)
	if none.Code != http.StatusUnauthorized {
		t.Errorf("no credentials = %d, want 401", none.Code)
	}
}

func TestTwilioScopesMessagesByAccount(t *testing.T) {
	s, auth := twilio(t)

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	rec := request(t, s, http.MethodGet, "/2010-04-01/Accounts/ACcauldron00000000000000000000000/Messages.json", "", auth)

	var wrapped map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &wrapped)

	// Three since small-account grew a message for the paging case. Worth
	// saying while looking at it: every message in that fixture belongs to
	// this one account, so what this asserts is that the scoped listing
	// answers them -- not that another account's are kept out, which it has
	// nothing to exclude.
	if list, _ := wrapped["messages"].([]any); len(list) != 3 {
		t.Fatalf("got %d messages, want 3", len(list))
	}
}

// The failure people cannot stage against a real Twilio sandbox.
func TestTwilioCanFailTheWayRealTwilioDoes(t *testing.T) {
	s, auth := twilio(t)

	if err := s.Arm(Fault{Error: "unreachable_destination"}); err != nil {
		t.Fatalf("Arm: %v", err)
	}

	rec := request(t, s, http.MethodPost,
		"/2010-04-01/Accounts/ACcauldron00000000000000000000000/Messages.json",
		"To=%2B441234567890&From=%2B441234567891&Body=hello", auth)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "21612") {
		t.Errorf("the provider's own error code should be returned\n%s", rec.Body)
	}
}

func TestEveryShippedRecipeBootsAndServes(t *testing.T) {
	for _, name := range recipe.Bundled() {
		t.Run(name, func(t *testing.T) {
			s := sandboxFor(t, name)

			if len(s.Fixtures()) == 0 {
				t.Fatal("no fixtures")
			}

			for _, fixture := range s.Fixtures() {
				if err := s.Seed(fixture); err != nil {
					t.Errorf("seeding %q: %v", fixture, err)
				}
			}

			// Every Recipe must be able to reset cleanly, or it cannot be
			// reused between tests.
			if err := s.Reset(); err != nil {
				t.Errorf("reset: %v", err)
			}

			if len(s.Errors()) == 0 {
				t.Error("a Recipe with no injectable failures cannot test the paths that matter")
			}
		})
	}
}

func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))

	for key := range m {
		out = append(out, key)
	}

	return out
}
