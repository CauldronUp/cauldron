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

	// A spelled number of any length, because they get longer. This capture
	// was one word until the count reached a hundred, then two until it
	// reached a hundred and two. Widening it each time is the wrong repair --
	// see spelled below.
	stated := regexp.MustCompile(`([A-Z][a-z-]*(?: (?:and )?[a-z-]+)*) Recipes send at least one identifier as a number`).FindStringSubmatch(string(raw))
	if stated == nil {
		t.Fatal("the backlog no longer states how many Recipes send a numeric identifier; update this test with it")
	}

	if spelled(stated[1]) != counted {
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
	"Forty": 40, "Forty-one": 41, "Forty-two": 42, "Forty-three": 43,
	"Forty-four": 44, "Forty-five": 45, "Forty-six": 46, "Forty-seven": 47,
	"Forty-eight": 48, "Forty-nine": 49, "Fifty": 50,
	"Fifty-one": 51, "Fifty-two": 52, "Fifty-three": 53,
	"Fifty-four": 54, "Fifty-five": 55, "Fifty-six": 56,
	"Fifty-seven": 57, "Fifty-eight": 58, "Fifty-nine": 59,
	"Sixty": 60, "Sixty-one": 61, "Sixty-two": 62,
	"Sixty-three": 63, "Sixty-four": 64, "Sixty-five": 65,
	"Sixty-six": 66, "Sixty-seven": 67, "Sixty-eight": 68,
	"Sixty-nine": 69, "Seventy": 70, "Seventy-one": 71,
	"Seventy-two": 72, "Seventy-three": 73, "Seventy-four": 74,
	"Seventy-five": 75, "Seventy-six": 76, "Seventy-seven": 77,
	"Seventy-eight": 78, "Seventy-nine": 79, "Eighty": 80,
	"Eighty-one": 81, "Eighty-two": 82, "Eighty-three": 83,
	"Eighty-four": 84, "Eighty-five": 85, "Eighty-six": 86,
	"Eighty-seven": 87, "Eighty-eight": 88, "Eighty-nine": 89,
	"Ninety": 90, "Ninety-one": 91, "Ninety-two": 92,
	"Ninety-three": 93, "Ninety-four": 94, "Ninety-five": 95,
	"Ninety-six": 96, "Ninety-seven": 97, "Ninety-eight": 98,
	"Ninety-nine": 99, "One hundred": 100,
}

