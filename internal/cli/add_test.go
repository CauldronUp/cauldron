package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/project"
)

// Detection covers the common case and cannot cover all of it. A project that
// talks to a provider over raw HTTP has no dependency to find, and neither has
// one using any of the two dozen Recipes no package maps to yet. Without
// somewhere to write it down the answer was retyping the names on every
// command.

func projectDir(t *testing.T, manifest string) string {
	t.Helper()

	dir := t.TempDir()

	if manifest != "" {
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(manifest), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}
	}

	return dir
}

func TestAddWritesTheFileAndKeepsWhatWasDetected(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	stdout, _, code := run(t, "add", "--dir", dir, "mercury")
	if code != 0 {
		t.Fatalf("exit code = %d: %s", code, stdout)
	}

	config, existed, err := project.Load(dir)
	if err != nil || !existed {
		t.Fatalf("load: %v, existed %v", err, existed)
	}

	// The first write copies what detection found, because the file replaces
	// detection once it exists and starting one must not quietly drop the
	// providers already being emulated.
	if !config.Has("stripe") {
		t.Errorf("the detected provider was dropped: %v", config.Recipes)
	}

	if !config.Has("mercury") {
		t.Errorf("the added provider is missing: %v", config.Recipes)
	}
}

func TestAddIsIdempotent(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	run(t, "add", "--dir", dir, "mercury")
	stdout, _, code := run(t, "add", "--dir", dir, "mercury")

	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}

	if !strings.Contains(stdout, "already there") {
		t.Errorf("stdout = %q", stdout)
	}

	config, _, _ := project.Load(dir)

	count := 0

	for _, name := range config.Recipes {
		if name == "mercury" {
			count++
		}
	}

	if count != 1 {
		t.Errorf("mercury appears %d times: %v", count, config.Recipes)
	}
}

// A typo in the third argument must not leave the first two half-applied.
func TestAddWritesNothingWhenAnyNameIsUnknown(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	_, stderr, code := run(t, "add", "--dir", dir, "mercury", "stripo")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "did you mean stripe") {
		t.Errorf("stderr = %q", stderr)
	}

	if _, existed, _ := project.Load(dir); existed {
		t.Error("a rejected add wrote the file anyway")
	}
}

func TestPlanPrefersTheProjectFileOverDetection(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	if err := project.Save(dir, &project.Config{Recipes: []string{"mercury"}}); err != nil {
		t.Fatalf("save: %v", err)
	}

	mount, _, err := plan(serveOptions{dir: dir})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	// The file is the answer once it exists, because a list that only added
	// to detection could never take anything away, and the surprise being
	// fixed is usually something detection got wrong.
	if len(mount) != 1 || mount[0] != "mercury" {
		t.Errorf("mount = %v, want [mercury]", mount)
	}
}

func TestPlanStillPrefersNamedRecipesOverTheFile(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	if err := project.Save(dir, &project.Config{Recipes: []string{"mercury"}}); err != nil {
		t.Fatalf("save: %v", err)
	}

	mount, _, err := plan(serveOptions{dir: dir, recipes: []string{"notion"}})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	if len(mount) != 1 || mount[0] != "notion" {
		t.Errorf("mount = %v, want [notion]", mount)
	}
}

func TestAnEmptyProjectFileSaysWhatToDo(t *testing.T) {
	dir := projectDir(t, `{"dependencies":{"stripe":"^17.0.0"}}`)

	if err := os.WriteFile(project.Path(dir), []byte("recipes: []\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, _, err := plan(serveOptions{dir: dir})
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "cauldron add") {
		t.Errorf("err = %v", err)
	}
}
