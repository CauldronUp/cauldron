package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A description that wraps its payload is the common case, not the exception.
// Cloudflare wraps under result, Xero wraps under a plural and puts one object
// inside an array, and Qdrant wraps under result and then again under a name
// that differs per endpoint.
//
// Only listings were unwrapped, so every one of those Recipes had every field
// of every resource reported as undeclared: the comparison ran against the
// envelope, whose properties are things like result and status, and no
// resource has fields called those. A report that cries wolf is one nobody
// reads the second time.

const wrappedSpec = `
openapi: 3.0.0
info: {title: Wrapped, version: "1"}
paths:
  /v1/things/{id}:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  status: {type: string}
                  result: {$ref: "#/components/schemas/Thing"}
        "4XX":
          description: error
          content:
            application/json:
              schema: {type: object}
  /v1/things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  status: {type: string}
                  result:
                    type: object
                    properties:
                      things:
                        type: array
                        items: {$ref: "#/components/schemas/Summary"}
        "4XX":
          description: error
          content:
            application/json:
              schema: {type: object}
components:
  schemas:
    Thing:
      type: object
      properties:
        name: {type: string}
        colour: {type: string}
    Summary:
      type: object
      properties:
        name: {type: string}
`

const wrappedRecipe = `
recipe: wrapped
capability: commerce
version: 0.1.0
upstream:
  api: "1"
responses:
  resource:
    style: wrapped
    key: result
resources:
  thing:
    collection: result.things
    id:
      style: opaque
      field: name
    fields:
      colour:
        type: string
routes:
  - method: GET
    path: /v1/things/{id}
    resource: thing
    operation: get
  - method: GET
    path: /v1/things
    resource: thing
    operation: list
    returns:
      - id
errors:
  nope:
    status: 404
    message: "no"
conformance:
  - name: a thing has a colour
    source: https://docs.wrapped.test
    request:
      method: GET
      path: /v1/things/t1
    expect:
      status: 200
      body:
        result.colour: green
`

func disagreements(t *testing.T, spec, source string) []string {
	t.Helper()

	doc := parseSpec(t, spec)

	r, err := recipe.Parse([]byte(source))
	if err != nil {
		t.Fatalf("parsing the Recipe: %v", err)
	}

	var said []string

	for _, finding := range Check(r, doc, "") {
		if finding.Severity == Disagrees {
			said = append(said, finding.Where+": "+finding.What)
		}
	}

	return said
}

func TestAWrappedResourceIsComparedAgainstWhatItWraps(t *testing.T) {
	for _, said := range disagreements(t, wrappedSpec, wrappedRecipe) {
		if strings.Contains(said, `field "colour"`) {
			t.Errorf("colour is declared one level down and was reported: %s", said)
		}
	}
}

func TestARouteThatReturnsLessIsHeldToLess(t *testing.T) {
	// The listing sends a name and nothing else, which is the point of the
	// listing. Holding the whole record to that schema reports the very
	// difference the declaration exists to describe.
	for _, said := range disagreements(t, wrappedSpec, wrappedRecipe) {
		if strings.Contains(said, "/v1/things") && strings.Contains(said, "colour") {
			t.Errorf("the listing trims colour and was reported anyway: %s", said)
		}
	}
}

func TestADeclaredRangeCoversItsBlock(t *testing.T) {
	// A description saying 4XX has said a 404 is expected. Reading only the
	// numeric codes beside it reported every error entry as undeclared.
	for _, said := range disagreements(t, wrappedSpec, wrappedRecipe) {
		if strings.Contains(said, "nope") {
			t.Errorf("404 is covered by the declared 4XX and was reported: %s", said)
		}
	}
}

func TestDefaultIsNotARange(t *testing.T) {
	// "default" means every status the operation did not name, which is a
	// description that has stopped enumerating rather than one saying a
	// particular failure is expected.
	if block, ok := statusBlock("default"); ok {
		t.Errorf("default read as the %dxx block", block)
	}

	for _, bad := range []string{"", "20", "2000", "2X", "XXX", "0XX", "6XX", "200"} {
		if _, ok := statusBlock(bad); ok {
			t.Errorf("%q read as a range", bad)
		}
	}

	for _, good := range []string{"4XX", "5xx", "2Xx"} {
		if _, ok := statusBlock(good); !ok {
			t.Errorf("%q is a range and was not read as one", good)
		}
	}
}
