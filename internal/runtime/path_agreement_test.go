package runtime

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The two path matchers have to agree.
//
// The router decides which route serves a request. recipe.PathMatches decides
// which route a conformance case is about, and four counters are built on it:
// shown listings, sent paging parameters, described records, asserted field
// names. They are separate implementations of one idea, and they had drifted
// three ways at once -- a placeholder embedded in a segment, a greedy
// {name...}, and the whole-segment case both got right.
//
// The drift was invisible because each is green on its own. It surfaced as a
// generator sweep that would not converge: cases written, passing, and counted
// by nothing. This compares them directly, on every path every bundled Recipe's
// cases actually ask for.
//
// The router is the authority. A path it routes to a route is a path about that
// route, whatever the second implementation thinks.
func TestPathMatchesAgreesWithTheRouter(t *testing.T) {
	checked := 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		router := newRouter(r)

		for _, c := range r.Conformance {
			asked := c.Request.Path
			if i := strings.Index(asked, "?"); i >= 0 {
				asked = asked[:i]
			}

			if asked == "" {
				continue
			}

			parts := splitPath(asked)

			for _, rt := range router.routes {
				if rt.method != c.Request.Method {
					continue
				}

				if _, _, ok := rt.matches(parts); !ok {
					continue
				}

				checked++

				if !recipe.PathMatches(asked, rt.spec.Path) {
					t.Errorf("%s: the router routes %q to %q and PathMatches says it is a different route",
						name, asked, rt.spec.Path)
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("no path pairs were compared, so this test proved nothing")
	}

	t.Logf("compared %d routed paths", checked)
}

// And the other direction, with one stated exception.
//
// PathMatches trims a trailing slash from both sides and the router does not,
// because for Django-shaped APIs the slash is part of the path -- sixty-one
// routes across five Recipes declare one. So PathMatches is deliberately the
// more forgiving of the two there, and nowhere else.
func TestPathMatchesIsNoMoreForgivingThanTheRouter(t *testing.T) {
	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		router := newRouter(r)

		for _, c := range r.Conformance {
			asked := c.Request.Path
			if i := strings.Index(asked, "?"); i >= 0 {
				asked = asked[:i]
			}

			if asked == "" {
				continue
			}

			parts := splitPath(asked)

			for _, rt := range router.routes {
				if rt.method != c.Request.Method || !recipe.PathMatches(asked, rt.spec.Path) {
					continue
				}

				if _, _, ok := rt.matches(parts); ok {
					continue
				}

				if rt.slash || strings.HasSuffix(asked, "/") {
					continue
				}

				t.Errorf("%s: PathMatches says %q is about %q and the router does not route it there",
					name, asked, rt.spec.Path)
			}
		}
	}
}
