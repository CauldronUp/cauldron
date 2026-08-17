package recipe

import (
	"strings"
	"testing"
)

const minimal = `
recipe: stripe
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
    fields:
      email:
        type: string
        required: true
routes:
  - method: POST
    path: /v1/customers
    resource: customer
    operation: create
`

func parse(t *testing.T, yaml string) (*Recipe, error) {
	t.Helper()

	return Parse([]byte(yaml))
}

// problems asserts that parsing fails and returns the combined message.
func problems(t *testing.T, yaml string) string {
	t.Helper()

	_, err := parse(t, yaml)
	if err == nil {
		t.Fatal("expected validation to fail, but it passed")
	}

	return err.Error()
}

func TestParseMinimalRecipe(t *testing.T) {
	r, err := parse(t, minimal)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if r.Name != "stripe" {
		t.Errorf("Name = %q", r.Name)
	}

	if r.Version != "0.1.0" {
		t.Errorf("Version = %q", r.Version)
	}

	if got := r.Resources["customer"].ID.Prefix; got != "cus_" {
		t.Errorf("id prefix = %q", got)
	}

	if len(r.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(r.Routes))
	}
}

func TestUnknownFieldIsRejected(t *testing.T) {
	got := problems(t, minimal+"\nnonsense: true\n")

	if !strings.Contains(got, "not valid YAML") {
		t.Errorf("expected an unknown-field error, got %q", got)
	}
}

func TestNameMustBeASlug(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "recipe: stripe", "recipe: Stripe Payments", 1))

	if !strings.Contains(got, "lowercase") {
		t.Errorf("got %q", got)
	}
}

func TestVersionMustLookLikeSemver(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "version: 0.1.0", "version: v1", 1))

	if !strings.Contains(got, "1.2.3") {
		t.Errorf("got %q", got)
	}
}

func TestUpstreamAPIIsRequired(t *testing.T) {
	yaml := `
recipe: stripe
version: 0.1.0
resources:
  customer:
    id: {prefix: cus_}
    fields: {email: {type: string}}
routes:
  - {method: POST, path: /v1/customers, resource: customer, operation: create}
`

	got := problems(t, yaml)

	if !strings.Contains(got, "upstream.api") {
		t.Errorf("a Recipe must record which API version it targets; got %q", got)
	}
}

func TestResourceMustDeclareAnIDPrefix(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "      prefix: cus_\n", "", 1))

	if !strings.Contains(got, "id.prefix") {
		t.Errorf("got %q", got)
	}
}

func TestRouteMustReferenceAKnownResource(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "resource: customer\n    operation: create", "resource: invoice\n    operation: create", 1))

	if !strings.Contains(got, `unknown resource "invoice"`) {
		t.Errorf("got %q", got)
	}
}

func TestDuplicateRoutesAreRejected(t *testing.T) {
	yaml := minimal + `  - method: POST
    path: /v1/customers
    resource: customer
    operation: create
`

	got := problems(t, yaml)

	if !strings.Contains(got, "duplicate route") {
		t.Errorf("got %q", got)
	}
}

func TestOperationMustBeKnown(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "operation: create", "operation: upsert", 1))

	if !strings.Contains(got, `operation "upsert"`) {
		t.Errorf("got %q", got)
	}
}

func TestPathMustBeAbsolute(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "path: /v1/customers", "path: v1/customers", 1))

	if !strings.Contains(got, "must start with /") {
		t.Errorf("got %q", got)
	}
}

func TestSigningRequiresAHeader(t *testing.T) {
	yaml := minimal + `
webhooks:
  events: [customer.created]
  signing:
    scheme: hmac-sha256
`

	got := problems(t, yaml)

	if !strings.Contains(got, "signing.header is required") {
		t.Errorf("got %q", got)
	}
}

func TestFixtureMustSeedAKnownResource(t *testing.T) {
	yaml := minimal + `
fixtures:
  small-shop:
    invoice:
      - id: in_1
`

	got := problems(t, yaml)

	if !strings.Contains(got, `unknown resource "invoice"`) {
		t.Errorf("got %q", got)
	}
}

