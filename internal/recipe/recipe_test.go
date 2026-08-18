package recipe

import (
	"strings"
	"testing"
)

const minimal = `
recipe: stripe
capability: payments
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
capability: payments
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
capability: payments
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
capability: payments
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
capability: payments
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
capability: payments
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
capability: payments
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

// A case whose every claim repeats a value it sent proves nothing. A create
// that posts {"name": "x"} and asserts {"name": "x"} is testing that the
// request survived the round trip, which every fake does by construction.
//
// Four of these were written before anything checked, and each was only found
// by deliberately breaking the Recipe and noticing the case still passed. One
// could not be salvaged at all and had to be replaced. The rule is here so the
// fifth is caught by the validator rather than by somebody remembering.
func TestACaseThatOnlyEchoesItsRequestIsRejected(t *testing.T) {
	echo := Case{
		Name:    "a create returns what it was given",
		Source:  "https://example.com/docs",
		Request: Request{Method: "POST", Path: "/things", JSON: map[string]any{"name": "Ada"}},
		Expect:  Expectation{Status: 201, Body: map[string]any{"name": "Ada"}},
	}

	if !echoesOnly(echo) {
		t.Error("a body that only repeats the request should be rejected")
	}

	// One claim the emulator decides is enough, because then something real is
	// under test.
	withGenerated := echo
	withGenerated.Expect.Matches = map[string]string{"id": "^thing_"}

	if echoesOnly(withGenerated) {
		t.Error("a matches claim cannot be satisfied by echoing, so the case stands")
	}

	withAbsence := echo
	withAbsence.Expect.Absent = []string{"deleted_at"}

	if echoesOnly(withAbsence) {
		t.Error("an absence cannot be satisfied by echoing either")
	}

	// A body claim about a field the request never mentioned is the ordinary
	// case and must not be flagged.
	withDecision := echo
	withDecision.Expect.Body = map[string]any{"name": "Ada", "status": "active"}

	if echoesOnly(withDecision) {
		t.Error("a field the request did not send is the emulator's answer")
	}

	// A read sends nothing, so nothing can be an echo.
	read := Case{
		Name:    "a read",
		Request: Request{Method: "GET", Path: "/things/1"},
		Expect:  Expectation{Status: 200, Body: map[string]any{"name": "Ada"}},
	}

	if echoesOnly(read) {
		t.Error("a request with no body cannot echo anything")
	}
}

// Form values arrive as text, so a form sending "2" and a body asserting the
// number 2 is still an echo. Requiring the YAML types to match would let the
// weakest cases through on a technicality.
func TestAnEchoIsCaughtAcrossFormEncoding(t *testing.T) {
	c := Case{
		Name:    "a form create",
		Request: Request{Method: "POST", Path: "/things", Form: map[string]string{"quantity": "2"}},
		Expect:  Expectation{Status: 200, Body: map[string]any{"quantity": 2}},
	}

	if !echoesOnly(c) {
		t.Error("a form string and a body number are the same claim")
	}
}

// A field name the Recipe chooses is only a claim if a case asserts it where
// the value exists. Asserting its absence on a last page holds whatever the
// field happens to be called, so a Recipe could declare a cursor, have it
// renamed to anything at all, and no case would notice.
//
// Twenty-one of these shipped before anything checked. The check is loose on
// purpose: a nested claim pins a name as well as a dotted one does, and what
// it refuses is silence.
func TestADeclaredFieldNameMustBeAsserted(t *testing.T) {
	// A dotted path in matches.
	dotted := []Case{{Expect: Expectation{Matches: map[string]string{"meta.next_page": "."}}}}
	if !assertsName(dotted, "meta.next_page") {
		t.Error("a dotted matches path should pin the name")
	}

	// The same claim written as nested maps in body.
	nested := []Case{{Expect: Expectation{Body: map[string]any{
		"meta": map[string]any{"pagination": map[string]any{"next": "abc"}},
	}}}}
	if !assertsName(nested, "meta.pagination.next") {
		t.Error("a nested body claim should pin the name too")
	}

	// An index in the path must not hide the field.
	indexed := []Case{{Expect: Expectation{Body: map[string]any{"items[0].next_cursor": "abc"}}}}
	if !assertsName(indexed, "next_cursor") {
		t.Error("an indexed path should still pin the name")
	}

	// The bug this exists for: only an absence, which holds whatever the
	// field is called.
	absenceOnly := []Case{{Expect: Expectation{
		Body:   map[string]any{"id": "1"},
		Absent: []string{"next_cursor"},
	}}}
	if assertsName(absenceOnly, "next_cursor") {
		t.Error("an absence does not pin a name, because it holds under any name")
	}

	// And a Recipe with no cases at all cannot have pinned anything.
	if assertsName(nil, "next_cursor") {
		t.Error("no cases cannot assert a name")
	}
}

// A path segment the runtime cannot walk used to become a literal key. A field
// nested under "to[0]" was emitted under a property actually spelled "to[0]" —
// a shape no provider sends, produced in silence, and invisible to a
// conformance suite with no case naming it.
//
// That is the third time a key like this has shipped in a different disguise,
// which is why it is a rule now rather than a habit. An index the runtime does
// honour must still be accepted, or the rule would forbid the fix as well as
// the bug.
func TestAnInPathSegmentMustBeSomethingTheRuntimeCanWalk(t *testing.T) {
	nested := func(in string) string {
		return strings.Replace(minimal, `      email:
        type: string
        required: true`, `      email:
        type: string
        required: true
      recipient:
        type: string
        in: `+in+`
        as: phoneNumber`, 1)
	}

	// Hyphens are ordinary in JSON keys and this rule refused the first one it
	// ever met: npm sends dist-tags. A pattern written against the Recipes
	// that happened to exist will do that.
	for _, in := range []string{"to[0]", "to", "message.to[0]", "a.b.c", "to[0][1]", "dist-tags", "x-request-id", "a.dist-tags[0]"} {
		if _, err := parse(t, nested(in)); err != nil {
			t.Errorf("in: %q was refused: %v", in, err)
		}
	}

	for _, in := range []string{"to[]", "to[a]", "to[-1]", "to.", "to[0", "0to", "to bar", "-to", "to.bar."} {
		_, err := parse(t, nested(in))
		if err == nil {
			t.Errorf("in: %q was accepted, and would be sent as a property spelled that", in)
			continue
		}

		if !strings.Contains(err.Error(), "literally spelled") {
			t.Errorf("in: %q was refused for the wrong reason: %v", in, err)
		}
	}
}

// A capability is a fixed word, not a free string.
//
// The value of a category is that two people reaching for it independently
// land on the same one. A free string gives you "payments", "payment",
// "billing" and "money" within a month, and then the grouping is worth less
// than no grouping at all, because it looks authoritative and is not.
func TestACapabilityMustComeFromTheList(t *testing.T) {
	with := func(capability string) string {
		if capability == "" {
			return strings.Replace(minimal, "capability: payments\n", "", 1)
		}

		return strings.Replace(minimal, "capability: payments", "capability: "+capability, 1)
	}

	if _, err := parse(t, with("payments")); err != nil {
		t.Errorf("a declared capability was refused: %v", err)
	}

	for _, wrong := range []string{"", "payment", "billing", "Payments", "money", "misc"} {
		_, err := parse(t, with(wrong))
		if err == nil {
			t.Errorf("capability %q was accepted", wrong)
			continue
		}

		if !strings.Contains(err.Error(), "capability") {
			t.Errorf("capability %q was refused for the wrong reason: %v", wrong, err)
		}
	}
}

// Every bundled Recipe declares one, or the grouping has holes in exactly the
// places somebody is looking.
func TestEveryBundledRecipeSaysWhatItDoes(t *testing.T) {
	summaries, err := Summarise()
	if err != nil {
		t.Fatalf("summarise: %v", err)
	}

	for _, s := range summaries {
		if s.Capability == "" {
			t.Errorf("%s does not say what it does", s.Name)
		}
	}
}

// A resource whose identifier never reaches the wire is defensible when it is
// fetched one at a time by a key that lives in the path. On a collection it is
// not: a page of records with nothing to tell them apart cannot be what any
// provider sends, because nobody could address the second one.
func TestAHiddenIdentifierIsRefusedOnACollection(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
      field: "-"
    fields:
      email:
        type: string
routes:
  - method: GET
    path: /v1/customers
    resource: customer
    operation: list
`

	got := problems(t, yaml)

	if !strings.Contains(got, "nothing can address") {
		t.Errorf("got %q", got)
	}
}

