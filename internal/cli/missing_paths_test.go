package cli

import (
	"testing"

	"github.com/CauldronUp/cauldron/internal/openapi"
)

// Every route missing is not a Recipe full of invented paths. It is the
// signature of the wrong document, and three of the thirty-four descriptions
// already fetched turned out to be one: CircleCI's declares /api/v1 while the
// Recipe models v2, and Pipedrive's and Snyk's are the same mistake.
//
// Without the notice those arrive as a list of findings indistinguishable
// from real ones, and somebody works through them.
func TestMissingPathsAreCounted(t *testing.T) {
	findings := []openapi.Finding{
		{What: "the description declares no such path"},
		{What: "answers 202 and the description declares 200"},
		{What: "the description declares no such path"},
	}

	if got := countMissingPaths(findings); got != 2 {
		t.Errorf("counted %d, want 2", got)
	}
}

func TestOtherDisagreementsAreNotCounted(t *testing.T) {
	findings := []openapi.Finding{
		{What: "answers 429 and the description declares no 429 on any operation"},
		{What: `is sent as "clientId" and no property of that name is declared`},
	}

	if got := countMissingPaths(findings); got != 0 {
		t.Errorf("counted %d, want 0: neither of those is a missing path", got)
	}
}

// A path declared in another file is not a path the description does not
// have, and the two are reported differently on purpose.
func TestAPathInAnotherFileIsNotMissing(t *testing.T) {
	findings := []openapi.Finding{{What: "the path is declared in another file"}}

	if got := countMissingPaths(findings); got != 0 {
		t.Errorf("counted %d, want 0", got)
	}
}
