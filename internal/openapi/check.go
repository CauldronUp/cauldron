package openapi

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

		match, ok := declared.find(route.Path)
		if !ok {
			findings = append(findings, Finding{
				Where:    where,
				What:     "the description declares no such path",
				Severity: Disagrees,
			})

			continue
		}

		seen[match.template] = true

		op := operationFor(doc.Paths[match.template], route.Method)
		if op == nil {
			findings = append(findings, Finding{
				Where:    where,
				What:     fmt.Sprintf("the description declares %s but not %s on it", strings.Join(methodsOf(doc.Paths[match.template]), ", "), route.Method),
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

		match, ok := declared.find(route.Path)
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
	} else {
		schema = resourceSchema(doc, r, schema, route.Resource)
		if schema == nil {
			return nil
		}
	}

	known := map[string]bool{}
	for _, property := range doc.Properties(schema) {
		known[property.Name] = true
	}

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
			What:     fmt.Sprintf("answers %d and the description declares no %d on any operation", status, status),
			Severity: Disagrees,
		})
	}

	return findings
}

// collectionSchema digs the array of objects out of a list envelope.
func collectionSchema(doc *Document, envelope *Schema, keys ...string) *Schema {
	if envelope.Type == "array" && envelope.Items != nil {
		return doc.Resolve(envelope.Items)
	}

	for _, key := range keys {
		if key == "" {
			continue
		}

		if resolved := descend(doc, envelope, key); resolved != nil {
			if resolved.Type == "array" && resolved.Items != nil {
				return doc.Resolve(resolved.Items)
			}

			return resolved
		}
	}

	// Nothing named. Fall back to the first array-of-objects property, which
	// is what a description with one collection in its envelope looks like.
	for _, property := range doc.Properties(envelope) {
		resolved := doc.Resolve(property.Schema)
		if resolved != nil && resolved.Type == "array" && resolved.Items != nil {
			return doc.Resolve(resolved.Items)
		}
	}

	return nil
}

// statusBlock reads a range like 4XX and answers the hundreds it covers.
//
// "default" is deliberately not a range. It means every status the operation
// did not name, which is a description saying it has stopped enumerating
// rather than one saying a particular failure is expected, and treating it as
// coverage would silence this check on every description that writes it.
func statusBlock(code string) (int, bool) {
	if len(code) != 3 {
		return 0, false
	}

	if !strings.EqualFold(code[1:], "XX") {
		return 0, false
	}

	n, err := strconv.Atoi(code[:1])
	if err != nil || n < 1 || n > 5 {
		return 0, false
	}

	return n, true
}

// descend walks a dotted key through a schema's properties, resolving
// references as it goes, the same way the runtime nests a dotted name.
func descend(doc *Document, schema *Schema, key string) *Schema {
	current := schema

	for _, segment := range strings.Split(key, ".") {
		current = doc.Resolve(current)
		if current == nil || current.Properties == nil {
			return nil
		}

		current = current.Properties[segment]
	}

	return doc.Resolve(current)
}

// resourceSchema descends a response envelope to the object one resource
// actually arrives in.
//
// Only listings were unwrapped, so a Recipe whose single-resource responses
// are wrapped had every field of every resource reported as undeclared: the
// comparison ran against the envelope, whose properties are things like
// result, status and time, and no resource has fields called those. Cloudflare
// wraps under result, Xero wraps under a plural and puts one object inside an
// array, and Qdrant wraps under result and then again under a name that
// differs per endpoint.
//
// A report that cries wolf is one nobody reads the second time, which is the
// same reason deletes are skipped.
func resourceSchema(doc *Document, r *recipe.Recipe, envelope *Schema, name string) *Schema {
	if r.Responses.Resource.Style != "wrapped" {
		return envelope
	}

	key := r.Responses.Resource.Key
	if key == "" {
		key = name
	}

	found := descend(doc, envelope, key)
	if found == nil {
		// The description does not nest where the Recipe says it does. That
		// is worth reporting one day and is not this function's to report,
		// and comparing against the envelope instead would turn it into one
		// finding per field.
		return envelope
	}

	if found.Type == "array" && found.Items != nil {
		return doc.Resolve(found.Items)
	}

	return found
}

// pathIndex matches a Recipe's paths against a description's templates.
type pathIndex struct {
	base      string
	templates []pathTemplate
}

type pathTemplate struct {
	template string
	pattern  *regexp.Regexp
}

type pathMatch struct {
	template string
}

func indexPaths(doc *Document, base string) pathIndex {
	index := pathIndex{base: strings.TrimRight(base, "/")}

	for _, path := range doc.SortedPaths() {
		index.templates = append(index.templates, pathTemplate{
			template: path,
			pattern:  templatePattern(index.base + path),
		})
	}

	return index
}

// templatePattern turns /v1/customers/{id} into a pattern that matches a path
// with any parameter names in it, since a Recipe names its parameters for
// itself and a description names them for its own reasons.
func templatePattern(path string) *regexp.Regexp {
	var b strings.Builder

	b.WriteString("^")

	for _, segment := range strings.Split(strings.Trim(path, "/"), "/") {
		b.WriteString("/")

		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			b.WriteString(`\{[^/}]*\}`)

			continue
		}

		// A segment that is partly literal and partly a parameter, which is
		// how Shopify writes /orders/{id}.json.
		var piece strings.Builder

		for _, part := range splitTemplate(segment) {
			if strings.HasPrefix(part, "{") {
				piece.WriteString(`\{[^/}]*\}`)

				continue
			}

			piece.WriteString(regexp.QuoteMeta(part))
		}

		b.WriteString(piece.String())
	}

	b.WriteString("/?$")

	return regexp.MustCompile(b.String())
}

// splitTemplate breaks a segment into literal and {parameter} pieces.
func splitTemplate(segment string) []string {
	var (
		out     []string
		current strings.Builder
	)

	for i := 0; i < len(segment); i++ {
		if segment[i] != '{' {
			current.WriteByte(segment[i])

			continue
		}

		close := strings.IndexByte(segment[i:], '}')
		if close < 0 {
			current.WriteByte(segment[i])

			continue
		}

		if current.Len() > 0 {
			out = append(out, current.String())
			current.Reset()
		}

		out = append(out, segment[i:i+close+1])
		i += close
	}

	if current.Len() > 0 {
		out = append(out, current.String())
	}

	return out
}

func (p pathIndex) find(path string) (pathMatch, bool) {
	normalised := "/" + strings.Trim(path, "/")

	for _, candidate := range p.templates {
		if candidate.pattern.MatchString(normalised) {
			return pathMatch{template: candidate.template}, true
		}
	}

	return pathMatch{}, false
}

func operationFor(item PathItem, method string) *Operation {
	for _, mo := range item.Operations() {
		if mo.Method == strings.ToUpper(method) {
			return mo.Operation
		}
	}

	return nil
}

func methodsOf(item PathItem) []string {
	out := make([]string, 0, 7)
	for _, mo := range item.Operations() {
		out = append(out, mo.Method)
	}

	return out
}

func sortedFieldNames(fields map[string]recipe.Field) []string {
	out := make([]string, 0, len(fields))
	for name := range fields {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

func sortedErrorNames(errors map[string]recipe.Error) []string {
	out := make([]string, 0, len(errors))
	for name := range errors {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}
