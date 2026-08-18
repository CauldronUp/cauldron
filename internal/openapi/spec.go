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
	Ref         string             `yaml:"$ref"`
	Type        string             `yaml:"type"`
	Format      string             `yaml:"format"`
	Nullable    bool               `yaml:"nullable"`
	Enum        []any              `yaml:"enum"`
	Properties  map[string]*Schema `yaml:"properties"`
	Required    []string           `yaml:"required"`
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
	var doc Document

	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("not a readable OpenAPI description: %w", err)
	}

	if doc.OpenAPI == "" && doc.Swagger == "" {
		return nil, fmt.Errorf("no openapi or swagger version field, so this is probably not an OpenAPI description")
	}

	if doc.Swagger != "" && doc.OpenAPI == "" {
		return nil, fmt.Errorf("this is Swagger %s, and only OpenAPI 3 is read; convert it first", doc.Swagger)
	}

	if len(doc.Paths) == 0 {
		return nil, fmt.Errorf("the description declares no paths, so there is nothing to read")
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

		response := op.Responses[code]

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
