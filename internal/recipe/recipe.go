// Package recipe defines the Recipe format — the description of how Cauldron
// emulates one external dependency.
//
// A Recipe is not a mock. A mock returns a shape; a Recipe models behaviour:
// what resources exist, how they change, what the provider emits afterwards,
// and how it fails. The format is declarative on purpose, so that the majority
// of Recipes are data a contributor can read and review rather than code.
package recipe

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Recipe is a complete emulation description for one provider.
type Recipe struct {
	Name      string              `yaml:"recipe"`
	Version   string              `yaml:"version"`
	Upstream  Upstream            `yaml:"upstream"`
	Auth      Auth                `yaml:"auth"`
	Resources map[string]Resource `yaml:"resources"`
	Routes    []Route             `yaml:"routes"`
	Webhooks  Webhooks            `yaml:"webhooks"`
	Responses Responses           `yaml:"responses"`
	Errors    map[string]Error    `yaml:"errors"`
	Fixtures  map[string]Fixture  `yaml:"fixtures"`
	// Conformance is the evidence that this Recipe resembles the real provider.
	Conformance []Case `yaml:"conformance"`
}

// Case is one checkable claim about the provider's behaviour.
//
// The point of a conformance case is not that the emulator passes it. Any fake
// passes its own tests. The point is provenance: every case cites where the
// expectation came from, and records whether it was observed against the real
// API or only read in the documentation. A developer deciding whether to trust
// this emulator can then read the evidence rather than the marketing.
type Case struct {
	Name string `yaml:"name"`
	// Source cites the provider documentation or transcript the expectation
	// came from. Required: an uncited claim about someone else's API is a
	// guess wearing a test's clothing.
	Source string `yaml:"source"`
	// Verified is the date this case was last checked against the real API,
	// as YYYY-MM-DD. Empty means the expectation was read, not observed, and
	// the report says so rather than quietly counting it as proof.
	Verified string `yaml:"verified"`
	// Fixture is seeded before the case runs. Empty leaves the sandbox as it is,
	// which lets a group of cases build on each other in order.
	Fixture string      `yaml:"fixture"`
	Request Request     `yaml:"request"`
	Expect  Expectation `yaml:"expect"`
}

// Request is the call a conformance case makes.
type Request struct {
	Method  string            `yaml:"method"`
	Path    string            `yaml:"path"`
	Query   map[string]string `yaml:"query"`
	Headers map[string]string `yaml:"headers"`
	// Form sends application/x-www-form-urlencoded, which is what Stripe's own
	// SDKs send. JSON sends a JSON body. A case may set at most one.
	Form map[string]string `yaml:"form"`
	JSON map[string]any    `yaml:"json"`
}

// Expectation is what the provider is claimed to answer.
//
// Body matching is a subset: a case asserts the fields it is making a claim
// about and ignores the rest, so a Recipe can grow a field without invalidating
// every case ever written about it.
type Expectation struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
	// Matches holds dotted field paths to regular expressions, for values that
	// are correct in shape rather than exact, such as generated identifiers.
	Matches map[string]string `yaml:"matches"`
	// HeaderMatches holds response header names to regular expressions, for
	// headers that carry a generated value. A plain `headers` entry compares
	// substrings, which cannot assert that a header is merely present and
	// well-formed.
	HeaderMatches map[string]string `yaml:"header_matches"`
	// Absent lists fields that must not appear. Providers are as specific about
	// what they omit as what they send.
	Absent []string `yaml:"absent"`
}

// Upstream records which real API version this Recipe targets. Without it, a
// Recipe silently rots as the provider moves on.
type Upstream struct {
	API  string `yaml:"api"`
	Docs string `yaml:"docs"`
}

// Auth describes how the provider authenticates callers.
type Auth struct {
	// Scheme is one of: bearer, basic, header, none.
	Scheme string `yaml:"scheme"`
	// Header is the header carrying the credential, when scheme is header.
	Header string `yaml:"header"`
	// Prefix is stripped from the credential before comparison, e.g. "Bearer ".
	Prefix string `yaml:"prefix"`
	// Keys are the credentials the emulator accepts. Test keys only — a Recipe
	// must never carry a real secret.
	Keys []string `yaml:"keys"`
}

