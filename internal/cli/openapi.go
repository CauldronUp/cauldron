package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/CauldronUp/cauldron/internal/openapi"
	"github.com/CauldronUp/cauldron/internal/recipe"
)

// runImport drafts a Recipe from an OpenAPI description.
//
// It writes to a path the caller chooses and never into the bundled Recipe
// directory, because a draft that could be shipped by being left where it fell
// is a draft that will be. What it produces has no conformance cases and says
// at the top of the file that it is not finished.
func runImport(ctx *context, args []string) int {
	set := flag.NewFlagSet("import", flag.ContinueOnError)
	set.SetOutput(ctx.stderr)

	name := set.String("name", "", "the Recipe name to draft, defaulting to the description's title")
	out := set.String("out", "", "where to write the draft, defaulting to <name>-recipe.yaml")

	if err := set.Parse(args); err != nil {
		return 2
	}

	if set.NArg() != 1 {
		fmt.Fprintln(ctx.stderr, "cauldron: import needs one OpenAPI description\n\nUsage: cauldron import [--name x] [--out path] <spec.yaml|spec.json>")

		return 2
	}

	doc, err := openapi.Load(set.Arg(0))
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	recipeName := *name
	if recipeName == "" {
		recipeName = slug(doc.Info.Title)
	}

	if recipeName == "" {
		fmt.Fprintln(ctx.stderr, "cauldron: the description has no usable title, so pass --name")

		return 1
	}

	target := *out
	if target == "" {
		target = recipeName + "-recipe.yaml"
	}

	draft := openapi.Draft(doc, recipeName)

	// Refused rather than overwritten. A draft is something somebody is part
	// way through editing, and losing that to a second run of the same command
	// would be an unkind way to find out.
	if _, err := os.Stat(target); err == nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %s already exists, and a draft is not worth overwriting; pass --out\n", target)

		return 1
	}

	if err := os.WriteFile(target, []byte(draft), 0o600); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: writing %s: %v\n", target, err)

		return 1
	}

	fmt.Fprintf(ctx.stdout, "Drafted %s from %s.\n\n", target, set.Arg(0))
	fmt.Fprintf(ctx.stdout, "This is not a Recipe yet. It has the paths, the field names and the\n")
	fmt.Fprintf(ctx.stdout, "status codes, which is the half a description carries. It has no\n")
	fmt.Fprintf(ctx.stdout, "conformance cases, because a description does not say what an API lies\n")
	fmt.Fprintf(ctx.stdout, "about, and that is the half worth having.\n\n")
	fmt.Fprintf(ctx.stdout, "The file starts with what somebody has to decide before it ships.\n")

	return 0
}

