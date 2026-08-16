package detect

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// writeProject materialises a temporary project directory from a map of
// filename to contents.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()

	root := t.TempDir()

	for name, contents := range files {
		path := filepath.Join(root, name)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}

		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	return root
}

func requirement(t *testing.T, p *Project, kind Kind, name string) Requirement {
	t.Helper()

	for _, r := range p.Requirements {
		if r.Kind == kind && r.Name == name {
			return r
		}
	}

	t.Fatalf("expected a %s requirement named %q; got %+v", kind, name, p.Requirements)

	return Requirement{}
}

func TestDetectReturnsErrNoProjectForEmptyDirectory(t *testing.T) {
	root := t.TempDir()

	_, err := Detect(root)

	if !errors.Is(err, ErrNoProject) {
		t.Fatalf("expected ErrNoProject, got %v", err)
	}
}

func TestDetectLaravelProject(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{
			"require": {
				"php": "^8.5",
				"laravel/framework": "^13.0",
				"laravel/horizon": "^6.0",
				"predis/predis": "^2.2",
				"stripe/stripe-php": "^17.0",
				"shopify/shopify-api": "^5.4"
			}
		}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Framework != "Laravel" {
		t.Errorf("Framework = %q, want Laravel", p.Framework)
	}

	if got := requirement(t, p, KindRuntime, "php").Version; got != "8.5" {
		t.Errorf("php version = %q, want 8.5", got)
	}

	for _, name := range []string{"redis", "horizon"} {
		if !p.Has(KindService, name) {
			t.Errorf("expected service %q to be detected", name)
		}
	}

	for _, name := range []string{"stripe", "shopify"} {
		if !p.Has(KindRecipe, name) {
			t.Errorf("expected recipe %q to be detected", name)
		}
	}
}

func TestDetectRecordsEvidenceAndSource(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{"require": {"stripe/stripe-php": "^17.0"}}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	stripe := requirement(t, p, KindRecipe, "stripe")

	if stripe.Evidence != "stripe/stripe-php" {
		t.Errorf("Evidence = %q, want stripe/stripe-php", stripe.Evidence)
	}

	if stripe.Source != "composer.json" {
		t.Errorf("Source = %q, want composer.json", stripe.Source)
	}
}

func TestDetectNodeProject(t *testing.T) {
	root := writeProject(t, map[string]string{
		"package.json": `{
			"engines": {"node": ">=24.0.0"},
			"dependencies": {
				"stripe": "^18.0.0",
				"@shopify/shopify-api": "^11.0.0",
				"ioredis": "^5.4.1"
			},
			"devDependencies": {"vite": "^7.0.0"}
		}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if got := requirement(t, p, KindRuntime, "node").Version; got != "24.0.0" {
		t.Errorf("node version = %q, want 24.0.0", got)
	}

	for _, name := range []string{"stripe", "shopify"} {
		if !p.Has(KindRecipe, name) {
			t.Errorf("expected recipe %q", name)
		}
	}

	if !p.Has(KindService, "redis") {
		t.Error("expected redis service from ioredis")
	}

	if !p.Has(KindService, "vite") {
		t.Error("expected vite")
	}
}

func TestDetectGoProjectIgnoresIndirectDependencies(t *testing.T) {
	root := writeProject(t, map[string]string{
		"go.mod": `module example.com/app

go 1.26

require (
	github.com/stripe/stripe-go/v76 v76.25.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/slack-go/slack v0.15.0 // indirect
)
`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if got := requirement(t, p, KindRuntime, "go").Version; got != "1.26" {
		t.Errorf("go version = %q, want 1.26", got)
	}

	if !p.Has(KindRecipe, "stripe") {
		t.Error("expected stripe recipe from versioned module path")
	}

	if !p.Has(KindService, "redis") {
		t.Error("expected redis service")
	}

	if p.Has(KindRecipe, "slack") {
		t.Error("indirect dependencies must not produce recipes")
	}
}

func TestDetectReportsUnmatchedAPIClients(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{"require": {"acme/weather-api-client": "^1.0"}}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "acme/weather-api-client") {
		t.Error("an unknown package must never become a booted recipe")
	}

	if len(p.Unmatched) != 1 {
		t.Fatalf("Unmatched = %+v, want exactly one entry", p.Unmatched)
	}

	if p.Unmatched[0].Evidence != "acme/weather-api-client" {
		t.Errorf("Unmatched evidence = %q", p.Unmatched[0].Evidence)
	}
}

func TestDetectIsDeterministic(t *testing.T) {
	files := map[string]string{
		"composer.json": `{"require": {"php": "^8.5", "twilio/sdk": "^8.0", "stripe/stripe-php": "^17.0"}}`,
		"package.json":  `{"dependencies": {"@octokit/rest": "^21.0.0"}}`,
	}

	var first []Requirement

	for i := 0; i < 5; i++ {
		p, err := Detect(writeProject(t, files))
		if err != nil {
			t.Fatalf("Detect: %v", err)
		}

		if first == nil {
			first = p.Requirements
			continue
		}

		if len(p.Requirements) != len(first) {
			t.Fatalf("run %d produced %d requirements, first run produced %d", i, len(p.Requirements), len(first))
		}

		for j := range first {
			if p.Requirements[j].Kind != first[j].Kind || p.Requirements[j].Name != first[j].Name {
				t.Fatalf("run %d differs at index %d: %+v vs %+v", i, j, p.Requirements[j], first[j])
			}
		}
	}
}

func TestNormaliseVersion(t *testing.T) {
	cases := map[string]string{
		"^8.5":       "8.5",
		"~24.0.1":    "24.0.1",
		">=8.2 <9.0": "8.2",
		"*":          "",
		"":           "",
		"v76.25.0":   "76.25.0",
		">=24.0.0":   "24.0.0",
	}

	for input, want := range cases {
		if got := normaliseVersion(input); got != want {
			t.Errorf("normaliseVersion(%q) = %q, want %q", input, got, want)
		}
	}
}
