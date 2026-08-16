package cli

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/CauldronUp/cauldron/internal/engine"
)

// runLogs shows recent output from one of this project's services.
func runLogs(ctx *context, args []string) int {
	var (
		dir   string
		lines int
	)

	fs := flag.NewFlagSet("logs", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	fs.StringVar(&dir, "path", ".", "project directory")
	fs.IntVar(&lines, "n", 50, "how many lines to show")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	stdctx, cancel := deadline(30 * time.Second)
	defer cancel()

	eng := engine.New(projectName(dir))

	if err := eng.Available(stdctx); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	service := arg(positional, 0)

	if service == "" {
		running, err := eng.Running(stdctx)
		if err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			return 1
		}

		if len(running) == 0 {
			fmt.Fprint(ctx.stderr, "cauldron: nothing is running for this project. Start it with 'cauldron up'\n")
			return 1
		}

		fmt.Fprintf(ctx.stderr, "cauldron: which service? This project has: %s\n", strings.Join(running, ", "))

		return 1
	}

	out, err := eng.Logs(stdctx, service, lines)
	if err != nil {
		var missing *engine.ErrNotRunning
		if errors.As(err, &missing) {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", missing)
			return 1
		}

		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	if strings.TrimSpace(out) == "" {
		fmt.Fprintf(ctx.stdout, "%s has not logged anything yet.\n", service)
		return 0
	}

	fmt.Fprintln(ctx.stdout, strings.TrimRight(out, "\n"))

	return 0
}
