// Questions the rest of the program asks a Recipe.

package recipe

import (
	"sort"
)

// Verified reports how many conformance cases were observed against the real
// API, and how many rest on documentation alone. The distinction is the whole
// value of the suite, so it is reported rather than averaged away.
func (r *Recipe) Verified() (observed, documented int) {
	for _, c := range r.Conformance {
		if c.Verified != "" {
			observed++
			continue
		}

		documented++
	}

	return observed, documented
}

// Events returns the webhook event names this Recipe can emit.
func (r *Recipe) Events() []string {
	out := make([]string, len(r.Webhooks.Events))
	copy(out, r.Webhooks.Events)
	sort.Strings(out)

	return out
}
