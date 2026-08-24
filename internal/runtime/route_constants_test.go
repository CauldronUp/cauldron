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

// The third place the same assertion had to be fixed, and the one that says
// the pattern is the assertion rather than any caller.
//
// setPath handles an indexed segment so that a Recipe can put a record at
// output.completeTrackResults[0], which is where FedEx puts one: an array per
// tracking number asked about, each holding an array of the shipments that
// number matched. The branch that does it asserted map[string]any on the
// value, and a record is a store.Record, so a record placed at an indexed
// leaf skipped the branch entirely and fell through to the plain assignment
// underneath -- producing a key literally spelled "completeTrackResults[0]",
// brackets and all. A shape no provider sends, produced in silence.
//
// The comment beside that branch already records the same failure happening
// for a declared constant, and calls it the fourth time a key like that had
// been written. This is the fifth, from the same root cause.
func TestARecordMayBePlacedAtAnIndexedLeaf(t *testing.T) {
	r, err := recipe.Open("fedex")
	if err != nil {
		t.Fatalf("open fedex: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/track/v1/trackingnumbers/794832185981", nil)
	req.Header.Set("Authorization", "Bearer l7cauldronfedexaccesstoken0000")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	output, ok := body["output"].(map[string]any)
	if !ok {
		t.Fatalf("output = %#v, want an object", body["output"])
	}

	if _, literal := output["completeTrackResults[0]"]; literal {
		t.Fatal(`the body carries a key literally spelled "completeTrackResults[0]"`)
	}

	results, ok := output["completeTrackResults"].([]any)
	if !ok || len(results) == 0 {
		t.Fatalf("output.completeTrackResults = %#v, want a non-empty array", output["completeTrackResults"])
	}

	first, ok := results[0].(map[string]any)
	if !ok {
		t.Fatalf("output.completeTrackResults[0] = %#v, want the record in it", results[0])
	}

	if first["trackingNumber"] != "794832185981" {
		t.Errorf("output.completeTrackResults[0].trackingNumber = %v, want the record", first["trackingNumber"])
	}
}