// spelled reads a number this document writes out, including the ones past a
// hundred that no table would keep up with.
//
// The table above stopped at a hundred, and the check that read it was a map
// lookup, so the day a count reached 101 the test failed with "the backlog
// says and two Recipes..." -- a message about its own regex, not about the
// figure. It had already been widened once, from one captured word to two,
// when the count first hit a hundred. Widening it again would buy two more
// counts.
//
// So: split on the hundreds, look the remainder up in the table, and add.
// Everything below a hundred is still a table lookup, because English below a
// hundred is a table and not a grammar.
func spelled(phrase string) int {
	if value, ok := words[phrase]; ok {
		return value
	}

	hundreds, rest, found := strings.Cut(phrase, " hundred")
	if !found {
		return 0
	}

	count, ok := words[hundreds]
	if !ok {
		return 0
	}

	total := count * 100

	rest = strings.TrimPrefix(strings.TrimSpace(rest), "and ")
	if rest == "" {
		return total
	}

	// The remainder is written lowercase inside a sentence -- "One hundred and
	// two" -- and the table is keyed by the capitalised form it takes when it
	// starts one.
	remainder, ok := words[capitalised(rest)]
	if !ok {
		return 0
	}

	return total + remainder
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

// The test above compares a queued row to a Recipe's directory name exactly.
// Directory names are short and product names are not, so two providers sat in
// the queue for months after they shipped: Hetzner Cloud, under recipes/hetzner,
// and Fly.io, under recipes/fly. The Hetzner row was the pointed one. It asked
// for "every mutation returns an action object you have to poll, rather than the
// thing you changed", which is the sentence the shipped Recipe opens with.
//
// Widening the match is not free, and that is the whole difficulty. Of the nine
// pairs a plain substring test finds, two are the real ones and the other seven
// are not. Circle is the stablecoin company and not CircleCI, which is a
// homonym. GitHub Actions and Stripe Tax are separate surfaces of a vendor
// already here, and both rows say so themselves -- "probably its own",
// "probably a specialised extension". GoCardless sat here for the same reason
// until it shipped in its own right -- payments, where gocardlessbank is
// GoCardless Bank Account Data: one vendor, two APIs,
// different shapes. Cloudflare R2 is an S3-compatible object store that shares
// a brand with the Cloudflare Recipe and nothing else. And two are accidents of
// spelling -- "lob" sits inside "Azure Blob Storage", and "kit" inside
// "buildkite".
//
// So the machine cannot decide this and this test does not try. It finds
// candidates and requires each one to be resolved -- struck through, or written
// down below with the reason it is a different thing. A silent blind spot
// becomes a decision somebody has to record, which is the most a name
// comparison can honestly do.
func TestTheBacklogDoesNotQueueAShippedProviderUnderALongerName(t *testing.T) {
	// Candidates that are genuinely different products, and why. A row named
	// here must still be a candidate: an entry that stops matching is stale and
	// fails, for the same reason the queue itself does.
	distinct := map[string]string{
		"Sage":                   "Sage Business Cloud Accounting, a separate product from Sage Intacct, which ships",
		"Grafana Cloud":          "a stack-management API on top of the Grafana HTTP API, which is what ships",
		"Cloudflare R2":          "an S3-compatible object store sharing a brand with the Cloudflare Recipe",
		"GitHub Actions":         "a separate surface of GitHub, as the row itself says",
		"Stripe Tax":             "a specialised extension of Stripe, as the row itself says",
		"Hugging Face Inference": "the model-serving API, where huggingface is the Hub registry API",
		"Wise Business":          "the business banking product, where wise is the transfers API",
		"Codat":                  "accounting-data aggregation, an unrelated company that shares four letters with Coda",
		"Coinbase Commerce":      "the crypto-payments product, where coinbase is the exchange and market-data API",
	}

	raw, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	shipped := recipe.Bundled()

	header := regexp.MustCompile(`^\|\s*Provider\s*\|\s*Why\s*\|`)
	row := regexp.MustCompile(`^\| ([^|]+?) \|`)

	matched := map[string]bool{}

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

		for _, under := range namedBy(name, shipped) {
			matched[name] = true

			if _, known := distinct[name]; known {
				continue
			}

			t.Errorf("the backlog queues %q and recipes/%s ships; strike it through, or record here why they are different things", name, under)
		}
	}

	for name := range distinct {
		if !matched[name] {
			t.Errorf("%q is recorded as distinct from a shipped Recipe and no longer looks like one; drop the entry", name)
		}
	}
}

// namedBy returns the shipped Recipes whose directory name looks like the
// provider a queued row names.
//
// Exact-equal is the caller's job; this is the near miss. A row matches a
// Recipe when one of its words is that Recipe's whole name -- "Fly.io" against
// fly, "Hetzner Cloud" against hetzner -- or when either name is a prefix of
// the other, which is how "Circle" reaches circleci. The prefix arms need four
// characters to fire, because three-letter Recipe names turn up inside ordinary
// words: without the floor, lob matches "Azure Blob Storage".
func namedBy(row string, shipped []string) []string {
	whole := simplify(row)

	words := strings.FieldsFunc(strings.ToLower(row), func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})

	var found []string

	for _, name := range shipped {
		if name == whole {
			continue
		}

		var looksLike bool

		for _, word := range words {
			if word == name {
				looksLike = true
				break
			}
		}

		if len(name) >= 4 {
			looksLike = looksLike ||
				strings.HasPrefix(name, whole) ||
				strings.HasPrefix(whole, name)
		}

		if looksLike {
			found = append(found, name)
		}
	}

	return found
}

