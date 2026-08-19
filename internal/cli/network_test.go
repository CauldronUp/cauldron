package cli

import (
	"strings"
	"testing"
)

func TestNetworkArmsAgainstARunningServer(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "network", "stripe", "--latency", "800ms", "--jitter", "200ms", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	for _, want := range []string{"stripe", "latency 800ms", "200ms"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q:\n%s", want, stdout)
		}
	}
}

func TestNetworkAcceptsEveryFlag(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "network", "stripe",
		"--bandwidth", "50",
		"--timeout", "5s",
		"--limit", "1024",
		"--slice", "64",
		"--probability", "0.5",
		"--count", "3",
		"--for", "30s",
		"--path", "/v1/charges",
		"--url", url,
	)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}
}

func TestNetworkReset(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "network", "stripe", "--reset", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "reset") {
		t.Errorf("stdout should confirm the reset:\n%s", stdout)
	}
}

func TestNetworkClear(t *testing.T) {
	url := liveServer(t)

	if _, stderr, code := run(t, "network", "stripe", "--latency", "1s", "--url", url); code != 0 {
		t.Fatalf("arming failed: %s", stderr)
	}

	stdout, stderr, code := run(t, "network", "stripe", "--clear", "--url", url)

	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, "Cleared") {
		t.Errorf("stdout should confirm the clear:\n%s", stdout)
	}
}

func TestNetworkWithoutARecipeSaysWhatToType(t *testing.T) {
	stdout, stderr, code := run(t, "network", "--latency", "1s")

	if code == 0 {
		t.Fatalf("expected a non-zero exit, stdout = %s", stdout)
	}

	// A usage error that does not show the working form makes the reader go
	// and find the docs, which is the whole cost this avoids.
	if !strings.Contains(stderr, "cauldron network stripe") {
		t.Errorf("stderr should show a working example:\n%s", stderr)
	}
}

// Arming nothing is almost always a mistyped flag, so it has to be refused
// loudly rather than accepted as a no-op.
func TestNetworkWithNoConditionsIsRefused(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "network", "stripe", "--url", url)

	if code == 0 {
		t.Fatal("arming no conditions should fail")
	}

	if !strings.Contains(stderr, "nothing to degrade") {
		t.Errorf("stderr should explain what is missing:\n%s", stderr)
	}
}

func TestNetworkAgainstAServerThatIsNotRunning(t *testing.T) {
	_, stderr, code := run(t, "network", "stripe", "--latency", "1s", "--url", "http://127.0.0.1:1")

	if code == 0 {
		t.Fatal("expected a non-zero exit")
	}

	if !strings.Contains(stderr, "cauldron:") {
		t.Errorf("stderr should be a cauldron error:\n%s", stderr)
	}
}

// The command has to appear in help, or nobody finds it.
func TestNetworkIsListedInHelp(t *testing.T) {
	stdout, _, code := run(t, "help")

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, "network") {
		t.Errorf("help should list the network command:\n%s", stdout)
	}
}
