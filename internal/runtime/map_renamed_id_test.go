package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A keyed list keys on the identifier's real name, not on "id".
//
// The map style keys entries by the identifier so it does not repeat inside the
// value, and it read record["id"] literally. A resource declaring id.field has
// already had that key renamed by the time the list is shaped, so the lookup
// found nothing -- and finding nothing meant skipping the record, so the entry
// disappeared from the response entirely.
//
// That is the worst shape this bug could take. Nothing errored, nothing logged,
// and the caller got a smaller object rather than a failure. Electricity Maps
// found it: its zones are keyed by a zoneKey, and declaring that as the
// identifier emptied the response.
func TestAKeyedListKeysOnTheRenamedIdentifier(t *testing.T) {
	r := keyedZones()
	r.Resources["zone"] = recipe.Resource{
		ID:     recipe.ID{Field: "zoneKey", Style: "opaque"},
		Fields: map[string]recipe.Field{"zoneKey": {Type: "string"}, "carbon": {Type: "number"}},
	}

	body := zonesFrom(t, r)

	entry, ok := body["DE"].(map[string]any)
	if !ok {
		t.Fatalf("the zone keyed by its renamed identifier is missing entirely: %v", body)
	}

	if entry["carbon"] == nil {
		t.Error("the entry is there and its fields are not")
	}

	// The key does not repeat inside the value, which is the whole reason this
	// style exists.
	if _, repeated := entry["zoneKey"]; repeated {
		t.Error("the identifier repeats inside the value it keys")
	}
}

// The ordinary case still works, which is what keeps the fix from being a
// rename of the bug.
func TestAKeyedListStillKeysOnAPlainIdentifier(t *testing.T) {
	r := keyedZones()
	r.Resources["zone"] = recipe.Resource{
		ID:     recipe.ID{Style: "opaque"},
		Fields: map[string]recipe.Field{"id": {Type: "string"}, "carbon": {Type: "number"}},
	}

	body := zonesFrom(t, r)

	entry, ok := body["DE"].(map[string]any)
	if !ok {
		t.Fatalf("a plainly identified zone is missing: %v", body)
	}

	if _, repeated := entry["id"]; repeated {
		t.Error("the identifier repeats inside the value it keys")
	}
}

func zonesFrom(t *testing.T, r *recipe.Recipe) map[string]any {
	t.Helper()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("zones"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v3/zones", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not an object: %v\n%s", err, w.Body.String())
	}

	return body
}

// A Recipe whose listing is keyed rather than ordered.
func keyedZones() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "keyedzones",
		Capability: "infrastructure",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v3"},
		Auth:       recipe.Auth{Scheme: "none"},
		Responses:  recipe.Responses{List: recipe.ListResponse{Style: "map", Key: "-"}},
		Resources:  map[string]recipe.Resource{},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v3/zones", Resource: "zone", Operation: "list"},
		},
		Fixtures: map[string]recipe.Fixture{
			"zones": {"zone": {{"id": "DE", "carbon": 350}}},
		},
	}
}
