package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, args ...string) (string, string, int) {
	t.Helper()

	var stdout, stderr bytes.Buffer

	code := Run(args, &stdout, &stderr)

	return stdout.String(), stderr.String(), code
}

func laravelProject(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	manifest := `{"require": {
		"php": "^8.5",
		"laravel/framework": "^13.0",
		"predis/predis": "^2.2",
		"stripe/stripe-php": "^17.0",
		"acme/weather-api-client": "^1.0"
	}}`

	if err := os.WriteFile(filepath.Join(root, "composer.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	return root
}

func TestHelpListsCommands(t *testing.T) {
	stdout, _, code := run(t, "help")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	for _, want := range []string{"detect", "up", "version"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help output is missing %q\n%s", want, stdout)
		}
	}
}

func TestNoArgumentsPrintsHelp(t *testing.T) {
	stdout, _, code := run(t)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	if !strings.Contains(stdout, "Usage:") {
		t.Errorf("expected usage, got %q", stdout)
	}
}

func TestUnknownCommandFails(t *testing.T) {
	_, stderr, code := run(t, "summon")

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestVersion(t *testing.T) {
	stdout, _, code := run(t, "version")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	if !strings.HasPrefix(stdout, "cauldron ") {
		t.Errorf("stdout = %q", stdout)
	}
}

func TestDetectPrintsThePlan(t *testing.T) {
	stdout, _, code := run(t, "detect", laravelProject(t))

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	for _, want := range []string{
		"Laravel project",
		"Runtime",
		"php 8.5",
		"Services",
		"redis",
		"Recipes",
		"stripe",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("plan is missing %q\n%s", want, stdout)
		}
	}
}

func TestDetectWarnsAboutUnmatchedClients(t *testing.T) {
	stdout, _, code := run(t, "detect", laravelProject(t))

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	if !strings.Contains(stdout, "still reach the real network") {
		t.Errorf("expected an explicit warning about unfaked dependencies\n%s", stdout)
	}

	if !strings.Contains(stdout, "acme/weather-api-client") {
		t.Errorf("expected the unmatched package to be named\n%s", stdout)
	}
}

func TestDetectFailsOutsideAProject(t *testing.T) {
	_, stderr, code := run(t, "detect", t.TempDir())

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "no composer.json") {
		t.Errorf("stderr = %q", stderr)
	}
}

// `up` prints the plan before it acts, so a developer can see what Cauldron
// believes the project needs even when the run then fails.
//
// --no-services is mandatory here. No test may start a container: it would be
// slow, it would leave state on the developer's machine, and it would make the
// suite depend on a running Docker daemon.
func TestUpShowsThePlanBeforeActing(t *testing.T) {
	dir := projectWith(t, `{"require":{"php":"^8.5","predis/predis":"^2.2"}}`)

	stdout, _, code := run(t, "up", "--no-services", dir)

	if code != 1 {
		t.Fatalf("exit code = %d, want 1 because there is nothing to serve", code)
	}

	if !strings.Contains(stdout, "redis") {
		t.Errorf("up should show the detected plan\n%s", stdout)
	}

	if strings.Contains(stdout, "Starting services") {
		t.Errorf("--no-services must not touch Docker\n%s", stdout)
	}
}
