// Telling a file this could not read from a description that read fine and is
// not a valid one.

package openapi

import (
	"errors"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Codecov found this. api.codecov.io/api/v2/schema/ is well-formed YAML that
// declares, on the repositories listing, a parameter whose schema reads
// `type: array` with `items: string` -- where OpenAPI requires items to be a
// schema object. One node in 34 paths, on the endpoint a client is most likely
// to generate first, and it makes the whole document unreadable to a strict
// parser.
//
// The document arrives intact every time and no amount of waiting fixes it. So
// it belongs in the column that means "this will be the same tomorrow", not
// the one that means "the docs host had a bad afternoon" -- otherwise a
// scheduled scan shows a permanent fact in the column people learn to skip.
const codecovShapedSpec = `
openapi: 3.0.3
info:
  title: Malformed API
  version: 2.0.0
paths:
  /things/:
    get:
      parameters:
      - in: query
        name: names
        schema:
          type: array
          items: string
      responses:
        '200':
          description: ok
`

func TestADocumentThatParsesAndIsNotADescriptionIsAFormatError(t *testing.T) {
	t.Parallel()

	_, err := Parse([]byte(codecovShapedSpec))
	if err == nil {
		t.Fatal("a schema with items: string parsed, and OpenAPI requires an object there")
	}

	var format *FormatError
	if !errors.As(err, &format) {
		t.Fatalf("the error is %v, want a FormatError so drift reports it as unsupported rather than unreachable", err)
	}

	// The message has to say which of the two it is, because the remedy
	// differs: one is waiting, the other is telling the provider.
	if !strings.Contains(err.Error(), "not a valid OpenAPI description") {
		t.Errorf("the message does not say the document is the problem: %v", err)
	}
}

// And the other half of the distinction: bytes that are not a document at all
// stay unreachable, because an HTML error page or a truncated download says
// nothing permanent about the provider.
func TestBytesThatAreNotADocumentAtAllStayUnreadable(t *testing.T) {
	t.Parallel()

	// A tab where YAML requires a space, which no YAML parser accepts and
	// which is what a corrupted or wrongly-typed download tends to look like.
	_, err := Parse([]byte("openapi: 3.0.3\npaths:\n\t- broken\n  \tmore: [unclosed\n"))
	if err == nil {
		t.Fatal("unparseable bytes were accepted as a description")
	}

	var format *FormatError
	if errors.As(err, &format) {
		t.Errorf("a file that is not YAML was reported as a provider's format problem: %v", err)
	}
}

// The classification is what drift reads, so check it there rather than only
// at the parser.
func TestDriftPutsAMalformedDescriptionInThePermanentColumn(t *testing.T) {
	t.Parallel()

	spec := &recipe.Recipe{
		Name:     "malformed",
		Upstream: recipe.Upstream{Spec: "https://example.test/schema/"},
	}

	reports := Drift([]*recipe.Recipe{spec}, func(string) ([]byte, error) {
		return []byte(codecovShapedSpec), nil
	})

	if len(reports) != 1 {
		t.Fatalf("got %d reports, want 1", len(reports))
	}

	if got := reports[0].Status; got != Unsupported {
		t.Errorf("a malformed description is reported %q, want %q -- %q is the column that means try again tomorrow", got, Unsupported, Unreachable)
	}

	// And it is not drift, because nothing has been compared.
	if reports[0].Moved() {
		t.Error("a description that could not be read was reported as moved")
	}
}
