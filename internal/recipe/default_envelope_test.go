package recipe

import (
	"strings"
	"testing"
)

// A Recipe that declares no payload template sends the default envelope,
// which is this project's convention and not any provider's shape. A
// conformance case pinning a path inside it is evidence about Cauldron rather
// than about the provider, which is the one thing a conformance case is not
// for.
//
// Square had one and it passed. It asserted data.object.amount_money.amount
// -- a path Square has never sent -- while its name claimed it was checking
// the webhook matched the response. The claim was right, the path was
// Stripe's, and nothing could tell the difference.

const conventionRecipe = `
recipe: convention
capability: crm
version: 0.1.0
upstream:
  api: "1"
resources:
  contact:
    collection: contacts
    id:
      prefix: con_
      length: 14
    fields:
      email:
        type: string
      status:
        type: string
        default: active
routes:
  - method: POST
    path: /v1/contacts
    resource: contact
    operation: create
webhooks:
  events:
    - contact.created
  signing:
    scheme: none
conformance:
  - name: creating a contact announces it
    source: https://example.invalid/docs
    fixture: small
    request:
      method: POST
      path: /v1/contacts
      json:
        email: ada@example.com
    expect:
      status: 200
      webhook:
        event: contact.created
        body:
          data.object.status: active
fixtures:
  small: {}
`

func TestPinningTheDefaultEnvelopeIsRejected(t *testing.T) {
	if got := problems(t, conventionRecipe); !strings.Contains(got, "pins this project's convention") {
		t.Errorf("got %q", got)
	}
}

// Declaring one makes the same path a claim about the provider, so it is
// allowed -- the rule is about the absence of a template, not about the word
// "data", which plenty of providers really do use.
func TestTheSamePathIsFineOnceAnEnvelopeIsDeclared(t *testing.T) {
	yaml := strings.Replace(conventionRecipe,
		"  signing:\n    scheme: none",
		"  payload:\n    data:\n      object:\n        \"{object}\": true\n  signing:\n    scheme: none", 1)

	if _, err := parse(t, yaml); err != nil {
		t.Fatalf("declaring the envelope should make the case valid: %v", err)
	}
}

// A path outside the default envelope says nothing about the convention, so
// it is not what this rule is looking for.
func TestAPathOutsideTheDefaultEnvelopeIsUntouched(t *testing.T) {
	yaml := strings.Replace(conventionRecipe, "          data.object.status: active", "          status: active", 1)

	// parse rather than problems: this one is expected to be valid, and the
	// path being wrong for the default envelope is a case that fails when it
	// runs rather than a Recipe that fails to load.
	if _, err := parse(t, yaml); err != nil {
		t.Fatalf("a top-level path should not trip the rule: %v", err)
	}
}
