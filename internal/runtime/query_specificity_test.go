// Two routes matching one request, and which of them answers.

package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A route selected by query parameters used to score a flat bonus however many
// it matched, so two matching routes tied and the request went to whichever was
// declared first. Nothing in either route said it depended on the order, and
// nothing failed when the order was wrong -- the response was simply the other
// route's, with fields missing.
//
// Datamuse is what found it. /words?rel_rhy=blue answers a syllable count and
// no tags; adding md=dpsr answers both plus definitions. Modelled as two
// routes, the second request matched them both, and the Recipe that declared
// five fields served three.
//
// This uses a shipped Recipe on purpose: the overlap is the provider's.
func TestTheRouteMatchingMoreQueryParametersAnswers(t *testing.T) {
	r, err := recipe.Open("datamuse")
	if err != nil {
		t.Fatalf("open datamuse: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("slew"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	first := func(path string) map[string]any {
		t.Helper()

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d: %s", path, rec.Code, rec.Body.String())
		}

		var list []map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
			t.Fatalf("%s: decode: %v", path, err)
		}

		if len(list) == 0 {
			t.Fatalf("%s: no results", path)
		}

		return list[0]
	}

	// One parameter matched: the rhyme route, which carries a syllable count
	// nobody asked for and no tags at all.
	rhyme := first("/words?rel_rhy=blue&max=2")

	if _, ok := rhyme["numSyllables"]; !ok {
		t.Errorf("the rhyme route did not answer: %#v", rhyme)
	}

	if _, ok := rhyme["tags"]; ok {
		t.Errorf("the rhyme route answered tags, which that query does not carry: %#v", rhyme)
	}

	// Two parameters matched: the metadata route, which carries both, and it
	// has to win against the route above even though that one also matches.
	full := first("/words?rel_rhy=blue&md=dpsr&max=2")

	for _, field := range []string{"numSyllables", "tags", "defs"} {
		if _, ok := full[field]; !ok {
			t.Errorf("the metadata route did not answer -- %q is missing, so the less specific route won: %#v", field, full)
		}
	}
}
