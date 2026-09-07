package recipe

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Recipes cite each other constantly, and the citation is the point: a finding
// is worth more when it is the third provider to do something than when it is
// the first, and the link is what lets a reader check that claim.
//
// A link to a Recipe that does not ship is worse than no link. It reads as
// evidence, renders as a link on the repository page, and goes nowhere -- and
// it happens naturally while writing several Recipes at once, when a note cites
// a sibling that has not been committed yet or one whose directory ended up
// named something else. Five of them arrived that way in one batch.
//
// So every ../name in a Recipe's own README or recipe.yaml has to name a
// directory under recipes/.
func TestEverySiblingLinkNamesARecipeThatShips(t *testing.T) {
	names := Bundled()

	ships := make(map[string]bool, len(names))
	for _, name := range names {
		ships[name] = true
	}

	link := regexp.MustCompile(`\]\(\.\./([A-Za-z0-9_.-]+)\)`)

	for _, name := range names {
		for _, file := range []string{"README.md", "recipe.yaml"} {
			path := filepath.Join("..", "..", "recipes", name, file)

			body, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("read %s: %v", path, err)

				continue
			}

			for _, m := range link.FindAllStringSubmatch(string(body), -1) {
				target := m[1]

				// A link to the collection itself, rather than to a sibling.
				if target == "." || strings.HasPrefix(target, ".") {
					continue
				}

				if !ships[target] {
					t.Errorf("recipes/%s/%s links to ../%s and no such Recipe ships", name, file, target)
				}
			}
		}
	}
}
