package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/CauldronUp/cauldron/internal/client"
)

// urlFlag adds the shared --url flag. Every control command needs it, and a
// developer running two projects at once needs to point at the right one.
func urlFlag(fs *flag.FlagSet, into *string) {
	fs.StringVar(into, "url", "", "address of the running Cauldron server (default $CAULDRON_URL, then "+client.DefaultBase+")")
}

// parseFlags parses flags that may appear before, after or between positional
// arguments.
//
// Go's flag package stops at the first non-flag argument, which would make
// "cauldron seed stripe --fixture small-shop" silently ignore the fixture —
// the most natural way to type the command. Looping keeps both orders working.
func parseFlags(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string

	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}

		rest := fs.Args()
		if len(rest) == 0 {
			return positional, nil
		}

		positional = append(positional, rest[0])
		args = rest[1:]
	}
}

// fail prints an error the way a developer needs to read it: the server's own
// words, and a concrete next step when nothing is listening.
func fail(ctx *context, err error) int {
	var unreachable *client.ErrUnreachable

	if errors.As(err, &unreachable) {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", unreachable)
		return 1
	}

	fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

	return 1
}

func runStatus(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	if _, err := parseFlags(fs, args); err != nil {
		return 1
	}

	status, err := client.New(base).Status()
	if err != nil {
		return fail(ctx, err)
	}

	if len(status.Recipes) == 0 {
		fmt.Fprint(ctx.stdout, "No recipes are running.\n")
		return 0
	}

	fmt.Fprintf(ctx.stdout, "Sandbox time %s\n\n", status.Time)

	w := tabwriter.NewWriter(ctx.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "RECIPE\tVERSION\tFIXTURE\tREQUESTS\tFAULTS\tWEBHOOKS\tNETWORK\n")

	for _, r := range status.Recipes {
		fixture := r.Fixture
		if fixture == "" {
			fixture = "-"
		}

		// Degraded conditions are invisible until they bite, so status is the
		// one place they have to be stated plainly.
		network := "-"
		if len(r.Network) > 0 {
			network = strings.Join(r.Network, "; ")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n", r.Recipe, r.Version, fixture, r.Requests, r.Faults, r.Webhooks, network)
	}

	_ = w.Flush()

	return 0
}

func runRequests(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("requests", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe := arg(positional, 0)
	if recipe == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which recipe? e.g. 'cauldron requests stripe'\n")
		return 1
	}

	exchanges, err := client.New(base).Requests(recipe)
	if err != nil {
		return fail(ctx, err)
	}

	if len(exchanges) == 0 {
		fmt.Fprintf(ctx.stdout, "%s has not been called yet.\n", recipe)
		return 0
	}

	w := tabwriter.NewWriter(ctx.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "#\tMETHOD\tPATH\tSTATUS\tNOTE\n")

	for _, e := range exchanges {
		note := e.Op

		if e.Fault != "" {
			note = "fault: " + e.Fault
		}

		// A degraded request is often the one somebody is squinting at, so the
		// cause belongs in the log rather than only in status.
		if e.Network != "" {
			if note == "" {
				note = "network: " + e.Network
			} else {
				note += "  network: " + e.Network
			}
		}

		status := fmt.Sprintf("%d", e.Status)
		if e.Status == 0 {
			// A severed connection never produced one.
			status = "-"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", e.Seq, e.Method, e.Path, status, note)
	}

	_ = w.Flush()

	return 0
}

func runSeed(ctx *context, args []string) int {
	var (
		base    string
		fixture string
	)

	fs := flag.NewFlagSet("seed", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&fixture, "fixture", "", "fixture to load")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe := arg(positional, 0)
	if recipe == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which recipe? e.g. 'cauldron seed stripe --fixture small-shop'\n")
		return 1
	}

	if fixture == "" {
		fixture = arg(positional, 1)
	}

	if fixture == "" {
		fmt.Fprint(ctx.stderr, "cauldron: a fixture is required, e.g. --fixture small-shop\n")
		return 1
	}

	if err := client.New(base).Seed(recipe, fixture); err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Seeded %s with %s.\n", recipe, fixture)

	return 0
}

func runReset(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("reset", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe := arg(positional, 0)

	if err := client.New(base).Reset(recipe); err != nil {
		return fail(ctx, err)
	}

	if recipe == "" {
		fmt.Fprint(ctx.stdout, "Reset every recipe and the clock.\n")
	} else {
		fmt.Fprintf(ctx.stdout, "Reset %s.\n", recipe)
	}

	return 0
}

func runFault(ctx *context, args []string) int {
	var (
		base     string
		errName  string
		count    int
		every    int
		duration string
		path     string
	)

	fs := flag.NewFlagSet("fault", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&errName, "error", "", "named failure from the recipe, e.g. rate_limit")
	fs.IntVar(&count, "count", 0, "affect only this many requests")
	fs.IntVar(&every, "every", 0, "fail one request in N")
	fs.StringVar(&duration, "for", "", "expire after this long, e.g. 30s")
	fs.StringVar(&path, "path", "", "only affect paths containing this")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe := arg(positional, 0)
	if recipe == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which recipe? e.g. 'cauldron fault stripe --error rate_limit'\n")
		return 1
	}

	if errName == "" {
		errName = arg(positional, 1)
	}

	if errName == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which failure? e.g. --error rate_limit\n(see 'cauldron recipe info "+recipe+"')\n")
		return 1
	}

	fault := client.Fault{Error: errName, Count: count, Every: every, For: duration, Path: path}

	if err := client.New(base).Arm(recipe, fault); err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Armed %s on %s%s.\n", errName, recipe, describeFault(fault))

	return 0
}

