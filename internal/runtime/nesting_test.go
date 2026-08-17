package runtime

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// HubSpot puts every business attribute under "properties" and sends RFC 3339
// timestamps. Both are presentation concerns: the store stays flat and keyed
// by id, and only the wire shape changes.

func hubspotSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("hubspot")
	if err != nil {
		t.Fatalf("open hubspot: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-portal"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func hubspotCall(t *testing.T, s *Sandbox, method, target, body string) map[string]any {
	t.Helper()

	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	var req *http.Request
	if reader != nil {
		req = httptest.NewRequest(method, target, reader)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	req.Header.Set("Authorization", "Bearer pat-na1-cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return decode(t, rec)
}

func TestDeclaredFieldsAreNestedOnTheWire(t *testing.T) {
	s := hubspotSandbox(t)

	body := hubspotCall(t, s, http.MethodGet, "/crm/v3/objects/contacts/1", "")

	properties, _ := body["properties"].(map[string]any)
	if properties == nil {
		t.Fatalf("no properties in %v", body)
	}

	if properties["email"] != "ada@example.com" {
		t.Errorf("properties.email = %v", properties["email"])
	}

	// A client reading contact.email finds nothing, because HubSpot does not
	// put it there. The emulator must not be more generous than the provider.
	if _, flat := body["email"]; flat {
		t.Error("email should not also appear at the top level")
	}

	// id, archived and the timestamps stay where HubSpot keeps them.
	for _, top := range []string{"id", "archived", "createdAt"} {
		if _, present := body[top]; !present {
			t.Errorf("%s should stay at the top level", top)
		}
	}
}

func TestNestedFieldsAreFlattenedOnTheWayIn(t *testing.T) {
	s := hubspotSandbox(t)

	created := hubspotCall(t, s, http.MethodPost, "/crm/v3/objects/contacts",
		`{"properties":{"email":"alan@example.com","firstname":"Alan"}}`)

	properties, _ := created["properties"].(map[string]any)
	if properties == nil {
		t.Fatalf("no properties in %v", created)
	}

	if properties["email"] != "alan@example.com" {
		t.Errorf("properties.email = %v", properties["email"])
	}

	// A declared default applies to a nested field like any other, which only
	// works because the store sees it flat.
	if properties["lifecyclestage"] != "lead" {
		t.Errorf("lifecyclestage = %v, want the declared default", properties["lifecyclestage"])
	}

	// And it round-trips: fetching it back gives the same shape.
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("no id")
	}

	fetched := hubspotCall(t, s, http.MethodGet, "/crm/v3/objects/contacts/"+id, "")

	refetched, _ := fetched["properties"].(map[string]any)
	if refetched == nil || refetched["email"] != "alan@example.com" {
		t.Errorf("round trip lost the value: %v", fetched)
	}
}

func TestARequiredNestedFieldIsStillRequired(t *testing.T) {
	s := hubspotSandbox(t)

	body := hubspotCall(t, s, http.MethodPost, "/crm/v3/objects/contacts",
		`{"properties":{"firstname":"Nameless"}}`)

	// email is required and nested. Flattening happens before the check, so a
	// nested omission is caught rather than silently accepted.
	if message, _ := body["message"].(string); !strings.Contains(message, "email") {
		t.Errorf("body = %v, want a complaint naming email", body)
	}
}

var rfc3339 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func TestDatetimeFieldsAreStringsFromTheSandboxClock(t *testing.T) {
	s := hubspotSandbox(t)

	created := hubspotCall(t, s, http.MethodPost, "/crm/v3/objects/contacts",
		`{"properties":{"email":"time@example.com"}}`)

	stamp, _ := created["createdAt"].(string)
	if !rfc3339.MatchString(stamp) {
		t.Errorf("createdAt = %v, want an RFC 3339 string", created["createdAt"])
	}

	// The sandbox clock, not the wall clock: the same call on any machine at
	// any time produces the same stamp.
	if stamp != "2026-01-01T00:00:00Z" {
		t.Errorf("createdAt = %q, want the sandbox epoch", stamp)
	}
}

func TestANestedCursorAppearsUnderItsDeclaredPath(t *testing.T) {
	s := hubspotSandbox(t)

	body := hubspotCall(t, s, http.MethodGet, "/crm/v3/objects/contacts?limit=1", "")

	paging, _ := body["paging"].(map[string]any)
	if paging == nil {
		t.Fatalf("no paging in %v", body)
	}

	next, _ := paging["next"].(map[string]any)
	if next == nil || next["after"] == nil {
		t.Errorf("paging.next.after missing from %v", body)
	}
}

// A declared envelope constant and a computed value can land in the same
// nested object. Intercom's pages carries both a declared type and a next
// cursor, and whichever was written second used to erase the other.
func TestDeclaredConstantsMergeWithComputedValues(t *testing.T) {
	r, err := recipe.Open("intercom")
	if err != nil {
		t.Fatalf("open intercom: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-workspace"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/contacts?limit=2", nil)
	req.Header.Set("Authorization", "Bearer dG9rOmNhdWxkcm9u")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := decode(t, rec)

	pages, _ := body["pages"].(map[string]any)
	if pages == nil {
		t.Fatalf("no pages in %v", body)
	}

	if pages["type"] != "pages" {
		t.Errorf("the declared constant was lost: pages = %v", pages)
	}

	next, _ := pages["next"].(map[string]any)
	if next == nil || next["starting_after"] == nil {
		t.Errorf("the computed cursor was lost: pages = %v", pages)
	}
}

// Contentful keeps the identifier at sys.id rather than at the top level, so
// the rename has to nest the same way a declared field does. Code reading
// entry.id against the real API finds nothing.
func TestTheIdentifierCanBeRenamedIntoANestedObject(t *testing.T) {
	r, err := recipe.Open("contentful")
	if err != nil {
		t.Fatalf("open contentful: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-space"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/spaces/cauldron/environments/master/entries/entry000000000000000001", nil)
	req.Header.Set("Authorization", "Bearer cauldron-delivery-token")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := decode(t, rec)

	sys, _ := body["sys"].(map[string]any)
	if sys == nil {
		t.Fatalf("no sys block in %v", body)
	}

	if sys["id"] != "entry000000000000000001" {
		t.Errorf("sys.id = %v", sys["id"])
	}

	// The renamed identifier and the declared sys fields share one object,
	// so neither may erase the other.
	if sys["revision"] == nil || sys["type"] != "Entry" {
		t.Errorf("the declared sys fields were lost: %v", sys)
	}

	if _, flat := body["id"]; flat {
		t.Error("the identifier should not also appear at the top level")
	}

	if _, literal := body["sys.id"]; literal {
		t.Error("a dotted field name should nest, not become a literal key")
	}
}
