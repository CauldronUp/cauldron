package runtime

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A credential can fail in three ways and a Recipe can answer each differently.
//
// Until this existed the check was a bool, so every failure produced one
// message. That was fine for most providers and wrong for a dozen: Pipedream
// answers an absent credential with 404 "record not found" -- the record being
// the caller's own account -- and a junk one with a plain 401. Make has three
// distinct sentences. FireHydrant has two and had to reach the second by arming
// a fault, which is a workaround with a header comment apologising for it.
//
// The verdicts are decided by the carrier, not the value. A header that is not
// on the request is absent. A header that is there and cannot be a credential,
// because the declared prefix rules it out, is malformed. Anything of the right
// shape that is not a key is rejected.
func TestACredentialFailsInThreeDistinguishableWays(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	if r.Auth.Scheme != "bearer" || r.Auth.Prefix == "" {
		t.Fatalf("firehydrant no longer uses a prefixed bearer; pick another Recipe")
	}

	cases := []struct {
		name   string
		header string
		value  string
		want   recipe.Verdict
	}{
		{"no header at all", "", "", recipe.Absent},
		{"the wrong prefix", "Authorization", "Token fhb-nope", recipe.Malformed},
		{"no prefix at all", "Authorization", "fhb-nope", recipe.Malformed},
		{"prefixed and unknown", "Authorization", r.Auth.Prefix + "fhb-nope", recipe.Rejected},
		{"a key the Recipe holds", "Authorization", r.Auth.Prefix + r.Auth.Keys[0], recipe.Accepted},
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/v1/incidents", nil)
		if c.header != "" {
			req.Header.Set(c.header, c.value)
		}

		if got := s.credential(req); got != c.want {
			t.Errorf("%s: verdict %d, want %d", c.name, got, c.want)
		}
	}
}

// A Recipe that names an error per verdict serves that error.
//
// Constructed rather than read from a shipped Recipe, because the point is the
// mechanism: the same request path, three credentials, three different bodies
// and three different statuses. Pipedream's absent case really is a 404, so the
// status has to come from the named error rather than being pinned at 401.
func TestEachCredentialVerdictCanServeItsOwnError(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	split := *r
	split.Errors = map[string]recipe.Error{}

	for name, e := range r.Errors {
		split.Errors[name] = e
	}

	split.Errors["nothing_sent"] = recipe.Error{Status: 404, Message: "record not found"}
	split.Errors["not_a_token"] = recipe.Error{Status: 400, Message: "invalid number of segments"}
	split.Auth.AbsentError = "nothing_sent"
	split.Auth.MalformedError = "not_a_token"
	split.Errors["stale_key"] = recipe.Error{Status: 403, Message: "invalid or expired"}
	split.Auth.RejectedError = "stale_key"

	if err := split.Validate(); err != nil {
		t.Fatalf("the constructed Recipe does not validate: %v", err)
	}

	cases := []struct {
		name    string
		value   string
		status  int
		message string
	}{
		{"absent", "", 404, "record not found"},
		{"malformed", "Token fhb-nope", 400, "invalid number of segments"},
		{"rejected", split.Auth.Prefix + "fhb-nope", 403, "invalid or expired"},
	}

	for _, c := range cases {
		s, err := New(&split, Options{Seed: 1})
		if err != nil {
			t.Fatalf("%s: new sandbox: %v", c.name, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/incidents", nil)
		if c.value != "" {
			req.Header.Set("Authorization", c.value)
		}

		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		if w.Code != c.status {
			t.Errorf("%s: status %d, want %d", c.name, w.Code, c.status)
		}

		body, _ := io.ReadAll(w.Body)
		if !strings.Contains(string(body), c.message) {
			t.Errorf("%s: body %s does not carry %q", c.name, body, c.message)
		}

		var decoded any
		if err := json.Unmarshal(body, &decoded); err != nil {
			t.Errorf("%s: body is not JSON: %v", c.name, err)
		}
	}
}

// A Recipe declaring neither field answers all three verdicts identically.
//
// This is the compatibility claim, and it is the one worth a test: 394 Recipes
// were written against a check that could not tell these apart, and none of
// them should have changed. Byte for byte, not merely the same status.
func TestARecipeNamingNoVerdictErrorsAnswersThemAllTheSameWay(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	// Cleared explicitly rather than asserted, so that wiring FireHydrant to
	// the new fields cannot turn this test into a skip. What is being claimed
	// is about Recipes that name neither, and the way to test that is to hold
	// one, not to find one.
	silent := *r
	silent.Auth.AbsentError = ""
	silent.Auth.MalformedError = ""
	silent.Auth.RejectedError = ""

	var first []byte

	for i, value := range []string{"", "Token fhb-nope", silent.Auth.Prefix + "fhb-nope"} {
		s, err := New(&silent, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/incidents", nil)
		if value != "" {
			req.Header.Set("Authorization", value)
		}

		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		body, _ := io.ReadAll(w.Body)

		if i == 0 {
			first = body

			continue
		}

		if string(body) != string(first) {
			t.Errorf("credential %q answered %s, want %s", value, body, first)
		}
	}
}
