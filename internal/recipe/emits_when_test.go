package recipe

import (
	"strings"
	"testing"
)

// A conditional emission fails quietly in a way an unconditional one does not.
// When emits is wrong the event never arrives and a case asserting it goes
// red; when emits_when is wrong the event also never arrives, and that is
// indistinguishable from the field simply not having changed. Every one of
// these mistakes looks like correct behaviour from the outside, so the
// validator has to catch them at the declaration instead.

const conditionalRecipe = `
recipe: conditional
capability: support
version: 0.1.0
upstream:
  api: "1"
resources:
  ticket:
    collection: tickets
    id:
      prefix: tkt_
      length: 14
    fields:
      subject:
        type: string
      status:
        type: string
      tags:
        type: list
routes:
  - method: PATCH
    path: /v1/tickets/{id}
    resource: ticket
    operation: update
    emits_when:
      - event: ticket.status_changed
        field: status
webhooks:
  events:
    - ticket.status_changed
  signing:
    scheme: none
conformance:
  - name: moving the status announces it
    source: https://example.invalid/docs
    fixture: small
    request:
      method: PATCH
      path: /v1/tickets/tkt_00000000000001
      json:
        status: closed
    expect:
      status: 200
      webhook:
        event: ticket.status_changed
fixtures:
  small:
    ticket:
      - id: tkt_00000000000001
        subject: the printer again
        status: open
`

func TestAConditionalEmissionIsAccepted(t *testing.T) {
	if _, err := parse(t, conditionalRecipe); err != nil {
		t.Fatalf("the base recipe should be valid: %v", err)
	}
}

// Naming an event nothing declares is the same mistake emits already guards
// against, and it arrives here by the same route: a typo.
func TestAConditionalEventMustBeDeclared(t *testing.T) {
	yaml := strings.Replace(conditionalRecipe, "      - event: ticket.status_changed", "      - event: ticket.status_change", 1)

	if got := problems(t, yaml); !strings.Contains(got, "webhooks.events does not list it") {
		t.Errorf("got %q", got)
	}
}

// A field the resource does not have is compared absent-to-absent, which is
// equal, so the event never fires and nothing says why.
func TestAConditionalFieldMustBeDeclared(t *testing.T) {
	yaml := strings.Replace(conditionalRecipe, "        field: status", "        field: state", 1)

	if got := problems(t, yaml); !strings.Contains(got, "the comparison is between two absences") {
		t.Errorf("got %q", got)
	}
}

// The comparison needs a write to sit either side of, and only update makes
// one. A create has no before and a delete has no after.
func TestAConditionalEmissionNeedsAnUpdateRoute(t *testing.T) {
	yaml := strings.Replace(conditionalRecipe, "    operation: update", "    operation: create", 1)

	if got := problems(t, yaml); !strings.Contains(got, "which compares a field before and after a write that this route does not make") {
		t.Errorf("got %q", got)
	}
}

// Comparing a composite would announce a change every time a list came back
// in a different order, which is not what the provider calls a change.
func TestAConditionalFieldMustBeScalar(t *testing.T) {
	yaml := strings.Replace(conditionalRecipe, "        field: status", "        field: tags", 1)

	if got := problems(t, yaml); !strings.Contains(got, "these events key off a scalar moving") {
		t.Errorf("got %q", got)
	}
}

func TestAConditionalEmissionNeedsBothHalves(t *testing.T) {
	yaml := strings.Replace(conditionalRecipe, "        field: status", "", 1)

	if got := problems(t, yaml); !strings.Contains(got, "needing both event and field") {
		t.Errorf("got %q", got)
	}
}
