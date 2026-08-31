// Package cli implements the cauldron command surface.
//
// Everyday verbs stay conventional — up, down, status, logs — because muscle
// memory from Docker Compose is an asset. Cauldron only invents vocabulary for
// concepts that have no established name.
package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Version is stamped at build time via -ldflags.
var Version = "dev"

type command struct {
	name    string
	summary string
	usage   string
	run     func(ctx *context, args []string) int
}

// context carries everything a command needs, so commands stay testable
// without touching global state.
type context struct {
	stdout io.Writer
	stderr io.Writer
	// dir is the project directory commands operate on.
	dir string
}

func commands() map[string]command {
	return map[string]command{
		"detect": {
			name:    "detect",
			summary: "Show what Cauldron found in this project",
			usage:   "cauldron detect [path]",
			run:     runDetect,
		},
		"add": {
			name:    "add",
			summary: "Add a provider this project talks to",
			usage:   "cauldron add <recipe>...",
			run:     runAdd,
		},
		"up": {
			name:    "up",
			summary: "Boot the project and its dependencies",
			usage:   "cauldron up [path]",
			run:     runUp,
		},
		"down": {
			name:    "down",
			summary: "Stop this project's services",
			usage:   "cauldron down [--keep-data]",
			run:     runDown,
		},
		"doctor": {
			name:    "doctor",
			summary: "Check whether this machine and project can run",
			usage:   "cauldron doctor [path]",
			run:     runDoctor,
		},
		"logs": {
			name:    "logs",
			summary: "Show recent output from a running service",
			usage:   "cauldron logs <service> [-n 50]",
			run:     runLogs,
		},
		"open": {
			name:    "open",
			summary: "Open a running service in a browser",
			usage:   "cauldron open [service] [--print]",
			run:     runOpen,
		},
		"serve": {
			name:    "serve",
			summary: "Serve the fake APIs this project depends on",
			usage:   "cauldron serve [flags] [recipe...]",
			run:     runServe,
		},
		"recipe": {
			name:    "recipe",
			summary: "Inspect the Recipes Cauldron can emulate",
			usage:   "cauldron recipe list [capability] | cauldron recipe info <name>",
			run:     runRecipe,
		},
		"import": {
			name:    "import",
			summary: "Draft a Recipe from an OpenAPI description",
			usage:   "cauldron import [--name x] [--out path] <spec.yaml|spec.json>",
			run:     runImport,
		},
		"drift": {
			name:    "drift",
			summary: "Report where a provider's description has moved under a Recipe",
			usage:   "cauldron drift [-q] [--record] [recipe...]",
			run:     runDrift,
		},
		"discover": {
			name:    "discover",
			summary: "Look for a published description for the Recipes that name none",
			usage:   "cauldron discover [-a] [recipe...]",
			run:     runDiscover,
		},
		"check": {
			name:    "check",
			summary: "Report where a Recipe and an OpenAPI description disagree",
			usage:   "cauldron check [-a] [--base /v1] <recipe> <spec.yaml|spec.json>",
			run:     runCheck,
		},
		"verify": {
			name:    "verify",
			summary: "Check a Recipe against what the provider is documented to do",
			usage:   "cauldron verify [-v] [recipe...]",
			run:     runVerify,
		},
		"snapshot": {
			name:    "snapshot",
			summary: "Capture or restore the exact state of a sandbox",
			usage:   "cauldron snapshot save|load|list [name]",
			run:     runSnapshot,
		},
		"status": {
			name:    "status",
			summary: "Show what a running Cauldron is doing",
			usage:   "cauldron status [--url]",
			run:     runStatus,
		},
		"requests": {
			name:    "requests",
			summary: "Show what your code actually sent",
			usage:   "cauldron requests <recipe>",
			run:     runRequests,
		},
		"seed": {
			name:    "seed",
			summary: "Load a fixture into a running recipe",
			usage:   "cauldron seed <recipe> --fixture <name>",
			run:     runSeed,
		},
		"reset": {
			name:    "reset",
			summary: "Return a recipe, or everything, to its seeded state",
			usage:   "cauldron reset [recipe]",
			run:     runReset,
		},
		"fault": {
			name:    "fault",
			summary: "Make a provider fail on purpose",
			usage:   "cauldron fault <recipe> --error <name>",
			run:     runFault,
		},
		"network": {
			name:    "network",
			summary: "Make the connection to a provider slow, flaky or dead",
			usage:   "cauldron network <recipe> --latency 800ms [--jitter 200ms]",
			run:     runNetwork,
		},
		"emit": {
			name:    "emit",
			summary: "Fire a webhook at your application",
			usage:   "cauldron emit <recipe> <event>",
			run:     runEmit,
		},
		"clock": {
			name:    "clock",
			summary: "Move sandbox time forward",
			usage:   "cauldron clock advance <duration>",
			run:     runClock,
		},
		"version": {
			name:    "version",
			summary: "Print the Cauldron version",
			usage:   "cauldron version",
			run:     runVersion,
		},
	}
}

// Run executes a command and returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	ctx := &context{stdout: stdout, stderr: stderr, dir: "."}

	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}

	name := args[0]

	switch name {
	case "help", "-h", "--help":
		printHelp(stdout)
		return 0
	case "-v", "--version":
		return runVersion(ctx, nil)
	}

	cmd, ok := commands()[name]
	if !ok {
		fmt.Fprintf(stderr, "cauldron: unknown command %q\n\nRun 'cauldron help' to see what's available.\n", name)
		return 1
	}

	return cmd.run(ctx, args[1:])
}

func runVersion(ctx *context, _ []string) int {
	fmt.Fprintf(ctx.stdout, "cauldron %s\n", Version)

	return 0
}

func printHelp(w io.Writer) {
	fmt.Fprint(w, "Cauldron: your code, one command, every dependency.\n\nUsage:\n  cauldron <command> [arguments]\n\nCommands:\n")

	all := commands()

	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)

	width := 0
	for _, name := range names {
		if len(name) > width {
			width = len(name)
		}
	}

	for _, name := range names {
		fmt.Fprintf(w, "  %s%s  %s\n", name, strings.Repeat(" ", width-len(name)), all[name].summary)
	}

	writeControlHelp(w)

	fmt.Fprint(w, "\nRun 'cauldron <command> --help' for details.\n")
}
