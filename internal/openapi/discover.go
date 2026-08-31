// Finding the description a provider publishes, and proving it is theirs.

package openapi

import (
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Candidate is a description that might be a Recipe's, with the evidence for
// saying so.
//
// Every field here exists to be read by a person deciding whether to record
// the URL. Matched out of Routes is the case for it; Declared is the case
// against, because a document declaring two paths that happens to cover one
// route is a documentation platform's boilerplate rather than an API.
type Candidate struct {
	// URL is where the description was served from.
	URL string
	// Format is what it turned out to be, in the provider's own version.
	Format string
	// Declared is how many paths the description carries.
	Declared int
	// Matched is how many of the Recipe's routes it declares.
	Matched int
	// Routes is how many routes the Recipe has.
	Routes int
}

// wellKnown is where a description is served when nobody has said where.
//
// This list is the whole reason the search is cheap enough to run across a
// collection, and the whole reason it will never be complete. A provider that
// publishes its description in a GitHub repository -- which is how eleven of
// the first twelve here were found -- is not reachable this way and needs a
// person to say where it is.
var wellKnown = []string{
	"/openapi.json",
	"/openapi.yaml",
	"/swagger.json",
	"/swagger.yaml",
	"/.well-known/openapi.json",
	"/api-docs",
	"/api-docs.json",
	"/v1/openapi.json",
	"/api/openapi.json",
	"/api/v1/openapi.json",
	"/docs/openapi.json",
	"/openapi/v1.json",
	"/swagger/v1/swagger.json",
	"/spec.json",
}

// Discover looks for a description of the API a Recipe models, and returns one
// only where the document proves it is describing that API.
//
// The proof is the point, and it is here because guessing was tried first and
// produced nonsense twice. Matching a Recipe to a description by domain scored
// ninety-eight hits across this collection, of which a dozen were Recipes whose
// documentation happens to live on github.com being matched to GitHub's own
// description; matching on a directory's declared hosts instead paired Bugsnag
// with ClickSend, because both record an origin on a shared documentation host.
// A URL resolving says nothing. A 200 says nothing. What says something is the
// document declaring a path the Recipe already models, and nothing weaker is
// accepted here.
//
// Nothing is written. A Candidate is a proposal for a person to confirm,
// because recording the wrong URL does not fail loudly -- it reports drift
// against a document that was never this provider's, every day, until somebody
// works out why.
func Discover(r *recipe.Recipe, hosts []string, fetch func(string) ([]byte, error)) *Candidate {
	// With no routes there is nothing to prove a match against, and a
	// proposal on no evidence is the guessing this exists to avoid.
	if r == nil || len(r.Routes) == 0 {
		return nil
	}

	var best *Candidate

	for _, host := range searchHosts(hosts) {
		for _, path := range wellKnown {
			url := "https://" + strings.TrimSuffix(host, "/") + path

			raw, err := fetch(url)
			if err != nil || len(raw) == 0 {
				continue
			}

			doc, err := Parse(raw)
			if err != nil {
				continue
			}

			matched := MatchedRoutes(r, doc)
			if matched == 0 {
				continue
			}

			candidate := &Candidate{
				URL:      url,
				Format:   formatOf(doc),
				Declared: len(doc.Paths),
				Matched:  matched,
				Routes:   len(r.Routes),
			}

			// More of the Recipe accounted for is a better claim to being its
			// description. Where two account for the same, the smaller
			// document is the more specific one.
			if best == nil || candidate.Matched > best.Matched ||
				(candidate.Matched == best.Matched && candidate.Declared < best.Declared) {
				best = candidate
			}

			// A document that accounts for every route cannot be beaten,
			// and every further request is one somebody else's server did
			// not need to answer.
			if best.Matched == best.Routes {
				return best
			}
		}
	}

	return best
}

// searchHosts widens the verified hosts to the siblings a description is
// commonly published on, in order, without repeats.
//
// Attio is verified against api.attio.com and publishes at attio.com; Monday
// is verified against api.monday.com and publishes at monday.com. The host a
// response was recorded from names the organisation, and the description
// frequently lives on a sibling of it rather than on the API itself.
//
// On its own this would be reckless -- it is one step from the domain
// matching that paired Bugsnag with ClickSend, because both record an origin
// on a shared documentation host. It is safe only because of what happens to
// whatever it finds: a sibling's document still has to declare a path the
// Recipe models, or it is not proposed. Widening the search is free when the
// test at the end is on the document rather than on the address.
func searchHosts(verified []string) []string {
	var (
		out  []string
		seen = map[string]bool{}
	)

	add := func(host string) {
		if host != "" && !seen[host] {
			seen[host] = true
			out = append(out, host)
		}
	}

	// Every verified host first: where a provider serves its description
	// from the API itself, that is the answer and no sibling is fetched.
	for _, host := range verified {
		add(host)
	}

	for _, host := range verified {
		domain := registeredDomain(host)

		add(domain)
		add("developers." + domain)
		add("docs." + domain)
		add("developer." + domain)
	}

	return out
}

// registeredDomain is the domain a host belongs to.
//
// The two-label public suffixes are listed rather than guessed at because
// treating co.uk as the domain would widen the search to every API in the
// United Kingdom. This is not the public suffix list and does not need to be:
// a suffix it gets wrong costs a handful of requests that find nothing, since
// nothing is proposed without the proof.
func registeredDomain(host string) string {
	labels := strings.Split(strings.ToLower(host), ".")

	keep := 2
	if len(labels) >= 3 {
		switch labels[len(labels)-2] {
		case "co", "com", "org", "net", "gov", "ac", "edu":
			keep = 3
		}
	}

	if len(labels) <= keep {
		return strings.Join(labels, ".")
	}

	return strings.Join(labels[len(labels)-keep:], ".")
}

// MatchedRoutes is how many of a Recipe's routes a description declares.
//
// It uses the same index Check does, which matters more than it sounds: a
// description declares its paths relative to its own server, so /widgets under
// a server of /v2 is the Recipe's /v2/widgets. Comparing the strings would
// report a real description as describing none of this API -- the mistake this
// package has now made twice, once in the fingerprint and once in the routes
// serve --with-spec derives.
func MatchedRoutes(r *recipe.Recipe, doc *Document) int {
	declared := indexPaths(doc, BasePath(doc))

	var matched int

	for _, route := range r.Routes {
		if _, ok := declared.find(route.Path, route.Method, doc); ok {
			matched++
		}
	}

	return matched
}

func formatOf(doc *Document) string {
	if doc.Swagger != "" {
		return "Swagger " + doc.Swagger
	}

	return "OpenAPI " + doc.OpenAPI
}

// Hosts is where a Recipe was actually seen talking, most-used first.
//
// Only the sources of conformance cases are read. A Recipe's documentation
// link is not the API: reading every URL in the file is what paired a dozen
// Recipes with GitHub's description, because that is where their documentation
// lives. A case source is an address a response was recorded from, so it is
// the API by construction.
func Hosts(r *recipe.Recipe) []string {
	counts := map[string]int{}

	for _, c := range r.Conformance {
		if host := hostOf(c.Source); host != "" {
			counts[host]++
		}
	}

	hosts := make([]string, 0, len(counts))
	for host := range counts {
		hosts = append(hosts, host)
	}

	// Most-seen first, and alphabetically within a tie so two runs of the same
	// Recipe try the same host first.
	sort.SliceStable(hosts, func(i, j int) bool {
		if counts[hosts[i]] != counts[hosts[j]] {
			return counts[hosts[i]] > counts[hosts[j]]
		}

		return hosts[i] < hosts[j]
	})

	return hosts
}

func hostOf(raw string) string {
	rest, ok := strings.CutPrefix(raw, "https://")
	if !ok {
		if rest, ok = strings.CutPrefix(raw, "http://"); !ok {
			return ""
		}
	}

	if i := strings.IndexAny(rest, "/?#"); i >= 0 {
		rest = rest[:i]
	}

	return strings.ToLower(rest)
}
