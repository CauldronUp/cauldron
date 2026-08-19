package conform

import (
	"net/http"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// encoding/json decodes every number into a float64 unless told otherwise, and
// a float64 cannot hold an int64. So this comparator could not see a number
// larger than nine quadrillion, which is a problem for a tool whose entire job
// is noticing that an emulator sent the wrong one.
//
// Typesense made it visible. Its relevance score is an int64 up near
// 578730123365711993; the emulator put those exact digits on the wire, and the
// comparator read 578730123365712000. The hit ranked immediately below it,
// 578730123365711994, read as the same number, so no case could have told the
// two apart.

func checkOne(t *testing.T, body string, path string, want any) []string {
	t.Helper()

	// Through check, so the production decode is the one under test. A helper
	// that decoded the body itself would leave UseNumber unexercised, and two
	// mutations proved exactly that before this was written.
	expect := recipe.Expectation{
		Status: http.StatusOK,
		Body:   map[string]any{path: want},
	}

	return check(expect, &http.Response{StatusCode: http.StatusOK, Header: http.Header{}}, []byte(body))
}

func TestAnInt64SurvivesTheComparator(t *testing.T) {
	body := `{"text_match": 578730123365711993}`

	if failures := checkOne(t, body, "text_match", 578730123365711993); len(failures) != 0 {
		t.Errorf("the true value was rejected: %v", failures)
	}
}

func TestTwoAdjacentInt64sAreToldApart(t *testing.T) {
	// The pair that rounds to one float64. If these compare equal, the
	// comparator is reading a number nobody sent.
	body := `{"text_match": 578730123365711994}`

	failures := checkOne(t, body, "text_match", 578730123365711993)
	if len(failures) == 0 {
		t.Fatal("two scores one apart compared equal")
	}

	if !strings.Contains(failures[0], "578730123365711993") {
		t.Errorf("the failure should quote the digits, got %v", failures)
	}
}

func TestOrdinaryNumbersStillCompare(t *testing.T) {
	for _, c := range []struct {
		body string
		path string
		want any
	}{
		{`{"n": 3}`, "n", 3},
		{`{"n": 0}`, "n", 0},
		{`{"n": -7}`, "n", -7},
		{`{"n": 0.8552}`, "n", 0.8552},
		{`{"n": 60.288}`, "n", 60.288},
	} {
		if failures := checkOne(t, c.body, c.path, c.want); len(failures) != 0 {
			t.Errorf("%s: %v", c.body, failures)
		}
	}
}

func TestANumberIsStillNotAString(t *testing.T) {
	// Telling those apart is the whole point of several Recipes: Vonage says
	// a send succeeded with the string "0". UseNumber makes a number a string
	// type internally, and it must not make it a string here.
	if failures := checkOne(t, `{"n": 0}`, "n", "0"); len(failures) == 0 {
		t.Error("the number 0 matched the string \"0\"")
	}

	if failures := checkOne(t, `{"n": "0"}`, "n", 0); len(failures) == 0 {
		t.Error("the string \"0\" matched the number 0")
	}
}
