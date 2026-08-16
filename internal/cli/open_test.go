package cli

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/engine"
)

func TestOpenPrintsTheSandboxAddress(t *testing.T) {
	url := liveServer(t)

	stdout, stderr, code := run(t, "open", "--print", "--url", url, "-path", t.TempDir())
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	if !strings.Contains(stdout, url) {
		t.Errorf("stdout = %q, want the sandbox address", stdout)
	}
}

func TestOpenPrintsAMountedRecipe(t *testing.T) {
	url := liveServer(t)

	stdout, _, code := run(t, "open", "stripe", "--print", "--url", url, "-path", t.TempDir())
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, url+"/stripe") {
		t.Errorf("stdout = %q", stdout)
	}
}

func TestOpenNamesWhatIsAvailable(t *testing.T) {
	url := liveServer(t)

	_, stderr, code := run(t, "open", "netsuite", "--print", "--url", url, "-path", t.TempDir())

	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "stripe") {
		t.Errorf("an unknown target should list the real ones; got %q", stderr)
	}
}

func TestSortedTargetsSaysSoWhenThereIsNothing(t *testing.T) {
	got := sortedTargets(map[string]string{})

	if len(got) != 1 || !strings.Contains(got[0], "nothing is running") {
		t.Errorf("sortedTargets(empty) = %v", got)
	}
}

// Opening a browser at a database helps nobody, so only services with a page
// worth looking at are listed, and each names a port the service really has.
func TestBrowsableServicesPointAtRealPorts(t *testing.T) {
	for service, page := range browsable {
		spec, ok := engine.SpecFor(service)
		if !ok {
			t.Errorf("%s is browsable but not in the catalogue", service)
			continue
		}

		if page.port >= len(spec.Ports) {
			t.Errorf("%s wants port index %d but publishes %d", service, page.port, len(spec.Ports))
		}
	}
}