// Responses describes the envelopes a provider wraps its payloads in.
//
// This exists because providers genuinely disagree: Stripe returns
// {object, data, has_more}, GitHub returns a bare array, Shopify nests under a
// resource key. Hardcoding one of them would make every other provider a
// second-class citizen.
type Responses struct {
	List     ListResponse     `yaml:"list"`
	Error    ErrorResponse    `yaml:"error"`
	Resource ResourceResponse `yaml:"resource"`
	Success  SuccessResponse  `yaml:"success"`
}

// SuccessResponse describes what a provider adds to every successful body.
//
// Slack stamps {"ok": true} on everything and its clients check it before
// looking at anything else. A fake that omits it fails at the first line of
// every handler written against the real API.
type SuccessResponse struct {
	Fields map[string]any `yaml:"fields"`
}

// ResourceResponse describes how a single object comes back.
//
// Shopify wraps it under the singular resource name, so a client reads
// body.order.id. Stripe and GitHub return the object itself. Getting this
// wrong is not a cosmetic difference: every field access is one level out.
type ResourceResponse struct {
	// Style is bare (the default) or wrapped.
	Style string `yaml:"style"`
}

// ErrorResponse describes the envelope a provider puts failures in.
//
// Stripe nests under "error" with a type and a code; GitHub sends a flat object
// with a message and a documentation link. Code that unwraps one and receives
// the other does not report a helpful failure, it panics.
type ErrorResponse struct {
	// Style is nested (Stripe, the default), flat (GitHub) or list (SendGrid,
	// which sends {"errors": [{...}]} because one request can fail several
	// ways at once).
	Style string `yaml:"style"`
	// Key is the property holding the array when the style is list.
	Key string `yaml:"key"`
	// MessageField names the property carrying the human-readable message when
	// the style is flat. Empty means "message". Set it to "-" to omit the
	// message entirely, which Slack does: its errors are a code and nothing
	// else, and inventing prose the provider never sends is still infidelity.
	MessageField string `yaml:"message_field"`
	// CodeField names the property carrying the error code in a flat envelope.
	// Twilio sends one and its clients switch on it; GitHub does not send one
	// at all, so this stays empty unless a Recipe claims otherwise. A code that
	// is entirely digits is sent as a number, because Twilio's is. As with
	// MessageField, "-" omits the code, which the nested style needs too:
	// Airtable nests its error but sends only a type and a message.
	CodeField string `yaml:"code_field"`
	// StatusField names a property echoing the HTTP status inside the body,
	// which Twilio does.
	StatusField string `yaml:"status_field"`
	// Fields are constants the provider adds to every error, such as GitHub's
	// documentation_url.
	Fields map[string]any `yaml:"fields"`
}

// ListResponse describes how a collection is returned.
type ListResponse struct {
	// Style is one of: envelope (Stripe), bare (GitHub), wrapped (Shopify).
	// Empty means envelope, which keeps existing Recipes working.
	Style string `yaml:"style"`
	// Key is the wrapping property name, required when style is wrapped.
	Key string `yaml:"key"`
	// URL asks the envelope to echo the request path, which Stripe does.
	URL bool `yaml:"url"`
	// CursorField names a property carrying the next cursor. Most providers do
	// not send one: Stripe expects the caller to pass the last id back as
	// starting_after. Leaving it empty is therefore the faithful default, and
	// setting it is a deliberate claim that the provider really sends it.
	//
	// A dotted name nests, so Slack's response_metadata.next_cursor is
	// expressible without a second mechanism.
	CursorField string `yaml:"cursor_field"`
}

// Resource is an object type the provider exposes.
type Resource struct {
	// Collection is the plural name the provider wraps lists in, e.g. "orders"
	// for an order. Declared rather than derived: guessing English plurals is
	// exactly the kind of cleverness that produces a fake which is subtly
	// wrong for "person", "category" or "status".
	Collection string           `yaml:"collection"`
	ID         ID               `yaml:"id"`
	Fields     map[string]Field `yaml:"fields"`
	// Constants are fields the provider always sends with a fixed value, such
	// as Stripe's object discriminator and livemode flag. Unlike a default they
	// cannot be overridden by the caller, because the provider does not let you
	// override them either. Applications really do branch on these.
	Constants map[string]any `yaml:"constants"`
}

