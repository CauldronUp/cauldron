package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/CauldronUp/cauldron/internal/detect"
	"github.com/CauldronUp/cauldron/internal/project"
	"github.com/CauldronUp/cauldron/internal/recipe"
)

func runAdd(ctx *context, args []string) int {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)

	dir := fs.String("dir", ".", "project directory")

	fs.Usage = func() {
		fmt.Fprint(ctx.stderr, "Usage: cauldron add <recipe>...\n\nAdds providers to this project, so they are emulated without being named\nevery time. Writes "+project.FileName+".\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	names := fs.Args()
	if len(names) == 0 {
		fmt.Fprint(ctx.stderr, "cauldron: add what? e.g. 'cauldron add stripe'\n")

		return 1
	}

	// Every name checked before anything is written, so a typo in the third
	// argument does not leave the first two half-applied.
	for _, name := range names {
		if _, err := recipe.Open(name); err != nil {
			fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

			return 1
		}
	}

	config, existed, err := project.Load(*dir)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	// The first write copies what detection found, because the file replaces
	// detection once it exists and starting one should not quietly drop the
	// providers that were already being emulated.
	var seeded []string

	if !existed {
		config = &project.Config{}

		if found, err := detect.Detect(*dir); err == nil {
			for _, r := range found.Of(detect.KindRecipe) {
				if _, err := recipe.Open(r.Name); err != nil {
					continue
				}

				config.Recipes = append(config.Recipes, r.Name)
				seeded = append(seeded, r.Name)
			}
		}
	}

	var added, already []string

	for _, name := range names {
		if config.Has(name) {
			already = append(already, name)

			continue
		}

		config.Recipes = append(config.Recipes, name)
		added = append(added, name)
	}

	if err := project.Save(*dir, config); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	if !existed {
		fmt.Fprintf(ctx.stdout, "Wrote %s.\n", project.FileName)

		if len(seeded) > 0 {
			fmt.Fprintf(ctx.stdout, "  kept %s, which this project already needed\n", strings.Join(seeded, ", "))
		}
	}

	if len(added) > 0 {
		fmt.Fprintf(ctx.stdout, "  added %s\n", strings.Join(added, ", "))
	}

	for _, name := range already {
		fmt.Fprintf(ctx.stdout, "  %s was already there\n", name)
	}

	return 0
}
