package recipe

import (
	"testing"
)

// An error's wording is a claim, and nothing was checking most of them.
//
// A Recipe declares what a provider says when it refuses: a status, a code, a
// message, sometimes a header that tells two refusals apart. That wording is
// the part a client switches on and the part a support thread quotes, and it
// had no rule behind it -- rename Stripe's "resource_missing" to anything and
// every case stays green, exactly as with a paging parameter or a field name.
//
// An entry counts as shown when some case expects its status and asserts one of
// the things that tells it apart: the code, the message, or a header value.
func TestAnErrorNoCaseAssertsIsCounted(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"resource_missing": {Status: 404, Code: "not_found", Message: "Not Found"},
			"rate_limit":       {Status: 429, Code: "too_many", Message: "Slow down"},
		},
		Conformance: []Case{{
			Name:    "an unknown id is refused",
			Request: Request{Method: "GET", Path: "/v1/things/nope"},
			Expect:  Expectation{Status: 404, Body: map[string]any{"error.code": "not_found"}},
		}},
	}

	if n := r.UnshownError(); n != 1 {
		t.Fatalf("one of two errors asserted: counted %d unshown, want 1", n)
	}

	r.Conformance = append(r.Conformance, Case{
		Name:    "the rate limit says so",
		Arm:     "rate_limit",
		Request: Request{Method: "GET", Path: "/v1/things"},
		Expect:  Expectation{Status: 429, Body: map[string]any{"error.message": "Slow down"}},
	})

	if n := r.UnshownError(); n != 0 {
		t.Errorf("both errors asserted: counted %d unshown, want 0", n)
	}
}

// The status alone is not the wording. Two refusals commonly share one.
func TestAStatusAloneDoesNotShowAnError(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"authentication_error": {Status: 401, Code: "unauthorized", Message: "Unauthorized"},
		},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 401},
		}},
	}

	if n := r.UnshownError(); n != 1 {
		t.Errorf("a case asserting only the status: counted %d unshown, want 1", n)
	}
}

// A message with a placeholder in it is shown by the part that does not vary.
func TestAMessageIsMatchedByItsFixedPart(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"resource_missing": {Status: 404, Message: "Not Found: {detail}"},
		},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things/nope"},
			Expect:  Expectation{Status: 404, Body: map[string]any{"message": "Not Found: thing"}},
		}},
	}

	if n := r.UnshownError(); n != 0 {
		t.Errorf("a message asserted with its detail filled in: counted %d unshown, want 0", n)
	}
}

// Plain-text failures are asserted with a regular expression over the body.
func TestABodyMatchesShowsAnError(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"missing_signature_params": {Status: 400, Message: "expected string"},
		},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 400, BodyMatches: "^expected string$"},
		}},
	}

	if n := r.UnshownError(); n != 0 {
		t.Errorf("a plain-text failure asserted by regex: counted %d unshown, want 0", n)
	}
}
