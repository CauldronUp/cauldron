package conform

import (
	"net/http"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Dropbox names a field ".tag", where the leading dot is part of the name
// rather than a separator. Without an escape there is no way to assert on it,
// and the emulator was sending it correctly all along while every case that
// mentioned it failed.
func TestSplitFieldPathHonoursAnEscapedDot(t *testing.T) {
	cases := []struct {
		path string
		want []string
	}{
		{"data.id", []string{"data", "id"}},
		{`\.tag`, []string{".tag"}},
		{`entries[0].\.tag`, []string{"entries[0]", ".tag"}},
		{`a.\.b.c`, []string{"a", ".b", "c"}},
		{"plain", []string{"plain"}},
	}

	for _, c := range cases {
		got := splitFieldPath(c.path)

		if len(got) != len(c.want) {
			t.Errorf("splitFieldPath(%q) = %v, want %v", c.path, got, c.want)
			continue
		}

		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("splitFieldPath(%q) = %v, want %v", c.path, got, c.want)
				break
			}
		}
	}
}

func TestLookupFindsAFieldNamedWithALeadingDot(t *testing.T) {
	document := map[string]any{
		".tag": "file",
		"entries": []any{
			map[string]any{".tag": "folder", "name": "Documents"},
		},
	}

	value, ok := lookup(document, `\.tag`)
	if !ok || value != "file" {
		t.Errorf("lookup of the escaped tag = %v, %v", value, ok)
	}

	nested, ok := lookup(document, `entries[0].\.tag`)
	if !ok || nested != "folder" {
		t.Errorf("lookup of the nested escaped tag = %v, %v", nested, ok)
	}

	// An unescaped dot still separates, so nothing about ordinary paths
	// changes.
	if _, ok := lookup(document, ".tag"); ok {
		t.Error("an unescaped leading dot should not resolve")
	}
}

// A string "0" is not the number 0, and telling them apart is the whole point
// of several Recipes: Vonage reports a successful send with the string "0" and
// Docusign counts with strings. Comparing only the rendered form let those
// cases pass whichever the emulator sent.
func TestScalarsCompareByKindAsWellAsValue(t *testing.T) {
	document := map[string]any{
		"stringZero": "0",
		"numberZero": float64(0),
		"count":      "4",
	}

	if failures := compare("stringZero", "0", document); len(failures) != 0 {
		t.Errorf("a string against a string should match: %v", failures)
	}

	// The bug this test exists for: without a kind check these pass.
	if failures := compare("stringZero", float64(0), document); len(failures) == 0 {
		t.Error("the number 0 should not match the string \"0\"")
	}

	if failures := compare("numberZero", "0", document); len(failures) == 0 {
		t.Error("the string \"0\" should not match the number 0")
	}

	if failures := compare("count", 4, document); len(failures) == 0 {
		t.Error("the number 4 should not match the string \"4\"")
	}

	// Integer and float are the same kind, because YAML and JSON disagree
	// about which one a literal is and no provider distinguishes them.
	if failures := compare("numberZero", 0, document); len(failures) != 0 {
		t.Errorf("an int should match a float of the same value: %v", failures)
	}
}

// Salesforce answers an update with 204 and nothing at all, so a client calling
// .json() on the response throws rather than seeing that the update worked. A
// case claiming that has nothing to put in body, matches or absent — there is
// no body to put a claim against — and check returned early whenever all three
// were empty, so the claim was never checked. The case passed whatever the
// emulator sent, including a full record.
func TestAnEmptyBodyClaimIsCheckedInBothDirections(t *testing.T) {
	response := &http.Response{StatusCode: http.StatusNoContent, Header: http.Header{}}
	expect := recipe.Expectation{Status: http.StatusNoContent, NoBody: true}

	if failures := check(expect, response, nil); len(failures) != 0 {
		t.Errorf("an empty body should satisfy the claim: %v", failures)
	}

	// The bug this test exists for: without the check moved above the early
	// return, this passes.
	failures := check(expect, response, []byte(`{"Id":"001cauldron000001A"}`))
	if len(failures) == 0 {
		t.Error("a body should not satisfy a claim that there is no body")
	}
}

// Three states, not two. A field can be absent, present and null, or present
// and empty, and providers mean different things by each.
//
// AssemblyAI sends text: null on a queued transcript and utterances: null on
// a completed one that did not ask for speaker labels. Neither is an absence
// and neither is an empty value: the key is there and unusable, so a check for
// the key passes and the value cannot be used. A checker that collapsed any of
// these into the others would let a Recipe claim the wrong one.
func TestNullIsNeitherAbsentNorEmpty(t *testing.T) {
	document := map[string]any{
		"text":       nil,
		"utterances": nil,
		"empty":      "",
		"list":       []any{},
	}

	// Present and null is found, so an absence claim about it must fail.
	if _, ok := lookup(document, "text"); !ok {
		t.Error("a present null should be found, or absent would wrongly match it")
	}

	if _, ok := lookup(document, "never"); ok {
		t.Error("an absent key should not be found")
	}

	// Null is not the empty string.
	if failures := compare("text", "", document); len(failures) == 0 {
		t.Error(`null should not match ""`)
	}

	// Nor an empty list.
	if failures := compare("utterances", []any{}, document); len(failures) == 0 {
		t.Error("null should not match an empty array")
	}

	// And the empty forms are not null either, in the other direction.
	if failures := compare("empty", nil, document); len(failures) == 0 {
		t.Error(`"" should not match null`)
	}

	if failures := compare("list", nil, document); len(failures) == 0 {
		t.Error("an empty array should not match null")
	}

	// An absence claim is not satisfied by a present null.
	if failures := check(recipe.Expectation{Absent: []string{"text"}}, nil, []byte(`{"text":null}`)); len(failures) == 0 {
		t.Error("absent should not be satisfied by a present null")
	}
}
