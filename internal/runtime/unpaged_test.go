package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A listing that does not page can say so.
//
// Every listing is paged: one with no page size declared is given ten and reads
// "limit", which is a claim about the provider that 344 routes across 222
// Recipes never made. Nothing was truncated by it, because no fixture behind
// one of those holds more than ten records -- which is why it stayed invisible
// and is not a reason it is fine. The claim is about the provider and the
// fixture is not the provider.
//
// Several providers' own descriptions settle it outright. OpenAI's /v1/models,
// xAI's, Perplexity's, Supabase's /v1/projects, Upstash's Redis databases and
// a dozen more declare no query parameters at all: no limit, no offset, no
// cursor. Silence in the Recipe could not distinguish "this provider serves
// the whole collection" from "nobody has looked yet", and those are different
// things to say.
func TestAListingCanDeclareThatItDoesNotPage(t *testing.T) {
	r := unpaged()

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things")

	records, _ := body["things"].([]any)
	if len(records) != 25 {
		t.Errorf("an unpaged listing served %d of 25 records", len(records))
	}

}

// And the position parameter is not read either. A cursor sent against a
// listing that has none is an ordinary ignored parameter at the provider, and
// a fake that treats it as a position answers a short page or refuses the
// request outright.
func TestAnUnpagedListingIgnoresACursor(t *testing.T) {
	s, err := New(unpaged(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things?cursor=thing_10")

	records, _ := body["things"].([]any)
	if len(records) != 25 {
		t.Errorf("a cursor was honoured on an unpaged listing: %d records", len(records))
	}
}

// The page size parameter is not read either, because the provider does not
// have one. A caller sending it against the real API is ignored, and a fake
// that honours it lets a client page locally and discover in production that
// it never was.
func TestAnUnpagedListingIgnoresALimitParameter(t *testing.T) {
	s, err := New(unpaged(), Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things?limit=5")

	records, _ := body["things"].([]any)
	if len(records) != 25 {
		t.Errorf("limit=5 was honoured on an unpaged listing: %d records", len(records))
	}
}

// And a Recipe cannot say both. Declaring no paging beside a page size is two
// claims that contradict each other, and the useful moment to find that is
// before it boots.
func TestNoPagingBesideAPageSizeIsRefused(t *testing.T) {
	r := unpaged()
	r.Routes[0].Pagination.Limit = 50

	if err := r.Validate(); err == nil {
		t.Fatal("a route declaring no paging and a page size validated")
	}

	r = unpaged()
	r.Routes[0].Pagination.CursorParam = "after"

	if err := r.Validate(); err == nil {
		t.Fatal("a route declaring no paging and a cursor parameter validated")
	}
}

func listThings(t *testing.T, s *Sandbox, target string) map[string]any {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer the-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	return body
}

func seed(t *testing.T, s *Sandbox, _ int) {
	t.Helper()

	if err := s.Seed("many"); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func manyThings(n int) []map[string]any {
	records := make([]map[string]any, 0, n)

	for i := range n {
		records = append(records, map[string]any{
			"id":   fmt.Sprintf("thing_%02d", i),
			"name": fmt.Sprintf("thing %d", i),
		})
	}

	return records
}

func unpaged() *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "unpaged",
		Capability: "ai",
		Version:    "0.1.0",
		Upstream:   recipe.Upstream{API: "v1"},
		Auth: recipe.Auth{
			Scheme: "bearer",
			Prefix: "Bearer ",
			Keys:   []string{"the-key"},
		},
		Responses: recipe.Responses{
			List: recipe.ListResponse{Style: "wrapped"},
		},
		Resources: map[string]recipe.Resource{
			"thing": {
				Collection: "things",
				ID:         recipe.ID{Style: "opaque"},
				Fields: map[string]recipe.Field{
					"id":   {Type: "string"},
					"name": {Type: "string"},
				},
			},
		},
		Fixtures: map[string]recipe.Fixture{"many": {"thing": manyThings(25)}},
		Routes: []recipe.Route{
			{
				Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list",
				Pagination: recipe.Pagination{Style: "none"},
			},
		},
	}
}
