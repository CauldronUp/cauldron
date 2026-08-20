package runtime

import "testing"

// A count of records under a name that says pages is worse than an invented
// field: the name is real and the number is plausible. Documenso's envelope is
// {documents, totalPages}, and totalPages was the record count -- so three
// documents at ten per page reported three, and a client looping while
// page <= totalPages asked for two pages that do not exist and read them as
// empty results rather than as a mistake.

func TestPageCountRoundsUp(t *testing.T) {
	// A partial page is still a page. Rounding down loses the last one, which
	// is the records a client never fetches.
	for _, c := range []struct{ total, limit, want int }{
		{total: 3, limit: 10, want: 1},
		{total: 3, limit: 2, want: 2},
		{total: 3, limit: 1, want: 3},
		{total: 10, limit: 10, want: 1},
		{total: 11, limit: 10, want: 2},
		{total: 1, limit: 10, want: 1},
	} {
		if got := pageCount(c.total, c.limit); got != c.want {
			t.Errorf("pageCount(%d, %d) = %d, want %d", c.total, c.limit, got, c.want)
		}
	}
}

func TestAnEmptySetIsNoughtPages(t *testing.T) {
	// Providers differ about whether it is nought or one, and nought is the
	// reading that stops a loop rather than sending it after a page with
	// nothing in it.
	if got := pageCount(0, 10); got != 0 {
		t.Errorf("an empty set made %d pages", got)
	}
}

func TestPageCountSurvivesAnAbsurdLimit(t *testing.T) {
	// A page size of nought would divide by zero. It cannot arrive from a
	// Recipe -- the route default is applied first -- but a function that
	// panics on it is a function waiting for the first Recipe that can.
	if got := pageCount(3, 0); got != 0 {
		t.Errorf("a page size of nought made %d pages", got)
	}

	if got := pageCount(3, -1); got != 0 {
		t.Errorf("a negative page size made %d pages", got)
	}
}
