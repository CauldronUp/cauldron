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

// Shopify ships two Recipes and one of them cannot be detected from most of
// its own libraries, which is worth asserting rather than leaving to be
// rediscovered.
//
// A dependency maps to exactly one Recipe: the lookup returns the first match
// and stops. Shopify's two most-installed clients each serve both the REST
// Admin API and the GraphQL one -- @shopify/shopify-api exposes a REST client
// and a GraphQL client side by side, and @shopify/admin-api-client documents
// both -- so neither of them answers which emulator a project wants. They stay
// on the REST Recipe, where they already were, because moving them would
// silently change what existing projects are offered on a guess.
//
// @shopify/graphql-client is GraphQL and nothing else, so it is the one
// dependency that decides by itself.
func TestDetectShopifyGraphQLOnlyFromTheGraphQLOnlyClient(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"@shopify/graphql-client": "^1.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "shopifygraphql") {
		t.Errorf("@shopify/graphql-client did not detect shopifygraphql: %+v", p.Requirements)
	}

	// The ambiguous ones still answer REST, and the point is that this is a
	// deliberate default rather than an oversight.
	for _, pkg := range []string{"@shopify/shopify-api", "@shopify/admin-api-client"} {
		p, err := Detect(writeProject(t, map[string]string{
			"package.json": `{"dependencies": {"` + pkg + `": "^11.0.0"}}`,
		}))
		if err != nil {
			t.Fatalf("Detect(%s): %v", pkg, err)
		}

		if !p.Has(KindRecipe, "shopify") {
			t.Errorf("%s should still answer the REST Recipe: %+v", pkg, p.Requirements)
		}

		if p.Has(KindRecipe, "shopifygraphql") {
			t.Errorf("%s serves both APIs and should not claim to pick one: %+v", pkg, p.Requirements)
		}
	}
}

// The SP-API replaced Marketplace Web Service in 2022, and MWS was not an
// earlier version of it: XML over a signed query string, its own endpoints,
// its own vocabulary. Packages named after MWS are still installed and still
// have amazon in the name, so a rule that matched the word would offer a
// REST emulator to code that speaks none of it.
func TestDetectAmazonSPFromSellingPartnerClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"jlevers/selling-partner-api": "^6.0"}}`},
		{"composer.json": `{"require": {"amazon-php/sp-api-sdk": "^7.0"}}`},
		{"composer.json": `{"require": {"amzn-spapi/sdk": "^1.0"}}`},
		{"package.json": `{"dependencies": {"amazon-sp-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "amazonsp") {
			t.Errorf("%v did not detect amazonsp: %+v", manifest, p.Requirements)
		}
	}
}

