package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/openapi"
)

func discoverOutput(t *testing.T, proposals []proposal, tally discoverTally) (string, int) {
	t.Helper()

	out := &bytes.Buffer{}
	code := reportDiscover(&context{stdout: out, stderr: out}, proposals, tally)

	return out.String(), code
}

// Finding nothing is the ordinary case -- most providers publish nothing --
// and a command that goes red for the ordinary case is one nobody runs twice.
func TestDiscoveryNeverFails(t *testing.T) {
	for _, c := range []struct {
		name      string
		proposals []proposal
		tally     discoverTally
	}{
		{"nothing found", nil, discoverTally{Searched: 40}},
		{"something found", []proposal{{
			Recipe: "notion",
			Found:  &openapi.Candidate{URL: "https://x/openapi.json", Matched: 6, Routes: 6, Declared: 45, Format: "OpenAPI 3.1.0"},
		}}, discoverTally{Searched: 1}},
	} {
		if _, code := discoverOutput(t, c.proposals, c.tally); code != 0 {
			t.Errorf("%s exited %d, want 0", c.name, code)
		}
	}
}

// A URL is printed for a person to paste, never written. Recording the wrong
// one does not fail loudly: it reports drift against a document that was never
// this provider's, every morning, until somebody works out why.
func TestAProposalIsPrintedAsALineToPaste(t *testing.T) {
	out, _ := discoverOutput(t, []proposal{{
		Recipe: "notion",
		Found: &openapi.Candidate{
			URL: "https://developers.notion.com/openapi.json", Matched: 6, Routes: 6, Declared: 45, Format: "OpenAPI 3.1.0",
		},
	}}, discoverTally{Searched: 1})

	if !strings.Contains(out, "spec: https://developers.notion.com/openapi.json") {
		t.Errorf("no pasteable spec line:\n%s", out)
	}

	if !strings.Contains(out, "6 of 6 route(s)") {
		t.Errorf("the evidence was not printed:\n%s", out)
	}

	if !strings.Contains(out, "Nothing was written.") {
		t.Errorf("the output does not say nothing was written:\n%s", out)
	}
}

// Monday's marketing domain answers a two-path openapi.json, one of which
// happens to be the single route its Recipe models. That is as consistent with
// coincidence as with discovery, and saying so is the difference between a
// proposal somebody checks and one somebody pastes.
func TestAProposalRestingOnAlmostNothingSaysSo(t *testing.T) {
	out, _ := discoverOutput(t, []proposal{{
		Recipe: "monday",
		Found:  &openapi.Candidate{URL: "https://monday.com/openapi.json", Matched: 1, Routes: 1, Declared: 2, Format: "OpenAPI 3.1.0"},
	}}, discoverTally{Searched: 1})

	if !strings.Contains(out, "weak") {
		t.Errorf("a one-route match in a two-path document was not called weak:\n%s", out)
	}
}

// The warning has to stay rare or it stops being read. A Recipe with one route
// matched against a real description is not a coincidence.
func TestAOneRouteRecipeAgainstARealDescriptionIsNotWeak(t *testing.T) {
	out, _ := discoverOutput(t, []proposal{{
		Recipe: "sunrisesunset",
		Found:  &openapi.Candidate{URL: "https://x/openapi.json", Matched: 1, Routes: 1, Declared: 80, Format: "OpenAPI 3.1.0"},
	}}, discoverTally{Searched: 1})

	if strings.Contains(out, "weak") {
		t.Errorf("a match inside an 80-path description was called weak:\n%s", out)
	}
}

// Strongest first, so the reader spends their attention where the evidence is.
func TestProposalsAreOrderedByTheShareOfTheRecipeAccountedFor(t *testing.T) {
	out, _ := discoverOutput(t, []proposal{
		{Recipe: "weakest", Found: &openapi.Candidate{URL: "u", Matched: 1, Routes: 10, Declared: 200}},
		{Recipe: "strongest", Found: &openapi.Candidate{URL: "u", Matched: 6, Routes: 6, Declared: 45}},
		{Recipe: "middling", Found: &openapi.Candidate{URL: "u", Matched: 3, Routes: 6, Declared: 60}},
	}, discoverTally{Searched: 3})

	order := []string{"strongest", "middling", "weakest"}

	var at int

	for _, name := range order {
		i := strings.Index(out, name)
		if i < 0 {
			t.Fatalf("%s missing from output:\n%s", name, out)
		}

		if i < at {
			t.Errorf("%s came out of order:\n%s", name, out)
		}

		at = i
	}
}

// The counts have to distinguish "searched and found nothing" from "could not
// be searched at all", because the second is a Recipe with no verified case to
// take a host from and is fixed by different work entirely.
func TestTheTallySeparatesNotFoundFromNotSearched(t *testing.T) {
	out, _ := discoverOutput(t, nil, discoverTally{Searched: 40, NoHost: 3, Recorded: 12})

	for _, want := range []string{
		"0 proposed",
		"40 searched and nothing found",
		"3 with no verified host",
		"12 already recorded",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the tally does not say %q:\n%s", want, out)
		}
	}
}
