package runtime

import "testing"

// A form body is text. Turning it into the types a JSON client would have
// sent is right, and strconv is not the grammar that decides what those are:
// it accepts three things JSON does not, and each one silently rewrites a
// value a caller posted.
//
// The one that matters most is the leading plus. Every provider that takes a
// destination phone number in a form body takes it in E.164, and E.164 is the
// plus. Parsing it as an integer drops the character that made it a phone
// number and hands back something a client has to guess the country of.

func TestAPhoneNumberSurvivesAFormBody(t *testing.T) {
	if got := coerce("+15017122661"); got != "+15017122661" {
		t.Errorf("E.164 must stay a string, got %#v", got)
	}
}

func TestALeadingZeroIsNotDropped(t *testing.T) {
	// Account codes, sort codes and Adyen's "000" all begin with one, and
	// JSON has no syntax for a number that does.
	if got := coerce("007"); got != "007" {
		t.Errorf("leading zero must stay a string, got %#v", got)
	}
}

func TestRealNumbersAreStillNumbers(t *testing.T) {
	for _, c := range []struct {
		in   string
		want any
	}{
		{"0", int64(0)},
		{"42", int64(42)},
		{"-3", int64(-3)},
		{"12.5", 12.5},
		{"-0.25", -0.25},
		{"1e3", 1000.0},
	} {
		if got := coerce(c.in); got != c.want {
			t.Errorf("coerce(%q) = %#v, want %#v", c.in, got, c.want)
		}
	}
}

func TestTheThingsJSONCannotHoldStayText(t *testing.T) {
	// NaN and the infinities are the dangerous ones: encoding/json refuses to
	// write them, so storing one made every later read of that collection
	// answer with an empty body until the sandbox was reset.
	for _, in := range []string{"NaN", "Inf", "+Inf", "-Inf", "Infinity", "1 2", "12abc", "", "0x1f"} {
		if got := coerce(in); got != in {
			t.Errorf("coerce(%q) = %#v, want it left alone", in, got)
		}
	}
}

func TestBooleansAreStillBooleans(t *testing.T) {
	if got := coerce("true"); got != true {
		t.Errorf(`coerce("true") = %#v`, got)
	}

	if got := coerce("false"); got != false {
		t.Errorf(`coerce("false") = %#v`, got)
	}
}
