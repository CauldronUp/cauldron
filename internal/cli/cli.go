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
		"up": {
			name:    "up",
			summary: "Boot the project and its dependencies",
			usage:   "cauldron up [path]",
			run:     runUp,
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
	fmt.Fprint(w, "Cauldron — your code, one command, every dependency.\n\nUsage:\n  cauldron <command> [arguments]\n\nCommands:\n")

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

	fmt.Fprint(w, "\nRun 'cauldron <command> --help' for details.\n")
}
