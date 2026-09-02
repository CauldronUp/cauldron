package recipe_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The README states how many conformance cases have been run against a live
// provider, and a test has held that total since it drifted once. Underneath
// it the same figure is broken out by provider -- five OpenRouter's, five the
// npm registry's, and so on -- and nothing held the parts.
//
// A total alone is the wrong thing to hold. Adding a verified case to GitHub
// moves the total, somebody corrects the total, and the sentence saying three
// are GitHub's goes on saying three. The parts are also where the reader
// actually looks: the total is a number, the split is the evidence.
func TestTheREADMEBreakdownOfLiveCasesMatchesTheRecipes(t *testing.T) {
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	text := string(readme)

	const opens = "are the cases whose provider can be asked without a key"
	const closes = "Every other provider needs an account"

	start := strings.Index(text, opens)
	end := strings.Index(text, closes)

	if start < 0 || end < start {
		t.Fatal("the README no longer breaks the live-verified cases down by provider in the form it did; update this test with it")
	}

	paragraph := text[start:end]

	observed := map[string]int{}
	total := 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, c := range r.Conformance {
			if c.Verified != "" {
				observed[name]++
				total++
			}
		}
	}

	// A dot is part of a name. crates.io is the first provider here whose name
	// carries one, and without it in the class the claim "Eleven are crates.io's"
	// matched nothing -- the sum caught that too, which is the second time it
	// has earned its place.
	//
	// And the possessive s is optional, because a name already ending in one
	// takes a bare apostrophe: "Nine are RubyGems'". That is the third time,
	// and the pattern is now that every new provider tests a corner of English
	// the last one did not.
	//
	// The paragraph is wrapped, so a claim can be split across a line break:
	// "Five are" ends one line and "OpenRouter's" begins the next. Matching a
	// literal space silently found four of the five claims and the sum caught it,
	// which is the whole argument for checking the sum as well as the parts.
	//
	// And a hyphen is part of a name, because Open-Meteo has one. That is the
	// fourth corner of English a new provider has found in this one expression,
	// after the dot in crates.io, the bare apostrophe on RubyGems' and the line
	// break in the middle of a claim -- and the fourth time the sum was what
	// noticed, because a claim that does not match is simply not counted.
	// Digits belong in the name half, including at the front. n8n,
	// what3words and 100ms are all real providers, and none of them could be
	// attributed while the pattern allowed only letters -- the paragraph
	// would name them, the parser would skip them, and the total would
	// silently fail to add up with nothing saying which provider was
	// missing. 100ms needed the leading position too, since its name begins
	// with a digit and the product is called that.
	claims := regexp.MustCompile(`([A-Z][a-z]+)[[:space:]]+(?:is|are)[[:space:]]+(?:the[[:space:]]+)?([A-Za-z0-9][A-Za-z0-9.[:space:]-]*?)'s?`).FindAllStringSubmatch(paragraph, -1)
	if len(claims) == 0 {
		t.Fatal("no provider counts were found in the paragraph; update this test with the form it takes now")
	}

	bundled := map[string]string{}
	for _, name := range recipe.Bundled() {
		bundled[simplify(name)] = name
	}

	claimed := 0

	for _, claim := range claims {
		count, ok := words[claim[1]]
		if !ok {
			t.Errorf("the README says %q cases are %s's and this test cannot read that as a number", claim[1], claim[2])
			continue
		}

		name, ok := bundled[simplify(claim[2])]
		if !ok {
			t.Errorf("the README attributes %d live-verified cases to %s and no such Recipe ships", count, claim[2])
			continue
		}

		if observed[name] != count {
			t.Errorf("the README says %d of the live-verified cases are %s's and %d of them are", count, claim[2], observed[name])
		}

		claimed += count
	}

	// The parts have to account for the whole. This is what notices a live
	// case added to a provider the paragraph does not mention: every named
	// count still agrees, and the sentence introducing them is short.
	if claimed != total {
		t.Errorf("the README breaks the live-verified cases into %d and there are %d", claimed, total)
	}
}
