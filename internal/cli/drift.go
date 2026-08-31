// Whether the providers still describe themselves the way the Recipes were
// written against.

package cli

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/CauldronUp/cauldron/internal/openapi"
	"github.com/CauldronUp/cauldron/internal/recipe"
)

// specFetchLimit is how much of a description will be read.
//
// Some of these documents are enormous -- Stripe's is tens of megabytes -- and
// a scan that can be made to allocate without bound by a host answering with
// an infinite stream is a scan that eventually stops being run.
const specFetchLimit = 64 << 20

// fetchSpec is how the command reaches a description in earnest. It is a
// variable so the test can answer without a network, because everything
// interesting about this command is what it does when a host misbehaves.
var fetchSpec = func(url string) ([]byte, error) {
	return fetchSpecWithin(url, 60*time.Second)
}

// fetchProbe is how discovery guesses at a URL.
//
// The deadline is the whole difference. drift fetches twelve descriptions
// somebody wrote down, and waiting a minute for one is right. Discovery
// fetches thousands of addresses nobody promised anything about, and almost
// all of them are wrong -- so a minute spent on each is hours spent learning
// that a host does not serve a file. A description either answers promptly
// or is not there, and a provider slow enough to miss this is one somebody
// can record by hand.
var fetchProbe = func(url string) ([]byte, error) {
	return fetchSpecWithin(url, 8*time.Second)
}

func fetchSpecWithin(url string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Some docs hosts answer HTML to a client that does not ask for anything
	// else, and an HTML error page that parses as YAML is the failure this
	// whole command is built to not report as drift.
	req.Header.Set("Accept", "application/json, application/yaml, text/yaml;q=0.9, */*;q=0.1")
	req.Header.Set("User-Agent", "cauldron-drift/1 (+https://github.com/CauldronUp/cauldron)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s answered %s", url, resp.Status)
	}

	return io.ReadAll(io.LimitReader(resp.Body, specFetchLimit))
}

func runDrift(ctx *context, args []string) int {
	set := flag.NewFlagSet("drift", flag.ContinueOnError)
	set.SetOutput(ctx.stderr)

	record := set.Bool("record", false, "print the upstream lines to paste into each Recipe whose fingerprint is unrecorded or has moved")
	quiet := set.Bool("q", false, "report only the Recipes that moved or could not be reached")

	if err := set.Parse(args); err != nil {
		return 2
	}

	names := set.Args()
	if len(names) == 0 {
		names = recipe.Bundled()
	}

	recipes := make([]*recipe.Recipe, 0, len(names))

	for _, name := range names {
		r, err := recipe.Open(name)
		if err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

			return 1
		}

		recipes = append(recipes, r)
	}

	reports := openapi.Drift(recipes, fetchSpec)

	return reportDrift(ctx, reports, *record, *quiet)
}

// reportDrift prints the scan and returns the exit code.
//
// Only a moved fingerprint fails. An unreachable description is printed and
// does not fail, because a docs host answering 503 has said nothing about
// whether the provider changed anything, and a scan that goes red for that is
// a scan somebody switches off inside a fortnight.
func reportDrift(ctx *context, reports []openapi.DriftReport, record, quiet bool) int {
	counts := map[string]int{}

	sort.Slice(reports, func(i, j int) bool { return reports[i].Recipe < reports[j].Recipe })

	for _, report := range reports {
		counts[report.Status]++

		// Unsupported is left out of the quiet view on purpose. It is a
		// permanent fact, it will not change tomorrow, and a scheduled job
		// repeating it every day is how the actionable lines get lost.
		if quiet && report.Status != openapi.Moved && report.Status != openapi.Unreachable {
			continue
		}

		switch report.Status {
		case openapi.Undeclared:
			if !quiet {
				fmt.Fprintf(ctx.stdout, "  %-24s no description to check\n", report.Recipe)
			}
		case openapi.Unrecorded:
			fmt.Fprintf(ctx.stdout, "  %-24s no fingerprint recorded, nothing compared\n", report.Recipe)
		case openapi.Unchanged:
			fmt.Fprintf(ctx.stdout, "  %-24s unchanged since %s\n", report.Recipe, report.Seen)
		case openapi.Unreachable:
			fmt.Fprintf(ctx.stdout, "  %-24s could not be read: %v\n", report.Recipe, report.Err)
		case openapi.Unsupported:
			fmt.Fprintf(ctx.stdout, "  %-24s a format this cannot read: %v\n", report.Recipe, report.Err)
		case openapi.Moved:
			fmt.Fprintf(ctx.stdout, "  %-24s MOVED since %s\n", report.Recipe, report.Seen)
			fmt.Fprintf(ctx.stdout, "  %-24s   %s\n", "", report.Spec)
		}

		if record && report.Now != "" && report.Status != openapi.Unchanged {
			// Printed rather than written. Editing a Recipe in place would
			// have to preserve the comments, which are most of what a Recipe
			// is worth, and a scan that quietly rewrites the evidence it was
			// asked to check is the wrong shape of tool.
			fmt.Fprintf(ctx.stdout, "  %-24s   spec_hash: %s\n", "", report.Now)
			fmt.Fprintf(ctx.stdout, "  %-24s   spec_seen: \"%s\"\n", "", time.Now().Format("2006-01-02"))
		}
	}

	fmt.Fprintf(ctx.stdout, "\n%d unchanged, %d moved, %d unreachable, %d in a format this cannot read, %d unrecorded, %d with no description to check.\n",
		counts[openapi.Unchanged], counts[openapi.Moved], counts[openapi.Unreachable],
		counts[openapi.Unsupported], counts[openapi.Unrecorded], counts[openapi.Undeclared])

	// Said plainly, because the alternative reading of a green scan is that
	// every Recipe was checked, and most of them were not.
	if counts[openapi.Undeclared]+counts[openapi.Unsupported] > 0 {
		fmt.Fprintf(ctx.stdout, "A Recipe with no description this can read is not verified by this. It is unexamined.\n")
	}

	if counts[openapi.Moved] > 0 {
		return 1
	}

	return 0
}
