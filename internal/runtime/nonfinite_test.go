package runtime

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Go's ParseFloat accepts NaN, Inf, Infinity and +Inf, case-insensitively, and
// coerce ran it over every form value. encoding/json cannot encode any of
// them, and writeJSON wrote the status line before asking the encoder, so the
// failure arrived as a 200 with Content-Type: application/json and zero bytes.
//
// The record stayed in the collection, so every later read of that collection
// answered the same way until something reset the sandbox. One test posting
// the word "Infinity" made every later test that listed that resource fail on
// res.json(), a long way from the cause, with nothing in any log.
func TestAFormValueCannotPoisonACollection(t *testing.T) {
	for _, value := range []string{"NaN", "Inf", "-Inf", "Infinity", "+Inf", "nan", "infinity"} {
		s := twilioSandbox(t)

		post := httptest.NewRequest(http.MethodPost, "/2010-04-01/Accounts/ACcauldron/Messages.json",
			strings.NewReader("To=%2B15551234567&From=%2B15559876543&Body="+value))
		post.SetBasicAuth("ACcauldron00000000000000000000000", "cauldron-twilio-auth-token")
		post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, post)

		if rec.Body.Len() == 0 {
			t.Errorf("%q: the create answered %d with an empty body", value, rec.Code)
		}

		// Whatever the create did, a later read of the collection must still
		// produce a body. That is the part that used to stay broken.
		list := httptest.NewRequest(http.MethodGet, "/2010-04-01/Accounts/ACcauldron/Messages.json", nil)
		list.SetBasicAuth("ACcauldron00000000000000000000000", "cauldron-twilio-auth-token")

		listed := httptest.NewRecorder()
		s.ServeHTTP(listed, list)

		if listed.Body.Len() == 0 {
			t.Errorf("%q: listing the collection afterwards answered %d with an empty body", value, listed.Code)
		}
	}
}

// coerce turns form strings into the types a JSON client would have sent. A
// non-finite float is not one of those types, because no JSON client can send
// one, so the string is what a form carrying that word actually means.
func TestCoerceKeepsANonFiniteFloatAsText(t *testing.T) {
	for _, value := range []string{"NaN", "Inf", "-Inf", "Infinity", "+Inf", "nan"} {
		got := coerce(value)

		if text, ok := got.(string); !ok || text != value {
			t.Errorf("coerce(%q) = %#v, want the string back", value, got)
		}
	}

	// The ordinary cases must be untouched.
	for value, want := range map[string]any{
		"42":    int64(42),
		"-7":    int64(-7),
		"3.5":   3.5,
		"true":  true,
		"false": false,
		"hello": "hello",
	} {
		if got := coerce(value); got != want {
			t.Errorf("coerce(%q) = %#v, want %#v", value, got, want)
		}
	}
}

// The second half of the same failure: even if something non-encodable reaches
// the writer, the caller must not be told 200 and handed nothing. A truncated
// success is the worst available answer, because it is the one a client parses
// before it fails.
func TestWriteJsonDoesNotAnswerSuccessWithNothing(t *testing.T) {
	rec := httptest.NewRecorder()

	(&Sandbox{}).writeJSON(rec, http.StatusOK, map[string]any{"amount": math.NaN()})

	if rec.Code == http.StatusOK {
		t.Errorf("an unencodable body answered 200")
	}

	if rec.Body.Len() == 0 {
		t.Error("an unencodable body produced no body at all, which is what the caller then tries to parse")
	}
}

func twilioSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("twilio")
	if err != nil {
		t.Fatalf("open twilio: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	return s
}
