package openapi

import (
	"strings"
	"testing"
)

// A single object can come back as a list of one, under the plural name.
// Xero does it for every resource: fetching one invoice answers with
// {"Invoices": [ ... ]}, and the runtime knows -- when a wrapped resource is
// an array, the key defaults to the collection name rather than the
// resource name, because a list of one is still a collection.
//
// check did not know, and looked for a key called "invoice". It never found
// one, fell back to comparing every field against the envelope -- whose
// properties are Invoices, Warnings and pagination -- and reported 61
// disagreements against a Recipe that was right.

const arrayWrapped = `
openapi: 3.0.0
info: {title: Wrapped, version: "1"}
paths:
  /v1/Invoices:
    post:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  Warnings: {type: array, items: {type: string}}
                  Invoices:
                    type: array
                    items: {$ref: "#/components/schemas/Invoice"}
components:
  schemas:
    Invoice:
      type: object
      properties:
        InvoiceID: {type: string}
        AmountDue: {type: number}
`

const arrayWrappedRecipe = `
recipe: wrapped-array
capability: accounting
version: 0.1.0
upstream:
  api: "1"
responses:
  resource:
    style: wrapped
    array: true
resources:
  invoice:
    collection: Invoices
    id:
      style: uuid
      field: InvoiceID
    fields:
      AmountDue:
        type: number
routes:
  - method: POST
    path: /v1/Invoices
    resource: invoice
    operation: create
    returns:
      - id
      - AmountDue
errors: {}
conformance:
  - name: an invoice has an amount due
    source: https://docs.wrapped.test
    request:
      method: POST
      path: /v1/Invoices
      json:
        AmountDue: 10
    expect:
      status: 200
      body:
        Invoices[0].AmountDue: 10
`

func TestAResourceWrappedAsAListOfOneUsesTheCollectionName(t *testing.T) {
	for _, said := range disagreements(t, arrayWrapped, arrayWrappedRecipe) {
		if strings.Contains(said, "AmountDue") {
			t.Errorf("AmountDue is declared inside the wrapped list and was reported: %s", said)
		}
	}
}
