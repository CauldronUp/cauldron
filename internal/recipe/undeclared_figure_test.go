package recipe

import (
	"os"
	"regexp"
	"strconv"
	"testing"
)

// The README states, twice, how many Recipes name no machine-readable
// description. Both numbers were wrong, and had been for a long time.
//
// `undeclared` is a drift status and it means one thing: Upstream.Spec is
// empty. The README's figure for it had been tracking the *detection* coverage
// figure instead -- 501 of 627, then 506 of 651, then 521 of 668, always equal
// to how many Recipes a dependency maps to, because every Recipe that shipped
// bumped both numbers whether or not it declared a description. The real count
// sat 58 lower and nothing noticed, over at least eighty commits.
//
// The detection guard next door checks every "N of the M Recipes" in the file
// and its own comment says why it looks at all of them rather than the first:
// "a check that looks at one occurrence of a repeated claim makes the other one
// safe to forget." Both undeclared sentences are phrased without the word
// Recipes -- "521 of the 668, and falling", "521 of the 668 here" -- so that
// regex never reached either of them, and the claim they make was the one
// nobody was checking.
//
// So this counts the Recipes with no spec and holds the prose to it.
func TestTheREADMEStatesHowManyRecipesDeclareNoDescription(t *testing.T) {
	names := Bundled()

	undeclared := 0

	for _, name := range names {
		r, err := Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if r.Upstream.Spec == "" {
			undeclared++
		}
	}

	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	// Both phrasings, and the count of them, so that adding a third place the
	// figure is stated without adding it here is itself a failure rather than
	// a silent gap.
	claims := regexp.MustCompile(`(\d+) of the (\d+)(?:, and falling| here)`).
		FindAllStringSubmatch(string(readme), -1)

	if len(claims) != 2 {
		t.Fatalf("the README states the undeclared figure %d times, want 2; if a place was added or reworded, update this test with it", len(claims))
	}

	for _, claim := range claims {
		if got, _ := strconv.Atoi(claim[1]); got != undeclared {
			t.Errorf("the README says %d Recipes declare no description and %d do", got, undeclared)
		}

		if got, _ := strconv.Atoi(claim[2]); got != len(names) {
			t.Errorf("the README states that out of %d Recipes and %d ship", got, len(names))
		}
	}
}
