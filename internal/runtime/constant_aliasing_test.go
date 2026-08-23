package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

const aliasingRecipe = `
recipe: aliasing
capability: crm
version: 0.1.0
upstream:
  api: "1"
responses:
  list:
    style: wrapped
    limit_field: meta.limit
    page_field: meta.page
    fields:
      meta:
        source: cauldron
resources:
  thing:
    collection: things
    id:
      prefix: th_
      length: 8
    fields:
      name:
        type: string
routes:
  - method: GET
    path: /things
    resource: thing
    operation: list
    pagination:
      style: page
      limit: 10
      limit_param: limit
      cursor_param: page
  - method: GET
    path: /others
    resource: thing
    operation: list
    pagination:
      style: page
      limit: 10
      limit_param: "-"
      cursor_param: "-"
fixtures:
  small:
    thing:
      - id: th_aaaaaaaa
        name: one
      - id: th_bbbbbbbb
        name: two
conformance: []
`

// A declared constant went into the response by reference, and the echo
// fields are written into that same path afterwards. So serving one request
// rewrote the Recipe:
//
//	before: {"source": "cauldron"}
//	after:  {"source": "cauldron", "page": 3, "limit": 1}
//
// Two things follow, and the second is the worse one. A later request on a
// different route carried the earlier request's numbers -- including a route
// declaring "-" for both, which means it sends neither. And Reset does not
// undo it: it rewinds the store, the clock, the faults, the log and the
// webhooks, and never touches the Recipe. One request poisoned a long-lived
// serve for the rest of its life.
func TestAConstantIsNotAliasedIntoTheResponse(t *testing.T) {
	r, err := recipe.Parse([]byte(aliasingRecipe))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	get := func(path string) map[string]any {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s = %d\n%s", path, rec.Code, rec.Body)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}

		return body
	}

	before := fmt.Sprint(r.Responses.List.Fields["meta"])

	if meta, _ := get("/things?page=3&limit=1")["meta"].(map[string]any); fmt.Sprint(meta["page"]) != "3" {
		t.Errorf("the echo did not reach the response: %v", meta)
	}

	if after := fmt.Sprint(r.Responses.List.Fields["meta"]); after != before {
		t.Fatalf("serving a request rewrote the Recipe's own constant:\n before %s\n after  %s", before, after)
	}

	// A second route, which accepts no paging parameters at all. It still
	// serves a page and still reports one -- that part is right. What it must
	// not report is the numbers from the request before it, on a different
	// route, which is what the aliasing produced.
	meta, _ := get("/others")["meta"].(map[string]any)

	if got := fmt.Sprint(meta["page"]); got != "1" {
		t.Errorf("page = %s, want 1: this route was asked for no particular page", got)
	}

	if got := fmt.Sprint(meta["limit"]); got != "10" {
		t.Errorf("limit = %s, want its own declared 10, not the 1 the previous request asked for", got)
	}

	if meta["source"] != "cauldron" {
		t.Errorf("the constant itself was lost: %v", meta)
	}
}
