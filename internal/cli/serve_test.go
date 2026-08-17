package cli

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func projectWith(t *testing.T, composer string) string {
	t.Helper()

	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "composer.json"), []byte(composer), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	return root
}

func TestServeFlagDefaults(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseServeFlags(nil, &stderr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if opts.port != defaultPort {
		t.Errorf("port = %d, want %d", opts.port, defaultPort)
	}

	// A fake must not squat on the port the application under development
	// is most likely to want.
	for _, forbidden := range []int{80, 3000, 8000, 8080} {
		if opts.port == forbidden {
			t.Errorf("default port must not be %d", forbidden)
		}
	}

	if opts.seed != 1 {
		t.Errorf("seed = %d, want 1", opts.seed)
	}
}

func TestServeFlagsAreParsed(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseServeFlags([]string{"-port", "5000", "-seed", "99", "-fixture", "small-shop", "stripe"}, &stderr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if opts.port != 5000 || opts.seed != 99 || opts.fixture != "small-shop" {
		t.Errorf("opts = %+v", opts)
	}

	if len(opts.recipes) != 1 || opts.recipes[0] != "stripe" {
		t.Errorf("recipes = %v", opts.recipes)
	}
}

func TestPlanUsesNamedRecipes(t *testing.T) {
	mount, missing, err := plan(serveOptions{recipes: []string{"stripe"}})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	if len(mount) != 1 || mount[0] != "stripe" {
		t.Errorf("mount = %v", mount)
	}

	if len(missing) != 0 {
		t.Errorf("missing = %v", missing)
	}
}

func TestPlanDetectsRecipesFromTheProject(t *testing.T) {
	dir := projectWith(t, `{"require":{"stripe/stripe-php":"^17.0"}}`)

	mount, _, err := plan(serveOptions{dir: dir})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	if len(mount) != 1 || mount[0] != "stripe" {
		t.Errorf("mount = %v, want [stripe]", mount)
	}
}

// A detected provider Cauldron cannot emulate must be reported, never silently
// dropped, or the developer believes it is faked when it is not.
//
// OpenAI is used here precisely because no Recipe ships for it. If one ever
// does, this test should fail and be repointed rather than deleted, which is
// what happened to the Adyen it used to name.
func TestPlanReportsDetectedRecipesItCannotServe(t *testing.T) {
	dir := projectWith(t, `{"require":{"stripe/stripe-php":"^17.0","openai-php/client":"^0.10"}}`)

	mount, missing, err := plan(serveOptions{dir: dir})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	if len(mount) != 1 || mount[0] != "stripe" {
		t.Errorf("mount = %v", mount)
	}

	if len(missing) != 1 || missing[0] != "openai" {
		t.Errorf("missing = %v, want [openai]", missing)
	}
}

func TestPlanFailsWhenNothingCanBeServed(t *testing.T) {
	dir := projectWith(t, `{"require":{"openai-php/client":"^0.10"}}`)

	_, missing, err := plan(serveOptions{dir: dir})
	if err == nil {
		t.Fatal("expected an error")
	}

	if len(missing) != 1 {
		t.Errorf("the unservable providers should still be reported; got %v", missing)
	}
}

func TestPlanFailsOutsideAProject(t *testing.T) {
	_, _, err := plan(serveOptions{dir: t.TempDir()})
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "cauldron serve stripe") {
		t.Errorf("the error should suggest naming a recipe, got %q", err)
	}
}

func TestServeBannerShowsAddressesAndGaps(t *testing.T) {
	var out bytes.Buffer

	writeServeBanner(&out, "127.0.0.1:4600", []string{"stripe"}, []string{"twilio"}, serveOptions{fixture: "small-shop"})

	text := out.String()

	for _, want := range []string{
		"http://127.0.0.1:4600/stripe",
		"small-shop",
		"still reach the real network",
		"twilio",
		"_cauldron/status",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("banner is missing %q\n%s", want, text)
		}
	}
}

func TestServeIsListedInHelp(t *testing.T) {
	stdout, _, code := run(t, "help")

	if code != 0 {
		t.Fatalf("exit = %d", code)
	}

	if !strings.Contains(stdout, "serve") {
		t.Errorf("help should list serve\n%s", stdout)
	}
}

// Providers name their fixtures differently: stripe ships small-shop, github
// ships small-repo. A single global name must therefore be a convenience
// applied where it fits, not an instruction that fails the whole run.
func TestFixtureChoiceParsesBothForms(t *testing.T) {
	global := parseFixture("small-shop")

	if fixture, explicit := global.forRecipe("stripe"); fixture != "small-shop" || explicit {
		t.Errorf("global: got %q explicit=%v", fixture, explicit)
	}

	pairs := parseFixture("stripe=small-shop,github=small-repo")

	if fixture, explicit := pairs.forRecipe("github"); fixture != "small-repo" || !explicit {
		t.Errorf("pair: got %q explicit=%v", fixture, explicit)
	}

	if fixture, _ := pairs.forRecipe("twilio"); fixture != "" {
		t.Errorf("a recipe with no pair and no global should get nothing; got %q", fixture)
	}
}

