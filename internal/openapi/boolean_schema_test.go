package openapi

import (
	"strings"
	"testing"
)

// JSON Schema allows a boolean where a schema is expected: false means
// nothing validates, true means anything does. Descriptions use it for
// additionalProperties constantly, and for items when a tuple is closed.
//
// Meilisearch's description has fourteen of them and could not be read at
// all -- "cannot unmarshal !!bool `false` into openapi.Schema", at a line
// number, with no hint that a perfectly ordinary construct was the cause.

const booleanSchemas = `
openapi: 3.0.0
info: {title: Strict, version: "1"}
paths:
  /v1/indexes:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                additionalProperties: false
                properties:
                  uid: {type: string}
                  swaps:
                    type: array
                    items: false
`

func TestABooleanWhereASchemaIsExpectedIsRead(t *testing.T) {
	doc, err := Parse([]byte(booleanSchemas))
	if err != nil {
		t.Fatalf("a description using boolean schemas was refused: %v", err)
	}

	op := doc.Paths["/v1/indexes"].Get
	if op == nil {
		t.Fatal("the operation did not survive")
	}

	schema, _ := doc.Success(op)
	if schema == nil {
		t.Fatal("no success schema")
	}

	props := doc.Properties(schema)
	if len(props) != 2 {
		t.Fatalf("read %d properties, want 2: %v", len(props), props)
	}
}

func TestAnHTMLPageSaysSoRatherThanFailingAsYAML(t *testing.T) {
	// Five providers this week served a documentation page from the URL
	// their docs give for the description. The YAML error for an HTML file
	// says "mapping values are not allowed in this context", which sends the
	// reader looking for a syntax error in a file that is not a description
	// at all.
	_, err := Parse([]byte("<!DOCTYPE html><html lang=\"en\"><body>docs</body></html>"))
	if err == nil {
		t.Fatal("expected a failure")
	}

	if !strings.Contains(err.Error(), "HTML") {
		t.Errorf("got %q", err)
	}
}
