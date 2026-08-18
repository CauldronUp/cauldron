package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/CauldronUp/cauldron/internal/detect"
	projectconfig "github.com/CauldronUp/cauldron/internal/project"
)

func runDetect(ctx *context, args []string) int {
	dir := ctx.dir
	if len(args) > 0 && args[0] != "" {
		dir = args[0]
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

	writePlan(ctx.stdout, project)

	// What detection found is not what will be mounted once the project has
	// said otherwise, and printing the one without the other is how somebody
	// spends an afternoon wondering why a provider they can see here is not
	// running.
	config, existed, err := projectconfig.Load(dir)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	if existed {
		fmt.Fprintf(ctx.stdout, "\n%s decides what runs, not the above: %s\n",
			projectconfig.FileName, strings.Join(config.Recipes, ", "))
	}

	return 0
}

// writePlan renders a detected project as the boot plan.
func writePlan(w io.Writer, p *detect.Project) {
	if p.Framework != "" {
		fmt.Fprintf(w, "Detected a %s project.\n\n", p.Framework)
	} else {
		fmt.Fprint(w, "Detected project.\n\n")
	}

	section(w, "Runtime", p.Of(detect.KindRuntime))
	section(w, "Services", p.Of(detect.KindService))
	section(w, "Recipes", p.Of(detect.KindRecipe))

	if len(p.Unmatched) > 0 {
		fmt.Fprint(w, "No recipe yet. These will still reach the real network:\n")

		for _, r := range p.Unmatched {
			fmt.Fprintf(w, "  !  %s  (%s)\n", r.Evidence, r.Source)
		}

		fmt.Fprint(w, "\n")
	}
}

func section(w io.Writer, heading string, reqs []detect.Requirement) {
	if len(reqs) == 0 {
		return
	}

	fmt.Fprintf(w, "%s\n", heading)

	for _, r := range reqs {
		version := r.Version
		if version != "" {
			version = " " + version
		}

		fmt.Fprintf(w, "  +  %s%s  (%s)\n", r.Name, version, r.Evidence)
	}

	fmt.Fprint(w, "\n")
}
