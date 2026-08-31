package openapi

import (
	"errors"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

func discoverRecipe(routes ...string) *recipe.Recipe {
	r := &recipe.Recipe{Name: "example"}

	for _, path := range routes {
		r.Routes = append(r.Routes, recipe.Route{Method: "GET", Path: path, Resource: "thing"})
	}

	return r
}

// serveURLs answers a fixed map and errors for anything else, which is how a
// host that has no description at a given path behaves.
func serveURLs(bodies map[string]string) func(string) ([]byte, error) {
	return func(url string) ([]byte, error) {
		if body, ok := bodies[url]; ok {
			return []byte(body), nil
		}

		return nil, errors.New("404")
	}
}

const discoverSpec = `
openapi: 3.0.0
info: {title: Example, version: "1"}
servers: [{url: "https://api.example.com"}]
paths:
  /widgets:
    get:
      responses:
        "200": {description: ok}
  /widgets/{id}:
    get:
      responses:
        "200": {description: ok}
  /gadgets:
    get:
      responses:
        "200": {description: ok}
`

// The whole point of the command. A description that declares what the Recipe
// models is that Recipe's description.
func TestADescriptionThatDeclaresTheRecipesRoutesIsProposed(t *testing.T) {
	r := discoverRecipe("/widgets", "/widgets/{id}")

	found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": discoverSpec,
	}))

	if found == nil {
		t.Fatal("nothing was proposed")
	}

	if found.URL != "https://api.example.com/openapi.json" {
		t.Errorf("URL = %q", found.URL)
	}

	if found.Matched != 2 || found.Routes != 2 {
		t.Errorf("matched %d of %d, want 2 of 2", found.Matched, found.Routes)
	}

	if found.Declared != 3 {
		t.Errorf("declared = %d, want 3", found.Declared)
	}
}

// This is the test the command exists to pass. Documentation platforms serve a
// generic openapi.json at the vendor's marketing domain: ramp.com and
// circleci.com both answer one, with two and five paths, describing the docs
// site rather than the API. A discovery that counted a 200 as a find would
// record those as the provider's own description and then report drift against
// them forever.
func TestADescriptionThatDeclaresNoneOfTheRoutesIsNotProposed(t *testing.T) {
	r := discoverRecipe("/v1/payments", "/v1/payments/{id}")

	stub := `
openapi: 3.0.0
info: {title: Docs, version: "1"}
paths:
  /search:
    get:
      responses:
        "200": {description: ok}
`

	if found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": stub,
	})); found != nil {
		t.Errorf("a description sharing nothing was proposed: %+v", found)
	}
}

// The mistake already made twice in this package: a description declares its
// paths relative to its own server, so /widgets under a server of /v2 is the
// Recipe's /v2/widgets. Comparing literally finds nothing and reports a real
// description as no description at all.
func TestTheBasePathIsHonouredWhenMatching(t *testing.T) {
	r := discoverRecipe("/v2/widgets")

	spec := `
openapi: 3.0.0
info: {title: Example, version: "1"}
servers: [{url: "https://api.example.com/v2"}]
paths:
  /widgets:
    get:
      responses:
        "200": {description: ok}
`

	found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": spec,
	}))

	if found == nil {
		t.Fatal("a description whose server carries the prefix was not matched")
	}

	if found.Matched != 1 {
		t.Errorf("matched = %d, want 1", found.Matched)
	}
}

// Swagger 2.0 is read now, so discovery has to find it too -- most of the
// providers this was built to reach publish 2.0 and nothing else.
func TestASwagger2DescriptionIsDiscovered(t *testing.T) {
	r := discoverRecipe("/v1/things")

	spec := `
swagger: "2.0"
info: {title: Old, version: "1"}
host: api.example.com
basePath: /v1
paths:
  /things:
    get:
      responses:
        "200": {description: ok}
`

	found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/swagger.json": spec,
	}))

	if found == nil {
		t.Fatal("a Swagger 2.0 description was not discovered")
	}

	if found.Format != "Swagger 2.0" {
		t.Errorf("format = %q, want Swagger 2.0", found.Format)
	}
}

