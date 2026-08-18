package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A listing is usually a summary. Braze's /campaigns/list hands back five
// properties and /campaigns/details hands back a different, larger set, so
// code reading channels off a list entry gets undefined against the real API.
//
// The trim used to live inside the create path, which meant a Recipe could
// declare returns on a get or a list and be describing nothing: the validator
// accepted it, the emulator ignored it, and the full record went out anyway.

func brazeSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("braze")
	if err != nil {
		t.Fatalf("open braze: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-workspace"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func brazeCall(t *testing.T, s *Sandbox, method, target, body string) map[string]any {
	t.Helper()

	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	req.Header.Set("Authorization", "Bearer cauldron-braze-rest-api-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return decode(t, rec)
}

func TestAListingIsTrimmedToTheFieldsItDeclares(t *testing.T) {
	s := brazeSandbox(t)

	body := brazeCall(t, s, http.MethodGet, "/campaigns/list", "")

	campaigns, _ := body["campaigns"].([]any)
	if len(campaigns) != 2 {
		t.Fatalf("want 2 campaigns, got %v", body)
	}

	first, _ := campaigns[0].(map[string]any)
	if first == nil {
		t.Fatalf("first campaign is not an object: %v", campaigns[0])
	}

	// Naming the absent keys is not enough on its own, because a trim that
	// kept one extra field would still satisfy any list of names that did not
	// happen to include it. Five properties, and these five.
	if len(first) != 5 {
		t.Errorf("a listing entry should carry five properties, got %v", first)
	}

	for _, kept := range []string{"id", "name", "is_api_campaign", "tags", "last_edited"} {
		if _, present := first[kept]; !present {
			t.Errorf("%s should be in the listing: %v", kept, first)
		}
	}

	for _, dropped := range []string{"channels", "description", "schedule_type", "created_at", "updated_at", "first_sent", "message"} {
		if _, present := first[dropped]; present {
			t.Errorf("%s belongs to the detail endpoint, not the listing: %v", dropped, first)
		}
	}
}

func TestAGetIsTrimmedToTheFieldsItDeclares(t *testing.T) {
	s := brazeSandbox(t)

	body := brazeCall(t, s, http.MethodGet,
		"/campaigns/details?campaign_id=0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b", "")

	// The detail endpoint has the larger set, and still not the whole record:
	// last_edited and is_api_campaign are listing-only, and a client that saw
	// them here would build on something the provider never sends.
	if body["updated_at"] != "2026-08-04T16:20:11Z" {
		t.Errorf("updated_at = %v", body["updated_at"])
	}

	for _, dropped := range []string{"last_edited", "is_api_campaign"} {
		if _, present := body[dropped]; present {
			t.Errorf("%s belongs to the listing, not the detail: %v", dropped, body)
		}
	}
}

func TestATrimKeepsDeclaredConstants(t *testing.T) {
	s := brazeSandbox(t)

	// Braze stamps "message": "success" on every object, and it is the only
	// field a Braze client reliably checks, because a failure puts its prose
	// in the same place. A trim that dropped constants would remove it.
	body := brazeCall(t, s, http.MethodPost, "/users/export/segment",
		`{"segment_id":"2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d","fields_to_export":["email"]}`)

	if body["message"] != "success" {
		t.Errorf("message = %v, want success", body["message"])
	}

	// A prefix, and nothing that resembles the data. The export is
	// asynchronous and the file arrives somewhere else entirely.
	if prefix, _ := body["object_prefix"].(string); prefix == "" {
		t.Errorf("no object_prefix in %v", body)
	}

	for _, absent := range []string{"users", "data", "segment_id", "fields_to_export"} {
		if _, present := body[absent]; present {
			t.Errorf("%s should not be in an export response: %v", absent, body)
		}
	}
}
