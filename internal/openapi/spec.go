// Package openapi reads an OpenAPI 3 description and relates it to a Recipe.
//
// It exists for two jobs, and it is worth being clear about which is which,
// because the difference is the whole reason this package is shaped as it is.
//
// The first job is drafting. A specification carries the mechanical half of a
// Recipe — the paths, the methods, the field names, the types, the status
// codes — and typing that out by hand is work without judgement in it. So
// Draft produces a starting point.
//
// The second job is checking. A specification cannot tell you which state
// lies, which field is absent when its absence is the signal, or which
// success is not a success. That is what a Recipe is for, and no importer
// will ever produce it. What a specification can do is disagree: it knows the
// field is spelled amount_money and not amount, that the endpoint answers 202
// and not 200, that the path has a version segment nobody typed. Check finds
// those, and turns "read from the documentation" into something a machine can
// re-read.
//
// Only the subset that serves those two jobs is modelled. An OpenAPI document
// is a large format and most of it — callbacks, links, discriminators,
// XML shapes, the whole of 3.1's JSON Schema alignment — has no bearing on
// either job, so it is skipped rather than half-supported.
package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Document is the part of an OpenAPI description this package reads.
type Document struct {
	OpenAPI      string                `yaml:"openapi"`
	Swagger      string                `yaml:"swagger"`
	Info         Info                  `yaml:"info"`
	Servers      []Server              `yaml:"servers"`
	ExternalDocs *ExternalDocs         `yaml:"externalDocs"`
	Paths        map[string]PathItem   `yaml:"paths"`
	Components   Components            `yaml:"components"`
	Security     []map[string][]string `yaml:"security"`
}

// Info is the description's own metadata.
type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// Server is one base URL the API is served from.
type Server struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

// ExternalDocs points at the human documentation.
type ExternalDocs struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

// PathItem is the set of operations available on one path.
type PathItem struct {
	// Ref is a path declared in another file. A description may split itself
	// up, and Lob's does: all fifty-eight of its paths are one-line
	// references to resources/<thing>/<thing>.yml. Those files are not
	// fetched, so the path is known to exist and nothing else about it is.
	Ref        string      `yaml:"$ref"`
	Get        *Operation  `yaml:"get"`
	Put        *Operation  `yaml:"put"`
	Post       *Operation  `yaml:"post"`
	Delete     *Operation  `yaml:"delete"`
	Patch      *Operation  `yaml:"patch"`
	Head       *Operation  `yaml:"head"`
	Options    *Operation  `yaml:"options"`
	Parameters []Parameter `yaml:"parameters"`
}

// Operations returns the declared operations by upper-case HTTP method, in a
// stable order so anything built from them is reproducible.
func (p PathItem) Operations() []MethodOperation {
	candidates := []MethodOperation{
		{"GET", p.Get}, {"POST", p.Post}, {"PUT", p.Put},
		{"PATCH", p.Patch}, {"DELETE", p.Delete},
		{"HEAD", p.Head}, {"OPTIONS", p.Options},
	}

	out := make([]MethodOperation, 0, len(candidates))

	for _, candidate := range candidates {
		if candidate.Operation != nil {
			out = append(out, candidate)
		}
	}

	return out
}

// MethodOperation pairs an HTTP method with the operation declared under it.
type MethodOperation struct {
	Method    string
	Operation *Operation
}

// Operation is one method on one path.
type Operation struct {
	OperationID string                `yaml:"operationId"`
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description"`
	Deprecated  bool                  `yaml:"deprecated"`
	Tags        []string              `yaml:"tags"`
	Parameters  []Parameter           `yaml:"parameters"`
	RequestBody *RequestBody          `yaml:"requestBody"`
	Responses   map[string]Response   `yaml:"responses"`
	Security    []map[string][]string `yaml:"security"`
}

// Parameter is a path, query or header parameter.
type Parameter struct {
	Name     string  `yaml:"name"`
	In       string  `yaml:"in"`
	Required bool    `yaml:"required"`
	Schema   *Schema `yaml:"schema"`
	Ref      string  `yaml:"$ref"`
}

// RequestBody is what an operation accepts.
type RequestBody struct {
	Required bool                 `yaml:"required"`
	Content  map[string]MediaType `yaml:"content"`
}

// Response is one declared outcome.
type Response struct {
	Description string               `yaml:"description"`
	Content     map[string]MediaType `yaml:"content"`
	Headers     map[string]Header    `yaml:"headers"`
	Ref         string               `yaml:"$ref"`
}

// Header is a declared response header.
type Header struct {
	Description string  `yaml:"description"`
	Schema      *Schema `yaml:"schema"`
}

// MediaType carries the schema for one content type.
type MediaType struct {
	Schema *Schema `yaml:"schema"`
}

