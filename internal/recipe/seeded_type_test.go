package recipe

import (
	"strings"
	"testing"
)

// A field's type is documentation. Nothing in the runtime reads it, so the
// value a fixture seeds is the value that goes on the wire with whatever type
// the YAML gave it. A Recipe can therefore declare a field an integer, seed it
// with a quoted number, and serve a string to everything that reads it -- and
// its own conformance cases will assert the string and pass, because they were
// written against what came out.
//
// Quoting is the live risk in both directions. Backblaze sends part sizes as
// quoted strings and file sizes as bare numbers, on the same API, so a Recipe
// describing it has to get the quoting right twice and has no help from the
// type declaration in doing so.

func recipeWith(field, seeded string) string {
	return `
recipe: sizes
capability: storage
version: 0.1.0
upstream:
  api: "1"
resources:
  thing:
    collection: things
    id:
      style: opaque
    fields:
      ` + field + `
routes:
  - method: GET
    path: /v1/things
    resource: thing
    operation: list
fixtures:
  one:
    thing:
      - id: t1
        ` + seeded + `
conformance:
  - name: a thing has a size
    source: https://docs.sizes.test
    fixture: one
    request:
      method: GET
      path: /v1/things
    expect:
      status: 200
      body:
        things[0].id: t1
`
}

func TestAQuotedNumberUnderAnIntegerFieldIsRefused(t *testing.T) {
	_, err := Parse([]byte(recipeWith("size:\n        type: integer", `size: "284617"`)))
	if err == nil {
		t.Fatal("an integer field seeded with a quoted number was accepted")
	}

	if !strings.Contains(err.Error(), "a string") {
		t.Errorf("the message should say what went on the wire, got: %v", err)
	}
}

func TestABareNumberUnderAStringFieldIsRefused(t *testing.T) {
	_, err := Parse([]byte(recipeWith("size:\n        type: string", "size: 100000000")))
	if err == nil {
		t.Fatal("a string field seeded with a bare number was accepted")
	}
}

func TestTheRightTypesAreAccepted(t *testing.T) {
	for _, ok := range []struct{ field, seeded string }{
		{"size:\n        type: integer", "size: 284617"},
		{"size:\n        type: string", `size: "100000000"`},
		{"size:\n        type: number", "size: 1.5"},
		{"size:\n        type: number", "size: 2"},
		{"size:\n        type: boolean", "size: true"},
		{"size:\n        type: timestamp_ms", "size: 1755422400000"},
		{"size:\n        type: datetime", `size: "2026-08-17T00:00:00Z"`},
		// An explicit null is a distinct state from both an absent field and
		// a wrong one, and several providers send one.
		{"size:\n        type: integer", "size: null"},
	} {
		if _, err := Parse([]byte(recipeWith(ok.field, ok.seeded))); err != nil {
			t.Errorf("%s with %s was refused: %v", ok.field, ok.seeded, err)
		}
	}
}

func TestOffIsAStringNotABoolean(t *testing.T) {
	// yaml.v3 reads the YAML 1.2 core schema, where only true and false are
	// booleans. YAML 1.1 also took off, on, yes and no, and a checker written
	// against a 1.1 parser reports a DigitalOcean droplet whose status is off
	// as serving false. It does not: the wire carries the word.
	if seeded, ok := seededAs("string", "off"); !ok {
		t.Errorf("off under a string field was reported as %s", seeded)
	}
}
