// Two declarations, one wire key.

package recipe

import (
	"strings"
	"testing"
)

// A Recipe can name the key an identifier is emitted under, and it can rename
// a field onto a key with "as". Nothing stopped it doing both with the same
// key, and when it did the runtime picked one and dropped the other in
// silence -- so the response was byte-identical whichever way the Recipe was
// written, and no conformance case could tell them apart.
//
// The National Weather Service Recipe is what found it. A point is addressed
// at /points/39.7456,-97.0892 and its body carries the whole URL under "id",
// so the Recipe keys the record on the coordinate pair, hides the minted
// identifier with id.field "-", and renames a field onto "id". Un-hiding the
// minted identifier -- one character, and the kind of edit that happens by
// accident -- changed nothing in any of its nine responses.
func TestAFieldCannotTakeTheKeyTheIdentifierIsAlreadySentUnder(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  point:
    id:
      style: opaque
    fields:
      uri:
        type: string
        as: id
routes:
  - method: GET
    path: /points/{id}
    resource: point
    operation: get
`

	got := problems(t, yaml)

	for _, want := range []string{`field "uri"`, `sent as "id"`, "silently dropped"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected the collision to be reported; %q is missing from:\n%s", want, got)
		}
	}
}

// The same collision under a name the provider chose. Twilio calls the
// identifier "sid" everywhere, so a field renamed onto "sid" collides there
// exactly as one renamed onto "id" collides by default.
func TestTheCollisionIsFoundUnderAProvidersOwnIdentifierName(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  message:
    id:
      style: opaque
      field: sid
    fields:
      reference:
        type: string
        as: sid
routes:
  - method: GET
    path: /messages/{id}
    resource: message
    operation: get
`

	if got := problems(t, yaml); !strings.Contains(got, `sent as "sid"`) {
		t.Errorf("expected the collision under sid to be reported, got:\n%s", got)
	}
}

// And the shapes that are not a collision, because refusing any of them would
// turn away something a real provider sends.
//
// Two of these are held by the comparison itself rather than by a guard. A
// hidden identifier compares as the literal "-" and a nested one as "sys.id",
// and no top-level field is sent under either name -- so they fall out without
// a condition, and the conditions that were there for them were removed after
// this test showed that changing them changed no answer.
func TestTheIdentifierKeyRuleLeavesTheHonestShapesAlone(t *testing.T) {
	for name, yaml := range map[string]string{
		// The identifier is not emitted at all, which is the whole reason a
		// field is free to take the key.
		"hidden identifier": `
resources:
  point:
    id:
      style: opaque
      field: "-"
    fields:
      uri:
        type: string
        as: id
`,
		// The identifier nests, as Contentful's does at sys.id, so a
		// top-level "id" is a different key entirely.
		"nested identifier": `
resources:
  entry:
    id:
      style: opaque
      field: sys.id
    fields:
      uri:
        type: string
        as: id
`,
		// The field is simply called what the identifier is called, with no
		// "as" at all. That is how a Recipe declares the identifier's own type
		// and lines it up against an OpenAPI schema -- three fixtures in
		// internal/openapi do exactly this -- so it is a declaration about one
		// key rather than two claims on it.
		"field named after the identifier": `
resources:
  thing:
    id:
      style: opaque
      field: name
    fields:
      name:
        type: string
`,
		// The field nests under a parent, so it is not beside the identifier
		// either.
		"field nested under a parent": `
resources:
  point:
    id:
      style: opaque
    fields:
      uri:
        type: string
        as: id
        in: properties
`,
	} {
		t.Run(name, func(t *testing.T) {
			full := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
` + strings.TrimPrefix(yaml, "\n") + `routes:
  - method: GET
    path: /things/{id}
    resource: ` + resourceIn(yaml) + `
    operation: get
`

			r, err := parse(t, full)
			if err != nil && strings.Contains(err.Error(), "silently dropped") {
				t.Errorf("%s was refused as a key collision and is not one:\n%v", name, err)
			}

			_ = r
		})
	}
}

// resourceIn reads the single resource name out of the fragments above, so the
// route in each case addresses the resource that case declares.
func resourceIn(yaml string) string {
	for _, line := range strings.Split(yaml, "\n") {
		if strings.HasPrefix(line, "  ") && strings.HasSuffix(line, ":") &&
			!strings.HasPrefix(line, "   ") {
			return strings.TrimSuffix(strings.TrimSpace(line), ":")
		}
	}

	return ""
}
