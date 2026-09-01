package recipe

import (
	"strings"
	"testing"
)

func rawOnlyRecipe(raw *RawBody) *Recipe {
	return &Recipe{
		Name:       "example",
		Capability: "docs",
		Version:    "0.1.0",
		Upstream:   Upstream{API: "v1"},
		Auth:       Auth{Scheme: "none"},
		Routes:     []Route{{Method: "GET", Path: "/api/query", Raw: raw}},
		Fixtures:   map[string]Fixture{"empty": {}},
	}
}

// The rule this replaces produced the thing it was meant to prevent. arXiv
// sends Atom and nothing else, so its Recipe declared a resource purely to
// satisfy the schema -- one that no route could reach, because reaching it
// would have rendered an XML feed as a JSON object. A declaration written to
// satisfy a validator and understood by nothing is worse than its absence.
func TestARecipeThatOnlySendsBytesOwesNoResource(t *testing.T) {
	r := rawOnlyRecipe(&RawBody{ContentType: "application/atom+xml", Text: "<feed/>"})

	if err := r.Validate(); err != nil {
		t.Errorf("a Recipe whose every route is raw should validate: %v", err)
	}
}

// The relaxation is narrow on purpose. One ordinary route and the resource is
// owed again, or the rule would be gone rather than qualified.
func TestOneOrdinaryRouteBringsTheResourceRuleBack(t *testing.T) {
	r := rawOnlyRecipe(&RawBody{ContentType: "text/plain", Text: "OK"})
	r.Routes = append(r.Routes, Route{Method: "GET", Path: "/things", Operation: "list", Resource: "thing"})

	err := r.Validate()
	if err == nil {
		t.Fatal("a Recipe with an ordinary route and no resources should not validate")
	}

	if !strings.Contains(err.Error(), "at least one resource is required") {
		t.Errorf("error does not mention the missing resource: %v", err)
	}
}

// A raw route renders no record, so an operation or a resource on one is a
// declaration that cannot mean anything -- the same reason an error route
// refuses both.
func TestARawRouteRefusesAnOperationOrResource(t *testing.T) {
	for _, c := range []struct {
		name  string
		apply func(*Route)
	}{
		{"an operation", func(rt *Route) { rt.Operation = "list" }},
		{"a resource", func(rt *Route) { rt.Resource = "thing" }},
	} {
		r := rawOnlyRecipe(&RawBody{ContentType: "text/plain", Text: "OK"})
		c.apply(&r.Routes[0])

		err := r.Validate()
		if err == nil {
			t.Errorf("%s: a raw route declaring it should not validate", c.name)

			continue
		}

		if !strings.Contains(err.Error(), "raw body and an operation or resource") {
			t.Errorf("%s: unexpected error %v", c.name, err)
		}
	}
}

// A route cannot both send recorded bytes and be a declared failure. Allowing
// it would leave which one wins to the order of the checks in the handler.
func TestARawRouteRefusesAnError(t *testing.T) {
	r := rawOnlyRecipe(&RawBody{ContentType: "text/plain", Text: "OK"})
	r.Routes[0].Error = "boom"
	r.Errors = map[string]Error{"boom": {Status: 500, Message: "boom"}}

	err := r.Validate()
	if err == nil {
		t.Fatal("a raw route declaring an error should not validate")
	}

	if !strings.Contains(err.Error(), "raw body and an error") {
		t.Errorf("unexpected error: %v", err)
	}
}

// An empty body nobody meant serves 200 and nothing, which is exactly what a
// provider that answers with nothing looks like. The two have to be told
// apart, so saying nothing has to be said out loud.
func TestAnEmptyRawBodyMustBeDeclaredDeliberate(t *testing.T) {
	r := rawOnlyRecipe(&RawBody{ContentType: "text/plain"})

	err := r.Validate()
	if err == nil {
		t.Fatal("an empty raw body with no empty flag should not validate")
	}

	if !strings.Contains(err.Error(), "raw.empty") {
		t.Errorf("the error does not point at the flag that fixes it: %v", err)
	}

	r.Routes[0].Raw.Empty = true

	if err := r.Validate(); err != nil {
		t.Errorf("a deliberately empty raw body should validate: %v", err)
	}
}
