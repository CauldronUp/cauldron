// Conformance cases: the evidence a Recipe resembles its provider, and the
// inspection each case is held to.

package recipe

import (
	"fmt"
	"reflect"
	"strings"
)

// omitsHeader reports whether any case reaches a route without the header and
// is refused for it.
//
// A case that omits the header and still succeeds proves the opposite, so the
// status has to be a failure for it to count.
func omitsHeader(cases []Case, header string) bool {
	for _, c := range cases {
		if c.Expect.Status < 400 {
			continue
		}

		sent := false

		for h := range c.Request.Headers {
			if strings.EqualFold(h, header) {
				sent = true
				break
			}
		}

		if !sent {
			return true
		}
	}

	return false
}

// sendsParam reports whether any case sends a query parameter by name.
//
// A case may carry its query in the map or written into the path, and both are
// used: Alpaca's filter cases request "/v2/orders?status=closed" directly. A
// check reading only the map would report every one of those as unexercised,
// which is how this was first got wrong.
func sendsParam(cases []Case, param string) bool {
	for _, c := range cases {
		if _, ok := c.Request.Query[param]; ok {
			return true
		}

		if _, query, found := strings.Cut(c.Request.Path, "?"); found {
			for _, pair := range strings.Split(query, "&") {
				if name, _, _ := strings.Cut(pair, "="); name == param {
					return true
				}
			}
		}
	}

	return false
}

// assertsName reports whether any case claims something about a field by name.
//
// The comparison is on the last segment of a dotted path, and it is deliberately
// loose: a Recipe declaring data.pagination.next may be asserted as a dotted
// path or as nested maps, and either pins the name down. What it will not
// accept is silence.
func assertsName(cases []Case, field string) bool {
	segments := strings.Split(field, ".")
	leaf := segments[len(segments)-1]

	for _, c := range cases {
		for path := range c.Expect.Matches {
			if mentions(path, leaf) {
				return true
			}
		}

		if mentionedIn(c.Expect.Body, leaf) {
			return true
		}
	}

	return false
}

// mentions reports whether a dotted path names this field at any depth.
func mentions(path, leaf string) bool {
	for _, segment := range strings.Split(path, ".") {
		// Trim any [0] so items[0].next matches next.
		if index := strings.IndexByte(segment, '['); index >= 0 {
			segment = segment[:index]
		}

		if segment == leaf {
			return true
		}
	}

	return false
}

// mentionedIn walks a body assertion looking for the field, at any depth,
// because a nested claim is as good as a dotted one.
func mentionedIn(node any, leaf string) bool {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			if mentions(key, leaf) || mentionedIn(value, leaf) {
				return true
			}
		}
	case []any:
		for _, value := range typed {
			if mentionedIn(value, leaf) {
				return true
			}
		}
	}

	return false
}

// echoesOnly reports whether a case's only claims repeat its own request.
//
// A create that sends {"name": "x"} and asserts {"name": "x"} is testing that
// the request body survived the round trip, which every fake does by
// construction. It reads as evidence and is not, and the failure is invisible
// until somebody breaks the Recipe on purpose and watches the case pass.
func echoesOnly(c Case) bool {
	if len(c.Expect.Body) == 0 {
		return false
	}

	// Either of these puts something under test that cannot be echoed.
	if len(c.Expect.Matches) > 0 || len(c.Expect.Absent) > 0 {
		return false
	}

	sent := map[string]any{}

	for name, value := range c.Request.JSON {
		sent[name] = value
	}

	for name, value := range c.Request.Form {
		sent[name] = value
	}

	if len(sent) == 0 {
		return false
	}

	for name, want := range c.Expect.Body {
		value, present := sent[name]
		if !present {
			return false
		}

		// Form values arrive as text, so compare rendered forms rather than
		// requiring the YAML types to match.
		if !reflect.DeepEqual(value, want) && fmt.Sprint(value) != fmt.Sprint(want) {
			return false
		}
	}

	return true
}

