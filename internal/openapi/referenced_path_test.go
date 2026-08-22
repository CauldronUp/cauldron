package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A description may split itself across files, with every path a $ref to
// another document. Lob's does: all fifty-eight of its paths are one-line
// references to resources/<thing>/<thing>.yml, and none of those files is
// fetched here.
//
// The path matched, the path item was empty, and every route was reported as
// a method the description does not declare -- with an empty list of the
// methods it supposedly does declare: "the description declares  but not GET
// on it". That reads like a broken description rather than one this tool has
// only half read, and it is missing evidence rather than a disagreement.

const referencedPaths = `
openapi: 3.0.0
info: {title: Split, version: "1"}
paths:
  /v1/letters:
    $ref: "resources/letters/letters.yml"
  /v1/postcards:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema: {type: object, properties: {id: {type: string}}}
`

const referencedRecipe = `
recipe: split-files
capability: email
version: 0.1.0
upstream:
  api: "1"
resources:
  letter:
    id:
      prefix: ltr_
      length: 12
    fields:
      to:
        type: string
routes:
  - method: GET
    path: /v1/letters
    resource: letter
    operation: list
errors: {}
conformance:
  - name: a letter has a recipient
    source: https://docs.split.test
    request:
      method: GET
      path: /v1/letters
    expect:
      status: 200
      body:
        "[0]": {}
`

func TestAPathReferredToAnotherFileSaysSo(t *testing.T) {
	doc := parseSpec(t, referencedPaths)
	if doc.Paths["/v1/letters"].Ref == "" {
		t.Fatal("the path-level $ref was not read")
	}

	var said []string
	for _, f := range Check(mustRecipe(t, referencedRecipe), doc, "") {
		said = append(said, f.Severity+": "+f.Where+": "+f.What)
	}

	var found bool
	for _, s := range said {
		if strings.Contains(s, "another file") {
			found = true

			if !strings.HasPrefix(s, Missing) {
				t.Errorf("an unread file is missing evidence, not a disagreement: %s", s)
			}
		}

		if strings.Contains(s, "declares  but") {
			t.Errorf("an empty method list was reported: %s", s)
		}
	}

	if !found {
		t.Errorf("nothing said the path is declared elsewhere: %v", said)
	}
}

func mustRecipe(t *testing.T, source string) *recipe.Recipe {
	t.Helper()

	r, err := recipe.Parse([]byte(source))
	if err != nil {
		t.Fatalf("parsing the Recipe: %v", err)
	}

	return r
}
