package openapi

import (
	"strings"
	"testing"
)

// A route may answer with no body at all, and say so with empty_body. Ably
// publish does: it acknowledges with 204 and nothing else, which is why the
// Recipe declares it.
//
// Comparing the resource's fields against that route asks which of them the
// description declares, and the answer does not matter, because none of them
// is sent. Every field of the message was reported as undeclared on a route
// that emits nothing.

const emptyBodySpec = `
openapi: 3.0.0
info: {title: Ack, version: "1"}
paths:
  /v1/messages:
    post:
      responses:
        "2XX":
          description: published
          content:
            application/json:
              schema:
                type: object
                properties:
                  channel: {type: string}
                  messageId: {type: string}
`

const emptyBodyRecipe = `
recipe: ack-only
capability: chat
version: 0.1.0
upstream:
  api: "1"
resources:
  message:
    id:
      style: opaque
      length: 12
    fields:
      data:
        type: string
      clientId:
        type: string
routes:
  - method: POST
    path: /v1/messages
    resource: message
    operation: create
    status: 204
    empty_body: true
errors: {}
conformance:
  - name: publishing acknowledges and says nothing
    source: https://docs.ack.test
    request:
      method: POST
      path: /v1/messages
      json:
        data: hello
    expect:
      status: 204
      absent:
        - data
`

func TestARouteThatSendsNoBodyHasNoFieldsToCompare(t *testing.T) {
	for _, said := range disagreements(t, emptyBodySpec, emptyBodyRecipe) {
		if strings.Contains(said, "field") {
			t.Errorf("the route sends no body, so no field of it can disagree: %s", said)
		}
	}
}
