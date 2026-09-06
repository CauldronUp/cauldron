package cli

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/conform"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/runtime"
	"github.com/CauldronUp/cauldron/internal/server"
)

// runVerify runs each Recipe's conformance cases against the emulator.
//
// It runs in process and needs no credentials, because a check that only the
// maintainer can run is a check that stops being run. The report distinguishes
// cases observed against the real API from cases read in the documentation, so
// "verified" never quietly means "we tested our fake against itself".
func runVerify(ctx *context, args []string) int {
	var verbose bool

	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	fs.BoolVar(&verbose, "v", false, "list every case, not only the failures")

	names, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	if len(names) == 0 {
		names = recipe.Bundled()
	}

	sort.Strings(names)

	var (
		cases, passed, observed, documented int
		guessed, guessedRecipes             int
		unstated, unstatedRecipes           int
		unsent, unsentRecipes               int
		failedRecipes                       []string
	)

	for _, name := range names {
		srv := server.New()

		if err := srv.Mount(name, 1, ""); err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
			return 1
		}

		sandbox, ok := srv.Sandbox(name)
		if !ok {
			fmt.Fprintf(ctx.stderr, "cauldron: %s did not mount\n", name)
			return 1
		}

		report := conform.Run(sandbox.Recipe(), srv, "/"+name, func(fixture string) error {
			return srv.SeedRecipe(name, fixture)
		}, func(errorName string) error {
			sandbox.ClearFaults()

			if errorName == "" {
				return nil
			}

			// Count 1, so the failure covers this case's single request and
			// nothing after it even if the clear above is ever missed.
			return sandbox.Arm(runtime.Fault{Error: errorName, Count: 1})
		}, func() []conform.Delivery {
			recorded := sandbox.Webhooks().Deliveries()
			out := make([]conform.Delivery, 0, len(recorded))

			for _, d := range recorded {
				out = append(out, conform.Delivery{Event: d.Event, Payload: d.Payload, SignatureHeader: d.SignatureHeader, Signature: d.Signature, Headers: d.Headers})
			}

			return out
		})

		writeReport(ctx, report, verbose)

		// Paging nobody checked, counted the same way the evidence is. A
		// Recipe can be entirely green and still be guessing at how its
		// provider pages, because a case that never sends a page size
		// cannot notice the wrong name being read.
		if n := sandbox.Recipe().GuessedPagination(); n > 0 {
			guessed += n
			guessedRecipes++

			fmt.Fprintf(ctx.stdout, "  %s\n", guessedLine(n))
		}

		// The other half of the same omission, and the larger one. A listing
		// that declares no paging is still paged: the runtime gives it ten and
		// reads "limit". The count above cannot see these, because it starts
		// from a declared page size.
		if n := sandbox.Recipe().UnstatedPagination(); n > 0 {
			unstated += n
			unstatedRecipes++
		}

		// And the third kind, which the other two cannot see: a name that is
		// declared and that no case ever sends. The Recipe format already
		// refuses a response field name nothing asserts, on the grounds that
		// it could be renamed to anything unnoticed. The request half had no
		// such rule, and it is the half a client gets wrong -- asserting the
		// response to a listing says nothing about which parameter produced
		// it.
		if n := sandbox.Recipe().UnsentPagingParam(); n > 0 {
			unsent += n
			unsentRecipes++
		}

		recipeObserved, recipeDocumented, _ := report.Provenance()

		cases += len(report.Results)
		passed += report.Passed()
		observed += recipeObserved
		documented += recipeDocumented

		if len(report.Failed()) > 0 {
			failedRecipes = append(failedRecipes, name)
		}
	}

	if cases == 0 {
		fmt.Fprint(ctx.stdout, "No conformance cases yet. A Recipe without them is a claim without evidence.\n")
		return 0
	}

	fmt.Fprintf(ctx.stdout, "\n%d of %d cases passed across %d recipe(s).\n", passed, cases, len(names))
	fmt.Fprintf(ctx.stdout, "%s\n", provenanceLine(observed, documented))

	if guessed > 0 {
		fmt.Fprintf(ctx.stdout, "%d route(s) across %d recipe(s) page by a parameter nobody named.\n",
			guessed, guessedRecipes)
	}

	if unstated > 0 {
		fmt.Fprintf(ctx.stdout, "%d more across %d recipe(s) declare no paging at all, and are paged at ten reading \"limit\".\n",
			unstated, unstatedRecipes)
	}

	if unsent > 0 {
		fmt.Fprintf(ctx.stdout, "%d paging parameter name(s) across %d recipe(s) are declared and sent by no case, so renaming them would break nothing.\n",
			unsent, unsentRecipes)
	}

	if len(failedRecipes) > 0 {
		fmt.Fprintf(ctx.stderr, "\nFailing: %s\n", strings.Join(failedRecipes, ", "))
		return 1
	}

	return 0
}

func writeReport(ctx *context, report conform.Report, verbose bool) {
	observed, documented, latest := report.Provenance()

	fmt.Fprintf(ctx.stdout, "%s %s\n", report.Recipe, report.Version)

	if len(report.Results) == 0 {
		fmt.Fprint(ctx.stdout, "  no conformance cases\n")
		return
	}

	fmt.Fprintf(ctx.stdout, "  %d of %d cases passed\n", report.Passed(), len(report.Results))
	fmt.Fprintf(ctx.stdout, "  %s\n", provenanceLine(observed, documented))

	if latest != "" {
		fmt.Fprintf(ctx.stdout, "  last checked against the real API on %s\n", latest)
	}

	for _, result := range report.Results {
		if result.Passed() && !verbose {
			continue
		}

		mark := "fail"
		if result.Passed() {
			mark = "ok  "
		}

		fmt.Fprintf(ctx.stdout, "  %s %s\n", mark, result.Case.Name)

		if result.Passed() {
			continue
		}

		for _, failure := range result.Failures {
			fmt.Fprintf(ctx.stdout, "         %s\n", failure)
		}

		fmt.Fprintf(ctx.stdout, "         claim: %s\n", result.Case.Source)
	}
}

// provenanceLine states what the evidence rests on. Documentation-only cases
// are worth having, but they are not the same as having watched the provider
// do it, and the difference should not need digging for.
func provenanceLine(observed, documented int) string {
	switch {
	case observed == 0 && documented == 0:
		return "no evidence recorded"
	case observed == 0:
		return fmt.Sprintf("%d from documentation only, none checked against the real API", documented)
	case documented == 0:
		return fmt.Sprintf("all %d checked against the real API", observed)
	default:
		return fmt.Sprintf("%d checked against the real API, %d from documentation only", observed, documented)
	}
}

// guessedLine reports paging the runtime has to guess at.
//
// A declared page size with no style and no parameter name means the runtime
// reads "limit", which is right for some providers and wrong for plenty. The
// wrongness is invisible from the cases: the page size is ignored, one full
// page comes back, and a paging loop written against it runs once and passes.
// Saying so beside the evidence is the same bargain the provenance line makes.
func guessedLine(n int) string {
	if n == 1 {
		return "1 route pages by a parameter nobody named"
	}

	return fmt.Sprintf("%d routes page by a parameter nobody named", n)
}
