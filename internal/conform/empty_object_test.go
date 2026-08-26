package conform

import "testing"

// An empty object in an expectation means an empty object.
//
// An object expectation is deliberately partial: {a: 1} says the response
// carries a and says nothing about its neighbours, which is what lets a case
// name the fields it cares about without freezing the rest of the body.
//
// That degenerates at zero. {} named nothing, so it required nothing, and
// asserting it meant no more than "this is an object" -- while the list beside
// it, [], is length-checked and really does mean empty. Two containers, two
// strengths, and nothing in a Recipe to tell them apart.
//
// It was found by a mutation that could not be killed: seeding a retirement
// into Hex's retirements map left the case asserting `retirements: {}` green.
// Four cases in the collection assert an empty object -- Asana's delete
// receipt, DigitalOcean's links when everything fits on one page, Jira's error
// map and Hex's retirements -- and all four say in their own names that they
// mean empty. None of them was checking it.
func TestAnEmptyObjectExpectationMeansEmpty(t *testing.T) {
	for _, c := range []struct {
		name  string
		want  any
		got   any
		fails bool
	}{
		{
			name:  "empty against empty passes",
			want:  map[string]any{},
			got:   map[string]any{},
			fails: false,
		},
		{
			name:  "empty against populated fails",
			want:  map[string]any{},
			got:   map[string]any{"1.8.0": map[string]any{"reason": "security"}},
			fails: true,
		},
		{
			name:  "empty against a list still fails on the type",
			want:  map[string]any{},
			got:   []any{},
			fails: true,
		},
		{
			// The partial rule is untouched: naming one key says nothing about
			// the others, which is what every ordinary case relies on.
			name:  "a named key ignores its neighbours",
			want:  map[string]any{"a": float64(1)},
			got:   map[string]any{"a": float64(1), "b": float64(2)},
			fails: false,
		},
		{
			name:  "and still fails when the named key is wrong",
			want:  map[string]any{"a": float64(1)},
			got:   map[string]any{"a": float64(2), "b": float64(2)},
			fails: true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			failures := compareValue("body", c.want, c.got)

			if got := len(failures) > 0; got != c.fails {
				t.Errorf("failures = %v, want a failure: %v", failures, c.fails)
			}
		})
	}
}
