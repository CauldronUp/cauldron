package openapi

import (
	"strings"
	"testing"
)

// Two ways for check to report a Recipe that is right, both found by running
// it over Recipes whose emitted shape had already been read against the
// description by hand.
//
// The first is composition. A description is free to build a response out of
// allOf members rather than declaring properties on it, and Neon describes a
// branch exactly that way: the 200 has no properties of its own, only two
// members, one carrying the branch and one carrying an annotation. Document
// merges allOf when asked for a schema's properties, but the descent into a
// wrapped resource read the map directly and so walked into an object that
// looked empty, gave up, and compared every field against the envelope.
//
// The second is a field that never goes on the wire. "in: -" says the record
// holds it and the response does not carry it, which is how a partition that
// lives in the path is modelled -- Attio's object slug, Fly's app name. A
// field the client never sees cannot contradict a description of what the
// client sees, and reporting it puts a finding against the Recipes that took
// the trouble to say so.

const composedSpec = `
openapi: 3.0.0
info: {title: Composed, version: "1"}
paths:
  /v1/branches/{id}:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                allOf:
                  - $ref: "#/components/schemas/BranchResponse"
                  - $ref: "#/components/schemas/AnnotationResponse"
components:
  schemas:
    BranchResponse:
      type: object
      properties:
        branch: {$ref: "#/components/schemas/Branch"}
    AnnotationResponse:
      type: object
      properties:
        annotation: {type: object}
    Branch:
      type: object
      properties:
        name: {type: string}
        current_state: {type: string}
`

const composedRecipe = `
recipe: composed
capability: infrastructure
version: 0.1.0
upstream:
  api: "1"
responses:
  resource:
    style: wrapped
    key: branch
resources:
  branch:
    id:
      style: opaque
      field: name
    fields:
      current_state:
        type: string
      project_id:
        type: string
        in: "-"
routes:
  - method: GET
    path: /v1/branches/{id}
    resource: branch
    operation: get
errors: {}
conformance:
  - name: a branch has a state
    source: https://docs.composed.test
    request:
      method: GET
      path: /v1/branches/b1
    expect:
      status: 200
      body:
        branch.current_state: ready
`

func TestAResourceComposedWithAllOfIsFound(t *testing.T) {
	// The field is declared, one allOf member down and then one key in. A
	// descent that reads the properties map directly finds nothing at the
	// top, falls back to the envelope, and reports every field of the
	// resource -- which is what made a correct Neon Recipe look broken.
	for _, said := range disagreements(t, composedSpec, composedRecipe) {
		if strings.Contains(said, `field "current_state"`) {
			t.Errorf("current_state is declared inside an allOf member and was reported: %s", said)
		}
	}
}

func TestAFieldThatNeverGoesOnTheWireIsNotCompared(t *testing.T) {
	// "in: -" is a Recipe saying this field stays in the record. Holding a
	// description of the response to it asks the description to declare
	// something the response does not contain.
	for _, said := range disagreements(t, composedSpec, composedRecipe) {
		if strings.Contains(said, `field "project_id"`) {
			t.Errorf("project_id is declared as never sent and was reported: %s", said)
		}
	}
}
