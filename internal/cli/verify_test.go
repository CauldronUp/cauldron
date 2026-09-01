package cli

import (
	"strings"
	"testing"
)

func TestVerifyRunsEveryBundledRecipe(t *testing.T) {
	stdout, stderr, code := run(t, "verify")

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s\n%s", code, stderr, stdout)
	}

	for _, want := range []string{"stripe", "github", "shopify", "twilio", "cases passed"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("report is missing %q\n%s", want, stdout)
		}
	}
}

// The provenance line is the honest part of the report. If it ever disappears,
// "verified" quietly starts meaning "we tested our fake against itself".
func TestVerifyStatesWhereTheEvidenceCameFrom(t *testing.T) {
	stdout, _, code := run(t, "verify", "stripe")

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, "documentation only") && !strings.Contains(stdout, "checked against the real API") {
		t.Errorf("the report must say what the evidence rests on\n%s", stdout)
	}
}

func TestVerifyListsEveryCaseWhenAsked(t *testing.T) {
	quiet, _, _ := run(t, "verify", "stripe")
	loud, _, _ := run(t, "verify", "stripe", "-v")

	if len(loud) <= len(quiet) {
		t.Errorf("-v should name the passing cases too\n%s", loud)
	}

	if !strings.Contains(loud, "ok") {
		t.Errorf("-v output = %q", loud)
	}
}

func TestVerifyRejectsAnUnknownRecipe(t *testing.T) {
	_, stderr, code := run(t, "verify", unshipped)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if stderr == "" {
		t.Error("an unknown recipe should explain itself")
	}
}

func TestProvenanceLineDistinguishesEvidence(t *testing.T) {
	cases := []struct {
		observed, documented int
		want                 string
	}{
		{0, 0, "no evidence"},
		{0, 5, "none checked against the real API"},
		{5, 0, "all 5 checked against the real API"},
		{3, 2, "3 checked against the real API, 2 from documentation only"},
	}

	for _, c := range cases {
		got := provenanceLine(c.observed, c.documented)

		if !strings.Contains(got, c.want) {
			t.Errorf("provenanceLine(%d, %d) = %q, want it to contain %q", c.observed, c.documented, got, c.want)
		}
	}
}
