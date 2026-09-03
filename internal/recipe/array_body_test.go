package recipe

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

// A conformance case can send a JSON array, because a JSON body is not always
// an object.
//
// Mixpanel's track and import endpoints declare an array with at least one
// member, and every batch ingest API in this collection is that shape. While
// this field was a mapping, the one behaviour those Recipes exist to describe
// was the one thing no case could send: a Recipe could assert that a bare
// object is refused and could not assert that the array its own clients post is
// accepted.
func TestACaseCanSendAJSONArrayBody(t *testing.T) {
	var r Request

	const document = `
method: POST
path: /track
json:
  - event: a-name
    properties:
      token: a-token
`

	if err := yaml.Unmarshal([]byte(document), &r); err != nil {
		t.Fatalf("an array body would not decode: %v", err)
	}

	if !r.SendsJSON() {
		t.Fatal("an array body does not count as a body")
	}

	encoded, err := json.Marshal(r.JSON)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(encoded); got[0] != '[' {
		t.Errorf("the body went out as %s, want a JSON array", got)
	}

	// An array has no top-level field names, so nothing is echoed from it.
	if fields := r.JSONFields(); fields != nil {
		t.Errorf("an array reported object fields: %v", fields)
	}
}

// An object body still decodes, still counts, and still reports its fields --
// which is what keeps the loosening from being a loss.
func TestAnObjectBodyStillCarriesItsFields(t *testing.T) {
	var r Request

	const document = `
method: POST
path: /v1/things
json:
  name: a-name
  count: 2
`

	if err := yaml.Unmarshal([]byte(document), &r); err != nil {
		t.Fatalf("an object body would not decode: %v", err)
	}

	if !r.SendsJSON() {
		t.Fatal("an object body does not count as a body")
	}

	fields := r.JSONFields()
	if fields["name"] != "a-name" {
		t.Errorf("fields = %v, want the name it was given", fields)
	}
}

// No body at all is still no body. The distinction matters because a provider
// can answer differently to an empty document than to nothing, and Coveralls
// already does.
func TestAnAbsentBodyIsNotABody(t *testing.T) {
	var r Request

	if err := yaml.Unmarshal([]byte("method: GET\npath: /v1/things\n"), &r); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if r.SendsJSON() {
		t.Error("a request with no json key reported a body")
	}
}

// An empty array and an empty object are both bodies, and deliberately so.
func TestAnEmptyBodyIsStillABody(t *testing.T) {
	for _, document := range []string{"json: []\n", "json: {}\n"} {
		var r Request

		if err := yaml.Unmarshal([]byte(document), &r); err != nil {
			t.Fatalf("%q: %v", document, err)
		}

		if !r.SendsJSON() {
			t.Errorf("%q reported no body, and sending it is a claim", document)
		}
	}
}

// The body marker scan reads inside an array, so a route selected by something
// in the payload is still reachable from an array-bodied case.
func TestTheBodyMarkerScanReadsInsideAnArray(t *testing.T) {
	cases := []Case{{
		Request: Request{
			Method: "POST",
			Path:   "/track",
			JSON:   []any{map[string]any{"event": "a-name"}},
		},
	}}

	if !sendsBodyMarker(cases, "event") {
		t.Error("a marker inside an array body was not found")
	}

	if sendsBodyMarker(cases, "not-in-there") {
		t.Error("a marker that is not in the body was reported as present")
	}
}
