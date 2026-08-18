package recipe

import (
	"strings"
	"testing"
)

// A hundred and twenty-seven names is too many to print. It was reasonable at
// four, and the message that once listed everything now buries the answer:
// somebody who typed "stripo" had to read a wall of text to find "stripe".

func TestSuggestFindsTheObviousTypo(t *testing.T) {
	for typed, want := range map[string]string{
		"stripo":        "stripe",
		"moderntreasry": "moderntreasury",
		"secretmanager": "secretsmanager",
		"githb":         "github",
		"zendesk ":      "zendesk",
	} {
		got := Suggest(strings.TrimSpace(typed))
		if len(got) == 0 {
			t.Errorf("%q suggested nothing, want %q", typed, want)

			continue
		}

		if got[0] != want {
			t.Errorf("%q suggested %v, want %q first", typed, got, want)
		}
	}
}

// A name somebody stopped typing is a match whatever the edit distance.
// "secrets" is seven edits from "secretsmanager", well past any budget a
// short name could tolerate, and obviously what they meant.
//
// "gocardless" would not have shown this: it is four edits from
// gocardlessbank and the budget happens to be four, so the distance rule
// finds it with or without the prefix rule. A test that passes either way is
// not a test of the thing it names.
func TestSuggestFindsAPrefixThatIsTooFarToMeasure(t *testing.T) {
	got := Suggest("secrets")

	if len(got) == 0 || got[0] != "secretsmanager" {
		t.Errorf("got %v, want secretsmanager first", got)
	}
}

func TestSuggestSaysNothingForNonsense(t *testing.T) {
	if got := Suggest("qqqqqqqqqqqq"); len(got) != 0 {
		t.Errorf("got %v, want nothing", got)
	}
}

// Three at most. The point is to answer the question, not to print a shorter
// version of the same wall of text.
//
// "sen" is the input that shows it: five names match it and the cap is what
// stops all five being printed. "se" matched two and "mail" matched exactly
// three, so both passed with the cap removed, which made the first two
// versions of this test assertions about nothing.
func TestSuggestIsShort(t *testing.T) {
	got := Suggest("sen")

	if len(got) != 3 {
		t.Errorf("got %d suggestions, want exactly 3: %v", len(got), got)
	}
}

func TestAnUnknownRecipeErrorCarriesTheSuggestion(t *testing.T) {
	err := &ErrNotBundled{Name: "stripo"}

	if !strings.Contains(err.Error(), "did you mean stripe") {
		t.Errorf("got %q", err.Error())
	}
}

// And points somewhere when there is nothing close, rather than saying only
// that the name is wrong.
func TestAnUnknownRecipeErrorPointsSomewhereWhenNothingIsClose(t *testing.T) {
	err := &ErrNotBundled{Name: "qqqqqqqqqqqq"}

	if !strings.Contains(err.Error(), "cauldron recipe list") {
		t.Errorf("got %q", err.Error())
	}
}

func TestAnUnknownRecipeErrorNeverPrintsEverything(t *testing.T) {
	// The failure this replaced. Whatever else the message does, it must not
	// be a hundred and twenty-seven names long.
	for _, name := range []string{"stripo", "qqqqqqqqqqqq", "s"} {
		err := &ErrNotBundled{Name: name}

		if commas := strings.Count(err.Error(), ","); commas > 2 {
			t.Errorf("%q produced a list of %d: %q", name, commas+1, err.Error())
		}
	}
}
