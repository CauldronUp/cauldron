package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Absent and null are different on the wire, and providers disagree about
// which they use. Bandwidth leaves errorCode out of the messages that worked,
// so a client testing for null never sees one. Alpaca sends every timestamp on
// every order and leaves the ones that have not happened as null, so a client
// testing for the key's existence finds it and reads nothing out of it. Each
// of those is a real bug in code written against the other assumption, so the
// format has to be able to say both.

func alpacaCall(t *testing.T, target, body string) map[string]any {
	t.Helper()

	r, err := recipe.Open("alpaca")
	if err != nil {
		t.Fatalf("open alpaca: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	method := http.MethodGet

	var req *http.Request
	if body != "" {
		method = http.MethodPost
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	req.Header.Set("APCA-API-KEY-ID", "cauldron-alpaca-key-id")
	req.Header.Set("APCA-API-SECRET-KEY", "cauldron-alpaca-secret-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return decode(t, rec)
}

func TestAnUnsetFieldCanBePresentAndNull(t *testing.T) {
	body := alpacaCall(t, "/v2/orders",
		`{"symbol":"TSLA","qty":"4","side":"buy","type":"market","order_type":"market","time_in_force":"day"}`)

	for _, name := range []string{"filled_at", "canceled_at", "expired_at"} {
		value, present := body[name]
		if !present {
			t.Errorf("%s should be on the wire as null, not left out: %v", name, body)

			continue
		}

		// Not a stamped timestamp. An emulator that invented one would report
		// a fill time for an order that has not filled.
		if value != nil {
			t.Errorf("%s = %v, want null", name, value)
		}
	}
}

func alpacaOrders(t *testing.T, target string) []any {
	t.Helper()

	r, err := recipe.Open("alpaca")
	if err != nil {
		t.Fatalf("open alpaca: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("APCA-API-KEY-ID", "cauldron-alpaca-key-id")
	req.Header.Set("APCA-API-SECRET-KEY", "cauldron-alpaca-secret-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out []any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("not a JSON array: %v %s", err, rec.Body.String())
	}

	return out
}

func orderStatuses(t *testing.T, records []any) []string {
	t.Helper()

	out := make([]string, 0, len(records))

	for _, record := range records {
		object, _ := record.(map[string]any)
		if object == nil {
			t.Fatalf("not an object: %v", record)
		}

		status, _ := object["status"].(string)
		out = append(out, status)
	}

	return out
}

func TestAFilterValueCanCoverSeveralFieldValues(t *testing.T) {
	// "open" is not a status any order holds. new and partially_filled are
	// open, and a filter matching the word literally would hide the partially
	// filled order, which is the one that most needs to be visible because it
	// is a real position.
	statuses := orderStatuses(t, alpacaOrders(t, "/v2/orders"))

	if len(statuses) != 2 {
		t.Fatalf("the default listing should be the two open orders, got %v", statuses)
	}

	want := map[string]bool{"partially_filled": true, "new": true}
	for _, status := range statuses {
		if !want[status] {
			t.Errorf("%s is not an open status: %v", status, statuses)
		}
	}
}

func TestTheOtherBucketIsTheOtherHalf(t *testing.T) {
	statuses := orderStatuses(t, alpacaOrders(t, "/v2/orders?status=closed"))

	if len(statuses) != 1 || statuses[0] != "filled" {
		t.Errorf("status=closed should be the filled order alone, got %v", statuses)
	}
}
