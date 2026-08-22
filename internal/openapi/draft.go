package openapi

import (
	"fmt"
	"sort"
	"strings"
)

// Draft renders a starting point for a Recipe from an OpenAPI description.
//
// What comes back is a draft and the file says so at the top, in the first
// line, where nobody can miss it. That is not modesty. A Recipe earns its
// place by reproducing behaviour worth reproducing — a third state nobody
// branches on, a field whose absence is the signal, a success that is not a
// success — and a description contains none of that. It contains paths, names
// and types.
//
// So this removes the typing and leaves the thinking, and it deliberately
// produces something that cannot be shipped by accident: no conformance cases,
// no fixture beyond the empty one, and a list at the top of what a person has
// to decide. A generated Recipe with a green suite and nothing in it would be
// exactly the confident, empty artefact this project exists to argue against.
func Draft(doc *Document, name string) string {
	var b strings.Builder

	writeDraftHeader(&b, doc, name)

	fmt.Fprintf(&b, "recipe: %s\n", name)

	guessed, confident := guessCapability(doc)

	if confident {
		fmt.Fprintf(&b, "# Guessed from the description's title. Check it.\ncapability: %s\n", guessed)
	} else {
		fmt.Fprintf(&b, `# A guess, and a poor one: nothing in the description said. This is one of
# the things at the top of this file a person has to decide, and it decides
# whether anybody looking for this provider by what it does will find it.
capability: %s
`, guessed)
	}

	b.WriteString("version: 0.1.0\n\n")

	b.WriteString("upstream:\n")
	fmt.Fprintf(&b, "  api: %q\n", doc.Info.Version)

	if docs := documentationURL(doc); docs != "" {
		fmt.Fprintf(&b, "  docs: %s\n", docs)
	}

	b.WriteString("\n")

	writeDraftAuth(&b, doc)
	writeDraftResponses(&b)

	resources := deriveResources(doc)
	writeDraftResources(&b, doc, resources)
	writeDraftRoutes(&b, doc, resources)
	writeDraftErrors(&b, doc)

	b.WriteString("fixtures:\n  empty: {}\n\n")
	b.WriteString(`# No fixture and no conformance case is generated, on purpose.
#
# A fixture is a claim about what a plausible account looks like, and a case is
# a claim about what the provider does. Neither is in the description this was
# built from, and generating a case that asserts whatever the emulator happens
# to produce would be worse than generating none: it would pass, for ever,
# while proving nothing at all.
conformance: []
`)

	return b.String()
}

func writeDraftHeader(b *strings.Builder, doc *Document, name string) {
	title := doc.Info.Title
	if title == "" {
		title = name
	}

	fmt.Fprintf(b, `# DRAFT %s Recipe, generated from an OpenAPI description.
#
# This is not a Recipe yet. It is the mechanical half of one: the paths, the
# field names, the types and the status codes, which a description does carry
# and which are tedious to type. Everything that makes a Recipe worth having
# is missing, because a description does not contain it.
#
# Before this ships, somebody has to decide:
#
#   1. What does this API lie about? A status that reads like success and is
#      not, a field whose absence is the signal, an amount in minor units with
#      nothing saying so, a list that returns the series rather than the
#      occurrence. That is the reason a Recipe exists, and the reason to keep
#      or discard this one.
#   2. Which endpoints must not be here at all? Anything that moves money,
#      sends a message, employs somebody or deletes something permanently is
#      left out of the Recipes in this collection on purpose, and the header
#      says why. Check what has been generated below against that.
#   3. What does a plausible account look like? The fixture below is empty. A
#      fixture carrying one of each interesting state is most of what makes
#      the emulator useful.
#   4. What is not modelled? Every Recipe here states its gaps in its own
#      header rather than leaving them to be discovered. This one has none
#      written yet, and it certainly has gaps.
#
# Then delete this block, write the header properly, and write the cases.
#
# Generated from: %s%s
#
`, title, title, versionSuffix(doc))
}

func versionSuffix(doc *Document) string {
	if doc.Info.Version == "" {
		return ""
	}

	return " " + doc.Info.Version
}

func documentationURL(doc *Document) string {
	if doc.ExternalDocs != nil && doc.ExternalDocs.URL != "" {
		return doc.ExternalDocs.URL
	}

	if len(doc.Servers) > 0 {
		return doc.Servers[0].URL
	}

	return ""
}

