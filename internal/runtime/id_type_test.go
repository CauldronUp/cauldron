package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Identifiers are minted and stored as strings, because that is the only form
// every style shares and the only form a path parameter arrives in. What they
// are on the wire is a separate question and the two answers disagree more
// often than they agree: GitHub sends an issue id as the number 1, HubSpot
// sends a contact id as the string "1".
//
// It is not cosmetic. id === 1 fails against a string, typeof id === "number"
// fails, and a schema declaring an integer rejects a quoted one outright. An
// emulator answering with a string where the provider answers with a number
// commits the exact class of bug it exists to catch, and Cauldron did for
// every Recipe until this existed.

func meiliBody(t *testing.T, target string) map[string]any {
	t.Helper()

	r, err := recipe.Open("meilisearch")
	if err != nil {
		t.Fatalf("open meilisearch: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-instance"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer cauldron-meilisearch-master-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	return decode(t, rec)
}

func TestANumericIdentifierIsAJSONNumber(t *testing.T) {
	body := meiliBody(t, "/tasks/42")

	// json.Unmarshal into any gives float64 for every JSON number, so the
	// type assertion is the whole assertion: a string would fail it.
	uid, ok := body["uid"].(float64)
	if !ok {
		t.Fatalf("uid should be a JSON number, got %T (%v)", body["uid"], body["uid"])
	}

	if uid != 42 {
		t.Errorf("uid = %v, want 42", uid)
	}
}

func TestANumericIdentifierIsANumberInAListingToo(t *testing.T) {
	// The list path renames and retypes through the same presentation pass,
	// and it would be easy for one to be handled and not the other.
	body := meiliBody(t, "/tasks")

	results, _ := body["results"].([]any)
	if len(results) == 0 {
		t.Fatalf("no results in %v", body)
	}

	for i, record := range results {
		object, _ := record.(map[string]any)
		if object == nil {
			t.Fatalf("result %d is not an object", i)
		}

		if _, ok := object["uid"].(float64); !ok {
			t.Errorf("result %d uid should be a JSON number, got %T", i, object["uid"])
		}
	}
}

func TestAnUnrenamedNumericIdentifierIsRetypedInAListing(t *testing.T) {
	// The task listing renames uid as well as retyping it, so it reaches the
	// presentation pass through the rename check and says nothing about
	// whether a retype alone is enough to get there. A document keeps the
	// name "id" and only changes type, which is the case that would have
	// slipped through.
	body := meiliBody(t, "/indexes/movies/documents")

	results, _ := body["results"].([]any)
	if len(results) == 0 {
		t.Fatalf("no results in %v", body)
	}

	first, _ := results[0].(map[string]any)
	if first == nil {
		t.Fatalf("not an object: %v", results[0])
	}

	if _, ok := first["id"].(float64); !ok {
		t.Errorf("a document id should be a JSON number, got %T (%v)", first["id"], first["id"])
	}
}

func TestAStringIdentifierStaysAString(t *testing.T) {
	// The default is unchanged, because changing it silently would rewrite the
	// wire shape of every shipped Recipe at once.
	body := meiliBody(t, "/indexes/movies")

	if _, ok := body["uid"].(string); !ok {
		t.Errorf("an index uid should stay a string, got %T (%v)", body["uid"], body["uid"])
	}
}

func TestANumericIdentifierIsStillAddressableByPath(t *testing.T) {
	// The store keeps the string form, so the conversion has to happen at the
	// edge or the lookup that found this record would stop working.
	body := meiliBody(t, "/tasks/43")

	if body["status"] != "failed" {
		t.Errorf("wrong record: %v", body)
	}
}

func TestANumericIdentifierSurvivesAnEncodeDecodeRoundTrip(t *testing.T) {
	body := meiliBody(t, "/tasks/41")

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Not "uid":"41". A quoted identifier is what a schema rejects.
	if got := string(encoded); !strings.Contains(got, `"uid":41`) {
		t.Errorf("uid should be unquoted on the wire: %s", got)
	}
}
