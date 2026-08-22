package recipe

import (
	"strings"
	"testing"
)

// A required header and a filter parameter are both claims about a name, and
// both were checkable only by whatever cases happened to exist.
//
// Enforcing Notion-Version is one of the reasons this project exists:
// forgetting it is the classic Notion integration bug, and a fake that waves
// it through lets code ship that fails on its first real call. A Recipe could
// declare the header, never enforce it, and stay green -- every case sends it.
//
// All eighteen required headers and all three filters shipping today already
// pass. These rules exist so the next one has to.

const claimingRecipe = `
recipe: claiming
capability: crm
version: 0.1.0
upstream:
  api: "1"
required_headers:
  X-Api-Version: missing_version
errors:
  missing_version:
    status: 400
    message: a version header is required
resources:
  contact:
    collection: data
    id:
      prefix: con_
      length: 14
    fields:
      status:
        type: string
routes:
  - method: GET
    path: /v1/contacts
    resource: contact
    operation: list
    returns:
      - id
      - status
    filters:
      - param: status
        field: status
conformance:
  - name: a listing narrows itself when asked
    source: https://example.invalid/docs
    fixture: small
    request:
      method: GET
      path: /v1/contacts?status=active
      headers:
        X-Api-Version: "1"
    expect:
      status: 200
      body:
        data[0].status: active
  - name: a request without the version header is refused
    source: https://example.invalid/docs
    fixture: small
    request:
      method: GET
      path: /v1/contacts
    expect:
      status: 400
      body:
        message: a version header is required
fixtures:
  small:
    contact:
      - id: con_00000000000001
        status: active
`

func TestARequiredHeaderNoCaseOmitsIsRejected(t *testing.T) {
	// Give the omission case the header, so every case sends it and nothing
	// shows the header is enforced.
	yaml := strings.Replace(claimingRecipe,
		"      path: /v1/contacts\n    expect:\n      status: 400",
		"      path: /v1/contacts\n      headers:\n        X-Api-Version: \"1\"\n    expect:\n      status: 400", 1)

	if got := problems(t, yaml); !strings.Contains(got, "no conformance case omits it and is refused") {
		t.Errorf("got %q", got)
	}
}

// Omitting the header and succeeding proves the opposite of enforcement, so
// the case only counts when the request is refused.
func TestOmittingTheHeaderAndSucceedingDoesNotCount(t *testing.T) {
	yaml := strings.Replace(claimingRecipe, "    expect:\n      status: 400", "    expect:\n      status: 200", 1)

	if got := problems(t, yaml); !strings.Contains(got, "no conformance case omits it and is refused") {
		t.Errorf("got %q", got)
	}
}

func TestAFilterNoCaseSendsIsRejected(t *testing.T) {
	yaml := strings.Replace(claimingRecipe, "      - param: status", "      - param: state", 1)

	if got := problems(t, yaml); !strings.Contains(got, `declares a filter on "state"`) {
		t.Errorf("got %q", got)
	}
}

// A case may carry its query in the map or written into the path, and both
// are used: Alpaca requests "/v2/orders?status=closed" directly. A check
// reading only the map reports every one of those as unexercised, which is
// how this was first got wrong.
func TestAFilterSentInTheQueryMapAlsoCounts(t *testing.T) {
	yaml := strings.Replace(claimingRecipe,
		"      path: /v1/contacts?status=active",
		"      path: /v1/contacts\n      query:\n        status: active", 1)

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("a filter sent through the query map should count: %v", err)
	}
}

func TestTheRecipeAsWrittenIsAccepted(t *testing.T) {
	if _, err := parse(t, claimingRecipe); err != nil {
		t.Errorf("the unmutated Recipe should pass: %v", err)
	}
}
