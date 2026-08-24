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

// The check above was doing nothing for more than half the resources it
// applied to, and this is the case that was missed.
//
// A fixture seeds the identifier under "id" whatever the provider calls it on
// the wire -- id.field renames the field in the response, not the key a
// fixture is written with. The lookup used the wire name alone, so for every
// resource that renames it the fixture id was simply not found and the shape
// went unchecked: Twilio's sid, Asana's gid, Calendly's uri, Auth0's user_id,
// 94 resources in all.
//
// Turning it on found four real ones. Adyen seeded upper-case ids under a hex
// declaration, which mints lower case, so seeded and created records had two
// different shapes in one collection. Gusto's payroll fixtures led with p and
// its job fixtures with j, neither of which is a hex digit, so they were not
// UUIDs at all. Shippo's rates, shipments and transactions led with r, s and
// t for the same reason. And Auth0's fixture held a Google user on purpose,
// which is the one case where the fixture was right and the format could not
// say so.
func TestASeededIDIsCheckedWhenTheProviderRenamesTheField(t *testing.T) {
	_, err := Parse([]byte(recipeWithID("style: hex\n      length: 32\n      field: object_id", `"0000000000000000000000000000grc2"`)))
	if err == nil {
		t.Fatal("a non-hexadecimal id was accepted because the resource renames its id field")
	}

	if !strings.Contains(err.Error(), "not hexadecimal") {
		t.Errorf("the message should say what is wrong with the id, got: %v", err)
	}
}

// Auth0's user_id encodes the connection the user came from: auth0|abc is a
// database user and google-oauth2|123 signed in with Google. Code that parses
// the identifier assuming auth0| breaks on the first social login, which is
// the first thing that Recipe says, and its fixture holds a Google user on
// purpose.
//
// Minting still uses one prefix, because there is one shape a new record
// takes. These are the others a real account already contains.
func TestASeededIDMayCarryAPrefixTheRecipeDoesNotMint(t *testing.T) {
	declared := "prefix: \"auth0|\"\n      length: 24\n      field: user_id\n      other_prefixes: [\"google-oauth2|\"]"

	if _, err := Parse([]byte(recipeWithID(declared, `"google-oauth2|107329581004"`))); err != nil {
		t.Errorf("a declared alternative prefix was refused: %v", err)
	}

	// And an undeclared one still is, so the field widens the claim rather
	// than removing it.
	if _, err := Parse([]byte(recipeWithID(declared, `"samlp|107329581004"`))); err == nil {
		t.Error("an undeclared prefix was accepted, so other_prefixes turned the check off rather than widening it")
	}
}
