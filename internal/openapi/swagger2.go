// Reading a Swagger 2.0 description by rewriting it into the OpenAPI 3 this
// package is written against.

package openapi

import "strings"

// convertSwagger2 rewrites a decoded Swagger 2.0 document in place into the
// shape the rest of this package reads.
//
// This exists because refusing 2.0 was refusing real providers. GitLab,
// LaunchDarkly, Netlify, Polygon, Postmark, Royal Mail, PDBe and eBay all
// publish 2.0 and nothing else, and every one of them was counted as a
// provider that publishes nothing readable -- which was true of this reader
// and not true of them. Discovery that found those documents and then threw
// them away would be worse than not looking.
//
// The rewrite is mechanical and deliberately partial. Only the differences
// that change what this package can read are handled: where the API is served,
// where the reusable schemas live, how a response carries a shape, how a
// request body is spelled, and how a parameter declares its type. The rest of
// what changed between 2.0 and 3 -- the callbacks and links that were added,
// the form-data encodings, the collection formats -- is skipped for the same
// reason the 3 reader skips them, which is that neither drafting nor checking
// asks about them.
//
// Working on the decoded map rather than the typed Document is the point. A
// converter written against the structs would have to model 2.0 as a second
// set of types, and the two would drift. This produces a document that is
// OpenAPI 3 and is then read by exactly the same code as any other, so
// everything downstream -- Draft, Check, Fingerprint, Augment -- gets 2.0
// support without knowing 2.0 exists.
func convertSwagger2(doc map[string]any) error {
	version, _ := doc["swagger"].(string)
	if version == "" {
		return nil
	}

	// The declared version is kept. A Recipe that fingerprints a description
	// should be able to say what it actually read, and "this was 2.0" is part
	// of that.
	doc["openapi"] = "3.0.0"

	convertServers(doc)
	convertDefinitions(doc)
	convertSecurityDefinitions(doc)

	// Media types are declared once for the whole document in 2.0 and per
	// response in 3, so the document-level default is what an operation
	// inherits when it says nothing.
	produces := stringsOf(doc["produces"])
	consumes := stringsOf(doc["consumes"])

	paths, _ := doc["paths"].(map[string]any)
	for _, item := range paths {
		pathItem, ok := item.(map[string]any)
		if !ok {
			continue
		}

		convertParameters(pathItem)

		for _, method := range []string{"get", "put", "post", "delete", "patch", "head", "options"} {
			operation, ok := pathItem[method].(map[string]any)
			if !ok {
				continue
			}

			convertOperation(operation, produces, consumes)
		}
	}

	rewriteRefs(doc)

	delete(doc, "definitions")
	delete(doc, "securityDefinitions")
	delete(doc, "host")
	delete(doc, "basePath")
	delete(doc, "schemes")
	delete(doc, "produces")
	delete(doc, "consumes")

	return nil
}

// convertServers turns host + basePath + schemes into the one server URL they
// describe.
//
// Getting this wrong is not a small loss. BasePath reads servers, and a
// description whose prefix is not recovered reports every path the Recipe
// models as a path it does not declare -- and, since #502, mounts every
// derived route at an address no client would call.
func convertServers(doc map[string]any) {
	if _, already := doc["servers"]; already {
		return
	}

	host, _ := doc["host"].(string)
	base, _ := doc["basePath"].(string)
	base = strings.TrimSuffix(base, "/")

	if host == "" {
		// A description that does not say where it is served still says what
		// prefix its paths carry, and the prefix is the half that matters
		// here. A relative server is how OpenAPI 3 spells exactly that.
		if base != "" {
			doc["servers"] = []any{map[string]any{"url": base}}
		}

		return
	}

	scheme := "https"
	if schemes := stringsOf(doc["schemes"]); len(schemes) > 0 {
		scheme = schemes[0]

		// Where a description offers both, https is the one a client should
		// use and the one every recorded case here was taken over.
		for _, s := range schemes {
			if s == "https" {
				scheme = "https"

				break
			}
		}
	}

	doc["servers"] = []any{map[string]any{"url": scheme + "://" + host + base}}
}

