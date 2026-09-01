package openapi

import (
	"errors"
	"strings"
	"testing"
)

// A self-referential YAML anchor is refused, not crashed on.
//
// Windmill publishes a 1.5MB description with 720 anchors and 2784 aliases, and
// one of them contains itself. Decoding that into this package's own types
// walked the node graph through Schema's custom unmarshaler and recursed until
// the goroutine stack passed a gigabyte and the process died -- no error, no
// Recipe name, nothing in the output to say which description had done it. The
// whole scan simply stopped.
//
// The document below is the same shape in eleven lines. It has to be refused
// before the typed decode, because the typed decode is what recurses.
func TestASelfReferentialAnchorIsRefusedRatherThanCrashedOn(t *testing.T) {
	const document = `
openapi: 3.0.0
info:
  title: Recursive
  version: "1"
components:
  schemas:
    Loop: &loop
      type: object
      properties:
        self: *loop
paths: {}
`

	_, err := Parse([]byte(document))
	if err == nil {
		t.Fatal("a description whose anchor contains itself parsed")
	}

	// The column it lands in matters as much as the fact that it fails. A
	// document that will never be readable is not a host having a bad
	// afternoon, and reporting it as unreachable would put it under "try
	// again tomorrow" for as long as the provider publishes it.
	var format *FormatError
	if !errors.As(err, &format) {
		t.Fatalf("refused as %T, want a FormatError so drift files it under the format it cannot read", err)
	}

	// The anchor's own name has to survive into the message. Without it the
	// reader of a failing scan has 1.5MB and no starting point.
	if !strings.Contains(err.Error(), "loop") {
		t.Errorf("the failure does not name the anchor at fault: %s", err)
	}
}

// An ordinary description with anchors and aliases still parses.
//
// The guard is a cycle check, not an anchor check. Plenty of descriptions use
// anchors to avoid repeating themselves and there is nothing wrong with that;
// refusing them would trade a crash for a false negative.
func TestAnchorsWithoutACycleStillParse(t *testing.T) {
	const document = `
openapi: 3.0.0
info:
  title: Shared
  version: "1"
components:
  schemas:
    Timestamps: &timestamps
      type: object
      properties:
        created_at:
          type: string
    Order:
      allOf:
        - *timestamps
paths:
  /orders:
    get:
      responses:
        "200":
          description: ok
`

	doc, err := Parse([]byte(document))
	if err != nil {
		t.Fatalf("a description sharing a schema by anchor was refused: %v", err)
	}

	if doc.OpenAPI != "3.0.0" {
		t.Errorf("parsed version %q, want 3.0.0", doc.OpenAPI)
	}
}
