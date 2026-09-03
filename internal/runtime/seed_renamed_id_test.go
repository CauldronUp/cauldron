package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A fixture may write the identifier under the name the provider uses for it.
//
// The store keys on "id" and knows nothing about that rename, so a record
// carrying only the renamed field was handed a freshly minted identifier while
// its own was kept as an ordinary value. The record then served one identifier
// and was addressable by another, and nothing reported it.
//
// Canvas's key set found it: keyed by kid, the fixture's real value never
// became the key, and the workaround was writing the same value twice in a file
// that should say it once.
func TestAFixtureCanNameTheIdentifierTheWayTheProviderDoes(t *testing.T) {
	s, err := New(renamedID("kid"), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("keys"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Served under the provider's own name, holding the fixture's own value.
	body := listOf(t, s, "/v1/keys")
	if got := body[0]["kid"]; got != "a-real-key-id" {
		t.Errorf("served kid %v, want the one the fixture pinned", got)
	}

	// And addressable by it, which is the half that was silently wrong: the
	// record used to answer to a minted identifier nobody wrote down.
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/keys/a-real-key-id", nil))

	if w.Code != http.StatusOK {
		t.Errorf("fetching by the pinned identifier answered %d, want 200: %s", w.Code, w.Body.String())
	}
}

// Writing both still works, because eleven Recipes already do and their fixtures
// are not wrong -- just repetitive.
func TestAFixtureMayStillWriteBothNames(t *testing.T) {
	r := renamedID("kid")
	r.Fixtures["keys"] = recipe.Fixture{
		"key": {{"id": "a-real-key-id", "kid": "a-real-key-id", "alg": "RS256"}},
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("keys"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if got := listOf(t, s, "/v1/keys")[0]["kid"]; got != "a-real-key-id" {
		t.Errorf("served kid %v", got)
	}
}

// An ordinary resource is untouched, which is what keeps the fix from reaching
// past the case it was written for.
func TestAPlainlyIdentifiedFixtureIsUnaffected(t *testing.T) {
	r := renamedID("")
	r.Fixtures["keys"] = recipe.Fixture{
		"key": {{"id": "k_1", "alg": "RS256"}},
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("keys"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if got := listOf(t, s, "/v1/keys")[0]["id"]; got != "k_1" {
		t.Errorf("served id %v, want k_1", got)
	}
}

func listOf(t *testing.T, s *Sandbox, path string) []map[string]any {
	t.Helper()

	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

	if w.Code != http.StatusOK {
		t.Fatalf("%s answered %d: %s", path, w.Code, w.Body.String())
	}

	// Bare or enveloped, whichever this Recipe's defaults produce -- the claim
	// here is about the identifier, not about the wrapper around it.
	var body []map[string]any

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		var enveloped struct {
			Data []map[string]any `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &enveloped); err != nil {
			t.Fatalf("%s did not answer a list: %v\n%s", path, err, w.Body.String())
		}

		body = enveloped.Data
	}

	if len(body) == 0 {
		t.Fatalf("%s answered an empty list", path)
	}

	return body
}

// A Recipe whose identifier goes out under a name of the provider's choosing.
func renamedID(field string) *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "renamedid",
		Capability: "auth",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth:       recipe.Auth{Scheme: "none"},
		Responses:  recipe.Responses{List: recipe.ListResponse{Key: "-"}},
		Resources: map[string]recipe.Resource{
			"key": {
				ID:     recipe.ID{Field: field, Style: "opaque"},
				Fields: map[string]recipe.Field{"alg": {Type: "string"}},
			},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/keys", Resource: "key", Operation: "list"},
			{Method: "GET", Path: "/v1/keys/{id}", Resource: "key", Operation: "get"},
		},
		Fixtures: map[string]recipe.Fixture{
			"keys": {"key": {{"kid": "a-real-key-id", "alg": "RS256"}}},
		},
	}
}
