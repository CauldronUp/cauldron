package recipe

import (
	"strings"
	"testing"
)

// returns names the fields a route keeps, and the identifier is not one of
// them unless it is asked for. Leaving it out is easy, silent, and produces a
// listing nothing can be fetched from: every entry arrives with no id, every
// case still passes, because no case asserted an identifier on a trimmed
// route.
//
// Three Recipes shipped that way. Asana's task listing handed back nothing
// but a name, Supabase's create stopped telling anyone the project's ref, and
// Discourse's topic lost the id it is addressed by.

const trimmedList = `
recipe: trimmed
capability: crm
version: 0.1.0
upstream:
  api: "1"
resources:
  contact:
    collection: data
    id:
      prefix: con_
      length: 14
    fields:
      name:
        type: string
      email:
        type: string
routes:
  - method: GET
    path: /v1/contacts
    resource: contact
    operation: list
    returns:
      - name
`

func TestAListingThatTrimsAwayItsIdentifierIsRejected(t *testing.T) {
	got := problems(t, trimmedList)

	if !strings.Contains(got, "returns does not name id") {
		t.Errorf("got %q", got)
	}
}

func TestNamingTheIdentifierSatisfiesIt(t *testing.T) {
	yaml := strings.Replace(trimmedList, "    returns:\n      - name", "    returns:\n      - id\n      - name", 1)

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("naming id should be enough: %v", err)
	}
}

func TestACreateMustSayWhatItMade(t *testing.T) {
	yaml := strings.Replace(trimmedList, "operation: list", "operation: create", 1)

	if got := problems(t, yaml); !strings.Contains(got, "returns does not name id") {
		t.Errorf("got %q", got)
	}
}

func TestAGetMayLeaveTheIdentifierOut(t *testing.T) {
	// Braze does exactly this: its details endpoint answers without the
	// identifier the caller just used to ask for it. The caller already has
	// it, so nothing is unreachable.
	yaml := strings.Replace(trimmedList, "operation: list", "operation: get", 1)
	yaml = strings.Replace(yaml, "path: /v1/contacts", "path: /v1/contacts/{id}", 1)

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("a get may answer without the identifier: %v", err)
	}
}

func TestAProviderThatSendsNoIdentifierIsNotHeldToIt(t *testing.T) {
	// PostHog's capture endpoint accepts anything and confirms nothing, and
	// Hightouch's rejected rows are the caller's own rows with no Hightouch
	// identifier at all. Saying so once, on the resource, is the way out --
	// not leaving it implied by the route's returns.
	yaml := strings.Replace(trimmedList, "      prefix: con_\n      length: 14", "      prefix: con_\n      length: 14\n      field: \"-\"", 1)

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("a resource with no identifier on the wire should pass: %v", err)
	}
}