// Case is one checkable claim about the provider's behaviour.
//
// The point of a conformance case is not that the emulator passes it. Any fake
// passes its own tests. The point is provenance: every case cites where the
// expectation came from, and records whether it was observed against the real
// API or only read in the documentation. A developer deciding whether to trust
// this emulator can then read the evidence rather than the marketing.
type Case struct {
	Name string `yaml:"name"`
	// Source cites the provider documentation or transcript the expectation
	// came from. Required: an uncited claim about someone else's API is a
	// guess wearing a test's clothing.
	Source string `yaml:"source"`
	// Verified is the date this case was last checked against the real API,
	// as YYYY-MM-DD. Empty means the expectation was read, not observed, and
	// the report says so rather than quietly counting it as proof.
	Verified string `yaml:"verified"`
	// Fixture is seeded before the case runs. Empty leaves the sandbox as it is,
	// which lets a group of cases build on each other in order.
	Fixture string `yaml:"fixture"`
	// Arm names an entry in the Recipe's errors table to install before this
	// case's request, and only for it.
	//
	// Without this a Recipe's error table is a list of unverified claims.
	// Every failure a conformance suite could reach was one the runtime
	// produces on its own: a 404 for a missing record, a 401 for a bad
	// credential. The interesting entries, the ones describing a declined
	// card or an expired sync token or a rate limit, were declared and never
	// once exercised, so a field could be renamed, a status changed or a
	// nested detail dropped and nothing anywhere would notice.
	//
	// The fault is armed for exactly one request and cleared afterwards, so a
	// case cannot leak a failure into the next one.
	Arm     string      `yaml:"arm"`
	Request Request     `yaml:"request"`
	Expect  Expectation `yaml:"expect"`
}

// Request is the call a conformance case makes.
type Request struct {
	Method  string            `yaml:"method"`
	Path    string            `yaml:"path"`
	Query   map[string]string `yaml:"query"`
	Headers map[string]string `yaml:"headers"`
	// Form sends application/x-www-form-urlencoded, which is what Stripe's own
	// SDKs send. JSON sends a JSON body. A case may set at most one.
	Form map[string]string `yaml:"form"`
	JSON map[string]any    `yaml:"json"`
}

// Expectation is what the provider is claimed to answer.
//
// Body matching is a subset: a case asserts the fields it is making a claim
// about and ignores the rest, so a Recipe can grow a field without invalidating
// every case ever written about it.
type Expectation struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
	// Matches holds dotted field paths to regular expressions, for values that
	// are correct in shape rather than exact, such as generated identifiers.
	Matches map[string]string `yaml:"matches"`
	// HeaderMatches holds response header names to regular expressions, for
	// headers that carry a generated value. A plain `headers` entry compares
	// substrings, which cannot assert that a header is merely present and
	// well-formed.
	HeaderMatches map[string]string `yaml:"header_matches"`
	// AbsentHeaders lists response headers that must not appear.
	//
	// The absence of a header is a claim as real as its presence, and for
	// paging it is the one that terminates the loop: a provider advertises
	// the next page in Link and sends no Link on the last page, so a client
	// that keeps following one until it is gone stops exactly there. An
	// emulator that sent Link on every page would loop forever, and there
	// was no way to write that down.
	AbsentHeaders []string `yaml:"absent_headers"`
	// Absent lists fields that must not appear. Providers are as specific about
	// what they omit as what they send.
	Absent []string `yaml:"absent"`
	// BodyMatches is a regular expression applied to the raw response body,
	// without parsing it.
	//
	// A provider whose failures are plain text has no assertable body
	// otherwise: matches walks a decoded document, so the only thing Trello's
	// text-error case could pin down was its Content-Type. The prose is the
	// part support threads quote and the part a client ends up regex-matching
	// in anger, so it is worth being able to claim.
	BodyMatches string `yaml:"body_matches"`
	// NoBody asserts the response body is empty.
	//
	// This is a positive claim, not the absence of one. Salesforce answers an
	// update with 204 and nothing at all, so a client calling .json() on it
	// throws rather than seeing that the update worked, and an emulator that
	// helpfully returned an object would hide that. An `absent` list cannot
	// express it: absences are vacuously true against an empty body, so a case
	// built from them would pass whatever the emulator sent.
	NoBody bool `yaml:"no_body"`
	// Webhook asserts what the request emitted, which nothing could assert
	// before.
	//
	// Webhook payloads were the largest unverified surface in the project: 85
	// Recipes emit them, the record went in raw rather than shaped, and no
	// case could look. An application's handler written against the emulator
	// could read a field the provider never sends and be entirely green.
	Webhook *WebhookExpectation `yaml:"webhook"`
}

// WebhookExpectation is what a case claims about the webhook its request
// emitted.
//
// The last delivery is the one examined, because a request emits at most one
// event and asserting on "the one this caused" is the only reading that stays
// true as a Recipe grows.
type WebhookExpectation struct {
	// Event is the type the delivery must carry.
	Event string `yaml:"event"`
	// Body asserts dotted paths in the payload, envelope included, so a
	// Recipe declaring its own envelope can pin that too.
	Body map[string]any `yaml:"body"`
	// Matches asserts regular expressions against payload paths.
	Matches map[string]string `yaml:"matches"`
	// Absent names paths the payload must not carry. This is the half that
	// catches an internal field name leaking into a payload.
	Absent []string `yaml:"absent"`
	// None claims the request emitted nothing at all, which is worth being
	// able to say: an event that fires when it should not is as wrong as one
	// that does not fire.
	None bool `yaml:"none"`
}
