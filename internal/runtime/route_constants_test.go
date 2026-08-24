package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route's declared constants have to survive two things that were both
// dropping them, and JSON:API is what needs both at once: it puts a type
// beside every single record, at data.type, with the record itself at data.
//
// The first was a duplicate. get built its response body inline instead of
// going through writeRecord, and writeRecord is where a route's constants are
// applied -- so fields: on a get route was accepted by the validator, written
// into a Recipe, and silently ignored. No shipped Recipe declared any, which
// is why nothing caught it.
//
// The second was a type assertion. store.Record is a named map[string]any and
// asserting the unnamed one does not match it, so setPath could not descend
// into a value that happened to be a record -- and it did not fail politely:
// it replaced the record with a fresh empty map. Applying data.type erased the
// order sitting at data. The runtime already handled this exact trap one layer
// up, in the handler; this is the same trap one layer down.
func TestARoutesConstantsSurviveIntoASingleRecordsEnvelope(t *testing.T) {
	r, err := recipe.Open("lemonsqueezy")
	if err != nil {
		t.Fatalf("open lemonsqueezy: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-store"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/orders/800001", nil)
	req.Header.Set("Authorization", "Bearer eyJ0cGNhdWxkcm9ubGVtb25zcXVlZXp5")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data = %#v, want an object", body["data"])
	}

	// The constant landed.
	if data["type"] != "orders" {
		t.Errorf(`data.type = %v, want "orders" -- a get route's constants were dropped`, data["type"])
	}

	// And it did not cost the record. This is the half that failed loudly
	// once the first was fixed: the whole order vanished and only the
	// constant was left.
	if data["id"] != "800001" {
		t.Errorf("data.id = %v, want 800001 -- the constant replaced the record", data["id"])
	}

	attributes, ok := data["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("data.attributes = %#v, want the record nested under it", data["attributes"])
	}

	if attributes["order_number"] == nil {
		t.Error("data.attributes is empty, so the constant erased what was under data")
	}

	// A constant outside the envelope key was never affected by either bug,
	// and still works.
	jsonapi, ok := body["jsonapi"].(map[string]any)
	if !ok || jsonapi["version"] != "1.0" {
		t.Errorf("jsonapi = %#v, want version 1.0 beside the data", body["jsonapi"])
	}
}
