package cli

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/CauldronUp/cauldron/internal/client"
	"github.com/CauldronUp/cauldron/internal/engine"
)

// browsable maps a service to the page worth looking at, and which of its
// published ports serves it. Only services with something to look at appear:
// opening a browser at Postgres helps nobody.
var browsable = map[string]struct {
	port  int
	label string
}{
	"mailpit":     {1, "inbox"},
	"meilisearch": {0, "search"},
	"minio":       {1, "console"},
}

// runOpen opens a running service in a browser.
func runOpen(ctx *context, args []string) int {
	var (
		dir  string
		base string
		show bool
	)

	fs := flag.NewFlagSet("open", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	fs.StringVar(&dir, "path", ".", "project directory")
	fs.BoolVar(&show, "print", false, "print the address instead of opening a browser")
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	targets := openable(dir, base)

	name := arg(positional, 0)
	if name == "" {
		name = "sandbox"
	}

	target, ok := targets[name]
	if !ok {
		fmt.Fprintf(ctx.stderr, "cauldron: nothing to open for %q. Available: %s\n", name, strings.Join(sortedTargets(targets), ", "))
		return 1
	}

	if show {
		fmt.Fprintln(ctx.stdout, target)
		return 0
	}

	if err := openInBrowser(target); err != nil {
		// Failing to launch a browser is not a reason to withhold the address.
		fmt.Fprintf(ctx.stdout, "%s\n", target)
		fmt.Fprintf(ctx.stderr, "cauldron: could not open a browser: %v\n", err)

		return 1
	}

	fmt.Fprintf(ctx.stdout, "Opened %s\n", target)

	return 0
}

// openable lists what is worth opening right now, which depends on what is
// actually running rather than what the catalogue could run.
func openable(dir, base string) map[string]string {
	targets := map[string]string{}

	if status, err := client.New(base).Status(); err == nil {
		sandbox := client.New(base).Base()
		targets["sandbox"] = sandbox + "/_cauldron/status"

		for _, r := range status.Recipes {
			targets[r.Recipe] = sandbox + "/" + r.Recipe
		}
	}

	stdctx, cancel := deadline(15 * time.Second)
	defer cancel()

	eng := engine.New(projectName(dir))

	if err := eng.Available(stdctx); err != nil {
		return targets
	}

	running, err := eng.Running(stdctx)
	if err != nil {
		return targets
	}

	for _, service := range running {
		page, ok := browsable[service]
		if !ok {
			continue
		}

		spec, ok := engine.SpecFor(service)
		if !ok || page.port >= len(spec.Ports) {
			continue
		}

		targets[service] = fmt.Sprintf("http://127.0.0.1:%d", spec.Ports[page.port].Host)
	}

	return targets
}

func sortedTargets(targets map[string]string) []string {
	out := make([]string, 0, len(targets))

	for name := range targets {
		out = append(out, name)
	}

	sort.Strings(out)

	if len(out) == 0 {
		return []string{"nothing is running"}
	}

	return out
}

func openInBrowser(target string) error {
	switch runtime.GOOS {
	case "windows":
		// rundll32 avoids `cmd /c start`, which treats the first quoted
		// argument as a window title and would open a stray console.
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
	case "darwin":
		return exec.Command("open", target).Start()
	default:
		return exec.Command("xdg-open", target).Start()
	}
}