func TestAHiddenIdentifierIsAcceptedOnASingleFetch(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
      field: "-"
    fields:
      email:
        type: string
routes:
  - method: GET
    path: /v1/customers/{id}
    resource: customer
    operation: get
`

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("a hidden identifier on a get should be allowed: %v", err)
	}
}

// A constant is stamped onto the record and trimmed with everything else, so a
// route that answers with one has to be able to name it. Braze stamps
// "message": "success" on every object and it is the only field a Braze client
// reliably checks, because a failure puts its prose in the same place.
func TestReturnsMayNameAConstant(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
    constants:
      message: success
    fields:
      email:
        type: string
routes:
  - method: POST
    path: /v1/customers
    resource: customer
    operation: create
    returns: [message, id]
`

	if _, err := parse(t, yaml); err != nil {
		t.Errorf("returns should accept a declared constant: %v", err)
	}
}

func TestReturnsStillRefusesAnUndeclaredName(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
    constants:
      message: success
    fields:
      email:
        type: string
routes:
  - method: POST
    path: /v1/customers
    resource: customer
    operation: create
    returns: [nickname]
`

	got := problems(t, yaml)

	if !strings.Contains(got, "not a field on resource") {
		t.Errorf("got %q", got)
	}
}

// The identifier is held as "id" whatever the provider calls it on the wire,
// and returns names the record's own keys because the trim runs before the
// rename. Writing the wire name is the obvious mistake, so the message says
// what to write instead of reporting an unknown field.
func TestReturnsNamingTheWireIdentifierSaysSo(t *testing.T) {
	yaml := `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  customer:
    id:
      prefix: cus_
      length: 14
      field: customerId
    fields:
      email:
        type: string
routes:
  - method: POST
    path: /v1/customers
    resource: customer
    operation: create
    returns: [customerId]
`

	got := problems(t, yaml)

	if !strings.Contains(got, "what the identifier is called on the wire") {
		t.Errorf("got %q", got)
	}
}

// A default that excludes nothing is a default nobody can see. It is the same
// rule as a scope: the behaviour is only described when a fixture holds a
// record the filter must hide, because otherwise the listing looks identical
// with the filter and without it, and a mutation removing it passes.
const filterBase = `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  invoice:
    id:
      prefix: in_
      length: 14
    fields:
      status:
        type: string
        default: open
fixtures:
  small:
    invoice:
      - id: in_0000000000001
        status: open
routes:
  - method: GET
    path: /v1/invoices
    resource: invoice
    operation: list
    filters:
      - param: status
        field: status
        default: open
`

// The same Recipe with a paid invoice in the fixture, which is the record the
// default has to hide for the default to be describable at all.
const filterWithExcluded = `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  invoice:
    id:
      prefix: in_
      length: 14
    fields:
      status:
        type: string
        default: open
fixtures:
  small:
    invoice:
      - id: in_0000000000001
        status: open
      - id: in_0000000000002
        status: paid
routes:
  - method: GET
    path: /v1/invoices
    resource: invoice
    operation: list
    filters:
      - param: status
        field: status
        default: open
`

func TestADefaultedFilterNeedsSomethingToHide(t *testing.T) {
	got := problems(t, filterBase)

	if !strings.Contains(got, "the default would hide") {
		t.Errorf("got %q", got)
	}
}

func TestADefaultedFilterIsAcceptedWhenAFixtureExcludesSomething(t *testing.T) {
	if _, err := parse(t, filterWithExcluded); err != nil {
		t.Errorf("a filter with something to hide should be allowed: %v", err)
	}
}

func TestAFilterMustNameADeclaredField(t *testing.T) {
	got := problems(t, strings.Replace(filterWithExcluded, "field: status", "field: nickname", 1))

	if !strings.Contains(got, "not a field on resource") {
		t.Errorf("got %q", got)
	}
}

func TestAnEscapeValueEqualToTheDefaultIsRefused(t *testing.T) {
	got := problems(t, strings.Replace(filterWithExcluded,
		"        field: status\n        default: open", "        field: status\n        default: open\n        all: open", 1))

	if !strings.Contains(got, "never applies") {
		t.Errorf("got %q", got)
	}
}

func TestFiltersOnlyNarrowAListing(t *testing.T) {
	got := problems(t, strings.Replace(
		strings.Replace(filterWithExcluded, "path: /v1/invoices", "path: /v1/invoices/{id}", 1),
		"operation: list", "operation: get", 1))

	if !strings.Contains(got, "only narrow a listing") {
		t.Errorf("got %q", got)
	}
}

// A filter value that expands to nothing would hide every record, and one that
// is also the escape value contradicts itself. Both are the sort of mistake
// that reads as working.
const filterBucket = `
recipe: stripe
capability: payments
version: 0.1.0
upstream:
  api: "2026-06-30"
resources:
  invoice:
    id:
      prefix: in_
      length: 14
    fields:
      status:
        type: string
        default: open
fixtures:
  small:
    invoice:
      - id: in_0000000000001
        status: open
      - id: in_0000000000002
        status: paid
routes:
  - method: GET
    path: /v1/invoices
    resource: invoice
    operation: list
    filters:
      - param: status
        field: status
        default: unpaid
        all: every
        values:
          unpaid: [open, past_due]
`

func TestABucketDefaultIsCheckedAgainstItsMembers(t *testing.T) {
	// "unpaid" is not a status any invoice holds, so checking the word against
	// the fixture would report that nothing is hidden when the paid invoice
	// plainly is.
	if _, err := parse(t, filterBucket); err != nil {
		t.Errorf("a bucket default with something to hide should be allowed: %v", err)
	}
}

func TestABucketDefaultThatHidesNothingIsRefused(t *testing.T) {
	got := problems(t, strings.Replace(filterBucket, "unpaid: [open, past_due]", "unpaid: [open, paid]", 1))

	if !strings.Contains(got, "the default would hide") {
		t.Errorf("got %q", got)
	}
}

func TestAnEmptyBucketIsRefused(t *testing.T) {
	got := problems(t, strings.Replace(filterBucket, "unpaid: [open, past_due]", "unpaid: []", 1))

	if !strings.Contains(got, "expands") {
		t.Errorf("got %q", got)
	}
}

func TestABucketNamedLikeTheEscapeValueIsRefused(t *testing.T) {
	got := problems(t, strings.Replace(filterBucket, "unpaid: [open, past_due]",
		"unpaid: [open, past_due]\n          every: [open, paid]", 1))

	if !strings.Contains(got, "turns the filter off") {
		t.Errorf("got %q", got)
	}
}

// null_when_unset says the key is sent with no value in it. A field that also
// has a default, or is required, is never unset, so the declaration would
// describe a case that cannot happen.
func TestNullWhenUnsetConflictsWithADefault(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "        required: true",
		"        required: false\n        default: nobody@example.com\n        null_when_unset: true", 1))

	if !strings.Contains(got, "never unset") {
		t.Errorf("got %q", got)
	}
}

func TestNullWhenUnsetConflictsWithRequired(t *testing.T) {
	got := problems(t, strings.Replace(minimal, "        required: true",
		"        required: true\n        null_when_unset: true", 1))

	if !strings.Contains(got, "never unset") {
		t.Errorf("got %q", got)
	}
}
