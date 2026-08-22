package recipe_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// backlogPath is the queue document, relative to this package.
const backlogPath = "../../docs/backlog.md"

// The backlog keeps a table of providers assessed and deliberately not done.
// It is the one list in this project that nothing ever rereads: a Recipe that
// ships gets struck off the queue because somebody is looking at the queue,
// and nobody looks at the list of things that cannot be done.
//
// Four of its five rows were wrong. Linear and Attio shipped and the table
// still gave the reasons they could not, while the README named both of them
// as deliberately excluded in a paragraph two hundred lines below the list of
// shipped Recipes that included them. New Relic and Railway were kept out for
// being GraphQL-only, which stopped being a reason the day a route learned to
// match on the field a query names.
//
// A reason that has quietly expired reads exactly like one that still holds,
// so this checks the table against what ships rather than trusting it.
func TestTheExclusionTableAgreesWithWhatShips(t *testing.T) {
	backlog, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	section := exclusionSection(string(backlog))
	if section == "" {
		t.Fatal("the backlog no longer has an 'Assessed and deliberately not done' section; update this test with it")
	}

	bundled := map[string]bool{}
	for _, name := range recipe.Bundled() {
		bundled[simplify(name)] = true
	}

	rows := regexp.MustCompile(`(?m)^\| (?:~~)?([^|~]+?)(?:~~)? \|`).FindAllStringSubmatch(section, -1)

	standing := 0

	for _, row := range rows {
		name := strings.TrimSpace(row[1])
		if name == "Provider" {
			continue
		}

		struck := strings.Contains(row[0], "~~")

		switch {
		case struck && !bundled[simplify(name)]:
			t.Errorf("the exclusion table strikes %s off as shipped and no such Recipe ships", name)
		case !struck && bundled[simplify(name)]:
			t.Errorf("the exclusion table says %s was deliberately not done and it ships", name)
		case !struck:
			standing++
		}
	}

	if standing == 0 {
		t.Fatal("no provider in the exclusion table is still excluded; the README paragraph needs rewriting, not this count")
	}

	// The README states the same figure in its own words, and the two drifted
	// apart once already: the table struck Linear and Attio off and the README
	// went on naming them, in the same file that lists them among the Recipes
	// that ship.
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	stated := regexp.MustCompile(`Not every provider fits, and (\w+) (?:is|are) left out deliberately`).FindStringSubmatch(string(readme))
	if stated == nil {
		t.Fatal("the README no longer states how many providers are left out, in the form it did; update this test with it")
	}

	if counted := words[capitalised(stated[1])]; counted != standing {
		t.Errorf("the README says %s provider(s) are left out deliberately and the backlog table stands %d", stated[1], standing)
	}
}

// exclusionSection returns the body of the deliberately-not-done table.
func exclusionSection(backlog string) string {
	const heading = "## Assessed and deliberately not done"

	start := strings.Index(backlog, heading)
	if start < 0 {
		return ""
	}

	rest := backlog[start+len(heading):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return rest[:end]
	}

	return rest
}

// capitalised upper-cases the first letter, so a number spelled mid-sentence
// finds the same entry as one spelled at the start of one.
func capitalised(word string) string {
	if word == "" {
		return word
	}

	return strings.ToUpper(word[:1]) + word[1:]
}
