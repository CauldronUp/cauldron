package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// rawRecipe is the shape a provider that never sends JSON needs: a route that
// answers with the bytes it was recorded sending, and no resource behind it.
func rawRecipe(t *testing.T, raw *recipe.RawBody) *recipe.Recipe {
	t.Helper()

	return &recipe.Recipe{
		Name:       "example",
		Capability: "docs",
		Version:    "0.1.0",
		Auth:       recipe.Auth{Scheme: "none"},
		Upstream:   recipe.Upstream{API: "v1"},
		Routes: []recipe.Route{
			{Method: "GET", Path: "/api/query", Raw: raw},
		},
		Fixtures: map[string]recipe.Fixture{"empty": {}},
	}
}

const atomFeed = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>arXiv Query</title>
</feed>`

// The whole point. arXiv, Canada Post, PrestaShop and the older USPS
// generation answer XML and nothing else, and a mock that answers JSON is
// wrong about every one of them in the way that matters most: a client parses
// the response, and the parse is the integration.
func TestARawRouteServesTheBytesItWasGiven(t *testing.T) {
	r := rawRecipe(t, &recipe.RawBody{
		ContentType: "application/atom+xml; charset=utf-8",
		Text:        atomFeed,
	})

	rec := serveRaw(t, r, http.MethodGet, "/api/query")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	if got, want := rec.Body.String(), atomFeed; got != want {
		t.Errorf("body =\n%s\nwant\n%s", got, want)
	}

	if got, want := rec.Header().Get("Content-Type"), "application/atom+xml; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

// A declared content type is the difference between a client parsing the
// response and a client throwing on the first character, so it is not
// something to default quietly. But a Recipe that says nothing should still
// not be served as JSON.
func TestARawRouteWithNoContentTypeIsNotServedAsJSON(t *testing.T) {
	rec := serveRaw(t, rawRecipe(t, &recipe.RawBody{Text: "OK (not found)"}), http.MethodGet, "/api/query")

	if got := rec.Header().Get("Content-Type"); strings.Contains(got, "json") {
		t.Errorf("Content-Type = %q, and a raw body is not JSON", got)
	}

	if got, want := rec.Body.String(), "OK (not found)"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

// Healthchecks' ping answers 200 with a body a client reads as a string, and
// the management API on the same Recipe answers 401. A raw route has to be
// able to carry a status other than 200 or half of that is unsayable.
func TestARawRouteHonoursItsDeclaredStatus(t *testing.T) {
	r := rawRecipe(t, &recipe.RawBody{Text: "invalid url format", ContentType: "text/plain"})
	r.Routes[0].Status = 400

	if rec := serveRaw(t, r, http.MethodGet, "/api/query"); rec.Code != 400 {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

// A route declares headers and must send them whatever it answers with. This
// is the same rule the error routes already follow, and it exists because a
// declaration that is silently inert is worse than no declaration.
func TestARawRouteSendsItsDeclaredHeaders(t *testing.T) {
	r := rawRecipe(t, &recipe.RawBody{Text: "OK", ContentType: "text/plain"})
	r.Routes[0].Headers = map[string]string{"Ping-Body-Limit": "100000"}

	rec := serveRaw(t, r, http.MethodGet, "/api/query")

	if got, want := rec.Header().Get("Ping-Body-Limit"), "100000"; got != want {
		t.Errorf("Ping-Body-Limit = %q, want %q", got, want)
	}
}

// An empty body is a real answer -- a 204, a ping that acknowledges and says
// nothing -- and has to be distinguishable from a route that forgot to
// declare one. Validation refuses the second; this covers the first.
func TestARawRouteMaySendNothingAtAll(t *testing.T) {
	r := rawRecipe(t, &recipe.RawBody{Text: "", ContentType: "text/plain"})
	r.Routes[0].Status = 204
	r.Routes[0].Raw.Empty = true

	rec := serveRaw(t, r, http.MethodGet, "/api/query")

	if rec.Code != 204 {
		t.Errorf("status = %d, want 204", rec.Code)
	}

	if got := rec.Body.String(); got != "" {
		t.Errorf("body = %q, want empty", got)
	}
}

func serveRaw(t *testing.T, r *recipe.Recipe, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	if err := r.Validate(); err != nil {
		t.Fatalf("recipe does not validate: %v", err)
	}

	sandbox, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	rec := httptest.NewRecorder()
	sandbox.ServeHTTP(rec, httptest.NewRequest(method, path, nil))

	return rec
}
