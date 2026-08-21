package recipe

import "testing"

// A Recipe can be entirely green and still be guessing at how its provider
// pages. A conformance case that never sends a page size cannot notice the
// wrong name being read: the parameter is ignored, one full page comes back,
// and the case asserts a listing that happens to fit on it.
//
// The count exists so that the number is visible beside the evidence rather
// than buried, and so that it falls as Recipes are checked.

func paginated(p Pagination) Recipe {
	return Recipe{Routes: []Route{{Method: "GET", Path: "/things", Pagination: p}}}
}

func TestADeclaredSizeWithNoNamesIsAGuess(t *testing.T) {
	// The runtime reads "limit", which is right for some providers and wrong
	// for plenty, and nothing here recorded which.
	if got := paginated(Pagination{Limit: 20}).GuessedPagination(); got != 1 {
		t.Errorf("a bare limit counted %d, want 1", got)
	}
}

func TestNamingTheParameterMarksItChecked(t *testing.T) {
	if got := paginated(Pagination{Limit: 20, LimitParam: "per_page"}).GuessedPagination(); got != 0 {
		t.Errorf("a named parameter counted %d, want 0", got)
	}
}

func TestNamingTheStyleMarksItChecked(t *testing.T) {
	// Either one is enough, because neither is a word anybody writes down by
	// accident.
	if got := paginated(Pagination{Limit: 20, Style: "cursor"}).GuessedPagination(); got != 0 {
		t.Errorf("a declared style counted %d, want 0", got)
	}
}

func TestARouteThatDoesNotPageIsNotAGuess(t *testing.T) {
	// No page size declared is a listing that answers with everything, which
	// is a great many of them and not a claim about paging at all.
	if got := paginated(Pagination{}).GuessedPagination(); got != 0 {
		t.Errorf("an unpaged route counted %d, want 0", got)
	}
}

func TestEveryGuessingRouteCounts(t *testing.T) {
	r := Recipe{Routes: []Route{
		{Path: "/a", Pagination: Pagination{Limit: 10}},
		{Path: "/b", Pagination: Pagination{Limit: 10}},
		{Path: "/c", Pagination: Pagination{Limit: 10, LimitParam: "count"}},
		{Path: "/d"},
	}}

	if got := r.GuessedPagination(); got != 2 {
		t.Errorf("counted %d guessing routes, want 2", got)
	}
}
