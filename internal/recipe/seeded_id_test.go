package recipe

import (
	"strings"
	"testing"
)

// A fixture that seeds an identifier in a shape the generator would never mint
// puts two populations in one collection: the records the Recipe shipped and
// the records a run created. Client code that validates an id -- a regex, a
// prefix match, a length check, all of which real code does -- passes against
// one and fails against the other.
//
// A surviving mutation found this. Monday declares ten-digit identifiers and
// seeds ten-digit identifiers, and changing the declaration to six broke
// nothing at all, because every conformance case reads a seeded record and the
// declared shape never reached the wire. It could have disagreed with the
// provider from the first commit and no case would have said so.

func recipeWithID(declared, seeded string) string {
	return `
recipe: ids
capability: storage
version: 0.1.0
upstream:
  api: "1"
resources:
  thing:
    collection: things
    id:
      ` + declared + `
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/things
    resource: thing
    operation: list
fixtures:
  one:
    thing:
      - id: ` + seeded + `
        name: one
conformance:
  - name: a thing has a name
    source: https://docs.ids.test
    fixture: one
    request:
      method: GET
      path: /v1/things
    expect:
      status: 200
      body:
        things[0].name: one
`
}

func TestASeededIDOfTheWrongLengthIsRefused(t *testing.T) {
	_, err := Parse([]byte(recipeWithID("style: digits\n      length: 19", `"9821740065"`)))
	if err == nil {
		t.Fatal("a ten-digit id under a nineteen-digit declaration was accepted")
	}

	if !strings.Contains(err.Error(), "10 digits, not the declared 19") {
		t.Errorf("the message should name both lengths, got: %v", err)
	}
}

func TestASeededIDOutsideTheAlphabetIsRefused(t *testing.T) {
	// Mailchimp's member id is the MD5 of the lower-cased email address, and
	// no MD5 produces a g or an r. Both were seeded, under a hex declaration,
	// for as long as the Recipe had existed.
	_, err := Parse([]byte(recipeWithID("style: hex\n      length: 32", `"0000000000000000000000000000grc2"`)))
	if err == nil {
		t.Fatal("a non-hexadecimal id under a hex declaration was accepted")
	}
}

func TestASeededSnowflakeMayNotBeginWithAZero(t *testing.T) {
	// A snowflake is a number written as text. The generator never mints a
	// leading zero because a number does not have one, and an id that does
	// would compare unequal to the same id parsed and printed again.
	if _, ok := seededID(ID{Style: "digits", Length: 10}, "0821740065"); ok {
		t.Error("a digits id beginning with a zero was accepted")
	}
}

func TestASeededIDMissingItsPrefixIsRefused(t *testing.T) {
	if _, ok := seededID(ID{Style: "prefixed", Prefix: "cus_"}, "sub_9f3a"); ok {
		t.Error("an id carrying another resource's prefix was accepted")
	}
}

func TestAnIdentifierTheProviderNeverSendsIsFree(t *testing.T) {
	// Cohere keys its embeddings e1 and e2 to find them again and emits
	// neither. An id no client can see cannot disagree with anything, so its
	// shape carries no claim to check.
	if _, err := Parse([]byte(recipeWithID("style: uuid\n      field: \"-\"", "e1"))); err != nil {
		t.Errorf("an unemitted identifier was checked anyway: %v", err)
	}
}

func TestTheRightIDShapesAreAccepted(t *testing.T) {
	for _, ok := range []struct{ declared, seeded string }{
		{"style: digits\n      length: 10", `"9821740065"`},
		{"style: hex\n      length: 24", `"000000000000000000000a01"`},
		{"style: uuid", `"3f2504e0-4f89-41d3-9a0c-0305e82c3301"`},
		{"style: numeric", `"42"`},
		{"style: prefixed\n      prefix: cus_", "cus_9f3a"},
		{"style: opaque\n      length: 14", "comp0000000001"},
	} {
		if _, err := Parse([]byte(recipeWithID(ok.declared, ok.seeded))); err != nil {
			t.Errorf("%s seeded with %s was refused: %v", ok.declared, ok.seeded, err)
		}
	}
}
