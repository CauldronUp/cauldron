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

// An error with no body is shown by asserting there is no body.
//
// Bitwarden answers 405 with a status line and nothing else, and declares that
// as empty: true with no code and no message. There is no wording to quote back,
// so requiring one asked for evidence that cannot exist -- and forty-odd entries
// across the corpus sat in the count with no way out of it.
//
// The claim a case can make is the one the format already has for this: status,
// and no_body. That is exactly what a client gets, and it is not nothing --
// calling .json() on it throws.
func TestAnEmptyErrorIsShownByNoBody(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"method_not_allowed": {Status: 405, Empty: true},
		},
		Conformance: []Case{{
			Request: Request{Method: "DELETE", Path: "/v1/things"},
			Expect:  Expectation{Status: 405},
		}},
	}

	if n := r.UnshownError(); n != 1 {
		t.Fatalf("a case asserting only the status: counted %d unshown, want 1", n)
	}

	r.Conformance[0].Expect.NoBody = true

	if n := r.UnshownError(); n != 0 {
		t.Errorf("a case asserting the empty body: counted %d unshown, want 0", n)
	}
}

// An empty error with a header is shown by the header, as before.
func TestAnEmptyErrorWithAHeaderIsShownByIt(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{
			"method_not_allowed": {Status: 405, Empty: true, Headers: map[string]string{"Allow": "GET, POST"}},
		},
		Conformance: []Case{{
			Request: Request{Method: "DELETE", Path: "/v1/things"},
			Expect:  Expectation{Status: 405, Headers: map[string]string{"Allow": "GET, POST"}},
		}},
	}

	if n := r.UnshownError(); n != 0 {
		t.Errorf("an empty error asserted by its Allow header: counted %d unshown, want 0", n)
	}
}

// An entry that declares no code is often carried under its own name.
//
// Bringg declares method_not_allowed as a bare 404 and its envelope answers
// {"error": "method_not_allowed"}, so the name is the code even though the
// Recipe never repeats it.
func TestAnErrorWithNoCodeIsShownByItsName(t *testing.T) {
	r := &Recipe{
		Errors: map[string]Error{"method_not_allowed": {Status: 404}},
		Conformance: []Case{{
			Request: Request{Method: "DELETE", Path: "/v1/things"},
			Expect: Expectation{Status: 404, Body: map[string]any{
				"error": "method_not_allowed",
			}},
		}},
	}

	if n := r.UnshownError(); n != 0 {
		t.Errorf("an error carried under its own name: counted %d unshown, want 0", n)
	}
}
