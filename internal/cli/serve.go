package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/CauldronUp/cauldron/internal/detect"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/server"
)

// defaultPort is deliberately not 8080 or 3000: a fake provider must never
// collide with the application under development.
const defaultPort = 4600

type serveOptions struct {
	port    int
	seed    int64
	fixture string
	dir     string
	recipes []string
}

func parseServeFlags(args []string, stderr io.Writer) (serveOptions, error) {
	opts := serveOptions{port: defaultPort, seed: 1, dir: "."}

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.IntVar(&opts.port, "port", opts.port, "port to listen on")
	fs.Int64Var(&opts.seed, "seed", opts.seed, "seed for generated identifiers")
	fs.StringVar(&opts.fixture, "fixture", "", "fixture to seed with, or recipe=fixture pairs")
	fs.StringVar(&opts.dir, "dir", ".", "project directory to detect recipes from")

	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: cauldron serve [flags] [recipe...]\n\nWith no recipes named, Cauldron detects them from the project.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	opts.recipes = fs.Args()

	return opts, nil
}

// plan works out which recipes to mount: the ones named, or the ones detected
// in the project.
func plan(opts serveOptions) ([]string, []string, error) {
	if len(opts.recipes) > 0 {
		return opts.recipes, nil, nil
	}

	project, err := detect.Detect(opts.dir)
	if err != nil {
		if errors.Is(err, detect.ErrNoProject) {
			return nil, nil, fmt.Errorf("no project found in %s. Name the recipes explicitly, e.g. 'cauldron serve stripe'", opts.dir)
		}

		return nil, nil, err
	}

	available := map[string]bool{}
	for _, name := range recipe.Bundled() {
		available[name] = true
	}

	var mount, missing []string

	for _, r := range project.Of(detect.KindRecipe) {
		if available[r.Name] {
			mount = append(mount, r.Name)
			continue
		}

		missing = append(missing, r.Name)
	}

	if len(mount) == 0 {
		return nil, missing, fmt.Errorf("no recipes to serve. This project uses no APIs Cauldron can emulate yet")
	}

	return mount, missing, nil
}

func runServe(ctx *context, args []string) int {
	opts, err := parseServeFlags(args, ctx.stderr)
	if err != nil {
		return 1
	}

	mount, missing, err := plan(opts)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	srv := server.New()
	choice := parseFixture(opts.fixture)

	var unseeded []string

	for _, name := range mount {
		if err := srv.Mount(name, opts.seed, ""); err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			return 1
		}

		fixture, explicit := choice.forRecipe(name)
		if fixture == "" {
			continue
		}

		if !explicit && !srv.HasFixture(name, fixture) {
			// A global fixture is a convenience, not an instruction. Say which
			// recipes it did not fit, rather than failing the whole run or,
			// worse, seeding nothing and saying nothing.
			unseeded = append(unseeded, name)

			continue
		}

		if err := srv.SeedRecipe(name, fixture); err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			return 1
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", opts.port))
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: cannot listen on port %d: %v\n", opts.port, err)
		return 1
	}

	writeServeBanner(ctx.stdout, listener.Addr().String(), srv.Names(), missing, opts)

	if len(unseeded) > 0 {
		fmt.Fprintf(ctx.stdout, "\nNot seeded, no %q fixture: %s. Name one with --fixture %s=<fixture>.\n",
			choice.global, strings.Join(unseeded, ", "), unseeded[0])
	}

	httpServer := &http.Server{
		Handler:           srv,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Serve until interrupted, then shut down cleanly so an in-flight request
	// is not cut off mid-response.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errs := make(chan error, 1)

	go func() {
		errs <- httpServer.Serve(listener)
	}()

	select {
	case <-stop:
		fmt.Fprint(ctx.stdout, "\nStopping.\n")
		_ = httpServer.Close()

		return 0
	case err := <-errs:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			return 1
		}

		return 0
	}
}

func writeServeBanner(w io.Writer, addr string, names, missing []string, opts serveOptions) {
	fmt.Fprintf(w, "Cauldron is listening on http://%s\n\n", addr)

	for _, name := range names {
		fmt.Fprintf(w, "  %-12s http://%s/%s\n", name, addr, name)
	}

	if opts.fixture != "" {
		fmt.Fprintf(w, "\nSeeded with the %q fixture.\n", opts.fixture)
	}

	if len(missing) > 0 {
		fmt.Fprintf(w, "\nNo recipe yet. These will still reach the real network:\n")

		for _, name := range missing {
			fmt.Fprintf(w, "  !  %s\n", name)
		}
	}

	fmt.Fprintf(w, "\nPoint your SDK's base URL at the address above.\nControl it with http://%s/_cauldron/status\n\nPress Ctrl+C to stop.\n", addr)
}

// fixtureChoice resolves the --fixture flag against the recipes being mounted.
//
// A bare name applies to every recipe that ships it. Providers name their
// fixtures differently (stripe has small-shop, github has small-repo), so a
// single global name cannot be mandatory without making the flag useless the
// moment a project uses two providers. Pairs give precision when it matters:
//
//	--fixture small-shop
//	--fixture stripe=small-shop,github=small-repo
type fixtureChoice struct {
	global string
	byName map[string]string
}

func parseFixture(value string) fixtureChoice {
	choice := fixtureChoice{byName: map[string]string{}}

	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		recipe, fixture, isPair := strings.Cut(part, "=")
		if isPair {
			choice.byName[strings.TrimSpace(recipe)] = strings.TrimSpace(fixture)
			continue
		}

		choice.global = part
	}

	return choice
}

// forRecipe returns the fixture to load and whether it was named explicitly.
// An explicit choice must fail loudly if it does not exist; a global one is
// applied only where it fits.
func (f fixtureChoice) forRecipe(name string) (string, bool) {
	if fixture, ok := f.byName[name]; ok {
		return fixture, true
	}

	return f.global, false
}
