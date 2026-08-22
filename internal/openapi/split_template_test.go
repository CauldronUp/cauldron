package openapi

import (
	"strings"
	"testing"
)

// Two templates can differ only in what the parameter is called, and a
// description is free to split one path across both. Vercel does: the get on
// a deployment is declared under {idOrUrl} and the delete on the same path
// under {id}. Both match /v13/deployments/{id}, the first one wins, and the
// delete was reported as a method the description does not declare -- on a
// path where it is declared, one entry further down.

const splitTemplates = `
openapi: 3.0.0
info: {title: Split, version: "1"}
paths:
  /v1/things/{idOrUrl}:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema: {type: object, properties: {name: {type: string}}}
  /v1/things/{id}:
    delete:
      responses:
        "204": {description: gone}
`

const splitRecipe = `
recipe: split
capability: infrastructure
version: 0.1.0
upstream:
  api: "1"
resources:
  thing:
    id:
      style: opaque
      field: name
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/things/{id}
    resource: thing
    operation: get
  - method: DELETE
    path: /v1/things/{id}
    resource: thing
    operation: delete
errors: {}
conformance:
  - name: a thing has a name
    source: https://docs.split.test
    request:
      method: GET
      path: /v1/things/t1
    expect:
      status: 200
      body:
        name: t1
`

func TestAMethodDeclaredOnASecondTemplateIsFound(t *testing.T) {
	for _, said := range disagreements(t, splitTemplates, splitRecipe) {
		if strings.Contains(said, "DELETE") {
			t.Errorf("the delete is declared, on the other template: %s", said)
		}
	}
}
