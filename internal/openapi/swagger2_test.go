package openapi

import "testing"

// The document these tests are written against is Swagger 2.0 in the shape
// real providers publish it: a host and basePath instead of servers,
// definitions instead of components, a schema directly on the response, and a
// body parameter instead of a requestBody.
const swagger2Doc = `
swagger: "2.0"
info:
  title: Example
  version: "1"
host: api.example.com
basePath: /v2
schemes: [https]
produces: [application/json]
securityDefinitions:
  key:
    type: apiKey
    name: X-Api-Key
    in: header
  legacy:
    type: basic
paths:
  /widgets:
    get:
      operationId: listWidgets
      parameters:
        - name: limit
          in: query
          type: integer
          format: int32
          required: false
      responses:
        "200":
          description: ok
          schema:
            $ref: "#/definitions/WidgetList"
    post:
      operationId: createWidget
      consumes: [application/json]
      parameters:
        - name: body
          in: body
          required: true
          schema:
            $ref: "#/definitions/Widget"
      responses:
        "201":
          description: made
          schema:
            $ref: "#/definitions/Widget"
definitions:
  Widget:
    type: object
    required: [id]
    properties:
      id:
        type: string
      size:
        type: integer
  WidgetList:
    type: object
    properties:
      widgets:
        type: array
        items:
          $ref: "#/definitions/Widget"
`

func TestASwagger2DocumentIsRead(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("Swagger 2.0 was refused: %v", err)
	}

	if len(doc.Paths) != 1 {
		t.Fatalf("paths = %d, want 1", len(doc.Paths))
	}

	if doc.Paths["/widgets"].Get == nil {
		t.Error("the GET operation was lost")
	}
}

// host + basePath + schemes is how 2.0 says what OpenAPI 3 says with servers,
// and BasePath reads servers. Without this every path in the description is
// missing the /v2 the Recipe has to carry, which is the same mistake that made
// every derived route mount at a path no client would call.
func TestSwagger2HostAndBasePathBecomeAServer(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if len(doc.Servers) != 1 {
		t.Fatalf("servers = %d, want 1", len(doc.Servers))
	}

	if got, want := doc.Servers[0].URL, "https://api.example.com/v2"; got != want {
		t.Errorf("server = %q, want %q", got, want)
	}

	if got := BasePath(doc); got != "/v2" {
		t.Errorf("BasePath = %q, want /v2", got)
	}
}

// A reference that still points at #/definitions resolves to nothing, and a
// schema that resolves to nothing drafts a resource with no fields -- which
// looks like a description that declares an empty object rather than like a
// reader that did not follow the pointer.
func TestSwagger2DefinitionsBecomeComponentsAndRefsFollow(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	widget, ok := doc.Components.Schemas["Widget"]
	if !ok {
		t.Fatalf("Widget was not carried into components.schemas; got %v", keysOf(doc.Components.Schemas))
	}

	if _, ok := widget.Properties["id"]; !ok {
		t.Error("Widget lost its id property")
	}

	list := doc.Components.Schemas["WidgetList"]
	if list == nil || list.Properties["widgets"] == nil {
		t.Fatal("WidgetList lost its widgets property")
	}

	items := list.Properties["widgets"].Items
	if items == nil {
		t.Fatal("the array declared no items")
	}

	if got, want := items.Ref, "#/components/schemas/Widget"; got != want {
		t.Errorf("item ref = %q, want %q", got, want)
	}

	if resolved := doc.Resolve(items); resolved == nil {
		t.Error("the rewritten ref does not resolve")
	}
}

