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

func TestNamingBothParametersMarksItChecked(t *testing.T) {
	checked := Pagination{Limit: 20, Style: "cursor", LimitParam: "per_page", CursorParam: "page_token"}

	if got := paginated(checked).GuessedPagination(); got != 0 {
		t.Errorf("a fully named block counted %d, want 0", got)
	}
}

// Naming one of the two is not naming them.
//
// The page size and the position are separate guesses, and a Recipe that has
// looked up one has not necessarily looked up the other. The runtime supplies
// whichever is missing -- "limit" for the size, the style's own word for the
// position.
func TestNamingOneParameterLeavesTheOtherAGuess(t *testing.T) {
	if got := paginated(Pagination{Limit: 20, LimitParam: "per_page"}).GuessedPagination(); got != 1 {
		t.Errorf("a named size with an unnamed position counted %d, want 1", got)
	}

	if got := paginated(Pagination{Limit: 20, CursorParam: "page_token"}).GuessedPagination(); got != 1 {
		t.Errorf("a named position with an unnamed size counted %d, want 1", got)
	}
}

// Naming the style is not naming the parameters either, and this counter said
// it was.
//
// Front is the counterexample that found it. Three of its listings declare
// cursor paging and name neither parameter, and Front's own description says
// they take "limit" and "page_token" -- so the runtime read "cursor", Front
// ignored it, and every request came back on page one. Ninety-four routes
// across forty Recipes were in that state, and neither this counter nor
// UnstatedPagination could see one of them: this required an empty style, and
// that one required an empty everything.
func TestDeclaringTheStyleIsNotNamingTheParameters(t *testing.T) {
	if got := paginated(Pagination{Limit: 20, Style: "cursor"}).GuessedPagination(); got != 1 {
		t.Errorf("a declared style with no parameter names counted %d, want 1", got)
	}
}

// A provider that accepts no name for a parameter is a decision, not a gap.
// "-" is how a Recipe says so, and Webflow's site-level listings say it.
func TestARefusedNameIsNotAGuess(t *testing.T) {
	refused := Pagination{Limit: 100, Style: "cursor", LimitParam: "-", CursorParam: "-"}

	if got := paginated(refused).GuessedPagination(); got != 0 {
		t.Errorf("a block naming both refusals counted %d, want 0", got)
	}
}

// And a listing that says it does not page has nothing to guess at.
func TestAListingThatDoesNotPageIsNotAGuess(t *testing.T) {
	if got := paginated(Pagination{Style: "none"}).GuessedPagination(); got != 0 {
		t.Errorf("a listing declaring no paging counted %d, want 0", got)
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
		// Half named is still a guess about the other half.
		{Path: "/c", Pagination: Pagination{Limit: 10, LimitParam: "count"}},
		{Path: "/d"},
		{Path: "/e", Pagination: Pagination{Limit: 10, LimitParam: "count", CursorParam: "after"}},
	}}

	if got := r.GuessedPagination(); got != 3 {
		t.Errorf("counted %d guessing routes, want 3", got)
	}
}
