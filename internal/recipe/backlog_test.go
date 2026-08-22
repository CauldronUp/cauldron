package recipe_test

import (
	"os"
	"regexp"
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
