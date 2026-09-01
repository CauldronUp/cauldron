package runtime

import (
	"net/http"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A parameter declared {name...} swallows the rest of the path.
//
// Three Recipes needed this and none of them could have it. A DOI is
// "10.1145/3510003" and Semantic Scholar takes it raw, so the identifier a
// caller holds cannot reach the route that addresses it. A Guardian content id
// is "world/2023/jan/01/some-headline" -- one identifier written with four
// slashes in it. A Hugging Face repo is "org/name". Percent-encoding is not an
// escape hatch: Semantic Scholar answers a %2F with an error naming the
// encoding.
//
// The spelling is Go's own, and so is the restriction to the final segment.
// Anywhere else there is nothing to say where the parameter stops.
func TestAGreedyParameterSwallowsTheRestOfThePath(t *testing.T) {
	r := &recipe.Recipe{
		Name:       "greedy",
		Capability: "search",
		Version:    "0.1.0",
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/paper/{id...}", Resource: "paper", Operation: "get"},
			{Method: "GET", Path: "/v1/paper/search", Resource: "paper", Operation: "list"},
		},
	}

	rt := newRouter(r)

	cases := []struct {
		path string
		want string
		op   string
	}{
		{"/v1/paper/649def34", "649def34", "get"},
		{"/v1/paper/DOI:10.1145/3510003", "DOI:10.1145/3510003", "get"},
		{"/v1/paper/world/2023/jan/01/some-headline", "world/2023/jan/01/some-headline", "get"},
		// The literal route has to win the request it spells out, or a greedy
		// sibling quietly takes over every path in the collection.
		{"/v1/paper/search", "", "list"},
	}

	for _, c := range cases {
		matched, vars, ok := rt.match(http.MethodGet, c.path)
		if !ok {
			t.Errorf("%s matched nothing", c.path)

			continue
		}

		if matched.spec.Operation != c.op {
			t.Errorf("%s reached the %s route, want %s", c.path, matched.spec.Operation, c.op)
		}

		if got := vars["id"]; got != c.want {
			t.Errorf("%s captured id %q, want %q", c.path, got, c.want)
		}
	}

	// Greedy means at least one segment, not zero: the route still needs
	// something after the literals to swallow.
	if _, _, ok := rt.match(http.MethodGet, "/v1/paper"); ok {
		t.Error("/v1/paper matched a route whose last segment has nothing to take")
	}
}

// A greedy parameter anywhere but the end is refused by the validator.
//
// It would compile, the router would accept it, and every request written for
// it would miss -- the parameter eating the segments that were meant to follow.
// Silence is the wrong answer to that, so it is a validation failure.
func TestAGreedyParameterIsRefusedBeforeTheEndOfAPath(t *testing.T) {
	r := &recipe.Recipe{
		Name:       "greedy",
		Capability: "search",
		Version:    "0.1.0",
		Routes: []recipe.Route{
			{Method: "GET", Path: "/v1/{repo...}/revision/{rev}", Resource: "repo", Operation: "get"},
		},
	}

	err := r.Validate()
	if err == nil {
		t.Fatal("a greedy parameter before the end of a path validated")
	}

	if got := err.Error(); !strings.Contains(got, "greedy parameter can only be the last segment") {
		t.Errorf("the failure does not say why: %s", got)
	}
}
