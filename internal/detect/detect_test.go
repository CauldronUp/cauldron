package detect

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// writeProject materialises a temporary project directory from a map of
// filename to contents.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()

	root := t.TempDir()

	for name, contents := range files {
		path := filepath.Join(root, name)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}

		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	return root
}

func requirement(t *testing.T, p *Project, kind Kind, name string) Requirement {
	t.Helper()

	for _, r := range p.Requirements {
		if r.Kind == kind && r.Name == name {
			return r
		}
	}

	t.Fatalf("expected a %s requirement named %q; got %+v", kind, name, p.Requirements)

	return Requirement{}
}

// BigCommerce ships two client packages under different names and neither is
// called bigcommerce on npm, which is exactly why the table is explicit: a
// prefix rule would match @bigcommerce/checkout-sdk and miss node-bigcommerce
// entirely, and sending somebody to a Recipe they are not using is worse than
// not detecting at all.
func TestDetectBigCommerceFromEitherClient(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"bigcommerce/api": "^3.0"}}`},
		{"package.json": `{"dependencies": {"node-bigcommerce": "^4.0.0"}}`},
		{"package.json": `{"dependencies": {"@bigcommerce/checkout-sdk": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "bigcommerce") {
			t.Errorf("%v did not detect bigcommerce: %+v", manifest, p.Requirements)
		}
	}
}

// Etsy is the case where a prefix rule would be actively wrong rather than
// merely incomplete. Etsy the company publishes phan and phpunit-extensions
// under the etsy/ vendor name -- a static analyser and a test library, which
// have nothing to do with the Etsy API and which plenty of projects depend on
// without selling anything. Matching etsy/ would offer a commerce emulator to
// a project that runs a linter.
//
// The other half is version. The Recipe models Open API v3, and the older
// wrappers on both registries speak v2, which is a different API with a
// different envelope and a different auth scheme. Detecting one of those and
// pointing it here would hand back shapes its client cannot read, so the
// table names only the two packages that say v3.
func TestDetectEtsyFromEitherV3Client(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"rhysnhall/etsy-php-sdk": "^1.0"}}`},
		{"package.json": `{"dependencies": {"etsy-ts": "^8.0.1"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "etsy") {
			t.Errorf("%v did not detect etsy: %+v", manifest, p.Requirements)
		}
	}
}

// The companion to the test above: the packages that must not reach the Etsy
// Recipe. phan is a static analyser, phpunit-extensions is a test library, and
// oauth1-etsy is for the v2-era OAuth 1.0 flow the v3 API replaced.
func TestDetectDoesNotOfferEtsyForEtsysOwnToolingOrV2Clients(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"etsy/phan": "^5.0"}}`},
		{"composer.json": `{"require": {"etsy/phpunit-extensions": "^4.0"}}`},
		{"composer.json": `{"require": {"y0lk/oauth1-etsy": "^1.0"}}`},
		{"package.json": `{"dependencies": {"etsy-js": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "etsy") {
			t.Errorf("%v offered the etsy Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// Magento is the case where one name covers three protocols. The Recipe
// models the REST API at /rest/V1, and packages called magento talk to it,
// or to the SOAP API it replaced, or to the GraphQL API PWA Studio uses --
// three different wire formats, all correctly described as Magento.
//
// Detecting the wrong one is not a near miss. A SOAP client pointed at a REST
// emulator gets a 404 on its endpoint; a GraphQL client gets one on its. So
// the table names only packages that say REST and say 2.
func TestDetectMagentoFromRestClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"springimport/swagger-magento2-client": "^1.0"}}`},
		{"composer.json": `{"require": {"zero1/magento2-rest-client": "^1.0"}}`},
		{"package.json": `{"dependencies": {"magento2-rest-client": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "magento") {
			t.Errorf("%v did not detect magento: %+v", manifest, p.Requirements)
		}
	}
}

