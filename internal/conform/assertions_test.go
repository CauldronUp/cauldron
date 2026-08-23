package conform

import (
	"net/http"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// An assertion that cannot fail is not an assertion, and a case built on one
// passes whatever the emulator does.
//
// Six of the ways a case can assert something were reached by no test: a case
// could use header_matches, absent_headers or body_matches and the only
// evidence they worked was that a Recipe somewhere used them and was green --
// which is exactly what a silently ignored field looks like.
//
// Each kind is given a response it should reject and one it should accept.
func TestEveryKindOfAssertionCanFail(t *testing.T) {
	respond := func(status int, headers map[string]string, body string) (*http.Response, []byte) {
		r := &http.Response{StatusCode: status, Header: http.Header{}}
		for name, value := range headers {
			r.Header.Set(name, value)
		}

		return r, []byte(body)
	}

	for _, c := range []struct {
		kind   string
		expect recipe.Expectation
		// The response that should be refused, and the one that should pass.
		badStatus, goodStatus   int
		badHeaders, goodHeaders map[string]string
		badBody, goodBody       string
		wants                   string
	}{
		{
			kind:      "status",
			expect:    recipe.Expectation{Status: 201},
			badStatus: 200, goodStatus: 201,
			badBody: "{}", goodBody: "{}",
			wants: "status",
		},
		{
			kind:      "headers",
			expect:    recipe.Expectation{Status: 200, Headers: map[string]string{"Content-Type": "application/json"}},
			badStatus: 200, goodStatus: 200,
			badHeaders:  map[string]string{"Content-Type": "text/plain"},
			goodHeaders: map[string]string{"Content-Type": "application/json"},
			badBody:     "{}", goodBody: "{}",
			wants: "header Content-Type",
		},
		{
			kind:      "header_matches",
			expect:    recipe.Expectation{Status: 200, HeaderMatches: map[string]string{"Link": `rel="next"`}},
			badStatus: 200, goodStatus: 200,
			badHeaders:  map[string]string{"Link": `<https://example.invalid/2>; rel="prev"`},
			goodHeaders: map[string]string{"Link": `<https://example.invalid/2>; rel="next"`},
			badBody:     "{}", goodBody: "{}",
			wants: "header Link",
		},
		{
			kind:      "absent_headers",
			expect:    recipe.Expectation{Status: 200, AbsentHeaders: []string{"Link"}},
			badStatus: 200, goodStatus: 200,
			badHeaders: map[string]string{"Link": "<https://example.invalid/2>; rel=\"next\""},
			badBody:    "{}", goodBody: "{}",
			wants: "Link",
		},
		{
			kind:      "body",
			expect:    recipe.Expectation{Status: 200, Body: map[string]any{"state": "closed"}},
			badStatus: 200, goodStatus: 200,
			badBody: `{"state":"open"}`, goodBody: `{"state":"closed"}`,
			wants: "state",
		},
		{
			kind:      "matches",
			expect:    recipe.Expectation{Status: 200, Matches: map[string]string{"id": "^cus_"}},
			badStatus: 200, goodStatus: 200,
			badBody: `{"id":"acct_1"}`, goodBody: `{"id":"cus_1"}`,
			wants: "id",
		},
		{
			kind:      "absent",
			expect:    recipe.Expectation{Status: 200, Absent: []string{"error"}},
			badStatus: 200, goodStatus: 200,
			badBody: `{"error":"nope"}`, goodBody: `{"ok":true}`,
			wants: "error",
		},
		{
			kind: "body_matches",
			// The raw body, quotes and all, which is why npm's own case for
			// this writes the opening quote into the pattern.
			expect:    recipe.Expectation{Status: 404, BodyMatches: `^"version not found`},
			badStatus: 404, goodStatus: 404,
			badBody: `"no such thing"`, goodBody: `"version not found: 99.99.99"`,
			wants: "body",
		},
		{
			kind:      "no_body",
			expect:    recipe.Expectation{Status: 204, NoBody: true},
			badStatus: 204, goodStatus: 204,
			badBody: `{"id":"1"}`, goodBody: "",
			wants: "body",
		},
	} {
		t.Run(c.kind, func(t *testing.T) {
			response, body := respond(c.badStatus, c.badHeaders, c.badBody)

			failures := check(c.expect, response, body)
			if len(failures) == 0 {
				t.Fatalf("%s accepted a response it should have refused", c.kind)
			}

			if !strings.Contains(strings.Join(failures, "; "), c.wants) {
				t.Errorf("%s refused it and said %q, which does not mention %q", c.kind, failures, c.wants)
			}

			response, body = respond(c.goodStatus, c.goodHeaders, c.goodBody)

			if failures := check(c.expect, response, body); len(failures) > 0 {
				t.Errorf("%s refused a response it should have accepted: %q", c.kind, failures)
			}
		})
	}
}
