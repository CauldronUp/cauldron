package openapi

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Fingerprint reduces a provider's own description to the parts one Recipe
// makes claims about, and hashes those.
//
// The obvious thing is to checksum the file, and the obvious thing is wrong.
// Providers republish these documents constantly: a reworded summary, a new
// example, an endpoint in a product this Recipe has never heard of. A checksum
// over the file calls every one of those drift, a scan that reports drift on
// every publish gets switched off, and then the change that mattered arrives
// unannounced. The noisy check and no check at all fail the same way, and the
// noisy one costs more on the way there.
//
// So the fingerprint covers the intersection and nothing else: the paths and
// methods the Recipe declares routes for, the response codes those operations
// answer with, and the types of the field names the Recipe itself names. That
// is the surface where a provider's change and a Recipe's claim can
// contradict each other.
//
// It is worth saying what this does not do, because the temptation is to read
// a stable fingerprint as a Recipe still being right. It is not. A description
// can say a payment has a status of type string on the day the provider starts
// answering "approved" for payments nobody was paid for, and the fingerprint
// will not move, because nothing the document says has changed. What it
// catches is the mechanical half -- a field renamed, a path moved, a status
// dropped -- and the mechanical half is precisely what a Recipe's own
// conformance suite is structurally unable to catch, because that suite
// asserts what the Recipe says rather than what the provider does.
//
// A Recipe whose fingerprint has not moved is un-contradicted by the
// description, on the parts it claims. That is a smaller sentence than
// "unchanged", and it is the true one.
//
// It is smaller still for a provider that wraps its responses. The field half
// of the fingerprint reads the success schema's own top-level properties --
// Document.Properties merges allOf and stops there, and never descends into a
// property's own schema. An API answering {"data": {...}} therefore has one
// top-level property called data, and every field name the Recipe declares is
// recorded as absent rather than as a type. The claims are still made and
// still compared, so a Recipe adding or removing a field still moves its
// fingerprint; what cannot move it is the provider changing one of those
// fields' types, because the type was never read.
//
// Nine shipped Recipes are in that position: Attio, Vercel, Turso, Temporal,
// ClickHouse, NocoDB, Scout APM, Livepeer and MX. For those, drift compares
// paths, methods and response codes and nothing else, which is a real check
// and a narrower one than the paragraph above promises. Found while
// establishing why Attio moved, where the answer turned out to be a response
// code, and where it could not have been a field type whatever the provider
// did.
//
// Teaching this to follow an envelope is the obvious fix and is not free: it
// would move the recorded fingerprint of every one of those nine at once, and
// a batch of moves nobody can attribute to a provider is the same
// switched-off-scanner failure the paragraph above is about. Worth doing
// deliberately, on its own, rather than as a side effect.
func Fingerprint(r *recipe.Recipe, doc *Document, basePath string) string {
	index := indexPaths(doc, basePath)

	var claims []string

	for _, route := range r.Routes {
		claims = append(claims, routeClaims(doc, index, r, route)...)
	}

	for _, name := range sortedErrorNames(r.Errors) {
		// A status is claimed by the Recipe rather than by one route, and the
		// description declares them per operation, so the claim recorded is
		// whether any operation the Recipe routes to answers with it.
		status := strconv.Itoa(r.Errors[name].Status)
		claims = append(claims, "status "+status+" "+strconv.FormatBool(anyOperationAnswers(doc, index, r, status)))
	}

	sort.Strings(claims)

	sum := sha256.Sum256([]byte(strings.Join(claims, "\n")))

	return hex.EncodeToString(sum[:])
}

// routeClaims renders one route's share of the fingerprint.
func routeClaims(doc *Document, index pathIndex, r *recipe.Recipe, route recipe.Route) []string {
	method := routeMethod(route)

	match, ok := index.find(route.Path, method, doc)
	if !ok {
		// A path the description does not have is itself a claim, and one
		// that has to survive the provider adding it later, so the absence is
		// recorded rather than skipped.
		return []string{"absent " + method + " " + route.Path}
	}

	op := operationFor(doc.Paths[match.template], method)
	if op == nil {
		return []string{"absent " + method + " " + match.template}
	}

	prefix := method + " " + match.template

	claims := []string{prefix}

	success, _ := doc.Success(op)

	declared := map[string]string{}
	for _, property := range doc.Properties(success) {
		declared[property.Name] = string(property.Schema.Type)
	}

	// Only the fields the Recipe names. A schema growing a sibling is not a
	// contradiction of anything the Recipe says, and calling it one is the
	// same noise as checksumming the file. An empty type is the field being
	// absent from the schema, which is a claim worth moving on.
	for _, field := range sortedFieldNames(r.Resources[route.Resource].Fields) {
		claims = append(claims, prefix+" field "+field+" "+declared[field])
	}

	for _, code := range sortedCodes(op.Responses) {
		claims = append(claims, prefix+" answers "+code)
	}

	return claims
}

// anyOperationAnswers reports whether an operation the Recipe routes to
// declares this status.
func anyOperationAnswers(doc *Document, index pathIndex, r *recipe.Recipe, status string) bool {
	for _, route := range r.Routes {
		method := routeMethod(route)

		match, ok := index.find(route.Path, method, doc)
		if !ok {
			continue
		}

		op := operationFor(doc.Paths[match.template], method)
		if op == nil {
			continue
		}

		if _, declared := op.Responses[status]; declared {
			return true
		}
	}

	return false
}

// routeMethod is the route's method, upper-cased, defaulting to GET the way
// the rest of the package reads it.
func routeMethod(route recipe.Route) string {
	if route.Method == "" {
		return "GET"
	}

	return strings.ToUpper(route.Method)
}