func convertDefinitions(doc map[string]any) {
	definitions, ok := doc["definitions"].(map[string]any)
	if !ok || len(definitions) == 0 {
		return
	}

	components := mapAt(doc, "components")

	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		schemas = map[string]any{}
	}

	for name, schema := range definitions {
		// A 3 document that also carries definitions is not something to
		// merge blindly; where both declare a name, the 3 one is the one this
		// reader was written against.
		if _, taken := schemas[name]; !taken {
			schemas[name] = schema
		}
	}

	components["schemas"] = schemas
}

func convertSecurityDefinitions(doc map[string]any) {
	defs, ok := doc["securityDefinitions"].(map[string]any)
	if !ok || len(defs) == 0 {
		return
	}

	components := mapAt(doc, "components")

	schemes, _ := components["securitySchemes"].(map[string]any)
	if schemes == nil {
		schemes = map[string]any{}
	}

	for name, raw := range defs {
		scheme, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		// 2.0 has a "basic" type; 3 spells the same thing as an http scheme,
		// and the auth reader is written against 3.
		if t, _ := scheme["type"].(string); t == "basic" {
			scheme["type"] = "http"
			scheme["scheme"] = "basic"
		}

		if _, taken := schemes[name]; !taken {
			schemes[name] = scheme
		}
	}

	components["securitySchemes"] = schemes
}

func convertOperation(operation map[string]any, produces, consumes []string) {
	if own := stringsOf(operation["produces"]); len(own) > 0 {
		produces = own
	}

	if own := stringsOf(operation["consumes"]); len(own) > 0 {
		consumes = own
	}

	convertParameters(operation)
	convertBodyParameter(operation, consumes)
	convertResponses(operation, produces)

	delete(operation, "produces")
	delete(operation, "consumes")
}

// convertResponses moves the schema off the response and under a media type.
//
// Everything downstream reads response.Content, so a response left with a bare
// schema is a response with no shape at all -- which drafts an endpoint that
// exists and answers nothing.
func convertResponses(operation map[string]any, produces []string) {
	responses, ok := operation["responses"].(map[string]any)
	if !ok {
		return
	}

	for _, raw := range responses {
		response, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		if headers, ok := response["headers"].(map[string]any); ok {
			convertHeaders(headers)
		}

		schema, ok := response["schema"]
		if !ok {
			continue
		}

		delete(response, "schema")

		if _, already := response["content"]; already {
			continue
		}

		content := map[string]any{}
		for _, mediaType := range mediaTypesOr(produces) {
			content[mediaType] = map[string]any{"schema": schema}
		}

		response["content"] = content
	}
}

// convertHeaders gives a declared header the schema wrapper 3 expects, since
// 2.0 puts the type on the header itself.
func convertHeaders(headers map[string]any) {
	for _, raw := range headers {
		header, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		liftTypeToSchema(header)
	}
}

// convertBodyParameter turns the body parameter into a requestBody.
//
// A body parameter left among the parameters is read as a query string called
// "body", which is a claim about the wire that no provider sends.
func convertBodyParameter(operation map[string]any, consumes []string) {
	raw, ok := operation["parameters"].([]any)
	if !ok {
		return
	}

	kept := make([]any, 0, len(raw))

	var (
		body     map[string]any
		formData []map[string]any
	)

	for _, entry := range raw {
		parameter, ok := entry.(map[string]any)
		if !ok {
			kept = append(kept, entry)

			continue
		}

		switch in, _ := parameter["in"].(string); in {
		case "body":
			body = parameter
		case "formData":
			formData = append(formData, parameter)
		default:
			kept = append(kept, entry)
		}
	}

	operation["parameters"] = kept

	if _, already := operation["requestBody"]; already {
		return
	}

	switch {
	case body != nil:
		requestBody := map[string]any{}
		if required, _ := body["required"].(bool); required {
			requestBody["required"] = true
		}

		if description, ok := body["description"].(string); ok && description != "" {
			requestBody["description"] = description
		}

		content := map[string]any{}
		for _, mediaType := range mediaTypesOr(consumes) {
			content[mediaType] = map[string]any{"schema": body["schema"]}
		}

		requestBody["content"] = content
		operation["requestBody"] = requestBody

	case len(formData) > 0:
		// Form fields are separate parameters in 2.0 and properties of one
		// schema in 3. Naming the media type matters more than the shape:
		// a client sending a form to an endpoint declared as JSON is a
		// disagreement worth reporting.
		properties := map[string]any{}

		var required []any

		for _, field := range formData {
			name, _ := field["name"].(string)
			if name == "" {
				continue
			}

			property := map[string]any{}
			for _, key := range []string{"type", "format", "description", "enum", "items"} {
				if value, ok := field[key]; ok {
					property[key] = value
				}
			}

			properties[name] = property

			if isRequired, _ := field["required"].(bool); isRequired {
				required = append(required, name)
			}
		}

		schema := map[string]any{"type": "object", "properties": properties}
		if len(required) > 0 {
			schema["required"] = required
		}

		mediaType := "application/x-www-form-urlencoded"
		for _, candidate := range consumes {
			if strings.HasPrefix(candidate, "multipart/") {
				mediaType = candidate

				break
			}
		}

		operation["requestBody"] = map[string]any{
			"content": map[string]any{mediaType: map[string]any{"schema": schema}},
		}
	}
}

