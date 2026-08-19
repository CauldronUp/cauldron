package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A GraphQL API is one path and one method, so the path cannot say which route
// should answer. The query can.
//
// Seven providers were unreachable without this -- Linear, Monday, Attio, New
// Relic, Railway, ShipHero and half of Fly.io -- and each had been recorded as
// its own judgement call rather than as one missing feature.

func shipheroQuery(t *testing.T, query string) map[string]any {
	t.Helper()

	r, err := recipe.Open("shiphero")
	if err != nil {
		t.Fatalf("open shiphero: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"query": ` + quote(query) + `}`

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer cauldron.shiphero.jwt.000000")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("not JSON: %v %s", err, rec.Body.String())
	}

	return out
}

func quote(s string) string {
	encoded, _ := json.Marshal(s)

	return string(encoded)
}

func dataOf(t *testing.T, body map[string]any) map[string]any {
	t.Helper()

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("no data envelope: %v", body)
	}

	return data
}

func TestTheQueryChoosesTheRoute(t *testing.T) {
	orders := dataOf(t, shipheroQuery(t, "query { orders { data { edges { node { id } } } } }"))
	if _, ok := orders["orders"]; !ok {
		t.Errorf("a query naming orders did not get the orders route: %v", orders)
	}

	if _, ok := orders["products"]; ok {
		t.Errorf("the orders query also answered with products: %v", orders)
	}

	products := dataOf(t, shipheroQuery(t, "query { products { data { edges { node { sku } } } } }"))
	if _, ok := products["products"]; !ok {
		t.Errorf("a query naming products did not get the products route: %v", products)
	}

	if _, ok := products["orders"]; ok {
		t.Errorf("the products query also answered with orders: %v", products)
	}
}

func TestEachRouteCarriesItsOwnEnvelopeConstants(t *testing.T) {
	// ShipHero puts request_id and complexity beside each connection, so the
	// key depends on which query was asked. A Recipe-wide constant would
	// stamp the orders metadata onto a products response.
	orders := dataOf(t, shipheroQuery(t, "query { orders { complexity } }"))

	connection, ok := orders["orders"].(map[string]any)
	if !ok {
		t.Fatalf("no orders connection: %v", orders)
	}

	if connection["complexity"] != float64(12) {
		t.Errorf("orders complexity is %v, want 12", connection["complexity"])
	}

	products := dataOf(t, shipheroQuery(t, "query { products { complexity } }"))

	other, ok := products["products"].(map[string]any)
	if !ok {
		t.Fatalf("no products connection: %v", products)
	}

	if other["complexity"] != float64(4) {
		t.Errorf("products complexity is %v, want 4", other["complexity"])
	}
}

func TestAQueryThatSelectsNothingKnownIsNotServed(t *testing.T) {
	// Every route on this path declares what it answers, so a query naming
	// none of them has no route. Answering anyway would be the helpful lie.
	r, err := recipe.Open("shiphero")
	if err != nil {
		t.Fatalf("open shiphero: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/graphql",
		strings.NewReader(`{"query": "query { warehouses { id } }"}`))
	req.Header.Set("Authorization", "Bearer cauldron.shiphero.jwt.000000")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Errorf("an unmodelled query was answered: %s", rec.Body.String())
	}
}

func TestRoutingOnThePathAloneStillWorks(t *testing.T) {
	// 158 of the 159 Recipes here route on the path and must keep doing so,
	// and reading the body to look for a query must not disturb the handlers
	// that read it again afterwards.
	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-repo"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/repos/octocat/hello-world/issues",
		strings.NewReader(`{"title": "Something is wrong"}`))
	req.Header.Set("Authorization", "Bearer ghp_cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var created map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("not JSON: %v", err)
	}

	if created["title"] != "Something is wrong" {
		t.Errorf("the body was consumed before the handler read it: %v", created)
	}
}

func TestASelectorMatchesAWholeFieldOnly(t *testing.T) {
	// A substring match is not good enough. selects "me" matched
	// `query { viewer { name email } }`, because "name" contains it, and
	// short root fields are common: me, user, node, team.
	for _, c := range []struct {
		query string
		field string
		want  bool
	}{
		{"query { viewer { name email } }", "me", false},
		{"query { me { name } }", "me", true},
		{"query { orders { id } }", "orders", true},
		{"query { orders { id } }", "order", false},
		{"query { preorders { id } }", "orders", false},
		{"query { issues(first: 5) { nodes { id } } }", "issues", true},
		{"", "issues", false},
		{"query { issues { id } }", "", false},
	} {
		if got := mentions(c.query, c.field); got != c.want {
			t.Errorf("mentions(%q, %q) = %v, want %v", c.query, c.field, got, c.want)
		}
	}
}