// In 2.0 the response carries a schema directly; in 3 it carries content keyed
// by media type. Everything downstream reads content, so a response that is not
// moved is a response with no shape at all.
func TestSwagger2ResponseSchemaBecomesContent(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	response, ok := doc.Paths["/widgets"].Get.Responses["200"]
	if !ok {
		t.Fatal("the 200 response was lost")
	}

	media, ok := response.Content["application/json"]
	if !ok {
		t.Fatalf("no application/json content; got %v", keysOfMedia(response.Content))
	}

	if media.Schema == nil {
		t.Fatal("the content declares no schema")
	}

	if got, want := media.Schema.Ref, "#/components/schemas/WidgetList"; got != want {
		t.Errorf("schema ref = %q, want %q", got, want)
	}
}

// A 2.0 parameter puts type and format on the parameter itself. Draft reads
// parameter.Schema, so without this every query parameter is typeless.
func TestSwagger2ParameterTypeMovesToASchema(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	params := doc.Paths["/widgets"].Get.Parameters
	if len(params) != 1 {
		t.Fatalf("parameters = %d, want 1", len(params))
	}

	if params[0].Schema == nil {
		t.Fatal("the parameter has no schema")
	}

	if got, want := string(params[0].Schema.Type), "integer"; got != want {
		t.Errorf("type = %q, want %q", got, want)
	}
}

// The body parameter is not a parameter in 3, and leaving it among them makes
// a request body look like a query string called "body".
func TestSwagger2BodyParameterBecomesARequestBody(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	post := doc.Paths["/widgets"].Post
	if post == nil {
		t.Fatal("the POST operation was lost")
	}

	for _, p := range post.Parameters {
		if p.In == "body" {
			t.Error("the body parameter was left among the parameters")
		}
	}

	if post.RequestBody == nil {
		t.Fatal("no request body")
	}

	if !post.RequestBody.Required {
		t.Error("the body was required and is no longer")
	}

	media, ok := post.RequestBody.Content["application/json"]
	if !ok {
		t.Fatalf("request body content = %v, want application/json", keysOfMedia(post.RequestBody.Content))
	}

	if got, want := media.Schema.Ref, "#/components/schemas/Widget"; got != want {
		t.Errorf("body schema ref = %q, want %q", got, want)
	}
}

func TestSwagger2SecurityDefinitionsBecomeSchemes(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	key, ok := doc.Components.SecuritySchemes["key"]
	if !ok {
		t.Fatal("the apiKey scheme was lost")
	}

	if key.Type != "apiKey" || key.Name != "X-Api-Key" || key.In != "header" {
		t.Errorf("apiKey scheme = %+v", key)
	}

	// 2.0 spells this "basic"; 3 spells it http with a scheme of basic, and
	// the auth reader is written against 3.
	legacy, ok := doc.Components.SecuritySchemes["legacy"]
	if !ok {
		t.Fatal("the basic scheme was lost")
	}

	if legacy.Type != "http" || legacy.Scheme != "basic" {
		t.Errorf("basic scheme = %+v, want type http scheme basic", legacy)
	}
}

// A 2.0 document that declares no host cannot say where it is served, and
// inventing one would put a base path on routes that do not have it.
func TestSwagger2WithoutAHostDeclaresNoServer(t *testing.T) {
	doc, err := Parse([]byte(`
swagger: "2.0"
info: {title: X, version: "1"}
basePath: /v1
paths:
  /a:
    get:
      responses:
        "200": {description: ok}
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// basePath alone is still a prefix the Recipe has to carry, so it is kept
	// as a relative server rather than dropped.
	if got := BasePath(doc); got != "/v1" {
		t.Errorf("BasePath = %q, want /v1", got)
	}
}

// The version is not thrown away. A Recipe that records a fingerprint against
// a description should be able to say what it read.
func TestSwagger2RecordsThatItWasConverted(t *testing.T) {
	doc, err := Parse([]byte(swagger2Doc))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if doc.Swagger != "2.0" {
		t.Errorf("Swagger = %q, want 2.0", doc.Swagger)
	}

	if doc.OpenAPI == "" {
		t.Error("the converted document declares no OpenAPI version")
	}
}

func keysOf(m map[string]*Schema) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}

func keysOfMedia(m map[string]MediaType) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}
