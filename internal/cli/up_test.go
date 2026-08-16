package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/engine"
)

func TestProjectNameComesFromTheDirectory(t *testing.T) {
	cases := map[string]string{
		filepath.Join("a", "b", "my-shop"):      "my-shop",
		filepath.Join("a", "b", "My_Shop v2"):   "my-shop-v2",
		filepath.Join("a", "b", "Acme.Widgets"): "acme-widgets",
	}

	for input, want := range cases {
		if got := projectName(input); got != want {
			t.Errorf("projectName(%q) = %q, want %q", input, got, want)
		}
	}
}

// Two checkouts of the same repository must not fight over one environment.
func TestProjectNameSeparatesSiblingCheckouts(t *testing.T) {
	a := projectName(filepath.Join("work", "shop-main"))
	b := projectName(filepath.Join("work", "shop-hotfix"))

	if a == b {
		t.Errorf("two checkouts produced the same project name: %q", a)
	}
}

func TestConnectionStringsAreUsable(t *testing.T) {
	for _, service := range []string{"postgres", "mysql", "redis", "meilisearch", "minio", "mailpit"} {
		spec, ok := engine.SpecFor(service)
		if !ok {
			t.Fatalf("no spec for %s", service)
		}

		got := connectionFor(spec)

		if got == "" {
			t.Errorf("%s has no connection description", service)
			continue
		}

		// A description a developer cannot paste somewhere is not much use.
		if !strings.Contains(got, "127.0.0.1") {
			t.Errorf("%s connection is missing a host: %q", service, got)
		}
	}
}

func TestConnectionStringsCarryCredentialsWhereNeeded(t *testing.T) {
	postgres, _ := engine.SpecFor("postgres")

	got := connectionFor(postgres)

	if !strings.Contains(got, "cauldron:cauldron") {
		t.Errorf("a database URL without credentials cannot be pasted into an app: %q", got)
	}
}

func TestUnknownServiceStillGetsAnAddress(t *testing.T) {
	got := connectionFor(engine.Spec{Service: "something", Ports: []engine.Port{{Host: 1234, Container: 1234}}})

	if got != "127.0.0.1:1234" {
		t.Errorf("got %q", got)
	}
}

func TestUpAcceptsItsFlags(t *testing.T) {
	// -no-services keeps the whole command usable on a machine without Docker,
	// which is the difference between a broken tool and a degraded one.
	stdout, stderr, code := run(t, "up", "--no-services", projectWith(t, `{"require":{"php":"^8.5"}}`))

	if code != 1 {
		t.Fatalf("exit = %d, want 1 (nothing to serve)", code)
	}

	if strings.Contains(stdout+stderr, "flag provided but not defined") {
		t.Errorf("--no-services was not accepted\n%s%s", stdout, stderr)
	}
}

func TestDownIsListedInHelp(t *testing.T) {
	stdout, _, code := run(t, "help")

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, "down") {
		t.Errorf("help should list down\n%s", stdout)
	}
}