func writeDraftAuth(b *strings.Builder, doc *Document) {
	b.WriteString("auth:\n")

	names := make([]string, 0, len(doc.Components.SecuritySchemes))
	for name := range doc.Components.SecuritySchemes {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		scheme := doc.Components.SecuritySchemes[name]

		switch {
		case scheme.Type == "http" && strings.EqualFold(scheme.Scheme, "bearer"):
			fmt.Fprintf(b, "  # From the %q security scheme.\n", name)
			b.WriteString("  scheme: bearer\n  header: Authorization\n  prefix: \"Bearer \"\n")
			b.WriteString("  keys:\n    - cauldron-token\n\n")

			return
		case scheme.Type == "http" && strings.EqualFold(scheme.Scheme, "basic"):
			fmt.Fprintf(b, "  # From the %q security scheme. Which half carries the secret is\n", name)
			b.WriteString("  # not in the description, so credential is a guess worth checking.\n")
			b.WriteString("  scheme: basic\n  credential: password\n  keys:\n    - cauldron-token\n\n")

			return
		case scheme.Type == "apiKey" && scheme.In == "header":
			fmt.Fprintf(b, "  # From the %q security scheme.\n", name)
			fmt.Fprintf(b, "  scheme: header\n  header: %s\n  keys:\n    - cauldron-token\n\n", scheme.Name)

			return
		case scheme.Type == "apiKey" && scheme.In == "query":
			fmt.Fprintf(b, "  # From the %q security scheme. A credential in the URL ends up in\n", name)
			b.WriteString("  # access logs and browser history, and reproducing that is the point.\n")
			fmt.Fprintf(b, "  scheme: query\n  param: %s\n  keys:\n    - cauldron-token\n\n", scheme.Name)

			return
		}
	}

	b.WriteString(`  # The description declares no scheme this can read, which usually means
  # OAuth. Cauldron does not model a token exchange, so pick the shape the
  # API actually accepts on a request and say so here.
  scheme: bearer
  header: Authorization
  prefix: "Bearer "
  keys:
    - cauldron-token

`)
}

func writeDraftResponses(b *strings.Builder) {
	b.WriteString(`responses:
  # Guessed, and worth checking against a real response before trusting.
  # Whether a list is bare, wrapped or enveloped is not something a schema
  # says plainly, and getting it wrong is the difference between a client that
  # works and one that reads undefined.
  list:
    style: bare
  resource:
    style: bare
  error:
    style: flat
    # message_field is deliberately not set here, and it is one of the things
    # to decide before this ships. A Recipe naming where a failure's prose
    # lives has to have a case asserting that name, or it could be called
    # anything -- so a draft that guessed "message" would arrive already
    # refused by the validator, and setting it means writing the case beside
    # it.

`)
}

// draftResource is one resource the draft will declare.
type draftResource struct {
	name   string
	schema *Schema
}

// deriveResources picks a resource per schema that a successful response
// returns, named after the schema rather than after the path, because that is
// the name the description's author chose for the thing.
func deriveResources(doc *Document) map[string]draftResource {
	out := map[string]draftResource{}

	for _, path := range doc.SortedPaths() {
		for _, mo := range doc.Paths[path].Operations() {
			schema, _ := doc.Success(mo.Operation)
			if schema == nil {
				continue
			}

			object := schema
			if operationKind(path, mo.Method) == "list" {
				object = collectionSchema(doc, schema)
				if object == nil {
					continue
				}
			}

			if len(doc.Properties(object)) == 0 {
				continue
			}

			name := resourceName(path)
			if _, taken := out[name]; taken {
				continue
			}

			out[name] = draftResource{name: name, schema: object}
		}
	}

	return out
}

// resourceName takes the last literal segment of a path, singularised crudely,
// as the resource name.
//
// Crudely, and the draft says so: guessing English plurals is exactly the kind
// of cleverness that produces a fake which is subtly wrong for "person",
// "category" or "status", which is why the format makes collection an explicit
// declaration in the first place.
func resourceName(path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")

	for i := len(segments) - 1; i >= 0; i-- {
		segment := segments[i]
		if strings.HasPrefix(segment, "{") || segment == "" {
			continue
		}

		segment = strings.NewReplacer(".json", "", "-", "_").Replace(segment)

		switch {
		case strings.HasSuffix(segment, "ies"):
			return segment[:len(segment)-3] + "y"
		case strings.HasSuffix(segment, "sses"), strings.HasSuffix(segment, "shes"):
			return segment[:len(segment)-2]
		case strings.HasSuffix(segment, "s") && !strings.HasSuffix(segment, "ss"):
			return segment[:len(segment)-1]
		}

		return segment
	}

	return "record"
}