func TestDetectDoesNotOfferAmazonSPForMWSOrTheAWSSDK(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"amazon-mws": "^0.3.0"}}`},
		{"package.json": `{"dependencies": {"aws-sdk": "^2.1000.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "amazonsp") {
			t.Errorf("%v offered the amazonsp Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// The unscoped npm package called easypost is not this EasyPost. It reads
// POST bodies from form submissions in Node and has nothing to do with
// shipping, which is exactly the collision a substring rule would fall into:
// a project parsing form data would be offered a carrier emulator.
func TestDetectEasyPostFromTheCarrierClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"easypost/easypost-php": "^7.0"}}`},
		{"package.json": `{"dependencies": {"@easypost/api": "^7.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "easypost") {
			t.Errorf("%v did not detect easypost: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"easypost": "^1.0.0"}}`},
		{"composer.json": `{"require": {"gebruederheitz/wp-easy-post-options": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "easypost") {
			t.Errorf("%v offered the easypost Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// aftership/ on Packagist holds both API clients and Magento storefront
// extensions. A shop installing the extension talks to Magento; AfterShip
// talks to AfterShip. A vendor-prefix rule would offer a tracking emulator to
// a shop that never calls one.
func TestDetectAfterShipFromTheSDKsAndNotTheStorefrontExtension(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"aftership/aftership-php-sdk": "^6.0"}}`},
		{"composer.json": `{"require": {"aftership/tracking-sdk": "^1.0"}}`},
		{"package.json": `{"dependencies": {"aftership": "^7.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "aftership") {
			t.Errorf("%v did not detect aftership: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"aftership/aftership-apps-magento2": "^1.0"}}`},
		{"composer.json": `{"require": {"mrmonsters/magento2-aftership": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "aftership") {
			t.Errorf("%v offered the aftership Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// Three unrelated vendors publish a package called ship-station on Packagist
// and all three are real ShipStation API wrappers, so the vendor half of the
// name carries no information and each one has to be listed. There is no
// unscoped npm package called shipstation at all.
func TestDetectShipStationFromEveryVendorsWrapper(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"campo/laravel-shipstation": "^2.0"}}`},
		{"composer.json": `{"require": {"dissolve/ship-station": "^1.0"}}`},
		{"composer.json": `{"require": {"zack6849/ship-station": "^1.0"}}`},
		{"package.json": `{"dependencies": {"node-shipstation": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "shipstation") {
			t.Errorf("%v did not detect shipstation: %+v", manifest, p.Requirements)
		}
	}
}

// The same shape as ShipStation, from the same company: several vendors
// publish a package called shipengine and all of them are API clients, so the
// vendor half of the name says who wrote the wrapper rather than what it
// talks to. ShipEngine's own is one entry among several rather than the only
// one worth having.
func TestDetectShipEngineFromEveryVendorsWrapper(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"shipengine/shipengine": "^1.0"}}`},
		{"composer.json": `{"require": {"always-open/shipengine": "^1.0"}}`},
		{"composer.json": `{"require": {"jsamhall/shipengine": "^1.0"}}`},
		{"package.json": `{"dependencies": {"shipengine": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "shipengine") {
			t.Errorf("%v did not detect shipengine: %+v", manifest, p.Requirements)
		}
	}

	// And it is not the ShipStation Recipe, which is a different API from the
	// same company with a different envelope, a different credential and
	// dates in a different shape.
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"shipengine": "^1.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "shipstation") {
		t.Errorf("shipengine offered the shipstation Recipe: %+v", p.Requirements)
	}
}

// DHL's own SDK is the wrong one for this Recipe, which is the case worth
// asserting. dhl/sdk-api-express requires ext-soap -- checked in its manifest
// rather than guessed from its name -- and this Recipe models the REST
// tracking API, so a SOAP client pointed at it gets a 404 on its endpoint.
//
// dhlparcel-php-api is excluded for a different reason again: DHL Parcel is a
// separate API from DHL Express with its own credentials and shapes, and this
// Recipe says in as many words that it models Express and nothing else.
// Detecting it would offer an emulator that serves none of its endpoints.
func TestDetectDHLFromTheRestClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"sonnenglas/mydhl-php-sdk": "^1.0"}}`},
		{"composer.json": `{"require": {"korbeil/dhl-express-php-api": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "dhl") {
			t.Errorf("%v did not detect dhl: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"dhl/sdk-api-express": "^1.0"}}`},
		{"composer.json": `{"require": {"dreipunktnull/dhl-express": "^1.0"}}`},
		{"composer.json": `{"require": {"alfallouji/dhl_api": "^1.0"}}`},
		{"composer.json": `{"require": {"mvdnbrk/dhlparcel-php-api": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "dhl") {
			t.Errorf("%v offered the dhl Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// "recharge" is one of the most overloaded words on either registry, and this
// is the case where a substring rule would be wrong in four different ways at
// once. The unscoped npm package is a React static site generator. Packagist's
// matches are an EV charging SDK, an electricity meter service and a
// Zimbabwean airtime top-up library. Subscription billing is what the word
// means fourth or fifth.
//
// The storefront client is excluded for the reason DHL Parcel is: it really is
// Recharge, and it is a different API. It talks to the customer portal with a
// session token, and this Recipe models the admin API behind a store key.
func TestDetectRechargeFromTheSubscriptionClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"huel-global/recharge": "^1.0"}}`},
		{"package.json": `{"dependencies": {"recharge-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "recharge") {
			t.Errorf("%v did not detect recharge: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"recharge": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@rechargeapps/storefront-client": "^1.0.0"}}`},
		{"composer.json": `{"require": {"shell/ev-recharge-sdk": "^1.0"}}`},
		{"composer.json": `{"require": {"tinosoft/hot-recharge": "^1.0"}}`},
		{"composer.json": `{"require": {"recharge-meter/recharge-meter-service": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "recharge") {
			t.Errorf("%v offered the recharge Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// The most-installed package under the lemonsqueezy vendor name is a UI
// component library, not a client: plain-ui-components has roughly twice the
// installs of the Laravel package that actually talks to the API. It is the
// same shape as etsy/phan -- a vendor publishing more than its SDK -- and a
// prefix rule would offer a commerce emulator to a project that only borrowed
// some Blade components.
func TestDetectLemonSqueezyFromTheClientsAndNotTheComponentLibrary(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"lemonsqueezy/laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"seisigmasrl/lemonsqueezy-php": "^1.0"}}`},
		{"package.json": `{"dependencies": {"@lemonsqueezy/lemonsqueezy.js": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "lemonsqueezy") {
			t.Errorf("%v did not detect lemonsqueezy: %+v", manifest, p.Requirements)
		}
	}

	p, err := Detect(writeProject(t, map[string]string{
		"composer.json": `{"require": {"lemonsqueezy/plain-ui-components": "^1.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "lemonsqueezy") {
		t.Errorf("a UI component library offered the lemonsqueezy Recipe: %+v", p.Requirements)
	}
}

// Nothing maps to the Toast Recipe, and that is the decision rather than an
// omission.
//
// "toast" means a notification popup on both registries. tall-toasts alone
// has over four hundred thousand installs, toastable and toastbox are the
// same idea, and the unscoped npm package is a different library again. There
// is no widely used client for the Toast POS API at all -- integrations are
// written straight against the REST endpoints -- so there is nothing to map
// that would not be a collision.
//
// This guards the decision rather than a behaviour: a naive mapping on the
// word would offer a restaurant point-of-sale emulator to every Laravel
// project that shows a toast when a form saves.
func TestDetectDoesNotOfferToastForToastNotificationLibraries(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"usernotnull/tall-toasts": "^2.0"}}`},
		{"composer.json": `{"require": {"acatech/toastable": "^1.0"}}`},
		{"composer.json": `{"require": {"awesometoast/toastbox": "^1.0"}}`},
		{"package.json": `{"dependencies": {"toast": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "toast") {
			t.Errorf("%v offered the toast Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// One vendor, two protocols, and the names do not say which is which.
// avalara/avatax describes itself as the AvaTax SOAP client and
// avalara/avataxclient is the REST v2 one -- checked in their manifests,
// where the first has no HTTP dependency at all and the second pulls in
// Guzzle. This Recipe models REST v2, so a SOAP client pointed at it gets a
// 404 on its endpoint.
//
// The other two exclusions are the shapes already seen elsewhere:
// avatax-magento is a storefront module rather than a client, like
// aftership-apps-magento2, and avalara-for-communications is a different
// Avalara product with its own API, like DHL Parcel.
func TestDetectAvalaraFromTheRestClientOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"avalara/avataxclient": "^24.0"}}`},
		{"package.json": `{"dependencies": {"avatax": "^24.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "avalara") {
			t.Errorf("%v did not detect avalara: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"avalara/avatax": "^1.0"}}`},
		{"composer.json": `{"require": {"avalara/avatax-magento": "^1.0"}}`},
		{"composer.json": `{"require": {"phoneburner/avalara-for-communications-php-sdk": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "avalara") {
			t.Errorf("%v offered the avalara Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// The fourth time in this table that the most-installed package for a
// provider is the wrong protocol -- after eBay's Trading SDKs, DHL's own SOAP
// client and avalara/avatax. php-fedex-api-wrapper has over a million
// installs and wraps the SOAP web services this REST API replaced, while the
// REST wrappers have tens of thousands between them.
//
// That is the pattern worth naming: when a provider moves from SOAP or XML to
// REST, the old client keeps its install base for years, so popularity is
// evidence of age rather than of correctness.
func TestDetectFedExFromTheRestWrappersOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"whatarmy/fedex-rest": "^1.0"}}`},
		{"composer.json": `{"require": {"shipstream/fedex-rest-sdk": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "fedex") {
			t.Errorf("%v did not detect fedex: %+v", manifest, p.Requirements)
		}
	}

	p, err := Detect(writeProject(t, map[string]string{
		"composer.json": `{"require": {"jeremy-dunn/php-fedex-api-wrapper": "^1.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "fedex") {
		t.Errorf("the SOAP wrapper offered the fedex Recipe: %+v", p.Requirements)
	}
}
