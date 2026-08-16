package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// No test here may start a container. Doctor's Docker checks degrade to a
// warning when Docker is absent, which is exactly the path CI exercises.

func TestCheckRecipesParsesEveryBundledRecipe(t *testing.T) {
	got := checkRecipes()

	if got.state != "ok" {
		t.Fatalf("state = %q, detail = %q", got.state, got.detail)
	}

	if !strings.Contains(got.detail, "all parsing") {
		t.Errorf("detail = %q", got.detail)
	}
}

func TestCheckProjectReadsAManifest(t *testing.T) {
	dir := t.TempDir()

	manifest := `{"require":{"php":"^8.3","laravel/framework":"^13.0","stripe/stripe-php":"^17.0"}}`
	if err := os.WriteFile(filepath.Join(dir, "composer.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := checkProject(dir)

	if got.state == "fail" {
		t.Fatalf("state = %q, detail = %q", got.state, got.detail)
	}

	if !strings.Contains(got.detail, "Laravel") {
		t.Errorf("detail should name the framework; got %q", got.detail)
	}

	if !strings.Contains(got.detail, "1 recipe(s)") {
		t.Errorf("detail should count the recipes; got %q", got.detail)
	}
}

// A dependency with no Recipe still reaches the real network, so doctor says so
// rather than reporting a clean bill of health.
func TestCheckProjectWarnsAboutDependenciesWithNoRecipe(t *testing.T) {
	dir := t.TempDir()

	manifest := `{"require":{"php":"^8.3","acme/weather-api-client":"^2.0"}}`
	if err := os.WriteFile(filepath.Join(dir, "composer.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := checkProject(dir)

	if got.state != "warn" {
		t.Fatalf("state = %q, want warn", got.state)
	}

	if !strings.Contains(got.fix, "real network") {
		t.Errorf("the warning should say what happens; got %q", got.fix)
	}
}

func TestCheckProjectOnADirectoryWithNothingInIt(t *testing.T) {
	got := checkProject(t.TempDir())

	if got.state != "warn" {
		t.Errorf("state = %q, want warn", got.state)
	}

	if !strings.Contains(got.fix, "cauldron serve") {
		t.Errorf("an unrecognised directory should still offer a way forward; got %q", got.fix)
	}
}

func TestCheckProjectOnAMissingPath(t *testing.T) {
	got := checkProject(filepath.Join(t.TempDir(), "nowhere"))

	if got.state != "fail" {
		t.Errorf("state = %q, want fail", got.state)
	}
}

func TestWriteChecksReportsTheWorstState(t *testing.T) {
	cases := []struct {
		checks []check
		want   string
	}{
		{[]check{{name: "a", state: "ok"}}, "ok"},
		{[]check{{name: "a", state: "ok"}, {name: "b", state: "warn"}}, "warn"},
		{[]check{{name: "a", state: "warn"}, {name: "b", state: "fail"}}, "fail"},
		{[]check{{name: "a", state: "fail"}, {name: "b", state: "warn"}}, "fail"},
	}

	for _, c := range cases {
		var out bytes.Buffer

		ctx := &context{stdout: &out, stderr: &out}

		if got := writeChecks(ctx, c.checks); got != c.want {
			t.Errorf("worst = %q, want %q", got, c.want)
		}
	}
}

func TestDoctorRunsAndSummarises(t *testing.T) {
	stdout, stderr, code := run(t, "doctor", t.TempDir())

	// Docker may or may not be present on the machine running the tests. Either
	// way doctor must produce a report and must not fail on its absence.
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %s\n%s", code, stderr, stdout)
	}

	for _, want := range []string{"recipes", "project"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("report is missing %q\n%s", want, stdout)
		}
	}
}