func TestFixtureChoiceMixesGlobalAndPairs(t *testing.T) {
	choice := parseFixture("small-shop,github=small-repo")

	if fixture, explicit := choice.forRecipe("github"); fixture != "small-repo" || !explicit {
		t.Errorf("an explicit pair must win over the global; got %q explicit=%v", fixture, explicit)
	}

	if fixture, explicit := choice.forRecipe("stripe"); fixture != "small-shop" || explicit {
		t.Errorf("stripe should fall back to the global; got %q explicit=%v", fixture, explicit)
	}
}

func TestFixtureChoiceIgnoresEmptyParts(t *testing.T) {
	choice := parseFixture(" , stripe=small-shop , ")

	if fixture, _ := choice.forRecipe("stripe"); fixture != "small-shop" {
		t.Errorf("got %q", fixture)
	}
}

// Cauldron's job is the emulated providers. The application, its database and
// its web server belong to whoever is already running them, so there has to be
// a way to get the providers and nothing else: no plan describing an
// environment Cauldron will not set up, no containers, and output a program
// can read rather than a banner addressed to a person.
func TestHeadlessServeReportsItselfAsJSON(t *testing.T) {
	var out bytes.Buffer

	writeServeJSON(&out, "127.0.0.1:4600", loopback, 4600,
		[]string{"stripe", "woocommerce"}, []string{"twilio"}, nil)

	line := out.String()

	if strings.Count(strings.TrimSpace(line), "\n") != 0 {
		t.Errorf("headless output must be one line, so a caller can read one and move on:\n%s", line)
	}

	var payload struct {
		Address string `json:"address"`
		Bind    string `json:"bind"`
		Port    int    `json:"port"`
		Control string `json:"control"`
		Mounted []struct {
			Recipe string `json:"recipe"`
			URL    string `json:"url"`
		} `json:"mounted"`
		Missing  []string `json:"missing"`
		Unseeded []string `json:"unseeded"`
	}

	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("headless output is not JSON: %v\n%s", err, line)
	}

	if payload.Address != "http://127.0.0.1:4600" {
		t.Errorf("address = %q", payload.Address)
	}

	if payload.Port != 4600 || payload.Bind != loopback {
		t.Errorf("bind = %q port = %d, want them reported separately from the address", payload.Bind, payload.Port)
	}

	if len(payload.Mounted) != 2 || payload.Mounted[0].URL != "http://127.0.0.1:4600/stripe" {
		t.Errorf("mounted = %+v", payload.Mounted)
	}

	// A provider Cauldron cannot emulate still reaches the real network.
	// Leaving that out of the machine-readable form and keeping it only in the
	// banner would hide it from exactly the callers that cannot see a banner.
	if len(payload.Missing) != 1 || payload.Missing[0] != "twilio" {
		t.Errorf("missing = %v, want the unemulated providers reported", payload.Missing)
	}

	// Empty rather than null, so a caller can iterate without checking.
	if payload.Unseeded == nil {
		t.Error("unseeded should be an empty array rather than null")
	}
}

// A wildcard bind reports itself as [::]:4600, which is not a URL. Handing
// that back as the address gives a caller something that fails the moment it
// is used, so the reported address is always one that works from this machine
// and the real bind travels beside it.
func TestAWildcardBindStillReportsAUsableAddress(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Skipf("cannot bind a wildcard address here: %v", err)
	}

	defer func() { _ = listener.Close() }()

	addr := dialable(listener.Addr())

	if strings.HasPrefix(addr, "[::]") || strings.HasPrefix(addr, "0.0.0.0") {
		t.Errorf("addr = %q, which is not something anybody can dial", addr)
	}

	if !strings.HasPrefix(addr, loopback+":") {
		t.Errorf("addr = %q, want loopback with the bound port", addr)
	}
}

func TestHeadlessAndHostAreParsed(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseServeFlags([]string{"-headless", "-host", "0.0.0.0", "stripe"}, &stderr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if !opts.headless {
		t.Error("headless should be set")
	}

	if opts.host != "0.0.0.0" {
		t.Errorf("host = %q", opts.host)
	}

	// Loopback by default: a fake provider that answers the network by
	// accident is worse than one that is awkward to reach.
	def, err := parseServeFlags(nil, &stderr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if def.host != loopback {
		t.Errorf("default host = %q, want loopback", def.host)
	}

	if def.headless {
		t.Error("headless must be opt-in")
	}
}
