package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A listing whose records carry a field named like one of the envelope's own
// keys. Val Town is the provider: each val has a `links` object with self and
// html in it, and each page of vals has a `links` object with next and prev.
const collidingSpec = `
openapi: 3.0.0
info: {title: Colliding, version: "1"}
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      type: object
                      properties:
                        id: {type: string}
                        name: {type: string}
                        links: {type: object}
                  links:
                    type: object
`

// The envelope check asks whether this parser got inside the response at all,
// and used to answer yes as soon as one field the Recipe names appeared at the
// top level of the success schema. A record field sharing a name with an
// envelope key satisfies that on its own, so the reading was declared to be
// inside an envelope it had not opened, and every other field on the record
// was reported as a gap the description had -- three sentences, none of them
// true, on a description that declares all three fields.
//
// Looking through the envelope to the record removes the collision instead of
// tuning the guess.
func TestARecordFieldNamedLikeAnEnvelopeKeyIsNotAGap(t *testing.T) {
	doc, err := Parse([]byte(collidingSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	r := driftRecipe()
	r.Routes = []recipe.Route{{Method: "GET", Path: "/things", Resource: "thing", Operation: "list"}}
	r.Resources["thing"] = recipe.Resource{Fields: map[string]recipe.Field{
		"id": {Type: "string"}, "name": {Type: "string"}, "links": {Type: "map"},
	}}

	for _, gap := range Unbacked(r, doc, "") {
		if strings.Contains(gap, "does not declare") {
			t.Errorf("every field is declared inside the records, and this calls one a gap: %s", gap)
		}
	}

	// And the check still works once a field really does go missing, so the
	// repair above is not simply silence.
	dropped := strings.Replace(collidingSpec, "                        name: {type: string}\n", "", 1)

	moved, err := Parse([]byte(dropped))
	if err != nil {
		t.Fatalf("parse dropped: %v", err)
	}

	var found bool

	for _, gap := range Unbacked(r, moved, "") {
		if strings.Contains(gap, "does not declare name") {
			found = true
		}
	}

	if !found {
		t.Error("name is gone from the record schema and no gap says so")
	}
}
