// The comparisons themselves: status codes, field names and error statuses.

package openapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// checkStatus compares a route's declared status against the description's.
func checkStatus(doc *Document, route recipe.Route, op *Operation, where string) []Finding {
	declared := route.Status
	if declared == 0 {
		switch {
		case route.Operation == "create":
			// Nothing to compare against: the Recipe is taking the runtime's
			// default rather than making a claim.
			return nil
		case route.Operation == "delete" && route.DeletedBody == "":
			declared = 204
		default:
			// Including a delete that answers with something, which the
			// runtime gives 200. Assuming 204 for every delete reported a
			// disagreement on Stripe where there was none, which is the worst
			// thing a checker can do: a report that cries wolf is one nobody
			// reads the second time.
			declared = 200
		}
	}

	_, code := doc.Success(op)
	if code == "" {
		return nil
	}

	// A description that says "2XX" is declining to be specific, and holding
	// it to a number would be inventing precision it did not offer.
	if strings.ContainsAny(code, "Xx") {
		return nil
	}

	expected, err := strconv.Atoi(code)
	if err != nil || expected == declared {
		return nil
	}

	return []Finding{{
		Where:    where,
		What:     fmt.Sprintf("answers %d and the description declares %s", declared, code),
		Severity: Disagrees,
	}}
}

// checkFields compares the wire names a resource sends against the schema of
// the operation's successful response.
func checkFields(r *recipe.Recipe, doc *Document, route recipe.Route, op *Operation, where string) []Finding {
	if route.Resource == "" {
		return nil
	}

	spec, ok := r.Resources[route.Resource]
	if !ok {
		return nil
	}

	// A delete answers with a receipt or with nothing, and neither is the
	// resource. Comparing the resource's fields against it reported every
	// field of every resource as undeclared, which is noise rather than a
	// finding.
	if route.Operation == "delete" {
		return nil
	}

	// A route that says it sends no body sends no fields, so none of them can
	// contradict anything. Ably's publish acknowledges with 204 and nothing
	// else, and every field of its message was reported as undeclared on a
	// route that emits none of them.
	if route.EmptyBody {
		return nil
	}

	schema, _ := doc.Success(op)
	if schema == nil {
		return nil
	}

	// A listing's schema is the envelope, so the object to compare against is
	// whatever the collection holds.
	if route.Operation == "list" {
		schema = collectionSchema(doc, schema, r.Responses.List.Key, spec.Collection)
		if schema == nil {
			return nil
		}

		// Each entry wrapped under the resource's own name, which is a second
		// envelope inside the first. Chargebee's listing answers
		// {"list": [{"subscription": {...}}], "next_offset": "..."}, so the
		// collection's items are objects holding the key rather than the
		// subscription itself.
		//
		// Reading an item as the resource reported every field the Recipe
		// declares as one the description does not: eleven findings with the
		// shape of a Recipe that invented its whole model, against a Recipe
		// that had the nesting exactly right and said so with entry_style.
		if r.Responses.List.EntryStyle == "wrapped" {
			key := r.Responses.List.EntryField
			if key == "" {
				key = route.Resource
			}

			if entry := descend(doc, schema, key); entry != nil {
				schema = doc.Resolve(entry)
			}
		}
	} else {
		schema = resourceSchema(doc, r, r.EnvelopeFor(route), schema, route.Resource)
		if schema == nil {
			return nil
		}
	}

	// Every name any variant declares, not only the ones allOf merges.
	//
	// Properties merges allOf and deliberately leaves oneOf alone, because
	// merging alternatives into one object describes something that does not
	// exist. The question here is narrower -- whether any schema declares
	// this name at all -- and for that the union is the right answer. A DNS
	// record is an A record or a CNAME or an MX; Cloudflare nests that three
	// deep, an allOf of an anyOf of two oneOfs with twenty-one variants
	// between them, and every field of its record was reported as undeclared
	// while all twenty-one declared it.
	known := declaredNames(doc, schema, map[*Schema]bool{})

	// A schema with no properties at all is a description declining to say,
	// not a description saying the object is empty.
	if len(known) == 0 {
		return nil
	}

	var findings []Finding

	// A route that declares returns answers with less than the record holds,
	// and comparing the whole record against that route's schema reports the
	// difference the declaration exists to describe. Qdrant's collection
	// listing sends a name and nothing else, which is the point of the
	// listing and would otherwise be eight findings.
	emitted := sortedFieldNames(spec.Fields)

	if len(route.Returns) > 0 {
		named := map[string]bool{}
		for _, name := range route.Returns {
			named[name] = true
		}

		kept := emitted[:0:0]

		for _, field := range emitted {
			if named[field] {
				kept = append(kept, field)
			}
		}

		emitted = kept
	}

	for _, field := range emitted {
		f := spec.Fields[field]

		// A field the wire never carries cannot contradict a description of
		// what the wire carries. "in: -" is the Recipe saying the record
		// holds it and the response does not -- a partition that lives in
		// the path, like Attio's object slug or Fly's app name. Reporting it
		// put a finding against precisely the Recipes that took the trouble
		// to say so.
		if f.In == "-" {
			continue
		}

		// Only the outermost name can be compared without walking the schema
		// the same way the runtime walks the record, which is more machinery
		// than the value justifies. A nested field is checked by its parent.
		name := f.WireName(field)
		if f.In != "" {
			name, _, _ = strings.Cut(f.In, ".")
			name, _, _ = strings.Cut(name, "[")
		}

		if known[name] {
			continue
		}

		findings = append(findings, Finding{
			Where:    fmt.Sprintf("resource %q field %q", route.Resource, field),
			What:     fmt.Sprintf("is sent as %q and no property of that name is declared for %s %s", name, route.Method, route.Path),
			Severity: Disagrees,
		})
	}

	return findings
}

