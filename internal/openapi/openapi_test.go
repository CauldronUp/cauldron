package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

const widgetSpec = `
openapi: 3.0.3
info:
  title: Widget Service
  version: "2026-01-01"
externalDocs:
  url: https://docs.widget.test
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
  schemas:
    Widget:
      type: object
      required: [name]
      properties:
        id: {type: string}
        name: {type: string}
        status:
          type: string
          enum: [pending, active, retired]
        price_cents: {type: integer}
        created_at: {type: string, format: date-time}
        tags:
          type: array
          items: {type: string}
    WidgetList:
      type: object
      properties:
        data:
          type: array
          items: {$ref: '#/components/schemas/Widget'}
        has_more: {type: boolean}
paths:
  /v1/widgets:
    get:
      responses:
        "200":
          description: A list of widgets
          content:
            application/json:
              schema: {$ref: '#/components/schemas/WidgetList'}
        "429":
          description: Too many requests
    post:
      responses:
        "201":
          description: The created widget
          content:
            application/json:
              schema: {$ref: '#/components/schemas/Widget'}
        "400":
          description: Invalid request
  /v1/widgets/{widgetId}:
    get:
      responses:
        "200":
          description: A widget
          content:
            application/json:
              schema: {$ref: '#/components/schemas/Widget'}
        "404":
          description: No such widget
    delete:
      responses:
        "204":
          description: Deleted
`

func parseSpec(t *testing.T, raw string) *Document {
	t.Helper()

	doc, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("parsing the description: %v", err)
	}

	return doc
}

// The one thing a generated draft absolutely has to do is load. A draft that
// does not parse is worse than no draft: it looks like progress and it is a
// file somebody now has to debug before they can start the work the generator
// was supposed to save them.
func TestADraftIsAValidRecipe(t *testing.T) {
	drafted := Draft(parseSpec(t, widgetSpec), "widget")

	r, err := recipe.Parse([]byte(drafted))
	if err != nil {
		t.Fatalf("the draft does not load as a Recipe: %v\n\n%s", err, drafted)
	}

	if r.Name != "widget" {
		t.Errorf("name = %q", r.Name)
	}

	if r.Auth.Scheme != "bearer" {
		t.Errorf("auth scheme = %q, want the one the description declares", r.Auth.Scheme)
	}

	// Four operations, mapped by the only signal available: whether the path
	// ends in a parameter.
	kinds := map[string]string{}
	for _, route := range r.Routes {
		kinds[route.Method+" "+route.Path] = route.Operation
	}

	for path, want := range map[string]string{
		"GET /v1/widgets":         "list",
		"POST /v1/widgets":        "create",
		"GET /v1/widgets/{id}":    "get",
		"DELETE /v1/widgets/{id}": "delete",
	} {
		if got := kinds[path]; got != want {
			t.Errorf("%s = %q, want %q", path, got, want)
		}
	}

	widget, ok := r.Resources["widget"]
	if !ok {
		t.Fatalf("no widget resource: %v", r.Resources)
	}

	if widget.Fields["price_cents"].Type != "integer" {
		t.Errorf("price_cents type = %q", widget.Fields["price_cents"].Type)
	}

	if !widget.Fields["name"].Required {
		t.Error("name is required in the description and not in the draft")
	}

	// id is the identifier, not a field, so declaring it as both would make
	// the emulator send it twice.
	if _, present := widget.Fields["id"]; present {
		t.Error("id was declared as a field as well as an identifier")
	}
}

// The draft has to be visibly unfinished, because the failure mode of a
// generator like this is somebody shipping its output. A file with a green
// suite and nothing in it would be exactly the confident, empty artefact this
// project argues against.
func TestADraftAdmitsItIsNotFinished(t *testing.T) {
	drafted := Draft(parseSpec(t, widgetSpec), "widget")

	if !strings.HasPrefix(drafted, "# DRAFT ") {
		t.Errorf("the first line does not say DRAFT: %q", firstLineOf(drafted))
	}

	r, err := recipe.Parse([]byte(drafted))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if len(r.Conformance) != 0 {
		t.Errorf("the draft carries %d conformance cases, and a description cannot justify one", len(r.Conformance))
	}

	// Exactly one fixture, and it is empty. A generated fixture would be a
	// claim about a plausible account that nothing supports.
	if len(r.Fixtures) != 1 {
		t.Errorf("fixtures = %v, want only the empty one", r.Fixtures)
	}

	for _, phrase := range []string{
		"What does this API lie about",
		"Which endpoints must not be here at all",
		"states its gaps",
	} {
		if !strings.Contains(drafted, phrase) {
			t.Errorf("the header does not ask %q", phrase)
		}
	}
}

