// A map with no wrapper, and the fields that cannot ride beside it.

package recipe

import (
	"strings"
	"testing"
)

// The bare array's sibling. CoinGecko answers /simple/price with
// {"bitcoin": {...}} at the top level, so the response is the collection and
// there is no object left to hang a property on. Declaring Recipe-wide success
// fields as well asks for two shapes at once, and the runtime would have to
// drop one of them in silence -- or worse, write a property into a map whose
// every other key is an identifier, where a client reading body[id] would find
// it and take it for a record.
func TestAMapWithNoWrapperCannotAlsoCarryTheSuccessFields(t *testing.T) {
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
  list:
    style: map
    key: "-"
resources:
  price:
    id:
      style: opaque
      field: id
    fields:
      usd:
        type: number
routes:
  - method: GET
    path: /simple/price
    resource: price
    operation: list
`

	if got := problems(t, body); !strings.Contains(got, "nowhere to go beside a collection") {
		t.Errorf("expected the pair to be refused, got:\n%s", got)
	}

	// A wrapper gives the fields somewhere to live, which is what Pusher
	// already ships: {"channels": {...}} with room beside it.
	wrapped := strings.Replace(body, "    key: \"-\"", "    key: channels", 1)

	if _, err := parse(t, wrapped); err != nil {
		t.Errorf("a wrapped map has room for the fields and was still refused: %v", err)
	}

	// And dropping the fields leaves the map free to be the whole response.
	dropped := strings.Replace(body, "  success:\n    fields:\n      ok: true\n", "", 1)

	if _, err := parse(t, dropped); err != nil {
		t.Errorf("a map with no fields to place was still refused: %v", err)
	}
}
