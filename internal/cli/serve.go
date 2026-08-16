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
	fs.StringVar(&opts.fixture, "fixture", "", "fixture to seed each recipe with")
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
			return nil, nil, fmt.Errorf("no project found in %s — name the recipes explicitly, e.g. 'cauldron serve stripe'", opts.dir)
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
		return nil, missing, fmt.Errorf("no recipes to serve — this project uses no APIs Cauldron can emulate yet")
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

	for _, name := range mount {
		if err := srv.Mount(name, opts.seed, opts.fixture); err != nil {
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
		fmt.Fprintf(w, "\nNo recipe yet — these will still reach the real network:\n")

		for _, name := range missing {
			fmt.Fprintf(w, "  !  %s\n", name)
		}
	}

	fmt.Fprintf(w, "\nPoint your SDK's base URL at the address above.\nControl it with http://%s/_cauldron/status\n\nPress Ctrl+C to stop.\n", addr)
}
