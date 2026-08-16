package cli

import (
	"strings"
	"testing"
)

func TestRecipeListShowsBundledRecipes(t *testing.T) {
	stdout, _, code := run(t, "recipe", "list")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	for _, want := range []string{"RECIPE", "VERSION", "stripe"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("listing is missing %q\n%s", want, stdout)
		}
	}
}

func TestRecipeInfoDescribesTheRecipe(t *testing.T) {
	stdout, _, code := run(t, "recipe", "info", "stripe")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	for _, want := range []string{
		"stripe 0.1.0",
		"targets API",
		"Resources",
		"customer",
		"payment_intent",
		"Routes",
		"/v1/customers",
		"Webhook events",
		"payment_intent.succeeded",
		"Injectable failures",
		"rate_limit",
		"Fixtures",
		"small-shop",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("info output is missing %q\n%s", want, stdout)
		}
	}
}

func TestRecipeInfoSuggestsAlternativesWhenUnknown(t *testing.T) {
	_, stderr, code := run(t, "recipe", "info", "netsuite")

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, `no recipe named "netsuite"`) {
		t.Errorf("stderr = %q", stderr)
	}

	if !strings.Contains(stderr, "available:") {
		t.Errorf("an unknown name should list what is available; got %q", stderr)
	}
}

func TestRecipeInfoRequiresAName(t *testing.T) {
	_, stderr, code := run(t, "recipe", "info")

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "needs a name") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestRecipeWithoutSubcommandExplainsItself(t *testing.T) {
	_, stderr, code := run(t, "recipe")

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "cauldron recipe list") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestRecipeRejectsUnknownSubcommand(t *testing.T) {
	_, stderr, code := run(t, "recipe", "brew")

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr, "unknown recipe subcommand") {
		t.Errorf("stderr = %q", stderr)
	}
}
