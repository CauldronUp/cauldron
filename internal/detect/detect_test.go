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

// The rare provider whose name carries the information its packages do.
// Nothing else on either registry is called easyship, so there is no
// collision to avoid and no protocol split to pick a side of -- unlike
// easypost, recharge, toast and fedex, which needed all three. The clients
// are small, which is not a reason to leave them out: the bar is whether a
// dependency identifies the provider, not how many people use it.
func TestDetectEasyshipFromItsClients(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"jrebs/easyship-php": "^1.0"}}`},
		{"composer.json": `{"require": {"jrebs/easyship-laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"tigusigalpa/easyship-php": "^1.0"}}`},
		{"package.json": `{"dependencies": {"easyship": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "easyship") {
			t.Errorf("%v did not detect easyship: %+v", manifest, p.Requirements)
		}
	}

	// And it is not easypost, whose name is one letter away and whose Recipe
	// models a different carrier aggregator entirely.
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"easyship": "^1.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "easypost") {
		t.Errorf("easyship offered the easypost Recipe: %+v", p.Requirements)
	}
}

// The worst collision in this table, and the reason there is no PHP mapping
// for Clover at all: on Packagist the word is a code-coverage report format.
// It is what PHPUnit emits, so nette/tester, phpunit-coverage-check and
// several others carry it and have tens of millions of installs between
// them. Mapping any of them would offer a card-payments emulator to every
// project that measures its own test coverage.
//
// The unscoped npm package is a toy language, and remote-pay-cloud is
// Clover's own device SDK -- it drives a card terminal over a socket and
// never calls the REST API this Recipe models, which is the DHL Parcel
// exclusion again.
func TestDetectCloverFromTheEcommerceSdkOnly(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"clover-ecomm-sdk": "^1.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "clover") {
		t.Errorf("clover-ecomm-sdk did not detect clover: %+v", p.Requirements)
	}

	for _, manifest := range []map[string]string{
		{"composer.json": `{"require-dev": {"nette/tester": "^2.0"}}`},
		{"composer.json": `{"require": {"rregeer/phpunit-coverage-check": "^0.3"}}`},
		{"composer.json": `{"require": {"jaschilz/php-coverage-badger": "^1.0"}}`},
		{"package.json": `{"dependencies": {"clover": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"remote-pay-cloud": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "clover") {
			t.Errorf("%v offered the clover Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// Lightspeed is four unrelated products sharing a brand, not the two the
// backlog said: R-Series is retail, K-Series is restaurants, X-Series is the
// former Vend and eCom is the former SEOshop. Four APIs, four credentials,
// four vocabularies, and a client for one does not partly work against
// another -- it does not reach it.
//
// The two nearest misses are not Lightspeed at all. The unscoped npm package
// is the speed of light as a constant, and Packagist's most-installed match
// is a Magento performance module that borrowed the word.
func TestDetectLightspeedFromTheRetailClientOnly(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"composer.json": `{"require": {"timothydc/laravel-lightspeed-retail-api": "^1.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "lightspeed") {
		t.Errorf("the R-Series client did not detect lightspeed: %+v", p.Requirements)
	}

	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"timothydc/laravel-lightspeed-ecom-api": "^1.0"}}`},
		{"composer.json": `{"require": {"seoshop/seoshop-php": "^1.0"}}`},
		{"composer.json": `{"require": {"anytech/lightspeed-x-series-api": "^1.0"}}`},
		{"composer.json": `{"require": {"elgentos/module-lightspeed": "^1.0"}}`},
		{"package.json": `{"dependencies": {"lightspeed": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "lightspeed") {
			t.Errorf("%v offered the lightspeed Recipe and should not have: %+v", manifest, p.Requirements)
		}
	}
}

// printful/php-gettext-cms is Printful's own and is a translation management
// backend rather than an API client -- the same shape as etsy/phan and
// lemonsqueezy/plain-ui-components, and with a six-figure install count. A
// vendor prefix would offer a print-on-demand emulator to a project managing
// its own translations.
func TestDetectPrintfulFromTheApiClientsOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"printful/php-api-sdk": "^1.0"}}`},
		{"composer.json": `{"require": {"fwartner/printful": "^1.0"}}`},
		{"package.json": `{"dependencies": {"printful": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"printful-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "printful") {
			t.Errorf("%v did not detect printful: %+v", manifest, p.Requirements)
		}
	}

	p, err := Detect(writeProject(t, map[string]string{
		"composer.json": `{"require": {"printful/php-gettext-cms": "^1.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "printful") {
		t.Errorf("a translation CMS offered the printful Recipe: %+v", p.Requirements)
	}
}

// A different reason to exclude something than any other row in this table.
// @medusajs/medusa is not a client of the API -- it is the commerce engine
// that serves it. A project depending on it is the store rather than a caller
// of one, and offering it an emulator of itself is not useful: that project's
// tests want a database, not a fake of the thing they are running.
//
// @medusajs/js-sdk is the client, and it is the only mapping.
func TestDetectMedusaFromTheClientAndNotTheEngine(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"@medusajs/js-sdk": "^2.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "medusa") {
		t.Errorf("the SDK did not detect medusa: %+v", p.Requirements)
	}

	p, err = Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"@medusajs/medusa": "^2.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "medusa") {
		t.Errorf("the engine offered the medusa Recipe: %+v", p.Requirements)
	}
}

// The npm package called printify is a print-CSS helper -- "easy HTML print
// formatting" -- and has nothing to do with the print-on-demand company. The
// word means putting ink on paper before it means putting ink on a t-shirt.
//
// A test also asserts printify does not offer the printful Recipe. The two
// are the same category and their Recipes model different shapes: Printful
// puts cost and retail side by side on the order, Printify puts cost only on
// the line items.
func TestDetectPrintifyFromThePhpSdksOnly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"jacob-hyde/printify": "^1.0"}}`},
		{"composer.json": `{"require": {"garissman/printify": "^1.0"}}`},
		{"composer.json": `{"require": {"ahsan/printify-laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"8fold/php-printify-sdk": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "printify") {
			t.Errorf("%v did not detect printify: %+v", manifest, p.Requirements)
		}

		if p.Has(KindRecipe, "printful") {
			t.Errorf("%v offered the printful Recipe: %+v", manifest, p.Requirements)
		}
	}

	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"printify": "^1.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if p.Has(KindRecipe, "printify") {
		t.Errorf("a print-CSS helper offered the printify Recipe: %+v", p.Requirements)
	}
}

// Shopware's detection has to exclude the provider from its own name twice
// over.
//
// The four most-installed Packagist matches -- shopware/core,
// shopware/storefront, shopware/administration and shopware/platform -- are
// the shop, not a client of one. That is the exclusion Medusa already needed,
// at four to six million installs apiece rather than at Medusa's scale.
//
// shopware/shopware is the one worth a test of its own, because it is a shape
// nothing else here has: same vendor, same word, three quarters of a million
// installs, and a completely different product. Shopware 5's API is
// /api/articles behind basic auth, with no sales channel and no sw-access-key
// anywhere in it. This Recipe is 6.7's Store API and describes nothing that
// shop does.
//
// @shopware-ag/meteor-admin-sdk is excluded because it is for extensions
// running inside the administration interface, which is a Vue application
// rather than an HTTP API.
func TestDetectShopwareFromItsClientsAndNotFromTheShopItself(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"vin-sw/shopware-sdk": "^1.0"}}`},
		{"package.json": `{"dependencies": {"@shopware/api-client": "^1.5.0"}}`},
		{"package.json": `{"dependencies": {"@shopware-pwa/api-client": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@shopware/api-gen": "^1.5.0"}}`},
		{"package.json": `{"dependencies": {"shopware-api-client": "^1.6.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "shopware") {
			t.Errorf("%v did not detect shopware: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The shop itself, four ways.
		{"composer.json": `{"require": {"shopware/core": "^6.7"}}`},
		{"composer.json": `{"require": {"shopware/storefront": "^6.7"}}`},
		{"composer.json": `{"require": {"shopware/administration": "^6.7"}}`},
		{"composer.json": `{"require": {"shopware/platform": "^6.7"}}`},
		// And the previous generation, which is a different API.
		{"composer.json": `{"require": {"shopware/shopware": "^5.7"}}`},
		// An administration extension, which never speaks to the Store API.
		{"package.json": `{"dependencies": {"@shopware-ag/meteor-admin-sdk": "^6.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "shopware") {
			t.Errorf("%v offered the shopware Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// commercetools ships two SDKs for two different APIs under one npm scope, so
// its exclusions have to be made inside the vendor's own namespace.
//
// @commercetools/importapi-sdk is the Import API: a separate service on a
// separate host, built around import containers and asynchronous import
// operations, with no version on a resource and no update actions. A project
// holding it is not calling the API this Recipe describes, and the only thing
// that says so is the middle of the package name.
//
// The Merchant Center toolkits are excluded for the ordinary reason: they are
// React components for building screens inside commercetools' own back office
// rather than callers of the HTTP API.
func TestDetectCommercetoolsAndNotItsOtherApiOrItsFrontend(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"commercetools/commercetools-sdk": "^7.0"}}`},
		{"composer.json": `{"require": {"commercetools/php-sdk": "^2.0"}}`},
		{"composer.json": `{"require": {"commercetools/symfony-bundle": "^3.0"}}`},
		{"composer.json": `{"require": {"bestit/commercetools-odm": "^3.0"}}`},
		{"package.json": `{"dependencies": {"@commercetools/platform-sdk": "^8.0.0"}}`},
		{"package.json": `{"dependencies": {"@commercetools/sdk-client-v2": "^3.0.0"}}`},
		{"package.json": `{"dependencies": {"@commercetools/ts-client": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "commercetools") {
			t.Errorf("%v did not detect commercetools: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// A different API from the same vendor, one scope apart.
		{"package.json": `{"dependencies": {"@commercetools/importapi-sdk": "^5.0.0"}}`},
		// And the back office's own component libraries.
		{"package.json": `{"dependencies": {"@commercetools-frontend/application-shell": "^23.0.0"}}`},
		{"package.json": `{"dependencies": {"@commercetools-uikit/text": "^19.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "commercetools") {
			t.Errorf("%v offered the commercetools Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// VTEX has no widely installed external client, and the packages that look
// like one are the platform's own toolchain.
//
// npm "vtex" is the Toolbelt (repository vtex/toolbelt): a command-line tool
// for building and publishing IO apps, which calls nothing on a shop's
// behalf. @vtex/api (vtex/node-vtex-api) and @vtex/clients (vtex/io-clients)
// are the IO runtime's clients, and an app using them runs inside VTEX's
// infrastructure with a platform token rather than sending
// X-VTEX-API-AppKey and X-VTEX-API-AppToken -- which is the entire credential
// this Recipe models.
//
// So the mapped packages are the small external ones. That is the honest
// result rather than a shortfall: a four-figure package making these exact
// requests is a better signal than a five-figure one making different ones.
func TestDetectVtexFromExternalClientsAndNotFromTheIoToolchain(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"brandlive/vtex-api": "^1.0"}}`},
		{"composer.json": `{"require": {"juanfeltrin/vtex-sdk-php": "^1.0"}}`},
		{"composer.json": `{"require": {"daygarcia/laravel-vtex": "^1.0"}}`},
		{"composer.json": `{"require": {"grebo87/laravel-vtex-api": "^1.0"}}`},
		{"package.json": `{"dependencies": {"vtex-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "vtex") {
			t.Errorf("%v did not detect vtex: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The Toolbelt, which is a build tool.
		{"package.json": `{"devDependencies": {"vtex": "^3.0.0"}}`},
		// And the IO runtime clients, which authenticate a different way.
		{"package.json": `{"dependencies": {"@vtex/api": "^7.4.0"}}`},
		{"package.json": `{"dependencies": {"@vtex/clients": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "vtex") {
			t.Errorf("%v offered the vtex Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Saleor publishes four packages under one scope and only two of them make
// requests to the API this Recipe describes.
//
// @saleor/sdk is the client and @saleor/auth-sdk is the half that gets a
// token. npm "saleor" is the CLI, which creates and deploys projects.
// @saleor/macaw-ui is a React component library.
//
// @saleor/app-sdk is the exclusion worth stating: it is for building Saleor
// Apps, which is code Saleor calls rather than code that calls Saleor. An app
// does query the API eventually, with a token handed to it at install time
// and through whatever client it picks, so the dependency says "this is an
// extension" rather than "this makes these requests".
func TestDetectSaleorFromItsClientAndNotItsCliOrAppFramework(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@saleor/sdk": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@saleor/auth-sdk": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "saleor") {
			t.Errorf("%v did not detect saleor: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		{"package.json": `{"devDependencies": {"saleor": "^3.0.0"}}`},
		{"package.json": `{"dependencies": {"@saleor/app-sdk": "^0.50.0"}}`},
		{"package.json": `{"dependencies": {"@saleor/macaw-ui": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "saleor") {
			t.Errorf("%v offered the saleor Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Allegro's two exclusions are shapes this collection has recorded before,
// arriving together and with larger numbers than usual.
//
// allegro/php-protobuf is the most-installed match at seventy-six thousand,
// is published from github.com/allegro/php-protobuf by the same company, and
// is Google's Protocol Buffers for PHP. The vendor prefix is genuine and the
// package has nothing to do with the marketplace.
//
// The npm package called allegro describes itself as an "Allegro.pl WebAPI
// client", lists soap among its keywords and depends on soap-js. It is a
// client for the legacy SOAP interface rather than the REST API this Recipe
// models -- the same shape as eBay's Trading SDKs and DHL's own SOAP client.
func TestDetectAllegroFromRestClientsAndNotProtobufOrSoap(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"imper86/php-allegro-api": "^1.0"}}`},
		{"composer.json": `{"require": {"wiatrogon/allegro-rest-api": "^1.0"}}`},
		{"composer.json": `{"require": {"asocial-media/allegro-api": "^1.0"}}`},
		{"composer.json": `{"require": {"sebastianpozoga/php-allegroapi": "^1.0"}}`},
		{"composer.json": `{"require": {"ircykk/allegro-api": "^1.0"}}`},
		{"composer.json": `{"require": {"zoondo/allegro-api": "^1.0"}}`},
		{"composer.json": `{"require": {"macopedia/magento2-allegro": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "allegro") {
			t.Errorf("%v did not detect allegro: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's own unrelated library, and the top match by installs.
		{"composer.json": `{"require": {"allegro/php-protobuf": "^0.13"}}`},
		// And a client for the protocol this API replaced.
		{"package.json": `{"dependencies": {"allegro": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "allegro") {
			t.Errorf("%v offered the allegro Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Akeneo is where the vendor-prefix rule fails hardest.
//
// Four packages under the akeneo namespace have seven-figure install counts
// and only one of them is a client. akeneo/php-coupling-detector is an
// architecture linter at 1.1 million; akeneo/batch-bundle is part of the
// PIM's own Symfony stack at 1.06 million;
// akeneo/phpspec-skip-example-extension is a test-framework extension at 997
// thousand. None of the three has anything to do with the API.
//
// akeneo/pim-community-dev is excluded for the other reason: it is the PIM,
// so a project holding it serves this API rather than calls it.
//
// Three vendors now -- Etsy, Allegro, Akeneo -- have had their own prefix
// carry packages about something else entirely, which is enough to treat it
// as a shape rather than a run of coincidences.
func TestDetectAkeneoFromClientsAndNotTheVendorsOtherLibraries(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"akeneo/api-php-client": "^7.0"}}`},
		{"composer.json": `{"require": {"akeneo/module-magento2-connector-community": "^100.0"}}`},
		{"composer.json": `{"require": {"webgriffe/sylius-akeneo-plugin": "^2.0"}}`},
		{"composer.json": `{"require": {"synolia/sylius-akeneo-plugin": "^4.0"}}`},
		{"composer.json": `{"require": {"justbetter/laravel-akeneo-client": "^1.0"}}`},
		{"composer.json": `{"require": {"justbetter/magento2-akeneo-bundle": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "akeneo") {
			t.Errorf("%v did not detect akeneo: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's own libraries, none of them about the API.
		{"composer.json": `{"require-dev": {"akeneo/php-coupling-detector": "^0.5"}}`},
		{"composer.json": `{"require-dev": {"akeneo/phpspec-skip-example-extension": "^5.0"}}`},
		{"composer.json": `{"require": {"akeneo/batch-bundle": "^7.0"}}`},
		{"composer.json": `{"require": {"akeneo-labs/spreadsheet-parser": "^1.0"}}`},
		// And the PIM itself.
		{"composer.json": `{"require": {"akeneo/pim-community-dev": "^7.0"}}`},
		{"composer.json": `{"require": {"akeneo/pim-community-standard": "^7.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "akeneo") {
			t.Errorf("%v offered the akeneo Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Voucherify's near miss is a different company rather than a different word.
//
// voucherly/voucherly-php-sdk belongs to Voucherly, an Italian payments
// business whose name differs from Voucherify by two letters, so a prefix or
// substring rule matches it and a reader skims past it.
// fastwebmedia/laravel-vouchering is the ordinary kind, matching on the word.
//
// rspective/voucherify is mapped although the prefix is not the product name:
// rspective is Voucherify's parent company and the package describes itself
// as the "Voucherify promotion engine REST API".
func TestDetectVoucherifyAndNotVoucherly(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"rspective/voucherify": "^3.0"}}`},
		{"package.json": `{"dependencies": {"@voucherify/sdk": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"voucherify": "^4.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "voucherify") {
			t.Errorf("%v did not detect voucherify: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// A different company, two letters away.
		{"composer.json": `{"require": {"voucherly/voucherly-php-sdk": "^1.0"}}`},
		// And a different thing entirely, matching on the word.
		{"composer.json": `{"require": {"fastwebmedia/laravel-vouchering": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "voucherify") {
			t.Errorf("%v offered the voucherify Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Apideck's near miss is one letter, three times over.
//
// appdeck publishes sql, sampa and filter -- a PDO wrapper, an MVC framework
// and an input filter -- and none of them has anything to do with a unified
// commerce API. Voucherify had Voucherly two letters away; this has appdeck
// at one. muammarsiddiqui/apideck-laravel is the ordinary sort: an API
// playground that took the name.
//
// The mapping is also deliberately imprecise, and the comment beside it says
// so: one SDK covers every unified API Apideck offers, so the dependency
// proves the caller uses Apideck rather than that they use this one.
func TestDetectApideckAndNotAppdeck(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"apideck-libraries/php-sdk": "^0.10"}}`},
		{"composer.json": `{"require": {"apideck-libraries/sdk-php": "^1.0"}}`},
		{"package.json": `{"dependencies": {"@apideck/unify": "^0.20.0"}}`},
		{"package.json": `{"dependencies": {"@apideck/node": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"apideck": "^2.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "apideck") {
			t.Errorf("%v did not detect apideck: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// One letter away, and three different libraries.
		{"composer.json": `{"require": {"appdeck/sql": "^1.0"}}`},
		{"composer.json": `{"require": {"appdeck/sampa": "^1.0"}}`},
		{"composer.json": `{"require": {"appdeck/filter": "^1.0"}}`},
		// And a Laravel API playground that took the name.
		{"composer.json": `{"require": {"muammarsiddiqui/apideck-laravel": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "apideck") {
			t.Errorf("%v offered the apideck Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Royal Mail is mapped to nothing on purpose, and this guards the decision the
// way the Toast test guards its own.
//
// The company runs three separate APIs on three hosts, and every client
// package on Packagist targets one of the other two.
// elliotjreed/royal-mail-tracking is the tracking service;
// turtledesign/royalmail-php, zvps/royal-mail-shipping-rest-api-client and
// mobi-market/royalmail-shipping-v3 are all the Shipping API. Click & Drop --
// the API this Recipe describes, at api.parcel.royalmail.com -- has no client
// of its own.
//
// Offering its emulator to a project using a different Royal Mail API would be
// wrong about the host, the credential and the vocabulary at once, which is a
// worse answer than offering nothing.
func TestRoyalMailIsMappedToNothingOnPurpose(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"elliotjreed/royal-mail-tracking": "^4.0"}}`},
		{"composer.json": `{"require": {"turtledesign/royalmail-php": "^1.0"}}`},
		{"composer.json": `{"require": {"zvps/royal-mail-shipping-rest-api-client": "^1.0"}}`},
		{"composer.json": `{"require": {"mobi-market/royalmail-shipping-v3": "^1.0"}}`},
		{"composer.json": `{"require": {"octolize/wp-royal-mail-shipping-method": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "royalmail") {
			t.Errorf("%v offered the royalmail Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Lago's near miss is a substring of a bigger project's name, which is a shape
// none of the others here has.
//
// Lagoon is amazee.io's open-source Kubernetes application delivery platform,
// and "lago" sits inside it. A substring rule matches uselagoon/lagoon-php-sdk,
// steveworley/lagoon-php-sdk, amazeeio/lagoon-logs,
// salsadigitalauorg/wp_lagoon_logs and bryangruneberg/laragoon -- and offering
// a usage-billing emulator to every Drupal site deployed on amazee.io is not a
// near miss, it is a different industry.
//
// The npm package called lago is the ordinary sort: a data structures and
// algorithms library.
func TestDetectLagoAndNotLagoon(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"g-portal/lago-php-client": "^1.0"}}`},
		{"package.json": `{"dependencies": {"lago-javascript-client": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "lago") {
			t.Errorf("%v did not detect lago: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// A different product whose name contains this one.
		{"composer.json": `{"require": {"uselagoon/lagoon-php-sdk": "^1.0"}}`},
		{"composer.json": `{"require": {"steveworley/lagoon-php-sdk": "^1.0"}}`},
		{"composer.json": `{"require": {"amazeeio/lagoon-logs": "^2.0"}}`},
		{"composer.json": `{"require": {"salsadigitalauorg/wp_lagoon_logs": "^1.0"}}`},
		{"composer.json": `{"require": {"bryangruneberg/laragoon": "^1.0"}}`},
		// And an algorithms library that took the word.
		{"package.json": `{"dependencies": {"lago": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "lago") {
			t.Errorf("%v offered the lago Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Polar meets prefix containment, which Lago met first with Lagoon.
//
// "polar" starts "polaris" and "polarising", and both out-install the real
// client on Packagist: eduard9969/blade-polaris-icons is Shopify's design
// system at eighty-nine thousand and polarising/bcrypt is a password hashing
// library at sixty-eight. The first is worth naming twice over -- a prefix
// rule would offer a payments emulator to a project using Shopify's icons,
// and Shopify has two Recipes of its own here.
//
// npm's polar is the ordinary sort: an express boilerplate.
func TestDetectPolarAndNotPolarisOrPolarising(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"polar-sh/sdk": "^0.4"}}`},
		{"composer.json": `{"require": {"danestves/laravel-polar": "^1.0"}}`},
		{"package.json": `{"dependencies": {"@polar-sh/sdk": "^0.32.0"}}`},
		{"package.json": `{"dependencies": {"@polar-sh/nextjs": "^0.4.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "polar") {
			t.Errorf("%v did not detect polar: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Shopify's design system, and the bigger match by installs.
		{"composer.json": `{"require": {"eduard9969/blade-polaris-icons": "^1.0"}}`},
		{"composer.json": `{"require": {"rwsite/moonshine-polaris-theme": "^1.0"}}`},
		// A password hashing library whose vendor starts with the word.
		{"composer.json": `{"require": {"polarising/bcrypt": "^1.0"}}`},
		// And an express boilerplate that took the name outright.
		{"package.json": `{"dependencies": {"polar": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "polar") {
			t.Errorf("%v offered the polar Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Gumroad's exclusion is "sign in with" rather than "talk to".
//
// socialiteproviders/gumroad is a Laravel Socialite provider and
// alofoxx/oauth2-gumroad is the same idea for the PHP League. A project
// holding either lets people log in with a Gumroad account and may never
// touch a product or a licence -- the shape Shopware's meteor-admin-sdk and
// Saleor's app-sdk have.
//
// And for once the bare word is the real thing: npm's gumroad describes
// itself as "API client for Gumroad.", where the same plain name turned out
// to be an algorithms library for Lago and an express boilerplate for Polar.
func TestDetectGumroadClientsAndNotItsLoginProviders(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"serenity_technologies/gumroad": "^1.0"}}`},
		{"composer.json": `{"require": {"serenity_technologies/cashier-gumroad": "^1.0"}}`},
		{"composer.json": `{"require": {"diskopete/laravel-gumroad": "^1.0"}}`},
		{"composer.json": `{"require": {"flowframe/laravel-gumroad": "^1.0"}}`},
		{"package.json": `{"dependencies": {"gumroad": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"gumroad-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "gumroad") {
			t.Errorf("%v did not detect gumroad: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Logging in with Gumroad, not calling it.
		{"composer.json": `{"require": {"socialiteproviders/gumroad": "^4.0"}}`},
		{"composer.json": `{"require": {"alofoxx/oauth2-gumroad": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "gumroad") {
			t.Errorf("%v offered the gumroad Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Ecwid's two exclusions are shapes already recorded here, arriving together
// on one small provider.
//
// mugnate/oauth2-ecwid is a login provider for the PHP League -- signing in
// with an Ecwid account rather than calling the store API, which is what
// socialiteproviders/gumroad was.
//
// And npm's bare "ecwid" describes itself as "ecwid" and lists no repository,
// which is the state the backlog records against Marqeta: a package of that
// name exists and cannot be verified, so it is left alone rather than guessed
// at.
func TestDetectEcwidClientsAndNotItsLoginProviderOrTheUnverifiableName(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"dspacelabs/ecwid": "^1.0"}}`},
		{"composer.json": `{"require": {"dspacelabs/ecwid-client": "^1.0"}}`},
		{"composer.json": `{"require": {"strappberry/ecwid-api-wrapper": "^1.0"}}`},
		{"package.json": `{"dependencies": {"ecwid-api": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "ecwid") {
			t.Errorf("%v did not detect ecwid: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Logging in with Ecwid, not calling it.
		{"composer.json": `{"require": {"mugnate/oauth2-ecwid": "^1.0"}}`},
		// And the name nobody can verify.
		{"package.json": `{"dependencies": {"ecwid": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "ecwid") {
			t.Errorf("%v offered the ecwid Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Squarespace publishes a great deal of npm and none of it calls the Commerce
// API.
//
// The @squarespace/ scope is template tooling for building Squarespace sites,
// and @squarespace/core describes itself as "The frontend JS API for
// Squarespace templates" -- a package with the word API in its description
// that is not this API. eslint-config-squarespace is the company's own lint
// rules.
//
// unikapps/laravel-socialite-squarespace signs people in rather than talking
// to the store. beloop/squarespace is a component of an unrelated LMS.
// @tryghost/mg-squarespace-xml and is-squarespace never make a request at all:
// one reads an export file to move a site off the platform and the other looks
// at a page and guesses what built it.
func TestDetectSquarespaceCommerceClientsAndNotItsTemplateTooling(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"squarespace-node-api": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"marktera-squarespace": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@corte-so/commerce-squarespace": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@florasync/squarespace-mcp": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "squarespace") {
			t.Errorf("%v did not detect squarespace: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's own template tooling, including the one that calls
		// itself an API.
		{"package.json": `{"dependencies": {"@squarespace/core": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@squarespace/template-engine": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@squarespace/toolbelt": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"eslint-config-squarespace": "^1.0.0"}}`},
		// Moving content between platforms rather than calling the API.
		{"package.json": `{"dependencies": {"@tryghost/mg-squarespace-xml": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"seed-squarespace": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"is-squarespace": "^1.0.0"}}`},
		// Logging in with Squarespace, not calling it.
		{"composer.json": `{"require": {"unikapps/laravel-socialite-squarespace": "^1.0"}}`},
		// And somebody else's LMS.
		{"composer.json": `{"require": {"beloop/squarespace": "^1.0"}}`},
		{"composer.json": `{"require": {"beloop/squarespace-bundle": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "squarespace") {
			t.Errorf("%v offered the squarespace Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// "wix" is three different things on two registries, which makes this the
// busiest exclusion list in the file.
//
// It is a homonym: WiX is the Windows Installer XML toolset, and
// electron-wix-msi and @electron-forge/maker-wix build MSI packages with it.
// No shape recorded here covered that -- Lago met Lagoon and Polar met
// Polaris, which are longer words containing a provider's name, where this is
// the same word meaning something else.
//
// It is a prefix: wixel/gump is a PHP input validator with over a million
// downloads, and wixiweb, wixnit and wixet are three more vendors whose names
// begin with it.
//
// And the vendor's own scope is mostly not this API: @wix/design-system and
// wix-style-react are a React component library, @wix-pilot/core is an
// LLM-driven test runner, @wix/wix-code-types is types for Velo.
func TestDetectWixCommerceClientsAndNotTheInstallerToolsetOrTheDesignSystem(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@wix/sdk": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@wix/ecom": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@wix/api-client": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@wix/headless-ecom": "^1.0.0"}}`},
		{"composer.json": `{"require": {"chkltlabs/wix-client": "^1.0"}}`},
		{"composer.json": `{"require": {"storessuite/wix": "^1.0"}}`},
		{"composer.json": `{"require": {"storessuite/wix-connect": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "wix") {
			t.Errorf("%v did not detect wix: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The Windows Installer XML toolset, which is a different product
		// with the same name.
		{"package.json": `{"devDependencies": {"electron-wix-msi": "^5.0.0"}}`},
		{"package.json": `{"devDependencies": {"@electron-forge/maker-wix": "^7.0.0"}}`},
		{"package.json": `{"devDependencies": {"@mongodb-js/electron-wix-msi": "^5.0.0"}}`},
		// The vendor's own front end, which never calls the REST API.
		{"package.json": `{"dependencies": {"@wix/design-system": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"wix-style-react": "^10.0.0"}}`},
		{"package.json": `{"dependencies": {"wix-ui-core": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@wix/wix-ui-icons-common": "^3.0.0"}}`},
		// The vendor's own test runner, which is the most installed package
		// with the word in its name.
		{"package.json": `{"devDependencies": {"@wix-pilot/core": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@wix-pilot/detox": "^1.0.0"}}`},
		// Types for the language sites are scripted in.
		{"package.json": `{"devDependencies": {"@wix/wix-code-types": "^1.0.0"}}`},
		// The App Market SDK, published by an individual.
		{"package.json": `{"dependencies": {"wix": "^1.0.0"}}`},
		// A vendor whose name merely starts with it.
		{"composer.json": `{"require": {"wixel/gump": "^2.0"}}`},
		{"composer.json": `{"require": {"wixiweb/wixiweb-laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"wixnit/core": "^1.0"}}`},
		// Logging in with Wix, not calling it.
		{"composer.json": `{"require": {"lezhnev74/laravel-socialite-wix": "^1.0"}}`},
		// The retired contacts API.
		{"composer.json": `{"require": {"epicformbuilder/wixhive-php-api": "^1.0"}}`},
		// And Wix's own Symfony plumbing.
		{"composer.json": `{"require": {"wix/framework-bundle": "^1.0"}}`},
		{"composer.json": `{"require": {"wix/framework-component": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "wix") {
			t.Errorf("%v offered the wix Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Metronome's name is an everyday object, which produces more collisions than
// any shape recorded here so far.
//
// @dougflip/metronome, expo-precision-metronome, strumming-metronome and
// react-native-metronome-module are all what the word means: things that click
// in time. Wix was a homonym for a build toolset; this is a homonym for a piece
// of hardware.
//
// metronome-contracts and metronome-wallet-core belong to a different company
// with the same name -- the Metronome ERC-20 token. The bare npm "metronome" is
// a Node application framework, which is the shape Lago and Polar both met.
//
// And the closest miss is diego-ninja/laravel-metronome, "a metric aggregator
// and processor for Laravel": a metric aggregator named metronome that is not
// the metering provider named Metronome.
func TestDetectMetronomeSDKAndNotTheInstrumentOrTheToken(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@metronome/sdk": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@metronome/mcp": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "metronome") {
			t.Errorf("%v did not detect metronome: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The instrument.
		{"package.json": `{"dependencies": {"@dougflip/metronome": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"expo-precision-metronome": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"strumming-metronome": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"react-native-metronome-module": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@pickupmusic/metronome": "^1.0.0"}}`},
		// A different company with the same name.
		{"package.json": `{"dependencies": {"metronome-contracts": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"metronome-wallet-core": "^1.0.0"}}`},
		// Somebody's design system, and a weather-visualisation timer.
		{"package.json": `{"dependencies": {"@b12/metronome": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@opengeoweb/metronome": "^1.0.0"}}`},
		// The bare name, which is a Node application framework.
		{"package.json": `{"dependencies": {"metronome": "^1.0.0"}}`},
		// And four PHP packages, none of them this API -- including a metric
		// aggregator that is not the metering provider.
		{"composer.json": `{"require": {"diego-ninja/laravel-metronome": "^1.0"}}`},
		{"composer.json": `{"require": {"rwslinkman/metronome": "^1.0"}}`},
		{"composer.json": `{"require": {"satunnaisuus/metronome-bundle": "^1.0"}}`},
		{"composer.json": `{"require": {"rvxlab/laravel-metronome": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "metronome") {
			t.Errorf("%v offered the metronome Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Confluence brings two shapes this file had not recorded, and both are about
// a name being right while the thing behind it is not this API.
//
// The sibling product: Atlassian sells Jira and Confluence to the same
// customers behind the same credential, and a project holding a Jira client is
// not calling Confluence.
//
// The other deployment: @atlassian-dc-mcp/confluence says Data Center in its
// own description -- the self-hosted product, whose REST API is not Cloud's
// v2. That is Shopware 5's shape arriving as a deployment choice rather than a
// version.
//
// And the sharpest miss is what this Recipe is about:
// @shogobg/markdown2confluence converts Markdown into Confluence markup and
// never makes a request -- the storage-format trap as a pure function,
// installed a quarter of a million times a month.
func TestDetectConfluenceCloudClientsAndNotItsDesignSystemOrItsDataCentre(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"confluence.js": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"confluence-api": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"confluence-cli": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@theholocron/confluence-client": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@aashari/mcp-server-atlassian-confluence": "^1.0.0"}}`},
		{"composer.json": `{"require": {"lesstif/confluence-rest-api": "^1.0"}}`},
		{"composer.json": `{"require": {"cloudplaydev/confluence-php-client": "^1.0"}}`},
		{"composer.json": `{"require": {"artemeon/confluence": "^1.0"}}`},
		{"composer.json": `{"require": {"laravel-fans/confluence": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "confluence") {
			t.Errorf("%v did not detect confluence: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's own design system, installed more than every real
		// client here put together.
		{"package.json": `{"dependencies": {"@atlaskit/embedded-confluence": "^4.0.0"}}`},
		{"package.json": `{"dependencies": {"@atlaskit/editor-confluence-transformer": "^7.0.0"}}`},
		// The in-product app runtime rather than the REST API.
		{"package.json": `{"dependencies": {"@forge/confluence-bridge": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@forge/ui-confluence": "^1.0.0"}}`},
		// The other deployment, whose API is not this one.
		{"package.json": `{"dependencies": {"@atlassian-dc-mcp/confluence": "^1.0.0"}}`},
		// The storage-format trap as a pure function, with no request in it.
		{"package.json": `{"dependencies": {"@shogobg/markdown2confluence": "^0.5.0"}}`},
		{"composer.json": `{"require": {"xen3r0/adf-converter": "^1.0"}}`},
		// Connect framework and auth middleware, not a content client.
		{"composer.json": `{"require": {"thecatontheflat/atlassian-connect-bundle": "^3.0"}}`},
		{"composer.json": `{"require": {"brezzhnev/atlassian-connect-core": "^1.0"}}`},
		{"composer.json": `{"require": {"adlogix/guzzle-atlassian-connect-middleware": "^1.0"}}`},
		// And the name nobody can verify.
		{"package.json": `{"dependencies": {"confluence": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "confluence") {
			t.Errorf("%v offered the confluence Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Service Management inverts the sibling problem Confluence has, and the
// inversion is why this mapping is small on purpose.
//
// A Jira client is not a Confluence client. But jira.js is "a library for the
// Jira Cloud platform API, Jira Agile API, and Jira Service Management API",
// so one dependency speaks three products. A package maps to one Recipe here,
// so jira.js stays with Jira: claiming it would have made every Jira project
// look like a Service Management project.
func TestDetectJiraServiceManagementWithoutClaimingTheSharedJiraClient(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"@pipedream/jira_service_desk": "^0.1.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "jiraservicemanagement") {
		t.Errorf("the dedicated client did not detect jiraservicemanagement: %+v", p.Requirements)
	}

	// The shared client stays with Jira, and says nothing about this API.
	shared, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"dependencies": {"jira.js": "^6.0.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !shared.Has(KindRecipe, "jira") {
		t.Errorf("jira.js stopped offering the jira Recipe: %+v", shared.Requirements)
	}

	if shared.Has(KindRecipe, "jiraservicemanagement") {
		t.Errorf("jira.js claimed jiraservicemanagement, which would make every Jira project look like a Service Management one: %+v", shared.Requirements)
	}

	for _, manifest := range []map[string]string{
		// A stylesheet for a portal theme, and somebody else's chat widget.
		{"package.json": `{"dependencies": {"jira-service-desk-resources": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"react-spartez-support-chat-widget": "^1.0.0"}}`},
		// Data Center, whose REST API is not this one.
		{"package.json": `{"dependencies": {"@atlassian-dc-mcp/jira": "^1.0.0"}}`},
		// A Laravel package that implements a service desk rather than
		// calling one.
		{"composer.json": `{"require": {"jeffersongoncalves/laravel-service-desk": "^1.0"}}`},
	} {
		got, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if got.Has(KindRecipe, "jiraservicemanagement") {
			t.Errorf("%v offered the jiraservicemanagement Recipe: %+v", manifest, got.Requirements)
		}
	}
}

// OpenAI has the biggest ecosystem in this file, and the exclusions matter
// more than the mappings.
//
// @ai-sdk/openai-compatible is "a foundation for implementing providers" that
// speak OpenAI's request format, and deepseek-php-client is one of the
// providers it means: a shape that has outgrown its provider, which is a
// different thing from Wix's homonym or Metronome's everyday object.
//
// @azure/openai addresses deployments on a different host with an api-version
// parameter -- the same models behind a different API. @openai/codex is a
// coding agent, not a client, and it is installed sixty-seven million times a
// month. And yethee/tiktoken counts tokens without making a request, which is
// the arithmetic this Recipe says will disagree with usage.
func TestDetectOpenAIClientsAndNotTheShapeEverybodyElseCopied(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"openai": "^4.0.0"}}`},
		{"package.json": `{"dependencies": {"@ai-sdk/openai": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@langchain/openai": "^0.3.0"}}`},
		{"composer.json": `{"require": {"openai-php/client": "^0.10"}}`},
		{"composer.json": `{"require": {"openai-php/laravel": "^0.10"}}`},
		{"composer.json": `{"require": {"orhanerday/open-ai": "^5.0"}}`},
		{"composer.json": `{"require": {"tectalic/openai": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "openai") {
			t.Errorf("%v did not detect openai: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The shape without the provider.
		{"package.json": `{"dependencies": {"@ai-sdk/openai-compatible": "^1.0.0"}}`},
		{"composer.json": `{"require": {"deepseek-php/deepseek-php-client": "^1.0"}}`},
		// Another deployment, addressing deployments rather than models.
		{"package.json": `{"dependencies": {"@azure/openai": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"azure-openai": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@azure/openai-assistants": "^1.0.0"}}`},
		// The vendor's own coding agent, which is not a client.
		{"package.json": `{"devDependencies": {"@openai/codex": "^0.5.0"}}`},
		// Instrumentation wraps a client rather than being one.
		{"package.json": `{"dependencies": {"@opentelemetry/instrumentation-openai": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@traceloop/instrumentation-openai": "^0.1.0"}}`},
		// And a tokeniser, which counts without asking anybody.
		{"composer.json": `{"require": {"yethee/tiktoken": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "openai") {
			t.Errorf("%v offered the openai Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// ElevenLabs brings a shape this file had not recorded: the vendor's own most
// installed packages are for a different transport.
//
// @elevenlabs/react and @elevenlabs/client are between them installed five
// million times a month and both are for the Agents platform -- a browser
// holding a conversation over a socket, not a server calling this REST API.
//
// @elevenlabs/types is sharper: "AsyncAPI contracts and generated TypeScript
// types". Not OpenAPI, and types alone, so it makes no request at all.
func TestDetectElevenLabsRestClientsAndNotItsRealtimeSDKs(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@elevenlabs/elevenlabs-js": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"@ai-sdk/elevenlabs": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"elevenlabs-js": "^1.0.0"}}`},
		{"composer.json": `{"require": {"ardagnsrn/elevenlabs-laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"georgehadjisavva/elevenlabs-api-client": "^1.0"}}`},
		{"composer.json": `{"require": {"onramplab/elevenlabs-api-client": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "elevenlabs") {
			t.Errorf("%v did not detect elevenlabs: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// A different transport for the same product.
		{"package.json": `{"dependencies": {"@elevenlabs/react": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@elevenlabs/client": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@elevenlabs/react-native": "^0.1.0"}}`},
		// A different specification format, and types alone.
		{"package.json": `{"dependencies": {"@elevenlabs/types": "^0.1.0"}}`},
		// Working with the output rather than calling the API, and tooling
		// for other platforms.
		{"package.json": `{"dependencies": {"@remotion/elevenlabs": "^4.0.0"}}`},
		{"package.json": `{"devDependencies": {"@elevenlabs/cli": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@livekit/agents-plugin-elevenlabs": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@mastra/voice-elevenlabs": "^0.1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "elevenlabs") {
			t.Errorf("%v offered the elevenlabs Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// PlanetScale is the clearest instance of the shape ElevenLabs introduced: the
// vendor's most installed package is for a different protocol.
//
// @planetscale/database is "A Fetch API-compatible PlanetScale database
// driver", installed a million times a month, and it runs SQL over the
// serverless query protocol. This Recipe is the management API, which creates
// databases and merges deploy requests. Same company, same account, no
// endpoint in common -- and everything built on the driver inherits the
// exclusion.
func TestDetectPlanetScaleManagementClientsAndNotItsQueryDriver(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@distilled.cloud/planetscale": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@sst-provider/planetscale": "^1.0.0"}}`},
		{"composer.json": `{"require": {"x7media/laravel-planetscale": "^1.0"}}`},
		{"composer.json": `{"require": {"bellows-app/plugin-planetscale": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "planetscale") {
			t.Errorf("%v did not detect planetscale: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The query driver, and everything built on it.
		{"package.json": `{"dependencies": {"@planetscale/database": "^1.19.0"}}`},
		{"package.json": `{"dependencies": {"@prisma/adapter-planetscale": "^6.0.0"}}`},
		{"package.json": `{"dependencies": {"kysely-planetscale": "^1.4.0"}}`},
		{"package.json": `{"dependencies": {"@mattrax/mysql-planetscale": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@storecraft/database-planetscale": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"planetscale-stream-ts": "^1.0.0"}}`},
		// And the bare name, which describes itself as connecting rather than
		// managing and is installed two thousand times a month.
		{"package.json": `{"dependencies": {"planetscale": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "planetscale") {
			t.Errorf("%v offered the planetscale Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Dropbox Sign brings three exclusions, and one of them is a rename the
// ecosystem has not followed: hellosign/hellosign-php-sdk is installed seven
// times more than dropbox/sign, so both are mapped.
//
// The npm package "dropbox" is the file storage SDK -- the sibling product,
// across a wider gulf than Jira and Confluence.
//
// And hellosign-embedded is an iframe in a browser fed a sign_url some server
// already fetched. It makes no call to this API and is installed eight hundred
// thousand times a month, against the official Node client's four hundred and
// sixty.
func TestDetectDropboxSignClientsAndNotItsEmbedOrTheFileStorage(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@dropbox/sign": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@memberjunction/esignature-dropboxsign": "^2.0.0"}}`},
		{"composer.json": `{"require": {"dropbox/sign": "^1.0"}}`},
		{"composer.json": `{"require": {"hellosign/hellosign-php-sdk": "^3.0"}}`},
		{"composer.json": `{"require": {"industrious/hellosign-laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"bukashk0zzz/hellosign-bundle": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "dropboxsign") {
			t.Errorf("%v did not detect dropboxsign: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// An iframe fed a URL somebody else fetched.
		{"package.json": `{"dependencies": {"hellosign-embedded": "^2.0.0"}}`},
		{"package.json": `{"devDependencies": {"@types/hellosign-embedded": "^2.0.0"}}`},
		{"package.json": `{"devDependencies": {"@types/hellosign-sdk": "^1.0.0"}}`},
		// The sibling product: file storage, not signatures.
		{"package.json": `{"dependencies": {"dropbox": "^10.0.0"}}`},
		{"package.json": `{"dependencies": {"@uppy/dropbox": "^4.0.0"}}`},
		// And Packagist noise that begins with the word.
		{"composer.json": `{"require": {"hellosign/hawk": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "dropboxsign") {
			t.Errorf("%v offered the dropboxsign Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Vimeo's largest near miss is the vendor's own unrelated library, at a scale
// nothing else in this file approaches -- and it is the second time that
// library is a PHP static analyser.
//
// vimeo/psalm has eighty-five million downloads and finds errors in PHP.
// Etsy's phan is the same tool doing the same job under the same kind of
// vendor name.
//
// And @vimeo/player -- an iframe controller that never calls this API -- is
// installed twenty times more than the official Node client.
func TestDetectVimeoAPIClientsAndNotItsPlayerOrItsStaticAnalyser(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@vimeo/vimeo": "^3.0.0"}}`},
		{"composer.json": `{"require": {"vimeo/vimeo-api": "^3.0"}}`},
		{"composer.json": `{"require": {"vimeo/laravel": "^5.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "vimeo") {
			t.Errorf("%v did not detect vimeo: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor name on tools that are not about video at all.
		{"composer.json": `{"require-dev": {"vimeo/psalm": "^5.0"}}`},
		{"composer.json": `{"require-dev": {"vimeo/php-mysql-engine": "^0.1"}}`},
		// The player, and everything wrapping it.
		{"package.json": `{"dependencies": {"@vimeo/player": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"vimeo-video-element": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@u-wave/react-vimeo": "^0.9.0"}}`},
		{"package.json": `{"dependencies": {"vue-vimeo-player": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"react-native-vimeo-iframe": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"lite-vimeo-embed": "^0.1.0"}}`},
		{"package.json": `{"devDependencies": {"@types/vimeo__player": "^2.0.0"}}`},
		// Watching the player rather than calling anything.
		{"package.json": `{"dependencies": {"@snowplow/browser-plugin-vimeo-tracking": "^4.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "vimeo") {
			t.Errorf("%v offered the vimeo Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Keycloak brings a shape this file had not recorded: the other half of the
// same product.
//
// Not a different deployment and not a different transport -- the thing almost
// everybody uses Keycloak for is getting a token or checking one, and the
// Admin REST API this Recipe models is the side door. stevenmaguire
// /oauth2-keycloak signs people in and has six and three quarter million
// downloads, against the PHP admin client's seven hundred and ninety thousand.
//
// keycloakify is the sharpest: "Framework to create custom Keycloak UIs", a
// build tool that produces login pages and makes no request of any kind.
func TestDetectKeycloakAdminClientsAndNotTheHalfEverybodyUses(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@keycloak/keycloak-admin-client": "^26.0.0"}}`},
		{"package.json": `{"dependencies": {"@s3pweb/keycloak-admin-client-cjs": "^26.0.0"}}`},
		{"package.json": `{"dependencies": {"@pulumi/keycloak": "^6.0.0"}}`},
		{"composer.json": `{"require": {"mohammad-waleed/keycloak-admin-client": "^0.1"}}`},
		{"composer.json": `{"require": {"fschmtt/keycloak-rest-api-client-php": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "keycloak") {
			t.Errorf("%v did not detect keycloak: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Signing people in, which is the other half of the product.
		{"composer.json": `{"require": {"stevenmaguire/oauth2-keycloak": "^5.0"}}`},
		{"composer.json": `{"require": {"socialiteproviders/keycloak": "^4.0"}}`},
		// Validating a token somebody else issued.
		{"composer.json": `{"require": {"robsontenorio/laravel-keycloak-guard": "^3.0"}}`},
		{"composer.json": `{"require": {"vizir/laravel-keycloak-web-guard": "^3.0"}}`},
		{"package.json": `{"dependencies": {"keycloak-connect": "^26.0.0"}}`},
		{"package.json": `{"dependencies": {"keycloak-angular": "^16.0.0"}}`},
		{"package.json": `{"dependencies": {"@react-keycloak/web": "^3.4.0"}}`},
		{"package.json": `{"dependencies": {"nest-keycloak-connect": "^1.10.0"}}`},
		{"package.json": `{"dependencies": {"passport-keycloak-bearer": "^1.0.0"}}`},
		// The browser adapter under the bare name.
		{"package.json": `{"dependencies": {"keycloak": "^1.0.0"}}`},
		// A build tool that produces login pages and calls nothing.
		{"package.json": `{"devDependencies": {"keycloakify": "^11.0.0"}}`},
		// And test helpers that log a user in.
		{"package.json": `{"devDependencies": {"cypress-keycloak": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "keycloak") {
			t.Errorf("%v offered the keycloak Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// FusionAuth confirms the shape Keycloak's entry names: the other half of the
// same product. Everything popular is about getting a token or checking one,
// and the administration API is the smaller half again.
//
// The sharpest detail is that fusionauth/jwt-auth-webtoken-provider is under
// the vendor's own namespace and is still a JWT library plug-in rather than a
// client, so the vendor prefix is useless as a signal here. And
// nodebb-theme-fusionauth is a forum theme.
func TestDetectFusionAuthAdminClientsAndNotTheLoginHalf(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@fusionauth/typescript-client": "^1.50.0"}}`},
		{"package.json": `{"dependencies": {"pulumi-fusionauth": "^4.0.0"}}`},
		{"package.json": `{"devDependencies": {"@fusionauth/cli": "^0.1.0"}}`},
		{"composer.json": `{"require": {"fusionauth/fusionauth-client": "^1.50"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "fusionauth") {
			t.Errorf("%v did not detect fusionauth: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Signing people in.
		{"composer.json": `{"require": {"socialiteproviders/fusionauth": "^4.0"}}`},
		{"composer.json": `{"require": {"jerryhopper/oauth2-fusionauth": "^1.0"}}`},
		{"package.json": `{"dependencies": {"passport-fusionauth": "^1.0.0"}}`},
		// Validating a token somebody else issued -- including one under the
		// vendor's own namespace.
		{"composer.json": `{"require": {"danilopolani/laravel-fusionauth-jwt": "^1.0"}}`},
		{"composer.json": `{"require": {"fusionauth/jwt-auth-webtoken-provider": "^1.0"}}`},
		{"composer.json": `{"require": {"werk365/jwtauthroles": "^1.0"}}`},
		// Logging people in from a browser.
		{"package.json": `{"dependencies": {"@fusionauth/react-sdk": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"@fusionauth/angular-sdk": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@fusionauth/vue-sdk": "^1.0.0"}}`},
		// And a forum theme.
		{"package.json": `{"dependencies": {"nodebb-theme-fusionauth": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "fusionauth") {
			t.Errorf("%v offered the fusionauth Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Atlas is an everyday word, and one collision is a prefix nobody would
// predict: Atlassian begins with it. damienharper/adf-tools is "Atlassian
// Document Format PHP Tools" at one point eight million installs a month, and
// it is about Confluence and Jira rather than databases.
//
// The PHP persistence family atlas/pdo, atlas/query and atlas/orm outnumbers
// every real client here, and Packagist has no client for this API at all.
//
// MongoDB's own scope contains the deployment shape again: @mongodb-js
// /atlas-local manages containers rather than cloud projects.
func TestDetectMongoDBAtlasAdminClientsAndNotTheWordAtlas(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"mongodb-atlas-api-client": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"awscdk-resources-mongodbatlas": "^3.0.0"}}`},
		{"package.json": `{"dependencies": {"mcp-mongodb-atlas": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "mongodbatlas") {
			t.Errorf("%v did not detect mongodbatlas: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Atlassian begins with Atlas.
		{"composer.json": `{"require": {"damienharper/adf-tools": "^3.0"}}`},
		// A PHP persistence family that happens to be called it.
		{"composer.json": `{"require": {"atlas/orm": "^3.0"}}`},
		{"composer.json": `{"require": {"atlas/query": "^1.0"}}`},
		{"composer.json": `{"require": {"atlas/pdo": "^1.0"}}`},
		{"composer.json": `{"require": {"atlas-php/atlas": "^1.0"}}`},
		{"composer.json": `{"require": {"grazulex/laravel-atlas": "^1.0"}}`},
		// The vendor's own local deployment, which is not the cloud API.
		{"package.json": `{"dependencies": {"@mongodb-js/atlas-local": "^1.0.0"}}`},
		// A different interface, and a driver-level plugin.
		{"package.json": `{"dependencies": {"mongodb-data-api": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"mongoose-atlas-search": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@mongodb-js/search-index-schema": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "mongodbatlas") {
			t.Errorf("%v offered the mongodbatlas Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Drive follows the scoped-package pattern Gmail and Calendar already use
// here: the unscoped googleapis covers every Google API at once, and a package
// maps to one Recipe.
//
// The most installed Drive packages on Packagist are Flysystem adapters, and
// they are mapped because they really do call this API -- while their whole
// purpose is to make Drive look like a filesystem, on an API whose own
// documentation says a name "isn't necessarily unique within a folder".
//
// Everything else with the name is a picker, and one of the exclusions is
// google-drive-mock: another emulator doing this Recipe's job.
func TestDetectGoogleDriveClientsAndNotThePickers(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@googleapis/drive": "^8.0.0"}}`},
		{"composer.json": `{"require": {"masbug/flysystem-google-drive-ext": "^2.0"}}`},
		{"composer.json": `{"require": {"nao-pon/flysystem-google-drive": "^1.1"}}`},
		{"composer.json": `{"require": {"yaza/laravel-google-drive-storage": "^3.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "googledrive") {
			t.Errorf("%v did not detect googledrive: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Pickers: browser widgets on a different API.
		{"package.json": `{"dependencies": {"@uppy/google-drive": "^4.0.0"}}`},
		{"package.json": `{"dependencies": {"react-google-drive-picker": "^1.2.0"}}`},
		{"package.json": `{"dependencies": {"@googleworkspace/drive-picker-element": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@googleworkspace/drive-picker-react": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"google-drive-picker": "^1.0.0"}}`},
		// Typings, which make no request.
		{"package.json": `{"devDependencies": {"@maxim_mazurok/gapi.client.drive-v3": "^0.0.1"}}`},
		// Another emulator, doing this Recipe's job.
		{"package.json": `{"devDependencies": {"google-drive-mock": "^1.0.0"}}`},
		// And a different Google product entirely.
		{"composer.json": `{"require": {"spatie/laravel-google-cloud-storage": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "googledrive") {
			t.Errorf("%v offered the googledrive Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Gemini is also a cryptocurrency exchange -- the homonym shape Wix brought,
// and this is the second instance. ccxt/ccxt is a trading API supporting a
// hundred-odd venues and Gemini is one of them, so a trading bot matches the
// word without ever meeting a model.
//
// @google-cloud/vertexai is the same models through a different host with a
// different credential, which is @azure/openai one provider along. And
// @google/gemini-cli is a coding agent rather than a client, the shape
// @openai/codex has.
func TestDetectGeminiClientsAndNotTheExchangeOrVertex(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@google/genai": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@google/generative-ai": "^0.21.0"}}`},
		{"package.json": `{"dependencies": {"@google-ai/generativelanguage": "^2.0.0"}}`},
		{"package.json": `{"dependencies": {"@langchain/google-genai": "^0.1.0"}}`},
		{"composer.json": `{"require": {"google-gemini-php/client": "^1.0"}}`},
		{"composer.json": `{"require": {"google-gemini-php/laravel": "^1.0"}}`},
		{"composer.json": `{"require": {"gemini-api-php/client": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "gemini") {
			t.Errorf("%v did not detect gemini: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The cryptocurrency exchange of the same name.
		{"composer.json": `{"require": {"ccxt/ccxt": "^4.0"}}`},
		// The same models through a different host and credential.
		{"package.json": `{"dependencies": {"@google-cloud/vertexai": "^1.0.0"}}`},
		// A coding agent, not a client.
		{"package.json": `{"devDependencies": {"@google/gemini-cli": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@google/gemini-cli-core": "^0.1.0"}}`},
		// A different vendor matching on the phrase.
		{"package.json": `{"dependencies": {"@ibm-generative-ai/node-sdk": "^3.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "gemini") {
			t.Errorf("%v offered the gemini Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Cloud Storage carries the largest numbers in this file: google/cloud-storage
// has a hundred and five million downloads, and the Flysystem and Laravel
// adapters between them nearly sixty million more.
//
// Those adapters are mapped because they call this API, and it is worth
// noticing -- as it was on Drive -- that the dominant use of an object store
// is to make it look like a filesystem, on an API whose own documentation says
// directory-like mode is a pretence.
//
// And mock-gcs is the second emulator this file has excluded in two Recipes,
// after google-drive-mock.
func TestDetectCloudStorageClientsAndNotTheUmbrellaOrTheOtherMock(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@google-cloud/storage": "^7.0.0"}}`},
		{"composer.json": `{"require": {"google/cloud-storage": "^1.40"}}`},
		{"composer.json": `{"require": {"league/flysystem-google-cloud-storage": "^3.0"}}`},
		{"composer.json": `{"require": {"superbalist/flysystem-google-storage": "^7.2"}}`},
		{"composer.json": `{"require": {"spatie/laravel-google-cloud-storage": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "googlecloudstorage") {
			t.Errorf("%v did not detect googlecloudstorage: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The umbrella client for every Google Cloud product at once.
		{"composer.json": `{"require": {"google/cloud": "^0.200"}}`},
		// Another emulator, doing this Recipe's job.
		{"package.json": `{"devDependencies": {"mock-gcs": "^1.0.0"}}`},
		// Upload sinks that write through the official client.
		{"package.json": `{"dependencies": {"multer-cloud-storage": "^3.0.0"}}`},
		{"package.json": `{"dependencies": {"@tus/gcs-store": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@payloadcms/storage-gcs": "^3.0.0"}}`},
		// And build tooling that happens to upload.
		{"package.json": `{"devDependencies": {"webpack-google-cloud-storage-plugin": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"nx-remotecache-gcs": "^5.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "googlecloudstorage") {
			t.Errorf("%v offered the googlecloudstorage Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// YouTube brings a shape this file had not recorded, and the biggest number on
// the page produced it: matched on a simile.
//
// Packagist's top results for "youtube" are hashids/hashids at fifty-four
// million downloads and vinkla/hashids at fifteen, and sqids/sqids says
// outright what all three are -- "Generate short YouTube-looking IDs from
// numbers". They share neither the name, the vendor nor the domain; they
// describe their output as looking like something YouTube makes.
//
// And norkunas/youtube-dl-php wraps youtube-dl and yt-dlp, which get at videos
// by scraping rather than by asking the Data API -- the tool people reach for
// instead of this one.
func TestDetectYouTubeDataClientsAndNotASimileOrAScraper(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@googleapis/youtube": "^15.0.0"}}`},
		{"package.json": `{"dependencies": {"simple-youtube-api": "^5.2.0"}}`},
		{"package.json": `{"dependencies": {"youtube-node": "^1.3.0"}}`},
		{"composer.json": `{"require": {"alaouy/youtube": "^3.0"}}`},
		{"composer.json": `{"require": {"madcoda/php-youtube-api": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "youtube") {
			t.Errorf("%v did not detect youtube: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Matched on a simile: libraries whose output merely looks like a
		// YouTube identifier.
		{"composer.json": `{"require": {"hashids/hashids": "^5.0"}}`},
		{"composer.json": `{"require": {"vinkla/hashids": "^12.0"}}`},
		{"composer.json": `{"require": {"sqids/sqids": "^0.4"}}`},
		// The tool people use instead of this API, by scraping.
		{"composer.json": `{"require": {"norkunas/youtube-dl-php": "^2.0"}}`},
		// And oEmbed, which reads a public page rather than the Data API.
		{"composer.json": `{"require": {"mpratt/embera": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "youtube") {
			t.Errorf("%v offered the youtube Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Vault is the most contested name in this file. Four vendors ship a product
// called one, and the biggest is not HashiCorp's: @azure/keyvault-common has
// twenty-three million downloads against node-vault's 1.7 million, Oracle ships
// oci-vault, and @googleapis/vault is Google Vault -- e-discovery for
// Workspace, holding no secrets at all.
//
// Then the words that are not products: spryker/vault is a commerce framework
// module at 2.2 million, the file-vault family encrypts files, ansible-vault is
// a file format, and the top Packagist result for the word is ccxt/ccxt, a
// cryptocurrency exchange library.
//
// Two exclusions are deliberate. @testcontainers/vault starts a real Vault in
// Docker, which is what somebody reaches for instead of an emulator, and
// @pulumi/vault provisions the server rather than reading secrets from it.
func TestDetectVaultClientsAndNotTheOtherFourVaults(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"node-vault": "^0.10.2"}}`},
		{"package.json": `{"dependencies": {"hashi-vault-js": "^0.4.14"}}`},
		{"package.json": `{"dependencies": {"@litehex/node-vault": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"vault-api": "^1.0.0"}}`},
		{"composer.json": `{"require": {"csharpru/vault-php": "^4.2"}}`},
		{"composer.json": `{"require": {"jippi/vault-php-sdk": "^2.4"}}`},
		{"composer.json": `{"require": {"mittwald/vault-php": "^1.0"}}`},
		{"composer.json": `{"require": {"tokenly/laravel-vault": "^0.1"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "vault") {
			t.Errorf("%v did not detect vault: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Three other vendors' Vaults, and the largest numbers of the lot.
		{"package.json": `{"dependencies": {"@azure/keyvault-secrets": "^4.8.0"}}`},
		{"package.json": `{"dependencies": {"oci-vault": "^2.90.0"}}`},
		{"package.json": `{"dependencies": {"@googleapis/vault": "^11.0.0"}}`},
		{"composer.json": `{"require": {"keboola/azure-key-vault-client": "^3.0"}}`},
		// Words that are not products.
		{"composer.json": `{"require": {"spryker/vault": "^1.0"}}`},
		{"composer.json": `{"require": {"soarecostin/file-vault": "^1.0"}}`},
		{"composer.json": `{"require": {"dotenv-org/phpdotenv-vault": "^0.1"}}`},
		{"package.json": `{"dependencies": {"ansible-vault": "^1.1.1"}}`},
		{"package.json": `{"dependencies": {"@apideck/vault-react": "^2.0.0"}}`},
		// The real server in Docker, which is the alternative to a Recipe.
		{"package.json": `{"devDependencies": {"@testcontainers/vault": "^10.0.0"}}`},
		// And infrastructure-as-code, which configures Vault rather than
		// reading secrets from it.
		{"package.json": `{"dependencies": {"@pulumi/vault": "^6.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "vault") {
			t.Errorf("%v offered the vault Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Unleash brings the near-miss shape next door to YouTube's simile: matched on
// an imperative. "Unleash" is a verb package taglines love, so the biggest npm
// result for it is precinct at twenty million downloads -- a dependency parser
// described as "Unleash the detectives" -- with @remirror/extension-code-block
// and veewee/reflecta close behind on the same trick.
//
// Three exclusions are deliberate rather than accidental: unleash-server is the
// product itself, @pulumiverse/unleash provisions flags rather than reading
// them, and @unleash/mcp is an agent tool rather than an SDK.
func TestDetectUnleashClientsAndNotAnImperative(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"unleash-client": "^6.1.0"}}`},
		{"package.json": `{"dependencies": {"unleash-proxy-client": "^3.7.0"}}`},
		{"package.json": `{"dependencies": {"@unleash/proxy-client-react": "^4.3.0"}}`},
		{"package.json": `{"dependencies": {"@unleash/nextjs": "^1.4.0"}}`},
		{"composer.json": `{"require": {"unleash/client": "^2.6"}}`},
		{"composer.json": `{"require": {"unleash/symfony-client-bundle": "^2.5"}}`},
		{"composer.json": `{"require": {"mikefrancis/laravel-unleash": "^7.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "unleash") {
			t.Errorf("%v did not detect unleash: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Matched on an imperative: packages whose taglines start "Unleash".
		{"package.json": `{"dependencies": {"precinct": "^12.0.0"}}`},
		{"package.json": `{"dependencies": {"@remirror/extension-code-block": "^2.0.0"}}`},
		{"composer.json": `{"require": {"veewee/reflecta": "^0.5"}}`},
		// A different vendor holding the same word.
		{"composer.json": `{"require-dev": {"unleashedtech/php-coding-standard": "^4.0"}}`},
		// The product itself, which needs no emulator of itself.
		{"package.json": `{"dependencies": {"unleash-server": "^6.0.0"}}`},
		// Infrastructure-as-code, and an agent tool.
		{"package.json": `{"dependencies": {"@pulumiverse/unleash": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@unleash/mcp": "^0.1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "unleash") {
			t.Errorf("%v offered the unleash Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Pinecone settles two exclusions that have now come up three Recipes running.
//
// Infrastructure-as-code: @pinecone-database/pulumi joins @pulumi/vault and
// @pulumiverse/unleash. A project holding one provisions the provider rather
// than calling it at runtime. Agent tooling: @pinecone-database/mcp beside
// @unleash/mcp. Both are published under the vendor's own scope, which is why
// the name alone cannot decide it.
//
// And a new kind: @traceloop/instrumentation-pinecone is OpenTelemetry
// instrumentation that wraps the official SDK. It observes a client rather than
// being one, and any project holding it already holds the client.
func TestDetectPineconeClientsAndNotItsToolingOrATranspiler(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"@pinecone-database/pinecone": "^6.0.0"}}`},
		{"package.json": `{"dependencies": {"@langchain/pinecone": "^0.2.0"}}`},
		{"package.json": `{"dependencies": {"@mastra/pinecone": "^0.3.0"}}`},
		{"composer.json": `{"require": {"probots-io/pinecone-php": "^1.0"}}`},
		{"composer.json": `{"require": {"symfony/ai-pinecone-store": "^0.1"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "pinecone") {
			t.Errorf("%v did not detect pinecone: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's own tooling, which is not a client.
		{"package.json": `{"dependencies": {"@pinecone-database/pulumi": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@pinecone-database/mcp": "^1.0.0"}}`},
		// Instrumentation that wraps the client rather than being one.
		{"package.json": `{"dependencies": {"@traceloop/instrumentation-pinecone": "^0.14.0"}}`},
		// The bare name is a JavaScript-to-Lua converter.
		{"package.json": `{"devDependencies": {"pinecone": "^0.2.0"}}`},
		// The generic abstraction, without the Pinecone bridge.
		{"composer.json": `{"require": {"symfony/ai-store": "^0.1"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "pinecone") {
			t.Errorf("%v offered the pinecone Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Hookdeck is a small ecosystem where nearly everything under the vendor's own
// npm scope is something other than a client, so the scope is useless as a
// signal and the exclusions are the whole job.
//
// @hookdeck/outpost-sdk is the sharpest, and a new shape for this file: Outpost
// is Hookdeck's open-source, self-hosted outbound-webhook product with an API
// of its own. Same vendor, same scope, different surface.
//
// The rest are rules three Recipes have already settled: Pulumi bridges and a
// declarative deploy CLI are infrastructure-as-code, two packages are agent
// tooling, and one generates documentation.
func TestDetectHookdeckClientsAndNotTheVendorsOtherProduct(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"devDependencies": {"hookdeck-cli": "^0.10.0"}}`},
		{"package.json": `{"dependencies": {"@hookdeck/vercel": "^0.2.0"}}`},
		{"package.json": `{"dependencies": {"@hookdeck/pubsub": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"n8n-nodes-hookdeck": "^0.1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "hookdeck") {
			t.Errorf("%v did not detect hookdeck: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's other product, under the vendor's own scope.
		{"package.json": `{"dependencies": {"@hookdeck/outpost-sdk": "^0.5.0"}}`},
		// Infrastructure-as-code, for the fourth Recipe running.
		{"package.json": `{"dependencies": {"@nectarsocial/hookdeck": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"@sst-provider/hookdeck": "^0.1.0"}}`},
		{"package.json": `{"devDependencies": {"@toppy/hookdeck-deploy-cli": "^1.0.0"}}`},
		// Agent tooling, for the third.
		{"package.json": `{"devDependencies": {"@stackcurious/hookdeck-mcp": "^0.1.0"}}`},
		// And a documentation generator.
		{"package.json": `{"devDependencies": {"@hookdeck/eventcatalog-generator": "^0.1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "hookdeck") {
			t.Errorf("%v offered the hookdeck Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Grafana has the thinnest mapping in this file, and the reason matters more
// than the mapping: eighteen npm results for the word and not one of them calls
// this API.
//
// The whole @grafana/ scope is plugin-development tooling -- panels that run
// inside Grafana and reach its internals directly. On Packagist the largest
// numbers belong to Loki, a different product from the same vendor, beside
// Prometheus exporters that publish metrics to be scraped and call nobody.
//
// The Foundation SDK builds dashboard documents as code and does not send them.
func TestDetectGrafanaClientsAndNotItsPluginToolkitOrLoki(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"saschahemleb/php-grafana-api-client": "^0.4"}}`},
		{"composer.json": `{"require": {"saschahemleb/laravel-grafana": "^0.3"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "grafana") {
			t.Errorf("%v did not detect grafana: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The vendor's scope, which is entirely plugin tooling.
		{"package.json": `{"dependencies": {"@grafana/data": "^11.0.0"}}`},
		{"package.json": `{"dependencies": {"@grafana/ui": "^11.0.0"}}`},
		{"package.json": `{"dependencies": {"@grafana/runtime": "^11.0.0"}}`},
		{"package.json": `{"devDependencies": {"@grafana/plugin-e2e": "^1.0.0"}}`},
		// A builder rather than a client.
		{"package.json": `{"dependencies": {"@grafana/grafana-foundation-sdk": "^11.0.0"}}`},
		{"composer.json": `{"require": {"grafana/foundation-sdk": "^11.0"}}`},
		// A different product from the same vendor, and the biggest numbers.
		{"composer.json": `{"require": {"itspire/monolog-loki": "^2.0"}}`},
		{"composer.json": `{"require": {"cebe/yii2-loki-log-target": "^1.0"}}`},
		// Metrics that get scraped rather than sent.
		{"composer.json": `{"require": {"iamfarhad/laravel-prometheus": "^1.0"}}`},
		{"composer.json": `{"require": {"renoki-co/horizon-exporter": "^3.0"}}`},
		// Infrastructure-as-code, for the fifth Recipe running.
		{"package.json": `{"dependencies": {"@pulumiverse/grafana": "^0.4.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "grafana") {
			t.Errorf("%v offered the grafana Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// ConfigCat is the starkest version of a shape this file has now met three
// Recipes running: the vendor's popular client speaks a different surface.
//
// configcat/configcat-client has 1.29 million downloads and calls a CDN, not
// api.configcat.com -- a ConfigCat SDK fetches a static configuration file and
// evaluates flags locally. The Public Management API this Recipe serves is what
// dashboards and CI scripts use, and its clients have about five thousand
// downloads between them.
func TestDetectConfigCatManagementClientsAndNotItsSDKs(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"devDependencies": {"configcat-publicapi-node-client": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"ng-configcat-publicapi": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "configcat") {
			t.Errorf("%v did not detect configcat: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The SDKs, which read a configuration file from the CDN.
		{"composer.json": `{"require": {"configcat/configcat-client": "^9.0"}}`},
		{"package.json": `{"dependencies": {"configcat-common": "^9.0.0"}}`},
		{"package.json": `{"dependencies": {"@configcat/sdk": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"configcat-react": "^4.0.0"}}`},
		// Wrappers around those SDKs, so they are on the CDN side too.
		{"package.json": `{"dependencies": {"@openfeature/config-cat-provider": "^0.7.0"}}`},
		{"composer.json": `{"require": {"pod-point/laravel-configcat": "^2.0"}}`},
		// Infrastructure-as-code, for the sixth Recipe running.
		{"package.json": `{"dependencies": {"@pulumiverse/configcat": "^0.1.0"}}`},
		// And agent tooling, even though it names this exact API.
		{"package.json": `{"devDependencies": {"@configcat/mcp-server": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "configcat") {
			t.Errorf("%v offered the configcat Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// SonarCloud brings an exclusion kind this file had not recorded, and it owns
// the biggest number on the page that is not the scanner: packages that write a
// file for another program to upload.
//
// vitest-sonar-reporter has eight hundred and forty thousand downloads and never
// makes a network request -- it serialises test results to disk in a format the
// scanner later reads. The scanners themselves are bootstrappers that download
// the sonar-scanner Java CLI and run it, submitting analyses over a protocol
// this Recipe does not serve.
//
// sonarqube-verify is mapped despite the neighbourhood, because what it does
// after an analysis is poll the Compute Engine task and read the quality gate.
func TestDetectSonarWebAPIClientsAndNotReportersOrScanners(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"sonarqube-api-client": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"sonarqube-verify": "^1.0.0"}}`},
		{"composer.json": `{"require": {"forgeqc/sonarqube-api-client": "^1.0"}}`},
		{"composer.json": `{"require": {"spirit-dev/php-sonarqube-api": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "sonarcloud") {
			t.Errorf("%v did not detect sonarcloud: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Reporters, which write a file and call nobody.
		{"package.json": `{"devDependencies": {"vitest-sonar-reporter": "^2.0.0"}}`},
		{"package.json": `{"devDependencies": {"karma-sonarqube-reporter": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"cypress-sonarqube-reporter": "^5.0.0"}}`},
		{"composer.json": `{"require-dev": {"symbiote/phpcs-sonar": "^1.0"}}`},
		// Scanners, which download a Java CLI and run it.
		{"package.json": `{"devDependencies": {"sonarqube-scanner": "^4.0.0"}}`},
		{"package.json": `{"devDependencies": {"@sonar/scan": "^4.0.0"}}`},
		{"composer.json": `{"require-dev": {"rogervila/php-sonarqube-scanner": "^1.0"}}`},
		// A lint configuration, and the top Packagist result for the word.
		{"package.json": `{"devDependencies": {"eslint-config-sonarqube": "^1.0.0"}}`},
		{"composer.json": `{"require-dev": {"nimut/phpunit-merger": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "sonarcloud") {
			t.Errorf("%v offered the sonarcloud Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// GrowthBook repeats ConfigCat's split immediately afterwards, which makes it a
// rule about feature-flag platforms rather than a fact about one vendor: there
// is a delivery surface that SDKs read and a management API that everything else
// uses, and they are different APIs on different hosts.
//
// @growthbook/growthbook has 3.49 million downloads and the PHP SDK 3.92
// million, and neither calls the REST API this Recipe serves. Exactly one
// package here does, and it says so in its own description: growthbook,
// "Command-line interface for the GrowthBook REST API".
//
// Unleash is the contrast: its client and admin APIs share a host, so one Recipe
// covers both and its SDKs are mapped.
func TestDetectGrowthBookRestClientsAndNotItsSDKs(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"package.json": `{"devDependencies": {"growthbook": "^0.2.0"}}`,
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "growthbook") {
		t.Errorf("the REST CLI did not detect growthbook: %+v", p.Requirements)
	}

	for _, manifest := range []map[string]string{
		// The SDKs, which read feature definitions from a CDN or a proxy.
		{"package.json": `{"dependencies": {"@growthbook/growthbook": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@growthbook/growthbook-react": "^1.0.0"}}`},
		{"composer.json": `{"require": {"growthbook/growthbook": "^1.2"}}`},
		// Wrappers around those SDKs.
		{"package.json": `{"dependencies": {"@flags-sdk/growthbook": "^0.1.0"}}`},
		{"package.json": `{"dependencies": {"nuxt-growthbook": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@growthbook/edge-cloudflare": "^1.0.0"}}`},
		{"composer.json": `{"require": {"gathern/growthbook-openfeature-pro": "^1.0"}}`},
		// The product's own server component, and agent tooling.
		{"package.json": `{"dependencies": {"@growthbook/proxy": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@growthbook/mcp": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "growthbook") {
			t.Errorf("%v offered the growthbook Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// PyPI carries the largest number this file has ever excluded, and it is not
// even a homonym: composer/installers, at a hundred and forty-eight million
// downloads, is the top Packagist result for "pypi" and is a Composer library
// installer. A registry's ranking is not evidence of anything.
//
// The real split is between reading and publishing. semantic-release-pypi and
// its neighbours upload a built distribution over twine's protocol, which is
// not this API -- the same line drawn at the Sonar scanners.
func TestDetectPyPIReadersAndNotPublishers(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"pypi-info": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"pypi": "^0.1.0"}}`},
		{"package.json": `{"devDependencies": {"muaddib-scanner": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"bismar": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "pypi") {
			t.Errorf("%v did not detect pypi: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Publishers, which speak twine's upload protocol.
		{"package.json": `{"devDependencies": {"semantic-release-pypi": "^3.0.0"}}`},
		{"package.json": `{"devDependencies": {"@lets-release/pypi": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"foundry-release-pypi": "^1.0.0"}}`},
		{"composer.json": `{"require-dev": {"hiqdev/hidev-pypi": "^0.1"}}`},
		// And the biggest number on the page, which is a Composer installer.
		{"composer.json": `{"require": {"composer/installers": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "pypi") {
			t.Errorf("%v offered the pypi Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// crates.io brings two near-miss kinds at once, and one of them is new.
//
// The homonym is a single letter: every Packagist result for "crates" is
// CrateDB, a distributed SQL database that lives at crate.io. crate/crate-pdo
// has ninety-six thousand downloads and nothing to do with the Rust registry.
//
// The new kind is matched on a badge. ts-gettext-extractor at fifty-five
// thousand and the Tauri plugins beside it surface because their READMEs carry
// a shields.io badge whose URL contains crates.io. The badge is honest -- there
// really is a crate -- and the match is not, because nothing in the npm package
// calls the registry API.
func TestDetectCratesIOClientsAndNotCrateDBOrABadge(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"crates.io": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "cratesio") {
			t.Errorf("%v did not detect cratesio: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// CrateDB, which is crate.io and a database.
		{"composer.json": `{"require": {"crate/crate-pdo": "^2.0"}}`},
		{"composer.json": `{"require": {"crate/crate-dbal": "^2.0"}}`},
		{"composer.json": `{"require": {"ratkor/laravel-crate.io": "^1.0"}}`},
		// Matched on a README badge.
		{"package.json": `{"devDependencies": {"ts-gettext-extractor": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"tauri-plugin-keyring-api": "^1.0.0"}}`},
		// Publishers, which speak cargo's upload protocol.
		{"package.json": `{"devDependencies": {"@auto-it/crates": "^11.0.0"}}`},
		{"package.json": `{"devDependencies": {"@releasekit/publish": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "cratesio") {
			t.Errorf("%v offered the cratesio Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// RubyGems brings another new near-miss kind, and it owns the biggest number on
// the page: matched on a semantic. @snyk/ruby-semver at forty-two thousand
// downloads is "a node-semver compatible API with RubyGems semantics" -- a
// JavaScript implementation of how RubyGems compares versions. It knows more
// about this registry than any client here does and never sends it a request.
//
// A badge match, which crates.io records, means the package is published
// somewhere. A semantic match means it reimplements the provider's rules.
func TestDetectRubyGemsClientsAndNotItsSemantics(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"rubygems": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"gem-count": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"supply-chain-guard": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@agntn/registries": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "rubygems") {
			t.Errorf("%v did not detect rubygems: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Matched on a semantic: RubyGems' version rules, in JavaScript.
		{"package.json": `{"dependencies": {"@snyk/ruby-semver": "^3.0.0"}}`},
		// Publishers, which push a built gem.
		{"package.json": `{"devDependencies": {"@webhippie/semantic-release-ruby": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"@tegami/gem": "^1.0.0"}}`},
		// Agent tooling, for the fifth Recipe running.
		{"package.json": `{"devDependencies": {"@pipeworx/mcp-rubygems": "^1.0.0"}}`},
		// And Packagist noise, which mentions no Ruby at all.
		{"composer.json": `{"require": {"judev/php-htmltruncator": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "rubygems") {
			t.Errorf("%v offered the rubygems Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Packagist is the first registry in this run with real clients on its own
// language's registry: spatie/packagist-api at 2.4 million downloads and
// knplabs/packagist-api at 1.5 million.
//
// The exclusion that took the most care is private-packagist/api-client.
// Private Packagist is the same company's commercial product -- a hosted
// private repository at packagist.com with an API of its own -- so it is the
// Hookdeck/Outpost shape again, and sharper, because these two share a vendor
// prefix as well as a brand.
//
// And the badge match recorded at crates.io recurs, which makes it a pattern:
// two of the three packages surfacing on it are Angular projects with no PHP.
func TestDetectPackagistClientsAndNotPrivatePackagist(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"composer.json": `{"require": {"spatie/packagist-api": "^2.0"}}`},
		{"composer.json": `{"require": {"knplabs/packagist-api": "^1.0"}}`},
		{"package.json": `{"dependencies": {"@agonyz/packagist-api-client": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"gatsby-source-packagist": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "packagist") {
			t.Errorf("%v did not detect packagist: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The same company's commercial product, with an API of its own.
		{"composer.json": `{"require": {"private-packagist/api-client": "^1.0"}}`},
		{"composer.json": `{"require": {"private-packagist/bitbucket-api": "^1.0"}}`},
		// Matched on a README badge, for the second Recipe running.
		{"package.json": `{"dependencies": {"angular-dashboard-framework": "^0.14.0"}}`},
		// And fuzzy matches with nothing to do with the registry.
		{"composer.json": `{"require": {"khill/lavacharts": "^3.1"}}`},
		{"composer.json": `{"require": {"square/connect": "^2.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "packagist") {
			t.Errorf("%v offered the packagist Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Hex's homonym is the fixture this Recipe is written about. The npm package
// called phoenix, at 5.3 million downloads, is the official JavaScript client
// for the Phoenix framework's channels -- same name, same project, and it
// speaks a websocket protocol to an application rather than HTTP to a registry.
//
// The badge match turns up for the third Recipe running, which makes it
// reliable enough to expect on any provider whose users publish elsewhere.
func TestDetectHexClientsAndNotThePhoenixChannelClient(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"hex-api-client": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"cerebro-hex": "^1.0.0"}}`},
		{"package.json": `{"dependencies": {"@api-hooks/hex": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "hexpm") {
			t.Errorf("%v did not detect hexpm: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// The Phoenix channel client, and the biggest number on the page.
		{"package.json": `{"dependencies": {"phoenix": "^1.8.0"}}`},
		// Matched on a README badge, for the third Recipe running.
		{"package.json": `{"dependencies": {"@mrdotb/live-react": "^0.1.0"}}`},
		// Agent tooling, for the sixth.
		{"package.json": `{"devDependencies": {"@pipeworx/mcp-hex-pm": "^1.0.0"}}`},
		// And Packagist prefix collisions, which share four letters.
		{"composer.json": `{"require": {"hexpang/ssh-client": "^1.0"}}`},
		{"composer.json": `{"require": {"hexmedia/yaml-linter": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "hexpm") {
			t.Errorf("%v offered the hexpm Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// Almost everything the word returns on npm shells out to nuget.exe: nuget-bin
// downloads the binary, and the grunt and gulp tasks drive it to pack, push and
// restore. That is the line drawn at the Sonar scanners and PyPI's publishers.
//
// Two near-miss kinds recur and are now established. @snyk/nuget-semver is a
// semantic version parser -- matched on a semantic, the second after
// @snyk/ruby-semver and from the same publisher. The @cratis/arc packages carry
// an img.shields.io/nuget badge -- matched on a badge, for the fourth Recipe
// running.
func TestDetectNuGetClientsAndNotTheExeWrappers(t *testing.T) {
	for _, manifest := range []map[string]string{
		{"package.json": `{"dependencies": {"nuget": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"snyk-nuget-plugin": "^2.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if !p.Has(KindRecipe, "nuget") {
			t.Errorf("%v did not detect nuget: %+v", manifest, p.Requirements)
		}
	}

	for _, manifest := range []map[string]string{
		// Wrappers that run nuget.exe.
		{"package.json": `{"devDependencies": {"nuget-bin": "^1.0.0"}}`},
		{"package.json": `{"devDependencies": {"grunt-nuget": "^0.2.0"}}`},
		{"package.json": `{"devDependencies": {"gulp-nuget-restore": "^0.1.0"}}`},
		// Matched on a semantic, for the second Recipe.
		{"package.json": `{"dependencies": {"@snyk/nuget-semver": "^2.0.0"}}`},
		// Matched on a badge, for the fourth.
		{"package.json": `{"dependencies": {"@cratis/arc": "^1.0.0"}}`},
		// A repository server, and a JSON validator spelled with two g's.
		{"composer.json": `{"require": {"melonsmasher/chocolatier": "^1.0"}}`},
		{"composer.json": `{"require": {"activerules/nugget": "^1.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "nuget") {
			t.Errorf("%v offered the nuget Recipe: %+v", manifest, p.Requirements)
		}
	}
}

// The Go proxy carries the sharpest matched-on-a-semantic yet, and the third:
// golang.org/x/mod implements module.EscapePath, the exact rule this Recipe is
// about, and never opens a connection. After @snyk/ruby-semver and
// @snyk/nuget-semver the pattern is established -- a provider's rules get
// reimplemented as a library more often than its API gets called, and the
// library is the more popular package every time.
//
// The rest of the neighbourhood is the toolchain: go-mod-upgrade shells out to
// go list, which reaches the proxy through cmd/go rather than over an API a
// test could redirect. gomajor is mapped because it fetches @v/list itself.
func TestDetectGoProxyClientsAndNotTheEscapingLibrary(t *testing.T) {
	p, err := Detect(writeProject(t, map[string]string{
		"go.mod": "module example.com/app\n\nrequire github.com/icholy/gomajor v0.13.0\n",
	}))
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if !p.Has(KindRecipe, "goproxy") {
		t.Errorf("gomajor did not detect goproxy: %+v", p.Requirements)
	}

	for _, manifest := range []map[string]string{
		// The library that implements the escaping and calls nobody.
		{"go.mod": "module example.com/app\n\nrequire golang.org/x/mod v0.21.0\n"},
		// The toolchain wrapper.
		{"go.mod": "module example.com/app\n\nrequire github.com/oligot/go-mod-upgrade v0.10.0\n"},
		// A GOPROXY implementation is the product, not a client of it.
		{"package.json": `{"dependencies": {"gosub-goproxy": "^1.0.0"}}`},
		// And agent tooling, for the seventh Recipe running.
		{"package.json": `{"devDependencies": {"@pipeworx/mcp-goproxy": "^1.0.0"}}`},
	} {
		p, err := Detect(writeProject(t, manifest))
		if err != nil {
			t.Fatalf("Detect(%v): %v", manifest, err)
		}

		if p.Has(KindRecipe, "goproxy") {
			t.Errorf("%v offered the goproxy Recipe: %+v", manifest, p.Requirements)
		}
	}
}
