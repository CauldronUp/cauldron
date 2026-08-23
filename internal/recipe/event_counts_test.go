package recipe_test

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The README states how many events these Recipes declare and how many of them
// are not lifecycle events. Both are the kind of figure that has drifted every
// time this project has looked: the case total, the detection coverage, the
// count of cases nobody has watched, the numeric identifiers, the unnamed
// paging. Seven for seven so far.
//
// The second number rests on a word list, which makes this test the working
// definition of "lifecycle event" rather than an independent check of one.
// That is worth saying plainly. Its job is to stop the README's figure going
// stale without anybody noticing, not to settle what the category means.
func TestTheREADMECountsDeclaredEvents(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	lifecycle := []string{"creat", "updat", "delet", "chang", "add", "remov", "new_", "destroy"}

	declared, notLifecycle := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, event := range r.Webhooks.Events {
			declared++

			low := strings.ToLower(event)

			found := false
			for _, word := range lifecycle {
				if strings.Contains(low, word) {
					found = true

					break
				}
			}

			if !found {
				notLifecycle++
			}
		}
	}

	stated := regexp.MustCompile(`Of (\d+) events declared across these Recipes, (\d+) are`).FindStringSubmatch(string(readme))
	if stated == nil {
		t.Fatal("the README no longer states the event figures in the form it did; update this test with them")
	}

	if n, _ := strconv.Atoi(stated[1]); n != declared {
		t.Errorf("the README says %d events are declared and %d are", n, declared)
	}

	if n, _ := strconv.Atoi(stated[2]); n != notLifecycle {
		t.Errorf("the README says %d are not lifecycle events and %d are not", n, notLifecycle)
	}
}
