// Where a provider's own description lives, and which version of it a Recipe
// describes.

package recipe

import (
	"strings"
	"testing"
)

// upstream.api has always held the provider's version label -- v1, v2,
// 2026-06-30 -- and 95 of the Recipes here say v1 while 50 say v2. What it has
// never held is whose. Two Recipes both saying "v1" are describing two
// different APIs, and nothing in the format relates a provider's v1 to its v2
// or records that one replaced the other.
//
// That absence is not academic. The MBTA's v2 clients reach
// realtime.mbta.com, ClinicalTrials.gov's reach clinicaltrials.gov/ct2 and
// DataCite's minting clients reach mds.datacite.org -- three retired
// interfaces whose clients must not be offered a Recipe for the live one. All
// three exclusions currently live in prose, in a comment beside a detection
// entry, where nothing can check them.
func recipeWithUpstream(t *testing.T, upstream string) error {
	t.Helper()

	body := `
recipe: things
capability: docs
version: 0.1.0
upstream:
` + upstream + `
auth:
  scheme: none
resources:
  thing:
    id:
      style: opaque
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/things
    resource: thing
    operation: list
fixtures:
  empty: {}
`

	_, err := Parse([]byte(body))

	return err
}

func TestASpecHashWithoutASpecHashesNothing(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v1
  spec_hash: 0123456789abcdef
  spec_seen: "2026-08-30"`)
	if err == nil {
		t.Fatal("a Recipe recording a fingerprint of a description it does not name was accepted")
	}

	if !strings.Contains(err.Error(), "upstream.spec") {
		t.Errorf("the complaint does not name upstream.spec: %v", err)
	}
}

// Every other piece of evidence in this project carries the date it was taken.
// A fingerprint without one cannot be reasoned about: it says a description
// once looked like this, and not when, so nobody can tell a Recipe checked
// yesterday from one checked in 2023.
func TestAFingerprintWithoutADateIsUndatedEvidence(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v1
  spec: https://example.test/openapi.json
  spec_hash: 0123456789abcdef`)
	if err == nil {
		t.Fatal("a Recipe recording an undated fingerprint was accepted")
	}

	if !strings.Contains(err.Error(), "upstream.spec_seen") {
		t.Errorf("the complaint does not name upstream.spec_seen: %v", err)
	}
}

func TestASpecDateMustBeADate(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v1
  spec: https://example.test/openapi.json
  spec_hash: 0123456789abcdef
  spec_seen: last Tuesday`)
	if err == nil {
		t.Fatal("a Recipe dating its fingerprint in prose was accepted")
	}

	if !strings.Contains(err.Error(), "upstream.spec_seen") {
		t.Errorf("the complaint does not name upstream.spec_seen: %v", err)
	}
}

// A superseded version has to name what identifies it, or it is a note rather
// than a fact anything can act on. The host is what detection has to go on:
// realtime.mbta.com is how an MBTA v2 client is told from a v3 one.
func TestASupersededVersionMustNameWhatIdentifiesIt(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  provider: mbta
  supersedes:
    - version: v2`)
	if err == nil {
		t.Fatal("a superseded version naming nothing that identifies it was accepted")
	}

	if !strings.Contains(err.Error(), "host") {
		t.Errorf("the complaint does not mention the host: %v", err)
	}
}

func TestASupersededVersionCannotBeTheOneBeingDescribed(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  provider: mbta
  supersedes:
    - version: v3
      host: realtime.mbta.com`)
	if err == nil {
		t.Fatal("a Recipe declaring it superseded itself was accepted")
	}

	if !strings.Contains(err.Error(), "supersedes") {
		t.Errorf("the complaint does not name supersedes: %v", err)
	}
}

// Naming a superseded version without naming the provider leaves the two
// versions unrelated, which is the state the format was already in.
func TestSupersedingRequiresSayingWhoseAPIItIs(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  supersedes:
    - version: v2
      host: realtime.mbta.com`)
	if err == nil {
		t.Fatal("a Recipe relating two versions of an unnamed provider was accepted")
	}

	if !strings.Contains(err.Error(), "upstream.provider") {
		t.Errorf("the complaint does not name upstream.provider: %v", err)
	}
}

// And the whole point: a complete declaration is accepted.
func TestAProviderWithAVersionAndASpecIsAccepted(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  provider: mbta
  docs: https://api-v3.mbta.com/docs/swagger/index.html
  spec: https://api-v3.mbta.com/docs/swagger/swagger.json
  spec_hash: 0123456789abcdef
  spec_seen: "2026-08-30"
  supersedes:
    - version: v2
      host: realtime.mbta.com
      note: the retired realtime interface`)
	if err != nil {
		t.Fatalf("a complete upstream declaration was rejected: %v", err)
	}
}

// "-" is the format's word for "the provider does not send it", used by
// id.field, message_field and the list key already. Here it means the
// description exists and cannot be fingerprinted: the MBTA publishes Swagger
// 2.0, which this reads not at all, and will publish it again tomorrow.
//
// The alternative was to leave the URL out, which loses the fact that the
// provider publishes anything, and leaves the Recipe indistinguishable from
// the 285 whose providers publish nothing.
func TestAFingerprintOfNothingIsWrittenAsAHyphen(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  spec: https://api-v3.mbta.com/docs/swagger/swagger.json
  spec_hash: "-"
  spec_seen: "2026-08-30"`)
	if err != nil {
		t.Fatalf("a Recipe recording an unreadable description was rejected: %v", err)
	}
}

// And it still owes a date. "We cannot read this" is a claim about a document
// at a moment, and a provider that moves to OpenAPI 3 makes it wrong.
func TestAnUnreadableDescriptionStillOwesADate(t *testing.T) {
	err := recipeWithUpstream(t, `  api: v3
  spec: https://api-v3.mbta.com/docs/swagger/swagger.json
  spec_hash: "-"`)
	if err == nil {
		t.Fatal("a Recipe recording an undated unreadable description was accepted")
	}

	if !strings.Contains(err.Error(), "upstream.spec_seen") {
		t.Errorf("the complaint does not name upstream.spec_seen: %v", err)
	}
}
