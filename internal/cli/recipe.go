package cli

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

func runRecipe(ctx *context, args []string) int {
	if len(args) == 0 {
		fmt.Fprint(ctx.stderr, "cauldron: recipe needs a subcommand\n\n  cauldron recipe list\n  cauldron recipe info <name>\n")
		return 1
	}

	switch args[0] {
	case "list", "ls":
		return recipeList(ctx, args[1:])
	case "info", "show":
		return recipeInfo(ctx, args[1:])
	default:
		fmt.Fprintf(ctx.stderr, "cauldron: unknown recipe subcommand %q\n", args[0])
		return 1
	}
}

func recipeList(ctx *context, args []string) int {
	summaries, err := recipe.Summarise()
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	if len(summaries) == 0 {
		fmt.Fprint(ctx.stdout, "No recipes are bundled with this build.\n")
		return 0
	}

	// A filter, because at a hundred Recipes the useful question is "what can
	// emulate a payment provider" rather than "what is there".
	if len(args) > 0 && args[0] != "" {
		wanted := strings.ToLower(args[0])
		kept := make([]recipe.Summary, 0, len(summaries))

		for _, s := range summaries {
			if s.Capability == wanted {
				kept = append(kept, s)
			}
		}

		if len(kept) == 0 {
			fmt.Fprintf(ctx.stderr, "cauldron: no recipe does %q\n\nTry one of: %s\n",
				args[0], strings.Join(capabilitiesOf(summaries), ", "))

			return 1
		}

		summaries = kept
	}

	w := tabwriter.NewWriter(ctx.stdout, 0, 0, 2, ' ', 0)

	fmt.Fprint(w, "RECIPE\tDOES\tVERSION\tAPI\tRESOURCES\tROUTES\tEVENTS\n")

	for _, s := range summaries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n",
			s.Name, s.Capability, s.Version, s.API, s.Resources, s.Routes, s.Events)
	}

	_ = w.Flush()

	return 0
}

// capabilitiesOf lists the categories present, in a stable order, for the
// message a caller gets when they ask for one that is not.
func capabilitiesOf(summaries []recipe.Summary) []string {
	seen := map[string]bool{}

	var out []string

	for _, s := range summaries {
		if s.Capability == "" || seen[s.Capability] {
			continue
		}

		seen[s.Capability] = true

		out = append(out, s.Capability)
	}

	sort.Strings(out)

	return out
}

func recipeInfo(ctx *context, args []string) int {
	if len(args) == 0 {
		fmt.Fprint(ctx.stderr, "cauldron: recipe info needs a name, e.g. 'cauldron recipe info stripe'\n")
		return 1
	}

	name := args[0]

	r, err := recipe.Open(name)
	if err != nil {
		var notBundled *recipe.ErrNotBundled

		if asNotBundled(err, &notBundled) {
			fmt.Fprintf(ctx.stderr, "cauldron: no recipe named %q\n", name)

			if available := recipe.Bundled(); len(available) > 0 {
				fmt.Fprintf(ctx.stderr, "available: %s\n", strings.Join(available, ", "))
			}

			return 1
		}

		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	fmt.Fprintf(ctx.stdout, "%s %s\n", r.Name, r.Version)
	fmt.Fprintf(ctx.stdout, "  targets API   %s\n", r.Upstream.API)

	if r.Upstream.Docs != "" {
		fmt.Fprintf(ctx.stdout, "  docs          %s\n", r.Upstream.Docs)
	}

	if r.Auth.Scheme != "" {
		fmt.Fprintf(ctx.stdout, "  auth          %s\n", r.Auth.Scheme)
	}

	fmt.Fprint(ctx.stdout, "\nResources\n")

	for _, name := range sortedNames(r.Resources) {
		res := r.Resources[name]
		fmt.Fprintf(ctx.stdout, "  %s  (ids %s…, %d fields)\n", name, res.ID.Prefix, len(res.Fields))
	}

	fmt.Fprint(ctx.stdout, "\nRoutes\n")

	for _, route := range r.Routes {
		fmt.Fprintf(ctx.stdout, "  %-6s %-28s %s\n", route.Method, route.Path, route.Operation)
	}

	if events := r.Events(); len(events) > 0 {
		fmt.Fprint(ctx.stdout, "\nWebhook events\n")

		for _, event := range events {
			fmt.Fprintf(ctx.stdout, "  %s\n", event)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Fprint(ctx.stdout, "\nInjectable failures\n")

		for _, name := range sortedNames(r.Errors) {
			fmt.Fprintf(ctx.stdout, "  %-22s %d\n", name, r.Errors[name].Status)
		}
	}

	if len(r.Fixtures) > 0 {
		fmt.Fprintf(ctx.stdout, "\nFixtures\n  %s\n", strings.Join(sortedNames(r.Fixtures), ", "))
	}

	return 0
}

// asNotBundled keeps the errors.As call in one place; the CLI only needs to
// distinguish "no such recipe" from "this recipe is broken".
func asNotBundled(err error, target **recipe.ErrNotBundled) bool {
	notBundled, ok := err.(*recipe.ErrNotBundled)
	if ok {
		*target = notBundled
	}

	return ok
}

func sortedNames[V any](m map[string]V) []string {
	names := make([]string, 0, len(m))

	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
