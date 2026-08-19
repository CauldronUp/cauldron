package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A listing reached by POST usually carries its paging in the JSON body, and
// for as long as the parameters were read from the query string it was never
// there. What that produced was worse than an error: the caller's limit fell
// back to the route's default, the default was larger than any fixture, so the
// first response held the whole collection and reported no next page. A paging
// loop written against it ran exactly once, took neither branch, and passed.
//
// Dropbox shipped that way, and its own conformance case sent ?limit=1 -- a
// parameter Dropbox does not read -- because the case was written against what
// came out rather than against the provider.

func listFolder(t *testing.T, body string) map[string]any {
	t.Helper()

	r, err := recipe.Open("dropbox")
	if err != nil {
		t.Fatalf("open dropbox: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/2/files/list_folder", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer sl.cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v %s", err, rec.Body.String())
	}

	return out
}

func entryCount(t *testing.T, out map[string]any) int {
	t.Helper()

	entries, ok := out["entries"].([]any)
	if !ok {
		t.Fatalf("entries is not an array: %T", out["entries"])
	}

	return len(entries)
}

func TestALimitInTheBodyIsHonoured(t *testing.T) {
	out := listFolder(t, `{"path":"/documents","limit":1}`)

	if n := entryCount(t, out); n != 1 {
		t.Errorf("asked for one entry and got %d", n)
	}

	if out["has_more"] != true {
		t.Error("a truncated listing must say there is more")
	}

	if cursor, _ := out["cursor"].(string); cursor == "" {
		t.Error("a truncated listing must carry the cursor to resume from")
	}
}

func TestACursorInTheBodyFetchesTheNextPage(t *testing.T) {
	first := listFolder(t, `{"path":"/documents","limit":1}`)

	cursor, _ := first["cursor"].(string)
	if cursor == "" {
		t.Fatal("no cursor to resume from")
	}

	second := listFolder(t, `{"path":"/documents","limit":1,"cursor":"`+cursor+`"}`)

	if n := entryCount(t, second); n != 1 {
		t.Fatalf("the second page holds %d entries", n)
	}

	firstName := first["entries"].([]any)[0].(map[string]any)["name"]
	secondName := second["entries"].([]any)[0].(map[string]any)["name"]

	if firstName == secondName {
		t.Errorf("the cursor was ignored and page one came back again: %v", secondName)
	}

	if second["has_more"] != false {
		t.Error("the last page must say there is no more, or a loop never ends")
	}
}

func TestAQueryParameterTheProviderDoesNotReadIsInert(t *testing.T) {
	// The faithful half of the same change. Dropbox does not accept ?limit,
	// and an emulator that quietly honours it lets code that sends the wrong
	// thing work locally and fail in production -- which is what the Recipe's
	// own case was doing.
	r, err := recipe.Open("dropbox")
	if err != nil {
		t.Fatalf("open dropbox: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/2/files/list_folder?limit=1",
		strings.NewReader(`{"path":"/documents"}`))
	req.Header.Set("Authorization", "Bearer sl.cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if n := entryCount(t, out); n < 2 {
		t.Errorf("a query parameter Dropbox does not read was honoured anyway: %d entries", n)
	}
}

func TestANestedBodyParameterIsRead(t *testing.T) {
	// Plaid keeps count and offset inside options, so the name that finds
	// them is a dotted one.
	r, err := recipe.Open("plaid")
	if err != nil {
		t.Fatalf("open plaid: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-item"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/transactions/get",
		strings.NewReader(`{"access_token":"access-sandbox-cauldron","options":{"count":2,"offset":2}}`))
	req.Header.Set("PLAID-CLIENT-ID", "cauldron")
	req.Header.Set("PLAID-SECRET", "cauldron-sandbox-secret")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v %s", err, rec.Body.String())
	}

	transactions, ok := out["transactions"].([]any)
	if !ok {
		t.Fatalf("transactions is not an array: %s", rec.Body.String())
	}

	// Three in the fixture, two skipped, one left.
	if len(transactions) != 1 {
		t.Fatalf("options.offset was not read: %d transactions on the third page", len(transactions))
	}
}
