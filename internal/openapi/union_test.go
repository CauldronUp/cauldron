package openapi

import (
	"strings"
	"testing"
)

// One object, several shapes. A DNS record is an A record or a CNAME or an
// MX, and a description says so with oneOf: every variant declares name and
// ttl and content, and each adds its own.
//
// Properties merges allOf and deliberately does not merge oneOf, because
// merging alternatives into one object describes something that does not
// exist. But the question checkFields asks is narrower than that -- whether
// any schema declares this name at all -- and for that the union is the
// right answer. Cloudflare nests it three deep, allOf of an anyOf of two
// oneOfs of twenty-one variants between them, and every field of its DNS
// record was reported as undeclared while every one of the twenty-one
// declared it.

const unionSpec = `
openapi: 3.0.0
info: {title: Union, version: "1"}
paths:
  /v1/records:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  result:
                    type: array
                    items:
                      allOf:
                        - oneOf:
                            - $ref: "#/components/schemas/ARecord"
                            - $ref: "#/components/schemas/MXRecord"
                        - type: object
                          properties:
                            id: {type: string}
components:
  schemas:
    ARecord:
      type: object
      properties:
        name: {type: string}
        content: {type: string}
    MXRecord:
      type: object
      properties:
        name: {type: string}
        content: {type: string}
        priority: {type: integer}
`

const unionRecipe = `
recipe: union-shapes
capability: hosting
version: 0.1.0
upstream:
  api: "1"
responses:
  list:
    style: wrapped
resources:
  record:
    collection: result
    id:
      style: opaque
      length: 12
    fields:
      name:
        type: string
      content:
        type: string
      priority:
        type: integer
routes:
  - method: GET
    path: /v1/records
    resource: record
    operation: list
    returns:
      - id
      - name
      - content
      - priority
errors: {}
conformance:
  - name: a record has a name
    source: https://docs.union.test
    request:
      method: GET
      path: /v1/records
    expect:
      status: 200
      body:
        result: []
`

func TestAFieldDeclaredInOneVariantOfAUnionIsDeclared(t *testing.T) {
	for _, said := range disagreements(t, unionSpec, unionRecipe) {
		if strings.Contains(said, "field") {
			t.Errorf("every variant declares these and they were reported: %s", said)
		}
	}
}
