package cli

import (
	stdctx "context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CauldronUp/cauldron/internal/detect"
	"github.com/CauldronUp/cauldron/internal/engine"
)

// projectName derives a stable container prefix from the project directory.
//
// It is the directory name rather than anything from a manifest, so two
// checkouts of the same repository side by side get separate environments
// instead of fighting over one.
func projectName(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}

	base := strings.ToLower(filepath.Base(abs))
	base = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")

	if base == "" {
		return "project"
	}

	return base
}

// connection describes how to reach a started service.
var connections = map[string]func(engine.Spec) string{
	"postgres": func(s engine.Spec) string {
		return fmt.Sprintf("postgres://cauldron:cauldron@127.0.0.1:%d/cauldron", s.Ports[0].Host)
	},
	"mysql": func(s engine.Spec) string {
		return fmt.Sprintf("mysql://cauldron:cauldron@127.0.0.1:%d/cauldron", s.Ports[0].Host)
	},
	"redis": func(s engine.Spec) string {
		return fmt.Sprintf("redis://127.0.0.1:%d", s.Ports[0].Host)
	},
	"mailpit": func(s engine.Spec) string {
		return fmt.Sprintf("smtp 127.0.0.1:%d · inbox http://127.0.0.1:%d", s.Ports[0].Host, s.Ports[1].Host)
	},
	"meilisearch": func(s engine.Spec) string {
		return fmt.Sprintf("http://127.0.0.1:%d (key: cauldron)", s.Ports[0].Host)
	},
	"minio": func(s engine.Spec) string {
		return fmt.Sprintf("http://127.0.0.1:%d · console http://127.0.0.1:%d (cauldron / cauldron123)", s.Ports[0].Host, s.Ports[1].Host)
	},
}

func connectionFor(spec engine.Spec) string {
	if describe, ok := connections[spec.Service]; ok {
		return describe(spec)
	}

	if len(spec.Ports) > 0 {
		return fmt.Sprintf("127.0.0.1:%d", spec.Ports[0].Host)
	}

	return ""
}

