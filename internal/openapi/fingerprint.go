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
	claims := claimsFor(r, doc, basePath)

	texts := make([]string, 0, len(claims))
	for _, c := range claims {
		texts = append(texts, c.text)
	}

	sort.Strings(texts)

	sum := sha256.Sum256([]byte(strings.Join(texts, "\n")))

	return hex.EncodeToString(sum[:])
}

// Unbacked names the claims a Recipe makes that the description does not
// support: a path it does not declare, a field its success schema does not
// have, a status none of the Recipe's operations answers with.
//
// It exists because "MOVED" on its own is not actionable. The fingerprint is a
// hash of the claim list below, so a moved hash says that something in that
// list changed and nothing about what -- and the cheapest response to a line
// nobody can read is to re-record it, which is precisely the reflex that turns
// this check into a weekly rubber stamp. An absence is the half a reader can
// act on, and it is computable from the document in hand without keeping any
// history of previous ones.
//
// This is deliberately not the same question as "did it move", and the two
// answers can disagree in both directions. A gap can be old: a Recipe pinned to
// one of several reference files a provider splits its API across has had a
// route missing from that document since the day it was written, and its
// fingerprint has been stable throughout. A move can have no gap at all: a
// provider adding a 429 to an operation moves the hash and takes nothing away.
// Both are worth printing, and only at the moment somebody is already being
// asked to look.
func Unbacked(r *recipe.Recipe, doc *Document, basePath string) []string {
	var gaps []string

	for _, c := range claimsFor(r, doc, basePath) {
		if c.gap != "" {
			gaps = append(gaps, c.gap)
		}
	}

	sort.Strings(gaps)

	return gaps
}

// claim is one line of the fingerprint, and what is missing when the
// description does not back it.
//
// The text is carried rather than reconstructed because the recorded
// fingerprints of every Recipe that names a description are hashes of these
// exact strings. Recovering the gaps by parsing them back out would work until
// somebody reworded one, and rewording one moves every recorded hash in the
// collection at once.
type claim struct {
	text string
	// gap is empty when the description backs the claim, and a sentence
	// naming what is missing when it does not.
	gap string
}

func claimsFor(r *recipe.Recipe, doc *Document, basePath string) []claim {
	index := indexPaths(doc, basePath)

	var claims []claim

	for _, route := range r.Routes {
		claims = append(claims, routeClaims(doc, index, r, route)...)
	}

	for _, name := range sortedErrorNames(r.Errors) {
		// A status is claimed by the Recipe rather than by one route, and the
		// description declares them per operation, so the claim recorded is
		// whether any operation the Recipe routes to answers with it.
		status := strconv.Itoa(r.Errors[name].Status)
		answered := anyOperationAnswers(doc, index, r, status)

		gap := ""
		if !answered {
			// Named, because two of a Recipe's errors can share a status and
			// the claim texts for those are identical by design. Without the
			// name the report prints the same sentence twice and reads like a
			// bug in the reporter.
			gap = "no operation the Recipe routes to answers " + status + ", which " + name + " declares"
		}

		claims = append(claims, claim{
			text: "status " + status + " " + strconv.FormatBool(answered),
			gap:  gap,
		})
	}

	return claims
}

// routeClaims renders one route's share of the fingerprint.
func routeClaims(doc *Document, index pathIndex, r *recipe.Recipe, route recipe.Route) []claim {
	method := routeMethod(route)

	match, ok := index.find(route.Path, method, doc)
	if !ok {
		// A path the description does not have is itself a claim, and one
		// that has to survive the provider adding it later, so the absence is
		// recorded rather than skipped.
		return []claim{{
			text: "absent " + method + " " + route.Path,
			gap:  "the description does not declare " + method + " " + route.Path,
		}}
	}

	op := operationFor(doc.Paths[match.template], method)
	if op == nil {
		return []claim{{
			text: "absent " + method + " " + match.template,
			gap:  "the description does not declare " + method + " " + match.template,
		}}
	}

	prefix := method + " " + match.template

	claims := []claim{{text: prefix}}

	success, _ := doc.Success(op)

	declared := map[string]string{}
	for _, property := range doc.Properties(success) {
		declared[property.Name] = string(property.Schema.Type)
	}

	// Only the fields the Recipe names. A schema growing a sibling is not a
	// contradiction of anything the Recipe says, and calling it one is the
	// same noise as checksumming the file. An empty type is the field being
	// absent from the schema, which is a claim worth moving on.
	fields := sortedFieldNames(r.Resources[route.Resource].Fields)

	// Whether this reading got inside the response at all. Document.Properties
	// merges allOf and stops at the top level, so for a provider that wraps
	// its records -- {"data": {...}} -- the only property here is data, every
	// field the Recipe names is absent, and none of those absences is a fact
	// about the provider. Nine shipped Recipes are in that position, and the
	// header of this file names them.
	//
	// The claims are still recorded, because a Recipe adding or removing a
	// field must still move its fingerprint. What is suppressed is calling
	// them gaps: a reader told that Vercel's description declares none of the
	// forty fields across its seven routes learns only that this parser did
	// not open the envelope, and a report that says that every week is a
	// report nobody finishes reading.
	inside := false

	for _, field := range fields {
		if _, named := declared[field]; named {
			inside = true

			break
		}
	}

	for _, field := range fields {
		kind, named := declared[field]

		gap := ""
		if !named && inside {
			gap = prefix + ": the success schema does not declare " + field
		}

		claims = append(claims, claim{text: prefix + " field " + field + " " + kind, gap: gap})
	}

	for _, code := range sortedCodes(op.Responses) {
		claims = append(claims, claim{text: prefix + " answers " + code})
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