// Where several URLs answer, the one that describes more of the Recipe is the
// one that is actually this API. A vendor serving a docs stub at one path and
// the real description at another is the common case.
func TestTheStrongestCandidateWins(t *testing.T) {
	r := discoverRecipe("/widgets", "/widgets/{id}")

	weak := `
openapi: 3.0.0
info: {title: Stub, version: "1"}
paths:
  /widgets:
    get:
      responses:
        "200": {description: ok}
`

	found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": weak,
		"https://api.example.com/swagger.json": discoverSpec,
	}))

	if found == nil {
		t.Fatal("nothing was proposed")
	}

	if found.Matched != 2 {
		t.Errorf("matched = %d, want the stronger candidate's 2", found.Matched)
	}
}

// A host that fails must not stop the others being tried, or one dead docs
// domain hides a description that is served perfectly well elsewhere.
func TestAFailingHostDoesNotStopTheRest(t *testing.T) {
	r := discoverRecipe("/widgets")

	found := Discover(r, []string{"dead.example.com", "api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": discoverSpec,
	}))

	if found == nil {
		t.Fatal("a live host was not tried after a dead one")
	}
}

// An HTML page that happens to parse is the failure mode this whole area keeps
// hitting, and a description with no paths is a login redirect.
func TestAnHTMLPageIsNotADescription(t *testing.T) {
	r := discoverRecipe("/widgets")

	if found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": "<!DOCTYPE html><html><body>docs</body></html>",
	})); found != nil {
		t.Errorf("an HTML page was proposed: %+v", found)
	}
}

// A Recipe with no routes has nothing to prove a match against, and proposing
// on no evidence is the guessing this is built to avoid.
func TestARecipeWithNoRoutesProposesNothing(t *testing.T) {
	if found := Discover(&recipe.Recipe{Name: "empty"}, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://api.example.com/openapi.json": discoverSpec,
	})); found != nil {
		t.Errorf("a Recipe with no routes was given a description: %+v", found)
	}
}

// Attio is verified against api.attio.com and publishes its description at
// attio.com; Monday is verified against api.monday.com and publishes at
// monday.com. The verified host names the organisation, and the description
// frequently lives on a sibling of it.
//
// Widening the search this way would be reckless on its own -- it is one step
// from the domain matching that paired Bugsnag with ClickSend. It is safe here
// only because the proof does the work: whatever is served at the sibling still
// has to declare a path this Recipe models.
func TestSiblingHostsOfTheVerifiedHostAreTried(t *testing.T) {
	r := discoverRecipe("/widgets")

	for name, url := range map[string]string{
		"the bare domain":  "https://example.com/openapi.json",
		"a docs host":      "https://docs.example.com/openapi.json",
		"a developer host": "https://developers.example.com/openapi.json",
	} {
		found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{url: discoverSpec}))
		if found == nil {
			t.Errorf("%s: nothing found at %s", name, url)

			continue
		}

		if found.URL != url {
			t.Errorf("%s: URL = %q, want %q", name, found.URL, url)
		}
	}
}

// A sibling that serves something unrelated is still refused, which is the
// guarantee that makes widening the search safe at all.
func TestASiblingServingSomethingElseIsStillRefused(t *testing.T) {
	r := discoverRecipe("/v1/payments")

	unrelated := "openapi: 3.0.0\ninfo: {title: Marketing, version: \"1\"}\npaths:\n  /blog:\n    get:\n      responses:\n        \"200\": {description: ok}\n"

	if found := Discover(r, []string{"api.example.com"}, serveURLs(map[string]string{
		"https://example.com/openapi.json": unrelated,
	})); found != nil {
		t.Errorf("a sibling serving an unrelated description was proposed: %+v", found)
	}
}

// A two-label public suffix must not be mistaken for the registered domain, or
// the search widens to every API hosted in the United Kingdom.
func TestAPublicSuffixIsNotTreatedAsTheDomain(t *testing.T) {
	for _, c := range []struct{ host, want string }{
		{"api.example.com", "example.com"},
		{"api.example.co.uk", "example.co.uk"},
		{"www.ebi.ac.uk", "ebi.ac.uk"},
		{"example.com", "example.com"},
	} {
		if got := registeredDomain(c.host); got != c.want {
			t.Errorf("registeredDomain(%q) = %q, want %q", c.host, got, c.want)
		}
	}
}