// Check exists to find the mechanical disagreements a conformance suite cannot,
// because a suite written by the same person on the same day asserts whatever
// the Recipe produced.
func TestCheckFindsWhatADescriptionContradicts(t *testing.T) {
	doc := parseSpec(t, widgetSpec)

	wrong := `
recipe: widget
capability: commerce
version: 0.1.0
upstream:
  api: "2026-01-01"
resources:
  widget:
    id:
      style: opaque
    fields:
      name:
        type: string
      colour:
        type: string
routes:
  - method: GET
    path: /v1/widgets/{id}
    resource: widget
    operation: get
  - method: GET
    path: /v1/gadgets
    resource: widget
    operation: list
  - method: PATCH
    path: /v1/widgets/{id}
    resource: widget
    operation: update
  - method: POST
    path: /v1/widgets
    resource: widget
    operation: create
    status: 200
errors:
  teapot:
    status: 418
    message: "No"
conformance:
  - name: a widget has a name
    source: https://docs.widget.test
    request:
      method: GET
      path: /v1/widgets/w1
    expect:
      status: 404
      body:
        message: "Not found"
`

	r, err := recipe.Parse([]byte(wrong))
	if err != nil {
		t.Fatalf("parsing the deliberately wrong Recipe: %v", err)
	}

	var said []string

	// Unsaid as well as Disagrees. A status the description declares nowhere
	// is still reported -- it is reported as something the description has no
	// opinion about rather than as a contradiction, and this test is about
	// what gets said rather than under which heading.
	for _, finding := range Check(r, doc, "") {
		if finding.Severity == Disagrees || finding.Severity == Unsaid {
			said = append(said, finding.Where+": "+finding.What)
		}
	}

	joined := strings.Join(said, "\n")

	for _, want := range []string{
		// A path the description does not have.
		"/v1/gadgets",
		// A method that path does not take.
		"PATCH",
		// A status the operation does not answer with.
		"answers 200 and the description declares 201",
		// A field no schema declares.
		`"colour"`,
		// A failure status nothing declares, which is reported as unsaid
		// rather than as a contradiction.
		"418",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("Check did not report %q\n\nwhat it did report:\n%s", want, joined)
		}
	}

	// And it must not report the things that are right, or the report is
	// noise and nobody reads the next one. Checked against each finding's own
	// subject rather than the whole report, because a correct finding may
	// perfectly well name a route that is fine: the colour field is reported
	// against every operation that returns a widget, and one of those is the
	// route that is otherwise correct.
	for _, finding := range Check(r, doc, "") {
		if finding.Severity != Disagrees {
			continue
		}

		for _, fine := range []string{`field "name"`, "route 1 (GET /v1/widgets/{id})"} {
			if strings.Contains(finding.Where, fine) {
				t.Errorf("Check reported %s, which the description agrees with: %s", finding.Where, finding.What)
			}
		}
	}
}

// A Recipe that models less than an API is doing what every Recipe here does
// on purpose, so that has to read differently from getting something wrong.
func TestCheckSeparatesModellingLessFromModellingWrongly(t *testing.T) {
	doc := parseSpec(t, widgetSpec)

	partial := `
recipe: widget
capability: commerce
version: 0.1.0
upstream:
  api: "2026-01-01"
resources:
  widget:
    id:
      style: opaque
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/widgets/{id}
    resource: widget
    operation: get
conformance:
  - name: a widget has a name
    source: https://docs.widget.test
    request:
      method: GET
      path: /v1/widgets/w1
    expect:
      status: 404
      body:
        message: "Not found"
`

	r, err := recipe.Parse([]byte(partial))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var disagreements, omissions int

	for _, finding := range Check(r, doc, "") {
		if finding.Severity == Disagrees {
			disagreements++
			continue
		}

		omissions++
	}

	if disagreements != 0 {
		t.Errorf("a Recipe that models one endpoint correctly produced %d disagreements", disagreements)
	}

	// The collection endpoints it does not model.
	if omissions == 0 {
		t.Error("nothing was reported as unmodelled, and this Recipe models one path of two")
	}
}

// A path prefix the description leaves out is common enough to be worth a
// flag: plenty of descriptions put the version in the server URL and plenty
// of Recipes put it in the path.
func TestCheckHonoursABasePath(t *testing.T) {
	doc := parseSpec(t, strings.ReplaceAll(widgetSpec, "/v1/widgets", "/widgets"))

	r, err := recipe.Parse([]byte(`
recipe: widget
capability: commerce
version: 0.1.0
upstream:
  api: "2026-01-01"
resources:
  widget:
    id:
      style: opaque
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /v1/widgets/{id}
    resource: widget
    operation: get
conformance:
  - name: a widget has a name
    source: https://docs.widget.test
    request:
      method: GET
      path: /v1/widgets/w1
    expect:
      status: 404
      body:
        message: "Not found"
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if findings := disagreementsOf(Check(r, doc, "")); len(findings) == 0 {
		t.Error("without a base path the version segment should not match")
	}

	if findings := disagreementsOf(Check(r, doc, "/v1")); len(findings) != 0 {
		t.Errorf("with the base path it should match: %v", findings)
	}
}

func disagreementsOf(findings []Finding) []string {
	var out []string

	for _, finding := range findings {
		if finding.Severity == Disagrees {
			out = append(out, finding.Where+": "+finding.What)
		}
	}

	return out
}

func firstLineOf(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}

	return s
}

// A description this package cannot read has to say so rather than produce an
// empty draft, which would look like a provider with no endpoints.
func TestParseRefusesWhatItCannotRead(t *testing.T) {
	for name, raw := range map[string]string{
		"not a description": "hello: world\n",
		"swagger 2":         "swagger: \"2.0\"\ninfo: {title: x, version: y}\npaths: {/a: {get: {responses: {}}}}\n",
		"no paths":          "openapi: 3.0.0\ninfo: {title: x, version: y}\npaths: {}\n",
	} {
		if _, err := Parse([]byte(raw)); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}