func writeDraftResources(b *strings.Builder, doc *Document, resources map[string]draftResource) {
	if len(resources) == 0 {
		b.WriteString("resources: {}\n\n")

		return
	}

	b.WriteString("resources:\n")

	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintf(b, "  %s:\n    id:\n      style: opaque\n    fields:\n", name)

		for _, property := range doc.Properties(resources[name].schema) {
			if property.Name == "id" {
				continue
			}

			if !validFieldName(property.Name) {
				fmt.Fprintf(b, "      # %q is not a name this format can declare, and it is skipped\n      # rather than mangled.\n", property.Name)

				continue
			}

			resolved := doc.Resolve(property.Schema)

			fmt.Fprintf(b, "      %s:\n", property.Name)

			if kind := fieldType(resolved); kind != "" {
				fmt.Fprintf(b, "        type: %s\n", kind)
			}

			if property.Required {
				b.WriteString("        required: true\n")
			}

			if resolved != nil && len(resolved.Enum) > 0 {
				fmt.Fprintf(b, "        # Declared values: %s. Which of these is the one that\n        # reads like success and is not is the question worth answering.\n", enumList(resolved.Enum))
			}
		}

		b.WriteString("\n")
	}
}

// validFieldName reports whether a property name can be a Recipe field key.
func validFieldName(name string) bool {
	if name == "" {
		return false
	}

	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}

	return true
}

func fieldType(s *Schema) string {
	if s == nil {
		return ""
	}

	switch s.Type {
	case "string":
		if s.Format == "date-time" {
			return "datetime"
		}

		return "string"
	case "integer":
		return "integer"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "list"
	}

	return ""
}

func enumList(values []any) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, fmt.Sprint(value))
	}

	sort.Strings(out)

	if len(out) > 8 {
		return strings.Join(out[:8], ", ") + ", and more"
	}

	return strings.Join(out, ", ")
}

func writeDraftRoutes(b *strings.Builder, doc *Document, resources map[string]draftResource) {
	b.WriteString("routes:\n")

	for _, path := range doc.SortedPaths() {
		for _, mo := range doc.Paths[path].Operations() {
			kind := operationKind(path, mo.Method)
			if kind == "" {
				fmt.Fprintf(b, "  # %s %s is not a create, read, update, delete or list, so it has no\n  # operation in this format. Model it deliberately or leave it out.\n", mo.Method, path)

				continue
			}

			resource := resourceName(path)
			if _, ok := resources[resource]; !ok {
				fmt.Fprintf(b, "  # %s %s returns no object this could read a shape from.\n", mo.Method, path)

				continue
			}

			fmt.Fprintf(b, "  - method: %s\n    path: %s\n    resource: %s\n    operation: %s\n",
				mo.Method, recipePath(path), resource, kind)

			if _, code := doc.Success(mo.Operation); code != "" && code != "200" && !strings.ContainsAny(code, "Xx") {
				fmt.Fprintf(b, "    status: %s\n", code)
			}

			b.WriteString("\n")
		}
	}
}

// recipePath renames the last path parameter to {id}, which is what the
// runtime looks for, and leaves the others alone so they can be scoped.
func recipePath(path string) string {
	segments := strings.Split(path, "/")

	for i := len(segments) - 1; i >= 0; i-- {
		if strings.HasPrefix(segments[i], "{") && strings.HasSuffix(segments[i], "}") {
			segments[i] = "{id}"

			break
		}
	}

	return strings.Join(segments, "/")
}

// operationKind maps a method and path onto one of the five operations, by the
// only signal available: whether the path ends in a parameter.
func operationKind(path, method string) string {
	last := ""

	for _, segment := range strings.Split(strings.Trim(path, "/"), "/") {
		if segment != "" {
			last = segment
		}
	}

	item := strings.HasPrefix(last, "{") || strings.Contains(last, "{")

	switch strings.ToUpper(method) {
	case "GET":
		if item {
			return "get"
		}

		return "list"
	case "POST":
		if item {
			return ""
		}

		return "create"
	case "PUT", "PATCH":
		if item {
			return "update"
		}

		return ""
	case "DELETE":
		if item {
			return "delete"
		}

		return ""
	}

	return ""
}

