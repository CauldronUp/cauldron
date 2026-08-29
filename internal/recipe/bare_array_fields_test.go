// A bare array of one, and the fields that cannot ride in it.

package recipe

import (
	"strings"
	"testing"
)

// A bare array of one is a real provider shape -- PoetryDB answers
// /title/Ozymandias with [{...}] -- and the Recipe-wide success fields are
// object properties. Both at once asks for two shapes: the runtime would have
// to drop one of them in silence, so validation refuses the pair instead.
func TestABareArrayCannotAlsoCarryTheSuccessFields(t *testing.T) {
	body := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
responses:
  success:
    fields:
      ok: true
resources:
  poem:
    id:
      style: opaque
      field: "-"
      carried_by: title
    fields:
      title:
        type: string
routes:
  - method: GET
    path: /title/{id}
    resource: poem
    operation: get
    envelope:
      array: true
`

	if got := problems(t, body); !strings.Contains(got, "nowhere to go in a list") {
		t.Errorf("expected the pair to be refused, got:\n%s", got)
	}

	// Wrapping the array gives the fields somewhere to live, so that pair is
	// fine and is what Xero and Ghost already ship.
	wrapped := strings.Replace(body, "      array: true", "      style: wrapped\n      key: poems\n      array: true", 1)

	if _, err := parse(t, wrapped); err != nil && strings.Contains(err.Error(), "nowhere to go in a list") {
		t.Errorf("a wrapped array was refused and should not be:\n%v", err)
	}
}
