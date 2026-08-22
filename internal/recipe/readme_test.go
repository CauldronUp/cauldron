package recipe_test

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// readmePath is the repository README, relative to this package.
const readmePath = "../../README.md"

// The README makes three numeric claims about the project: how many Recipes
// ship, which ones, and how many conformance cases they carry. Nothing checked
// them, which is an odd gap in a project whose whole argument is that its
// claims are verified rather than asserted.
//
// They have been maintained by hand on every pull request so far and have
// stayed accurate, but "accurate so far" is exactly the kind of confidence
// this project is meant to replace with a check.
func TestTheREADMECountsMatchTheRecipes(t *testing.T) {
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	text := string(readme)
	bundled := recipe.Bundled()

	shipped := regexp.MustCompile(`\| Recipes shipped \| (\d+): (.*?) \|`).FindStringSubmatch(text)
	if shipped == nil {
		t.Fatal("the README no longer states a recipe count in the form it did; update this test with it")
	}

	stated, err := strconv.Atoi(shipped[1])
	if err != nil {
		t.Fatalf("recipe count is not a number: %v", err)
	}

	if stated != len(bundled) {
		t.Errorf("the README says %d Recipes ship and %d do", stated, len(bundled))
	}

	listed := strings.Split(shipped[2], ",")
	if len(listed) != len(bundled) {
		t.Errorf("the README names %d Recipes and %d ship", len(listed), len(bundled))
	}

	// Display names are not directory names: "AWS SQS" is sqs, "Help Scout" is
	// helpscout, "Google Pub/Sub" is pubsub. Comparing them stripped of case
	// and punctuation is enough to catch a Recipe that shipped without being
	// added to the list, which is the mistake worth catching.
	var flattened []string
	for _, name := range listed {
		flattened = append(flattened, simplify(name))
	}

	for _, name := range bundled {
		found := false

		for _, listedName := range flattened {
			if strings.Contains(listedName, simplify(name)) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("%s ships and the README does not list it", name)
		}
	}

	cases := 0
	for _, name := range bundled {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		cases += len(r.Conformance)
	}

	statedCases := regexp.MustCompile(`Working\. (\d+) cases`).FindStringSubmatch(text)
	if statedCases == nil {
		t.Fatal("the README no longer states a case count in the form it did; update this test with it")
	}

	counted, err := strconv.Atoi(statedCases[1])
	if err != nil {
		t.Fatalf("case count is not a number: %v", err)
	}

	if counted != cases {
		t.Errorf("the README says %d conformance cases and there are %d", counted, cases)
	}

	// The same number appears twice in the README, and the two have to agree
	// with each other as well as with the Recipes.
	//
	// So does the split beside it. That sentence is the project's own answer
	// to "how do you know it behaves like the real thing", and one of its
	// three numbers was checked by nothing: the total was maintained on every
	// pull request and the count of cases nobody has watched a provider
	// perform drifted to 1548 while the truth reached 1680. A stale figure in
	// the paragraph about not overstating evidence is the worst place to keep
	// one, so all three are counted here now.
	repeated := regexp.MustCompile(`Recipe: (\d+) cases, (\d+) run against\s+a live account and (\d+) not`).FindStringSubmatch(text)
	if repeated == nil {
		t.Fatal("the README no longer states the case split in the form it did; update this test with it")
	}

	if repeated[1] != statedCases[1] {
		t.Errorf("the README states the case count twice and they disagree: %s and %s", statedCases[1], repeated[1])
	}

	observed := 0
	for _, name := range bundled {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, c := range r.Conformance {
			if c.Verified != "" {
				observed++
			}
		}
	}

	if statedObserved, _ := strconv.Atoi(repeated[2]); statedObserved != observed {
		t.Errorf("the README says %d cases ran against a live account and %d carry a verified date", statedObserved, observed)
	}

	if statedRead, _ := strconv.Atoi(repeated[3]); statedRead != cases-observed {
		t.Errorf("the README says %d cases were never run against a live account and %d were not", statedRead, cases-observed)
	}

	// The sample output above that sentence is a claim about what the command
	// prints, and it had drifted too: it showed Stripe with eight cases when
	// Stripe had nine. A README that demonstrates a verification tool with
	// output the tool does not produce is making the reader's point for them.
	sample := regexp.MustCompile(`stripe 0\.1\.0\s+(\d+) of (\d+) cases passed`).FindStringSubmatch(text)
	if sample == nil {
		t.Fatal("the README no longer shows sample verify output in the form it did; update this test with it")
	}

	stripe, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	if shown, _ := strconv.Atoi(sample[2]); shown != len(stripe.Conformance) {
		t.Errorf("the README shows stripe with %d cases and it has %d", shown, len(stripe.Conformance))
	}

	if sample[1] != sample[2] {
		t.Errorf("the README's sample shows %s of %s passing, which is not the output of a green run", sample[1], sample[2])
	}
}

// simplify reduces a name to letters and digits, so "Google Pub/Sub" and
// "pubsub" can be compared.
func simplify(name string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}

	return b.String()
}