// convertParameters lifts a parameter's type onto a schema, which is where 3
// puts it and where Draft looks for it.
func convertParameters(holder map[string]any) {
	raw, ok := holder["parameters"].([]any)
	if !ok {
		return
	}

	for _, entry := range raw {
		parameter, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		if in, _ := parameter["in"].(string); in == "body" {
			continue
		}

		liftTypeToSchema(parameter)
	}
}

func liftTypeToSchema(holder map[string]any) {
	if _, already := holder["schema"]; already {
		return
	}

	schema := map[string]any{}

	for _, key := range []string{"type", "format", "enum", "items", "default", "maximum", "minimum", "pattern"} {
		if value, ok := holder[key]; ok {
			schema[key] = value

			delete(holder, key)
		}
	}

	if len(schema) > 0 {
		holder["schema"] = schema
	}
}

// rewriteRefs repoints every local reference at where this converter moved the
// thing it names.
//
// A reference left pointing at #/definitions resolves to nothing, and Resolve
// returning nothing is indistinguishable from a description that declares an
// empty object -- so the draft would be a resource with no fields and no
// complaint. That silence is the reason this walks the whole document rather
// than the places references are expected.
func rewriteRefs(node any) {
	switch value := node.(type) {
	case map[string]any:
		if ref, ok := value["$ref"].(string); ok {
			value["$ref"] = rewriteRef(ref)
		}

		for _, child := range value {
			rewriteRefs(child)
		}

	case []any:
		for _, child := range value {
			rewriteRefs(child)
		}
	}
}

func rewriteRef(ref string) string {
	for from, to := range map[string]string{
		"#/definitions/": "#/components/schemas/",
		"#/responses/":   "#/components/responses/",
		"#/parameters/":  "#/components/parameters/",
	} {
		if rest, ok := strings.CutPrefix(ref, from); ok {
			return to + rest
		}
	}

	return ref
}

// mediaTypesOr is the declared media types, or JSON when a description says
// nothing.
//
// Defaulting is a guess, and it is the right one: a description that declares
// no media type is describing an API this project only ever recorded JSON
// from. A response with no content at all would be worse, because it reads as
// an endpoint that answers with nothing.
func mediaTypesOr(declared []string) []string {
	out := make([]string, 0, len(declared))

	for _, mediaType := range declared {
		if mediaType = strings.TrimSpace(mediaType); mediaType != "" {
			out = append(out, mediaType)
		}
	}

	if len(out) == 0 {
		return []string{"application/json"}
	}

	return out
}

func stringsOf(value any) []string {
	raw, ok := value.([]any)
	if !ok {
		return nil
	}

	out := make([]string, 0, len(raw))

	for _, entry := range raw {
		if s, ok := entry.(string); ok {
			out = append(out, s)
		}
	}

	return out
}

func mapAt(doc map[string]any, key string) map[string]any {
	existing, ok := doc[key].(map[string]any)
	if ok {
		return existing
	}

	created := map[string]any{}
	doc[key] = created

	return created
}
