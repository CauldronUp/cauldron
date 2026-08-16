package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CauldronUp/cauldron/internal/client"
	"github.com/CauldronUp/cauldron/internal/detect"
	"github.com/CauldronUp/cauldron/internal/engine"
	"github.com/CauldronUp/cauldron/internal/recipe"
)

// check is one thing doctor looked at.
type check struct {
	name string
	// state is ok, warn or fail. A warning is something that will not stop the
	// project starting but will surprise someone later.
	state  string
	detail string
	// fix is what to do about it, when there is something to do.
	fix string
}

// runDoctor reports whether this machine and this project can actually run.
//
// The failures worth catching are the boring ones: Docker not running, a port
// already taken by another checkout, a Recipe that stopped parsing. Each is
// obvious once you know, and each costs twenty minutes when you do not.
func runDoctor(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	dir := arg(positional, 0)
	if dir == "" {
		dir = "."
	}

	checks := []check{
		checkProject(dir),
		checkRecipes(),
	}

	checks = append(checks, checkDocker(dir)...)
	checks = append(checks, checkServer(base))

	worst := writeChecks(ctx, checks)

	if worst == "fail" {
		return 1
	}

	return 0
}

func writeChecks(ctx *context, checks []check) string {
	worst := "ok"

	marks := map[string]string{"ok": "ok  ", "warn": "warn", "fail": "fail"}

	for _, c := range checks {
		fmt.Fprintf(ctx.stdout, "%s  %-22s %s\n", marks[c.state], c.name, c.detail)

		if c.fix != "" {
			fmt.Fprintf(ctx.stdout, "      %-22s %s\n", "", c.fix)
		}

		if c.state == "fail" || (c.state == "warn" && worst == "ok") {
			worst = c.state
		}
	}

	switch worst {
	case "fail":
		fmt.Fprint(ctx.stdout, "\nSomething needs fixing before this project will start.\n")
	case "warn":
		fmt.Fprint(ctx.stdout, "\nUsable, with the caveats above.\n")
	default:
		fmt.Fprint(ctx.stdout, "\nEverything checks out.\n")
	}

	return worst
}

func checkProject(dir string) check {
	if _, err := os.Stat(dir); err != nil {
		return check{"project", "fail", "cannot read " + dir, "check the path"}
	}

	found, err := detect.Detect(dir)
	if err != nil {
		return check{
			"project", "warn",
			"no manifests recognised in " + dir,
			"Cauldron reads composer.json, package.json and go.mod. 'cauldron serve <recipe>' works anyway.",
		}
	}

	described := found.Framework
	if described == "" {
		described = "a project"
	}

	detail := fmt.Sprintf("%s, %d service(s), %d recipe(s)",
		described, len(found.Of(detect.KindService)), len(found.Of(detect.KindRecipe)))

	if len(found.Unmatched) > 0 {
		// Worth a warning rather than silence: these are the dependencies that
		// will still reach the real network.
		return check{
			"project", "warn", detail,
			fmt.Sprintf("%d dependency(ies) have no Recipe and will reach the real network", len(found.Unmatched)),
		}
	}

	return check{"project", "ok", detail, ""}
}

// checkRecipes parses every bundled Recipe. A Recipe that stopped parsing is
// invisible until someone tries to use it, which is usually mid-task.
func checkRecipes() check {
	var broken []string

	for _, name := range recipe.Bundled() {
		if _, err := recipe.Open(name); err != nil {
			broken = append(broken, name)
		}
	}

	if len(broken) > 0 {
		return check{
			"recipes", "fail",
			"these do not parse: " + strings.Join(broken, ", "),
			"run 'cauldron recipe info <name>' to see why",
		}
	}

	return check{"recipes", "ok", fmt.Sprintf("%d bundled, all parsing", len(recipe.Bundled())), ""}
}

func checkDocker(dir string) []check {
	ctx, cancel := deadline(20 * time.Second)
	defer cancel()

	eng := engine.New(projectName(dir))

	if err := eng.Available(ctx); err != nil {
		return []check{{
			"docker", "warn",
			"not available",
			"Backing services need it. Recipes do not, so 'cauldron serve' still works.",
		}}
	}

	out := []check{{"docker", "ok", "running", ""}}

	running, err := eng.Running(ctx)
	if err == nil && len(running) > 0 {
		out = append(out, check{"services", "ok", strings.Join(running, ", "), ""})
	}

	// A port held by another project is the failure that looks like a Cauldron
	// bug and is not one, so name the holder rather than the symptom.
	held, err := eng.PublishedPorts(ctx)
	if err != nil {
		return out
	}

	var clashes []string

	for _, name := range engine.Catalogued() {
		spec, ok := engine.SpecFor(name)
		if !ok {
			continue
		}

		for _, port := range spec.Ports {
			owner, taken := held[port.Host]
			if !taken || strings.HasPrefix(owner, eng.ContainerName("")) {
				continue
			}

			clashes = append(clashes, fmt.Sprintf("%d wanted by %s, held by %s", port.Host, name, owner))
		}
	}

	if len(clashes) > 0 {
		out = append(out, check{
			"ports", "warn",
			strings.Join(clashes, "; "),
			"stop the other container, or start this project without that service",
		})
	}

	return out
}

func checkServer(base string) check {
	status, err := client.New(base).Status()
	if err != nil {
		return check{"sandbox", "ok", "not running", "start one with 'cauldron serve' or 'cauldron up'"}
	}

	names := make([]string, 0, len(status.Recipes))
	for _, r := range status.Recipes {
		names = append(names, r.Recipe)
	}

	if len(names) == 0 {
		return check{"sandbox", "warn", "running with no recipes mounted", "name one: 'cauldron serve stripe'"}
	}

	return check{"sandbox", "ok", "serving " + strings.Join(names, ", "), ""}
}
