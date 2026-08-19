package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some endpoints carry no identifier at all and answer about whoever asked.
// GitHub's /user, Stripe's /v1/account, Slack's auth.test and Backblaze's
// b2_authorize_account are all this shape.
//
// The format could previously say the id was in the path, the query or the
// body, and all three name a place to read from. The whole point of these
// routes is that there is nothing to read, so none of them could be described
// at all: a Recipe either invented a query parameter the provider ignores, or
// modelled a single object as a collection of one and got the wire shape
// wrong.

func TestARouteCanBeIdentifiedByItsCredentials(t *testing.T) {
	r, err := recipe.Open("backblaze")
	if err != nil {
		t.Fatalf("open backblaze: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("one-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/b2api/v3/b2_authorize_account", nil)
	req.Header.Set("Authorization", "cauldron_b2_authorization_token")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("not JSON: %v %s", err, rec.Body.String())
	}

	if body["accountId"] != "b2a41c8f0e19" {
		t.Errorf("accountId is %v, want the seeded account", body["accountId"])
	}

	// A single object, not a collection of one. The distinction is the whole
	// reason this could not be modelled as a listing.
	if _, wrapped := body["authorizations"]; wrapped {
		t.Errorf("answered with a collection: %s", rec.Body.String())
	}
}

func TestAnIdentityRouteRefusesAnAmbiguousFixture(t *testing.T) {
	// Two records and nothing to choose between them with. The runtime would
	// pick the first, which is the seeding order, which is not a decision
	// anybody made.
	source := `
recipe: whoami
capability: auth
version: 0.1.0
upstream:
  api: "1"
resources:
  account:
    collection: accounts
    id:
      style: opaque
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/account
    resource: account
    operation: get
    id_from: auth
fixtures:
  two:
    account:
      - id: a1
        name: one
      - id: a2
        name: two
conformance:
  - name: an account has a name
    source: https://docs.whoami.test
    fixture: two
    request:
      method: GET
      path: /v1/account
    expect:
      status: 200
      body:
        name: one
`

	_, err := recipe.Parse([]byte(source))
	if err == nil {
		t.Fatal("two identities behind one credential-shaped route was accepted")
	}
}

func TestAnIdentityRouteRefusesAPathIdentifier(t *testing.T) {
	// If the path carries an identifier then the request does carry one, and
	// saying otherwise is describing a different endpoint.
	source := `
recipe: whoami
capability: auth
version: 0.1.0
upstream:
  api: "1"
resources:
  account:
    collection: accounts
    id:
      style: opaque
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/accounts/{id}
    resource: account
    operation: get
    id_from: auth
fixtures:
  one:
    account:
      - id: a1
        name: one
conformance:
  - name: an account has a name
    source: https://docs.whoami.test
    fixture: one
    request:
      method: GET
      path: /v1/accounts/a1
    expect:
      status: 200
      body:
        name: one
`

	_, err := recipe.Parse([]byte(source))
	if err == nil {
		t.Fatal("id_from auth on a path that carries an id was accepted")
	}
}
