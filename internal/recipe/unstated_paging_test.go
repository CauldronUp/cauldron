package recipe_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A listing that declares no paging, or names only half of it, must say why.
//
// 344 listings across 222 Recipes once declared nothing, and silence covered two
// completely different situations: "this provider serves the whole collection"
// and "nobody has looked yet". `pagination.style: none` separated the first out.
// This test is about what is left of the second.
//
// A Recipe that could not establish a provider's paging has something worth
// writing down anyway -- which addresses 404'd, which reference renders from
// JavaScript, which host answers 403 to anything that is not a browser, which
// endpoint answered 429. Three of those searches turned up something larger than
// the paging: Basiq's own cited spec does not contain the route its Recipe
// models, Midtrans's listing does not appear in its provider's documentation at
// all, and DHL's changelog confirms it paginates while no page names the
// parameters.
//
// The next person to look starts from what already failed rather than repeating
// it, which is the whole value. So an undeclared listing needs a note, and this
// test is what stops one arriving without one.
func TestEveryUnstatedListingSaysWhatWasTried(t *testing.T) {
	// The phrases a Recipe uses to open that record. Both are load-bearing:
	// the first for a provider that could not be read, the second for the one
	// case where the obstacle is this format rather than the provider.
	markers := []string{
		"Paging is still unstated here",
		"Paging is deliberately not declared",
		"deliberately not declared on any of",
		"Still unstated, and this is what was tried",
		"Still unnamed, and this is what was tried",
		// Salesforce's is the one shape neither phrase fits: its position
		// does come back, as a whole nextRecordsUrl, and refusing the
		// parameter with "-" refused the next page with it. The Recipe
		// explains that where the declaration is.
		"cursor_param is left alone rather than set to",
	}

	names := recipe.Bundled()
	if len(names) == 0 {
		t.Fatal("no bundled Recipes, so this test is checking nothing")
	}

	for _, name := range names {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if r.UnstatedPagination() == 0 && r.GuessedPagination() == 0 {
			continue
		}

		path := filepath.Join("..", "..", "recipes", name, "recipe.yaml")

		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		said := false

		for _, marker := range markers {
			if strings.Contains(string(contents), marker) {
				said = true

				break
			}
		}

		if !said {
			t.Errorf("%s has a listing that declares no paging, or names only half of it, "+
				"and no record of what was tried. Declare the parameters, or write down "+
				"which addresses failed and when", name)
		}
	}
}
