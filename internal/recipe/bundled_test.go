package recipe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// recipesDir is the first-party Recipe collection, relative to this package.
const recipesDir = "../../recipes"

// TestBundledRecipesAreValid is the gate that keeps shipped Recipes honest.
// Every Recipe in the repository must parse and validate, so a broken Recipe
// cannot reach a user through a merge.
func TestBundledRecipesAreValid(t *testing.T) {
	entries, err := os.ReadDir(recipesDir)
	if err != nil {
		t.Fatalf("read %s: %v", recipesDir, err)
	}

	found := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(recipesDir, entry.Name(), "recipe.yaml")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s has no recipe.yaml", entry.Name())
			continue
		}

		found++

		t.Run(entry.Name(), func(t *testing.T) {
			r, err := recipe.Load(path)
			if err != nil {
				t.Fatalf("%s: %v", path, err)
			}

			// The directory name is the Recipe's identity everywhere else —
			// `cauldron recipe add stripe` has to find it.
			if r.Name != entry.Name() {
				t.Errorf("recipe name %q does not match directory %q", r.Name, entry.Name())
			}

			if len(r.Fixtures) == 0 {
				t.Error("a bundled Recipe should ship at least one fixture")
			}

			if _, ok := r.Fixtures["empty"]; !ok {
				t.Error(`every Recipe should ship an "empty" fixture so tests can start from nothing`)
			}
		})
	}

	if found == 0 {
		t.Fatal("no bundled recipes found — this test would pass vacuously")
	}
}

// A Recipe without conformance cases is a claim without evidence. This is the
// gate that stops a new Recipe shipping on assertion alone.
func TestEveryBundledRecipeCarriesConformanceCases(t *testing.T) {
	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		if len(r.Conformance) == 0 {
			t.Errorf("%s ships no conformance cases", name)
		}

		for _, c := range r.Conformance {
			if c.Source == "" {
				t.Errorf("%s: case %q cites no source", name, c.Name)
			}
		}
	}
}
