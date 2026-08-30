package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A fingerprint over a whole description is noise. Providers republish their
// OpenAPI documents constantly -- a reworded summary, a new example, an added
// endpoint in a product this Recipe has never heard of -- and a checksum over
// the file reports every one of those as drift. A scan that cries wolf on
// every publish gets switched off, and then the one change that mattered
// arrives unannounced.
//
// So the fingerprint covers the intersection: the paths and methods the Recipe
// declares routes for, the response codes it claims, and the field names it
// names. A change outside that is not this Recipe's business, and a change
// inside it is exactly the thing a Recipe cannot notice on its own, because
// its own conformance suite asserts whatever the Recipe already says.
const driftSpec = `
openapi: 3.0.0
info: {title: Things, version: "1"}
paths:
  /v1/things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  id: {type: string}
                  name: {type: string}
        "404": {description: gone}
  /v1/unrelated:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  whatever: {type: string}
`

func driftRecipe() *recipe.Recipe {
	return &recipe.Recipe{
		Name: "things",
		Resources: map[string]recipe.Resource{
			"thing": {Fields: map[string]recipe.Field{
				"name": {Type: "string"},
			}},
		},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
		Errors: map[string]recipe.Error{"not_found": {Status: 404}},
	}
}

func fingerprintOf(t *testing.T, raw string) string {
	t.Helper()

	doc, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	return Fingerprint(driftRecipe(), doc, "")
}

func TestTheFingerprintIsStableAcrossReadings(t *testing.T) {
	first := fingerprintOf(t, driftSpec)
	second := fingerprintOf(t, driftSpec)

	if first != second {
		t.Errorf("the same document fingerprinted twice gave %q and %q", first, second)
	}

	if first == "" {
		t.Error("the fingerprint is empty, so every comparison would pass")
	}
}

// The reason to fingerprint the intersection rather than the file. A provider
// adding an endpoint this Recipe says nothing about has not changed anything
// this Recipe claims, and reporting it as drift is how a scan becomes noise
// nobody reads.
func TestAPathTheRecipeDoesNotClaimDoesNotMoveTheFingerprint(t *testing.T) {
	grown := strings.Replace(driftSpec, `  /v1/unrelated:`, `  /v1/added:
    post:
      responses:
        "201": {description: made}
  /v1/unrelated:`, 1)

	if before, after := fingerprintOf(t, driftSpec), fingerprintOf(t, grown); before != after {
		t.Errorf("adding an unclaimed path moved the fingerprint from %q to %q", before, after)
	}
}

// And the same for a field: the Recipe names "name", so a sibling appearing
// beside it in a schema the Recipe already models less of is not drift.
func TestAFieldTheRecipeDoesNotNameDoesNotMoveTheFingerprint(t *testing.T) {
	grown := strings.Replace(driftSpec,
		"                  name: {type: string}",
		"                  name: {type: string}\n                  colour: {type: string}", 1)

	if before, after := fingerprintOf(t, driftSpec), fingerprintOf(t, grown); before != after {
		t.Errorf("adding an unclaimed field moved the fingerprint from %q to %q", before, after)
	}
}

// The changes that must move it. Each of these is a thing the Recipe asserts
// and the provider has altered underneath it, and each is invisible to the
// Recipe's own conformance suite, which asserts what the Recipe says rather
// than what the provider does.
func TestAClaimedFieldChangingMovesTheFingerprint(t *testing.T) {
	for _, change := range []struct {
		what string
		from string
		to   string
	}{
		{
			"a claimed field changes type",
			"                  name: {type: string}",
			"                  name: {type: integer}",
		},
		{
			"a claimed field is renamed",
			"                  name: {type: string}",
			"                  title: {type: string}",
		},
		{
			"a claimed path is renamed",
			"  /v1/things:",
			"  /v1/items:",
		},
		{
			"a claimed status disappears",
			`        "404": {description: gone}`,
			`        "410": {description: gone}`,
		},
	} {
		moved := strings.Replace(driftSpec, change.from, change.to, 1)

		if moved == driftSpec {
			t.Fatalf("%s: the test's own replacement did not apply", change.what)
		}

		if before, after := fingerprintOf(t, driftSpec), fingerprintOf(t, moved); before == after {
			t.Errorf("%s: the fingerprint did not move from %q", change.what, before)
		}
	}
}

// Two Recipes over one description claim different things, so they cannot
// share a fingerprint -- otherwise a Recipe could be handed another's recorded
// value and report itself unchanged.
func TestTwoRecipesOverOneDocumentFingerprintDifferently(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	other := driftRecipe()
	other.Routes = []recipe.Route{
		{Method: "GET", Path: "/v1/unrelated", Resource: "thing", Operation: "list"},
	}

	if mine, theirs := Fingerprint(driftRecipe(), doc, ""), Fingerprint(other, doc, ""); mine == theirs {
		t.Errorf("two Recipes claiming different paths share the fingerprint %q", mine)
	}
}