// Components holds the reusable pieces.
type Components struct {
	Schemas         map[string]*Schema        `yaml:"schemas"`
	Responses       map[string]Response       `yaml:"responses"`
	Parameters      map[string]Parameter      `yaml:"parameters"`
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes"`
}

// SecurityScheme is one way of authenticating.
type SecurityScheme struct {
	Type         string `yaml:"type"`
	Scheme       string `yaml:"scheme"`
	BearerFormat string `yaml:"bearerFormat"`
	In           string `yaml:"in"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
}

// Schema is the subset of JSON Schema an OpenAPI description uses for the
// shapes this package cares about.
type Schema struct {
	Ref string `yaml:"$ref"`
	// Type is a word in OpenAPI 3.0 and may be a sequence in 3.1, where
	// ["string", "null"] is how a nullable field is written. Telnyx's
	// published description uses the sequence form and this package refused to
	// read the whole file because of it, which is a poor reason to be unable
	// to check a Recipe.
	Type        SchemaType         `yaml:"type"`
	Format      string             `yaml:"format"`
	Nullable    bool               `yaml:"nullable"`
	Enum        []any              `yaml:"enum"`
	Properties  map[string]*Schema `yaml:"properties"`
	Required    RequiredNames      `yaml:"required"`
	Items       *Schema            `yaml:"items"`
	Description string             `yaml:"description"`
	// AllOf is followed because it is how most descriptions express "this
	// object plus these fields", and ignoring it loses half the properties of
	// anything built that way.
	AllOf []*Schema `yaml:"allOf"`
	// OneOf and AnyOf are read only far enough to notice that a field can be
	// several shapes, which is worth reporting and not worth merging.
	OneOf []*Schema `yaml:"oneOf"`
	AnyOf []*Schema `yaml:"anyOf"`
}

// UnmarshalYAML accepts a boolean where a schema is expected.
//
// JSON Schema says a schema may be true or false: true validates anything,
// false validates nothing. Descriptions use it constantly for
// additionalProperties, and for items when a tuple is closed. Meilisearch's
// has fourteen, and the whole file was unreadable with "cannot unmarshal
// !!bool `false` into openapi.Schema" -- a line number and no hint that a
// perfectly ordinary construct was the cause.
//
// Neither value describes an object with properties, so both read as a schema
// that declares nothing, which is what the rest of this package already treats
// as "the description declining to say".
func (s *Schema) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && node.Tag == "!!bool" {
		*s = Schema{}

		return nil
	}

	// A sequence where a schema belongs is tuple validation: draft-4 JSON
	// Schema lets items hold a list, one schema per position. OpenAPI 3.0
	// dropped it and descriptions carry it anyway. Webflow's has thirty-two,
	// and all 4.4MB of it was refused with "cannot unmarshal !!seq into
	// openapi.plain" repeated down the screen -- over a construct nothing here
	// reads a description for.
	//
	// The first position is taken and the rest dropped. What this package asks
	// of an array is what shape one element has, and a tuple whose positions
	// disagree has no single answer; reading the first is closer than refusing
	// the file.
	if node.Kind == yaml.SequenceNode {
		if len(node.Content) == 0 {
			*s = Schema{}

			return nil
		}

		return node.Content[0].Decode(s)
	}

	type plain Schema

	var raw plain

	if err := node.Decode(&raw); err != nil {
		return err
	}

	*s = Schema(raw)

	return nil
}

// Load reads a description from a file.
//
// JSON is accepted without a separate parser because JSON is a subset of YAML,
// which is the one place that equivalence is worth relying on: an OpenAPI
// description is as likely to be one as the other, and requiring the caller to
// know which would be a worse interface than accepting both.
func Load(path string) (*Document, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	return Parse(raw)
}

// Parse reads a description from bytes.
func Parse(raw []byte) (*Document, error) {
	// The URL a provider's docs give for its description often serves the
	// documentation page instead. Segment, Tradier, ClickUp, Mailchimp and
	// Cal.com all did, and the YAML error for an HTML file -- "mapping values
	// are not allowed in this context" -- sends the reader looking for a
	// syntax error in a file that is not a description at all.
	if looksHTML(raw) {
		return nil, fmt.Errorf("this is an HTML page, not an OpenAPI description; the URL probably serves documentation rather than the description itself")
	}

	var doc Document

	if err := yaml.Unmarshal(raw, &doc); err != nil {
		// YAML 1.2 is a superset of JSON; yaml.v3 implements 1.1, which is
		// not. The difference that matters in practice is the surrogate pair
		// JSON uses for anything outside the basic plane: one thumbs-up in
		// one example of Twilio's description made all hundred and twenty-one
		// of its paths unreadable, and the error said "invalid Unicode
		// character escape code" at a line number in the middle of 1.8MB,
		// which reads like a corrupt download.
		//
		// So a description that really is JSON gets decoded as JSON and
		// handed back through YAML, whose tags the Document is written
		// against. Only on the failing path, so nothing that already parses
		// pays for it, and a file that is neither still reports the YAML
		// error rather than a confusing JSON one.
		converted, jsonErr := viaJSON(raw)
		if jsonErr != nil {
			return nil, fmt.Errorf("not a readable OpenAPI description: %w", err)
		}

		if err := yaml.Unmarshal(converted, &doc); err != nil {
			return nil, fmt.Errorf("not a readable OpenAPI description: %w", err)
		}
	}

	if doc.OpenAPI == "" && doc.Swagger == "" {
		// A document that announces a format this does not read is a
		// permanent fact about the provider, and worth saying as one. An
		// unmarked document might be anything, including a login page,
		// so it gets the vaguer answer it deserves.
		if format := detectForeignFormat(raw); format != "" {
			return nil, &FormatError{Format: format}
		}

		return nil, fmt.Errorf("no openapi or swagger version field, so this is probably not an OpenAPI description")
	}

	// Swagger 2.0 is rewritten rather than refused. Refusing it was
	// refusing real providers -- GitLab, LaunchDarkly, Netlify, Polygon,
	// Postmark, Royal Mail and PDBe publish 2.0 and nothing else, and
	// every one of them was counted as publishing nothing readable,
	// which was true of this reader and not of them.
	if doc.Swagger != "" && doc.OpenAPI == "" {
		converted, err := parseSwagger2(raw)
		if err != nil {
			return nil, &FormatError{Format: "Swagger " + doc.Swagger, Err: err}
		}

		doc = *converted
	}

	if len(doc.Paths) == 0 {
		return nil, fmt.Errorf("the description declares no paths, so there is nothing to read")
	}

	return &doc, nil
}

// parseSwagger2 decodes a 2.0 description, rewrites it, and reads the
// result as the OpenAPI 3 the rest of this package is written against.
//
// The round trip through YAML is deliberate. The alternative is a second
// set of types modelling 2.0, which would have to be kept in step with the
// first set forever; this way the converted document is read by exactly the
// same code as any other, so Draft, Check, Fingerprint and Augment support
// 2.0 without knowing that 2.0 exists.
func parseSwagger2(raw []byte) (*Document, error) {
	var generic map[string]any

	if err := yaml.Unmarshal(raw, &generic); err != nil {
		converted, jsonErr := viaJSON(raw)
		if jsonErr != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(converted, &generic); err != nil {
			return nil, err
		}
	}

	if err := convertSwagger2(generic); err != nil {
		return nil, err
	}

	rewritten, err := yaml.Marshal(generic)
	if err != nil {
		return nil, err
	}

	var doc Document

	if err := yaml.Unmarshal(rewritten, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// Resolve follows a local $ref to the schema it names.
//
// Only local references are followed. A description that points at another
// file or another host is describing something this package cannot see, and
// guessing at it would produce a draft that looks complete and is not, which
// is the failure this whole project is about.
func (d *Document) Resolve(s *Schema) *Schema {
	seen := map[string]bool{}

	for s != nil && s.Ref != "" {
		if seen[s.Ref] {
			return nil
		}

		seen[s.Ref] = true

		name, ok := strings.CutPrefix(s.Ref, "#/components/schemas/")
		if !ok {
			return nil
		}

		s = d.Components.Schemas[name]
	}

	return s
}

// Properties returns a schema's own properties merged with everything its
// allOf members contribute, resolved, and sorted so the result is stable.
//
// Sorted because a map is not, and a draft that changes shape between runs of
// the same command over the same file is not a draft anybody can review.
func (d *Document) Properties(s *Schema) []Property {
	collected := map[string]*Schema{}
	required := map[string]bool{}

	d.collect(s, collected, required, map[*Schema]bool{})

	names := make([]string, 0, len(collected))
	for name := range collected {
		names = append(names, name)
	}

	sort.Strings(names)

	out := make([]Property, 0, len(names))

	for _, name := range names {
		out = append(out, Property{Name: name, Schema: collected[name], Required: required[name]})
	}

	return out
}

func (d *Document) collect(s *Schema, into map[string]*Schema, required map[string]bool, seen map[*Schema]bool) {
	s = d.Resolve(s)
	if s == nil || seen[s] {
		return
	}

	seen[s] = true

	for _, member := range s.AllOf {
		d.collect(member, into, required, seen)
	}

	for name, property := range s.Properties {
		into[name] = property
	}

	for _, name := range s.Required {
		required[name] = true
	}
}

// Property is one field of an object schema.
type Property struct {
	Name     string
	Schema   *Schema
	Required bool
}

// SortedPaths returns the declared paths in a stable order.
func (d *Document) SortedPaths() []string {
	out := make([]string, 0, len(d.Paths))
	for path := range d.Paths {
		out = append(out, path)
	}

	sort.Strings(out)

	return out
}

// Success returns the schema of an operation's first successful response, and
// the status code it belongs to.
func (d *Document) Success(op *Operation) (*Schema, string) {
	for _, code := range sortedCodes(op.Responses) {
		if code < "200" || code >= "300" {
			continue
		}

		response := d.ResolveResponse(op.Responses[code])

		for _, media := range []string{"application/json", "application/hal+json", "application/vnd.api+json"} {
			if content, ok := response.Content[media]; ok && content.Schema != nil {
				return d.Resolve(content.Schema), code
			}
		}

		// A success with no JSON body at all, which is a real answer worth
		// reporting rather than skipping.
		return nil, code
	}

	return nil, ""
}

// Failures returns the non-success status codes an operation declares.
func (d *Document) Failures(op *Operation) []string {
	var out []string

	for _, code := range sortedCodes(op.Responses) {
		if code >= "400" && code < "600" {
			out = append(out, code)
		}
	}

	return out
}

func sortedCodes(responses map[string]Response) []string {
	out := make([]string, 0, len(responses))
	for code := range responses {
		out = append(out, code)
	}

	sort.Strings(out)

	return out
}

// ResolveParameter follows a local $ref to the parameter it names.
func (d *Document) ResolveParameter(p Parameter) Parameter {
	seen := map[string]bool{}

	for p.Ref != "" && !seen[p.Ref] {
		seen[p.Ref] = true

		name, ok := strings.CutPrefix(p.Ref, "#/components/parameters/")
		if !ok {
			return p
		}

		resolved, ok := d.Components.Parameters[name]
		if !ok {
			return p
		}

		p = resolved
	}

	return p
}

// RequiredNames is the list of properties a schema requires.
//
// It tolerates a boolean, which is not what the format says and is what Spotify
// publishes:
//
//	schema:
//	  type: string
//	  format: byte
//	  required: true
//
// That is what required means on a parameter, one level up, and on a schema it
// should name properties. Refusing the document over it was the worse answer:
// 290KB and every path Spotify publishes became unreadable because of one
// boolean in one endpoint, and drift could not watch the provider at all.
//
// A boolean is read as naming nothing, which is what it conveys. The rest of
// the schema survives.
type RequiredNames []string

// UnmarshalYAML accepts the list, and forgives the boolean.
func (r *RequiredNames) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return nil
	}

	var names []string

	if err := node.Decode(&names); err != nil {
		return err
	}

	*r = names

	return nil
}

// SchemaType is a schema's declared type, which may be written as a word or,
// in OpenAPI 3.1, as a sequence of them.
type SchemaType string

// UnmarshalYAML accepts either spelling.
//
// A sequence keeps its first non-null member, because ["string", "null"] means
// a nullable string and the useful half is the string. Nothing here needs to
// know a field is nullable, and pretending the type is "null" would be worse
// than picking.
func (t *SchemaType) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.SequenceNode {
		for _, member := range node.Content {
			if member.Value != "" && member.Value != "null" {
				*t = SchemaType(member.Value)

				return nil
			}
		}

		return nil
	}

	var single string

	if err := node.Decode(&single); err != nil {
		return err
	}

	*t = SchemaType(single)

	return nil
}

// ResolveResponse follows a local $ref to the response it names.
//
// A description that declares a shared response once and references it
// everywhere is common, and Clerk's does exactly that for its delete: without
// following the reference the schema looked absent, which reads as "the
// description declines to say" rather than "the description says plainly".
func (d *Document) ResolveResponse(r Response) Response {
	seen := map[string]bool{}

	for r.Ref != "" && !seen[r.Ref] {
		seen[r.Ref] = true

		name, ok := strings.CutPrefix(r.Ref, "#/components/responses/")
		if !ok {
			return r
		}

		resolved, ok := d.Components.Responses[name]
		if !ok {
			return r
		}

		r = resolved
	}

	return r
}

// viaJSON re-encodes a JSON description as YAML.
//
// The Document is written against yaml tags, so the round trip is what makes
// the JSON readable without a second set of struct tags to keep in step with
// the first.
func viaJSON(raw []byte) ([]byte, error) {
	var any any

	if err := json.Unmarshal(raw, &any); err != nil {
		return nil, err
	}

	return yaml.Marshal(any)
}

// looksHTML reports whether these bytes open like a web page.
func looksHTML(raw []byte) bool {
	head := raw
	if len(head) > 512 {
		head = head[:512]
	}

	trimmed := strings.TrimSpace(string(head))
	lower := strings.ToLower(trimmed)

	return strings.HasPrefix(lower, "<!doctype html") ||
		strings.HasPrefix(lower, "<html") ||
		strings.HasPrefix(lower, "<?xml") && strings.Contains(lower, "<html")
}
