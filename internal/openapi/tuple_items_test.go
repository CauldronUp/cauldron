package openapi

import "testing"

// OpenAPI 3.0 says items holds one schema. Draft-4 JSON Schema says it may
// hold a list, one schema per position, and descriptions written against the
// older draft carry that through.
//
// Webflow's does, thirty-two times, and the whole 4.4MB file was refused with
// "cannot unmarshal !!seq into openapi.plain" repeated down the screen. The
// construct has nothing to do with what this package reads a description for.
const tupleItems = `
openapi: 3.0.0
info: {title: Tuple, version: "1"}
paths:
  /v2/sites:
    get:
      parameters:
        - name: limit
          in: query
        - name: offset
          in: query
      responses:
        "200":
          content:
            application/json:
              schema:
                type: object
                properties:
                  details:
                    type: array
                    items:
                      - type: string
                      - type: object
`

func TestTupleFormItemsDoesNotRefuseTheFile(t *testing.T) {
	doc, err := Parse([]byte(tupleItems))
	if err != nil {
		t.Fatalf("a description using tuple-form items was refused: %v", err)
	}

	if _, ok := doc.Paths["/v2/sites"]; !ok {
		t.Fatalf("the path did not survive: %v", doc.Paths)
	}
}

// The first entry is what gets read. A tuple whose positions disagree has no
// single answer to "what shape is an element", and reading the first is closer
// than reading nothing.
func TestTupleFormItemsReadsTheFirstPosition(t *testing.T) {
	doc, err := Parse([]byte(tupleItems))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	schema := doc.Paths["/v2/sites"].Get.Responses["200"].Content["application/json"].Schema
	details := schema.Properties["details"]

	if details == nil || details.Items == nil {
		t.Fatal("the array's items were dropped entirely")
	}

	if details.Items.Type != "string" {
		t.Errorf("read the items as %q, want the first position's string", details.Items.Type)
	}
}
