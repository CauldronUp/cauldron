package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A resource addressed by a key the provider never echoes. Marqeta's balance
// is fetched at /v3/balances/{token} and the body that comes back mentions no
// token at all, so a balance held on its own cannot say whose it is. Cauldron
// still keys the record internally, because it has to be found somehow, and
// id.field "-" keeps that key off the wire.

func marqetaSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("marqeta")
	if err != nil {
		t.Fatalf("open marqeta: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-program"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func marqetaGet(t *testing.T, s *Sandbox, target string) map[string]any {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.SetBasicAuth("cauldron_marqeta_app_token", "cauldron_marqeta_admin_access_token")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return decode(t, rec)
}

func TestAHiddenIdentifierIsStillAddressable(t *testing.T) {
	s := marqetaSandbox(t)

	body := marqetaGet(t, s, "/v3/balances/7fd2a640-b5c3-4e19-a80f-3c62d94b1e57")

	gpa, _ := body["gpa"].(map[string]any)
	if gpa == nil {
		t.Fatalf("no gpa in %v", body)
	}

	// The right record came back, so suppressing the identifier did not
	// break the lookup that uses it.
	if gpa["available_balance"] != 41.25 {
		t.Errorf("gpa.available_balance = %v, want 41.25", gpa["available_balance"])
	}

	// Naming the absent keys is not enough: a suppression implemented as a
	// rename to some other string would leave that string on the wire and
	// still satisfy a list of names nobody thought to include. The balance
	// carries one property, so say so.
	if len(body) != 1 {
		t.Errorf("balance should carry gpa and nothing else, got %v", body)
	}

	for _, name := range []string{"id", "token", "user_token"} {
		if _, present := body[name]; present {
			t.Errorf("%s should not be on the wire: %v", name, body)
		}
	}
}

func TestARenamedIdentifierStillAppears(t *testing.T) {
	s := marqetaSandbox(t)

	body := marqetaGet(t, s, "/v3/transactions/6f1c0a34-9e2b-4f57-8f3a-2b1d5c7e0a11")

	// The same presentation pass handles both, and suppressing one resource's
	// identifier must not suppress another's.
	if body["token"] != "6f1c0a34-9e2b-4f57-8f3a-2b1d5c7e0a11" {
		t.Errorf("token = %v", body["token"])
	}

	if _, present := body["id"]; present {
		t.Errorf("id should have been renamed away: %v", body)
	}
}
