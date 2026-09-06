package recipe

import (
	"testing"
)

// A placeholder is not always a whole segment.
//
// SEC EDGAR serves a company at /submissions/CIK{id}.json and iTunes looks a
// record up at /lookup. The first of those has a parameter in the middle of a
// segment, and every rule in this file that asks "is this case about that
// route?" compared segments and treated one as a placeholder only when it
// began with { and ended with }. So a case asking for CIK0000320193.json was
// not recognised as being about the route that serves it, and four counters
// -- shown listings, sent paging parameters, described records, asserted field
// names -- all read the same Recipe as having no evidence it plainly has.
//
// It showed up as a sweep that would not converge: the generator kept writing
// the same cases, they kept passing, and the count never moved.
func TestAPlaceholderInsideASegmentMatches(t *testing.T) {
	for _, tt := range []struct {
		asked, declared string
		want            bool
	}{
		{"/submissions/CIK0000320193.json", "/submissions/CIK{id}.json", true},
		{"/submissions/CIK0000320193.txt", "/submissions/CIK{id}.json", false},
		{"/submissions/0000320193.json", "/submissions/CIK{id}.json", false},
		{"/api/formula/wget.json", "/api/formula/{id}.json", true},
		{"/v1/things/7", "/v1/things/{id}", true},
		{"/v1/things/7/parts/3", "/v1/things/{id}/parts/{part}", true},
		{"/v1/a-7-b", "/v1/a-{id}-b", true},
		{"/v1/a-7-c", "/v1/a-{id}-b", false},
		{"/v1/things", "/v1/things", true},
		{"/v1/things?limit=1", "/v1/things", true},
		{"/v1/other", "/v1/things", false},
	} {
		if got := samePath(tt.asked, tt.declared); got != tt.want {
			t.Errorf("samePath(%q, %q) = %v, want %v", tt.asked, tt.declared, got, tt.want)
		}
	}
}

// A placeholder matches within its own segment and not across a slash, or
// /v1/things/{id} would swallow /v1/things/7/parts.
func TestAPlaceholderDoesNotCrossASlash(t *testing.T) {
	if samePath("/v1/things/7/parts", "/v1/things/{id}") {
		t.Error("a placeholder matched across a slash")
	}

	if samePath("/files/a/b.json", "/files/{name}.json") {
		t.Error("a placeholder inside a segment matched across a slash")
	}
}
