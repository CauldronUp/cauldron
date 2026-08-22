package detect

import (
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A mapping may name a Recipe that does not ship. That is how a recognised
// dependency Cauldron cannot emulate gets reported by name rather than falling
// through to the "looks like an API client" heuristic, and OpenAI is there for
// exactly that reason.
//
// What must not happen is a typo passing for an intention. A mapping naming
// "postmarkk" would quietly become a warning about a provider nobody has heard
// of, so every name that does not ship has to be declared as deliberate.
func TestEveryMappedRecipeShipsOrIsDeclaredUnshipped(t *testing.T) {
	shipped := map[string]bool{}
	for _, name := range recipe.Bundled() {
		shipped[name] = true
	}

	for _, p := range providers() {
		if shipped[p.recipe] {
			if _, declared := unshipped[p.recipe]; declared {
				t.Errorf("%q ships now, so remove it from unshipped", p.recipe)
			}

			continue
		}

		if _, declared := unshipped[p.recipe]; !declared {
			t.Errorf("%q is mapped from a dependency, does not ship, and is not declared unshipped", p.recipe)
		}
	}
}

// A package listed under two Recipes resolves to whichever happens to be
// first, and the order is an implementation detail of a slice literal. A wrong
// guess here is worse than no guess: booting the wrong fake sends a developer
// chasing a bug that does not exist.
func TestNoPackageMapsToTwoRecipes(t *testing.T) {
	for _, ecosystem := range []struct {
		name string
		of   func(provider) []string
	}{
		{"composer", func(p provider) []string { return p.composer }},
		{"npm", func(p provider) []string { return p.npm }},
		{"gomod", func(p provider) []string { return p.gomod }},
	} {
		owner := map[string]string{}

		for _, p := range providers() {
			for _, pkg := range ecosystem.of(p) {
				key := strings.ToLower(pkg)

				if first, taken := owner[key]; taken {
					t.Errorf("%s package %q maps to both %q and %q", ecosystem.name, pkg, first, p.recipe)

					continue
				}

				owner[key] = p.recipe
			}
		}
	}
}

func TestNoRecipeIsMappedTwice(t *testing.T) {
	seen := map[string]bool{}

	for _, p := range providers() {
		if seen[p.recipe] {
			t.Errorf("%q has two entries, and only the first is ever reached", p.recipe)
		}

		seen[p.recipe] = true
	}
}

func TestEveryMappingCarriesAtLeastOnePackage(t *testing.T) {
	for _, p := range providers() {
		if len(p.composer)+len(p.npm)+len(p.gomod) == 0 {
			t.Errorf("%q maps from nothing, so it can never be detected", p.recipe)
		}
	}
}

// Detection coverage is the difference between the claim on the front of the
// README and what actually happens in somebody's repository, so the number is
// worth having in front of whoever runs the suite -- and worth checking,
// which it was not.
//
// The README said 91 of 117 while the truth was 91 of 167. The numerator had
// been maintained and the denominator was the Recipe count from whenever the
// line was written, so the claim drifted in the flattering direction: 78 per
// cent coverage stated, 54 per cent real. A capability table that improves on
// its own while nobody is looking is worse than one that says nothing.
func TestReportDetectionCoverage(t *testing.T) {
	mapped := map[string]bool{}
	for _, p := range providers() {
		mapped[p.recipe] = true
	}

	var missing []string

	for _, name := range recipe.Bundled() {
		if !mapped[name] {
			missing = append(missing, name)
		}
	}

	sort.Strings(missing)

	shipped := len(recipe.Bundled())

	reachable := shipped - len(missing)

	t.Logf("%d of %d Recipes are reachable from a dependency", reachable, shipped)

	if len(missing) > 0 {
		t.Logf("no dependency maps to: %s", strings.Join(missing, " "))
	}

	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	// Every place the README states it, not the first. The figure appears
	// twice and correcting one of them left the other saying 91 of 117 for
	// as long as it took to notice, which is the same drift again in a
	// smaller way: a check that looks at one occurrence of a repeated claim
	// makes the other one safe to forget.
	pattern := regexp.MustCompile(`(\d+) of (?:the )?(\d+) Recipes`)

	found := pattern.FindAllStringSubmatch(string(readme), -1)
	if len(found) == 0 {
		t.Fatal("the README no longer states detection coverage in the form it did; update this test with it")
	}

	for _, stated := range found {
		if got, _ := strconv.Atoi(stated[1]); got != reachable {
			t.Errorf("the README says %d Recipes are reachable and %d are", got, reachable)
		}

		if got, _ := strconv.Atoi(stated[2]); got != shipped {
			t.Errorf("the README says coverage is out of %d Recipes and %d ship", got, shipped)
		}
	}

	// And the remainder it quotes beside them.
	remainder := regexp.MustCompile(`The other (\d+) ship`).FindStringSubmatch(string(readme))
	if remainder == nil {
		t.Fatal("the README no longer states how many Recipes nothing maps to; update this test with it")
	}

	if got, _ := strconv.Atoi(remainder[1]); got != len(missing) {
		t.Errorf("the README says nothing maps to %d Recipes and nothing maps to %d", got, len(missing))
	}
}
