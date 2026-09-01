package runtime

import "testing"

// A percent-escape is a delimiter, so a marker sitting behind one is a whole
// word.
//
// The whole-word check exists so that selecting on "customer" does not match
// "customerId", and it decides a boundary by looking at the character either
// side. In a form-encoded body every structural character is escaped, so a
// marker immediately after a quote is preceded by the "2" of %22 -- a digit,
// which is a name character, which made the check refuse a match that is
// plainly a whole word.
//
// Expensify found it. Its API takes a JSON document inside a form field, so
// every brace and quote in that document arrives escaped, and three of its
// cases had to route by a JSON body the client does not really send.
func TestAPercentEscapeIsAWordBoundary(t *testing.T) {
	// requestJobDescription={"type":"report", ...}, form-encoded.
	const form = "requestJobDescription=%7B%22type%22%3A%22report%22%2C%22limit%22%3A10%7D"

	for _, field := range []string{"type", "report", "limit"} {
		if !mentions(form, field) {
			t.Errorf("%q is a whole word in the encoded body and was not matched", field)
		}
	}
}

// And the check still refuses a partial word, which is the whole point of it.
//
// The risk in loosening a boundary rule is loosening it into uselessness. A
// marker that is genuinely part of a longer name must still miss, escaped or
// not.
func TestAPercentEscapeDoesNotMakeEveryPartialWordMatch(t *testing.T) {
	const form = "requestJobDescription=%7B%22reportList%22%3A%22x%22%7D"

	if mentions(form, "report") {
		t.Error(`"report" matched inside "reportList"`)
	}

	if !mentions(form, "reportList") {
		t.Error(`"reportList" is the whole word and did not match`)
	}

	// Not an escape at all -- a literal percent followed by two letters that
	// are not hex digits. The character before the match is still a name
	// character, so this is still not a boundary.
	if mentions("xx%zzreport", "report") {
		t.Error(`"report" matched after "%zz", which is not a percent-escape`)
	}
}

// An unescaped body is unaffected, which is the compatibility claim.
//
// Every Recipe using selects_body today sends JSON or GraphQL, where the
// delimiters are literal. None of them should notice this.
func TestAnUnescapedBodyMatchesExactlyAsBefore(t *testing.T) {
	const body = `{"query":"{ viewer { login } }"}`

	if !mentions(body, "viewer") {
		t.Error(`"viewer" is a whole word and did not match`)
	}

	if mentions(body, "view") {
		t.Error(`"view" matched inside "viewer"`)
	}
}