func runUp(ctx *context, args []string) int {
	var (
		port       int
		host       string
		fixture    string
		skipDocker bool
		headless   bool
	)

	fs := flag.NewFlagSet("up", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	fs.IntVar(&port, "port", defaultPort, "port for the emulated providers")
	fs.StringVar(&fixture, "fixture", "", "fixture to seed each recipe with")
	fs.StringVar(&host, "host", loopback, "interface to bind, e.g. 0.0.0.0 to be reachable from a container")
	fs.BoolVar(&skipDocker, "no-services", false, "skip containers and only run the emulated providers")
	fs.BoolVar(&headless, "headless", false, "emulate the providers only: no containers, no plan, one line of JSON")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	dir := arg(positional, 0)
	if dir == "" {
		dir = ctx.dir
	}

	project, err := detect.Detect(dir)
	if err != nil {
		if errors.Is(err, detect.ErrNoProject) {
			fmt.Fprintf(ctx.stderr, "cauldron: no composer.json, package.json or go.mod in %s\n", dir)
			return 1
		}

		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	// Headless means the emulated providers and nothing else. No containers,
	// and no plan describing an environment Cauldron is not going to set up:
	// the database, the queue and the web server are already somebody else's,
	// and printing a plan for them invites the belief that Cauldron took them
	// over.
	if !headless {
		writePlan(ctx.stdout, project)
	}

	if !skipDocker && !headless {
		if code := startServices(ctx, dir, project); code != 0 {
			return code
		}
	}

	serveArgs := []string{"-dir", dir, "-port", fmt.Sprint(port), "-host", host}
	if fixture != "" {
		serveArgs = append(serveArgs, "-fixture", fixture)
	}

	if headless {
		serveArgs = append(serveArgs, "-headless")
	}

	return runServe(ctx, serveArgs)
}

// startServices boots the containers a project needs and waits for them.
func startServices(ctx *context, dir string, project *detect.Project) int {
	name := projectName(dir)
	docker := engine.New(name)

	background, cancel := deadline(3 * time.Minute)
	defer cancel()

	if err := docker.Available(background); err != nil {
		// Docker being absent is not fatal. The emulated providers are pure Go
		// and still work, so say what was skipped and carry on.
		fmt.Fprintf(ctx.stdout, "Skipping services: %v\n", err)
		fmt.Fprint(ctx.stdout, "The emulated providers below do not need Docker.\n\n")

		return 0
	}

	var (
		specs   []engine.Spec
		skipped []string
	)

	for _, requirement := range project.Of(detect.KindService) {
		spec, ok := engine.SpecFor(requirement.Name)
		if !ok {
			skipped = append(skipped, requirement.Name)
			continue
		}

		specs = append(specs, spec)
	}

	// Mail is wanted by essentially every application and is never a declared
	// dependency, so it is offered whenever anything else is being booted.
	if len(specs) > 0 && !project.Has(detect.KindService, "mailpit") {
		if spec, ok := engine.SpecFor("mailpit"); ok {
			specs = append(specs, spec)
		}
	}

	if len(specs) == 0 {
		if len(skipped) > 0 {
			fmt.Fprintf(ctx.stdout, "No containers to start. Cauldron has no image for: %s\n\n", strings.Join(skipped, ", "))
		}

		return 0
	}

	if err := engine.CheckPorts(specs); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		fmt.Fprint(ctx.stderr, "Stop whatever is holding it, or run 'cauldron up --no-services'.\n")

		return 1
	}

	if err := docker.EnsureNetwork(background); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	fmt.Fprint(ctx.stdout, "Starting services\n")

	for _, spec := range specs {
		fmt.Fprintf(ctx.stdout, "  ·  %s", spec.Service)

		if err := docker.Start(background, spec); err != nil {
			fmt.Fprintf(ctx.stdout, "\n")
			fmt.Fprintf(ctx.stderr, "cauldron: starting %s: %v\n", spec.Service, err)

			return 1
		}

		if err := docker.WaitHealthy(background, spec.Service, 90*time.Second); err != nil {
			fmt.Fprintf(ctx.stdout, "\n")
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			fmt.Fprintf(ctx.stderr, "Inspect it with: docker logs %s\n", docker.ContainerName(spec.Service))

			return 1
		}

		fmt.Fprintf(ctx.stdout, "\r  +  %-12s %s\n", spec.Service, connectionFor(spec))
	}

	if len(skipped) > 0 {
		fmt.Fprintf(ctx.stdout, "\nNo image for: %s. Start these yourself.\n", strings.Join(skipped, ", "))
	}

	fmt.Fprint(ctx.stdout, "\nRuntimes are not containerised yet, so run your application as you normally do.\n\n")

	return 0
}

func runDown(ctx *context, args []string) int {
	var keepData bool

	fs := flag.NewFlagSet("down", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	fs.BoolVar(&keepData, "keep-data", false, "keep volumes so databases survive")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	dir := arg(positional, 0)
	if dir == "" {
		dir = ctx.dir
	}

	name := projectName(dir)
	docker := engine.New(name)

	background, cancel := deadline(2 * time.Minute)
	defer cancel()

	if err := docker.Available(background); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	running, err := docker.Running(background)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	if err := docker.Stop(background, keepData); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	if len(running) == 0 {
		fmt.Fprintf(ctx.stdout, "Nothing was running for %s.\n", name)
		return 0
	}

	fmt.Fprintf(ctx.stdout, "Stopped %s.\n", strings.Join(running, ", "))

	if keepData {
		fmt.Fprint(ctx.stdout, "Volumes kept.\n")
	}

	return 0
}

// deadline returns a context bounded by d.
//
// The stdlib context package is imported under an alias because this package
// already has its own context type for command wiring, and shadowing that
// would be a trap for the next person adding a command.
func deadline(d time.Duration) (stdctx.Context, stdctx.CancelFunc) {
	return stdctx.WithTimeout(stdctx.Background(), d)
}
