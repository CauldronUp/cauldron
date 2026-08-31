// Adding a provider's own description to a Recipe, without letting it
// overwrite one.

package recipe

import (
	"fmt"
	"sort"
	"strings"
)

// AugmentReport is what happened, in numbers a person can read out.
//
// It exists because the merged Recipe cannot be allowed to look like a written
// one. A user has to be able to see, without asking, how much of what they are
// talking to was observed and how much was inferred from a schema.
type AugmentReport struct {
	// Kept is routes the Recipe already had, which the description did not
	// replace.
	Kept int
	// Added is routes the description had and the Recipe did not.
	Added int
	// Skipped is routes the description had that could not be added safely.
	Skipped int
	// Reasons says why, one line per skipped route, sorted so two runs of the
	// same inputs read the same.
	Reasons []string
}

// Augment returns the Recipe with routes from a description added for paths the
// Recipe does not model. The Recipe always wins.
//
// The bargain is worth stating plainly, because the temptation is to read a
// merged Recipe as a better one.
//
// A description carries breadth. Box declares 186 paths where its Recipe models
// 7; OpenAI declares 182 where its Recipe models 5. Hand-writing the remainder
// will never happen, and a caller who only needs an endpoint to exist and
// answer something plausibly shaped is well served by a schema.
//
// A description does not carry truth, and cannot. Kraken's says error is an
// array of string; there is no way for it to say that the array is empty on
// success and that an empty array is true in every language a client is written
// in. Deezer's declares 200: Artist; it cannot declare that every failure is
// also a 200. Asana's omits completed and due_on from the task listing
// altogether -- they are opt_fields, real and undeclared -- so a mock built
// from that description serves tasks without them, and the Recipe that pins
// them is right exactly where the description is wrong.
//
// So the rule is one-directional: the description may add what the Recipe is
// silent about, and may never change what the Recipe says. Everything it adds
// is marked Derived, counted, and reported.
func Augment(r, extra *Recipe) (*Recipe, AugmentReport) {
	merged := *r

	var report AugmentReport

	if extra == nil || len(extra.Routes) == 0 {
		return &merged, report
	}

	have := map[string]bool{}
	for _, route := range r.Routes {
		have[routeKey(route)] = true
	}

	// Copied rather than shared, so adding to the merged Recipe cannot reach
	// back into the one that was passed in.
	merged.Routes = append([]Route(nil), r.Routes...)
	merged.Resources = copyResources(r.Resources)
	merged.Errors = copyErrors(r.Errors)

	// What the written Recipe calls its collections, when it is unanimous.
	shared := sharedCollection(r)

	for _, route := range extra.Routes {
		if have[routeKey(route)] {
			report.Kept++

			continue
		}

		// A resource name that means one thing here and another there cannot
		// be reconciled, and guessing which was meant would put a field on the
		// wire that no provider sends.
		if _, clash := r.Resources[route.Resource]; clash {
			report.Skipped++
			report.Reasons = append(report.Reasons, fmt.Sprintf(
				"%s %s: the description's resource %q is not the one this Recipe means by that name",
				route.Method, route.Path, route.Resource))

			continue
		}

		if resource, ok := extra.Resources[route.Resource]; ok {
			// A drafted resource names no collection, and the Recipe it is
			// joining may wrap its listings without a key -- in which case
			// every resource owes one. The name is read from the path rather
			// than guessed from the resource, because the path is what the
			// provider itself calls the collection.
			if resource.Collection == "" {
				resource.Collection = shared
			}

			if resource.Collection == "" {
				resource.Collection = collectionFromPath(route.Path)
			}

			merged.Resources[route.Resource] = resource
		} else {
			report.Skipped++
			report.Reasons = append(report.Reasons, fmt.Sprintf(
				"%s %s: the description declares no resource %q to answer with",
				route.Method, route.Path, route.Resource))

			continue
		}

		route.Derived = true
		merged.Routes = append(merged.Routes, route)
		report.Added++
	}

	// Errors are added only where the Recipe is silent. An error is a claim
	// about what a failure looks like, and the Recipe's claim was checked
	// against the provider where the description's was not.
	for name, err := range extra.Errors {
		if _, taken := merged.Errors[name]; !taken {
			merged.Errors[name] = err
		}
	}

	sort.Strings(report.Reasons)

	if report.Added > 0 {
		merged.DerivedFrom = extra.Upstream.Spec
	}

	// The guarantee the whole feature rests on: reaching for a description
	// can add to a Recipe and can never break one. A Recipe that will not
	// validate will not serve, so a merge that produced one would leave a
	// caller who asked for more endpoints with none at all.
	//
	// Everything above tries to make the addition work. This is what happens
	// when it could not, and it is deliberately blunt: the written Recipe,
	// unchanged, and a reason.
	if err := merged.Validate(); err != nil {
		return r, AugmentReport{
			Kept:    len(r.Routes),
			Skipped: report.Added + report.Skipped,
			Reasons: []string{
				"the merged Recipe would not validate, so nothing was added: " + err.Error(),
			},
		}
	}

	return &merged, report
}

// routeKey is the pair a route is identified by for this purpose: two routes
// on the same method and path are the same endpoint, whatever else they say.
func routeKey(route Route) string {
	method := strings.ToUpper(route.Method)
	if method == "" {
		method = "GET"
	}

	return method + " " + route.Path
}

// sharedCollection is the collection name every written resource uses, when
// they all use the same one, and empty otherwise.
//
// Asana wraps every listing in "data" and declares no list key, so the
// wrapper is each resource's collection name and all three of its written
// resources say "data". A derived resource naming itself from its path
// therefore served {"workspaces": []} where Asana sends {"data": [...]} --
// a wrong envelope rather than a missing one, which is worse for anything
// that parses.
//
// Unanimity is the whole condition. Where every observed resource agrees,
// that agreement is evidence about the provider and a derived resource
// should follow it. Where they disagree there is nothing to follow, and
// picking one would be inventing a convention the Recipe never claimed.
func sharedCollection(r *Recipe) string {
	// One resource is unanimous with itself, which says nothing. A house
	// convention needs at least two resources to have agreed on it.
	if len(r.Resources) < 2 {
		return ""
	}

	var shared string

	for _, resource := range r.Resources {
		if resource.Collection == "" {
			return ""
		}

		if shared == "" {
			shared = resource.Collection

			continue
		}

		if resource.Collection != shared {
			return ""
		}
	}

	return shared
}

// collectionFromPath is the last literal segment of a path, which is the
// name the provider gives the collection.
func collectionFromPath(path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")

	for i := len(segments) - 1; i >= 0; i-- {
		if segment := segments[i]; segment != "" && !strings.HasPrefix(segment, "{") {
			return segment
		}
	}

	return ""
}

func copyResources(in map[string]Resource) map[string]Resource {
	out := make(map[string]Resource, len(in))
	for name, resource := range in {
		out[name] = resource
	}

	return out
}

func copyErrors(in map[string]Error) map[string]Error {
	out := make(map[string]Error, len(in))
	for name, err := range in {
		out[name] = err
	}

	return out
}
