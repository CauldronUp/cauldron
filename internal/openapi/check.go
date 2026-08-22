package openapi

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Finding is one disagreement between a Recipe and a provider's own
// description of itself.
type Finding struct {
	// Where names the part of the Recipe the finding is about, in the same
	// shape the validator uses, so the two read alike.
	Where string
	// What is the disagreement, in a sentence.
	What string
	// Severity is "disagrees" when the Recipe claims something the
	// description contradicts, and "missing" when the Recipe is silent about
	// something the description declares.
	//
	// The distinction matters because only the first is a bug. A Recipe is
	// allowed to model less than an API — every Recipe here does, and the
	// headers say so — but it is not allowed to model it wrongly.
	Severity string
}

const (
	// Disagrees is a claim the description contradicts.
	Disagrees = "disagrees"
	// Missing is something the description has and the Recipe does not.
	Missing = "missing"
	// Unsaid is a claim the description neither supports nor contradicts,
	// because it says nothing on the subject at all.
	//
	// The difference matters most for statuses that belong to the transport
	// rather than to an operation. Nineteen of the thirty-one descriptions
	// fetched here declare no 429 anywhere, and among them are Stripe,
	// Twilio, Slack and Square, every one of which rate limits. Reporting a
	// Recipe's rate_limit as contradicted by that is reporting the fact that
	// rate limits are usually documented in prose.
	Unsaid = "unsaid"
)

// Check compares a Recipe against an OpenAPI description of the same API.
//
// What this can find is worth being precise about, because the temptation is
// to read more into a green result than it holds.
//
// It can find a path that does not exist, a method that path does not take, a
// status code the operation never answers with, and a field name no schema
// declares. Those are mechanical facts, and a Recipe that gets one wrong is
// wrong in a way that a conformance suite written by the same person on the
// same day will not notice, because it will assert whatever the Recipe
// produced.
//
// It cannot find a lie. A description will happily tell you a payment has a
// status field of type string and never once mention that "approved" means
// nobody has been paid. That is what a Recipe is for, and no amount of
// checking against a schema will produce it.
//
// So a Recipe with no findings is not verified. It is un-contradicted, which
// is a smaller and more honest claim, and the report says so in those words.
func Check(r *recipe.Recipe, doc *Document, basePath string) []Finding {
	var findings []Finding

	declared := indexPaths(doc, basePath)
	seen := map[string]bool{}

	for i, route := range r.Routes {
		where := fmt.Sprintf("route %d (%s %s)", i+1, route.Method, route.Path)

		match, ok := declared.find(route.Path, route.Method, doc)
		if !ok {
			findings = append(findings, Finding{
				Where:    where,
				What:     "the description declares no such path",
				Severity: Disagrees,
			})

			continue
		}

		seen[match.template] = true

		item := doc.Paths[match.template]

		op := operationFor(item, route.Method)
		if op == nil {
			declared := methodsOf(item)

			// A description may put a path in another file. Lob does it for
			// every one of its fifty-eight paths, and those files are not
			// fetched, so the path item is empty and every route was reported
			// as a method the description does not declare -- beside an empty
			// list of the ones it supposedly does: "declares  but not GET on
			// it". That reads like a broken description rather than one only
			// half read, and it is missing evidence, not a disagreement.
			if len(declared) == 0 && item.Ref != "" {
				findings = append(findings, Finding{
					Where:    where,
					What:     fmt.Sprintf("the description puts this path in another file (%s), which is not read", item.Ref),
					Severity: Missing,
				})

				continue
			}

			if len(declared) == 0 {
				findings = append(findings, Finding{
					Where:    where,
					What:     "the description declares this path with no methods at all",
					Severity: Disagrees,
				})

				continue
			}

			findings = append(findings, Finding{
				Where:    where,
				What:     fmt.Sprintf("the description declares %s but not %s on it", strings.Join(declared, ", "), route.Method),
				Severity: Disagrees,
			})

			continue
		}

		if op.Deprecated {
			findings = append(findings, Finding{
				Where:    where,
				What:     "the description marks this operation deprecated",
				Severity: Missing,
			})
		}

		findings = append(findings, checkStatus(doc, route, op, where)...)
		findings = append(findings, checkFields(r, doc, route, op, where)...)
	}

	findings = append(findings, checkErrors(r, doc)...)

	// Paths the description has and the Recipe does not. Reported, and as
	// missing rather than as a disagreement, because modelling less than an
	// API is what every Recipe here does on purpose.
	for _, path := range doc.SortedPaths() {
		if seen[path] {
			continue
		}

		findings = append(findings, Finding{
			Where:    "recipe " + r.Name,
			What:     fmt.Sprintf("the description declares %s (%s) and no route models it", path, strings.Join(methodsOf(doc.Paths[path]), ", ")),
			Severity: Missing,
		})
	}

	return findings
}

// Paging reports the query parameters a description declares for the routes a
// Recipe models, so a paging declaration can be verified rather than guessed.
//
// This exists because 146 routes shipped declaring a paging style with no
// parameter names on it. That was harmless while the style was read by
// nothing, and became a claim the moment it was implemented: without names the
// runtime reads "limit" and the style's own word, which is right for some
// providers and wrong for plenty. Filling in sixty-one providers' names from
// memory is exactly the guessing this project refuses to do anywhere else, and
// a description states them outright.
func Paging(r *recipe.Recipe, doc *Document, basePath string) []PagingReport {
	declared := indexPaths(doc, basePath)

	var out []PagingReport

	for _, route := range r.Routes {
		if route.Operation != "list" {
			continue
		}

		match, ok := declared.find(route.Path, route.Method, doc)
		if !ok {
			out = append(out, PagingReport{Path: route.Path, Found: false})

			continue
		}

		op := operationFor(doc.Paths[match.template], route.Method)
		if op == nil {
			out = append(out, PagingReport{Path: route.Path, Found: false})

			continue
		}

		report := PagingReport{
			Path:     route.Path,
			Found:    true,
			Declared: route.Pagination,
		}

		for _, parameter := range append(append([]Parameter{}, doc.Paths[match.template].Parameters...), op.Parameters...) {
			// Resolved, because a description that declares its paging
			// parameters once and references them everywhere is the common
			// case rather than the exception. DigitalOcean's reported no
			// parameters at all until this was followed, which would have made
			// the tool answer "the description does not say" about a
			// description that says it plainly.
			parameter = doc.ResolveParameter(parameter)

			if parameter.In != "query" {
				continue
			}

			report.Query = append(report.Query, parameter.Name)
		}

		sort.Strings(report.Query)

		out = append(out, report)
	}

	return out
}

// PagingReport is what a description says about one listing's parameters.
type PagingReport struct {
	Path     string
	Found    bool
	Query    []string
	Declared recipe.Pagination
}