func writeDraftErrors(b *strings.Builder, doc *Document) {
	statuses := map[string]string{}

	for _, path := range doc.SortedPaths() {
		for _, mo := range doc.Paths[path].Operations() {
			for _, code := range doc.Failures(mo.Operation) {
				if _, seen := statuses[code]; !seen {
					statuses[code] = mo.Operation.Responses[code].Description
				}
			}
		}
	}

	b.WriteString("errors:\n")

	if len(statuses) == 0 {
		b.WriteString(`  # The description declares no failures at all, which is a statement about
  # the description rather than about the API. Every API fails.
  resource_missing:
    status: 404
    message: "Not found"

`)

		return
	}

	codes := make([]string, 0, len(statuses))
	for code := range statuses {
		codes = append(codes, code)
	}

	sort.Strings(codes)

	for _, code := range codes {
		name := errorName(code)
		message := strings.TrimSpace(statuses[code])

		if message == "" {
			message = "Request failed"
		}

		fmt.Fprintf(b, "  %s:\n    status: %s\n    message: %q\n", name, code, firstLine(message))
	}

	b.WriteString(`
  # The description says which statuses happen and never why. A failure worth
  # modelling is one somebody has to branch on: a declined card that is worth
  # retrying against one that is not, a rate limit that names its window.
  # Those are read from the prose, not from the schema.

`)
}

func errorName(code string) string {
	switch code {
	case "400":
		return "parameter_missing"
	case "401":
		return "authentication_error"
	case "403":
		return "insufficient_scope"
	case "404":
		return "resource_missing"
	case "409":
		return "conflict"
	case "422":
		return "unprocessable"
	case "429":
		return "rate_limit"
	}

	return "status_" + code
}

func firstLine(s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		return strings.TrimSpace(s[:i])
	}

	return s
}

// guessCapability picks a category from words in the description's title, and
// reports whether it had any evidence for the choice.
//
// A description says what the endpoints are and never what the product is, so
// this is a keyword match and it is labelled as one in the file it writes. The
// alternative was leaving the field out, and a draft that does not load is
// worse than a draft with a wrong word in it that the header tells you to
// check.
func guessCapability(doc *Document) (string, bool) {
	haystack := strings.ToLower(doc.Info.Title + " " + doc.Info.Description)

	for _, pair := range []struct{ word, capability string }{
		{"payment", "payments"}, {"checkout", "payments"}, {"billing", "payments"},
		{"subscription", "payments"}, {"invoice", "accounting"}, {"accounting", "accounting"},
		{"bank", "banking"}, {"payroll", "payroll"}, {"tax", "tax"},
		{"email", "email"}, {"mail", "email"}, {"sms", "sms"}, {"messaging", "sms"},
		{"chat", "chat"}, {"voice", "voice"}, {"telephony", "voice"}, {"push", "push"},
		{"auth", "auth"}, {"identity", "identity"}, {"crm", "crm"}, {"support", "support"},
		{"ticket", "support"}, {"marketing", "marketing"}, {"commerce", "commerce"},
		{"store", "commerce"}, {"shipping", "shipping"}, {"shipment", "shipping"},
		{"storage", "storage"}, {"database", "database"}, {"queue", "queue"},
		{"search", "search"}, {"cdn", "cdn"}, {"hosting", "hosting"},
		{"monitor", "observability"}, {"logging", "observability"}, {"analytics", "analytics"},
		{"feature flag", "flags"}, {"pipeline", "ci"}, {"repositor", "vcs"}, {"git", "vcs"},
		{"issue", "issues"}, {"calendar", "calendar"}, {"file", "files"},
		{"video", "media"}, {"image", "media"}, {"transcri", "ai"}, {"model", "ai"},
		{"signature", "signing"}, {"sign", "signing"}, {"schedul", "scheduling"},
		{"form", "forms"}, {"content", "cms"}, {"cloud", "infrastructure"},
	} {
		if strings.Contains(haystack, pair.word) {
			return pair.capability, true
		}
	}

	// Something valid, so the draft loads, and the header says to change it.
	return "infrastructure", false
}