// The companion: magento-api is a SOAP wrapper, smalot/magento-client is SOAP
// v1, and @magento/peregrine is the PWA Studio runtime, which speaks GraphQL.
// None of them would get a usable answer from a REST emulator.
func TestDetectDoesNotOfferMagentoForSoapOrGraphQLClients(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"magento-api": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@magento/peregrine": "^14.0.0"}}`},
		{"composer.json": `{"require": {"smalot/magento-client": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "magento") {
			t.Errorf("%v offered the magento Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// eBay is the case where mapping by popularity would be worst. The two
// most-installed eBay packages on Packagist have well over a million installs
// between them and both speak the XML Trading, Finding and Shopping APIs --
// the ones the REST APIs this Recipe models were built to replace. Sending
// almost every PHP eBay project to an emulator that serves none of the
// endpoints it calls would be a worse outcome than never detecting eBay.
func TestDetectEBayFromRestClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"ebay/digital-signature-php-sdk": "^1.0"}}`},
		{"package.json": `{"dependencies": {"ebay-api": "^9.0.0"}}`},
		{"package.json": `{"dependencies": {"@hendt/ebay-api": "^9.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "ebay") {
			t.Errorf("%v did not detect ebay: %+v", manifest, p.Requirements)
		}
	}
}

// The companion, and the more important half here given the install counts.
func TestDetectDoesNotOfferEBayForTheXMLTradingSDKs(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"dts/ebay-sdk-php": "^19.0"}}`},
		{"composer.json": `{"require": {"benmorel/ebay-sdk-php": "^19.0"}}`},
		{"composer.json": `{"require": {"sapientpro/ebay-traditional-sdk": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "ebay") {
			t.Errorf("%v offered the ebay Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

func TestDetectReturnsErrNoProjectForEmptyDirectory(t *testing.T) {
	root := t.TempDir()

	_, err := Detect(root)

	if !errors.Is(err, ErrNoProject) {
		t.Fatalf("expected ErrNoProject, got %v", err)
	}
}

func TestDetectLaravelProject(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{
			"require": {
				"php": "^8.5",
				"laravel/framework": "^13.0",
				"laravel/horizon": "^6.0",
				"predis/predis": "^2.2",
				"stripe/stripe-php": "^17.0",
				"shopify/shopify-api": "^5.4"
			}
		}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Framework != "Laravel" {
		t.Errorf("Framework = %q, want Laravel", p.Framework)
	}

	if got := requirement(t, p, KindRuntime, "php").Version; got != "8.5" {
		t.Errorf("php version = %q, want 8.5", got)
	}

	for _, name := range []string{"redis", "horizon"} {
		if !p.Has(KindService, name) {
			t.Errorf("expected service %q to be detected", name)
		}
	}

	for _, name := range []string{"stripe", "shopify"} {
		if !p.Has(KindRecipe, name) {
			t.Errorf("expected recipe %q to be detected", name)
		}
	}
}

func TestDetectRecordsEvidenceAndSource(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{"require": {"stripe/stripe-php": "^17.0"}}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	stripe := requirement(t, p, KindRecipe, "stripe")

	if stripe.Evidence != "stripe/stripe-php" {
		t.Errorf("Evidence = %q, want stripe/stripe-php", stripe.Evidence)
	}

	if stripe.Source != "composer.json" {
		t.Errorf("Source = %q, want composer.json", stripe.Source)
	}
}

func TestDetectNodeProject(t *testing.T) {
	root := writeProject(t, map[string]string{
		"package.json": `{
			"engines": {"node": ">=24.0.0"},
			"dependencies": {
				"stripe": "^18.0.0",
				"@shopify/shopify-api": "^11.0.0",
				"ioredis": "^5.4.1"
			},
			"devDependencies": {"vite": "^7.0.0"}
		}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if got := requirement(t, p, KindRuntime, "node").Version; got != "24.0.0" {
		t.Errorf("node version = %q, want 24.0.0", got)
	}

	for _, name := range []string{"stripe", "shopify"} {
		if !p.Has(KindRecipe, name) {
			t.Errorf("expected recipe %q", name)
		}
	}

	if !p.Has(KindService, "redis") {
		t.Error("expected redis service from ioredis")
	}

	if !p.Has(KindService, "vite") {
		t.Error("expected vite")
	}
}

func TestDetectGoProjectIgnoresIndirectDependencies(t *testing.T) {
	root := writeProject(t, map[string]string{
		"go.mod": `module example.com/app

go 1.26

require (
	github.com/stripe/stripe-go/v76 v76.25.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/slack-go/slack v0.15.0 // indirect
)
`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if got := requirement(t, p, KindRuntime, "go").Version; got != "1.26" {
		t.Errorf("go version = %q, want 1.26", got)
	}

	if !p.Has(KindRecipe, "stripe") {
		t.Error("expected stripe recipe from versioned module path")
	}

	if !p.Has(KindService, "redis") {
		t.Error("expected redis service")
	}

	if p.Has(KindRecipe, "slack") {
		t.Error("indirect dependencies must not produce recipes")
	}
}

func TestDetectReportsUnmatchedAPIClients(t *testing.T) {
	root := writeProject(t, map[string]string{
		"composer.json": `{"require": {"acme/weather-api-client": "^1.0"}}`,
	})

	p, err := Detect(root)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "acme/weather-api-client") {
		t.Error("an unknown package must never become a booted recipe")
	}

	if len(p.Unmatched) != 1 {
		t.Fatalf("Unmatched = %+v, want exactly one entry", p.Unmatched)
	}

	if p.Unmatched[0].Evidence != "acme/weather-api-client" {
		t.Errorf("Unmatched evidence = %q", p.Unmatched[0].Evidence)
	}
}

func TestDetectIsDeterministic(t *testing.T) {
	files := map[string]string{
		"composer.json": `{"require": {"php": "^8.5", "twilio/sdk": "^8.0", "stripe/stripe-php": "^17.0"}}`,
		"package.json":  `{"dependencies": {"@octokit/rest": "^21.0.0"}}`,
	}

	var first []Requirement

	for i := 0; i < 5; i++ {
		p, err := Detect(writeProject(t, files))
		if err != nil {
			t.Fatalf("Detect: %v", err)
		}

		if first == nil {
			first = p.Requirements
			continue
		}

		if len(p.Requirements) != len(first) {
			t.Fatalf("run %d produced %d requirements, first run produced %d", i, len(p.Requirements), len(first))
		}

		for j := range first {
			if p.Requirements[j].Kind != first[j].Kind || p.Requirements[j].Name != first[j].Name {
				t.Fatalf("run %d differs at index %d: %+v vs %+v", i, j, p.Requirements[j], first[j])
			}
		}
	}
}

func TestNormaliseVersion(t *testing.T) {
	cases := map[string]string{
		"^8.5":       "8.5",
		"~24.0.1":    "24.0.1",
		">=8.2 <9.0": "8.2",
		"*":          "",
		"":           "",
		"v76.25.0":   "76.25.0",
		">=24.0.0":   "24.0.0",
	}

	for input, want := range cases {
		if got := normaliseVersion(input); got != want {
			t.Errorf("normaliseVersion(%q) = %q, want %q", input, got, want)
		}
	}
}
