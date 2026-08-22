package recipe

import "testing"

// listing builds a Recipe with one list route carrying the given paging.
func listing(p Pagination) Recipe {
	return Recipe{Routes: []Route{{Operation: "list", Pagination: p}}}
}

// A listing that says nothing about paging is still paged. The runtime gives
// it ten and reads "limit", which is the same guess GuessedPagination exists
// to report -- and GuessedPagination cannot see it, because it starts from a
// declared page size.
//
// The verify report said "60 routes page by a parameter nobody named" while
// 108 more paged the same way, so the figure read as the whole answer and was
// a third of it.
func TestAListingDeclaringNothingIsCountedAsUnstated(t *testing.T) {
	if got := listing(Pagination{}).UnstatedPagination(); got != 1 {
		t.Errorf("counted %d, want 1: a listing with no paging block is still paged at ten", got)
	}
}

func TestADeclaredPageSizeIsNotUnstated(t *testing.T) {
	// This one is GuessedPagination's to report. Counting it twice would make
	// the two figures overlap and neither of them mean anything.
	if got := listing(Pagination{Limit: 20}).UnstatedPagination(); got != 0 {
		t.Errorf("counted %d, want 0: a declared page size is the other count's", got)
	}
}

func TestANamedParameterAloneClearsIt(t *testing.T) {
	if got := listing(Pagination{LimitParam: "per_page"}).UnstatedPagination(); got != 0 {
		t.Errorf("counted %d, want 0: the parameter is named", got)
	}

	if got := listing(Pagination{Style: "cursor"}).UnstatedPagination(); got != 0 {
		t.Errorf("counted %d, want 0: the style is named", got)
	}
}

// Only listings. A get or a create has no page to size, and counting them
// would inflate the figure with routes that could never carry paging.
func TestOnlyListingsAreCounted(t *testing.T) {
	r := Recipe{Routes: []Route{
		{Operation: "get"},
		{Operation: "create"},
		{Operation: "delete"},
	}}

	if got := r.UnstatedPagination(); got != 0 {
		t.Errorf("counted %d, want 0: only listings page", got)
	}
}

// The two counts must not overlap, across every Recipe that ships. If they
// did, the report would be adding a route to both totals and describing the
// same omission twice.
func TestTheTwoPagingCountsDoNotOverlap(t *testing.T) {
	for _, name := range Bundled() {
		r, err := Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		lists := 0
		for _, route := range r.Routes {
			if route.Operation == "list" {
				lists++
			}
		}

		if guessed, unstated := r.GuessedPagination(), r.UnstatedPagination(); guessed+unstated > lists {
			t.Errorf("%s: %d guessed plus %d unstated is more than its %d listings, so a route is in both", name, guessed, unstated, lists)
		}
	}
}
