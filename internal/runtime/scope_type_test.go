package runtime

import "testing"

// A path is text. A field is not.
//
// Help Scout's conversationId, Gorgias's ticket_id, Shortcut's epic_id,
// Documenso's documentId and Hetzner's resources[0].id are integers on the
// wire, and a record created through a scoped route carried the string out of
// the URL instead. So the same field was a number when a fixture seeded it and
// a string when a create produced it, in the same collection and the same
// response shape, and only one of those matched the provider.
//
// Nothing caught it because a conformance case that creates and then asserts
// is asserting whatever the create produced.

func TestAScopedIdentifierTakesItsDeclaredType(t *testing.T) {
	if got := scoped("integer", "42"); got != int64(42) {
		t.Errorf("integer scope stayed %#v", got)
	}

	if got := scoped("number", "1.5"); got != 1.5 {
		t.Errorf("number scope stayed %#v", got)
	}
}

func TestAScopedStringIsLeftAlone(t *testing.T) {
	// The common case by far: owner, repo, workspace, account.
	for _, in := range []string{"octocat", "hello-world", "42"} {
		if got := scoped("string", in); got != in {
			t.Errorf("scoped(string, %q) = %#v", in, got)
		}
	}

	// An undeclared field is not a claim, so nothing is converted.
	if got := scoped("", "42"); got != "42" {
		t.Errorf("scoped with no declared type = %#v", got)
	}
}

func TestAnIdentifierThatOnlyLooksNumericSurvives(t *testing.T) {
	// The JSON grammar decides, for the same reason it decides in a form
	// body. A leading zero or a plus is not something a JSON client could
	// have sent, so a key that begins with one is a key rather than a number
	// however the field is declared.
	for _, in := range []string{"007", "+42", "42abc", "", "0x1f"} {
		if got := scoped("integer", in); got != in {
			t.Errorf("scoped(integer, %q) = %#v, want it left alone", in, got)
		}
	}
}