// ID describes how the provider mints identifiers. Getting this right matters
// more than it looks: applications routinely parse or prefix-match IDs.
type ID struct {
	// Style is one of: prefixed (cus_abc123), numeric (1, 2, 3), timestamp
	// (1767225600.000100, which is how Slack identifies a message) or opaque
	// (a bare random string, which is what SendGrid returns as a message id).
	// Empty means prefixed.
	Style  string `yaml:"style"`
	Prefix string `yaml:"prefix"`
	Length int    `yaml:"length"`
	// Field is the property the provider returns the identifier in. Empty means
	// "id". Twilio calls it "sid" everywhere, and code that reads response.id
	// against Twilio gets nothing at all.
	Field string `yaml:"field"`
}

// Field is a single attribute on a resource.
type Field struct {
	// Type is string, integer, boolean, timestamp (a Unix integer, which is
	// what Stripe and Twilio send) or datetime (an RFC 3339 string, which is
	// what GitHub, HubSpot and most newer APIs send). The difference is not
	// cosmetic: one parses as a number and the other does not.
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Default  any    `yaml:"default"`
	// In nests this field under a sub-object on the wire. HubSpot puts every
	// business attribute under "properties" and leaves only id, timestamps and
	// archived at the top level, so a client reads contact.properties.email.
	// The store stays flat; only the shape on the wire changes, and requests
	// are flattened back on the way in.
	In string `yaml:"in"`
}

// Route binds an HTTP method and path to an operation on a resource.
type Route struct {
	Method   string `yaml:"method"`
	Path     string `yaml:"path"`
	Resource string `yaml:"resource"`
	// Operation is one of: create, get, list, update, delete.
	Operation string `yaml:"operation"`
	// Scope names the path parameters that partition this resource, e.g.
	// [owner, repo] for /repos/{owner}/{repo}/issues. Scoped requests only
	// ever see records whose matching fields agree, and creates stamp them.
	//
	// Path parameters left out of scope are ignored, which is how an API
	// version segment like /admin/api/{version}/orders stays a path parameter
	// without becoming a filter.
	Scope []string `yaml:"scope"`
	// Status is the success status this route returns. Empty means 200. GitHub
	// answers a create with 201, Stripe with 200, and a client checking for one
	// exact code is not being unreasonable.
	Status int `yaml:"status"`
	// IDFrom says where the identifier comes from when it is not a path
	// parameter: "query:channel" or "body:channel". Slack and every other
	// RPC-shaped API put it in the query string or the body, and without this
	// the format could only describe APIs that happen to be RESTful.
	IDFrom string `yaml:"id_from"`
	// EmptyBody sends no body at all. SendGrid accepts a send with 202 and
	// nothing else, and a client that calls .json() on that response throws.
	// An emulator that helpfully returns an object hides the bug.
	EmptyBody bool `yaml:"empty_body"`
	// Headers are response headers this route sets. "{id}" is replaced with
	// the record's identifier, which is how SendGrid hands back the message id
	// a client needs to correlate a later event with the send.
	Headers    map[string]string `yaml:"headers"`
	Pagination Pagination        `yaml:"pagination"`
}

// Pagination describes how a list endpoint pages.
type Pagination struct {
	// Style is one of: cursor, offset, page.
	Style string `yaml:"style"`
	Limit int    `yaml:"limit"`
}

// Webhooks describes what the provider sends back, and how it signs it.
type Webhooks struct {
	Events  []string `yaml:"events"`
	Signing Signing  `yaml:"signing"`
}

// Signing describes webhook payload signing.
type Signing struct {
	// Scheme is one of: hmac-sha256, none.
	Scheme string `yaml:"scheme"`
	Header string `yaml:"header"`
	Secret string `yaml:"secret"`
}

// Error is a named failure mode that `cauldron fault` can inject.
type Error struct {
	Status int    `yaml:"status"`
	Code   string `yaml:"code"`
	// Type is the provider's error category, which is often a much smaller set
	// than the codes. Stripe has four types and dozens of codes, and client
	// libraries switch on the type. Empty falls back to the code, which is
	// wrong often enough that every Recipe should set it.
	Type    string            `yaml:"type"`
	Message string            `yaml:"message"`
	Headers map[string]string `yaml:"headers"`
}