func TestErrorStatusMustBeValid(t *testing.T) {
	yaml := minimal + `
errors:
  rate_limit:
    status: 999
`

	got := problems(t, yaml)

	if !strings.Contains(got, "not a valid HTTP status") {
		t.Errorf("got %q", got)
	}
}

func TestAllProblemsAreReportedTogether(t *testing.T) {
	yaml := `
recipe: Stripe
version: nope
resources: {}
routes: []
`

	got := problems(t, yaml)

	for _, want := range []string{"lowercase", "1.2.3", "upstream.api", "at least one resource", "at least one route"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected every problem at once; %q is missing from:\n%s", want, got)
		}
	}
}

func TestValidationProblemsAreDeterministic(t *testing.T) {
	yaml := `
recipe: x
version: 0.1.0
upstream: {api: "2026-01-01"}
resources:
  zebra: {id: {prefix: z_}, fields: {}}
  alpha: {id: {prefix: a_}, fields: {}}
routes:
  - {method: GET, path: /a, resource: alpha, operation: get}
`

	first := problems(t, yaml)

	for i := 0; i < 10; i++ {
		if got := problems(t, yaml); got != first {
			t.Fatalf("validation output is not deterministic:\n%s\nvs\n%s", got, first)
		}
	}
}

func TestEventsAreSorted(t *testing.T) {
	yaml := minimal + `
webhooks:
  events:
    - customer.updated
    - customer.created
`

	r, err := parse(t, yaml)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	events := r.Events()

	if len(events) != 2 || events[0] != "customer.created" {
		t.Errorf("Events() = %v", events)
	}
}

// An absence is a claim. "Another project's issues are not visible" has no
// positive half, and rejecting it as evidence-free turned away a real case.
func TestAnAbsenceCountsAsAnAssertion(t *testing.T) {
	source := `
recipe: absence
version: 0.1.0
upstream:
  api: "1"
auth:
  scheme: none
resources:
  thing:
    id:
      prefix: thg_
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /things
    resource: thing
    operation: list
conformance:
  - name: nothing leaks from another tenant
    source: https://example.com/docs
    request:
      method: GET
      path: /things
    expect:
      status: 200
      absent:
        - '[0]'
`

	if _, err := Parse([]byte(source)); err != nil {
		t.Fatalf("a case asserting only an absence should be valid: %v", err)
	}
}

func TestACaseAssertingNothingIsStillRejected(t *testing.T) {
	source := `
recipe: empty-claim
version: 0.1.0
upstream:
  api: "1"
auth:
  scheme: none
resources:
  thing:
    id:
      prefix: thg_
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /things
    resource: thing
    operation: list
conformance:
  - name: this claims nothing at all
    source: https://example.com/docs
    request:
      method: GET
      path: /things
    expect:
      status: 200
`

	if _, err := Parse([]byte(source)); err == nil {
		t.Error("a case asserting only a 200 is not evidence of anything")
	}
}

// Only time fields are ever filled in from the sandbox clock, so declaring
// stamped on a string does nothing while reading as though it does. The first
// draft of the Shippo Recipe made exactly that mistake.
func TestStampedOnlyAppliesToTimeFields(t *testing.T) {
	source := `
recipe: stamping
version: 0.1.0
upstream:
  api: "1"
auth:
  scheme: none
resources:
  thing:
    id:
      prefix: thg_
    fields:
      label_url:
        type: string
        stamped: false
routes:
  - method: GET
    path: /things
    resource: thing
    operation: list
`

	_, err := Parse([]byte(source))
	if err == nil {
		t.Fatal("stamped on a string field should be refused")
	}

	if !strings.Contains(err.Error(), "only applies to") {
		t.Errorf("err = %q", err)
	}
}

func TestStampedIsAcceptedOnATimeField(t *testing.T) {
	source := `
recipe: stamping-ok
version: 0.1.0
upstream:
  api: "1"
auth:
  scheme: none
resources:
  thing:
    id:
      prefix: thg_
    fields:
      resolved_at:
        type: datetime
        stamped: false
routes:
  - method: GET
    path: /things
    resource: thing
    operation: list
`

	if _, err := Parse([]byte(source)); err != nil {
		t.Fatalf("stamped on a datetime field should be accepted: %v", err)
	}
}
