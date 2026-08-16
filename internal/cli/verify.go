package cli

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/conform"
	"github.com/CauldronUp/cauldron/internal/recipe"
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
		})

		writeReport(ctx, report, verbose)

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
