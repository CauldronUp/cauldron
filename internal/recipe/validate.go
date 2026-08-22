// Validation: what a Recipe has to say before it is allowed to ship.

package recipe

import (
	"fmt"
	"regexp"
	"strings"
)

// indexedSegment matches a path segment naming an array position, e.g. to[0].
var indexedSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*(\[\d+\])+$`)

// plainSegment matches an ordinary object key.
//
// Hyphens are allowed because JSON keys have them: npm sends dist-tags and
// plenty of providers send header-shaped names. This rule refused the first
// hyphenated key it ever met, which is what a pattern written against the
// Recipes that happened to exist will do. A dot is still refused, because the
// runtime splits on it and a dot inside a segment would silently mean two
// levels rather than one.
var plainSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	datePattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	validSchemes    = []string{"bearer", "basic", "header", "query", "none"}
	validOperations = []string{"create", "get", "list", "update", "delete"}
	validPagination = []string{"", "cursor", "offset", "page"}
	// The categories a Recipe may declare. Kept short on purpose: a taxonomy
	// with sixty words in it is a taxonomy nobody can hold in their head, and
	// the point of this field is that somebody looking for a payments
	// emulator finds all of them at once.
	validCapabilities = []string{
		"payments", "banking", "accounting", "tax", "payroll",
		"email", "sms", "chat", "push", "voice",
		"auth", "identity", "crm", "support", "marketing", "brokerage",
		"commerce", "shipping", "storage", "database", "queue",
		"search", "cdn", "hosting", "observability", "analytics",
		"flags", "ci", "vcs", "issues", "docs",
		"calendar", "files", "media", "ai", "signing",
		"scheduling", "hr", "forms", "cms", "infrastructure",
	}
	validDeletedBody = []string{"", "receipt", "record", "flagged", "id", "empty"}
	validCursorURL   = []string{"", "absolute", "path"}
	validSigning     = []string{"", "none", "hmac-sha256"}
	validListStyles  = []string{"", "envelope", "bare", "wrapped", "map"}
	validErrStyles   = []string{"", "nested", "flat", "list", "string_list", "text", "string"}
	validCodeTypes   = []string{"", "string", "number"}
	validIDTypes     = []string{"", "string", "number"}
	validIDStyles    = []string{"", "prefixed", "numeric", "timestamp", "opaque", "uuid", "hex", "digits"}
	validFieldTypes  = []string{"", "string", "integer", "number", "boolean", "timestamp", "timestamp_ms", "timestamp_ms_string", "datetime", "msdate", "list", "map"}
	// The types Cauldron fills in from the sandbox clock, and therefore the
	// only ones a stamped declaration can affect.
	timeFieldTypes = []string{"timestamp", "timestamp_ms", "timestamp_ms_string", "datetime", "msdate"}

	knownMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
)

// Parse decodes and validates a Recipe from YAML bytes.
// safeDigits is how many figures a double holds exactly. 2^53 has sixteen, so
// fifteen is safe with one to spare, and an identifier longer than that has to
// stay a string.
const safeDigits = 15

// Validate checks the Recipe is internally consistent.
func (r *Recipe) Validate() error {
	var problems []string

	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	if r.Name == "" {
		add("recipe name is required")
	} else if !namePattern.MatchString(r.Name) {
		add("recipe name %q must be lowercase words separated by hyphens", r.Name)
	}

	if r.Version == "" {
		add("version is required")
	} else if !versionPattern.MatchString(r.Version) {
		add("version %q must look like 1.2.3", r.Version)
	}

	if r.Upstream.API == "" {
		add("upstream.api is required so the Recipe records which API version it targets")
	}

	if r.Capability == "" {
		add("capability is required: a hundred Recipes are only findable by what they do")
	} else if !contains(validCapabilities, r.Capability) {
		add("capability %q must be one of %s", r.Capability, strings.Join(validCapabilities, ", "))
	}

	if r.Auth.Scheme != "" && !contains(validSchemes, r.Auth.Scheme) {
		add("auth.scheme %q must be one of %s", r.Auth.Scheme, strings.Join(validSchemes, ", "))
	}

	if r.Auth.Credential != "" && r.Auth.Credential != "username" && r.Auth.Credential != "password" {
		add("auth.credential %q must be username or password", r.Auth.Credential)
	}

	if r.Auth.Credential != "" && r.Auth.Scheme != "basic" {
		add("auth.credential only applies to the basic scheme")
	}

	if r.Auth.Scheme == "query" && r.Auth.Param == "" {
		add("auth.param is required when the scheme is query")
	}

	if r.Auth.Param != "" && r.Auth.Scheme != "query" {
		add("auth.param only applies to the query scheme")
	}

	r.validateResources(add)
	r.validateRoutes(add)
	for i, c := range r.Conformance {
		if c.Arm != "" {
			if _, ok := r.Errors[c.Arm]; !ok {
				add("conformance[%d] %q: arms %q, which is not in the errors table", i, c.Name, c.Arm)
			}

			// Arming something and then expecting a status it does not
			// produce means the fault did nothing, and the case passes while
			// proving the opposite of what it says. That is the same class of
			// mistake as a case asserting its own request back.
			//
			// This used to compare against 400, on the reasoning that an
			// armed thing is a failure and a failure is a 4xx or a 5xx. Not
			// every armed thing is a failure: Snowflake's SQL API answers the
			// same endpoint with a 202 and a statement handle when the query
			// is slow, which is an alternate path rather than an error, and a
			// case arming it has to expect the 202 it installs. Comparing
			// against the armed entry's own status catches the real mistake
			// -- arm a 429, expect a 200 -- and permits the honest one.
			if armed, ok := r.Errors[c.Arm]; ok && c.Expect.Status > 0 && c.Expect.Status != armed.Status {
				add("conformance[%d] %q: arms %q, which answers %d, and expects %d, so what it installed changed nothing",
					i, c.Name, c.Arm, armed.Status, c.Expect.Status)
			}
		}

		if !c.Expect.NoBody {
			continue
		}

		if len(c.Expect.Body) > 0 || len(c.Expect.Matches) > 0 || len(c.Expect.Absent) > 0 {
			add("conformance[%d] %q: expect.no_body claims the response is empty, so there is nothing for body, matches or absent to look at", i, c.Name)
		}
	}

	// A field name the Recipe chooses is only a claim if a case asserts it
	// where the value exists. Asserting its absence on a last page holds
	// whatever the field happens to be called, so a Recipe could declare
	// cursor_field: next_page_token, be renamed to anything, and no case
	// would notice.
	//
	// Twenty-one of these were shipped before anything checked, across twenty
	// Recipes. Two of them turned out to have no successful list case at all:
	// every case touching the collection was checking a failure, so nothing
	// asserted what a working listing looks like.
	for _, declared := range []struct{ what, name string }{
		{"cursor_field", r.Responses.List.CursorField},
		{"count_field", r.Responses.List.CountField},
		{"has_more_field", r.Responses.List.HasMoreField},
		{"complete_field", r.Responses.List.CompleteField},
		{"final_field", r.Responses.List.FinalField},
	} {
		if declared.name == "" || assertsName(r.Conformance, declared.name) {
			continue
		}

		add("responses.list.%s is %q and no conformance case asserts that name where the value exists, so renaming it would break nothing",
			declared.what, declared.name)
	}

	// The same rule for the names a failure travels under. A Recipe saying
	// prose lives at error_message and codes at error_code has described the
	// shape a client unwraps, and if no case asserts either then both could be
	// called anything.
	//
	// Sixteen Recipes named a message field nothing asserted, Plaid's
	// error_message and Jira's errorMessages[] among them -- exactly the
	// provider-specific names somebody has to get right, and exactly the ones
	// a common name like "message" would not have made obvious.
	//
	// "-" is excluded because it is not a name. It says the provider sends no
	// such field, which is a claim about absence and is asserted by asserting
	// absence.
	for _, declared := range []struct{ what, name string }{
		{"message_field", r.Responses.Error.MessageField},
		{"code_field", r.Responses.Error.CodeField},
	} {
		if declared.name == "" || declared.name == "-" || assertsName(r.Conformance, declared.name) {
			continue
		}

		add("responses.error.%s is %q and no conformance case asserts that name where the value exists, so renaming it would break nothing",
			declared.what, declared.name)
	}

	if r.Responses.List.HasMoreField != "" && r.Responses.List.CompleteField != "" {
		add("responses.list declares both has_more_field and complete_field, which are the same flag with opposite senses")
	}

	// The whole point of a final field is that it is not the cursor. Naming
	// them the same thing would emit one key that changes meaning depending on
	// which page you are looking at, which is worse than either field alone.
	if f := r.Responses.List.FinalField; f != "" && f == r.Responses.List.CursorField {
		add("responses.list declares cursor_field and final_field as the same name %q, so one key would mean a page token on one page and a sync token on the next", f)
	}

	// A required header is an enforcement claim, and the only thing that
	// proves it is a case that leaves the header off and is refused.
	//
	// Enforcing Notion-Version is one of the reasons this project exists:
	// forgetting it is the classic Notion integration bug, and a fake that
	// waves it through lets code ship that fails on its first real call. A
	// Recipe could declare the header, never enforce it, and every case would
	// still be green -- they all send it.
	//
	// This is the same shape as the rule above about a field name no case
	// asserts. A name nothing exercises is a name that could be anything.
	for header := range r.RequiredHeaders {
		if omitsHeader(r.Conformance, header) {
			continue
		}

		add("required_headers declares %s and no conformance case omits it and is refused, so nothing here shows it is enforced rather than merely listed", header)
	}

	// A route naming an error to raise has to name one that exists, or the
	// only sign of the typo is a 404 in the shape it was overriding.
	for _, route := range r.Routes {
		if route.NotFound == "" {
			continue
		}

		if _, ok := r.Errors[route.NotFound]; !ok {
			add("%s %s declares not_found: %q and no error by that name is defined, so it would answer in the shape it was written to replace",
				route.Method, route.Path, route.NotFound)
		}
	}

	// A filter is a claim that a listing narrows itself when asked, and the
	// parameter's name is half of it. Every case listing without the parameter
	// exercises the default; only a case sending it pins what it is called.
	for _, route := range r.Routes {
		for _, f := range route.Filters {
			if f.Param == "" || sendsParam(r.Conformance, f.Param) {
				continue
			}

			add("%s %s declares a filter on %q and no conformance case sends it, so the parameter could be named anything and every case would still pass", route.Method, route.Path, f.Param)
		}
	}

	if r.Responses.Error.Key == "-" && r.Responses.Error.Style != "list" {
		add("responses.error.key is \"-\", which removes the envelope, but that only means anything for the list style")
	}

	if r.Responses.Error.Key == "-" && len(r.Responses.Error.Fields) > 0 {
		add("responses.error.key is \"-\", so there is no envelope for responses.error.fields to sit beside")
	}

	if r.Auth.Pattern != "" {
		if _, err := regexp.Compile(r.Auth.Pattern); err != nil {
			add("auth.pattern is not a valid regular expression: %v", err)
		}

		if len(r.Auth.Keys) > 0 {
			add("auth declares both keys and a pattern, and only one can decide whether a credential is accepted")
		}
	}

	if !contains(validCodeTypes, r.Responses.Error.CodeType) {
		add("responses.error.code_type %q must be one of %s", r.Responses.Error.CodeType, strings.Join(validCodeTypes[1:], ", "))
	}

	if !contains(validErrStyles, r.Responses.Error.Style) {
		add("responses.error.style %q must be one of %s", r.Responses.Error.Style, strings.Join(validErrStyles[1:], ", "))
	}

	if r.Responses.List.EntryStyle != "" && r.Responses.List.EntryStyle != "wrapped" {
		add("responses.list.entry_style %q must be wrapped", r.Responses.List.EntryStyle)
	}

	if !contains(validCursorURL, r.Responses.List.CursorURL) {
		add("responses.list.cursor_url %q must be one of %s", r.Responses.List.CursorURL, strings.Join(validCursorURL[1:], ", "))
	}

	if !contains(validListStyles, r.Responses.List.Style) {
		add("responses.list.style %q must be one of %s", r.Responses.List.Style, strings.Join(validListStyles[1:], ", "))
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key != "" {
		// The key is a fallback for resources that do not name their own
		// collection. When every one of them does, it is unreachable, and an
		// unreachable declaration reads as a description of the provider
		// while describing nothing. Five Recipes shipped with one before this
		// rule existed, each found by mutating it and watching nothing fail.
		covered := len(r.Resources) > 0

		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				covered = false
			}
		}

		if covered {
			add("responses.list.key is %q and every resource names its own collection, so nothing reads it",
				r.Responses.List.Key)
		}
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key == "" {
		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				add("resource %q needs a collection name, or responses.list.key must be set, because the list style is wrapped", name)
			}
		}
	}

	if !contains(validSigning, r.Webhooks.Signing.Scheme) {
		add("webhooks.signing.scheme %q is not supported", r.Webhooks.Signing.Scheme)
	}

	if r.Webhooks.Signing.Scheme == "hmac-sha256" && r.Webhooks.Signing.Header == "" {
		add("webhooks.signing.header is required when signing is enabled")
	}

	for _, event := range r.Webhooks.Events {
		if strings.TrimSpace(event) == "" {
			add("webhooks.events contains an empty event name")
		}
	}

	for _, name := range sortedKeys(r.Errors) {
		e := r.Errors[name]

		if e.Status < 100 || e.Status > 599 {
			add("error %q has status %d, which is not a valid HTTP status", name, e.Status)
		}
	}

	for _, header := range sortedKeys(r.RequiredHeaders) {
		required := r.RequiredHeaders[header]

		if name := required.Error; name != "" {
			if _, ok := r.Errors[name]; !ok {
				add("required_headers[%s] names error %q, which is not declared", header, name)
			}
		}

		for _, method := range required.Methods {
			if !knownMethods[strings.ToUpper(method)] {
				add("required_headers[%s] limits to method %q, which is not an HTTP method", header, method)
			}
		}
	}

	r.validateFixtures(add)
	r.validateCases(add)
	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

// fetchedByID reports whether any route fetches one record of a resource on
// its own, which is what makes an identifier something a client can use.
func fetchedByID(r *Recipe, resource string) bool {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return true
		}
	}

	return false
}

// fetchRoute names the first route that fetches one, for the message.
func fetchRoute(r *Recipe, resource string) string {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return route.Method + " " + route.Path
		}
	}

	return "another route"
}

// onlyCreated reports whether every route touching a resource creates it.
//
// Such a resource is a receipt: the provider hands back an identifier and
// nothing else, and there is nowhere to go and read it. Knock's trigger is
// one. A resource that is read anywhere is not, and a resource with no fields
// that is read anywhere is describing a response nobody can use.
func onlyCreated(r *Recipe, resource string) bool {
	touched := false

	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		touched = true

		if route.Operation != "create" {
			return false
		}
	}

	return touched
}

// identifierFieldName renders the property an identifier goes out under, for a
// message, naming the default rather than reporting an empty string.
func identifierFieldName(field string) string {
	if field == "" {
		return "id"
	}

	return field
}

// styleName renders an id style for a message, naming the default rather than
// reporting an empty string nobody wrote.
func styleName(style string) string {
	if style == "" {
		return "prefixed"
	}

	return style
}

// sortedKeys keeps iteration deterministic, so validation problems always
// appear in the same order.
// seededAs reports whether a fixture value matches the type its field
// declares, and names what it is when it does not.
//
// A field's type is documentation. Nothing in the runtime reads it: the value
// a fixture seeds is the value that goes on the wire, with whatever type YAML
// gave it. So a Recipe can declare a field a string, seed it with something
// that is not one, and serve the wrong type for as long as no case looks.
//
// The reason to check it is YAML 1.1 rather than carelessness. An unquoted
// off, no, on or yes is a boolean, and those are all real values of real
// string fields: a DigitalOcean droplet's status is one of new, active, off
// and archive, and the fixture that seeded off to say a machine was powered
// down served false to everything that read it. The comment beside it
// explained that a cost report counting running machines would miss that
// droplet, and no case asserted on it, so the claim and the wire disagreed
// from the day it was written.
//
// Only the unambiguous directions are checked. A number field seeded with a
// whole number is fine, and a timestamp is a number whichever width it lands
// in.
func seededAs(declared string, value any) (string, bool) {
	if value == nil {
		return "", true
	}

	name := "something else"

	switch value.(type) {
	case bool:
		name = "a boolean"
	case string:
		name = "a string"
	case int, int64, uint64:
		name = "a whole number"
	case float64, float32:
		name = "a number"
	case []any:
		name = "a sequence"
	case map[string]any:
		name = "a mapping"
	}

	switch declared {
	case "string", "datetime":
		_, ok := value.(string)

		return name, ok
	case "boolean":
		_, ok := value.(bool)

		return name, ok
	case "integer", "number", "timestamp", "timestamp_ms":
		switch value.(type) {
		case int, int64, uint64, float64, float32:
			return name, true
		}

		return name, false
	}

	return name, true
}

// seededID reports whether a fixture's explicit identifier has the shape the
// same resource would mint. It returns what is wrong with it, and true when
// nothing is.
//
// A fixture that seeds an id in a shape the generator would never produce puts
// two shapes in one collection: the seeded records look one way and anything
// created during a run looks another. Client code that parses or
// prefix-matches an id -- which is common enough that ID exists as a concept
// here at all -- then works against the fixtures and fails against the
// records it made itself.
//
// It was a surviving mutation that found this. Monday declares ten-digit ids
// and seeds ten-digit ids, and changing the declaration to six broke nothing:
// every case reads a seeded record, so the declared shape was never on the
// wire and could have drifted from the provider without a single case
// noticing.
func seededID(id ID, value any) (string, bool) {
	text, ok := value.(string)
	if !ok {
		// A numeric id is checked as a number elsewhere; nothing to say here.
		return "", true
	}

	if text == "" {
		return "is empty", false
	}

	switch id.Style {
	case "uuid":
		if !looksUUID(text) {
			return "is not a UUID", false
		}

	case "hex":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isHexDigit) {
			return "is not hexadecimal", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d characters, not the declared %d", len(body), id.Length), false
		}

	case "digits":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isDecimalDigit) {
			return "is not all digits", false
		}

		// A snowflake is a number written as text, and a number does not
		// begin with a zero.
		if body[0] == '0' {
			return "begins with a zero, which a number written as text does not", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d digits, not the declared %d", len(body), id.Length), false
		}

	case "numeric":
		if !all(text, isDecimalDigit) {
			return "is not a number", false
		}

	case "", "prefixed":
		if id.Prefix != "" && !strings.HasPrefix(text, id.Prefix) {
			return "does not start with " + id.Prefix, false
		}
	}

	return "", true
}

func looksUUID(text string) bool {
	runs := strings.Split(text, "-")
	if len(runs) != 5 {
		return false
	}

	for i, want := range []int{8, 4, 4, 4, 12} {
		if len(runs[i]) != want || !all(runs[i], isHexDigit) {
			return false
		}
	}

	return true
}

// IsPath reports whether a name is a well-formed dotted path rather than a
// literal key that happens to contain a dot.
//
// The distinction is not academic. Dropbox names a field ".tag" -- the leading
// dot is part of the name, not a separator -- so treating every dotted name as
// a path turns it into an object under an empty key. A path is at least two
// segments and every one of them is a name.
func IsPath(name string) bool {
	segments := strings.Split(name, ".")
	if len(segments) < 2 {
		return false
	}

	for _, segment := range segments {
		if !plainSegment.MatchString(segment) && !indexedSegment.MatchString(segment) {
			return false
		}
	}

	return true
}
