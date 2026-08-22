package detect

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The backlog keeps a table of Recipes no dependency maps to, with a note on
// each saying what was looked for. It is the answer to "why is this one not
// detected", and it drifted badly once already: it named twenty-six while
// fifteen of them had been mapped, and quoted a coverage figure two passes
// out of date.
//
// A list of gaps that is itself wrong is worse than no list, because it sends
// somebody to check a provider that needs no checking -- so the table has to
// name exactly the Recipes nothing maps to.
func TestTheBacklogListsExactlyTheUnmappedRecipes(t *testing.T) {
	raw, err := os.ReadFile("../../docs/backlog.md")
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	mapped := map[string]bool{}
	for _, p := range providers() {
		mapped[p.recipe] = true
	}

	unmapped := map[string]bool{}
	for _, name := range recipe.Bundled() {
		if !mapped[name] {
			unmapped[simplify(name)] = true
		}
	}

	// The table under the Detection coverage heading, which is the last one
	// in that section and the only one whose second column is a note.
	text := string(raw)

	start := strings.Index(text, "## Detection coverage")
	if start < 0 {
		t.Fatal("the backlog no longer has a Detection coverage section; update this test with it")
	}

	end := strings.Index(text[start:], "**OpenAI is mapped")
	if end < 0 {
		t.Fatal("the Detection coverage section no longer ends where this test expects")
	}

	section := text[start : start+end]
	row := regexp.MustCompile(`(?m)^\| ([^|]+?) \| ([^|]*)\|`)

	listed := map[string]bool{}

	for _, m := range row.FindAllStringSubmatch(section, -1) {
		name := simplify(strings.TrimSpace(m[1]))
		if name == "recipe" || name == "" || strings.HasPrefix(name, "---") {
			continue
		}

		// The collided-names table names npm packages, not Recipes.
		if strings.Contains(m[2], "actually is") || strings.HasPrefix(strings.TrimSpace(m[1]), "`") {
			continue
		}

		listed[name] = true
	}

	for name := range unmapped {
		if !listed[name] {
			t.Errorf("nothing maps to %q and the backlog does not list it", name)
		}
	}

	for name := range listed {
		if !unmapped[name] {
			t.Errorf("the backlog lists %q as unmapped and a dependency maps to it", name)
		}
	}
}

// simplify reduces a name to letters and digits, so "Bill.com" and "billcom"
// and "incident.io" and "incidentio" compare equal.
func simplify(name string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}

	return b.String()
}
