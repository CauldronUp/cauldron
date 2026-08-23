package recipe_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The README says webhook payload envelopes are declarable and that Recipes
// declaring none fall back to Stripe's shape. That was true and unquantified,
// which in this project is half a claim: the sentence reads like a small
// caveat and the split is 5 to 94.
//
// It matters more than the other counts here, because it is the one figure a
// reader would use to decide whether a webhook payload from this emulator
// resembles the provider. Left as prose it can only get quietly worse as
// Recipes are added.
//
// The sentence deliberately avoids the phrasing "N of the M Recipes":
// internal/detect pins detection coverage by matching every occurrence of
// that shape in the README, on purpose, and a sentence about envelopes that
// happened to share it was read as a second statement of that figure.
func TestTheREADMEEnvelopeSplitMatchesTheRecipes(t *testing.T) {
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	stated := regexp.MustCompile(`Of the (\d+) Recipes that emit events, (\d+) declare one and (\d+) fall back`).FindStringSubmatch(string(readme))
	if stated == nil {
		t.Fatal("the README no longer states the envelope split in the form it did; update this test with it")
	}

	declared, emitting := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if len(r.Webhooks.Events) == 0 {
			continue
		}

		emitting++

		if len(r.Webhooks.Payload) > 0 {
			declared++
		}
	}

	if got, _ := strconv.Atoi(stated[1]); got != emitting {
		t.Errorf("the README says %d Recipes emit events and %d do", got, emitting)
	}

	if got, _ := strconv.Atoi(stated[2]); got != declared {
		t.Errorf("the README says %d Recipes declare a payload envelope and %d do", got, declared)
	}

	if got, _ := strconv.Atoi(stated[3]); got != emitting-declared {
		t.Errorf("the README says %d Recipes fall back to the default shape and %d do", got, emitting-declared)
	}
}
