package recipe_test

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The backlog is the queue of providers worth a Recipe, and a shipped one is
// struck through. Nothing checked that, so the file drifted into
// contradicting itself: six providers were listed as queued in one section
// and recorded as shipped in another, and the paragraph stating the standard
// named Linear and Attio as examples of providers left out while both had
// since been built.
//
// A queue nobody prunes stops being a queue, and this one is the answer to
// "what is worth doing next" -- which is exactly the question a stale entry
// wastes somebody's evening on.
func TestTheBacklogDoesNotQueueShippedProviders(t *testing.T) {
	raw, err := os.ReadFile("../../docs/backlog.md")
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	shipped := map[string]bool{}
	for _, name := range recipe.Bundled() {
		shipped[simplify(name)] = true
	}

	header := regexp.MustCompile(`^\|\s*Provider\s*\|\s*Why\s*\|`)
	row := regexp.MustCompile(`^\| ([^|]+?) \|`)

	var inQueue bool

	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimRight(line, "\r")

		if header.MatchString(line) {
			inQueue = true
			continue
		}

		if !strings.HasPrefix(line, "|") {
			inQueue = false
			continue
		}

		if !inQueue {
			continue
		}

		m := row.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		name := strings.TrimSpace(m[1])
		if strings.HasPrefix(name, "~~") || strings.HasPrefix(name, "---") {
			continue
		}

		if shipped[simplify(name)] {
			t.Errorf("the backlog queues %q and a Recipe for it ships; strike it through", name)
		}
	}
}

// The backlog counts the Recipes that send an identifier as a number. That
// count is the answer to "how much of this is done", and it was written when
// sixteen were and stayed there while ten more landed.
//
// Every number in these documents that nobody checks has drifted -- the
// README's case total, its detection coverage, its count of cases nobody has
// watched a provider perform, and this. The pattern is consistent enough to
// stop treating each one as a surprise.
func TestTheBacklogCountsNumericIdentifiers(t *testing.T) {
	raw, err := os.ReadFile("../../docs/backlog.md")
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	counted := 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, resource := range r.Resources {
			if resource.ID.Type == "number" {
				counted++
				break
			}
		}
	}

	stated := regexp.MustCompile(`([A-Za-z-]+) Recipes send at least one identifier as a number`).FindStringSubmatch(string(raw))
	if stated == nil {
		t.Fatal("the backlog no longer states how many Recipes send a numeric identifier; update this test with it")
	}

	if words[stated[1]] != counted {
		t.Errorf("the backlog says %s Recipes send a numeric identifier and %d do", stated[1], counted)
	}
}

// words are the numbers this document spells out, which is how it is written
// and not worth changing for a test's convenience.
var words = map[string]int{
	"One": 1, "Two": 2, "Three": 3, "Four": 4, "Five": 5,
	"Six": 6, "Seven": 7, "Eight": 8, "Nine": 9, "Ten": 10,
	"Eleven": 11, "Twelve": 12, "Thirteen": 13, "Fourteen": 14, "Fifteen": 15,
	"Sixteen": 16, "Seventeen": 17, "Eighteen": 18, "Nineteen": 19,
	"Twenty": 20, "Twenty-one": 21, "Twenty-two": 22, "Twenty-three": 23,
	"Twenty-four": 24, "Twenty-five": 25, "Twenty-six": 26, "Twenty-seven": 27,
	"Twenty-eight": 28, "Twenty-nine": 29, "Thirty": 30, "Thirty-one": 31,
	"Thirty-two": 32, "Thirty-three": 33, "Thirty-four": 34, "Thirty-five": 35,
	"Thirty-six": 36, "Thirty-seven": 37, "Thirty-eight": 38, "Thirty-nine": 39,
	"Forty": 40,
}

// The backlog states how many routes still page by a parameter nobody named.
// It is the figure that moves every time one is settled, which makes it the
// most likely of the lot to be left behind -- and the paging section already
// carries three figures from the sweep that closed it, all of them historical
// and none of them what is true now.
func TestTheBacklogCountsUnnamedPaging(t *testing.T) {
	backlog, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	routes, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.GuessedPagination(); n > 0 {
			routes += n
			recipes++
		}
	}

	stated := regexp.MustCompile(`\*\*(\d+) routes across (\d+) Recipes\*\* still page by a parameter nobody named`).FindStringSubmatch(string(backlog))
	if stated == nil {
		t.Fatal("the backlog no longer states the unnamed-paging figure in the form it did; update this test with it")
	}

	if n, _ := strconv.Atoi(stated[1]); n != routes {
		t.Errorf("the backlog says %d routes page by a parameter nobody named and %d do", n, routes)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog says those routes are spread across %d Recipes and they are across %d", n, recipes)
	}
}

// The larger of the two paging figures, and the one that was missing entirely
// until the counter behind it was written. Same reason as its neighbour: it
// moves whenever a listing gains a paging block, which is every time somebody
// settles one.
func TestTheBacklogCountsUnstatedPaging(t *testing.T) {
	backlog, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	routes, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnstatedPagination(); n > 0 {
			routes += n
			recipes++
		}
	}

	stated := regexp.MustCompile(`([0-9]+) more listings across ([0-9]+) Recipes`).FindStringSubmatch(string(backlog))
	if stated == nil {
		t.Fatal("the backlog no longer states the unstated-paging figure in the form it did; update this test with it")
	}

	if n, _ := strconv.Atoi(stated[1]); n != routes {
		t.Errorf("the backlog says %d listings declare no paging at all and %d do", n, routes)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog says those listings are spread across %d Recipes and they are across %d", n, recipes)
	}
}