// Fixture is a named seed dataset: resource name to a list of records.
type Fixture map[string][]map[string]any

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	datePattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	validSchemes    = []string{"bearer", "basic", "header", "none"}
	validOperations = []string{"create", "get", "list", "update", "delete"}
	validPagination = []string{"", "cursor", "offset", "page"}
	validSigning    = []string{"", "none", "hmac-sha256"}
	validListStyles = []string{"", "envelope", "bare", "wrapped"}
	validIDStyles   = []string{"", "prefixed", "numeric", "timestamp", "opaque"}
	validFieldTypes = []string{"", "string", "integer", "number", "boolean", "timestamp", "datetime"}
)

// Load reads and validates a Recipe from a YAML file.
func Load(path string) (*Recipe, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(contents)
}

// Parse decodes and validates a Recipe from YAML bytes.
func Parse(contents []byte) (*Recipe, error) {
	var r Recipe

	decoder := yaml.NewDecoder(strings.NewReader(string(contents)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&r); err != nil {
		return nil, fmt.Errorf("recipe is not valid YAML: %w", err)
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return &r, nil
}

// ValidationError collects every problem with a Recipe, so an author sees all
// of them at once instead of fixing one and rerunning.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 1 {
		return "invalid recipe: " + e.Problems[0]
	}

	return fmt.Sprintf("invalid recipe:\n  - %s", strings.Join(e.Problems, "\n  - "))
}

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

	if r.Auth.Scheme != "" && !contains(validSchemes, r.Auth.Scheme) {
		add("auth.scheme %q must be one of %s", r.Auth.Scheme, strings.Join(validSchemes, ", "))
	}

	if len(r.Resources) == 0 {
		add("at least one resource is required")
	}

	for _, name := range sortedKeys(r.Resources) {
		resource := r.Resources[name]

		if !contains(validIDStyles, resource.ID.Style) {
			add("resource %q has id.style %q, which must be one of %s", name, resource.ID.Style, strings.Join(validIDStyles[1:], ", "))
		}

		// A prefix is required for the prefixed style precisely because
		// applications prefix-match identifiers. Providers whose ids carry no
		// prefix say so by declaring the opaque style, rather than leaving it
		// ambiguous whether the omission was deliberate.
		if (resource.ID.Style == "" || resource.ID.Style == "prefixed") && resource.ID.Prefix == "" {
			add("resource %q must declare an id.prefix, or use id.style opaque if the provider's identifiers have none", name)
		}

		if len(resource.Fields) == 0 {
			add("resource %q has no fields", name)
		}

		for _, field := range sortedKeys(resource.Fields) {
			if !contains(validFieldTypes, resource.Fields[field].Type) {
				add("resource %q field %q has type %q, which must be one of %s",
					name, field, resource.Fields[field].Type, strings.Join(validFieldTypes[1:], ", "))
			}
		}
	}

	if len(r.Routes) == 0 {
		add("at least one route is required")
	}

	seen := map[string]bool{}

	for i, route := range r.Routes {
		where := fmt.Sprintf("route %d (%s %s)", i+1, route.Method, route.Path)

		if route.Method == "" {
			add("%s: method is required", where)
		} else if route.Method != strings.ToUpper(route.Method) {
			add("%s: method must be upper case", where)
		}

		if !strings.HasPrefix(route.Path, "/") {
			add("%s: path must start with /", where)
		}

		key := route.Method + " " + route.Path
		if seen[key] {
			add("%s: duplicate route", where)
		}
		seen[key] = true

		if route.Operation == "" {
			add("%s: operation is required", where)
		} else if !contains(validOperations, route.Operation) {
			add("%s: operation %q must be one of %s", where, route.Operation, strings.Join(validOperations, ", "))
		}

		if route.Resource == "" {
			add("%s: resource is required", where)
		} else if _, ok := r.Resources[route.Resource]; !ok {
			add("%s: unknown resource %q", where, route.Resource)
		}

		for _, name := range route.Scope {
			if name == "id" {
				add("%s: id cannot be a scope parameter", where)
				continue
			}

			if !strings.Contains(route.Path, "{"+name+"}") {
				add("%s: scope %q does not appear in the path", where, name)
			}

			if resource, ok := r.Resources[route.Resource]; ok {
				if _, declared := resource.Fields[name]; !declared {
					add("%s: scope %q is not a field on resource %q", where, name, route.Resource)
				}
			}
		}

		if !contains(validPagination, route.Pagination.Style) {
			add("%s: pagination.style %q is not supported", where, route.Pagination.Style)
		}

		if route.IDFrom != "" {
			source, name, ok := strings.Cut(route.IDFrom, ":")

			switch {
			case !ok || name == "":
				add("%s: id_from %q must look like query:channel or body:channel", where, route.IDFrom)
			case source != "query" && source != "body":
				add("%s: id_from source %q must be query or body", where, source)
			case strings.Contains(route.Path, "{id}"):
				add("%s: id_from and an {id} path parameter cannot both apply", where)
			}
		}

		if route.Operation != "list" && route.Operation != "create" &&
			route.IDFrom == "" && !strings.Contains(route.Path, "{id}") {
			add("%s: a %s needs an {id} in the path or an id_from", where, route.Operation)
		}
	}

	if !contains(validListStyles, r.Responses.List.Style) {
		add("responses.list.style %q must be one of %s", r.Responses.List.Style, strings.Join(validListStyles[1:], ", "))
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

	for _, fixtureName := range sortedKeys(r.Fixtures) {
		for _, resourceName := range sortedKeys(r.Fixtures[fixtureName]) {
			if _, ok := r.Resources[resourceName]; !ok {
				add("fixture %q seeds unknown resource %q", fixtureName, resourceName)
			}
		}
	}

	for i, c := range r.Conformance {
		where := fmt.Sprintf("conformance case %d", i+1)
		if c.Name != "" {
			where = fmt.Sprintf("conformance case %q", c.Name)
		} else {
			add("%s: name is required", where)
		}

		if c.Source == "" {
			add("%s: source is required, so a reader can check the claim against the provider", where)
		}

		if c.Verified != "" && !datePattern.MatchString(c.Verified) {
			add("%s: verified %q must be a date like 2026-08-15", where, c.Verified)
		}

		if c.Request.Method == "" {
			add("%s: request.method is required", where)
		} else if c.Request.Method != strings.ToUpper(c.Request.Method) {
			add("%s: request.method must be upper case", where)
		}

		if !strings.HasPrefix(c.Request.Path, "/") {
			add("%s: request.path must start with /", where)
		}

		if len(c.Request.Form) > 0 && len(c.Request.JSON) > 0 {
			add("%s: a request sends form or json, not both", where)
		}

		if c.Fixture != "" {
			if _, ok := r.Fixtures[c.Fixture]; !ok {
				add("%s: unknown fixture %q", where, c.Fixture)
			}
		}

		if c.Expect.Status == 0 {
			add("%s: expect.status is required", where)
		}

		for field, pattern := range c.Expect.Matches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.matches[%s] is not a valid regular expression: %v", where, field, err)
			}
		}

		for header, pattern := range c.Expect.HeaderMatches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.header_matches[%s] is not a valid regular expression: %v", where, header, err)
			}
		}

		if c.Expect.Status < 400 && len(c.Expect.Body) == 0 && len(c.Expect.Matches) == 0 &&
			len(c.Expect.Headers) == 0 && len(c.Expect.HeaderMatches) == 0 {
			add("%s: a case that asserts nothing about the response is not evidence of anything", where)
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

// Verified reports how many conformance cases were observed against the real
// API, and how many rest on documentation alone. The distinction is the whole
// value of the suite, so it is reported rather than averaged away.
func (r *Recipe) Verified() (observed, documented int) {
	for _, c := range r.Conformance {
		if c.Verified != "" {
			observed++
			continue
		}

		documented++
	}

	return observed, documented
}

// Events returns the webhook event names this Recipe can emit.
func (r *Recipe) Events() []string {
	out := make([]string, len(r.Webhooks.Events))
	copy(out, r.Webhooks.Events)
	sort.Strings(out)

	return out
}

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}

// sortedKeys keeps iteration deterministic, so validation problems always
// appear in the same order.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