// runCheck reports where a Recipe and an OpenAPI description disagree.
func runCheck(ctx *context, args []string) int {
	set := flag.NewFlagSet("check", flag.ContinueOnError)
	set.SetOutput(ctx.stderr)

	base := set.String("base", "", "a path prefix the description omits, e.g. /v1; read from the description's servers when not given")
	paging := set.Bool("paging", false, "report the query parameters each listing declares, for filling in paging names")
	all := set.Bool("a", false, "report what the description has and the Recipe does not, as well as disagreements")

	if err := set.Parse(args); err != nil {
		return 2
	}

	if set.NArg() != 2 {
		fmt.Fprintln(ctx.stderr, "cauldron: check needs a Recipe and a description\n\nUsage: cauldron check [-a] [--base /v1] <recipe> <spec.yaml|spec.json>")

		return 2
	}

	r, err := recipe.Open(set.Arg(0))
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	doc, err := openapi.Load(set.Arg(1))
	if err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	// A description declares its paths relative to its own server, and a
	// Recipe carries the whole path a client requests. Compared literally
	// those never match: Box's description says /files/{id} beside a server
	// of https://api.box.com/2.0, while the Recipe says /2.0/files/{id}, and
	// every route in it reads as a path the description does not have.
	//
	// It is written down in the description, so it is read from there rather
	// than left for the caller to notice. Saying so, because a prefix quietly
	// added to every path is exactly the kind of help that is worse than none
	// when it is wrong.
	prefix := *base
	if prefix == "" {
		if prefix = openapi.BasePath(doc); prefix != "" {
			fmt.Fprintf(ctx.stdout, "Using base %s, from the description's own servers. Pass --base to override.\n\n", prefix)
		}
	}

	if *paging {
		for _, report := range openapi.Paging(r, doc, prefix) {
			if !report.Found {
				fmt.Fprintf(ctx.stdout, "  %s: not in the description\n", report.Path)

				continue
			}

			fmt.Fprintf(ctx.stdout, "  %s\n    declares: style=%q limit=%q cursor=%q\n    query:    %s\n",
				report.Path, report.Declared.Style, report.Declared.LimitParam, report.Declared.CursorParam,
				strings.Join(report.Query, ", "))
		}

		return 0
	}

	findings := openapi.Check(r, doc, prefix)

	var disagreements, omissions []openapi.Finding

	for _, finding := range findings {
		if finding.Severity == openapi.Disagrees {
			disagreements = append(disagreements, finding)
			continue
		}

		omissions = append(omissions, finding)
	}

	fmt.Fprintf(ctx.stdout, "%s %s against %s\n\n", r.Name, r.Version, set.Arg(1))

	if len(disagreements) == 0 {
		fmt.Fprintln(ctx.stdout, "  Nothing in this Recipe is contradicted by the description.")
	}

	// Every route missing is not a Recipe full of invented paths. It is the
	// signature of the wrong document: a description of another version, or
	// another product, or a fragment covering none of what this Recipe
	// models.
	//
	// CircleCI is the one that showed it. The cached description declares
	// /api/v1 and the Recipe models v2, so all five of its routes came back
	// as paths the description does not have -- five findings that look like
	// five defects and are one mistake, made by whoever fetched the file.
	//
	// Said before the list rather than after it, because the list is what
	// somebody would otherwise start working through.
	if missing := countMissingPaths(disagreements); missing > 0 && missing == len(r.Routes) {
		fmt.Fprintln(ctx.stdout)
		fmt.Fprintf(ctx.stdout, "  All %d of this Recipe's routes are missing from the description, which\n", missing)
		fmt.Fprintln(ctx.stdout, "  usually means the two are not about the same API: a different version, a")
		fmt.Fprintln(ctx.stdout, "  different product, or a fragment that covers none of this. Worth checking")
		fmt.Fprintln(ctx.stdout, "  before reading the list below.")
		fmt.Fprintln(ctx.stdout)
	}
	// Counted apart from the omissions, because they are not omissions. A
	// path declared in another file has not been compared against anything,
	// and folding it into "things the description has and this Recipe does
	// not" would let a Recipe whose every path went unread report that
	// nothing contradicts it. Lob splits all fifty-eight of its paths out
	// that way.
	if unread := countUnread(omissions); unread > 0 {
		fmt.Fprintf(ctx.stdout, "\n  %d of this Recipe's routes are declared in files this did not read, and were not compared.\n", unread)
	}

	for _, finding := range disagreements {
		fmt.Fprintf(ctx.stdout, "  %s\n    %s\n", finding.Where, finding.What)
	}

	if *all && len(omissions) > 0 {
		fmt.Fprintf(ctx.stdout, "\n  %d thing(s) the description has and this Recipe does not:\n", len(omissions))

		for _, finding := range omissions {
			fmt.Fprintf(ctx.stdout, "    %s\n", finding.What)
		}
	} else if len(omissions) > 0 {
		fmt.Fprintf(ctx.stdout, "\n  %d thing(s) the description has and this Recipe does not. Pass -a to list them.\n", len(omissions))
	}

	// Said every time, in the same words, for the same reason the verify
	// report says its cases came from documentation: the absence of a
	// contradiction is a smaller claim than it looks, and it is the kind of
	// claim that grows in the retelling.
	fmt.Fprintln(ctx.stdout, `
  A description carries paths, names, types and status codes. It does not
  carry the thing a Recipe is for: which status reads like success and is
  not, which field's absence is the signal, which amount is in minor units
  with nothing saying so. Nothing above checks any of that.`)

	if len(disagreements) > 0 {
		return 1
	}

	return 0
}

// slug turns a description's title into a Recipe name.
func slug(title string) string {
	var b strings.Builder

	previousDash := true

	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)

			previousDash = false
		case !previousDash && b.Len() > 0:
			b.WriteRune('-')

			previousDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}

// countUnread counts routes whose path the description keeps in another file.
func countUnread(findings []openapi.Finding) int {
	n := 0

	for _, finding := range findings {
		if strings.Contains(finding.What, "in another file") {
			n++
		}
	}

	return n
}

// countMissingPaths counts routes the description has no path for.
func countMissingPaths(findings []openapi.Finding) int {
	n := 0

	for _, finding := range findings {
		if strings.Contains(finding.What, "declares no such path") {
			n++
		}
	}

	return n
}
