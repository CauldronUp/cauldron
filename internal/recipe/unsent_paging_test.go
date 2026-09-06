package recipe

import (
	"testing"
)

// A paging parameter name is a claim about the provider, and a claim nothing
// sends is not checked by anything.
//
// The format already refuses a *response* field name that no case asserts:
// declare cursor_field: next_page_token, rename it to anything, and no case
// notices, so validation stops it. The request half had no equivalent and is
// the half that matters more. A client sends the size and the position; getting
// either name wrong is the failure the whole paging vocabulary exists to
// prevent, and asserting the response to a listing says nothing about which
// parameter produced it.
//
// This is the mutation the count is named after: rename limit_param and see
// whether anything fails.
func TestAPagingParameterNoCaseSendsIsCounted(t *testing.T) {
	r := &Recipe{
		Name:       "counted",
		Capability: "ai",
		Version:    "0.1.0",
		Upstream:   Upstream{API: "v1"},
		Resources:  map[string]Resource{"thing": {Collection: "things"}},
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list",
			Pagination: Pagination{
				Style:       "offset",
				Limit:       10,
				LimitParam:  "per_page",
				CursorParam: "offset",
			},
		}},
	}

	if n := r.UnsentPagingParam(); n != 2 {
		t.Fatalf("two declared names, no cases: counted %d, want 2", n)
	}

	// A case that sends one of them settles that one and not the other.
	r.Conformance = []Case{{
		Request: Request{
			Method: "GET",
			Path:   "/v1/things",
			Query:  map[string]string{"per_page": "2"},
		},
	}}

	if n := r.UnsentPagingParam(); n != 1 {
		t.Errorf("one name sent: counted %d, want 1", n)
	}

	// The second one too, and the debt is gone.
	r.Conformance[0].Request.Query["offset"] = "2"

	if n := r.UnsentPagingParam(); n != 0 {
		t.Errorf("both names sent: counted %d, want 0", n)
	}
}

// The query may travel in the path, which is the older spelling and still
// common across this collection.
func TestAParameterInThePathCounts(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Operation: "list",
			Pagination: Pagination{Style: "offset", Limit: 10, LimitParam: "per_page"},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things?per_page=2&offset=2"},
		}},
	}

	if n := r.UnsentPagingParam(); n != 0 {
		t.Errorf("a parameter sent in the path counted as unsent: %d", n)
	}
}

// And a route whose paging travels in the body is checked against the body.
// Reading the query string there would count every one of them as unsent,
// which is the wrong answer for eleven Recipes in this collection.
func TestABodyParameterIsCheckedAgainstTheBody(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "POST", Path: "/v1/search", Operation: "list",
			Pagination: Pagination{
				Style: "cursor", Limit: 10,
				LimitParam: "options.count", In: "body",
			},
		}},
		Conformance: []Case{{
			Request: Request{
				Method: "POST",
				Path:   "/v1/search",
				JSON:   map[string]any{"options": map[string]any{"count": 2}},
			},
		}},
	}

	if n := r.UnsentPagingParam(); n != 0 {
		t.Errorf("a dotted body parameter counted as unsent: %d", n)
	}

	// And the same name absent from the body is counted.
	r.Conformance[0].Request.JSON = map[string]any{"options": map[string]any{"offset": 2}}

	if n := r.UnsentPagingParam(); n != 1 {
		t.Errorf("a body parameter nobody sends: counted %d, want 1", n)
	}
}

// "-" is not a name. It says the provider accepts none, which is a claim about
// absence -- a case cannot send a parameter that does not exist, so counting it
// would make an honest declaration look like debt.
func TestTheAbsentParameterIsNotCounted(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Operation: "list",
			Pagination: Pagination{
				Style: "cursor", Limit: 10,
				LimitParam: "-", CursorParam: "-",
			},
		}},
	}

	if n := r.UnsentPagingParam(); n != 0 {
		t.Errorf(`"-" counted as an unsent name: %d`, n)
	}
}

// A path parameter in the route matches whatever the case put in its position.
// Without this, every scoped listing in the collection counts as unsent because
// the case sends a real id where the route declares {app_id}.
func TestAPlaceholderMatchesTheValueACaseSends(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/apps/{app_id}/things", Operation: "list",
			Pagination: Pagination{Style: "offset", Limit: 10, LimitParam: "per_page"},
		}},
		Conformance: []Case{{
			Request: Request{
				Method: "GET",
				Path:   "/v1/apps/app_42/things",
				Query:  map[string]string{"per_page": "2"},
			},
		}},
	}

	if n := r.UnsentPagingParam(); n != 0 {
		t.Errorf("a case against a scoped listing counted as unsent: %d", n)
	}
}
