// Serving a Recipe with the provider's own description added to it.

package cli

import (
	"fmt"
	"io"

	"github.com/CauldronUp/cauldron/internal/openapi"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/server"
)

// specOutcome is what happened to one Recipe when the description was asked
// for, so the caller can report all of them together rather than interleaving
// prose with the mounting.
type specOutcome struct {
	Recipe string
	// Report is empty unless routes were actually added.
	Report recipe.AugmentReport
	// Why is the reason nothing was added, when nothing was.
	Why string
}

// augmentFromSpec fetches the description a Recipe names, drafts a Recipe from
// it, and merges the two with the written one winning.
//
// A Recipe that names no description comes back unchanged and says so. That is
// the common case by a long way -- 301 of the 313 Recipes here have no
// machine-readable description to fetch -- and it is not a failure. Neither is
// a host that will not answer: the Recipe still works, because the Recipe never
// needed the network.
func augmentFromSpec(r *recipe.Recipe, fetch func(string) ([]byte, error)) (*recipe.Recipe, specOutcome) {
	out := specOutcome{Recipe: r.Name}

	if r.Upstream.Spec == "" {
		out.Why = "publishes no description to read"

		return r, out
	}

	raw, err := fetch(r.Upstream.Spec)
	if err != nil {
		out.Why = fmt.Sprintf("description could not be fetched: %v", err)

		return r, out
	}

	doc, err := openapi.Parse(raw)
	if err != nil {
		out.Why = fmt.Sprintf("description could not be read: %v", err)

		return r, out
	}

	drafted, err := recipe.Parse([]byte(openapi.Draft(doc, r.Name)))
	if err != nil {
		// A description this cannot be drafted from is a fact about the
		// description, and the Recipe is unharmed by it.
		out.Why = fmt.Sprintf("description drafted into a Recipe that will not parse: %v", err)

		return r, out
	}

	drafted.Upstream.Spec = r.Upstream.Spec

	// A description declares its paths relative to its own server, and a
	// Recipe carries the whole path a client requests. Adyen says /sessions
	// beside a server of checkout-test.adyen.com/v71, and its Recipe says
	// /v71/... -- so without this every derived route mounts at a path no
	// client would call, reports as added, and answers nothing where anybody
	// would look for it. Serving a route at the wrong path is worse than not
	// serving it at all.
	if base := openapi.BasePath(doc); base != "" {
		for i := range drafted.Routes {
			drafted.Routes[i].Path = base + drafted.Routes[i].Path
		}
	}

	merged, report := recipe.Augment(r, drafted)
	out.Report = report

	if report.Added == 0 {
		out.Why = "description added nothing this Recipe does not already model"
	}

	return merged, out
}

// reportSpecOutcomes prints what the descriptions contributed.
//
// It is deliberately blunt about the distinction. A merged Recipe must not read
// like a written one: a route taken from a schema is un-contradicted, which is
// a smaller claim than observed, and the two answer a caller identically. The
// only place the difference can be made visible is here.
func reportSpecOutcomes(w io.Writer, outcomes []specOutcome) {
	var added, kept, skipped int

	for _, outcome := range outcomes {
		added += outcome.Report.Added
		kept += outcome.Report.Kept
		skipped += outcome.Report.Skipped
	}

	if added == 0 {
		fmt.Fprintln(w, "\nNo routes were added from any description.")
	} else {
		fmt.Fprintf(w, "\n%d route(s) added from provider descriptions, %d already modelled and left alone.\n",
			added, kept)
	}

	for _, outcome := range outcomes {
		switch {
		case outcome.Report.Added > 0:
			fmt.Fprintf(w, "  %-24s +%d derived\n", outcome.Recipe, outcome.Report.Added)
		case outcome.Why != "":
			fmt.Fprintf(w, "  %-24s %s\n", outcome.Recipe, outcome.Why)
		}

		for _, reason := range outcome.Report.Reasons {
			fmt.Fprintf(w, "  %-24s   skipped %s\n", "", reason)
		}
	}

	if skipped > 0 {
		fmt.Fprintf(w, "%d route(s) in the descriptions could not be added safely, listed above.\n", skipped)
	}

	if added > 0 {
		// Said every time, because it is the whole point.
		fmt.Fprintln(w, "\nA derived route is what the provider says it does.")
		fmt.Fprintln(w, "A Recipe route is what it was seen doing.")
		fmt.Fprintln(w, "Those are not the same claim, and only the second one catches a 200 that")
		fmt.Fprintln(w, "means failure, an empty array that is true, or a field the description")
		fmt.Fprintln(w, "forgot it sends.")
	}
}

// mountOne mounts a Recipe, reaching for its description first when asked to.
//
// A failure to reach a description is not a failure to mount: the Recipe is
// mounted regardless, and the reason is recorded for the report. Only a genuine
// mounting error comes back as one.
func mountOne(srv *server.Server, name string, opts serveOptions, outcomes *[]specOutcome) error {
	if !opts.withSpec {
		return srv.Mount(name, opts.seed, "")
	}

	r, err := recipe.Open(name)
	if err != nil {
		return err
	}

	merged, outcome := augmentFromSpec(r, fetchSpec)
	*outcomes = append(*outcomes, outcome)

	return srv.MountRecipe(merged, opts.seed, "")
}
