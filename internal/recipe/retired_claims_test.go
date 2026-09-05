package recipe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A Recipe may not assert a limitation this engine no longer has.
//
// Ten Recipes were found doing it in one afternoon, and the shape of the
// mistake is the same every time. Somebody reads a provider, finds behaviour
// the format cannot express, records that carefully in prose, and declines to
// model it. Later the format grows the mechanism. The Recipe is not touched,
// because nothing about it is failing -- and the sentence goes on telling
// every reader that something impossible is impossible.
//
// Three Recipes claimed Cauldron sends no Link header while declaring
// link_header in the same file. Alpaca skipped an endpoint because "the format
// has no way to describe an endpoint that answers with one object and takes no
// identifier", and the format had two. Unleash refused a credential Unleash
// accepts, with a case pinning the refusal. Jotform, NewsAPI, Mailchimp,
// Mezmo and Textline each served one carrier of a secret their provider takes
// several ways.
//
// A retired sentence may still appear where the Recipe is saying it is
// retired -- several now do, because the history is worth keeping. What it may
// not do is stand as the reason for a present decision, and the marker words
// below are how the two are told apart.
func TestNoRecipeAssertsARetiredLimitation(t *testing.T) {
	retired := []struct {
		phrase, mechanism string
	}{
		{"one scheme for a whole Recipe", "a route may carry its own auth block"},
		{"one setting for a whole Recipe", "a route may carry its own auth block"},
		{"one credential scheme for a whole Recipe", "a route may carry its own auth block"},
		{"one mechanism per Recipe", "auth.also, and route-scoped auth"},
		{"one transport per Recipe", "auth.also"},
		{"one pass-or-fail gate for a whole Recipe", "a route may carry its own auth block"},
		{"no per-error override", "errors.<name>.type_field"},
		{"does not send Link headers", "responses.list.link_header"},
		{"sends no Link header", "responses.list.link_header"},
		{"Cauldron does not send one", "responses.list.link_header"},
		{"has no way to describe an endpoint that answers with one object", "collapse_single, or id_from: auth"},
	}

	// Words a Recipe uses when it is recording that a claim has been overtaken
	// rather than making it.
	markers := []string{
		"used to", "no longer", "true when", "stopped being true",
		"overtaken", "until that field existed", "and the format had",
		"is not why", "not for the reason", "reason recorded here",
		"which is no longer", "and it was closed",
	}

	names := Bundled()
	if len(names) == 0 {
		t.Fatal("no bundled Recipes, so this test is checking nothing")
	}

	for _, name := range names {
		path := filepath.Join("..", "..", "recipes", name, "recipe.yaml")

		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		lines := strings.Split(string(contents), "\n")

		for i, line := range lines {
			for _, claim := range retired {
				if !strings.Contains(line, claim.phrase) {
					continue
				}

				// The retirement is stated near the claim, not necessarily on
				// the same line -- these are wrapped comment blocks.
				window := strings.ToLower(strings.Join(lines[max(0, i-8):min(len(lines), i+9)], " "))

				retiredHere := false
				for _, marker := range markers {
					if strings.Contains(window, marker) {
						retiredHere = true

						break
					}
				}

				if !retiredHere {
					t.Errorf("%s:%d says %q, and %s exists now. Model it, or say which reason still applies.",
						name, i+1, claim.phrase, claim.mechanism)
				}
			}
		}
	}
}
