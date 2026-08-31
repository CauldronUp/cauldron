// Looking for the descriptions nobody has recorded yet.

package cli

import (
	"flag"
	"fmt"
	"sort"
	"sync"

	"github.com/CauldronUp/cauldron/internal/openapi"
	"github.com/CauldronUp/cauldron/internal/recipe"
)

// discoverWorkers is how many Recipes are probed at once.
//
// Each Recipe costs up to a few dozen requests, so the collection is a few
// thousand and the whole thing is waiting on the network. Kept modest because
// the hosts being probed did not ask to be.
const discoverWorkers = 8

// proposal is one Recipe and the description found for it.
type proposal struct {
	Recipe string
	Found  *openapi.Candidate
}

// discoverTally is what the search did, for the line at the end.
type discoverTally struct {
	// Searched is how many Recipes were looked for.
	Searched int
	// NoHost is how many had no verified host to search from.
	NoHost int
	// Recorded is how many already name a description.
	Recorded int
}

// runDiscover looks for a published description for every Recipe that names
// none, and proposes the ones it can prove.
//
// Nothing is written. The output is lines to paste, the same bargain drift
// --record makes, and for a sharper reason: a wrong spec URL does not fail
// loudly. It reports drift against a document that was never this provider's,
// every morning, until somebody works out why. A person confirming each URL is
// cheap; the alternative is a collection quietly checking itself against the
// wrong documents.
func runDiscover(ctx *context, args []string) int {
	set := flag.NewFlagSet("discover", flag.ContinueOnError)
	set.SetOutput(ctx.stderr)

	all := set.Bool("a", false, "search even where a Recipe already names a description")

	if err := set.Parse(args); err != nil {
		return 2
	}

	names := set.Args()
	if len(names) == 0 {
		names = recipe.Bundled()
	}

	var (
		mu        sync.Mutex
		wg        sync.WaitGroup
		proposals []proposal
		tally     discoverTally
		failures  []string
		work      = make(chan string)
	)

	for range discoverWorkers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for name := range work {
				r, err := recipe.Open(name)
				if err != nil {
					mu.Lock()
					failures = append(failures, fmt.Sprintf("%s: %v", name, err))
					mu.Unlock()

					continue
				}

				// A Recipe that already names one is drift's business, not
				// this command's.
				if r.Upstream.Spec != "" && !*all {
					mu.Lock()
					tally.Recorded++
					mu.Unlock()

					continue
				}

				hosts := openapi.Hosts(r)
				if len(hosts) == 0 {
					mu.Lock()
					tally.NoHost++
					mu.Unlock()

					continue
				}

				found := openapi.Discover(r, hosts, fetchSpec)

				mu.Lock()
				tally.Searched++

				if found != nil {
					proposals = append(proposals, proposal{Recipe: r.Name, Found: found})
				}
				mu.Unlock()
			}
		}()
	}

	for _, name := range names {
		work <- name
	}

	close(work)
	wg.Wait()

	sort.Strings(failures)

	for _, failure := range failures {
		fmt.Fprintf(ctx.stderr, "cauldron: %s\n", failure)
	}

	return reportDiscover(ctx, proposals, tally)
}

// reportDiscover prints the proposals and returns the exit code.
//
// Always zero. Finding nothing is not a failure -- most providers publish
// nothing, and a command that goes red for the ordinary case is one nobody
// runs twice.
func reportDiscover(ctx *context, proposals []proposal, tally discoverTally) int {
	// Strongest first: the share of the Recipe a document accounts for is the
	// case for believing it is that provider's description.
	sort.SliceStable(proposals, func(i, j int) bool {
		a, b := proposals[i], proposals[j]
		if shareOf(a.Found) != shareOf(b.Found) {
			return shareOf(a.Found) > shareOf(b.Found)
		}

		if a.Found.Matched != b.Found.Matched {
			return a.Found.Matched > b.Found.Matched
		}

		return a.Recipe < b.Recipe
	})

	for _, p := range proposals {
		fmt.Fprintf(ctx.stdout, "  %-24s %d of %d route(s) declared, in %d path(s), %s\n",
			p.Recipe, p.Found.Matched, p.Found.Routes, p.Found.Declared, p.Found.Format)
		fmt.Fprintf(ctx.stdout, "  %-24s   spec: %s\n", "", p.Found.URL)

		if weak(p.Found) {
			fmt.Fprintf(ctx.stdout, "  %-24s   (weak: one route in a %d-path document -- check this is the API and not a docs site)\n",
				"", p.Found.Declared)
		}
	}

	if len(proposals) > 0 {
		fmt.Fprintln(ctx.stdout)
	}

	fmt.Fprintf(ctx.stdout,
		"%d proposed, %d searched and nothing found, %d with no verified host to search from, %d already recorded.\n",
		len(proposals), tally.Searched-len(proposals), tally.NoHost, tally.Recorded)

	if len(proposals) > 0 {
		fmt.Fprintln(ctx.stdout, "A proposal is a document that declares a path the Recipe already models.")
		fmt.Fprintln(ctx.stdout, "That is evidence it belongs to this provider, not proof. Nothing was written.")
	}

	return 0
}

// weak is a proposal resting on almost nothing.
//
// Documentation platforms serve a generic openapi.json at the vendor's
// marketing domain: ramp.com and circleci.com both answer one, with two and
// five paths, describing the docs site rather than the API. One route matching
// inside a document that small is as consistent with coincidence as with
// discovery.
//
// Said out loud rather than filtered out. The judgement belongs to whoever
// knows the provider, and a proposal silently dropped is one nobody can
// overrule.
func weak(c *openapi.Candidate) bool {
	return c.Matched == 1 && c.Declared <= 5
}

func shareOf(c *openapi.Candidate) float64 {
	if c.Routes == 0 {
		return 0
	}

	return float64(c.Matched) / float64(c.Routes)
}
