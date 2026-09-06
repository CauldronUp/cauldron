package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some APIs put the credential in the request body.
//
// Canny is one: every read is a POST and the key is a form field called apiKey
// alongside the filters. That is worth reproducing exactly rather than
// approximating with a header, for the same reason the query scheme exists --
// a header-based fake hides where the secret ends up. A body credential cannot
// be set once as a default, so every call site carries it, and it lands in
// anything that logs request bodies, which is the logging people switch on
// during an incident without thinking about secrets.
//
// Modelling it as query would be worse than not modelling it: the secret would
// go in the URL, which is a different exposure with different consequences, and
// the Recipe would be teaching the wrong lesson confidently.
func TestACredentialCanTravelInTheBody(t *testing.T) {
	auth := recipe.Auth{
		Scheme: "body",
		Param:  "apiKey",
		Keys:   []string{"cauldron-key"},
	}

	s := &Sandbox{}

	for _, tt := range []struct {
		name string
		body string
		want recipe.Verdict
	}{
		{"the key on its own", `{"apiKey": "cauldron-key"}`, recipe.Accepted},
		{"the key beside filters", `{"boardID": "b1", "apiKey": "cauldron-key"}`, recipe.Accepted},
		{"a key nobody issued", `{"apiKey": "not-the-key"}`, recipe.Rejected},
		{"no key at all", `{"boardID": "b1"}`, recipe.Absent},
		{"an empty body", `{}`, recipe.Absent},
	} {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/posts/list", strings.NewReader(tt.body))
		r.Header.Set("Content-Type", "application/json")

		got, _ := s.credential(r, auth)
		if got != tt.want {
			t.Errorf("%s: verdict %v, want %v", tt.name, got, tt.want)
		}
	}
}

// And reading it must leave the body for whoever handles the request.
func TestReadingABodyCredentialDoesNotConsumeTheBody(t *testing.T) {
	auth := recipe.Auth{Scheme: "body", Param: "apiKey", Keys: []string{"cauldron-key"}}
	s := &Sandbox{}

	r := httptest.NewRequest(http.MethodPost, "/api/v1/posts/list",
		strings.NewReader(`{"apiKey": "cauldron-key", "boardID": "b1"}`))
	r.Header.Set("Content-Type", "application/json")

	if got, _ := s.credential(r, auth); got != recipe.Accepted {
		t.Fatalf("verdict %v, want Accepted", got)
	}

	if board := jsonBody(r)["boardID"]; board != "b1" {
		t.Errorf("the body was consumed by the credential check: boardID = %v", board)
	}
}

// A body credential can be nested, and Authorize.Net's is.
//
// Every operation there POSTs to one URL and carries
// merchantAuthentication: {name, transactionKey} inside the request body, so
// the credential is two levels down. Paging already reads dotted names out of
// a body for the same reason -- Plaid keeps count and offset under options --
// and a credential is no different.
func TestABodyCredentialCanBeNested(t *testing.T) {
	auth := recipe.Auth{
		Scheme: "body",
		Param:  "merchantAuthentication.transactionKey",
		Keys:   []string{"cauldron-key"},
	}

	s := &Sandbox{}

	for _, tt := range []struct {
		name string
		body string
		want recipe.Verdict
	}{
		{"nested and right", `{"merchantAuthentication": {"name": "login", "transactionKey": "cauldron-key"}}`, recipe.Accepted},
		{"nested and wrong", `{"merchantAuthentication": {"name": "login", "transactionKey": "nope"}}`, recipe.Rejected},
		{"the object without the key", `{"merchantAuthentication": {"name": "login"}}`, recipe.Absent},
		{"no object at all", `{"createTransactionRequest": {}}`, recipe.Absent},
	} {
		r := httptest.NewRequest(http.MethodPost, "/xml/v1/request.api", strings.NewReader(tt.body))
		r.Header.Set("Content-Type", "application/json")

		if got, _ := s.credential(r, auth); got != tt.want {
			t.Errorf("%s: verdict %v, want %v", tt.name, got, tt.want)
		}
	}
}
