package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

const tinySpec = `
openapi: 3.0.0
info: {title: Provider, version: "1"}
servers:
  - url: https://api.example.test
paths:
  /widgets:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id: {type: string}
                    colour: {type: string}
  /widgets/{id}:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  id: {type: string}
                  colour: {type: string}
`

func recipeNaming(spec string) *recipe.Recipe {
	return &recipe.Recipe{
		Name:       "provider",
		Version:    "0.1.0",
		Capability: "docs",
		Upstream:   recipe.Upstream{API: "v1", Spec: spec},
		Auth:       recipe.Auth{Scheme: "none"},
		Resources: map[string]recipe.Resource{
			"thing": {
				ID:         recipe.ID{Style: "opaque"},
				Collection: "things",
				Fields:     map[string]recipe.Field{"name": {Type: "string"}},
			},
		},
		Routes:   []recipe.Route{{Method: "GET", Path: "/things", Resource: "thing", Operation: "list"}},
		Fixtures: map[string]recipe.Fixture{"empty": {}},
	}
}

// The common case, and by a long way: a Recipe whose provider publishes nothing
// to read. That is not a failure and must not be reported as one -- the Recipe
// works, because it never needed the network.
func TestARecipeWithNoDescriptionComesBackUnchanged(t *testing.T) {
	r := recipeNaming("")

	merged, outcome := augmentFromSpec(r, func(string) ([]byte, error) {
		t.Error("a Recipe naming no description was fetched anyway")

		return nil, nil
	})

	if len(merged.Routes) != 1 {
		t.Errorf("routes became %d, want 1", len(merged.Routes))
	}

	if outcome.Why == "" {
		t.Error("nothing was said about why nothing was added")
	}
}

// A host that will not answer leaves the Recipe exactly as it was. The whole
// value of a Recipe is that it works offline; reaching for a description must
// never be able to take that away.
func TestAnUnreachableDescriptionLeavesTheRecipeWorking(t *testing.T) {
	r := recipeNaming("https://example.test/openapi.json")

	merged, outcome := augmentFromSpec(r, func(string) ([]byte, error) {
		return nil, errors.New("503 Service Unavailable")
	})

	if len(merged.Routes) != 1 {
		t.Errorf("routes became %d, want 1", len(merged.Routes))
	}

	if !strings.Contains(outcome.Why, "503") {
		t.Errorf("the reason does not name the failure: %q", outcome.Why)
	}
}

func TestADescriptionAddsRoutesTheRecipeDoesNotModel(t *testing.T) {
	r := recipeNaming("https://example.test/openapi.json")

	merged, outcome := augmentFromSpec(r, func(string) ([]byte, error) {
		return []byte(tinySpec), nil
	})

	if outcome.Report.Added == 0 {
		t.Fatalf("nothing was added: %+v", outcome)
	}

	var derived, written int

	for _, route := range merged.Routes {
		if route.Derived {
			derived++
		} else {
			written++
		}
	}

	if written != 1 {
		t.Errorf("%d written routes survived, want 1", written)
	}

	if derived == 0 {
		t.Error("no route was marked as derived")
	}

	if merged.DerivedFrom != "https://example.test/openapi.json" {
		t.Errorf("the merged Recipe does not name its description: %q", merged.DerivedFrom)
	}
}

// The report has one job beyond counting: making sure nobody reads a derived
// route as an observed one.
func TestTheReportSaysWhatADerivedRouteIsWorth(t *testing.T) {
	out := &bytes.Buffer{}

	reportSpecOutcomes(out, []specOutcome{
		{Recipe: "provider", Report: recipe.AugmentReport{Added: 3, Kept: 1}},
	})

	text := out.String()

	for _, want := range []string{
		"3 route(s) added",
		"what the provider says it does",
		"what it was seen doing",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("the report does not say %q:\n%s", want, text)
		}
	}
}

// And when nothing was added it must not imply anything was.
func TestTheReportIsQuietWhenNothingWasAdded(t *testing.T) {
	out := &bytes.Buffer{}

	reportSpecOutcomes(out, []specOutcome{
		{Recipe: "provider", Why: "publishes no description to read"},
	})

	text := out.String()

	if !strings.Contains(text, "No routes were added") {
		t.Errorf("the report does not say nothing happened:\n%s", text)
	}

	if strings.Contains(text, "what the provider says it does") {
		t.Errorf("the report warned about derived routes when there are none:\n%s", text)
	}

	if !strings.Contains(text, "publishes no description to read") {
		t.Errorf("the report does not say why:\n%s", text)
	}
}

// A description declares its paths relative to its own server, and a Recipe
// carries the whole path a client requests. Adyen is what found it: its
// description says /sessions beside a server of https://checkout-test.adyen.com
// /v71, and its Recipe says /v71/... -- so derived routes were mounted at
// /sessions, a path no client would ever call, and answered nothing at the
// address they would.
//
// It is the same mistake the fingerprint made before a test caught it there.
// Serving a route at the wrong path is worse than not serving it: the endpoint
// reports as added, and is missing.
func TestDerivedRoutesKeepTheBasePathFromTheDescription(t *testing.T) {
	based := strings.Replace(tinySpec,
		"  - url: https://api.example.test",
		"  - url: https://api.example.test/v71", 1)

	if based == tinySpec {
		t.Fatal("the test's own replacement did not apply")
	}

	merged, outcome := augmentFromSpec(recipeNaming("https://example.test/openapi.json"),
		func(string) ([]byte, error) { return []byte(based), nil })

	if outcome.Report.Added == 0 {
		t.Fatalf("nothing was added: %+v", outcome)
	}

	for _, route := range merged.Routes {
		if !route.Derived {
			continue
		}

		if !strings.HasPrefix(route.Path, "/v71/") {
			t.Errorf("derived route %s %s does not carry the description's base path", route.Method, route.Path)
		}
	}
}
