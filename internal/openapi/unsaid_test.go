package openapi

import (
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A description that declares a status nowhere has not contradicted a Recipe
// that answers with it. It has declined to say.
//
// Nineteen of the thirty-one descriptions fetched for this project declare no
// 429 anywhere, and Stripe, Twilio, Slack and Square are among them -- every
// one of which rate limits. Reporting those Recipes as contradicted was
// reporting the fact that rate limits are documented in prose.
const documentsFailures = `
openapi: 3.0.0
info: {title: Failing, version: "1"}
paths:
  /v1/things:
    get:
      responses:
        "200": {description: ok}
        "404": {description: gone}
    post:
      responses:
        "201": {description: made}
        "400": {description: bad}
`

func statusFindings(t *testing.T, r *recipe.Recipe) []Finding {
	t.Helper()

	doc, err := Parse([]byte(documentsFailures))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	return Check(r, doc, "")
}

func TestAStatusDeclaredNowhereIsUnsaidRatherThanContradicted(t *testing.T) {
	r := &recipe.Recipe{
		Name:   "failing",
		Errors: map[string]recipe.Error{"rate_limit": {Status: 429}},
	}

	found := false

	for _, finding := range statusFindings(t, r) {
		if finding.Where != `error "rate_limit"` {
			continue
		}

		found = true

		if finding.Severity != Unsaid {
			t.Errorf("severity is %q, want %q: the description says nothing about 429", finding.Severity, Unsaid)
		}
	}

	if !found {
		t.Error("the undeclared status was not reported at all, which is the other way to be wrong about it")
	}
}

// A status the description does declare is not reported at all, unsaid or
// otherwise.
func TestADeclaredStatusIsNotReported(t *testing.T) {
	r := &recipe.Recipe{
		Name:   "failing",
		Errors: map[string]recipe.Error{"resource_missing": {Status: 404}},
	}

	for _, finding := range statusFindings(t, r) {
		if finding.Where == `error "resource_missing"` {
			t.Errorf("404 is declared on an operation and was reported anyway: %s", finding.What)
		}
	}
}
