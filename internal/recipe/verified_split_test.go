package recipe_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Every Recipe carries its own README, and it says how many of its cases were
// checked against the live provider.
//
// That count used to live in the top-level README, broken out by provider in
// one enormous run of paragraphs, and a test held every part of it. The parts
// were the right thing to hold -- adding a verified case to GitHub moves the
// total, somebody corrects the total, and the sentence saying three are
// GitHub's goes on saying three -- but the place was wrong. Three thousand
// lines of per-provider notes sat in a file nobody reads to the end, and the
// one note a person wants is the one for the provider they are about to fake.
//
// So the notes moved next to their Recipes and this moved with them. The claim
// is now local to each file, which makes it both easier to write and harder to
// forget: a Recipe with no README fails here, and a README claiming a count its
// Recipe does not have fails here too.
func TestEveryRecipeREADMEStatesItsOwnVerifiedCount(t *testing.T) {
	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		cases, verified := 0, 0

		for _, c := range r.Conformance {
			cases++

			if c.Verified != "" {
				verified++
			}
		}

		path := filepath.Join("..", "..", "recipes", name, "README.md")

		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s has no README beside its Recipe: %v", name, err)

			continue
		}

		text := string(body)

		// The wording differs with the shape of the answer -- all of them, some
		// of them, none of them -- so the check is on the numbers rather than
		// on one sentence, which leaves the prose free to read well.
		want := fmt.Sprintf("%d conformance case", cases)
		if cases > 0 && !strings.Contains(text, want) {
			t.Errorf("%s's README does not say it has %d conformance cases", name, cases)
		}

		switch {
		case verified == 0 && cases > 0:
			if !strings.Contains(text, "none checked against a live API") {
				t.Errorf("%s has no live-checked case and its README does not say so", name)
			}
		case verified == cases:
			if !strings.Contains(text, "all of them checked against the live API") {
				t.Errorf("%s checked all %d of its cases live and its README does not say so", name, cases)
			}
		default:
			claim := fmt.Sprintf("%d checked against the live API", verified)
			if !strings.Contains(text, claim) {
				t.Errorf("%s's README does not say %d of its cases were checked live", name, verified)
			}
		}
	}
}

// The top-level README keeps the total, and says where the parts went.
//
// Holding the total here is what stops the move from quietly losing the claim:
// the number still has to be right, and the sentence pointing at the per-Recipe
// files still has to be there, or somebody deletes the pointer and the notes
// become unreachable from the front page.
func TestTheREADMEPointsAtThePerRecipeNotes(t *testing.T) {
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	text := string(readme)

	if !strings.Contains(text, "recipes/<name>/README.md") {
		t.Error("the README no longer says where a provider's own findings live")
	}

	total := 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, c := range r.Conformance {
			if c.Verified != "" {
				total++
			}
		}
	}

	claim := fmt.Sprintf("All %d are the cases whose provider can be asked without a key", total)
	if !strings.Contains(text, claim) {
		t.Errorf("the README no longer says %d cases were checked live, in the form this reads", total)
	}
}
