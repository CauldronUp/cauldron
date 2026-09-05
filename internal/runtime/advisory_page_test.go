package runtime

import (
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A page size is not always a contract, and two providers say so out loud.
//
// Missive: "A page may return more [items] than limit", documented on both of
// its listings. Onfleet's task listing "will return up to 64 tasks but may
// return fewer", and takes no size parameter at all, so a caller cannot even
// name a number to compare against.
//
// The format could already say a provider serves less than was asked for --
// max_limit trims, over_limit refuses -- and had no way to say either of
// these. Both break the same loop, which is the one everybody writes:
//
//	while len(page) == limit: fetch(next)
//
// Overshooting breaks it on the first page, because 26 != 25. Undershooting
// breaks it mid-collection, because 24 != 25. Neither errors, neither is
// visibly short, and both report a partial collection as the whole one.
func TestAPageCanCarryMoreRecordsThanWereAskedFor(t *testing.T) {
	r := advisory()
	r.Routes[0].Pagination.MayOvershoot = true

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things?limit=5")

	records, _ := body["things"].([]any)
	if len(records) != 6 {
		t.Errorf("a provider that may overshoot served %d records for limit=5, want 6", len(records))
	}
}

// And the walk still terminates, which is the part worth checking: a page that
// carries one extra record must advance the position by what it served rather
// than by what was asked for, or the walk repeats a record on every page.
func TestAnOvershootingWalkDoesNotRepeatARecord(t *testing.T) {
	r := advisory()
	r.Routes[0].Pagination.MayOvershoot = true
	r.Routes[0].Pagination.Style = "offset"
	r.Routes[0].Pagination.CursorParam = "offset"

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	first := listThings(t, s, "/v1/things?limit=5")
	second := listThings(t, s, "/v1/things?limit=5&offset=6")

	a, _ := first["things"].([]any)
	b, _ := second["things"].([]any)

	if len(a) != 6 || len(b) != 6 {
		t.Fatalf("pages were %d and %d records, want 6 and 6", len(a), len(b))
	}

	last := a[len(a)-1].(map[string]any)["id"]
	if got := b[0].(map[string]any)["id"]; got == last {
		t.Errorf("the second page began with %v, the record the first page ended on", got)
	}
}

func TestAPageCanCarryFewerRecordsWithoutBeingTheLast(t *testing.T) {
	r := advisory()
	r.Routes[0].Pagination.MayUndershoot = true

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things?limit=5")

	records, _ := body["things"].([]any)
	if len(records) != 4 {
		t.Errorf("a provider that may undershoot served %d records for limit=5, want 4", len(records))
	}

	// And it is not the last page, which is the whole point. A client that
	// carries on gets more records; a client that read the thin page as the
	// end stopped twenty-one records early and was told nothing.
	last := records[len(records)-1].(map[string]any)["id"].(string)

	next := listThings(t, s, "/v1/things?limit=5&cursor="+last)
	if rest, _ := next["things"].([]any); len(rest) == 0 {
		t.Error("a thin page was the last one after all")
	}
}

// A page of one may not become a page of none. Zero records ends a walk for a
// reason that is not true, which is the bug this key exists to stop rather
// than a smaller version of it.
func TestUndershootingNeverServesAnEmptyPage(t *testing.T) {
	r := advisory()
	r.Routes[0].Pagination.MayUndershoot = true

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	seed(t, s, 25)

	body := listThings(t, s, "/v1/things?limit=1")

	records, _ := body["things"].([]any)
	if len(records) != 1 {
		t.Errorf("limit=1 with may_undershoot served %d records, want 1", len(records))
	}
}

// And a Recipe cannot claim both. They are opposite statements about one page,
// and a file carrying both describes no provider at all.
func TestOvershootingAndUndershootingTogetherAreRefused(t *testing.T) {
	r := advisory()
	r.Routes[0].Pagination.MayOvershoot = true
	r.Routes[0].Pagination.MayUndershoot = true

	if err := r.Validate(); err == nil {
		t.Fatal("a route declaring both may_overshoot and may_undershoot validated")
	}
}

// Nor beside style: none. A listing that serves the whole collection has one
// page, it is always the last, and there is no size for it to miss.
func TestAnUnpagedListingCannotOvershoot(t *testing.T) {
	r := unpaged()
	r.Routes[0].Pagination.MayOvershoot = true

	if err := r.Validate(); err == nil {
		t.Fatal("a route declaring no paging and may_overshoot validated")
	}

	r = unpaged()
	r.Routes[0].Pagination.MayUndershoot = true

	if err := r.Validate(); err == nil {
		t.Fatal("a route declaring no paging and may_undershoot validated")
	}
}

// The same Recipe as the unpaged tests, with a page size instead of a refusal
// to have one.
func advisory() *recipe.Recipe {
	r := unpaged()
	r.Name = "advisory"
	r.Routes[0].Pagination = recipe.Pagination{
		Style:       "cursor",
		Limit:       10,
		LimitParam:  "limit",
		CursorParam: "cursor",
	}

	return r
}
