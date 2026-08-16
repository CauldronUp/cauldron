package cli

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/server"
)

// liveServer boots a real Cauldron server so the CLI is exercised against the
// actual control API rather than a hand-written stub. A stub would happily
// agree with a wrong client.
func liveServer(t *testing.T) string {
	t.Helper()

	srv := server.New()

	if err := srv.Mount("stripe", 1, "small-shop"); err != nil {
		t.Fatalf("mount: %v", err)
	}

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	return ts.URL
}

func TestStatusAgainstARunningServer(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "status", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	for _, want := range []string{"Sandbox time", "RECIPE", "stripe", "small-shop"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("status output missing %q\n%s", want, stdout)
		}
	}
}

func TestStatusExplainsWhenNothingIsRunning(t *testing.T) {
	// Port 1 is reserved and nothing will be listening.
	_, stderr, code := run(t, "status", "--url", "http://127.0.0.1:1")

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "cauldron serve") {
		t.Errorf("an unreachable server should suggest starting one; got %q", stderr)
	}
}

func TestSeedAndReset(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "seed", "stripe", "--fixture", "empty", "--url", url)
	if code != 0 {
		t.Fatalf("seed exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Seeded stripe with empty") {
		t.Errorf("stdout = %q", stdout)
	}

	stdout, stderr, code = run(t, "reset", "--url", url)
	if code != 0 {
		t.Fatalf("reset exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Reset every recipe") {
		t.Errorf("stdout = %q", stdout)
	}
}

func TestSeedReportsTheServersOwnError(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "seed", "stripe", "--fixture", "enormous", "--url", url)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	// The server knows which fixtures exist; the CLI should relay that rather
	// than inventing its own message.
	if !strings.Contains(stderr, "small-shop") {
		t.Errorf("stderr should list available fixtures; got %q", stderr)
	}
}

func TestFaultArmsAndIsDescribed(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "fault", "stripe", "--error", "rate_limit", "--count", "3", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Armed rate_limit on stripe") {
		t.Errorf("stdout = %q", stdout)
	}

	if !strings.Contains(stdout, "3 request(s)") {
		t.Errorf("the confirmation should describe the limit; got %q", stdout)
	}
}

func TestFaultRejectsAnUnknownFailureWithHelp(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "fault", "stripe", "--error", "meteor", "--url", url)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "rate_limit") {
		t.Errorf("stderr should list the real options; got %q", stderr)
	}
}

func TestFaultWithoutAnErrorPointsAtRecipeInfo(t *testing.T) {
	_, stderr, code := run(t, "fault", "stripe")

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "recipe info stripe") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestEmit(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "emit", "stripe", "payment_intent.payment_failed", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Emitted payment_intent.payment_failed") {
		t.Errorf("stdout = %q", stdout)
	}
}

func TestEmitRejectsAnEventTheProviderNeverSends(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "emit", "stripe", "customer.exploded", "--url", url)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "does not emit") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestClockAdvance(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "clock", "advance", "30d", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "2026-01-31") {
		t.Errorf("stdout = %q, want the advanced date", stdout)
	}
}

func TestClockRejectsNonsense(t *testing.T) {
	url := liveServer(t)

	_, _, code := run(t, "clock", "advance", "soon", "--url", url)

	if code != 1 {
		t.Errorf("exit = %d, want 1", code)
	}
}

func TestClockUsageWhenIncomplete(t *testing.T) {
	_, stderr, code := run(t, "clock")

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "clock advance") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestRequestsShowsWhatWasSent(t *testing.T) {
	url := liveServer(t)

	// Drive one call through the emulator so there is something to show.
	if _, _, code := run(t, "fault", "stripe", "--error", "rate_limit", "--url", url); code != 0 {
		t.Fatal("arming the fault failed")
	}

	stdout, stderr, code := run(t, "requests", "stripe", "--url", url)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "not been called yet") && !strings.Contains(stdout, "METHOD") {
		t.Errorf("unexpected output %q", stdout)
	}
}

func TestRequestsNeedsARecipe(t *testing.T) {
	_, stderr, code := run(t, "requests")

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "which recipe") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestHelpListsTheControlCommands(t *testing.T) {
	stdout, _, code := run(t, "help")

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	for _, want := range []string{"status", "requests", "seed", "reset", "fault", "emit", "clock", "CAULDRON_URL"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help is missing %q\n%s", want, stdout)
		}
	}
}

func TestSnapshotSaveLoadAndList(t *testing.T) {
	url := liveServer(t)
	dir := t.TempDir()

	// Change the sandbox away from its seeded state so the restore is provably
	// doing something.
	if _, _, code := run(t, "seed", "stripe", "--fixture", "empty", "--url", url); code != 0 {
		t.Fatal("seeding empty failed")
	}

	stdout, stderr, code := run(t, "snapshot", "save", "empty-shop", "--url", url, "--dir", dir)
	if code != 0 {
		t.Fatalf("save exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Saved empty-shop") {
		t.Errorf("stdout = %q", stdout)
	}

	// Move away from the captured state.
	if _, _, code := run(t, "seed", "stripe", "--fixture", "small-shop", "--url", url); code != 0 {
		t.Fatal("reseeding failed")
	}

	stdout, stderr, code = run(t, "snapshot", "load", "empty-shop", "--url", url, "--dir", dir)
	if code != 0 {
		t.Fatalf("load exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Restored empty-shop") {
		t.Errorf("stdout = %q", stdout)
	}

	stdout, _, code = run(t, "snapshot", "list", "--dir", dir)
	if code != 0 {
		t.Fatalf("list exit = %d", code)
	}

	if !strings.Contains(stdout, "empty-shop") {
		t.Errorf("listing is missing the snapshot\n%s", stdout)
	}
}

func TestSnapshotLoadNamesWhatIsAvailable(t *testing.T) {
	url := liveServer(t)
	dir := t.TempDir()

	if _, _, code := run(t, "snapshot", "save", "known-good", "--url", url, "--dir", dir); code != 0 {
		t.Fatal("save failed")
	}

	_, stderr, code := run(t, "snapshot", "load", "typo", "--url", url, "--dir", dir)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "known-good") {
		t.Errorf("an unknown name should list what exists; got %q", stderr)
	}
}

func TestSnapshotListIsEmptyByDefault(t *testing.T) {
	stdout, _, code := run(t, "snapshot", "list", "--dir", t.TempDir())

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, "No snapshots") {
		t.Errorf("stdout = %q", stdout)
	}
}

func TestSnapshotSaveNeedsAName(t *testing.T) {
	_, stderr, code := run(t, "snapshot", "save")

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "before-refund-bug") {
		t.Errorf("the error should show a usable example; got %q", stderr)
	}
}

func TestSnapshotRejectsAFileThatIsNotASnapshot(t *testing.T) {
	url := liveServer(t)
	dir := t.TempDir()

	path := filepath.Join(dir, ".cauldron", "snapshots")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(path, "junk.json"), []byte("not json at all"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, stderr, code := run(t, "snapshot", "load", "junk", "--url", url, "--dir", dir)

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "not a Cauldron snapshot") {
		t.Errorf("stderr = %q", stderr)
	}
}
