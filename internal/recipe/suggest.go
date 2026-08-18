package recipe

import "strings"

// Suggest names the bundled Recipes closest to what somebody typed.
//
// A hundred and twenty-seven names is too many to print. It was a reasonable
// thing to do at four, and the message that once listed everything now fills a
// terminal and buries the answer, which is worse than saying nothing: somebody
// who typed "stripo" has to read a wall of text to find "stripe" in it.
//
// Names like gocardlessbank, moderntreasury, npmregistry and secretsmanager
// make a typo close to certain, so the useful answer is the two or three
// nearest ones and a pointer to the full list.
func Suggest(name string) []string {
	typed := strings.ToLower(name)

	type scored struct {
		name     string
		distance int
	}

	var close []scored

	for _, candidate := range Bundled() {
		distance := editDistance(typed, candidate)

		// A third of the length, so a short name tolerates one mistake and a
		// long one tolerates a few. A fixed threshold either misses
		// "gocardlesbank" or suggests half the list for "sms".
		budget := len(candidate)/3 + 1
		if budget > 4 {
			budget = 4
		}

		// A prefix somebody stopped typing is a match whatever the distance:
		// "gocardless" is not within budget of "gocardlessbank" and is
		// obviously what they meant.
		if strings.HasPrefix(candidate, typed) && len(typed) >= 3 {
			distance = 0
		}

		if distance <= budget {
			close = append(close, scored{name: candidate, distance: distance})
		}
	}

	// Nearest first, and alphabetically within a distance so the answer does
	// not depend on map ordering.
	for i := 1; i < len(close); i++ {
		for j := i; j > 0; j-- {
			if close[j].distance < close[j-1].distance {
				close[j], close[j-1] = close[j-1], close[j]

				continue
			}

			break
		}
	}

	out := make([]string, 0, 3)

	for _, s := range close {
		if len(out) == 3 {
			break
		}

		out = append(out, s.name)
	}

	return out
}

// editDistance is Levenshtein, which is enough for names this short. It counts
// a transposition as two edits, so "sripte" is further from "stripe" than a
// human would judge it, and the budget above absorbs that.
func editDistance(a, b string) int {
	if a == b {
		return 0
	}

	if a == "" {
		return len(b)
	}

	if b == "" {
		return len(a)
	}

	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)

	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current[0] = i

		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			current[j] = min3(
				current[j-1]+1,
				previous[j]+1,
				previous[j-1]+cost,
			)
		}

		copy(previous, current)
	}

	return previous[len(b)]
}

func min3(a, b, c int) int {
	if b < a {
		a = b
	}

	if c < a {
		a = c
	}

	return a
}
