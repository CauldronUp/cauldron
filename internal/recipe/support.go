// Small helpers shared across this package.

package recipe

import (
	"sort"
	"strings"
)

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// splitNonEmpty splits a dotted path, and returns nothing for an empty one so
// a field that does not nest is not asked to justify its parent.
func splitNonEmpty(path string) []string {
	if path == "" {
		return nil
	}

	return strings.Split(path, ".")
}

func all(text string, predicate func(byte) bool) bool {
	for i := 0; i < len(text); i++ {
		if !predicate(text[i]) {
			return false
		}
	}

	return true
}

func isDecimalDigit(b byte) bool { return b >= '0' && b <= '9' }

func isHexDigit(b byte) bool {
	return isDecimalDigit(b) || (b >= 'a' && b <= 'f')
}