func describeFault(f client.Fault) string {
	var parts []string

	if f.Count > 0 {
		parts = append(parts, fmt.Sprintf("for %d request(s)", f.Count))
	}

	if f.Every > 1 {
		parts = append(parts, fmt.Sprintf("one in %d", f.Every))
	}

	if f.For != "" {
		parts = append(parts, "for "+f.For)
	}

	if f.Path != "" {
		parts = append(parts, "on paths containing "+f.Path)
	}

	if len(parts) == 0 {
		return ""
	}

	return " " + strings.Join(parts, ", ")
}

// runNetwork arms degraded network conditions.
//
// The flag names are Toxiproxy's on purpose. Anyone who has run
// `toxiproxy-cli toxic add -t latency` against their database already knows
// what --latency and --jitter mean here, and making them learn a second
// vocabulary for the same idea would buy nothing.
func runNetwork(ctx *context, args []string) int {
	var (
		base        string
		latency     string
		jitter      string
		timeout     string
		duration    string
		path        string
		bandwidth   int
		limit       int
		slice       int
		count       int
		probability float64
		reset       bool
		clear       bool
	)

	fs := flag.NewFlagSet("network", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&latency, "latency", "", "delay every response, e.g. 800ms")
	fs.StringVar(&jitter, "jitter", "", "vary the delay by plus or minus this, e.g. 200ms")
	fs.IntVar(&bandwidth, "bandwidth", 0, "throttle the response body, in KB/s")
	fs.StringVar(&timeout, "timeout", "", "answer nothing, then close the connection after this long")
	fs.BoolVar(&reset, "reset", false, "close the connection immediately, with no response")
	fs.IntVar(&limit, "limit", 0, "close after this many bytes of body")
	fs.IntVar(&slice, "slice", 0, "write the body in chunks of roughly this many bytes")
	fs.Float64Var(&probability, "probability", 0, "affect this share of requests, 0 to 1 (default all)")
	fs.IntVar(&count, "count", 0, "affect only this many requests")
	fs.StringVar(&duration, "for", "", "expire after this long, e.g. 30s")
	fs.StringVar(&path, "path", "", "only affect paths containing this")
	fs.BoolVar(&clear, "clear", false, "remove every armed condition")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe := arg(positional, 0)
	if recipe == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which recipe? e.g. 'cauldron network stripe --latency 800ms'\n")
		return 1
	}

	c := client.New(base)

	if clear {
		if _, err := c.Degrade(recipe, client.Network{Clear: true}); err != nil {
			return fail(ctx, err)
		}

		fmt.Fprintf(ctx.stdout, "Cleared network conditions on %s.\n", recipe)

		return 0
	}

	network := client.Network{
		Latency:     latency,
		Jitter:      jitter,
		Bandwidth:   bandwidth,
		Timeout:     timeout,
		Reset:       reset,
		Limit:       limit,
		Slice:       slice,
		Probability: probability,
		Count:       count,
		For:         duration,
		Path:        path,
	}

	armed, err := c.Degrade(recipe, network)
	if err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Degraded %s: %s.\n", recipe, armed)

	return 0
}

func runEmit(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("emit", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	recipe, event := arg(positional, 0), arg(positional, 1)

	if recipe == "" || event == "" {
		fmt.Fprint(ctx.stderr, "cauldron: usage is 'cauldron emit <recipe> <event>'\n")
		return 1
	}

	id, err := client.New(base).Emit(recipe, event, nil)
	if err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Emitted %s from %s (%s).\n", event, recipe, id)

	return 0
}

func runClock(ctx *context, args []string) int {
	var base string

	fs := flag.NewFlagSet("clock", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	if arg(positional, 0) != "advance" || arg(positional, 1) == "" {
		fmt.Fprint(ctx.stderr, "cauldron: usage is 'cauldron clock advance <duration>', e.g. 30d\n")
		return 1
	}

	now, err := client.New(base).Advance(arg(positional, 1))
	if err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Sandbox time is now %s.\n", now)

	return 0
}

// writeControlHelp is used by the help listing.
func writeControlHelp(w io.Writer) {
	fmt.Fprint(w, "\nControl a running server with --url, $CAULDRON_URL, or the default "+client.DefaultBase+".\n")
}

// arg returns the nth positional argument, or an empty string.
func arg(positional []string, n int) string {
	if n < len(positional) {
		return positional[n]
	}

	return ""
}