// checkErrors compares the Recipe's error table against the status codes the
// description declares anywhere.
func checkErrors(r *recipe.Recipe, doc *Document) []Finding {
	declared := map[int]bool{}
	// Ranges, keyed by their hundreds block, for the descriptions that decline
	// to be specific about failures.
	ranges := map[int]bool{}

	operations, documented := 0, 0

	for _, item := range doc.Paths {
		for _, mo := range item.Operations() {
			operations++

			failures := doc.Failures(mo.Operation)
			if len(failures) > 0 {
				documented++
			}

			for _, code := range failures {
				if n, err := strconv.Atoi(code); err == nil {
					declared[n] = true

					continue
				}

				// A description may decline to be specific and declare a
				// range instead: 4XX covers every client error, which is what
				// Qdrant writes on every operation it has. Reading only the
				// numeric codes beside it left 503 as the entire set of known
				// failures, so a 400, a 401 and a 404 were each reported as
				// undeclared against a description that had already said all
				// three were expected.
				if block, ok := statusBlock(code); ok {
					ranges[block] = true
				}
			}
		}
	}

	// A description that names failures on almost none of its operations is
	// not describing failures, it is describing successes. Twilio's declares a
	// numeric 4xx on two operations out of a hundred and ninety-seven, and
	// holding an error table to those two reported every entry in it as
	// wrong — six findings, all noise, on a Recipe with nothing wrong.
	//
	// Half is a judgement rather than a fact, and it is here so the check
	// stays silent when it has nothing to say rather than filling a report
	// nobody will read twice.
	if (len(declared) == 0 && len(ranges) == 0) || operations == 0 || documented*2 < operations {
		return nil
	}

	var findings []Finding

	for _, name := range sortedErrorNames(r.Errors) {
		status := r.Errors[name].Status
		if status == 0 || declared[status] || ranges[status/100] {
			continue
		}

		findings = append(findings, Finding{
			Where:    fmt.Sprintf("error %q", name),
			What:     fmt.Sprintf("answers %d and the description declares no %d anywhere, so it neither supports nor contradicts this", status, status),
			Severity: Unsaid,
		})
	}

	return findings
}
