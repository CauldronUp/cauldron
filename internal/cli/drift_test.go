package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/openapi"
)

func driftOutput(t *testing.T, reports []openapi.DriftReport, record, quiet bool) (string, int) {
	t.Helper()

	out := &bytes.Buffer{}
	code := reportDrift(&context{stdout: out, stderr: out}, reports, record, quiet)

	return out.String(), code
}

// The property the scheduled scan rests on. If anything but a moved
// fingerprint can fail the build, the first flaky docs host turns the job red,
// and a job that is red for reasons nobody can act on gets muted -- taking the
// one signal that mattered with it.
func TestOnlyAMovedFingerprintFailsTheScan(t *testing.T) {
	for _, report := range []openapi.DriftReport{
		{Recipe: "a", Status: openapi.Undeclared},
		{Recipe: "b", Status: openapi.Unchanged, Seen: "2026-08-30"},
		{Recipe: "c", Status: openapi.Unrecorded, Now: "abc"},
		{Recipe: "d", Status: openapi.Unreachable, Err: errors.New("503 Service Unavailable")},
	} {
		if _, code := driftOutput(t, []openapi.DriftReport{report}, false, false); code != 0 {
			t.Errorf("%q exited %d, and only a moved fingerprint may fail", report.Status, code)
		}
	}

	moved := []openapi.DriftReport{{Recipe: "e", Status: openapi.Moved, Seen: "2026-08-30", Recorded: "a", Now: "b"}}

	if _, code := driftOutput(t, moved, false, false); code != 1 {
		t.Errorf("a moved fingerprint exited %d, want 1", code)
	}
}

// An unreachable description must still be printed. Not failing the build is
// not the same as saying nothing: a host that has been unreachable for a month
// is a Recipe nobody is checking, and the only way to notice is to see it in
// the output every time.
func TestAnUnreachableDescriptionIsPrintedWithItsReason(t *testing.T) {
	out, _ := driftOutput(t, []openapi.DriftReport{
		{Recipe: "somewhere", Status: openapi.Unreachable, Err: errors.New("503 Service Unavailable")},
	}, false, false)

	if !strings.Contains(out, "somewhere") || !strings.Contains(out, "503 Service Unavailable") {
		t.Errorf("the reason was not reported:\n%s", out)
	}
}

// The count of Recipes nothing can check has to be visible, or a green scan
// reads as a clean bill of health for the whole collection when most of it was
// never examined.
func TestTheScanSaysHowManyRecipesItCouldNotExamine(t *testing.T) {
	out, _ := driftOutput(t, []openapi.DriftReport{
		{Recipe: "a", Status: openapi.Undeclared},
		{Recipe: "b", Status: openapi.Undeclared},
		{Recipe: "c", Status: openapi.Unchanged, Seen: "2026-08-30"},
	}, false, true)

	if !strings.Contains(out, "2 with no description to check") {
		t.Errorf("the unexamined Recipes were not counted:\n%s", out)
	}

	if !strings.Contains(out, "It is unexamined.") {
		t.Errorf("the scan did not say what a missing description means:\n%s", out)
	}
}

// Quiet is for the scheduled job, which should print the Recipes somebody has
// to act on and not three hundred lines of "no description to check".
func TestQuietPrintsOnlyWhatSomebodyMustAct(t *testing.T) {
	out, _ := driftOutput(t, []openapi.DriftReport{
		{Recipe: "undeclared-one", Status: openapi.Undeclared},
		{Recipe: "unchanged-one", Status: openapi.Unchanged, Seen: "2026-08-30"},
		{Recipe: "moved-one", Status: openapi.Moved, Seen: "2026-08-30"},
	}, false, true)

	if strings.Contains(out, "undeclared-one") || strings.Contains(out, "unchanged-one") {
		t.Errorf("quiet printed a Recipe nobody has to act on:\n%s", out)
	}

	if !strings.Contains(out, "moved-one") {
		t.Errorf("quiet dropped the one Recipe that moved:\n%s", out)
	}
}

// Recording prints the lines to paste rather than editing the Recipe, because
// a Recipe's comments are most of what it is worth and a scan that rewrites
// the evidence it was asked to check is the wrong shape of tool.
func TestRecordingPrintsTheLinesRatherThanWritingThem(t *testing.T) {
	out, _ := driftOutput(t, []openapi.DriftReport{
		{Recipe: "moved-one", Status: openapi.Moved, Seen: "2026-08-30", Recorded: "old", Now: "newfingerprint"},
	}, true, false)

	if !strings.Contains(out, "spec_hash: newfingerprint") {
		t.Errorf("the fingerprint to record was not printed:\n%s", out)
	}

	if !strings.Contains(out, "spec_seen:") {
		t.Errorf("the date to record was not printed:\n%s", out)
	}
}

// An unchanged Recipe has nothing to record, and printing a line to paste
// beside it invites somebody to paste it and call that a check.
func TestRecordingSaysNothingAboutAnUnchangedRecipe(t *testing.T) {
	out, _ := driftOutput(t, []openapi.DriftReport{
		{Recipe: "steady", Status: openapi.Unchanged, Seen: "2026-08-30", Recorded: "same", Now: "same"},
	}, true, false)

	if strings.Contains(out, "spec_hash:") {
		t.Errorf("a line to paste was offered for an unchanged Recipe:\n%s", out)
	}
}
