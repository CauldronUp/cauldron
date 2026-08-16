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
	List ListResponse `yaml:"list"`
}

// ListResponse describes how a collection is returned.
type ListResponse struct {
	// Style is one of: envelope (Stripe), bare (GitHub), wrapped (Shopify).
	// Empty means envelope, which keeps existing Recipes working.
	Style string `yaml:"style"`
	// Key is the wrapping property name, required when style is wrapped.
	Key string `yaml:"key"`
}

// Resource is an object type the provider exposes.
type Resource struct {
	ID     ID               `yaml:"id"`
	Fields map[string]Field `yaml:"fields"`
}

// ID describes how the provider mints identifiers. Getting this right matters
// more than it looks: applications routinely parse or prefix-match IDs.
type ID struct {
	// Style is one of: prefixed (cus_abc123) or numeric (1, 2, 3).
	// Empty means prefixed.
	Style  string `yaml:"style"`
	Prefix string `yaml:"prefix"`
	Length int    `yaml:"length"`
}

// Field is a single attribute on a resource.
type Field struct {
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Default  any    `yaml:"default"`
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
	Scope      []string   `yaml:"scope"`
	Pagination Pagination `yaml:"pagination"`
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
	Status  int               `yaml:"status"`
	Code    string            `yaml:"code"`
	Message string            `yaml:"message"`
	Headers map[string]string `yaml:"headers"`
}

// Fixture is a named seed dataset: resource name to a list of records.
type Fixture map[string][]map[string]any

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	validSchemes    = []string{"bearer", "basic", "header", "none"}
	validOperations = []string{"create", "get", "list", "update", "delete"}
	validPagination = []string{"", "cursor", "offset", "page"}
	validSigning    = []string{"", "none", "hmac-sha256"}
	validListStyles = []string{"", "envelope", "bare", "wrapped"}
	validIDStyles   = []string{"", "prefixed", "numeric"}
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
			add("resource %q has id.style %q, which must be prefixed or numeric", name, resource.ID.Style)
		}

		if resource.ID.Style != "numeric" && resource.ID.Prefix == "" {
			add("resource %q must declare an id.prefix — applications routinely prefix-match identifiers", name)
		}

		if len(resource.Fields) == 0 {
			add("resource %q has no fields", name)
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
	}

	if !contains(validListStyles, r.Responses.List.Style) {
		add("responses.list.style %q must be one of %s", r.Responses.List.Style, strings.Join(validListStyles[1:], ", "))
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key == "" {
		add("responses.list.key is required when the list style is wrapped")
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

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
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