// The three counters written after the paging pair, guarded the same way.
//
// The reason is on the page above them: a table of routes stuck for want of a
// provider description sat there for weeks naming AssemblyAI, Mercury and
// QuickBooks, all three long since settled, and missing six Recipes that had
// joined the list. Nothing checked it, so nothing noticed. The two paging
// figures beside it were guarded and stayed right.
//
// Each of these reads the figure out of the block quoting verify's own output,
// so the page and the tool cannot disagree about what the tool says.
func TestTheBacklogCountsUnsentPagingNames(t *testing.T) {
	names, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnsentPagingParam(); n > 0 {
			names += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) paging parameter name\(s\) across (\d+) recipe\(s\) are declared and sent by no`)

	if n, _ := strconv.Atoi(stated[1]); n != names {
		t.Errorf("the backlog says %d paging names are sent by no case and %d are", n, names)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those names across %d Recipes and they are across %d", n, recipes)
	}
}

func TestTheBacklogCountsUnshownListings(t *testing.T) {
	listings, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnshownListing(); n > 0 {
			listings += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) listing\(s\) across (\d+) recipe\(s\) have no case that answers them`)

	if n, _ := strconv.Atoi(stated[1]); n != listings {
		t.Errorf("the backlog says %d listings are answered by nothing and %d are", n, listings)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those listings across %d Recipes and they are across %d", n, recipes)
	}
}

func TestTheBacklogCountsHollowListings(t *testing.T) {
	listings, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.HollowListing(); n > 0 {
			listings += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) more across (\d+) recipe\(s\) are answered only empty`)

	if n, _ := strconv.Atoi(stated[1]); n != listings {
		t.Errorf("the backlog says %d listings are answered only empty and %d are", n, listings)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those listings across %d Recipes and they are across %d", n, recipes)
	}
}

// backlogFigure pulls one two-number claim out of the backlog, and says so
// plainly when the sentence it is looking for has been rewritten -- a figure
// this file stops stating is a decision, not a failure, and the test that
// guards it should be changed deliberately rather than deleted in a panic.
func backlogFigure(t *testing.T, pattern string) []string {
	t.Helper()

	backlog, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("read backlog: %v", err)
	}

	found := regexp.MustCompile(pattern).FindStringSubmatch(string(backlog))
	if found == nil {
		t.Fatalf("the backlog no longer states the figure matching %q; update this test with the new wording", pattern)
	}

	return found
}

func TestTheBacklogCountsUnassertedFields(t *testing.T) {
	fields, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnassertedField(); n > 0 {
			fields += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) record field name\(s\) across (\d+) recipe\(s\) are declared and asserted`)

	if n, _ := strconv.Atoi(stated[1]); n != fields {
		t.Errorf("the backlog says %d field names are asserted by no case and %d are", n, fields)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those names across %d Recipes and they are across %d", n, recipes)
	}
}

func TestTheBacklogCountsUnservedFields(t *testing.T) {
	fields, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnservedField(); n > 0 {
			fields += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) of those across (\d+) recipe\(s\) no fixture sets`)

	if n, _ := strconv.Atoi(stated[1]); n != fields {
		t.Errorf("the backlog says %d field names no fixture sets and %d are", n, fields)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those names across %d Recipes and they are across %d", n, recipes)
	}
}

func TestTheBacklogCountsUnshownErrors(t *testing.T) {
	errors, recipes := 0, 0

	for _, name := range recipe.Bundled() {
		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if n := r.UnshownError(); n > 0 {
			errors += n
			recipes++
		}
	}

	stated := backlogFigure(t,
		`(\d+) declared error\(s\) across (\d+) recipe\(s\) have no case asserting`)

	if n, _ := strconv.Atoi(stated[1]); n != errors {
		t.Errorf("the backlog says %d errors have no case asserting them and %d do", n, errors)
	}

	if n, _ := strconv.Atoi(stated[2]); n != recipes {
		t.Errorf("the backlog spreads those errors across %d Recipes and they are across %d", n, recipes)
	}
}
