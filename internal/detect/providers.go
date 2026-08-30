package detect

import "strings"

// provider maps a dependency to the Recipe that emulates it.
//
// The mapping is deliberately explicit rather than fuzzy. A wrong guess here is
// worse than no guess: booting the wrong fake sends a developer chasing a bug
// that doesn't exist.
type provider struct {
	// recipe is the canonical Recipe name.
	recipe string
	// packages are exact dependency names, per ecosystem.
	composer []string
	npm      []string
	gomod    []string
}

func providers() []provider {
	return []provider{
		{
			recipe:   "stripe",
			composer: []string{"stripe/stripe-php"},
			npm:      []string{"stripe"},
			gomod:    []string{"github.com/stripe/stripe-go"},
		},
		{
			recipe:   "shopify",
			composer: []string{"shopify/shopify-api"},
			npm:      []string{"@shopify/shopify-api", "@shopify/admin-api-client"},
			gomod:    []string{"github.com/bold-commerce/go-shopify"},
		},
		{
			// Only the package that serves one API is mapped here.
			//
			// A dependency maps to exactly one Recipe -- the lookup returns
			// the first match -- and Shopify's two most-installed libraries
			// each serve both APIs. @shopify/shopify-api exposes a REST
			// client and a GraphQL client side by side, and
			// @shopify/admin-api-client documents both. Moving either one to
			// this Recipe would silently change which emulator existing
			// projects are offered, on a guess about which half of the
			// library they call.
			//
			// @shopify/graphql-client is GraphQL and nothing else, so it is
			// the one dependency that answers the question by itself.
			recipe: "shopifygraphql",
			npm:    []string{"@shopify/graphql-client"},
		},
		{
			recipe:   "bigcommerce",
			composer: []string{"bigcommerce/api"},
			npm:      []string{"node-bigcommerce", "@bigcommerce/checkout-sdk"},
		},
		{
			// The npm package called printify is a print-CSS helper -- "easy
			// HTML print formatting" -- and has nothing to do with the
			// print-on-demand company. The word means putting ink on paper
			// before it means putting ink on a t-shirt.
			//
			// The PHP clients are small and several, none official, which is
			// not a reason to leave them out: the bar is whether a dependency
			// identifies the provider.
			recipe: "printify",
			composer: []string{
				"jacob-hyde/printify",
				"garissman/printify",
				"ahsan/printify-laravel",
				"8fold/php-printify-sdk",
			},
		},
		{
			// The sibling problem, inverted -- and a mapping this file
			// deliberately does not make.
			//
			// Confluence's entry above notes that a Jira client is not a
			// Confluence client. Service Management is the other way round:
			// jira.js describes itself as a "library for the Jira Cloud
			// platform API, Jira Agile API, and Jira Service Management API",
			// and it is installed over a million times a month. One
			// dependency speaks three products.
			//
			// A package maps to one Recipe here, and TestNoPackageMapsToTwo
			// Recipes enforces it. So jira.js stays with Jira, where it
			// mostly belongs, and this API is reachable only through the
			// small dedicated client below -- which is the honest state of
			// the world rather than a shortfall in the search. Claiming
			// jira.js for Service Management would have made every Jira
			// project look like a Service Management project.
			//
			// The rest are near misses: jira-service-desk-resources is a
			// stylesheet for a portal theme, react-spartez-support-chat
			// -widget embeds somebody else's chat product into a portal page,
			// and @atlassian-dc-mcp/jira is Data Center, whose REST API is
			// not this one. Packagist has nothing at all -- searching it for
			// the product returns Freshdesk SDKs, Zoho Desk SDKs and a
			// Laravel package that implements a service desk rather than
			// calling one.
			recipe: "jiraservicemanagement",
			npm:    []string{"@pipedream/jira_service_desk"},
		},
		{
			// Two new shapes here, and both are about a name being right
			// while the thing behind it is not this API.
			//
			// The sibling product. Atlassian sells Jira and Confluence to the
			// same customers, behind the same credential, on the same host,
			// and a project holding a Jira client is not calling Confluence.
			// Nothing distinguishes them but the package name, so the Jira
			// Recipe's mappings and these are kept apart on purpose.
			//
			// The other deployment. @atlassian-dc-mcp/confluence says Data
			// Center in its own description: the self-hosted product, whose
			// REST API is a different one from Cloud's v2. That is the shape
			// Shopware 5 and epicformbuilder/wixhive-php-api have -- a
			// previous or parallel generation of the same vendor's API --
			// arriving here as a deployment choice rather than a version.
			//
			// Then the shapes already on file. @atlaskit/embedded-confluence
			// and @atlaskit/editor-confluence-transformer are Atlassian's own
			// design system, and the first has three hundred thousand monthly
			// installs, which is more than every real client here put
			// together. @forge/confluence-bridge and @forge/ui-confluence are
			// the in-product app runtime rather than the REST API.
			//
			// And the sharpest miss is the one this Recipe is about:
			// @shogobg/markdown2confluence converts Markdown into Confluence
			// markup and never makes a request. It is the storage-format trap
			// as a pure function, installed a quarter of a million times a
			// month by people who will never call this API at all.
			//
			// The bare npm "confluence" describes itself as "confluence
			// ==========" and is installed 76 times a month, which is the
			// unverifiable-name state the backlog records against Marqeta.
			recipe: "confluence",
			composer: []string{
				"lesstif/confluence-rest-api",
				"cloudplaydev/confluence-php-client",
				"artemeon/confluence",
				"rainflute/confluence-php-client",
				"laravel-fans/confluence",
				"sinkcup/confluence-sdk-php",
			},
			npm: []string{
				"confluence.js",
				"confluence-api",
				"confluence-cli",
				"@theholocron/confluence-client",
				"@aashari/mcp-server-atlassian-confluence",
			},
		},
		{
			// The name is an everyday object, and that is most of the
			// problem. @dougflip/metronome, expo-precision-metronome,
			// strumming-metronome, react-native-metronome-module,
			// @pickupmusic/metronome, brgndy-react-metronome,
			// @forestream/metronome and com.sonosthesia.metronome are all
			// what the word means: things that click in time. Wix was a
			// homonym for a build toolset and this is a homonym for a piece
			// of hardware, which is commoner and produces more collisions.
			//
			// Then a different company with the same name: metronome
			// -contracts, metronome-wallet-core and metronome-wallet-ui-logic
			// are the Metronome ERC-20 token and its Ethereum wallet. Nothing
			// in a manifest distinguishes one Metronome from the other except
			// the scope.
			//
			// The bare npm "metronome" is a Node application framework, which
			// is the shape Lago and Polar both met.
			//
			// And the closest miss of the lot is diego-ninja/laravel
			// -metronome, "a metric aggregator and processor for Laravel" --
			// a metric aggregator named metronome that is not the metering
			// provider named Metronome. rwslinkman/metronome is a Symfony
			// test utility, satunnaisuus/metronome-bundle is an admin panel
			// constructor and rvxlab/laravel-metronome runs the scheduler on
			// an event loop. Packagist has no client for this API at all.
			//
			// What is left is the scope, and it is the busiest client in this
			// file: @metronome/sdk is installed over a million times a month.
			recipe: "metronome",
			npm: []string{
				"@metronome/sdk",
				"@metronome/mcp",
			},
		},
		{
			// The busiest exclusion list here, and the reason is that "wix"
			// is three different things on two registries.
			//
			// It is a homonym. WiX is the Windows Installer XML toolset, and
			// electron-wix-msi, @electron-forge/maker-wix and
			// @mongodb-js/electron-wix-msi build MSI packages with it. They
			// have nothing to do with Wix.com, they belong to other people --
			// felixrieseberg/electron-wix-msi -- and between them they are
			// installed far more than any client here. No shape recorded in
			// this file covered that: Lago met Lagoon and Polar met Polaris,
			// which are longer words that contain a provider's name. This is
			// the same word meaning something else.
			//
			// It is a prefix. Packagist's most installed wix-something is
			// wixel/gump, a PHP input validator with over a million
			// downloads, and wixiweb, wixnit and wixet are three more vendors
			// whose names begin with it. That is Lago's shape at a scale
			// Lago never reached.
			//
			// And the vendor's own scope is mostly not this API.
			// @wix/design-system, wix-style-react, wix-ui-core, wix-ui-tpa
			// and @wix/wix-ui-icons-common are a React component library with
			// hundreds of thousands of installs; @wix-pilot/core and
			// @wix-pilot/detox are an LLM-driven test runner from
			// wix-incubator with over a million; @wix/wix-code-types is
			// types for Velo, the language sites are scripted in. The single
			// most installed package with the word in its name, on either
			// registry, is a testing tool.
			//
			// Then the shapes already recorded: lezhnev74/laravel-socialite
			// -wix signs people in rather than calling the store, which is
			// what unikapps/laravel-socialite-squarespace and
			// socialiteproviders/gumroad were. epicformbuilder/wixhive-php-api
			// speaks WixHive, Wix's retired contacts API, which is Shopware 5
			// and the SOAP allegro package. And the bare npm "wix" is Wix's
			// App Market SDK published by an individual, not a REST client.
			//
			// wix/framework-bundle and wix/framework-component are Wix's own
			// Symfony packages for "decoding Wix instances" -- the vendor
			// name is genuinely theirs and the packages are internal
			// plumbing.
			recipe: "wix",
			composer: []string{
				"chkltlabs/wix-client",
				"storessuite/wix",
				"storessuite/wix-connect",
			},
			npm: []string{
				"@wix/sdk",
				"@wix/ecom",
				"@wix/api-client",
				"@wix/headless-ecom",
			},
		},
		{
			// Four exclusions, and the biggest of them is the vendor's own
			// scope. Everything under @squarespace/ is template tooling for
			// building Squarespace sites -- template-engine, toolbelt,
			// server, polyfills, less-ts-cli -- and @squarespace/core calls
			// itself "The frontend JS API for Squarespace templates", which
			// is a package with the word API in its description that is not
			// this API. eslint-config-squarespace is the company's own lint
			// rules. That is Etsy's phan and Akeneo's three, on a larger
			// scale: the provider publishes a great deal and none of it
			// calls the Commerce API.
			//
			// unikapps/laravel-socialite-squarespace signs people in rather
			// than talking to the store, which is what socialiteproviders
			// /gumroad and mugnate/oauth2-ecwid were.
			//
			// beloop/squarespace and beloop/squarespace-bundle are components
			// of the Beloop LMS that happen to carry the word. They require
			// beloop/core, beloop/course and beloop/user, and nothing in
			// either one says Squarespace anywhere else.
			//
			// And @tryghost/mg-squarespace-xml, seed-squarespace and
			// is-squarespace never call the provider at all: the first two
			// read and write export files to move a site between platforms,
			// and the third looks at a page and guesses what built it. A
			// project holding any of them is doing something to Squarespace
			// content with no API key existing anywhere.
			//
			// What is left is four small npm packages and, on Packagist,
			// nothing whatsoever. Squarespace has no PHP client for its
			// Commerce API on the registry, which is a fact about the
			// provider rather than a gap in the search.
			recipe: "squarespace",
			npm: []string{
				"squarespace-node-api",
				"marktera-squarespace",
				"@corte-so/commerce-squarespace",
				"@florasync/squarespace-mcp",
			},
		},
		{
			// Two exclusions, both shapes already recorded here, arriving
			// together on one small provider.
			//
			// mugnate/oauth2-ecwid is a login provider for the PHP League --
			// signing in with an Ecwid account rather than calling the store
			// API, which is what socialiteproviders/gumroad was.
			//
			// And npm's bare "ecwid" describes itself as "ecwid" and lists no
			// repository, which is exactly the state the backlog records
			// against Marqeta: a package of that name exists and cannot be
			// verified, so it is left alone rather than guessed at.
			//
			// What is mapped is small -- three-figure install counts at best.
			// Ecwid has no widely installed client on either registry, which
			// is a fact about the provider rather than a gap in the search.
			recipe: "ecwid",
			composer: []string{
				"dspacelabs/ecwid",
				"dspacelabs/ecwid-client",
				"strappberry/ecwid-api-wrapper",
			},
			npm: []string{"ecwid-api"},
		},
		{
			// The exclusion here is "sign in with" rather than "talk to".
			// socialiteproviders/gumroad is a Laravel Socialite provider and
			// alofoxx/oauth2-gumroad is the same idea for the PHP League: a
			// project holding either one lets people log in with a Gumroad
			// account and may never touch a product or a licence. That is the
			// shape Shopware's meteor-admin-sdk and Saleor's app-sdk have --
			// an authentication or extension package rather than a client.
			//
			// Worth noting for once: npm's gumroad really is the client
			// ("API client for Gumroad."), where the same bare word turned
			// out to be an algorithms library for Lago and an express
			// boilerplate for Polar. A collision on the plain name is common
			// enough to check for and not common enough to assume.
			//
			// Every mapped package here is small. Gumroad has no widely
			// installed client on either registry, which is a fact about the
			// provider rather than a shortfall in the search.
			recipe: "gumroad",
			composer: []string{
				"serenity_technologies/gumroad",
				"serenity_technologies/cashier-gumroad",
				"diskopete/laravel-gumroad",
				"flowframe/laravel-gumroad",
			},
			npm: []string{"gumroad", "gumroad-api"},
		},
		{
			// Prefix containment, which Lago met first with Lagoon and which
			// is now common enough to be a shape rather than a coincidence.
			// "polar" is the start of "polaris" and of "polarising", and both
			// of those out-install the real client on Packagist:
			// eduard9969/blade-polaris-icons is Shopify's design system at
			// eighty-nine thousand, and polarising/bcrypt is a password
			// hashing library at sixty-eight.
			//
			// The Shopify one is worth naming twice over. A prefix rule would
			// offer a payments emulator to a project using Shopify's icons,
			// while Shopify itself has two Recipes here -- so the wrong
			// answer would be wrong about the vendor as well as the domain.
			//
			// npm's polar is the ordinary sort: an express boilerplate.
			recipe:   "polar",
			composer: []string{"polar-sh/sdk", "danestves/laravel-polar"},
			npm:      []string{"@polar-sh/sdk", "@polar-sh/nextjs"},
		},
		{
			// A fresh shape of near miss: this one is a substring of a
			// bigger project's name. Lagoon is amazee.io's open-source
			// Kubernetes application delivery platform, and "lago" is inside
			// it -- so uselagoon/lagoon-php-sdk, steveworley/lagoon-php-sdk,
			// amazeeio/lagoon-logs, salsadigitalauorg/wp_lagoon_logs and
			// bryangruneberg/laragoon all match a substring rule. A usage
			// billing emulator offered to every Drupal site deployed on
			// amazee.io is not a near miss, it is a different industry.
			//
			// The npm package called lago is the ordinary sort: a data
			// structures and algorithms library.
			//
			// What is left are two clients. lago-javascript-client is the
			// official JavaScript one, and g-portal/lago-php-client borrows
			// Lago's own API blurb for its description, which is a better
			// signal than its four thousand installs.
			recipe:   "lago",
			composer: []string{"g-portal/lago-php-client"},
			npm:      []string{"lago-javascript-client"},
		},
		{
			// A dependency on Apideck proves the caller uses Apideck and not
			// which of its unified APIs. One SDK covers ecommerce, CRM,
			// accounting, HRIS and the rest, so a project holding it may be
			// doing none of the shopping this Recipe describes. It is mapped
			// anyway, because offering the only Apideck Recipe there is beats
			// offering nothing, and the imprecision is recorded here rather
			// than hidden in a match.
			//
			// The near miss is one letter. appdeck publishes sql, sampa and
			// filter -- a PDO wrapper, an MVC framework and an input filter --
			// and none of the three has anything to do with this. Voucherify
			// had Voucherly two letters away; this has appdeck at one, three
			// times over.
			//
			// muammarsiddiqui/apideck-laravel is the ordinary sort: an API
			// playground for Laravel that took the name.
			recipe:   "apideck",
			composer: []string{"apideck-libraries/php-sdk", "apideck-libraries/sdk-php"},
			npm:      []string{"@apideck/unify", "@apideck/node", "apideck"},
		},
		{
			// The near miss here is a different company rather than a
			// different word. voucherly/voucherly-php-sdk belongs to
			// Voucherly, an Italian payments business, and the name differs
			// from Voucherify by two letters -- so a prefix or substring rule
			// matches it and a reader skims past it. It is excluded by name.
			//
			// fastwebmedia/laravel-vouchering is the ordinary kind: a
			// vouchering system for Laravel that has nothing to do with this
			// provider, matching on the word rather than the vendor.
			//
			// rspective/voucherify is mapped although the vendor prefix is
			// not the product name: rspective is Voucherify's parent company,
			// and the package's own description is "Voucherify promotion
			// engine REST API". Half a million installs of exactly these
			// requests.
			recipe:   "voucherify",
			composer: []string{"rspective/voucherify"},
			npm:      []string{"@voucherify/sdk", "voucherify"},
		},
		{
			// The vendor-prefix rule fails hardest here, and it is worth
			// recording how hard. Four packages under the akeneo namespace
			// have seven-figure install counts and only one of them is a
			// client: akeneo/php-coupling-detector is an architecture linter
			// at 1.1 million, akeneo/batch-bundle is part of the PIM's own
			// Symfony stack at 1.06 million, and
			// akeneo/phpspec-skip-example-extension is a test-framework
			// extension at 997 thousand. akeneo-labs/spreadsheet-parser reads
			// XLSX files.
			//
			// This is etsy/phan and allegro/php-protobuf a third time, and
			// three instances is enough to say the shape rather than the
			// example: a vendor that publishes prolifically makes its own
			// prefix a poor filter, because the packages carrying its name
			// are mostly not about its product.
			//
			// akeneo/pim-community-dev and -standard are excluded for the
			// other reason. They are the PIM, so a project holding them
			// serves this API rather than calls it.
			//
			// What is left are the client and the connectors that use it:
			// Akeneo's own Adobe Commerce connector, two Sylius plugins, a
			// Laravel client and a Magento bundle. All of them make these
			// requests on somebody's catalogue.
			recipe: "akeneo",
			composer: []string{
				"akeneo/api-php-client",
				"akeneo/module-magento2-connector-community",
				"webgriffe/sylius-akeneo-plugin",
				"synolia/sylius-akeneo-plugin",
				"justbetter/laravel-akeneo-client",
				"justbetter/magento2-akeneo-bundle",
			},
		},
		{
			// Two exclusions, and the collection has seen both shapes before
			// at smaller scale.
			//
			// allegro/php-protobuf is the most-installed match by a wide
			// margin -- seventy-six thousand -- and it is Google's Protocol
			// Buffers for PHP, published from github.com/allegro/php-protobuf
			// by the same company's engineering organisation. The vendor
			// prefix is real and the package has nothing to do with the
			// marketplace, which is etsy/phan and printful/php-gettext-cms
			// again with a bigger number on it.
			//
			// The npm package called allegro is the other kind. Its
			// description says "Allegro.pl WebAPI client", its keywords
			// include soap, and it depends on soap-js: it is a client for the
			// legacy SOAP interface, not for the REST API this Recipe
			// describes. That is the pattern eBay's Trading SDKs, DHL's own
			// SOAP client and avalara/avatax already established -- when a
			// provider moves protocol the old client keeps its name and its
			// install base, so popularity is evidence of age rather than of
			// correctness.
			//
			// asocial-media/allegro-api is mapped although it covers both
			// interfaces, because a project holding it does make these
			// requests.
			recipe: "allegro",
			composer: []string{
				"imper86/php-allegro-api",
				"wiatrogon/allegro-rest-api",
				"asocial-media/allegro-api",
				"sebastianpozoga/php-allegroapi",
				"ircykk/allegro-api",
				"zoondo/allegro-api",
				"macopedia/magento2-allegro",
			},
		},
		{
			// Four packages under one scope and only two of them talk to the
			// API. @saleor/sdk is the client and @saleor/auth-sdk is the half
			// of it that gets a token, so both belong here.
			//
			// npm "saleor" is the CLI -- its own readme leads with a
			// saleor-cli logo -- and a command-line tool creates and deploys
			// projects rather than calling a storefront. @saleor/macaw-ui is
			// "Saleor's UI component library", which is React components and
			// no HTTP at all.
			//
			// @saleor/app-sdk is the interesting exclusion. It is for
			// building Saleor Apps: code that Saleor calls rather than code
			// that calls Saleor, built around webhook manifests, request
			// verification and the app's own registration handshake. An app
			// does eventually query the API, but it does so with a token
			// Saleor hands it during install, through whatever client it
			// chooses -- so the dependency says "this is an extension" rather
			// than "this makes these requests".
			recipe: "saleor",
			npm:    []string{"@saleor/sdk", "@saleor/auth-sdk"},
		},
		{
			// VTEX is the first provider here with no widely installed
			// external client, and the reason is worth recording rather than
			// working around.
			//
			// Its popular packages are all the IO platform's own toolchain.
			// npm "vtex" is the Toolbelt, a command-line tool for building
			// and publishing IO apps -- a build tool, not a caller of
			// anything. @vtex/api and @vtex/clients are the IO runtime's
			// clients, and an app using them runs inside VTEX's own
			// infrastructure and is handed a platform token: it never sends
			// X-VTEX-API-AppKey or X-VTEX-API-AppToken, which is the whole
			// credential this Recipe describes. Mapping them would offer an
			// emulator of the external integration path to code that does not
			// take it.
			//
			// What is left are small. brandlive/vtex-api and the Laravel
			// wrappers have four-figure install counts between them, and the
			// npm one has fewer. They are mapped anyway, because a small
			// package making exactly these requests is a better signal than a
			// large one making different ones -- and because the alternative
			// is mapping nothing at all for a platform that runs a
			// significant share of Latin American retail.
			recipe: "vtex",
			composer: []string{
				"brandlive/vtex-api",
				"juanfeltrin/vtex-sdk-php",
				"daygarcia/laravel-vtex",
				"grebo87/laravel-vtex-api",
			},
			npm: []string{"vtex-api"},
		},
		{
			// A third shape of vendor exclusion, after Printful's unrelated
			// tool and Shopware's previous generation: here the vendor ships
			// two SDKs for two different APIs under one npm scope.
			// @commercetools/importapi-sdk is the Import API -- a separate
			// service on a separate host, built around import containers and
			// asynchronous import operations, with no version on a resource
			// and no update actions. A project holding it is not calling this
			// API at all, and the only thing that says so is the middle of
			// the package name.
			//
			// The Merchant Center toolkits are excluded for the ordinary
			// reason. @commercetools-frontend/application-shell and the
			// @commercetools-uikit packages are React components for building
			// screens inside commercetools' own back office, so they are a
			// user interface rather than a caller of the HTTP API.
			//
			// commercetools/php-sdk is deprecated and mapped anyway. Its
			// deprecation notice does not stop three hundred thousand
			// installs from making exactly these requests, and a Recipe is
			// useful to whoever is maintaining that code precisely because
			// nobody is rewriting it this week.
			recipe: "commercetools",
			composer: []string{
				"commercetools/commercetools-sdk",
				"commercetools/php-sdk",
				"commercetools/symfony-bundle",
				"bestit/commercetools-odm",
			},
			npm: []string{
				"@commercetools/platform-sdk",
				"@commercetools/sdk-client-v2",
				"@commercetools/ts-client",
			},
		},
		{
			// The engine-is-not-a-client exclusion Medusa already needed,
			// at a much larger scale: shopware/core, shopware/storefront,
			// shopware/administration and shopware/platform run between one
			// and six million installs each, and every one of them is the
			// shop rather than a caller of one.
			//
			// shopware/shopware is a new shape and the sharper exclusion,
			// because it is the same vendor and the same word. It is
			// Shopware 5 -- a different product with a different API
			// entirely: /api/articles, basic auth, no sales channel and no
			// sw-access-key anywhere in it. Three quarters of a million
			// installs of a shop this Recipe describes nothing about, and
			// the only thing separating it from a match is the second half
			// of the package name.
			//
			// @shopware-ag/meteor-admin-sdk is excluded for the third
			// reason: it is an SDK for extensions that run inside the
			// administration's own interface, so it talks to a Vue
			// application rather than over HTTP to the Store API.
			//
			// What is left are the four that really are clients.
			// vin-sw/shopware-sdk is the PHP one, @shopware/api-client the
			// official TypeScript one, @shopware-pwa/api-client its
			// deprecated predecessor -- still installed, still calling the
			// same endpoints -- and @shopware/api-gen the generator that
			// exists only to build a client for this API.
			recipe:   "shopware",
			composer: []string{"vin-sw/shopware-sdk"},
			npm: []string{
				"@shopware/api-client",
				"@shopware-pwa/api-client",
				"@shopware/api-gen",
				"shopware-api-client",
			},
		},
		{
			// A new reason to exclude something: @medusajs/medusa is not a
			// client, it is the commerce engine. A project depending on it is
			// the store rather than a caller of one, and emulating the thing
			// you are running is not useful -- the store's own tests want a
			// database, not a fake of itself.
			//
			// @medusajs/js-sdk is the client, and it is the only mapping.
			recipe: "medusa",
			npm:    []string{"@medusajs/js-sdk"},
		},
		{
			// printful/php-gettext-cms is Printful's, and it is a translation
			// management backend rather than an API client -- the same shape
			// as etsy/phan, and with a six-figure install count of its own. A
			// vendor prefix would offer a print-on-demand emulator to a
			// project managing its own translations.
			recipe:   "printful",
			composer: []string{"printful/php-api-sdk", "fwartner/printful"},
			npm:      []string{"printful", "printful-api"},
		},
		{
			// Lightspeed is four unrelated products sharing a brand, not two.
			// R-Series is retail, K-Series is restaurants, X-Series is the
			// former Vend, and eCom is the former SEOshop -- four APIs, four
			// credentials, four vocabularies. This Recipe is R-Series, so only
			// the R-Series client is mapped and the other three are excluded by
			// name rather than by hoping a prefix rule sorts them out.
			//
			// The two nearest misses are not Lightspeed at all: the unscoped npm
			// package is the speed of light as a constant, and Packagist's
			// most-installed match is a Magento performance module that borrowed
			// the word.
			recipe:   "lightspeed",
			composer: []string{"timothydc/laravel-lightspeed-retail-api"},
		},
		{
			// On Packagist, "clover" is a code-coverage report format. It is
			// what PHPUnit emits, so nette/tester, phpunit-coverage-check and
			// half a dozen others carry the word and have tens of millions of
			// installs between them -- and not one of them is the payments
			// company. There is no PHP mapping here for that reason.
			//
			// The unscoped npm package is a toy language. remote-pay-cloud is
			// Clover's, and is the device SDK rather than a REST client: it
			// drives a card terminal over a socket and never calls the API
			// this Recipe models.
			recipe: "clover",
			npm:    []string{"clover-ecomm-sdk"},
		},
		{
			// Small on both registries -- the PHP clients are one author's and
			// have four figures of installs between them -- but unambiguous,
			// which is the bar. Nothing else on either registry is called
			// easyship, so this is the rare provider whose name carries the
			// information its packages do.
			recipe:   "easyship",
			composer: []string{"jrebs/easyship-php", "jrebs/easyship-laravel", "tigusigalpa/easyship-php"},
			npm:      []string{"easyship"},
		},
		{
			// The most-installed FedEx package is the wrong protocol, for the
			// fourth time in this table -- after eBay's Trading SDKs, DHL's own
			// SOAP client and avalara/avatax. php-fedex-api-wrapper has over a
			// million installs and wraps the SOAP web services this REST API
			// replaced.
			//
			// Both mapped packages say REST in their own descriptions. No npm
			// client was found under any obvious name.
			recipe:   "fedex",
			composer: []string{"whatarmy/fedex-rest", "shipstream/fedex-rest-sdk"},
		},
		{
			// One vendor, two protocols, and the names do not say which is
			// which. avalara/avatax describes itself as the AvaTax SOAP client
			// and avalara/avataxclient is the REST v2 one -- checked in their
			// manifests, where the first has no HTTP dependency at all and the
			// second pulls in Guzzle. This Recipe models REST v2.
			//
			// avatax-magento is a Magento module rather than a client, and
			// avalara-for-communications is a different Avalara product with a
			// different API, so neither is mapped here.
			recipe:   "avalara",
			composer: []string{"avalara/avataxclient"},
			npm:      []string{"avatax"},
		},
		{
			// The most-installed package under the lemonsqueezy vendor name is
			// a UI component library, not a client -- plain-ui-components has
			// twice the installs of the Laravel package that actually talks
			// to the API. Same shape as etsy/phan: the vendor publishes more
			// than its SDK, and a prefix rule would offer a commerce emulator
			// to a project that only borrowed some Blade components.
			recipe:   "lemonsqueezy",
			composer: []string{"lemonsqueezy/laravel", "seisigmasrl/lemonsqueezy-php"},
			npm:      []string{"@lemonsqueezy/lemonsqueezy.js"},
		},
		{
			// "recharge" is one of the most overloaded words on either
			// registry. The unscoped npm package is a React static site
			// generator; Packagist's matches are an EV charging SDK, an
			// electricity meter service and a Zimbabwean airtime top-up
			// library. Subscription billing is what the word means fourth or
			// fifth, so nothing here can be inferred from the name.
			//
			// @rechargeapps/storefront-client is left out for the reason DHL
			// Parcel is: it is Recharge, and it is a different API. The
			// storefront client talks to the customer portal with a session
			// token, and this Recipe models the admin API behind a store key.
			recipe:   "recharge",
			composer: []string{"huel-global/recharge"},
			npm:      []string{"recharge-api"},
		},
		{
			// DHL's own SDK is the wrong one for this Recipe. dhl/sdk-api-express
			// requires ext-soap; this Recipe models the REST tracking API, and
			// a SOAP client pointed at it gets a 404 on its endpoint. So does
			// dreipunktnull/dhl-express, and alfallouji/dhl_api and
			// mtcmedia/dhl-api speak the older XML services.
			//
			// dhlparcel-php-api is excluded for a different reason: DHL Parcel
			// is a separate API from DHL Express with its own credentials and
			// its own shapes, and this Recipe says in as many words that it
			// models Express and nothing else.
			recipe:   "dhl",
			composer: []string{"sonnenglas/mydhl-php-sdk", "korbeil/dhl-express-php-api"},
		},
		{
			// Four vendors publish a package called shipengine on Packagist,
			// including ShipEngine's own, and all four are API clients. Same
			// shape as ship-station below: the vendor half of the name says
			// who wrote the wrapper rather than what it talks to.
			recipe:   "shipengine",
			composer: []string{"shipengine/shipengine", "always-open/shipengine", "jsamhall/shipengine"},
			npm:      []string{"shipengine"},
		},
		{
			// Three unrelated vendors publish a package called ship-station on
			// Packagist and all three are real API wrappers, so the vendor
			// prefix carries no information here and each has to be named.
			recipe:   "shipstation",
			composer: []string{"campo/laravel-shipstation", "dissolve/ship-station", "zack6849/ship-station"},
			npm:      []string{"node-shipstation"},
		},
		{
			// aftership-apps-magento2 is a Magento storefront extension rather
			// than a client: a shop installing it talks to Magento, and
			// AfterShip talks to AfterShip. Naming the two SDKs explicitly
			// keeps a vendor-prefix rule from offering a tracking emulator to
			// a shop that never calls one.
			recipe:   "aftership",
			composer: []string{"aftership/aftership-php-sdk", "aftership/tracking-sdk"},
			npm:      []string{"aftership"},
		},
		{
			// The unscoped npm package called easypost is not this EasyPost.
			// It reads POST bodies from form submissions in Node, has nothing
			// to do with shipping, and is exactly the sort of name collision
			// a substring rule would hand a shipping emulator to. The scoped
			// @easypost/api is the carrier client.
			recipe:   "easypost",
			composer: []string{"easypost/easypost-php"},
			npm:      []string{"@easypost/api"},
		},
		{
			// The SP-API replaced Marketplace Web Service, and MWS was a
			// different API entirely -- XML over a signed query string, with
			// its own endpoints and its own vocabulary. Packages named after
			// it are still installed and still named amazon, so they are
			// named here to be excluded rather than matched by a rule that
			// looks for the word.
			recipe: "amazonsp",
			composer: []string{
				"jlevers/selling-partner-api",
				"amazon-php/sp-api-sdk",
				"amzn-spapi/sdk",
			},
			npm: []string{"amazon-sp-api"},
		},
		{
			// eBay's most-installed PHP package is the wrong one. dts/ebay-sdk-php
			// and its fork have well over a million installs between them and
			// they speak the XML Trading, Finding and Shopping APIs -- the ones
			// the REST APIs modelled here were built to replace. Mapping by
			// popularity would send almost every PHP eBay project to an
			// emulator that does not serve a single endpoint it calls.
			//
			// ebay/digital-signature-php-sdk is eBay's own, and it exists to
			// sign REST requests, so a project carrying it is on the REST side.
			recipe:   "ebay",
			composer: []string{"ebay/digital-signature-php-sdk"},
			npm:      []string{"ebay-api", "@hendt/ebay-api"},
		},
		{
			// Magento is the case where the same name covers three protocols.
			// magento-api on npm is SOAP, @magento/peregrine is PWA Studio
			// and talks GraphQL, and smalot/magento-client is SOAP v1 -- all
			// of them Magento, none of them the REST API this Recipe models.
			// So the table names only packages that say REST and say 2.
			recipe:   "magento",
			composer: []string{"springimport/swagger-magento2-client", "zero1/magento2-rest-client"},
			npm:      []string{"magento2-rest-client"},
		},
		{
			// Both of these say Etsy API v3 in as many words, which is the
			// bar here: the Recipe models v3, and the older wrappers on both
			// registries speak v2. Pointing a v2 project at a v3 emulator
			// would answer every call with a shape its client cannot read.
			recipe:   "etsy",
			composer: []string{"rhysnhall/etsy-php-sdk"},
			npm:      []string{"etsy-ts"},
		},
		{
			recipe:   "github",
			composer: []string{"knplabs/github-api"},
			npm:      []string{"@octokit/rest", "@octokit/core", "octokit"},
			gomod:    []string{"github.com/google/go-github"},
		},
		{
			recipe:   "slack",
			composer: []string{"slack-php/slack-php-api"},
			npm:      []string{"@slack/web-api", "@slack/bolt"},
			gomod:    []string{"github.com/slack-go/slack"},
		},
		{
			recipe:   "twilio",
			composer: []string{"twilio/sdk"},
			npm:      []string{"twilio"},
			gomod:    []string{"github.com/twilio/twilio-go"},
		},
		{
			recipe:   "adyen",
			composer: []string{"adyen/php-api-library"},
			npm:      []string{"@adyen/api-library"},
		},
		{
			recipe:   "clerk",
			composer: []string{"clerk/clerk-sdk-php"},
			npm:      []string{"@clerk/clerk-sdk-node", "@clerk/backend"},
		},
		{
			recipe:   "hubspot",
			composer: []string{"hubspot/api-client"},
			npm:      []string{"@hubspot/api-client"},
		},
		{
			recipe:   "xero",
			composer: []string{"xeroapi/xero-php-oauth2"},
			npm:      []string{"xero-node"},
		},
		{
			recipe:   "mailgun",
			composer: []string{"mailgun/mailgun-php"},
			npm:      []string{"mailgun.js"},
		},
		{
			recipe:   "sendgrid",
			composer: []string{"sendgrid/sendgrid"},
			npm:      []string{"@sendgrid/mail"},
		},
		{
			// The sharpest matched-on-a-semantic yet, and the third:
			// golang.org/x/mod implements module.EscapePath, which is the exact
			// rule this Recipe is about -- every capital to an exclamation mark
			// and its lowercase -- and it never opens a connection. It knows the
			// protocol better than any client here and does not speak it.
			//
			// After @snyk/ruby-semver and @snyk/nuget-semver, the pattern is
			// established: a provider's rules get reimplemented as a library far
			// more often than its API gets called, and the library is the more
			// popular package every time.
			//
			// The rest of the neighbourhood is the toolchain. go-mod-upgrade and
			// its relatives shell out to go list, which reaches the proxy
			// through cmd/go rather than over an API a test could point
			// somewhere else -- the line drawn at the Sonar scanners and NuGet's
			// exe wrappers. gomajor is mapped because it fetches @v/list from
			// GOPROXY itself.
			//
			// On npm, gosub-goproxy is a GOPROXY implementation: the product
			// rather than a client of it, which is where unleash-server and
			// unleash-frontend went.
			recipe: "goproxy",
			gomod:  []string{"github.com/icholy/gomajor"},
		},
		{
			// The word is an ordinary English one, and on Packagist it returns
			// three unrelated products. bariew/maven is described as "maven API"
			// and points at app.invoice-maven.co.il -- an Israeli invoicing
			// service. sukohi/maven, the most installed of the three, is a Laravel
			// package for managing an FAQ. Only augmentedlogic/maven-central-api
			// is about the registry, and it builds the exact URL the cases below
			// send: search.maven.org/solrsearch/select?q=g:...&core=gav&wt=json.
			//
			// Matched on a semantic, the fourth: mvn-artifact-name-parser is
			// "Parse java artifact names" and contains no URL at all. After
			// @snyk/ruby-semver, @snyk/nuget-semver and golang.org/x/mod, a
			// provider's coordinate format keeps getting reimplemented as a
			// library that never opens a connection.
			//
			// Its sibling mvn-artifact-download does open one, to
			// repo1.maven.org/maven2 -- the repository rather than the search API,
			// and a host this Recipe says it does not serve. maven-central-mcp
			// goes further and reads maven.google.com beside it, which is a
			// different registry. Neither is mapped.
			//
			// Four MCP servers carry the name; mvn-repository-mcp-server scrapes
			// mvnrepository.com, a third-party site rather than the vendor. And
			// the badge kind is here again: @ephemeris/core and
			// @dev.atsushieno/mugene match on img.shields.io/maven-central/v/ in
			// their READMEs.
			//
			// node-java-maven and snyk-mvn-plugin drive mvn itself, which is the
			// line already drawn at the Sonar scanners and NuGet's exe wrappers.
			recipe:   "mavencentral",
			composer: []string{"augmentedlogic/maven-central-api"},
			npm:      []string{"mavencc", "mvns", "maven-api-client"},
			gomod: []string{
				"github.com/tamnd/mavencentral-cli",
				"github.com/nscuro/cdx-central",
			},
		},
		{
			// A new near-miss kind, and the cleanest example of it: the vendor
			// ships its database as a file, so the clients that use it never call
			// the API. @renovatebot/osv-offline downloads a packed copy and queries
			// it locally; the tarball contains no osv.dev URL at all, and
			// @mintmaker/osv-offline is a fork of it. Nothing an emulator serves
			// would ever be reached.
			//
			// Matched on a semantic, the fifth: github.com/ossf/osv-schema is the
			// schema itself -- twenty-one Go files, none of which import net/http.
			// It is the definition of the format this Recipe describes and it does
			// not speak to anything. After @snyk/ruby-semver, @snyk/nuget-semver,
			// golang.org/x/mod and mvn-artifact-name-parser.
			//
			// On Packagist the three letters are somebody's name. osvaldoabel/acl,
			// osvenson/mailgun-api and osvaldovictor/validate-docs-ao all match on
			// a vendor namespace beginning "Osv", and rajangdavis/osvc_php is
			// Oracle Service Cloud, which abbreviates to OSvC. Only
			// gumslone/laravel-vulns is about this provider, and it is mapped: it
			// takes a configurable base_url, calls v1/query and v1/querybatch, and
			// its own comments record the stub shape the cases below pin.
			//
			// osv.dev is the module path of the vendor's Go binding, whose
			// BaseHostURL is configurable and whose three endpoint constants are
			// exactly the three routes served here. osv-to-sarif reads
			// osv-scanner's output file rather than the API -- the file-writing
			// reporter kind, from the other side -- and osv-scanner-mcp wraps the
			// binary, which is where the Sonar scanners and NuGet's exe wrappers
			// went.
			recipe:   "osvdev",
			composer: []string{"gumslone/laravel-vulns"},
			npm: []string{
				"yarn-osv-audit",
				"@bun-security-scanner/osv",
				"@absolutejs/vulnerabilities-osv",
				"osv-ui",
			},
			gomod: []string{"osv.dev", "github.com/google/osv-scanner"},
		},
		{
			// Matched on a semantic for the sixth time, and the first time inside
			// the mapped module. deps.dev is the module path of Google's own Go
			// code, and it holds both deps.dev/api/v3 -- the API client -- and
			// deps.dev/util/semver and deps.dev/util/resolve, which compare
			// versions and resolve npm, PyPI and Maven manifests locally and
			// import net/http nowhere. A go.mod names the module rather than the
			// package, so the mapping is right and the reason to write it down is
			// that two of the module's most useful packages never open a
			// connection. After @snyk/ruby-semver, @snyk/nuget-semver,
			// golang.org/x/mod, mvn-artifact-name-parser and ossf/osv-schema.
			//
			// Packagist has no deps.dev client at all, the way it had none for
			// Hookdeck.
			//
			// On npm the three mapped packages all build the same base,
			// https://api.deps.dev/v3, and @corvalon/lichen documents the three
			// route shapes this Recipe serves in its own comments.
			// @agntn/registries calls itself a universal package registry client
			// and names this provider in its description, and its tarball contains
			// no deps.dev URL -- the description names it, the code does not call
			// it. @pipeworx/mcp-deps-dev is agent tooling from the same publisher
			// as @pipeworx/mcp-goproxy, which is where the Go proxy's near misses
			// ended, and the word "deps" alone returns tailwind configurations.
			recipe: "depsdev",
			npm: []string{
				"@toiroakr/argent",
				"@eastagile/dephold",
				"@corvalon/lichen",
			},
			gomod: []string{"deps.dev"},
		},
		{
			// The sharpest near miss in the collection, and a kind that has not
			// appeared before: the vendor's own client does not speak the vendor's
			// documented JSON API. The npm package openmeteo is Open-Meteo's, and
			// it sets format=flatbuffers on every request unconditionally, so it
			// receives a binary payload rather than the document this Recipe
			// describes. @openmeteo/sdk is nothing but the generated FlatBuffers
			// schema for it. Neither is mapped, because an emulator answering JSON
			// cannot answer either of them.
			//
			// @openmeteo/weather-map-layer and go-omfiles are a third wire format
			// again, the OMfile tiles. After @renovatebot/osv-offline, where the
			// database ships as a file, the pattern is that a provider's own
			// tooling is the least likely thing in its neighbourhood to call its
			// public API.
			//
			// Three MCP servers carry the name, one of them @pipeworx/mcp-open-meteo
			// -- the third Recipe running for that publisher, after mcp-goproxy and
			// mcp-deps-dev.
			//
			// Mapped: the clients that ask for JSON. symfony/ai-open-meteo-tool and
			// php-weather/open-meteo both request api.open-meteo.com/v1/forecast
			// directly, and omgo keeps its forecast, historical and ensemble URLs
			// in fields a caller can replace.
			recipe: "openmeteo",
			composer: []string{
				"symfony/ai-open-meteo-tool",
				"php-weather/open-meteo",
				"flibidi67/open-meteo",
			},
			npm: []string{
				"@agentic/open-meteo",
				"@molecule/api-weather-open-meteo",
			},
			gomod: []string{
				"github.com/hectormalot/omgo",
				"github.com/tpryan/openmeteogo",
			},
		},
		{
			// The same host serves three URL families and only one of them is this
			// API. dat-usgs-earthquakes reads
			// /earthquakes/feed/v1.0/summary/all_hour.geojson, the pre-baked
			// summary feeds: identical GeoJSON, no query parameters, and because
			// there is no limit to send they always carry metadata.count -- so the
			// behaviour this Recipe leads with cannot happen there at all.
			// usgs-earthquake-reports parses the pages under /earthquakes/.
			// Neither is mapped.
			//
			// Five MCP servers carry the name, the most for any Recipe here, and
			// one is @pipeworx/mcp-usgs-earthquakes -- the fourth Recipe running
			// for that publisher after goproxy, deps.dev and Open-Meteo.
			//
			// And a shape worth recording rather than mapping: four npm packages
			// carry the identical description "weather and gis aggregation library
			// with a simple API" under unrelated names, one of them the name of a
			// commercial device-detection product. That one does call this API,
			// from a file called quakes.js, beside an install-time script. A
			// package whose name, description and code refer to three different
			// things is not evidence that a project meant to reach this provider,
			// which is what detection is for, so none of the four is mapped.
			recipe:   "usgs",
			composer: []string{"okierazorback/usgs"},
			npm:      []string{"usgs-earthquake-api"},
			gomod:    []string{"github.com/jasonmoo/usgs"},
		},
		{
			// Two hosts, two version prefixes, and packages written for a version
			// that is not there. api.frankfurter.app 301s every path to
			// api.frankfurter.dev with a /v1 prefix added, and the unversioned path
			// on the new host is a 404 -- so frankfurter-cli, which hardcodes the
			// old host, works only because its HTTP client follows redirects, and a
			// PHP cURL client with the defaults would get an HTML redirect page
			// where it expected JSON. Meanwhile @pontx/frankfurter-v2 calls itself
			// an SDK for "the official Frankfurter v2 API" and go-finance defaults
			// to api.frankfurter.dev/v2/, and there is no v2: every path under it
			// answers 404. go-finance is mapped anyway, because its base URL is
			// configurable and the project plainly means to reach this provider.
			//
			// frankfurter-api-status-client is named after the provider and calls
			// frankfurter.instatus.com, a hosted status-page vendor rather than the
			// provider's own host at all. That is a new one: not the vendor's other
			// product, but somebody else's product about the vendor.
			//
			// On Packagist the word is a place and a surname. The service is named
			// after the city, so chuckcms-template-frankfurt is a CMS theme, and
			// frankforte/quantumphp is a JavaScript console logger under a vendor
			// namespace that merely begins the same way. And @pipeworx/mcp-frankfurter
			// is the fifth Recipe running for that publisher.
			recipe:   "frankfurter",
			composer: []string{"investbrainapp/frankfurter-client"},
			npm:      []string{"frankfurter-js", "frankfurter-api-client"},
			gomod: []string{
				"github.com/tdrn-org/go-finance",
				"github.com/tamnd/frankfurter-cli",
			},
		},
		{
			// Nominatim is software as well as a service, so the neighbourhood has
			// three kinds of thing that are not clients of the public instance.
			// @mailwoman/nominatim is "a Nominatim-compatible HTTP geocoding API"
			// -- it serves /search and /reverse and its tarball holds no
			// nominatim.openstreetmap.org at all, which is the product rather than
			// a client of it, where gosub-goproxy and unleash-server went.
			//
			// geo-golang/mapquest/nominatim is the same software hosted by
			// MapQuest, behind open.mapquestapi.com and an API key -- so neither
			// the missing-User-Agent 403 nor the no-credential reading this Recipe
			// is built on applies there. A third party running the vendor's
			// software is a new kind, and it is not the vendor.
			//
			// codemacher/tile_proxy is the same project's other service: raster
			// tiles from tile.openstreetmap.org, images rather than JSON, which is
			// the vendor's-other-product kind again.
			//
			// Mapped: geocoder-php/nominatim-provider, at 4.9 million installs the
			// most used client of this API anywhere, and the rest that name the
			// public host outright.
			recipe: "nominatim",
			composer: []string{
				"geocoder-php/nominatim-provider",
				"maxh/php-nominatim",
			},
			npm:   []string{"nominatim-client", "nominatim-ts"},
			gomod: []string{"github.com/doppiogancio/go-nominatim"},
		},
		{
			// The name is the two commonest words in a package description, and on
			// Packagist that is the whole result. Searching "open library" returns
			// PHPWord, openspout, bootstrap-icons, geophp, casbin and nelexa/zip --
			// eight libraries that are open source and none of them about this
			// provider. Only searching the word closed up finds anything real. That
			// is the ordinary-word collision at its limit: not a homonym, just two
			// adjectives every library in the world applies to itself.
			//
			// openlibrary-scraper does reach openlibrary.org, at /isbn/ and
			// /search, and parses what comes back with cheerio -- the HTML site
			// rather than the API, which is the same host and a surface this
			// Recipe does not serve. Not mapped.
			//
			// alexandria-worker is mapped and is worth naming: it builds
			// openlibrary.org${authorRef.key}.json to turn an author reference into
			// a name, which is exactly the second request this Recipe's headline is
			// about, written out in somebody's production code.
			recipe: "openlibrary",
			composer: []string{
				"beezus/openlibrary-php",
				"shiranse/laravel-open-library",
			},
			npm: []string{
				"open-library-client",
				"openlibrary-cli",
				"alexandria-worker",
			},
			gomod: []string{"github.com/exlibris-fed/openlibrary-go"},
		},
		{
			// Two APIs on one host, and almost everything uses the other one.
			// MediaWiki's Action API lives at /w/api.php and predates the REST API
			// this Recipe describes by a decade, so node-wikifetch,
			// symfony/ai-wikipedia-tool and sophivorus/easy-wiki all reach
			// wikipedia.org and none of them reach anything served here. That is
			// the vendor's-other-product kind, except the other product is the
			// vendor's older API and it is still the popular one.
			//
			// wikifox is the nearest miss of all: it does call /api/rest_v1/, at
			// page/pdf/{title}, which answers a PDF. The right host, the right API
			// family, and a response an emulator of this endpoint cannot give.
			//
			// webignition/internet-media-type matched because its description
			// cites an en.wikipedia.org URL as the definition of the thing it
			// models -- matched on a link rather than on a dependency, which is
			// new. wikimedia/wikipedia-preview is the vendor's own embeddable
			// widget. And "wikipedia api" on npm returns color-name.
			//
			// Mapped: the two that name /api/rest_v1/ outright. Nothing on
			// Packagist or in Go modules does, which is why the coverage table
			// counts this Recipe as reachable from npm alone.
			recipe: "wikipedia",
			npm:    []string{"wikipedia", "wikipedia-summary"},
		},
		{
			// The generic name belongs to somebody else's API of the same data.
			// The npm package called hacker-news-api -- the name anyone would try
			// first -- is a client for hn.algolia.com, the search service Algolia
			// runs over Hacker News, and its tarball holds no firebaseio.com at
			// all. Different host, different shapes, same stories. That is the
			// MapQuest-Nominatim kind and sharper, because here the third party has
			// the obvious name and the vendor's own API does not.
			//
			// Matched on a semantic, the seventh: hacker-news-api-types is one
			// index.d.ts. It cites firebaseio URLs in its README and declares the
			// shapes this Recipe serves, and there is nothing in it to run. After
			// @snyk/ruby-semver, @snyk/nuget-semver, golang.org/x/mod,
			// mvn-artifact-name-parser, ossf/osv-schema and deps.dev/util/semver.
			//
			// hnpwa-api deploys a Hacker News API on your own domain: the product
			// rather than a client of it, where gosub-goproxy, unleash-server and
			// @mailwoman/nominatim went. Two MCP servers carry the name.
			//
			// Packagist has no client that reaches this API, the way it had none
			// for Hookdeck or deps.dev.
			recipe: "hackernews",
			npm:    []string{"node-hn-api", "hn-ts"},
			gomod:  []string{"github.com/hermanschaaf/hackernews"},
		},
		{
			// pokeapi-sprites is the offline-data kind at its largest: 4237 files,
			// every sprite shipped as an asset, and the only pokeapi.co URLs in it
			// are in the README. It is the same reading as @renovatebot/osv-offline
			// -- the data as a package, so there is no request to redirect -- and
			// it lands doubly here, because the sprite URLs this API returns point
			// at raw.githubusercontent.com and were never on pokeapi.co either.
			//
			// Matched on a semantic, the eighth: pokeapi-types and
			// @bgoff1/pokeapi-types are four files each, declaring the shapes this
			// Recipe serves, with the API's address appearing only in a comment.
			//
			// @professorragna/pokeapi-cli is mapped and is worth naming: its own
			// description is "CLI for the PokeAPI (pokemon and pokemon-species
			// endpoints)", which is exactly the two this Recipe models.
			recipe:   "pokeapi",
			composer: []string{"danrovito/pokephp", "lmerotta/phpokeapi"},
			npm: []string{
				"pokeapi-js-wrapper",
				"pokeapi-typescript",
				"@professorragna/pokeapi-cli",
			},
			gomod: []string{"github.com/mtslzr/pokeapi-go"},
		},
		{
			// The Recipe's headline is the near miss. npm has two clients for one
			// API, one per wire format: "musicbrainz" describes itself as "a
			// MusicBrainz XML Web Service Version 2 client" and parses with xml2js,
			// and "nodebrainz" as "a MusicBrainz JSON Web Service Version 2 client"
			// and sends fmt=json. Only the second is mapped. The first would in
			// fact be answered by the XML default this Recipe serves, for the one
			// path it holds, and that is not emulating an XML client: every
			// resource and every case here is JSON.
			//
			// dhowden/tag/mbz is the sharpest thing in the neighbourhood.
			// http://musicbrainz.org appears in it as a constant, and it is never
			// fetched: it is the namespace string inside an ID3 UFID frame, so the
			// package reads MusicBrainz identifiers out of audio files without ever
			// asking MusicBrainz anything. The provider's own URL, as data.
			//
			// graphbrainz is a GraphQL server over the same database -- somebody
			// else's interface, and a product rather than a client, both at once.
			// mikealmond/coverartarchive is the sibling service on
			// coverartarchive.org. go.uploadedlobster.com/mbtypes is types and
			// /discid computes a disc identifier locally. And music-tagger is
			// AcoustID, which is a different service again.
			recipe:   "musicbrainz",
			composer: []string{"mikealmond/musicbrainz", "lachlan-00/musicbrainz"},
			npm:      []string{"nodebrainz", "musicbrainz-api"},
			gomod:    []string{"github.com/michiwend/gomusicbrainz"},
		},
		{
			// The counter-example to Open-Meteo. There the vendor's own client
			// spoke a wire format the vendor's documented API does not use; here
			// the vendor publishes an official client in all three ecosystems --
			// @openfoodfacts/openfoodfacts-nodejs, openfoodfacts/openfoodfacts-php
			// and github.com/openfoodfacts/openfoodfacts-go -- and every one of
			// them calls the JSON API this Recipe serves. The nodejs one covers
			// v1, v2 and v3 of /product from a single package.
			//
			// Six MCP servers carry the name, which is the most for any Recipe in
			// this collection: USGS had five. One of them is
			// @pipeworx/mcp-openfoodfacts, the sixth Recipe running for that
			// publisher after the Go proxy, deps.dev, Open-Meteo, USGS and
			// Frankfurter.
			//
			// @molecule/api-nutrition-database is a "nutrition + barcode lookup
			// core interface" over several providers rather than a client of this
			// one, which is the shape already recorded at @agntn/registries.
			recipe: "openfoodfacts",
			composer: []string{
				"openfoodfacts/openfoodfacts-php",
				"openfoodfacts/openfoodfacts-laravel",
			},
			npm: []string{
				"@openfoodfacts/openfoodfacts-nodejs",
				"openfoodfacts-nodejs",
				"openfoodfacts-cli",
			},
			gomod: []string{"github.com/openfoodfacts/openfoodfacts-go"},
		},
		{
			// "Cross-reference" is a technical term in four unrelated fields, and
			// the search returns all of them. ttree/crossreference and
			// carbon/crossreference manage cross-references between content nodes
			// in Neos; weaviate's crossref package is a reference between objects
			// in a vector database; @paperist/remark-crossref does pandoc-style
			// cross-references in Markdown; and only this provider means the
			// bibliographic one. That is the ordinary-word collision with four
			// distinct technical senses rather than one word doing double duty.
			//
			// bsobbe/ithenticate is the neighbouring service: iThenticate is what
			// Crossref's Similarity Check runs on, so a project can be reaching
			// the one because of the other and still never call this API.
			//
			// Five MCP servers carry the name.
			recipe:   "crossref",
			composer: []string{"renanbr/crossref-client"},
			npm: []string{
				"crossref",
				"crossref-cli",
				"@jamesgopsill/crossref-client",
			},
			gomod: []string{
				"github.com/caltechlibrary/crossrefapi",
				"github.com/cgxeiji/crossref",
			},
		},
		{
			// The thinnest neighbourhood in this collection, and the funniest near
			// miss in it. haoteam/hippo-buyer-php-sdk describes itself as a
			// "hippopotamus buyer php sdk" and matches on -potam-, the Greek root
			// for river that both names happen to carry. A substring of an animal.
			//
			// @pipeworx/mcp-zippopotam is the seventh Recipe running for that
			// publisher, after the Go proxy, deps.dev, Open-Meteo, USGS,
			// Frankfurter and Open Food Facts.
			//
			// Mapped: two PHP wrappers, a Go package under a ziplookup module, and
			// @eusilvio/zip-lookup, which calls itself an agnostic US ZIP lookup
			// and reaches this host for the data.
			recipe: "zippopotam",
			composer: []string{
				"zippopotamus/zippopotamus",
				"th3fallen/zippopotamus-api",
			},
			npm:   []string{"@eusilvio/zip-lookup"},
			gomod: []string{"github.com/johnnypivot/ziplookup"},
		},
		{
			// The vendor's other API, behind the vendor's own paywall.
			// @datafire/tvmaze is an integration for the TVmaze *user* API, which
			// lives under a key, covers watchlists and follows, and is reached
			// through www.tvmaze.com/dashboard and /premium. This Recipe is the
			// public keyless one, and nothing an emulator of it serves would
			// answer that package.
			//
			// tvmaze-api-ts calls itself "a tvmaze scraper and api wrapper" and is
			// both: it reaches api.tvmaze.com and also parses
			// www.tvmaze.com/episodes pages with cheerio. Half of it is a client
			// of this API and half of it is a client of the website, which is the
			// line already drawn at openlibrary-scraper -- and there the package
			// was only ever the second half.
			//
			// No agent tooling carries this name, which after six MCP servers on
			// Open Food Facts and five on USGS is worth writing down.
			recipe: "tvmaze",
			composer: []string{
				"joshpinkney/tv-maze-php-api",
				"tavy315/tvmaze-api",
			},
			npm: []string{
				"tvmaze",
				"tvmaze-node",
				"node-tvmaze",
			},
			gomod: []string{"github.com/tamnd/tvmaze-cli"},
		},
		{
			// Matched on a semantic, the ninth -- and the first time it is the
			// entire neighbourhood rather than one package in it. The sun's
			// position is a closed-form calculation, so almost everything that
			// answers this search computes the answer instead of asking for it:
			// suncalc, suncalc2, sunrise-sunset-js and @hebcal/noaa on npm;
			// ypid/suncalc, tuxonice/suncalc-php, martindilling/sunny,
			// andrmoel/astronomy-bundle and minube/solar-position on Packagist;
			// and go-sunrise, kelvins/sunrisesunset, lukasschwab/sunrisesunset,
			// sixdouglas/suncalc and keep94/sunrise in Go. Every one of them is
			// trigonometry. Both Go packages named exactly after the provider
			// compute rather than fetch, and there is no Go or PHP client of this
			// API anywhere.
			//
			// That is the sharpest form of the pattern first recorded at
			// @snyk/ruby-semver: a provider's answer gets reimplemented far more
			// often than its API gets called. Here it is not a library or two but
			// the whole field, and the reason is that the API is the unusual
			// choice rather than the obvious one.
			//
			// Mapped: the two on npm that actually send a request.
			recipe: "sunrisesunset",
			npm:    []string{"sunrise-sunset-api", "n8n-nodes-sunrise-sunset"},
		},
		{
			// Two Recipes running, and the same reason both times: a VIN is
			// defined by ISO 3779 and its structure can be read without asking
			// anybody. sunrise/vin, at 152,000 installs, says so in its own
			// description -- "based on ISO-3779" -- and never opens a connection.
			// @cardog/corgi sells itself on the point: "fast, offline VIN
			// decoding", and its tarball holds no nhtsa.dot.gov at all. That is
			// matched-on-a-semantic and the offline-data kind at once.
			//
			// What the API adds is the part ISO does not define: the make, the
			// model and the body class behind a manufacturer's own codes. So the
			// packages that call it are the ones that want more than the standard,
			// and there are very few.
			//
			// And "decoder" is a word. The same search returns opus-decoder,
			// mpg123-decoder, ogg-opus-decoder, text-decoder and
			// @wasm-audio-decoders/opus-ml -- five packages about audio and text
			// streams, none of them about vehicles.
			recipe: "nhtsa",
			npm:    []string{"@shaggytools/nhtsa-api-wrapper"},
		},
		{
			// The obvious name is a placeholder. The npm package called "wikidata"
			// is two files and its own description reads "this is a placeholder for
			// an incoming package". Nothing has arrived. That is a kind not seen
			// before: the name a project would reach for first is reserved and
			// empty.
			//
			// Then the same split as Wikipedia and MusicBrainz, and for the same
			// reason. wikibase-sdk, wikibase-edit, wikibase-cli, nodemw and
			// wikidata-entity-lookup all reach w/api.php -- the Action API, which
			// predates this one. Only @wmde/wikibase-rest-api and
			// wikibase-rest-api-ts call rest.php/wikibase, and the first of those
			// is by the team that builds Wikibase.
			//
			// Packagist's whole first page is the vendor's own PHP value-object
			// model: wikibase/data-model, data-values/geo, data-values/number,
			// serialization/serialization and diff/diff at 1.7 million installs.
			// They are the types and parsers lifted out of the Wikibase software
			// and they open no connection at all -- matched on a semantic, at
			// scale, by the provider itself.
			//
			// And three more that ship the data rather than fetch it: wikidata-lang
			// and wikidata-person-names are datasets built from SPARQL, and
			// wikidata-filter processes the newline-delimited dump.
			recipe: "wikidata",
			npm:    []string{"@wmde/wikibase-rest-api", "wikibase-rest-api-ts"},
		},
		{
			// A new way to collide: not a shared word, but a shared edit distance.
			// Packagist's three highest-ranked results for "gbif" are
			// intervention/gif, a PHP GIF codec at 34.9 million installs;
			// mollie/mollie-api-php, a payments library at 17.2 million; and
			// griffinledingham/php-apple-signin at 2.2 million. Fifty-four million
			// installs between them and not one about biodiversity. The word does
			// not mean anything else -- the search engine's fuzziness is what
			// matched, which is a first here.
			//
			// gbif-checklistbank and gbif-namefinder reach api.gbif.org at
			// /dataset and /lookup/name_usage: older endpoint families on the same
			// host, and not the v1 species matcher. gbif-map draws occurrence
			// density from /v2/map/, which answers tiles rather than JSON. Neither
			// is mapped.
			//
			// Three MCP servers carry the name, one of them @pipeworx/mcp-gbif --
			// the eighth Recipe running for that publisher, after the Go proxy,
			// deps.dev, Open-Meteo, USGS, Frankfurter, Open Food Facts and
			// Zippopotam.
			recipe:   "gbif",
			composer: []string{"restelae/php-gbif"},
			npm:      []string{"gbif-crawler"},
			gomod:    []string{"github.com/HannesOberreiter/gbif-extinct"},
		},
		{
			// The Sunrise-Sunset shape again: a holiday calendar is a rule set, so
			// the neighbourhood computes rather than asks. date-holidays and
			// date-holidays-parser carry the rules for every country;
			// poland-public-holidays says "compute" in its own description;
			// @18f/us-federal-holidays, @holiday-jp/holiday_jp and
			// colombian-holidays each ship one country's. None of them opens a
			// connection.
			//
			// Two more are somebody else's answer to the same question: the npm
			// package called public-holidays reads Google Calendar, and
			// @festivo-io/festivo-sdk is a commercial competitor's client.
			//
			// And two are written for versions that are not there.
			// @pontx/nager-date calls itself a "Nager.Date Community API v4 SDK"
			// and /api/v4/ answers 404; @thomasc93/holidates builds
			// /api/v2/publicholidays/ and that answers 404 too. The same shape as
			// Frankfurter's v2, twice on one provider.
			//
			// On Packagist "nager" is an edit distance again: nagarajanchinnasamy,
			// socialiteproviders/naver, donejeh/nuban and mercyware/nigerianstates.
			recipe:   "nagerdate",
			composer: []string{"rolle-marketplace/nager-date-laravel"},
		},
		{
			// A new near miss, and the sharpest one yet: a Composer vendor
			// namespace that is the vendor and still is not the product.
			// internetarchive/virustotalapi is published by the Internet Archive
			// and wraps VirusTotal. Matching on the vendor prefix would map
			// somebody else's API to this Recipe with the vendor's own name as
			// the evidence.
			//
			// A different company's product has the same name. @esri/wayback-core,
			// @rapideditor/wayback-core and maplibre-gl-esri-wayback are all
			// Esri's World Imagery Wayback -- versions of satellite imagery, on
			// wayback.maptiles.arcgis.com. Nothing to do with archived web pages.
			//
			// And the bare name on npm belongs to a third thing again: wayback is
			// "Manage sequences of arbitrary data with revisions", and it opens no
			// connection at all. chrisullyott/wayback-cache is the same metaphor
			// on Packagist -- "Caches and catalogs a history of API requests".
			//
			// The rest of the neighbourhood calls the archive and calls something
			// else on it. Four other endpoints, none of them availability:
			//
			//   - CDX (/cdx/search/cdx): wayback-machine-downloader,
			//     baraja-core/wayback, odinns/laravel-wayback-machine,
			//     survos/site-discovery-bundle -- and every Go package that turned
			//     up, unanimously. gowaybackgo, fgau and blackhorn-modules all
			//     build /cdx/search/cdx, so there is nothing to map on gomod.
			//   - Save Page Now (/save/): mapkyca/known-wayback-machine,
			//     matthiasott/craft-internetarchive, inwebo/save-page-now-2.
			//   - Items and metadata (advancedsearch.php, /metadata/):
			//     internetarchive-sdk-js, internet-archive-api,
			//     votemike/internet-archive-api, inwebo/item-metadata.
			//   - Undocumented internals (/__wb/calendarcaptures): dawood/wmb-scrapper.
			//
			// The vendor's own npm org is nearly all user interface --
			// @internetarchive/bookreader, ia-icons, modal-manager, ia-dropdown,
			// shared-resize-observer -- and the one piece of it named wayback,
			// @internetarchive/ia-wayback-search, is a search *form component*
			// that builds web.archive.org/web/ links for a browser to follow.
			//
			// Two more of the standing kinds. strip-wayback-toolbar processes the
			// archive's output -- it strips the toolbar out of an already-archived
			// page and never connects. And the format is shipped as a file:
			// warcio, warc, node-warc, @webrecorder/wabac and @harvard-lil/js-wacz
			// read WARC and WACZ off disk.
			recipe:   "waybackmachine",
			composer: []string{"f2face/wayback-machine"},
			npm: []string{
				"wayback.js",
				"wayback-machine",
				"wayback-machine-fixed",
			},
		},
		{
			// A new near miss: the description belongs to a different package.
			// latlon-utils is published as "A JavaScript library for obtaining
			// data from weather.gov and tidesandcurrents.noaa.gov" -- word for
			// word noaajs's own description -- and its code contains no
			// weather.gov anywhere. Detection that reads registry prose would map
			// it on the strength of a sentence its author copied.
			//
			// The abbreviation collides three ways. npm's bare "nws" launches a
			// static web server and opens no connection at all;
			// namirasoft-nws-region and its siblings are one company's products,
			// served from region.nws.namirasoft.com; and vrtnws reads a Belgian
			// broadcaster's news feed, where the letters are Dutch.
			//
			// NOAA is much larger than this API, and most of the neighbourhood is
			// its other services: tides and currents (uriweb/uri-tides), space
			// weather, marine buoys, climate stations, NCAT coordinate
			// conversion, aerodrome METAR (ballen/metar), and GFS model files
			// pulled from nomads.ncep.noaa.gov (noaa-gfs-js).
			//
			// weather.gov itself answers on more hosts than this one, and each has
			// its own clients. alerts.weather.gov serves CAP in XML, which is what
			// thtroyer/noaa-cap-alerts, weather-alerts-parser and
			// weather-alerts-geojson read. forecast.weather.gov serves the human
			// site, and weather-gov-graph-parse scrapes its markup rather than
			// calling anything. aviationweather.gov is a different centre again.
			//
			// And the predecessors are still in the results. amwhalen/noaa targets
			// weather.gov/xml/current_obs and graphical.weather.gov's NDFD
			// service, which is what this API replaced; origindesign/noaa still
			// documents forecast-v3.weather.gov, the hostname the API carried
			// before it moved, though its code calls the current one and it maps.
			//
			// Two of the standing kinds close it out. @hebcal/noaa computes
			// sunrise and sunset from a NOAA algorithm and never connects, which
			// is the third holiday-calendar shape in three Recipes. And a cluster
			// of near-identical Anax modules -- serbseb/weathermodule,
			// lioo19/weathermodule, lingul/ramverk1-weathermodule,
			// bjos19/anax-weathermodule, bjorn/weathermod -- is one course's
			// coursework; the one opened here calls Dark Sky, ipstack and Google
			// Maps, and Dark Sky has been off for years.
			//
			// On Go, goa.design/clue/example/weather/.../weathergov and its twin
			// under goa.design/pulse are demo applications vendored inside a
			// framework's own repository. They are importable and they are not a
			// dependency anybody adds on purpose.
			recipe: "nationalweatherservice",
			composer: []string{
				"chrisreedio/laravel-weathergov-sdk",
				"benjaminhansen/nws-php",
				"projectsaturnstudios/nws",
				"almostengr/cakephp-nwsapi",
				"origindesign/noaa",
			},
			npm: []string{
				"weathered",
				"nws-client-js",
				"@vavassor/nws-client",
				"@bdking71/nws-weather-wrapper",
				"@wkronmiller/weather-api-client",
				"@paxperscientiam/national-weather-service-api.ts",
				"@oliverhodge2/weather-wrapper",
				"@naviaero/winds-aloft",
				"signalk-nws-alerts",
				"nws-current-temperature",
				"noaajs",
			},
			gomod: []string{
				"github.com/icodealot/noaa",
				"github.com/adamgreenhall/noaa",
				"github.com/chrisdobbins/noaa",
			},
		},
		{
			// Packagist has no client for this API at all, and searching it for
			// "opentdb" returns OpenTBS five times: tinybutstrong/opentbs,
			// makhov/opentbs and three Symfony bundles wrapping it. One letter
			// apart, D against B, and it merges LibreOffice and Word documents.
			// Nothing on Packagist calls opentdb.com, so nothing is mapped there.
			//
			// The same search turns up Open Tibia -- renatorib/otinfo and
			// pandaac/pandaac -- where OT abbreviates a game server, and
			// react/promise-timer, which describes itself as "a trivial
			// implementation of timeouts" and matched because trivia is a prefix
			// of trivial.
			//
			// On npm the bare name "trivia" answers a different API on
			// pareshchouhan-trivia-v1.p.mashape.com, a Mashape gateway that has
			// not existed since it became RapidAPI. @trivia-api/fetch and
			// @trivia-api/models are a competitor with an almost identical name,
			// the-trivia-api.com. @apiverve/trivia is a commercial aggregator's
			// own endpoint, and imdb-trivia scrapes IMDb.
			//
			// Three carry no URL at all, because the questions ship with them:
			// discord-trivia, @helyx/module-trivia -- "from a saved question
			// catalogue", in its own description -- and hubot-trivia-game.
			//
			// And the sharpest one is named after the API and does not call it.
			// github.com/tamnd/opentdb-cli builds opentdb.com/browse.php, which is
			// the human website, so it reads the pages a person would read rather
			// than api.php. Only quizzoro calls the API on Go, and its two module
			// paths are a fork pair.
			//
			// opendtb-wrapper is mapped and the typo is in its published name:
			// openDTB, not openTDB. It calls the right host.
			recipe: "opentdb",
			npm: []string{
				"open-trivia-db",
				"opentdb-api",
				"opentriviasdk",
				"opendtb-wrapper",
				"triviadb",
				"tdb-api",
				"trivia-cli",
				"@brecert/trivia",
			},
			gomod: []string{
				"github.com/handybots/quizzoro",
				"github.com/DucTheVulpe/quizzoro",
			},
		},
		{
			// One word apart, and a different company: BreweryDB is a commercial
			// API on api.brewerydb.com that required a key and closed its free
			// tier. It has more clients than this one does. brewerydb-node,
			// node-brewerydb, brewerydb-graphql and the bare npm name "brewery"
			// all call it, and so do ethanclevenger91/brewerydb-php and
			// liamthegrey/brewerydb-php on Packagist.
			//
			// The pair to watch is openbrewerydb-node against brewerydb-node.
			// Same suffix, same shape, one word of difference, and two different
			// hosts -- so a prefix or substring match on "brewerydb" maps half the
			// neighbourhood to the wrong API.
			//
			// "Brewery" on npm is also a software project with no beer in it.
			// @brewery/cli, brewery-core, brewery-core-cli and generator-brewery
			// belong to a scaffolding tool called The Brewery, and @brewery/cli's
			// only hosts are AWS and a private GitLab. @spbrewery/log is one
			// company's logger, @house-agency/brewmail is a mailer, and
			// @brewnode/server automates brewing hardware.
			//
			// Two more are beer and still not this. @umbraculum/brewery-api-client
			// covers "recipes, water, brew sessions, inventory" -- running a
			// brewery rather than listing them -- and @unclick/untappd-mcp is
			// Untappd. On Packagist, brewerwall is a vendor whose name begins with
			// brewer and whose packages are a barcode generator, beer calculations
			// and a unit-conversion model; beercalc_php's only host is hbd.org.
			//
			// Nothing on Go modules calls this API at all.
			recipe: "openbrewerydb",
			composer: []string{
				"joeymckenzie/openbrewerydb-php-api",
				"alexjustesen/php-openbrewerydb",
			},
			npm: []string{
				"openbrewerydb-node",
				"breweries-cli",
			},
		},
		{
			// "Poetry" in a package registry usually means the Python packaging
			// tool. snyk-poetry-lockfile-parser reads poetry.lock and
			// pyproject.toml against python-poetry.org; @almapario/vercel-poetry is
			// a runtime for it; and four semantic-release plugins bump versions in
			// its files. None of them is about poems.
			//
			// It is also an institution's product. ec-europa/oe-poetry-client is a
			// "Client library for the European Commission Poetry Service", which is
			// how translation is requested inside the Commission, and
			// oe-poetry-behat tests it. Same word, different building.
			//
			// The bare npm name "poetry" ships with an empty description and no URL
			// at all, and destinylab's lottery-poetry packages are Chinese
			// divination verses -- a fortune-telling engine rather than a library.
			//
			// Two more are poems and still not this API. chinese-poetry is a corpus
			// of three hundred thousand of them shipped as files, which is the
			// offline shape again, and balkhiate scrapes stands4.com.
			//
			// Nothing on Packagist calls poetrydb.org. On Go the only result is a
			// package inside a gRPC tutorial repository, which is an example rather
			// than a dependency, so nothing is mapped there either.
			recipe: "poetrydb",
			npm: []string{
				"@refkit/provider-poetrydb",
				"poemist",
			},
		},
		{
			// Packagist answers "datamuse" with its most popular packages --
			// monolog, doctrine/orm, doctrine/dbal, spatie/laravel-backup, phinx --
			// because nothing matches and the search falls back to popularity. A
			// detector reading search results rather than names would map Monolog
			// to a rhyming dictionary.
			//
			// The one real candidate there is undeadyetii/datamuse-php-api-wrapper,
			// and it is deliberately not mapped: its GitHub repository is gone, so
			// the source cannot be opened, and the only evidence is the registry's
			// own description. That is exactly the evidence latlon-utils had -- it
			// carries noaajs's description word for word and calls nothing -- so a
			// sentence is not enough on its own.
			//
			// Wordnik is the other dictionary API, and it needs a key. Searching
			// for it returns the entire Swagger toolchain, because Swagger was
			// created at Wordnik: swagger-node-express is "Wordnik swagger
			// implementation for the express framework" and still points at
			// api.wordnik.com, and swagger-validation and two forks come with it.
			// The vendor invented a spec format, so the format's tooling carries
			// the vendor's name.
			//
			// Three more compute or look elsewhere. nounsing-pro is built on the
			// CMU Pronouncing Dictionary at speech.cs.cmu.edu and ships the data;
			// rhymes-with decides whether two words rhyme without opening a
			// connection; and alfred-spell calls onelook.com, which is the same
			// author's other site and a different API.
			//
			// polymnia and melpomene are two Muses and one package: identical
			// descriptions, identical lib/make_request.js, both mapped.
			recipe: "datamuse",
			npm: []string{
				"datamuse",
				"datamuse-node",
				"datamuse-word-utils",
				"homophones",
				"wordfuzz",
				"wordie",
				"polymnia",
				"melpomene",
			},
			gomod: []string{
				"github.com/kostaspt/go-datamuse/v2",
				"github.com/stefanoghinelli/asyn",
			},
		},
		{
			// "artic" is a substring of "article", and that is what the registries
			// return. Packagist's results are all article extraction:
			// j0k3r/php-readability pulls articles out of HTML, j0k3r/graby does
			// the same job, facebook/facebook-instant-articles-sdk-php formats
			// them, georgringer/news publishes them and mtownsend/read-time counts
			// how long one takes to read. Nothing there touches a museum.
			//
			// npm's bare "artic" is the same word cut short -- "Artic. It's for
			// articles." in its own description -- and its only hosts are two
			// personal blogs.
			//
			// It is also a design system. artic-ui is "a modern React component
			// library" and asgrid is "Artic Studio CSS grid"; neither opens a
			// connection at all. And artic-enuu1, artic-enuu3 and
			// artic-ennnnuuu3 are published with no description and no code worth
			// the name.
			//
			// art-data-renderer is art and still not this: its hosts are
			// archive.vogue.com and doi.org.
			//
			// The two Go modules are one author's tool under two names --
			// tamnd/artic-cli and tamnd/artinstitute-cli, both calling
			// api.artic.edu/api/v1 -- and the same author's opentdb-cli, mapped
			// nowhere in the Open Trivia DB entry above, calls that provider's
			// website instead of its API. Same person, two habits.
			recipe: "artinstituteofchicago",
			npm: []string{
				"@refkit/provider-artic",
				"art-institute-of-chicago-api-client",
				"artic-sdk",
			},
			gomod: []string{
				"github.com/tamnd/artic-cli",
				"github.com/tamnd/artinstitute-cli",
			},
		},
		{
			// The vendor's own commercial product is the near miss. Ideal
			// Postcodes runs postcodes.io -- @stevegoossens/postcodes.io names
			// "ideal-postcodes/postcodes.io" in its own description -- and also
			// sells a paid UK lookup that needs a key.
			// @ideal-postcodes/postcode-lookup and its bundled twin call
			// ideal-postcodes.co.uk, not this. One company, two APIs, and only one
			// of them free.
			//
			// "Postcode" also means South Korea. react-daum-postcode,
			// vue-daum-postcode, @clroot/react-kakao-postcode and
			// @types/daum-postcode are all the Daum/Kakao address widget, on
			// kakaocdn.net, and they outnumber the UK ones in a search for the
			// bare word.
			//
			// Two more are competitors: getaddress-find is GetAddress.io and
			// @venditan/address-lookup is Google Places wearing a postcode label.
			//
			// And three carry no URL at all, because a postcode is a string with a
			// grammar: the bare npm name "postcode" is "UK Postcode helper
			// methods", postcode-validator validates by country from patterns, and
			// aus-postcode-suburbs ships Australia's list as data.
			//
			// node-red-contrib-postcodes-io is left out on purpose. It names
			// postcodes.io and builds no api.postcodes.io URL anywhere its source
			// can be read, so there is nothing to verify -- the same standard that
			// left Datamuse's one Packagist candidate unmapped.
			recipe: "postcodesio",
			composer: []string{
				"juststeveking/laravel-postcodes",
				"jabranr/postcodes-io",
				"ammaar23/postcodes-io-sdk",
				"boxuk/postcodes-io-bundle",
			},
			npm: []string{
				"node-postcodes.io",
				"postcodesio-client",
				"@cuvva/postcodesio-client",
				"@stevegoossens/postcodes.io",
				"@opsydyn/effect-postcodes-client",
				"@openpets/uk-postcodes",
			},
		},
		{
			// Packagist has nothing at all for this API, and npm's results are
			// mostly MCP servers -- seven of them, three of which are the same
			// server republished under three scopes.
			//
			// The one worth naming calls itself official and is not.
			// @dackerman-stainless/met-museum-demo is "The official TypeScript
			// library for the Met Museum API", published under a personal scope
			// with "demo" in the package name by an SDK-generation company. It does
			// call the API, so it is mapped -- but the word "official" in a
			// description is worth nothing on its own.
			//
			// Two more are museums and not this one: archival-imagery-mcp is the
			// Wellcome Collection and open-museum-mcp spans several.
			//
			// On Go, aaronland/go-metmuseum-openaccess is named for the Open Access
			// CSV dump and calls the API as well, so it maps.
			// dixieflatline76/Spice does too and is left out: it is a desktop
			// wallpaper application that happens to be a Go module, not something
			// anybody adds as a dependency.
			recipe: "metmuseum",
			npm: []string{
				"metmuseum",
				"met-client",
				"@mariolazzari/met",
				"@refkit/provider-met",
				"@chancehl/met-wallpaper",
				"@dackerman-stainless/met-museum-demo",
			},
			gomod: []string{
				"github.com/aaronland/go-metmuseum-openaccess",
			},
		},
		{
			// The near miss here is a fandom rather than a homonym: the show has
			// more packages about it than clients of its API.
			//
			// rick-and-morty and rick-and-morty-cli both fetch GIFs from
			// i.giphy.com. rick-morty-components-lib is a React component library
			// themed after the show -- it names rickandmortyapi.com in its README
			// and its code calls fb.me, which is React -- so a detector reading
			// documentation rather than source would map a colour scheme.
			// @zsy1126/dsh-rick-and-morty is a skin for a web GUI. And "squanch" is
			// a text mangler named after a catchphrase.
			//
			// Two MCP servers are Frinkiac, which is screencap search across the
			// Simpsons, Futurama and this -- a different API that happens to cover
			// the same show.
			//
			// @giovanihdz/rickandmortyapi is mapped and is one of a small cohort of
			// Spanish-language coursework -- "Consumiendo API de rickandmorty",
			// with @imaciasd/rickandmortyapi-axios and @hiramnm14/rickandmortyapi_axios
			// beside it. It calls the API in simpleGet.js, so the standard that
			// mapped everything else maps it too; the other two were not opened.
			recipe: "rickandmorty",
			composer: []string{
				"nickbeen/rick-and-morty-api-php",
				"zmaglica/rick-and-morty-api-wrapper",
			},
			npm: []string{
				"rickmortyapi",
				"@giovanihdz/rickandmortyapi",
			},
		},
		{
			// A deck of cards is a data structure, so nearly every package with
			// the name is a local implementation that opens no connection:
			// deck-of-cards, cards.js, card-deck, deckjs, cardgames,
			// card-games-typescript, deck-of-cards-ts, real-deal-deck-of-cards,
			// npm-cards, preferans-deck-js, @mikeygough/deck-of-cards, and
			// @andrewscripts/deck-of-cards.js beside @andrewcreated/deck-of-cards.js,
			// which are the same library under two scopes.
			//
			// The bare name is one of them. npm's "deckofcards" is "a standard deck
			// of playing cards" and has nothing to do with this API, so the most
			// obvious name in the registry is the wrong one.
			//
			// And a new kind: a package about the finding rather than about the
			// provider. api-false-success is "Measured behaviour of public APIs on
			// requests that cannot succeed: which return HTTP 200 anyway, how to
			// detect it per API, and whether it is a defect or a documented
			// contract" -- a dataset describing exactly what this Recipe's headline
			// pins, which opens no connection to anything.
			//
			// Packagist has nothing at all.
			recipe: "deckofcards",
			npm: []string{
				"deckofcardsapi-client",
				"deckofcards-api",
				"node-deckofcards",
			},
		},
		{
			// Packagist reads "themealdb" as "theme". Its results are
			// parallax/filament-syntax-entry ("themeable syntax highlighting"),
			// yrizzz/adminkit ("themeable admin panel") and
			// nasirkhan/laravel-sharekit -- the first five letters of the provider
			// are a word, and that word is what the search matched. Same shape as
			// artic inside article.
			//
			// tamnd is here for the third time. mealdb-cli calls the API, artic-cli
			// and artinstitute-cli call the Art Institute's, and opentdb-cli --
			// unmapped in the Open Trivia DB entry above -- calls that provider's
			// website instead. One author, four CLIs, three providers, and one of
			// them scraped.
			//
			// koishi-plugin-cooke-tmd and koishi-plugin-cooke-mda are the same
			// Chinese-language bot plugin under two names, with identical
			// descriptions. Both call the API, so both map.
			//
			// a3_os_cooking_recipes is mapped and is not what its name suggests: an
			// "open-source library" for cooking recipes that turns out to be a
			// TheMealDB client with an Arma-3-shaped name.
			recipe:   "themealdb",
			composer: []string{"xxiv/mealdb"},
			npm: []string{
				"mealdb",
				"meal-api",
				"a3_os_cooking_recipes",
				"koishi-plugin-cooke-tmd",
				"koishi-plugin-cooke-mda",
			},
			gomod: []string{"github.com/tamnd/mealdb-cli"},
		},
		{
			// Packagist for "world bank" is entirely the word "world". Its results
			// are pn/swiftcodes ("Banks in the world"), two PayTR gateways where
			// World is a Turkish credit card brand, two purchasing-power-parity
			// packages, and laravel-zip -- twice -- described as "the world's
			// leading zip utility". Nothing there is this API, and one of them
			// matched on a marketing superlative.
			//
			// npm is mostly MCP servers, seven of them, including one that is only
			// an npx launcher for a hosted service and one scoped to Nigeria.
			//
			// world-bank-dataset is the data as a file: "A JSONified dataset of
			// projects funded by the World Bank", which opens no connection.
			//
			// @volks publishes six packages with placeholder descriptions -- "A
			// worldbank-indicator", "A worldbank-sources", "A worldbank-topics" --
			// and they are not all the same. countries and indicator call
			// api.worldbank.org and are mapped; the cli does not, and is not.
			recipe: "worldbank",
			npm: []string{
				"@volks/worldbank-countries",
				"@volks/worldbank-indicator",
			},
		},
		{
			// The clearest example of the data shipped as a file that this
			// collection has: simple-country-names is "The list of all the
			// countries in the world", and it is one static JavaScript module
			// holding disease.sh's country list -- names, iso2, iso3, lat, long,
			// and every flag URL on disease.sh's own host. The vendor's hostname
			// appears on every one of its two hundred entries and it never opens a
			// connection. Grepping for the host maps it; reading what the host is
			// doing there does not.
			//
			// novelcovid is mapped and is named after the API's former name: this
			// provider was NovelCOVID before it was disease.sh, so the client
			// carries the old one and the Recipe carries the new.
			//
			// Packagist has three packages called module-disease, from three
			// different vendors, with the same one-line description and the same
			// typo in it -- "data diasease in klinik or puskesmas" -- which is
			// Indonesian clinic software forked twice and nothing to do with this.
			//
			// And two npm results for "covid api" are icon sets:
			// @stacksjs/iconify-covid and @ngxi/covid are "Covid Icons".
			recipe: "diseasesh",
			npm: []string{
				"@aero/corona",
				"novelcovid",
				"covid19-vn",
			},
		},
		{
			// A near miss that is only a type definition:
			// @mgrzmil-org/api-types is "Generated TypeScript types + OpenAPI spec
			// for dog-api" and contains no reference to dog.ceo anywhere in it --
			// not in code, not in documentation. It describes the API's shapes
			// without naming the API, so there is nothing to verify and nothing to
			// map.
			//
			// dog-names is the other kind: "Get popular dog names", which is a list
			// shipped as data and opens no connection.
			//
			// The rest of what "dog api" returns is the word "api": fast-diff,
			// diff-match-patch, wrap-ansi, slice-ansi, mock-axios and frak all came
			// back for it, and none of them is about dogs at all.
			//
			// sabina/project-dogapi is mapped and says what it is in its own
			// description: "Task - use dogceo API to list all the dog breeds". A
			// homework assignment that calls the API is still a package that calls
			// the API.
			recipe: "dogceo",
			composer: []string{
				"jomingeorge/dog-ceo",
				"sabina/project-dogapi",
			},
			npm: []string{
				"dog-ceo",
				"random.dog.img",
				"gatsby-source-dog",
				"doggo-api-wrapper",
				"dog-image-generator",
			},
		},
		{
			// "Slip" is a networking standard. @serialport/parser-slip-encoder
			// implements SLIP, the Serial Line Internet Protocol, which has been
			// framing packets since 1988 and has nothing to do with advice.
			//
			// minicli/advice-slip is mapped and calls itself "Demo - Advice Slip",
			// which is what most of this neighbourhood is: a small API with an
			// obvious shape attracts tutorials, and four of the five packages that
			// call it are one person's example.
			recipe:   "adviceslip",
			composer: []string{"minicli/advice-slip"},
			npm: []string{
				"@draggem/adviceslip-wrapper",
				"advice-slip-cli",
				"random-advice",
				"@dennisk2025/random-advice-slip",
			},
		},
		{
			// One API, four hostnames, three response shapes -- and the packages
			// are split across all of them.
			//
			// swapi.co is the original and answers 301 now. swapi.dev and
			// swapi.info both answer the bare record this Recipe describes.
			// swapi.tech does not: it wraps everything in
			// {"message": "ok", "result": {"properties": {...}, "_id": "...",
			// "description": "A planet."}}, with the same field names three levels
			// down, a Mongo identifier beside them, and created and edited
			// timestamps that are today's date because they are generated when the
			// record is read. A client written against one and pointed at the other
			// reads body.name and gets undefined.
			//
			// So react-swapi is not mapped. It calls swapi.tech, which is a
			// different API wearing the same name and the same field names.
			// ts-swapi, star-wars-api and the bare npm "swapi" call swapi.co and
			// are mapped: the host is gone, the shape they were written against is
			// this one.
			//
			// And @swapi-finance/sdk and @swapi-finance/contracts are a
			// decentralised exchange -- "An SDK for building applications on top of
			// Swapi", swap plus API, with a bee for a logo. Its own scope, its own
			// word, nothing to do with any of this.
			//
			// xvrmallafre/swapi is on Packagist with a description that is entirely
			// about the neighbourhood: "This is a fork of rmasters/swapi that works
			// again."
			recipe: "swapi",
			composer: []string{
				"gale/swapi",
				"christiancocco/swapi",
				"jlgomes/swapi-php",
			},
			npm: []string{
				"swapi",
				"swapi-node",
				"swapi-handler",
				"nuxt-swapi",
				"ts-swapi",
				"star-wars-api",
			},
		},
		{
			// There are two Chuck Norris APIs and the neighbourhood splits between
			// them. ICNDb -- the Internet Chuck Norris Database, at icndb.com -- is
			// the older one and is gone. chucknorris-quotes, chucknorris-joke-node,
			// chuck-jokes and the npm package called "cn" all call it.
			//
			// Packagist has no client for this API at all. Its four results share
			// one description word for word -- "Create random Chuck Norris jokes" --
			// across mpociot, zaratedev, vikas5914 and ttnppedr, and the two that
			// were opened both call icndb.com. Four forks of one package, pointed at
			// a host that no longer answers.
			//
			// Three more open no connection at all. The bare npm name
			// "chucknorris" -- "Chuck Norris does not need a description" -- has no
			// host in it; chuck-norris-api is named for ICNDb in its own
			// description and calls nothing; and chuckscript is an esoteric
			// programming language whose source is written entirely in zeroes.
			recipe: "chucknorris",
			npm: []string{
				"chucknorris-io",
				"chuck-norris-jokes",
			},
		},
		{
			// One client, and it corroborates the Recipe: node-red-contrib-iss-location
			// hardcodes "http://api.open-notify.org" as its default domain, because
			// there is no https one to reach. It is also named for a package that was
			// never published -- npm has no iss-location.
			//
			// Packagist has nothing at all. Searching it for "iss" returns ten
			// packages about issues: composer-dependency-analyser, phpcompatibility,
			// php-deprecation-detector. The substring again, the way "artic" sat
			// inside "article". Its one open-notify hit, jjonline/wx-open-notify, is
			// a WeChat open-platform callback SDK: the ordinary words, not the API.
			//
			// Two MCP servers carry the name and neither is mapped, the way the
			// Open-Meteo and USGS ones were not. @pipeworx/mcp-open-notify is the
			// fifth Recipe running for that publisher, after goproxy, deps.dev,
			// Open-Meteo and USGS. @mcp-mk/iss calls this API and Where The ISS At
			// both, and says so in its own description.
			//
			// And the exe wrapper at a scale worth recording: PeopleInSpace ships as
			// five npm packages, one per CPU architecture, each between forty and
			// ninety-five megabytes. Every one of them holds the same two files --
			// compose-desktop-1.0-SNAPSHOT-all.jar and jar-runner.jar -- and the
			// string "open-notify" appears in none of them, because whatever they
			// call is compiled into the jar. A Java desktop demo, published to a
			// JavaScript registry five times over.
			//
			// The ordinary word is "space": jetbrains-space-api is a team
			// collaboration product, deta-space-client is a cloud platform,
			// ocr-space-api reads text out of images, and space-station is
			// "Atomic Mission Control for Parallel Operations".
			recipe: "opennotify",
			npm: []string{
				"node-red-contrib-iss-location",
			},
		},
		{
			// The siblings share a budget but not a hostname, and detection follows
			// the hostname. demografix calls all three -- api.agify.io,
			// api.genderize.io and api.nationalize.io -- and calls itself the
			// official client for them, which is the vendor confirming in public
			// what the shared counter shows. Its PHP twin is demografix/demografix,
			// whose Client.php holds HOST_AGIFY alongside the other two and which
			// ships a Quota model and a RateLimitError to go with them.
			//
			// Not mapped: the clients that call only a sibling. The bare npm name
			// "genderize" reaches api.genderize.io and nothing else;
			// nationality-guesser reaches api.nationalize.io and nothing else; and
			// pixelpeter, alexanderpavlov, magdy-hakam, javihgil and barttyrant all
			// wrap genderize.io on Packagist. One allowance, three hosts, and a
			// dependency on one of them is not a dependency on this one.
			//
			// The near miss is the scope that is not the vendor, the way
			// internetarchive/virustotalapi was. @agify/n8n-nodes-leadifyv2 and
			// @agify/n8n-nodes-leadifyv3 sit under the vendor's own npm scope and
			// are about "Leadify API for lead management": opened, they contain no
			// API host at all.
			//
			// And the edit distance, one character wide. Searching Packagist for
			// "agify" returns ten packages for Apify -- laravel-apify, megaads
			// /apify, php-apify-sdk, consoletvs/apify, dantepiazza/apify and
			// apify/apify-client, which describes itself as "AI-generated and
			// AI-maintained". A different company, a different product, and a name
			// that differs by a letter.
			//
			// Two MCP servers again, neither mapped: @pipeworx/mcp-agify is the
			// sixth Recipe running for that publisher, after goproxy, deps.dev,
			// Open-Meteo, USGS and Open Notify, and @mcp-mk/names covers all three
			// siblings at once.
			//
			// The competitor is Gender-API.com: gender-api.com-client on npm and
			// gender-api/client on Packagist are a different vendor answering the
			// same question. genderize-br is the ordinary word -- grammatical
			// gender inflection in Portuguese, no API in it at all.
			recipe: "agify",
			composer: []string{
				"demografix/demografix",
			},
			npm: []string{
				"demografix",
				"@demografix/n8n-nodes-agify",
				"@pipedream/agify",
			},
		},
		{
			// "Bible API" is the name of a category, not of this API, and the
			// neighbourhood is the largest yet for one Recipe. At least seven
			// different products answer to it and each has its own clients.
			//
			// One package reaches this one and says so: the-holy-bible-api declares
			// API_ROOT = "https://bible-api.com" in its compiled output.
			//
			// The near miss is one letter of top-level domain. The npm scope
			// @bible-api publishes a translation per package -- kjv1769, vulgate,
			// the Polish Biblia Gdanska of 1632 -- and every one of them points at
			// bible-api.io. The scope is named for this API and calls a different
			// one, which is the data-shipped-as-a-file kind and the vendor-namespace
			// kind at once.
			//
			// The others are simply other APIs: @helloao/cli, @helloao/tools and
			// free-use-bible-sdk reach bible.helloao.org; holy-bible-api reaches
			// holy-bible-api.com and is generated, carrying openapi-generator.tech
			// and a "your-api-endpoint.com" placeholder nobody replaced;
			// verse-of-the-day-cli scrapes biblegateway.com; and on Packagist
			// prevailexcel/laravel-bible is for API.Bible, ichthus-soft/bible-api is
			// the Romanian Biblia Cornilescu, bkuhl/bible-csb is the Christian
			// Standard Bible and apiverve/bible is APIVerve's.
			//
			// Three more open no connection at all. The bare npm name "bible-api",
			// whose whole description is "Bible API", contains no host; and
			// ngx-bible-api-client and axios-bible-api-client -- "A Bible Api Client
			// for Angular" and "for Axios" -- carry no host either, so there is
			// nothing in them to say which of the seven they mean.
			//
			// Packagist has no client for this API at all, the way it had none for
			// Chuck Norris or Open Notify. rootxs/sudobible is the closest thing and
			// it is the other side of the wire: "Open source Bible API", a server.
			recipe: "bibleapi",
			npm: []string{
				"the-holy-bible-api",
			},
		},
		{
			// A crowded neighbourhood for once, and most of it is real. Five npm
			// packages carry api.coingecko.com in their published output, the
			// official TypeScript library among them, and two Composer packages do
			// too: tigusigalpa/coingecko-php holds api.coingecko.com/api/v3 in its
			// config, and codenix-sv/coingecko-api names the host in its README.
			//
			// The MCP server is the vendor's own this time. @coingecko/coingecko-mcp
			// is published under the same scope as the official library, and it is
			// still not mapped: every MCP server so far has been somebody else's,
			// and the reason for leaving them out does not change when the publisher
			// is the vendor. Installing one means an agent can reach the API, not
			// that the project does.
			//
			// @types/coingecko-api is types with no client -- DefinitelyTyped
			// declarations for coingecko-api, carrying no host of their own, the way
			// @mgrzmil-org/api-types did.
			//
			// geckoterm is the vendor's other product: "Coin Gecko Terminal Open
			// Source SDK", and it reaches api.geckoterminal.com. Same company, a
			// different API, a different host.
			//
			// And the ordinary word is a browser engine. Searching npm for
			// "coin-gecko" returns ua-is-frozen and universal-user-agent, which
			// match on Gecko as in Mozilla's rendering engine, in the user-agent
			// string every browser still sends. @flatfile/gecko is a mascot -- "a
			// green gecko wearing glasses" -- and ccxt/ccxt is a trading library
			// spanning more than a hundred exchanges that mentions this one in
			// passing. escudo/api-coin-gecko-test has no description at all.
			recipe: "coingecko",
			composer: []string{
				"codenix-sv/coingecko-api",
				"tigusigalpa/coingecko-php",
			},
			npm: []string{
				"@coingecko/coingecko-typescript",
				"coingecko",
				"coingecko-api",
				"coingecko-api-v3",
				"coingecko-api-pro",
			},
		},
		{
			// The near miss here is an XML namespace. Apple defines an RSS extension
			// for podcasts -- itunes:author, itunes:category, itunes:explicit -- so
			// feed libraries carry the word without touching any HTTP API at all.
			// zetacomponents/feed parses RSS1, RSS2 and ATOM; marcw/rss-writer
			// writes them; silverorange/castanet generates podcast feeds
			// "including iTunes-specific fields"; and lukaswhite/itunes-categories
			// is a list of the category names that extension allows. None of them
			// is mapped, and none of them makes a request.
			//
			// The second cluster is Apple's other APIs, which is most of what
			// Packagist returns for the word. aporat/store-receipt-validator,
			// kartina-tv/store-receipt-validator, readdle/app-store-server-api,
			// yanlongli/app-store-server-api and
			// readdle/app-store-receipt-verification are all about receipts;
			// snscripts/itc-reporter downloads sales figures from iTunes Connect.
			// On npm, @types/apple-music-api and @yujinakayama/apple-music are for
			// the Apple Music API at api.music.apple.com, which needs a developer
			// token, and @expo/apple-utils says in its own description that it uses
			// "the unofficial iTunes APIs" -- a different set again.
			//
			// Two MCP servers, neither mapped. @pipeworx/mcp-itunes-search is the
			// seventh Recipe running for that publisher, after goproxy, deps.dev,
			// Open-Meteo, USGS, Open Notify and Agify -- and @pipeworx/mcp-movies
			// wraps this API and TVmaze together, which is two shipped Recipes in
			// one package.
			recipe: "itunes",
			composer: []string{
				"atomescrochus/laravel-itunes-search-api",
				"timlebrun/itunes-search",
			},
			npm: []string{
				"itunes",
				"itunes-search",
				"node-itunes-search",
				"itunes-search-api",
				"itunes-store-api",
				"app-store-scraper",
			},
		},
		{
			// Four npm packages carry api.first.org/data/v1/epss in their published
			// output and are mapped. The sharpest near miss is the one that uses the
			// same data without the API: npm-epss-audit reads the bulk CSV from
			// epss.cyentia.com and links www.first.org for the documentation, so it
			// is about EPSS in every sense except the one detection can act on.
			// vulnscope is the same shape with its database bundled from GitHub --
			// "OSV + KEV + EPSS", none of them fetched from here.
			//
			// Packagist has no client for this API at all, and what it does return
			// for "epss" is a thermal receipt printer: mike42/escpos-php, where the
			// letters sit inside ESC/POS. The rest of that page matches on "ess" --
			// essence/essence, essence/http, essence/dom, essa/api-tool-kit,
			// nunomaduro/essentials. Searching it for "first.org" returns Firestore
			// clients and a Magento first-order discount.
			//
			// Six MCP servers, the most for any Recipe here, and none mapped:
			// @pipeworx/mcp-epss is the eighth Recipe running for that publisher,
			// and cve-mcp, @datanexusmcp/mcp-server, @iflow-mcp/badchars-cve-mcp,
			// cve-fix-check and @mcpskillsio/server all list this API among their
			// sources.
			//
			// And a new kind of neighbour: cctx exists to occupy a typo. Its whole
			// description is "Looking for CCXT? You've mistyped the name. Bad guys
			// could exploit this" -- a package published so that nobody else can
			// publish it, which turned up here because it contains the word exploit.
			recipe: "epss",
			npm: []string{
				"@absolutejs/vulnerabilities-epss",
				"breach-gate",
				"niyantri",
				"fad-checker",
			},
		},
		{
			// The near miss is the vendor namespace again, and it is not even the
			// vendor's repository. Packagist has a chesscom namespace, and what it
			// publishes there is chesscom/honeybadger-php and
			// chesscom/honeybadger-bundle -- error monitoring for Honeybadger.io,
			// from github.com/ShonM. chesscom/chess-game is under the same namespace
			// and "represents a chess game as a php object": the domain, not the
			// API.
			//
			// Then a whole shelf of chess libraries that never make a request.
			// p-chess/chess, akondas/chess.php, ryanhs/chess.php and onspli/chess
			// generate and validate moves; amyboyd/pgn-parser reads the notation;
			// chess.js on npm is the same thing and links lichess.org. Move
			// generation is not a client of anything.
			//
			// And the fuzzy matches, which are the widest yet for one word.
			// wiserwebsolutions/laravel-chesco-gis is Chester County's ArcGIS
			// FeatureServer; miripiruni/csscomb sorts CSS properties;
			// chessio/module-matomo is Magento analytics; elioair/chesscaptcha is a
			// captcha; chescos/cs-proto has no description at all.
			//
			// Four MCP servers, none mapped -- @pipeworx/mcp-chess is the ninth
			// Recipe running for that publisher -- and two of them, along with
			// cli-chess-puzzle, cover Lichess as well, which is the other server for
			// the same game.
			//
			// Three clients that look right were left out on evidence:
			// chess-web-api, libchess and chesster all describe themselves as
			// wrappers for this API and none of them carries the host in its
			// published output. Packagist's two, flowtt/chesscom-php and
			// joebocock/chess-api-php, name it in neither README.
			recipe: "chesscom",
			npm: []string{
				"chess-com",
				"chess.com.js",
				"chesscom-ts-client",
				"chesscom-sdk",
				"chess-api-typescript",
				"online-chess-api",
			},
		},
		{
			// The near miss corroborates the Recipe. react-currency-localizer is "a
			// React hook to display prices in a user's local currency using
			// HTTPS-compatible IP geolocation" -- it advertises HTTPS compatibility
			// as a feature, which is only worth advertising because the obvious
			// choice is not, and it calls ipapi.co. One hyphen away from this host,
			// a different company, chosen for the trait this Recipe is about.
			//
			// And a new sense of the letters: ethernet-ip-cip is "a simple node.js
			// Ethernet/IP API focused on CIP implementation". EtherNet/IP is an
			// industrial automation protocol, and the IP in it is Industrial
			// Protocol rather than Internet Protocol. shen-rate-limit matches from
			// the other direction, on "IP/API Key rate-limiting" -- a slash where
			// the name has a hyphen.
			//
			// Then the parsing libraries, which never make a request: ip-regex,
			// is-ip, ip-address, internal-ip and the bare name ip all match on the
			// two letters and none of them asks anybody anything. public-ip does
			// make a request, to somewhere else.
			//
			// The rest is a market. MaxMind (@maxmind/geoip2-node, geoip2/geoip2),
			// IPinfo (node-ipinfo, ipinfo/ipinfo), IP2Location, ipwhois and
			// apiip.net all answer the same question from different hosts.
			//
			// pulkitjalan/geoip is mapped where ccxt/ccxt was not, and the
			// difference is proportion: it ships src/Drivers/IPApiDriver.php beside
			// two others, so this API is one of a handful it is built on rather
			// than one of a hundred it mentions.
			//
			// Left out on evidence: ipwhere calls itself a "Geo IP API framework"
			// and carries no host, and aungmyokyaw/ip-api's repository is gone.
			recipe: "ipapi",
			composer: []string{
				"maciejkrol/ipapicom",
				"pulkitjalan/geoip",
				"pulkitjalan/ip-geolocation",
			},
			npm: []string{
				"@juicyllama/nestjs-ip-api",
			},
		},
		{
			// Packagist has nothing at all, and what it does return is one person's
			// initials. Nine packages by tflori -- an ORM, a dependency injector, a
			// PHPUnit printer, a console library -- plus tflanagan/quickbase, all
			// matched on a vendor prefix that happens to spell the acronym.
			//
			// And the letters mean something else again on npm: tfl-lsp is a
			// "Language Server Protocol for TFL Fitch-style proofs", where TFL is
			// Truth-Functional Logic. A proof checker for a formal logic course,
			// matched on three letters shared with a railway.
			//
			// tfl-to-freeagent is the data as a file: it converts an exported CSV of
			// journeys into accounting software, so it is about this provider and
			// never calls it.
			//
			// Five MCP servers, the most for any Recipe here after EPSS, and none
			// mapped: mcp-tfl-journey, @daanrongen/tfl-mcp, @rmp/tfl-mcp,
			// london-transport-mcp and one whose name is nine words long.
			//
			// Left out on evidence: the bare name "tfl", tfl-api-wrapper and
			// @onemedia/tfl-data-module all describe themselves as clients of this
			// API and none of them carries the host in its published output.
			recipe: "tfl",
			npm: []string{
				"tfl.api",
				"tfl-ts",
				"tfl-api",
				"tfl-client",
				"tfl-unified-api",
			},
		},
		{
			// The whole npm neighbourhood is agent tooling. Every one of the twelve
			// results for "openfda" is an MCP server or a tool for one --
			// @ythalorossy/openfda, @cyanheads/openfda-mcp-server,
			// @pipeworx/mcp-openfda, @pipeworx/mcp-fda-devices, mcp-openfda,
			// bach-openfda, openfda-mcp-server, @sineai/drug-safety-mcp,
			// biopharma-catalyst-mcp, @shaddyt/clinical-reference-mcp and
			// medical-mcp. No ordinary client library appears at all, which is the
			// first time that has been true of a provider here.
			//
			// One package is not an MCP server and is mapped: @pharmatools/drug-data
			// carries api.fda.gov and resolves drug names against RxNorm beside it.
			//
			// Packagist answers "openfda" with two other open- projects.
			// open-feature/sdk, open-feature/flagd-provider,
			// launchdarkly/openfeature-server and configcat/openfeature-provider are
			// OpenFeature, the feature-flag standard; freedsx/ldap,
			// ldaptools/ldaptools and causal/ig_ldap_sso_auth are OpenLDAP. A prefix
			// shared with two unrelated specifications.
			//
			// And the acronym belongs to more than one country:
			// aghfatehi/laravel-saudi-fda is the Saudi Food and Drug Authority, a
			// different regulator with the same three letters.
			// montopolis/fda-nutrition-rounding-php implements this agency's
			// rounding rules for nutrition labels and makes no request at all.
			//
			// Left out on evidence: meisam-mulla/laravel-openfda is described as a
			// "Simplified OpenFDA API for Laravel" and its repository cannot be
			// read, and n8n-nodes-openfda-drug-scraper carries no host in its
			// published output.
			recipe: "openfda",
			npm: []string{
				"@pharmatools/drug-data",
			},
		},
		{
			// The Open- prefix collides with every other Open- project, and this is
			// the second Recipe in a row where that is the whole Packagist result.
			// openFDA got OpenFeature and OpenLDAP; openalex gets OpenExchangeRates
			// four times over -- qbil-software/openexchangeratesbundle,
			// qbil-software/openexchangerates, dandelionmood/openexchangerates and
			// mrzard/open-exchange-rates-service -- beside openagenda/sdk-php and
			// openclerk/country-list. A currency API, an events API and a list of
			// countries, for a catalogue of academic papers.
			//
			// Eight MCP servers carry the name, and none is mapped.
			//
			// Left out on evidence: @univ-lehavre/atlas-biblio-cli is a "CLI entry
			// point for OpenAlex bibliographic data validation" and carries no host
			// in its published output, and mbsoft31/laravel-openalex names the API
			// in its description but not in a README that can be read.
			recipe: "openalex",
			npm: []string{
				"openalex",
				"@sickrin/openalex-sdk",
				"openalex-skill",
			},
		},
		{
			// The near miss is the same agency's previous API, on its own host, and
			// the split is even. Three packages reach api-v3.mbta.com and are mapped;
			// boston-mbta, mbta-data and mbtapi reach realtime.mbta.com, which is the
			// v2 Realtime API -- a different host, a different shape, and one of them
			// says "v2" in its own description.
			//
			// The bare name is taken by the documentation. npm's "mbta" describes
			// itself as "an api for your api" and the only host in it is
			// developer.mbta.com, the portal rather than the service.
			//
			// Packagist has nothing, and what it answers with is metadata:
			// jms/metadata, symfony/property-info, composer/metadata-minifier and
			// cmb2/cmb2, a WordPress metabox library. Two more are mutation-testing
			// frameworks, infection/infection and pestphp/pest-plugin-mutate, matched
			// on the same four letters in a different order.
			//
			// Left out on evidence: @swittuth/mbta-tracker-sdk and where-is-mbta both
			// describe themselves as clients of this API and neither carries a host
			// in its published output. One MCP server, @pipeworx/mcp-mbta, is the
			// eleventh Recipe running for that publisher.
			recipe: "mbta",
			npm: []string{
				"mbta-client",
				"mbta-api",
				"mbta-vis",
			},
		},
		{
			// The ordinary word, at the largest scale it has reached here. Packagist
			// answers "carbon" with nesbot/carbon, one of the most installed PHP
			// packages there is and a DateTime library, beside
			// carbonphp/carbon-doctrine-types, cmixin/business-day,
			// kylekatarnls/laravel-carbon-2 and genealabs/laravel-null-carbon. npm
			// answers it with IBM's Carbon Design System -- @carbon/layout,
			// @carbon/styles, @carbon/icons, @carbon/icons-react and
			// @carbon/icon-helpers. A date library, a design system and an element,
			// sharing one word.
			//
			// co2-data is the data as a file again: "Carbon intensity of 66 countries
			// and a world average updated to 2020", a dataset with no host in it,
			// about the same quantity this API measures for one country in real time.
			//
			// Two MCP servers, neither mapped -- @pipeworx/mcp-carbon is the twelfth
			// Recipe running for that publisher.
			//
			// Left out on evidence: medigeek/php-carbon-intensity names
			// carbonintensity.org.uk in its own Packagist description and its
			// repository cannot be read.
			recipe: "carbonintensity",
			composer: []string{
				"sebrave/laravel-carbon-intensity",
			},
			npm: []string{
				"@serodigital/national-grid-api-client",
				"carbon-intensity-plugin",
				"homebridge-plugin-uk-nationalgrid-carbonintensity",
			},
		},
		{
			// The near miss is "police" as a metaphor. case-police is "Make the case
			// correct, PLEASE!" -- it fixes the capitalisation of words like
			// javascript -- and eslint-plugin-case-police wraps it for a linter.
			// openapi-police is OpenAPI validation. Three packages where policing
			// means enforcing a rule rather than a constabulary.
			//
			// And the country code means another country. dictionary-uk is a
			// Ukrainian spelling dictionary: uk is Ukraine's language code, and the
			// United Kingdom's is gb. Two letters, two countries, one search.
			//
			// Two MCP servers carry the name, neither mapped, and homedata-mcp is
			// a different UK property API entirely.
			//
			// Left out on evidence: bespoke-support/police-api is "API and Data for
			// Police UK" and its repository cannot be read, and
			// roadsigns/laravel-police-uk has no description at all.
			recipe: "ukpolice",
			composer: []string{
				"mightymango/data-police-uk",
			},
			npm: []string{
				"pluk",
				"ukpd",
				"@pipedream/data_police_uk",
				"koop-provider-ukcrime",
			},
		},
		{
			// Both words in the name pull in a different unrelated cluster. "Rating"
			// brings star widgets -- @zag-js/rating-group, @fluentui/react-rating,
			// @smastrom/react-rating, react-simple-star-rating and the bare
			// react-rating, all UI components for letting somebody click five stars.
			// "Hygiene" brings danger-plugin-pr-hygiene, which enforces good pull
			// request hygiene. Neither has anything to do with food.
			//
			// And the acronym belongs to a CSS framework. Packagist answers "fhrs"
			// with twbs/bootstrap, twitter/bootstrap, zurb/foundation,
			// foundation/foundation-sites and coreui/coreui -- five front-end
			// frameworks -- beside mjoc1985/fhr-laravel-client, which is an "FHR API
			// client (search, inventory, cart)" for something else entirely.
			//
			// Four MCP servers carry the name, and none is mapped.
			// @pipeworx/mcp-uk-food-hygiene is the thirteenth Recipe running for that
			// publisher; three more are @gera-services', one of which exists to let
			// an agent "check if a UK business is real".
			recipe: "fhrs",
			npm: []string{
				"food-hygiene-ratings",
				"@wmfs/food-hygiene-blueprint",
			},
		},
		{
			// Packagist has nothing for "bank of canada" and eight packages for
			// "valet", all of them Laravel Valet: laravel/valet, cpriego/valet-linux,
			// cretueusebiu/valet-windows, genesisweb/valet-linux-plus,
			// weprovide/valet-plus, ycodetech/valet-windows, cpriego/valet-ubuntu and
			// jacerider/valet, which is a Drupal navigation system. A local
			// development environment, eight ways.
			//
			// npm spreads the word further: @archway/valet is a CSS-in-JS engine,
			// @archway/create-valet-app scaffolds one, @caiokf/valet drives coding
			// CLIs, and @puffle/secret-valet is a password prompt. Valet is a name
			// for a helper, and everybody has one.
			//
			// boc-exchange-rate and @allratestoday/boc-exchange-rate are the data as
			// a file, and they call themselves "Official Bank of Canada (Canada)
			// daily exchange rates" while carrying no host at all -- the rates are
			// bundled, not fetched.
			//
			// @cybrid/cybrid-api-bank-angular and its TypeScript twin are a different
			// company's bank API entirely. Two MCP servers carry the name, and
			// @pipeworx/mcp-bank-of-canada is the fourteenth Recipe running for that
			// publisher.
			recipe: "bocvalet",
			npm: []string{
				"@varve/boc-valet",
			},
		},
		{
			// Two packages carry zenodo.org/api in their published output and are
			// mapped. Four more describe themselves as clients of this API and build
			// the URL from parts, so nothing in them says where they go: the bare
			// name "zenodo" ("Node.js library to access the Zenodo API"), zenodo-ts,
			// zenodraft and @citation-js/plugin-zenodo. Four of six named clients
			// with no host between them.
			//
			// The badge kind is here again: code-health-meter matches on a Zenodo DOI
			// badge in its README, which is a link to a citation rather than a call.
			//
			// Packagist has nothing, and answers the word three ways. zenodorus is a
			// vendor prefix -- zenodorus/strings, zenodorus/filesystem,
			// zenodorus/arrays and zenodorus/core, utility libraries by one author.
			// zenorocha/clipboardjs is a person's name, and clipboard.js is one of
			// the most installed front-end libraries there is. And zendframework
			// /zendoauth and zendframework/zendopenid are Zend Framework, where the
			// letters line up by accident.
			recipe: "zenodo",
			npm: []string{
				"zenodo-lib",
				"@iomeg/zenodo-upload",
			},
		},
		{
			// The near miss is a substring inside a word about protection. Packagist
			// answers "uniprot" with unisolutions/silverstripe-uniprotect and
			// gelysis/gs4-uniprotect, two spam-protection modules for Silverstripe,
			// beside soderlind/unprotect-protected-posts, cacing69/uxcel (which
			// protects and unprotects spreadsheets) and
			// dkplus/csrf-api-unprotection-bundle. iconic/uniproperty gets in the
			// same way, on a different word.
			//
			// The bare npm name is a file parser: "uniprot" is a "Uniprot text file
			// parser" and links www.uniprot.org, the website, rather than the REST
			// host -- the data as a file again. vsm-dictionary-uniprot reaches
			// www.uniprot.org too, which is the older query interface rather than
			// rest.uniprot.org, so it is left out for the reason the MBTA's v2
			// clients were.
			//
			// The rest of npm is one viewer: ProtVista and its six adapters and data
			// loaders, which display UniProt annotations and take their entries from
			// whatever is handed to them. None carries a host.
			recipe: "uniprot",
			npm: []string{
				"uniprotparserjs",
			},
		},
		{
			// Nine of the twelve npm results are MCP servers, which is the openFDA
			// shape again: clinicaltrialsgov-mcp-server, @pipeworx/mcp-clinicaltrials,
			// bach-clinical-trials, @akshay89/clinical-trials-mcp,
			// mcp-server-clinical-trials, heor-agent-mcp, xiaoyibao-clinical-trials
			// and two more. One of them, @skills-il/israel-clinical-trials-mcp, is
			// about Israeli hospitals rather than this registry at all.
			//
			// The near miss is the retired interface. clinical-trials-gov -- "a way
			// to interface with the clinicaltrials.gov api" -- reaches
			// clinicaltrials.gov/ct2, the classic site this API replaced, so it is
			// left out for the reason the MBTA's realtime.mbta.com clients were.
			// clinicaltrials-gov-scraper is an Apify actor and carries no host.
			//
			// Packagist has three results and none of them is a client of this.
			// tutor/clinicaltrials is a "ClinicalTrials Web Service Guzzle Client",
			// which is the SOAP service from before the REST one, and its repository
			// cannot be read. pixiekat/lumina-ui-symfony-bundle and
			// pixiekat/hmfp-search-tool-bundle are Symfony bundles about evaluation
			// entities and search tooling.
			recipe: "clinicaltrialsgov",
			npm: []string{
				"clinicaltrialsgov-ts",
			},
		},
		{
			// The near miss is the vendor's other API. DataCite runs two:
			// api.datacite.org, which this Recipe describes, and mds.datacite.org,
			// the older Metadata Store used for minting. illgrenoble/
			// datacite-doi-bundle reaches mds.datacite.org and is left out, and
			// dagstuhl/datacite takes its host from the caller -- a providerUrl
			// constructor argument with {{ CREDENTIALS }} interpolated into it --
			// so there is no host in the package to map at all.
			//
			// linuskohl/phpdatacite is types without a client: 24 PHP files, every
			// one a model under src/models, and nothing that makes a request.
			//
			// The vendor's own namespace is not the vendor's client, four times
			// over. datacite-components, datacite-footer, @datacite/datacite-
			// tracker ("Frontend Tracker for DataCite Analytics") and @datacite/
			// mastiff (a server, not a client) are all DataCite's own front-end
			// work. The sharpest of them is datacite-rest, which is named after
			// this API and described, in full, as "React components".
			//
			// @pipeworx/mcp-datacite is an MCP server, by the same publisher as
			// the ClinicalTrials.gov one.
			//
			// The rest of Packagist is the popularity fallback on "data":
			// datacreativa/laravel-presentable, xaben/datafilter, velocity-sports-
			// labs/datacenter-php-sdk, datacodetech/phpcs-ruleset and
			// yiisoft/yii-dataview, none of which mentions DOIs.
			recipe: "datacite",
			composer: []string{
				"vincentauger/datacite-php-sdk",
			},
			npm: []string{
				"datacite-ts",
			},
		},
		{
			// The near miss is new: a client that reaches the right host, for the
			// one endpoint this Recipe does not have.
			//
			// Deezer's login providers are not Deezer's API clients, but they are
			// not innocent of it either. socialiteproviders/deezer reaches exactly
			// https://api.deezer.com/user/me and nothing else; julienbornstein/
			// oauth2-deezer reaches that plus a user's flow and image. Both are on
			// this host, and every path they use is under /user, which this Recipe
			// has no resource for. Offering it to them would emulate an API that
			// answers none of their requests. The same goes for passport-deezer,
			// chrishemmings/oauth2-deezer, magarrent/deezer-laravel-auth,
			// hwi/oauth-bundle and chillerlan/php-oauth-providers.
			//
			// The Discord playback plugins are left out as a family:
			// erela.js-deezer, @distube/deezer, kazagumo-deezer, stone-deezer,
			// discord-player-deezer and play-dl resolve a link through the public
			// API on their way to audio. A project holding one wants a music bot,
			// and the search it makes is incidental to that.
			//
			// @types/deezer-sdk is types without a client twice over: no request
			// anywhere in it, and the thing it describes is the browser JS SDK
			// rather than this API at all.
			recipe: "deezer",
			composer: []string{
				"atomescrochus/laravel-deezer-api",
				"bhanwarpsrathore/deezer-api-php",
				"clemfromspace/deezer-php-api",
				"deezer/deezer-php-sdk",
				"gabrieljmj/deezer-web-api",
				"pouler/deezer-api",
				"wasymshykh/deezer-api-php-wrapper",
			},
			npm: []string{
				"deezer-lib",
				"deezer-public-api",
				"deezer-ts-api",
				"deezer-web-api",
			},
		},
		{
			// The near miss is the dataset rather than the service.
			// node-red-contrib-oscar-plants-classification and its tensorflow
			// sibling "classify plant images among 10K species from the iNaturalist
			// dataset" -- a model trained on iNaturalist's photographs, which never
			// asks iNaturalist anything. The name is the corpus, not the API.
			//
			// Four of the twelve npm results are MCP servers:
			// @pondlog/mcp-inaturalist, inaturalist-mcp, @pipeworx/mcp-inaturalist
			// and @skills-il/israel-nature-mcp, the last of which is about Israeli
			// nature and wraps GBIF beside this. Both of those publishers have
			// turned up before, on openFDA and ClinicalTrials.gov.
			//
			// @citation-js/plugin-zenodo is the search engine reaching: it is a
			// Zenodo plugin and mentions iNaturalist nowhere.
			//
			// Packagist has no result for "inaturalist" at all.
			recipe: "inaturalist",
			npm: []string{
				"@pondlog/source-inaturalist",
				"@richard-stovall/inat-typescript-client",
				"inaturalistjs",
				"react-inaturalist-observation",
			},
		},
		{
			// The near miss is the vendor's own npm scope, holding something with
			// nothing to do with the API. @europepmc/express-middleware-minio is
			// "an Express middleware that helps you store files in Minio" -- the
			// organisation's internal object-storage plumbing, published under the
			// name a client would search for.
			//
			// The literature MCP servers are a crowd by now: @pipeworx/mcp-europepmc,
			// biomcp, dsh-pubmed, @cyanheads/pubmed-mcp-server, @pharmatools/
			// pubcrawl, @drvibeai/clinical-apis-mcp, lit-search and
			// @mengbingrock/labee-protocol-searcher all wrap this beside PubMed,
			// Semantic Scholar or bioRxiv, and a project holding one wants an agent
			// rather than an emulator of this host.
			//
			// Packagist has no result for "europepmc" at all.
			recipe: "europepmc",
			npm: []string{
				"biojs-vis-pmccitation",
				"epmc",
				"getpapers",
				"node-europmc",
			},
		},
		{
			// Almost everything the word returns on npm shells out to nuget.exe.
			// nuget-bin downloads the binary; grunt-nuget, grunt-nuget-pack,
			// gulp-nuget and gulp-nuget-restore drive it to pack, push and
			// restore. That is the line already drawn at the Sonar scanners and
			// at PyPI's publishers: something that runs a tool is not a client
			// of the read API.
			//
			// Two near-miss kinds recur, and both are now established rather
			// than noticed. @snyk/nuget-semver is "a semantic version parser for
			// NuGet" -- matched on a semantic, the second after @snyk/ruby-semver
			// and, oddly, from the same publisher. The @cratis/arc packages, at
			// fifteen thousand downloads each, carry an img.shields.io/nuget
			// badge in their READMEs -- matched on a badge, for the fourth
			// Recipe running.
			//
			// Packagist has nothing. melonsmasher/chocolatier is a NuGet
			// repository server rather than a client of one, and
			// activerules/nugget is a JSON validator spelled with two g's.
			//
			// snyk-nuget-plugin is mapped despite being CLI tooling, because
			// what it does is resolve dependency versions against this API --
			// the same reading that mapped sonarqube-verify while excluding the
			// scanners beside it.
			recipe: "nuget",
			npm: []string{
				"nuget",
				"snyk-nuget-plugin",
			},
		},
		{
			// The homonym here is the fixture. This Recipe is written about the
			// phoenix package on Hex, and the npm package called phoenix -- at
			// 5.3 million downloads, the biggest number on the page -- is the
			// official JavaScript client for the Phoenix framework's channels.
			// Same name, same project, and it speaks a websocket protocol to an
			// application rather than HTTP to a registry.
			//
			// The badge match turns up for the third Recipe running:
			// @mrdotb/live-react and @leuchtturm/turboprop both surface on a
			// shields.io/hexpm badge in their READMEs. It is now reliable enough
			// to expect on any provider whose users publish packages elsewhere.
			//
			// Packagist has nothing, and the results are all prefix collisions:
			// hexpang/ssh-client, hexmedia/yaml-linter and sizuhiko/hexpress
			// share four letters and no subject.
			//
			// What is left is four small TypeScript clients and a search plugin.
			// The obvious client of Hex is a Mix project, which is a registry
			// this file does not read -- the same structural thinness recorded
			// at PyPI and crates.io.
			recipe: "hexpm",
			npm: []string{
				"hex-api-client",
				"cerebro-hex",
				"hex-name",
				"@api-hooks/hex",
			},
		},
		{
			// The first registry in this run with real clients on its own
			// language's registry, and the numbers say why: spatie/packagist-api
			// at 2.4 million and knplabs/packagist-api at 1.5 million are both
			// what they say they are.
			//
			// The exclusion that took the most care is private-packagist
			// /api-client at eighty-seven thousand. Private Packagist is the
			// same company's commercial product -- a hosted private repository
			// at packagist.com with an API of its own -- so it is the
			// Hookdeck/Outpost shape a second time, and sharper, because here
			// the two products share a vendor prefix as well as a brand.
			// private-packagist/bitbucket-api is the same vendor again and is
			// entirely about Bitbucket.
			//
			// The badge match recorded at crates.io recurs and is now a
			// pattern rather than an observation: maxfactor-laravel-checkout,
			// angular-dashboard-framework and opengate-angular-dashboard-frame
			// all surface because their READMEs carry an
			// img.shields.io/packagist badge. Two of the three are Angular
			// projects with no PHP in them at all.
			//
			// khill/lavacharts at 2.1 million and square/connect at 1.4 million
			// are fuzzy matches with nothing to do with the registry, and
			// square/connect is doubly odd: it is a deprecated client for a
			// provider this collection already ships.
			recipe: "packagist",
			composer: []string{
				"spatie/packagist-api",
				"knplabs/packagist-api",
			},
			npm: []string{
				"@agonyz/packagist-api-client",
				"php-packagist-api-client",
				"gatsby-source-packagist",
				"alfred-packagist",
			},
		},
		{
			// Another new near-miss kind, and it owns the biggest number here:
			// matched on a semantic. @snyk/ruby-semver, at forty-two thousand
			// downloads, is "a node-semver compatible API with RubyGems
			// semantics" -- a JavaScript implementation of how RubyGems compares
			// versions. It knows more about this registry than any client here
			// does, and it never sends it a request.
			//
			// That is worth separating from a badge match, which the crates.io
			// entry records. A badge means the package is published somewhere; a
			// semantic means the package reimplements the provider's rules. Both
			// are honest mentions of a provider by something that does not call
			// it.
			//
			// The publishers go out on the line drawn at PyPI:
			// @webhippie/semantic-release-ruby and @tegami/gem push a built gem,
			// which is a different protocol with a token. And two MCP servers,
			// for the fifth Recipe running.
			//
			// Packagist has nothing at all. The word returns
			// judev/php-htmltruncator at three hundred thousand, ory/client and
			// duncan3dc/exec -- none of which mentions Ruby -- so the empty
			// result and the wrong result look identical, which is the reason to
			// write down that it is empty.
			//
			// bismar belongs here too and is mapped to PyPI, for the reason
			// recorded there.
			recipe: "rubygems",
			npm: []string{
				"rubygems",
				"gem-count",
				"cerebro-rubygems",
				"supply-chain-guard",
				"@agntn/registries",
				"@push.rocks/smartregistry",
			},
		},
		{
			// Two near-miss kinds at once, and one of them is new.
			//
			// The homonym is a single letter. Every Packagist result for
			// "crates" is CrateDB, a distributed SQL database that lives at
			// crate.io: crate/crate-pdo at ninety-six thousand,
			// crate/crate-dbal at forty-seven, ratkor/laravel-crate.io and
			// michabbb/php-crate. Not one of them has anything to do with the
			// Rust registry at crates.io, and the domain differs by an s.
			//
			// The new kind is matched on a badge. ts-gettext-extractor at
			// fifty-five thousand, @fltsci/tauri-plugin-tracing,
			// tauri-plugin-keyring-api, @fuel-bridge/fungible-token and
			// @tauri-liesel/reg all surface because their READMEs carry a
			// shields.io badge whose URL contains crates.io. The package is
			// usually a JavaScript binding for something written in Rust, so
			// the badge is honest and the match is not: the crate exists, and
			// nothing in the npm package calls the registry API.
			//
			// The publishers are excluded on the line drawn at PyPI:
			// @auto-canary/crates, @auto-it/crates, @releasekit/publish and
			// putitoutthere upload a built crate, which is cargo's protocol
			// rather than this API.
			//
			// What is left is the package named after the registry: "A crates.io
			// API client".
			//
			// bismar walks npm, crates.io and PyPI together and would belong
			// here as much as it belongs to PyPI, which is exactly what a
			// package mapping to one Recipe forbids -- the same invariant that
			// stopped jira.js naming two. It stays with PyPI because that is
			// where it was mapped first, and this note is the record that the
			// choice was arbitrary rather than considered.
			//
			// As with PyPI, the coverage is structurally thin: the obvious
			// client of crates.io is a Rust crate, and this file reads
			// Composer, npm and Go modules.
			recipe: "cratesio",
			npm:    []string{"crates.io"},
		},
		{
			// The largest number this file has ever excluded, and it is not even
			// a homonym: composer/installers, at a hundred and forty-eight
			// million downloads, is the top Packagist result for "pypi". It is a
			// multi-framework Composer library installer. The word appears
			// somewhere in its metadata and the search engine did the rest,
			// which is a useful reminder that a registry's ranking is not
			// evidence of anything.
			//
			// The real split here is between reading and publishing.
			// semantic-release-pypi at forty-two thousand, @lets-release/pypi,
			// foundry-release-pypi and putitoutthere all upload a built
			// distribution, which is twine's protocol and not this API. Same
			// line drawn at the Sonar scanners: something that happens to upload
			// is not a client of the read API.
			//
			// What reads it is small and genuine. pypi and pypi-info fetch
			// project metadata, pypi-axi inspects packages, and bismar and
			// muaddib-scanner walk several registries at once -- npm, crates.io,
			// PyPI -- which is exactly the shape a supply-chain tool has.
			//
			// Python packages are the obvious clients of PyPI and this file
			// reads Composer, npm and Go modules, so the coverage here is
			// structurally thin rather than accidentally so. That is worth
			// saying plainly: the Recipe is reachable, and the dependency that
			// would most often reach it is one this detector cannot see.
			recipe: "pypi",
			npm: []string{
				"pypi",
				"pypi-info",
				"pypi-axi",
				"bismar",
				"muaddib-scanner",
			},
		},
		{
			// The same split as ConfigCat, immediately afterwards, and by now it
			// is a rule about the category rather than a fact about a vendor: a
			// feature-flag platform has a delivery surface that SDKs read and a
			// management API that everything else uses, and they are different
			// APIs on different hosts.
			//
			// @growthbook/growthbook has 3.49 million downloads,
			// @growthbook/growthbook-react 1.65 million and the PHP SDK 3.92
			// million. All of them fetch feature definitions from
			// cdn.growthbook.io or from a proxy, and none calls the REST API
			// this Recipe serves. The edge apps, the Nuxt module, the Flags SDK
			// provider, the OpenFeature provider and the DevTools plugin all
			// wrap one of those SDKs, so they are on the same side.
			//
			// Exactly one package here speaks this API, and it says so in its
			// own description: growthbook, "Command-line interface for the
			// GrowthBook REST API", at a hundred and ten thousand downloads
			// against nine million on the other side.
			//
			// @growthbook/proxy is excluded as the product's own server
			// component, on the reading that put unleash-server outside: a
			// project running it is operating GrowthBook rather than
			// integrating with it. @growthbook/mcp is agent tooling, for the
			// fourth Recipe running.
			//
			// Unleash is the useful contrast. Its client and admin APIs share a
			// host, so one Recipe covers both and its SDKs are mapped. Where the
			// split is by host, as here and at ConfigCat, it is two Recipes or
			// none.
			recipe: "growthbook",
			npm:    []string{"growthbook"},
		},
		{
			// A new exclusion kind, and it owns the biggest number on the page
			// that is not the scanner: packages that write a file for another
			// program to upload.
			//
			// vitest-sonar-reporter has eight hundred and forty thousand
			// downloads and never makes a network request. Nor do
			// karma-sonarqube-reporter, cypress-sonarqube-reporter,
			// mocha-reporter-sonarqube, node-reporter-sonarqube or
			// symbiote/phpcs-sonar. They serialise test or lint results to disk
			// in a format the scanner later reads. A project holding one is
			// integrating with SonarCloud and is not calling it.
			//
			// The scanners are excluded for a neighbouring reason.
			// sonarqube-scanner at two million and @sonar/scan at seven hundred
			// thousand are bootstrappers: they download the sonar-scanner
			// Java CLI and run it, and the CLI submits an analysis over a
			// protocol this Recipe does not serve. Same line as the build
			// tooling excluded for Cloud Storage -- something that happens to
			// upload is not a client of the read API.
			//
			// sonarqube-verify is mapped despite living in the same
			// neighbourhood, because what it does after the analysis is poll
			// the Compute Engine task and read the quality gate, which is
			// exactly the surface here.
			//
			// The homonym is the largest number of all: nimut/phpunit-merger at
			// 1.89 million, the top Packagist result for "sonarqube", merges
			// PHPUnit reports. eslint-config-sonarqube is a lint configuration.
			//
			// What is left speaks the Web API: two community JavaScript clients
			// and two PHP ones, with about twelve thousand downloads between
			// them.
			recipe: "sonarcloud",
			composer: []string{
				"forgeqc/sonarqube-api-client",
				"spirit-dev/php-sonarqube-api",
			},
			npm: []string{
				"sonarqube-api-client",
				"sonarqube",
				"sonarqube-verify",
			},
		},
		{
			// The starkest version of a shape this file has now met three
			// Recipes running: the vendor's popular client speaks a different
			// surface, and mapping it would offer an emulator of the wrong API.
			//
			// configcat/configcat-client has 1.29 million downloads and
			// configcat-common seven hundred thousand, and neither calls
			// api.configcat.com. A ConfigCat SDK fetches a static configuration
			// file from the CDN and evaluates flags locally; the Public
			// Management API this Recipe serves is what dashboards and CI
			// scripts use. The clients that do speak it --
			// configcat-publicapi-node-client and ng-configcat-publicapi -- have
			// about five thousand downloads between them, so the ratio between
			// the two surfaces is roughly two hundred and sixty to one.
			//
			// That makes the config-delivery surface worth a Recipe of its own
			// rather than worth folding into this one. It is a different host,
			// a different shape and a different job.
			//
			// Everything else follows from that. The OpenFeature providers, the
			// Laravel, Spryker and Vue wrappers and the Belvo client all wrap
			// the SDK, so they are on the CDN side too. @pulumiverse/configcat
			// is infrastructure-as-code, for the sixth Recipe running. And
			// @configcat/mcp-server is the awkward one: it says outright that it
			// exposes the Public Management API, which is exactly this surface,
			// and it is still agent tooling rather than something an
			// application calls -- the same line drawn at @unleash/mcp and
			// @pinecone-database/mcp.
			recipe: "configcat",
			npm: []string{
				"configcat-publicapi-node-client",
				"ng-configcat-publicapi",
			},
		},
		{
			// The thinnest mapping in this file, and the reason is worth more
			// than the mapping: eighteen npm results for "grafana" and not one
			// of them calls this API.
			//
			// The entire @grafana/ scope is plugin-development tooling.
			// @grafana/data, /ui, /runtime, /schema, /scenes, /plugin-ui,
			// /e2e-selectors and /plugin-e2e are for building panels that run
			// inside Grafana, which reach its internals directly and never make
			// an HTTP request to it. A vendor scope full of the vendor's name
			// and empty of clients is the Hookdeck lesson at scale.
			//
			// On Packagist the largest numbers belong to a different product
			// from the same vendor. itspire/monolog-loki at nine hundred
			// thousand, cebe/yii2-loki-log-target, two Laravel Loki drivers and
			// a Craft one all speak to Loki, Grafana Labs' log database, whose
			// push and query API has nothing in common with this one. Beside
			// them sit Prometheus exporters -- iamfarhad/laravel-prometheus,
			// machour/yii2-prometheus, renoki-co/horizon-exporter -- which
			// publish metrics to be scraped and call nobody.
			//
			// The Foundation SDK is the near miss that took the most thought.
			// It builds dashboard documents as code, in Go, TypeScript, PHP and
			// more; it does not send them. A project holding it probably does
			// POST the result somewhere, but the package itself is a builder,
			// and mapping a builder would offer the Recipe on the strength of a
			// guess about the line after.
			//
			// @pulumiverse/grafana is infrastructure-as-code, for the fifth
			// Recipe running.
			//
			// What is left is two small PHP packages by one author. That is the
			// honest coverage for an API most people reach through the UI, a
			// provisioning file, or Terraform.
			recipe: "grafana",
			composer: []string{
				"saschahemleb/php-grafana-api-client",
				"saschahemleb/laravel-grafana",
			},
		},
		{
			// A small ecosystem where nearly everything under the vendor's own
			// scope is something other than a client, which makes the scope
			// useless as a signal and the exclusions the whole job.
			//
			// @hookdeck/outpost-sdk is the sharpest: Outpost is Hookdeck's
			// open-source, self-hosted outbound-webhook product, with an API of
			// its own. Same vendor, same npm scope, different surface -- so
			// mapping it here would offer a Recipe for the wrong API to a
			// project that is not calling this one at all. That is a new shape
			// for this file, and it is the reason a scope cannot be trusted the
			// way @googleapis/ can.
			//
			// The rest are the rules three Recipes have already settled.
			// @nectarsocial/hookdeck and @sst-provider/hookdeck are Pulumi
			// bridges and @toppy/hookdeck-deploy-cli deploys infrastructure
			// from declarative config, all infrastructure-as-code;
			// @stackcurious/hookdeck-mcp and @iflow-mcp/hookdeck-hookdeck-cli
			// are agent tooling; @hookdeck/eventcatalog-generator produces
			// documentation.
			//
			// The workflow nodes are mapped, because they really do call this
			// API on a project's behalf, which is the same reading that mapped
			// @pipedream/jira_service_desk.
			//
			// Packagist has no Hookdeck client at all. Searching for the word
			// returns laracasts/holodeck, ffm/hookdocs and avirdz/hookdeploy,
			// none of which shares it -- fuzzy-match noise rather than
			// homonyms, and worth recording only because an empty result and a
			// wrong result look the same from a distance.
			recipe: "hookdeck",
			npm: []string{
				"hookdeck-cli",
				"@hookdeck/vercel",
				"@hookdeck/pubsub",
				"n8n-nodes-hookdeck",
				"@hookdeck/n8n-nodes-hookdeck",
				"@pipedream/hookdeck",
			},
		},
		{
			// Two exclusions here have now happened three times running, which
			// makes them a rule rather than a judgement call.
			//
			// The first is infrastructure-as-code. @pinecone-database/pulumi
			// joins @pulumi/vault and @pulumiverse/unleash: a project holding
			// one is provisioning the provider rather than calling it at
			// runtime, so an emulator of the runtime API is not what its tests
			// need. The second is agent tooling: @pinecone-database/mcp beside
			// @unleash/mcp. Both are published by the vendor under the vendor's
			// own scope, which is exactly why the name alone cannot decide it.
			//
			// @traceloop/instrumentation-pinecone at eight hundred thousand is
			// a third kind, and a new one: it is OpenTelemetry instrumentation
			// that wraps the official SDK to emit traces. It observes a client
			// rather than being one, and any project holding it already holds
			// @pinecone-database/pinecone, which is what actually offers the
			// Recipe.
			//
			// The homonyms are small but real. The npm package called plainly
			// "pinecone" is a JavaScript-to-Lua converter, and create-pc-app
			// scaffolds "a Pinecone addon app" for something else entirely.
			// symfony/ai-store is the generic vector-store abstraction, so only
			// its Pinecone bridge is mapped.
			recipe: "pinecone",
			composer: []string{
				"probots-io/pinecone-php",
				"scotteuser/pinecone-php",
				"symfony/ai-pinecone-store",
				"islambaraka90/pinecone-php-client",
				"mbvb1223/pinecone-php-client",
			},
			npm: []string{
				"@pinecone-database/pinecone",
				"@pinecone-database/connect",
				"@langchain/pinecone",
				"@mastra/pinecone",
				"genkitx-pinecone",
			},
		},
		{
			// Another near-miss shape, and again the biggest number on the page
			// produced it: matched on an imperative. "Unleash" is a verb that
			// package taglines love, so the largest npm result for it is
			// precinct at twenty million downloads -- a JavaScript dependency
			// parser whose description is "Unleash the detectives". Beside it,
			// @remirror/extension-code-block "Unleash the inner coder",
			// veewee/reflecta "Unleash the Power of Optics in your code!", and
			// an npm package literally called unleash: "Unleash your code into
			// the wild yonder". None of them is a feature flag.
			//
			// unleashedtech/php-coding-standard is a different vendor holding
			// the same word -- a CodeSniffer ruleset from a company called
			// Unleashed Technologies.
			//
			// Three exclusions are deliberate. unleash-server and
			// unleash-frontend are the product: a project depending on those is
			// running Unleash rather than calling it, and an emulator is not
			// what it needs. @pulumiverse/unleash provisions flags rather than
			// reading them, on the same line drawn for @pulumi/vault. And
			// @unleash/mcp is an agent tool for managing flags rather than an
			// SDK an application evaluates them with.
			//
			// One mapping is worth its own note: react-unleash-flags describes
			// itself as a component "for Unleash or Gitlab Feature Flags",
			// because GitLab implements Unleash's client API. A project holding
			// it may be pointing at GitLab and still speaking this protocol,
			// which is a reason to offer the Recipe rather than a reason not
			// to.
			recipe: "unleash",
			composer: []string{
				"unleash/client",
				"unleash/symfony-client-bundle",
				"j-webb/laravel-unleash",
				"mikefrancis/laravel-unleash",
				"huddlebv/laravel-unleash",
				"stogon/unleash-bundle",
				"bonch.dev/laravel-unleash",
				"emersonecologics/unleash-client-php",
			},
			npm: []string{
				"unleash-client",
				"unleash-proxy-client",
				"@unleash/proxy-client-react",
				"@unleash/proxy-client-vue",
				"@unleash/proxy-client-svelte",
				"@unleash/nextjs",
				"@unleash/unleash-react-native-sdk",
				"@unleash/proxy",
				"nestjs-unleash",
				"react-unleash-flags",
			},
		},
		{
			// The most contested name in this file. Four vendors ship a product
			// called Vault, and the biggest one is not this one:
			// @azure/keyvault-common has twenty-three million downloads and
			// @azure/keyvault-secrets nearly ten, against node-vault's 1.7
			// million. Oracle ships oci-vault, and @googleapis/vault is Google
			// Vault, which is e-discovery for Workspace and holds no secrets
			// at all.
			//
			// Then the words that are not products. spryker/vault at 2.2
			// million is a module of the Spryker commerce framework.
			// soarecostin/file-vault and its three relatives encrypt files.
			// ansible-vault is a file format. dotenv-org/phpdotenv-vault reads
			// .env.vault files. The DeFi packages are yield vaults.
			// @apideck/vault-react is the settings UI of a provider already in
			// this collection. And the top Packagist result for the word is
			// ccxt/ccxt, a cryptocurrency exchange library.
			//
			// Two exclusions are deliberate rather than accidental.
			// @testcontainers/vault starts a real Vault in Docker for tests,
			// which is the thing somebody reaches for instead of an emulator --
			// mapping it would offer this Recipe to the one project that has
			// already decided against it. And @pulumi/vault provisions Vault
			// rather than reading from it: infrastructure-as-code calls the API
			// at deploy time, but a project holding it is configuring the
			// server, not fetching secrets from it at runtime, which is what
			// this Recipe serves.
			//
			// csharpru/vault-php-guzzle6-transport is mapped even though it is
			// a transport rather than a client, because it exists only to carry
			// csharpru/vault-php and four million downloads say so.
			recipe: "vault",
			composer: []string{
				"csharpru/vault-php",
				"csharpru/vault-php-guzzle6-transport",
				"jippi/vault-php-sdk",
				"violuke/vault-php-sdk",
				"mittwald/vault-php",
				"tokenly/laravel-vault",
				"lucasberto/laravel-vault",
				"itk-dev/vault",
				"itk-dev/vault-bundle",
			},
			npm: []string{
				"node-vault",
				"node-vault-client",
				"hashi-vault-js",
				"@litehex/node-vault",
				"vault-api",
				"vault-auth-aws",
			},
		},
		{
			// A new shape, and the biggest number on the page produced it:
			// matched on a simile. Packagist's top results for "youtube" are
			// hashids/hashids at fifty-four million downloads and
			// vinkla/hashids at fifteen, and sqids/sqids says outright what
			// all three are -- "Generate short YouTube-looking IDs from
			// numbers". They do not share the name, the vendor or the
			// domain. They describe their output as looking like something
			// YouTube makes.
			//
			// And the second exclusion is the tool people reach for instead
			// of this API. norkunas/youtube-dl-php wraps youtube-dl and
			// yt-dlp, which get at videos by scraping the site rather than by
			// asking the Data API, so a project holding it is doing the thing
			// this API deliberately does not offer.
			//
			// The client follows the scoped pattern Gmail, Calendar and Drive
			// already use here. npm has no widely installed YouTube client
			// otherwise: the largest is simple-youtube-api at seven thousand
			// a month, which is a fact about the API rather than a gap in the
			// search -- most people who touch YouTube from a browser use the
			// iframe player, and most who touch it from a server use the
			// official client under its scope.
			recipe: "youtube",
			composer: []string{
				"alaouy/youtube",
				"madcoda/php-youtube-api",
			},
			npm: []string{
				"@googleapis/youtube",
				"simple-youtube-api",
				"youtube-node",
			},
		},
		{
			// The largest numbers in this file. google/cloud-storage has a
			// hundred and five million downloads and @google-cloud/storage
			// sixty-one million, and the Flysystem and Laravel adapters
			// between them come to nearly sixty million more.
			//
			// Those adapters are mapped, because they call this API -- and it
			// is worth noticing, as it was on Drive, that the dominant use of
			// an object store is to make it look like a filesystem. Here the
			// tension is sharper: GCS's own documentation says directory-like
			// mode is a pretence, with the folders arriving in a separate
			// array from the files.
			//
			// google/cloud is the umbrella client for every Google Cloud
			// product at once, so it stays unmapped for the reason jira.js
			// does: a package maps to one Recipe and that one speaks dozens.
			//
			// And mock-gcs is "a Mock implementation of the Google Cloud
			// Storage" -- the second emulator this file has had to exclude in
			// two Recipes, after google-drive-mock. Object stores attract
			// them, presumably because everybody wants their tests to stop
			// touching the network and nobody wants to run a bucket.
			//
			// The upload sinks are left alone: multer-cloud-storage,
			// @tus/gcs-store, @payloadcms/storage-gcs and the Strapi provider
			// write through the official client rather than being one, and
			// webpack-google-cloud-storage-plugin, nx-remotecache-gcs and
			// reg-publish-gcs-plugin are build tooling that happens to
			// upload.
			recipe: "googlecloudstorage",
			composer: []string{
				"google/cloud-storage",
				"league/flysystem-google-cloud-storage",
				"superbalist/flysystem-google-storage",
				"spatie/laravel-google-cloud-storage",
				"superbalist/laravel-google-cloud-storage",
				"spatie/flysystem-google-cloud-storage",
			},
			npm: []string{"@google-cloud/storage"},
		},
		{
			// Gemini is also a cryptocurrency exchange, which is the homonym
			// shape Wix brought and this is the second instance of: ccxt/ccxt
			// is "a cryptocurrency trading API" supporting a hundred-odd
			// venues and Gemini is one of them, so a trading bot matches the
			// word without ever meeting a model.
			//
			// The deployment shape again, and a large one.
			// @google-cloud/vertexai is the same models reached through
			// Vertex AI -- a different host, a different credential and
			// service-account auth rather than an API key -- installed nearly
			// four million times a month. That is @azure/openai one provider
			// along.
			//
			// The vendor's own CLI outweighs whole ecosystems: @google
			// /gemini-cli and @google/gemini-cli-core are a coding agent, not
			// a client, and between them they are installed more than every
			// PHP client here put together. @openai/codex is the same shape.
			//
			// And two official SDKs coexist at enormous scale:
			// @google/generative-ai at sixteen million a month is superseded
			// by @google/genai at seventy-four, so both are mapped -- the
			// older name is what a great many projects are still holding,
			// which is the rename Dropbox Sign's entry records at a smaller
			// size.
			//
			// @ibm-generative-ai/node-sdk is a different vendor that matches
			// on the phrase rather than the name, and Vercel's "ai" is a
			// generic SDK that mentions Gemini among many.
			recipe: "gemini",
			composer: []string{
				"google-gemini-php/client",
				"google-gemini-php/laravel",
				"gemini-api-php/client",
			},
			npm: []string{
				"@google/genai",
				"@google/generative-ai",
				"@google-ai/generativelanguage",
				"@langchain/google-genai",
			},
		},
		{
			// The scoped package, as Gmail and Calendar already do here: the
			// unscoped googleapis covers every Google API at once and a
			// package maps to one Recipe, so @googleapis/drive is the one
			// that means Drive and nothing else.
			//
			// The most installed Drive packages on Packagist are filesystem
			// adapters -- masbug/flysystem-google-drive-ext and nao-pon
			// /flysystem-google-drive, four and a half million downloads
			// between them -- and they are mapped, because they really do
			// call this API. It is worth noticing that their whole purpose is
			// to make Drive look like a filesystem, on an API whose own
			// documentation says a file's name "isn't necessarily unique
			// within a folder".
			//
			// Everything else with the name is a picker. @uppy/google-drive,
			// react-google-drive-picker, @googleworkspace/drive-picker
			// -element and -react, google-drive-picker and
			// @fyelci/react-google-drive-picker are between them installed
			// more than two and a half million times a month, and all of them
			// are browser widgets that let a person choose a file through the
			// Picker API rather than servers calling this one.
			//
			// Two more worth naming. @maxim_mazurok/gapi.client.drive-v3 is
			// third-party typings and makes no request. And google-drive-mock
			// is "a Mock-Server that simulates being google-drive" -- another
			// emulator, doing this Recipe's job, which is the first time
			// anything in this file has collided with the project itself.
			//
			// spatie/laravel-google-cloud-storage, at nearly eleven million
			// downloads, is Google Cloud Storage: a different Google product
			// that shares three quarters of its name with nothing here.
			recipe: "googledrive",
			composer: []string{
				"masbug/flysystem-google-drive-ext",
				"nao-pon/flysystem-google-drive",
				"yaza/laravel-google-drive-storage",
			},
			npm: []string{"@googleapis/drive"},
		},
		{
			// Atlas is an everyday word, and one of the collisions is a
			// prefix nobody would predict: Atlassian begins with it.
			// damienharper/adf-tools is "Atlassian Document Format PHP
			// Tools", installed one point eight million times a month, and it
			// is about Confluence and Jira rather than about databases.
			//
			// The everyday word does the rest. atlas/pdo, atlas/query and
			// atlas/orm are a PHP persistence family under the vendor name
			// atlas, and between them they outnumber every real client here.
			// atlas-php/atlas is an AI SDK for Laravel. grazulex/laravel
			// -atlas maps a Laravel project. Packagist has no client for this
			// API at all.
			//
			// And MongoDB's own scope contains the deployment shape again:
			// @mongodb-js/atlas-local manages "MongoDB Atlas Local
			// deployments", which run in a container rather than in the
			// cloud, so the administration API this Recipe models is not
			// involved. mongodb-data-api speaks the Data API, a different
			// interface for reading documents rather than administering
			// clusters, and mongoose-atlas-search reaches Atlas Search
			// through the database driver.
			recipe: "mongodbatlas",
			npm: []string{
				"mongodb-atlas-api-client",
				"awscdk-resources-mongodbatlas",
				"mcp-mongodb-atlas",
			},
		},
		{
			// Keycloak's entry below names a shape -- the other half of the
			// same product -- and FusionAuth confirms it is a property of
			// identity providers rather than of one vendor. Everything
			// popular here is about getting a token or checking one, and the
			// administration API is the smaller half again.
			//
			// socialiteproviders/fusionauth is the Socialite shape for the
			// fifth time in this file. jerryhopper/oauth2-fusionauth signs
			// people in for the PHP League. danilopolani/laravel-fusionauth
			// -jwt, werk365/jwtauthroles and fusionauth/jwt-auth-webtoken
			// -provider all validate a token FusionAuth issued -- and the
			// last of those is under the vendor's own namespace, which is
			// what makes the vendor prefix useless as a signal here.
			//
			// @fusionauth/react-sdk, angular-sdk and vue-sdk log people in
			// from a browser, and passport-fusionauth and nodebb-plugin
			// -fusionauth-oidc are authentication strategies.
			//
			// nodebb-theme-fusionauth is a forum theme. It is a stylesheet.
			recipe:   "fusionauth",
			composer: []string{"fusionauth/fusionauth-client"},
			npm: []string{
				"@fusionauth/typescript-client",
				"@fusionauth/cli",
				"fusionauth-cli",
				"@fusionauth/mcp-api",
				"pulumi-fusionauth",
			},
		},
		{
			// The other half of the same product, which is a new shape here.
			// Not a different deployment and not a different transport: the
			// thing almost everybody uses Keycloak for is getting a token or
			// checking one, and the Admin REST API this Recipe models is the
			// side door.
			//
			// stevenmaguire/oauth2-keycloak has six and three quarter million
			// downloads and signs people in. socialiteproviders/keycloak has
			// two and a half million and is the Socialite shape this file has
			// now recorded four times. robsontenorio/laravel-keycloak-guard,
			// vizir/laravel-keycloak-web-guard and two more like them
			// *validate* a token Keycloak issued and never call the admin API
			// at all.
			//
			// npm is the same: keycloak-connect, keycloak-angular, the
			// @react-keycloak packages, nest-keycloak-connect and
			// passport-keycloak-bearer are adapters that authenticate people,
			// and the bare "keycloak" is the browser adapter.
			//
			// keycloakify is the sharpest of them -- "Framework to create
			// custom Keycloak UIs", four hundred and thirty thousand installs
			// of a build tool that produces login pages and makes no request
			// of any kind. cypress-keycloak logs a test user in.
			//
			// What is mapped is the administration side, and it is much the
			// smaller half: the PHP admin client has seven hundred and ninety
			// thousand downloads against the OAuth provider's six and three
			// quarter million.
			recipe: "keycloak",
			composer: []string{
				"mohammad-waleed/keycloak-admin-client",
				"fschmtt/keycloak-rest-api-client-php",
			},
			npm: []string{
				"@keycloak/keycloak-admin-client",
				"@s3pweb/keycloak-admin-client-cjs",
				"@pulumi/keycloak",
			},
		},
		{
			// The vendor's own unrelated library, at a scale nothing else in
			// this file approaches -- and it is the second time the unrelated
			// library is a PHP static analyser.
			//
			// vimeo/psalm has eighty-five million downloads. It finds errors
			// in PHP and has nothing whatever to do with video. Etsy's phan
			// is the same tool doing the same job under the same kind of
			// vendor name, so two providers here now have static analysis as
			// their largest near miss. vimeo/php-mysql-engine is a MySQL
			// engine written in PHP, from the same vendor and the same
			// distance from this API.
			//
			// And the player, which is the transport shape at twenty to one.
			// @vimeo/player -- "Interact with and control an embedded Vimeo
			// Player" -- is installed five and a quarter million times a
			// month against the official Node client's two hundred and sixty
			// thousand. It drives an iframe in a browser and never calls this
			// API. vimeo-video-element, @u-wave/react-vimeo, vue-vimeo
			// -player, react-native-vimeo-iframe, lite-vimeo-embed and
			// @astro-community/astro-embed-vimeo are all the same thing
			// again, and @types/vimeo__player alone outnumbers the client six
			// to one.
			//
			// @snowplow/browser-plugin-vimeo-tracking watches the player
			// rather than calling anything, and plyr is a general HTML5
			// player that merely mentions Vimeo in its description.
			recipe: "vimeo",
			composer: []string{
				"vimeo/vimeo-api",
				"vimeo/laravel",
			},
			npm: []string{"@vimeo/vimeo"},
		},
		{
			// A rename the ecosystem has not followed. HelloSign became
			// Dropbox Sign, and on Packagist the old package is installed
			// seven times more than the new one -- hellosign/hellosign-php
			// -sdk against dropbox/sign -- so both are mapped and the older
			// name is the one most projects will actually be holding.
			//
			// The sibling product, across a wider gulf than Jira and
			// Confluence. The npm package "dropbox" is the file storage SDK
			// and @uppy/dropbox is an uploader for it: same company, same
			// brand, nothing to do with signatures. Nothing but the package
			// name connects them to this API.
			//
			// And the transport shape again, more installed than the client
			// it sits beside: hellosign-embedded "embeds HelloSign signature
			// requests and templates from within your app", which is an
			// iframe in a browser fed a sign_url that some server already
			// fetched. It makes no call to this API and is installed eight
			// hundred thousand times a month, against the official Node
			// client's four hundred and sixty.
			//
			// hellosign/hawk and lnn/hellowight are Packagist noise that
			// merely begins with the word.
			recipe: "dropboxsign",
			composer: []string{
				"dropbox/sign",
				"hellosign/hellosign-php-sdk",
				"industrious/hellosign-laravel",
				"bukashk0zzz/hellosign-bundle",
				"o0khoiclub0o/hellosign-php",
			},
			npm: []string{
				"@dropbox/sign",
				"@memberjunction/esignature-dropboxsign",
				"@molecule/api-esign-hellosign",
			},
		},
		{
			// The clearest instance yet of a shape ElevenLabs introduced: the
			// vendor's most installed package is for a different protocol.
			//
			// @planetscale/database is "A Fetch API-compatible PlanetScale
			// database driver", installed a million times a month, and it
			// speaks the serverless query protocol -- it runs SQL. This
			// Recipe is the management API, which creates databases and
			// merges deploy requests. Same company, same account, same
			// credential family, and no endpoint in common.
			//
			// Everything built on the driver inherits the exclusion:
			// @prisma/adapter-planetscale, kysely-planetscale,
			// @mattrax/mysql-planetscale, @storecraft/database-planetscale
			// and planetscale-stream-ts are all query-path packages, and
			// between them they outnumber every management client here by two
			// orders of magnitude.
			//
			// The bare npm "planetscale" is "a simple client for connecting
			// to PlanetScale" -- connecting, which is the driver's job -- and
			// is installed two thousand times a month, so it is left alone on
			// the Marqeta reading rather than guessed at.
			//
			// What is mapped is the management side: an SDK generated from
			// this very document, a Pulumi provider bridged from it, and two
			// PHP packages that exist to move migrations through branches and
			// deploy requests.
			recipe: "planetscale",
			composer: []string{
				"x7media/laravel-planetscale",
				"bellows-app/plugin-planetscale",
			},
			npm: []string{
				"@distilled.cloud/planetscale",
				"@sst-provider/planetscale",
			},
		},
		{
			// A new shape: the vendor's own most installed packages are for a
			// different transport.
			//
			// @elevenlabs/react and @elevenlabs/client are between them
			// installed five million times a month, and both are for the
			// Agents platform -- a browser holding a conversation over a
			// socket, not a server calling this REST API. @elevenlabs
			// /react-native says so in its own description. Squarespace's
			// @squarespace/* and Wix's @atlaskit-alikes were the vendor
			// publishing a *front end*; this is the vendor publishing a
			// different protocol for the same product.
			//
			// @elevenlabs/types is sharper still: "AsyncAPI contracts and
			// generated TypeScript types". Not OpenAPI -- AsyncAPI, the
			// event-driven description -- and types alone, so it makes no
			// request of any kind. Two and a half million installs a month of
			// a package generated from a different specification format for a
			// different half of the product.
			//
			// @remotion/elevenlabs works with "the output of the ElevenLabs
			// API" rather than calling it, and @elevenlabs/cli, the n8n nodes
			// and the LiveKit and Mastra plugins are integrations for other
			// platforms.
			recipe: "elevenlabs",
			composer: []string{
				"ardagnsrn/elevenlabs-laravel",
				"georgehadjisavva/elevenlabs-api-client",
				"onramplab/elevenlabs-api-client",
				"gridwb/laravel-elevenlabs",
				"samandar/laravel-elevenlabs",
			},
			npm: []string{
				"@elevenlabs/elevenlabs-js",
				"@ai-sdk/elevenlabs",
				"elevenlabs-js",
			},
		},
		{
			// No Recipe ships for Anthropic, and the mapping is here on
			// purpose -- the same reason OpenAI's was, before OpenAI's became
			// a Recipe. A dependency Cauldron recognises and cannot emulate
			// is reported by name, so the developer is told those calls will
			// still reach the real network. Without the mapping it falls
			// through to the "looks like an API client" heuristic and the
			// warning is vaguer than it needs to be.
			//
			// @anthropic-ai/bedrock-sdk is excluded: it is "the official
			// TypeScript library for the Anthropic Bedrock API", which is the
			// same models addressed through AWS rather than through
			// api.anthropic.com. That is the deployment shape @azure/openai
			// has, one provider along.
			recipe:   "anthropic",
			npm:      []string{"@anthropic-ai/sdk"},
			composer: []string{"anthropic-ai/sdk", "mozex/anthropic-php"},
			gomod:    []string{"github.com/anthropics/anthropic-sdk-go"},
		},
		{
			// This mapping predates the Recipe: it was here so that a
			// dependency Cauldron recognised and could not emulate was
			// reported by name rather than falling through to the "looks like
			// an API client" heuristic. A Recipe ships now, and the mapping
			// has grown the exclusions that go with the biggest ecosystem in
			// this file.
			//
			// A shape that has outgrown its provider.
			// @ai-sdk/openai-compatible, at twenty-three million installs a
			// month, exists "to provide a foundation for implementing
			// providers" that speak OpenAI's request format -- and
			// deepseek-php/deepseek-php-client is one of the providers it
			// means. A project holding either may never send a byte to
			// OpenAI. No shape recorded here covered that: Wix was one word
			// meaning two products and Metronome was one word meaning an
			// instrument, where this is one company's wire format used by
			// everybody else.
			//
			// Another deployment, again. @azure/openai and azure-openai talk
			// to {resource}.openai.azure.com, which addresses deployments
			// rather than models and versions itself with an api-version
			// query parameter. Same models, different API -- the shape
			// @atlassian-dc-mcp/confluence has.
			//
			// The vendor's own biggest package is not a client.
			// @openai/codex is a coding agent that runs on your machine, and
			// it is installed sixty-seven million times a month.
			//
			// And the sharpest miss is the subject of one of this Recipe's
			// own claims: yethee/tiktoken is OpenAI's tokeniser ported to
			// PHP, four and a half million installs of a library that counts
			// tokens and never makes a request -- which is exactly the
			// arithmetic the Recipe says will disagree with usage.
			//
			// Instrumentation is excluded too: @opentelemetry/instrumentation
			// -openai, @traceloop/instrumentation-openai and @langfuse/openai
			// wrap a client rather than being one, and a project holding any
			// of them holds the real client as well.
			recipe: "openai",
			composer: []string{
				"openai-php/client",
				"openai-php/laravel",
				"openai-php/symfony",
				"orhanerday/open-ai",
				"tectalic/openai",
			},
			npm: []string{
				"openai",
				"@ai-sdk/openai",
				"@langchain/openai",
			},
			gomod: []string{"github.com/sashabaranov/go-openai"},
		},
		{
			recipe:   "ably",
			npm:      []string{"ably"},
			composer: []string{"ably/ably-php"},
			gomod:    []string{"github.com/ably/ably-go"},
		},
		{
			recipe: "airtable",
			npm:    []string{"airtable"},
		},
		{
			recipe:   "algolia",
			npm:      []string{"algoliasearch"},
			composer: []string{"algolia/algoliasearch-client-php"},
			gomod:    []string{"github.com/algolia/algoliasearch-client-go"},
		},
		{
			recipe:   "asana",
			npm:      []string{"asana"},
			composer: []string{"asana/asana"},
		},
		{
			recipe: "assemblyai",
			npm:    []string{"assemblyai"},
			gomod:  []string{"github.com/AssemblyAI/assemblyai-go-sdk"},
		},
		{
			recipe:   "auth0",
			npm:      []string{"auth0", "@auth0/auth0-spa-js", "express-openid-connect"},
			composer: []string{"auth0/auth0-php"},
			gomod:    []string{"github.com/auth0/go-auth0"},
		},
		{
			recipe: "bitbucket",
			npm:    []string{"bitbucket"},
		},
		{
			recipe: "box",
			npm:    []string{"box-node-sdk", "box-typescript-sdk-gen"},
		},
		{
			recipe: "brex",
			npm:    []string{"@brexhq/protocol"},
		},
		{
			recipe:   "bugsnag",
			npm:      []string{"@bugsnag/js", "@bugsnag/node"},
			composer: []string{"bugsnag/bugsnag-laravel", "bugsnag/bugsnag"},
			gomod:    []string{"github.com/bugsnag/bugsnag-go"},
		},
		{
			recipe: "calendly",
			npm:    []string{"react-calendly"},
		},
		{
			recipe:   "chargebee",
			npm:      []string{"chargebee"},
			composer: []string{"chargebee/chargebee-php"},
			gomod:    []string{"github.com/chargebee/chargebee-go"},
		},
		{
			recipe: "circleci",
			npm:    []string{"circleci-api"},
		},
		{
			recipe: "clickup",
			npm:    []string{"clickup.js"},
		},
		{
			recipe:   "cloudflare",
			npm:      []string{"cloudflare"},
			composer: []string{"cloudflare/sdk"},
			gomod:    []string{"github.com/cloudflare/cloudflare-go"},
		},
		{
			recipe:   "cloudinary",
			npm:      []string{"cloudinary"},
			composer: []string{"cloudinary/cloudinary_php"},
			gomod:    []string{"github.com/cloudinary/cloudinary-go"},
		},
		{
			recipe:   "contentful",
			npm:      []string{"contentful", "contentful-management"},
			composer: []string{"contentful/contentful-management"},
			gomod:    []string{"github.com/labd/contentful-go"},
		},
		{
			recipe:   "datadog",
			npm:      []string{"@datadog/datadog-api-client", "dd-trace"},
			composer: []string{"datadog/php-datadogstatsd"},
			gomod:    []string{"github.com/DataDog/datadog-api-client-go"},
		},
		{
			recipe: "deepgram",
			npm:    []string{"@deepgram/sdk"},
			gomod:  []string{"github.com/deepgram/deepgram-go-sdk"},
		},
		{
			recipe: "digitalocean",
			npm:    []string{"do-wrapper"},
			gomod:  []string{"github.com/digitalocean/godo"},
		},
		{
			recipe:   "discord",
			npm:      []string{"discord.js", "@discordjs/rest"},
			composer: []string{"team-reflex/discord-php"},
			gomod:    []string{"github.com/bwmarrin/discordgo"},
		},
		{
			recipe:   "docusign",
			npm:      []string{"docusign-esign"},
			composer: []string{"docusign/esign-client"},
		},
		{
			recipe:   "dropbox",
			npm:      []string{"dropbox"},
			composer: []string{"spatie/flysystem-dropbox"},
		},
		{
			recipe: "dynamodb",
			npm:    []string{"@aws-sdk/client-dynamodb", "@aws-sdk/lib-dynamodb", "dynamoose"},
			gomod:  []string{"github.com/aws/aws-sdk-go-v2/service/dynamodb"},
		},
		{
			recipe: "firecrawl",
			npm:    []string{"@mendable/firecrawl-js"},
		},
		{
			recipe: "freshdesk",
			npm:    []string{"freshdesk-api"},
		},
		{
			recipe: "front",
			npm:    []string{"@frontapp/plugin-sdk"},
		},
		{
			recipe: "ghost",
			npm:    []string{"@tryghost/content-api", "@tryghost/admin-api"},
		},
		{
			recipe:   "gitlab",
			npm:      []string{"@gitbeaker/rest", "@gitbeaker/node"},
			composer: []string{"m4tthumphrey/php-gitlab-api"},
			gomod:    []string{"github.com/xanzy/go-gitlab"},
		},
		{
			recipe:   "gmail",
			npm:      []string{"@googleapis/gmail"},
			composer: []string{"google/apiclient"},
			gomod:    []string{"google.golang.org/api/gmail/v1"},
		},
		{
			recipe: "googlecalendar",
			npm:    []string{"@googleapis/calendar"},
			gomod:  []string{"google.golang.org/api/calendar/v3"},
		},
		{
			recipe:   "intercom",
			npm:      []string{"intercom-client"},
			composer: []string{"intercom/intercom-php"},
			gomod:    []string{"github.com/intercom/intercom-go"},
		},
		{
			recipe:   "jira",
			npm:      []string{"jira.js", "jira-client"},
			composer: []string{"lesstif/php-jira-rest-client"},
			gomod:    []string{"github.com/andygrunwald/go-jira"},
		},
		{
			recipe:   "klaviyo",
			npm:      []string{"klaviyo-api"},
			composer: []string{"klaviyo/api"},
		},
		{
			recipe:   "launchdarkly",
			npm:      []string{"launchdarkly-node-server-sdk", "@launchdarkly/node-server-sdk"},
			composer: []string{"launchdarkly/server-sdk"},
			gomod:    []string{"github.com/launchdarkly/go-server-sdk"},
		},
		{
			recipe:   "lob",
			npm:      []string{"lob"},
			composer: []string{"lob/lob-php"},
		},
		{
			recipe:   "mailchimp",
			npm:      []string{"@mailchimp/mailchimp_marketing", "@mailchimp/mailchimp_transactional"},
			composer: []string{"mailchimp/marketing"},
		},
		{
			recipe: "miro",
			npm:    []string{"@mirohq/miro-api"},
		},
		{
			recipe:   "mollie",
			npm:      []string{"@mollie/api-client"},
			composer: []string{"mollie/mollie-api-php", "mollie/laravel-mollie"},
		},
		{
			recipe:   "mux",
			npm:      []string{"@mux/mux-node"},
			composer: []string{"muxinc/mux-php"},
			gomod:    []string{"github.com/muxinc/mux-go"},
		},
		{
			recipe: "netlify",
			npm:    []string{"netlify"},
		},
		{
			recipe: "notion",
			npm:    []string{"@notionhq/client"},
			gomod:  []string{"github.com/jomei/notionapi"},
		},
		{
			recipe: "okta",
			npm:    []string{"@okta/okta-sdk-nodejs"},
			gomod:  []string{"github.com/okta/okta-sdk-golang"},
		},
		{
			recipe:   "onesignal",
			npm:      []string{"@onesignal/node-onesignal"},
			composer: []string{"onesignal/onesignal-php-api"},
		},
		{
			recipe:   "paddle",
			npm:      []string{"@paddle/paddle-node-sdk"},
			composer: []string{"laravel/cashier-paddle"},
		},
		{
			recipe: "pagerduty",
			npm:    []string{"@pagerduty/pdjs"},
			gomod:  []string{"github.com/PagerDuty/go-pagerduty"},
		},
		{
			recipe:   "paypal",
			npm:      []string{"@paypal/paypal-server-sdk", "@paypal/checkout-server-sdk"},
			composer: []string{"paypal/paypal-checkout-sdk", "paypal/rest-api-sdk-php"},
		},
		{
			recipe: "pipedrive",
			npm:    []string{"pipedrive"},
		},
		{
			recipe:   "plaid",
			npm:      []string{"plaid"},
			composer: []string{"plaid/plaid-php"},
			gomod:    []string{"github.com/plaid/plaid-go"},
		},
		{
			recipe:   "posthog",
			npm:      []string{"posthog-node", "posthog-js"},
			composer: []string{"posthog/posthog-php"},
			gomod:    []string{"github.com/posthog/posthog-go"},
		},
		{
			recipe:   "postmark",
			npm:      []string{"postmark"},
			composer: []string{"wildbit/postmark-php"},
		},
		{
			recipe:   "pubsub",
			npm:      []string{"@google-cloud/pubsub"},
			composer: []string{"google/cloud-pubsub"},
			gomod:    []string{"cloud.google.com/go/pubsub"},
		},
		{
			recipe:   "pusher",
			npm:      []string{"pusher", "pusher-js"},
			composer: []string{"pusher/pusher-php-server"},
			gomod:    []string{"github.com/pusher/pusher-http-go"},
		},
		{
			recipe:   "quickbooks",
			npm:      []string{"node-quickbooks", "intuit-oauth"},
			composer: []string{"quickbooks/v3-php-sdk"},
		},
		{
			recipe:   "recurly",
			npm:      []string{"recurly"},
			composer: []string{"recurly/recurly-client"},
			gomod:    []string{"github.com/recurly/recurly-client-go"},
		},
		{
			recipe:   "resend",
			npm:      []string{"resend"},
			composer: []string{"resend/resend-php", "resend/resend-laravel"},
			gomod:    []string{"github.com/resend/resend-go"},
		},
		{
			recipe:   "ringcentral",
			npm:      []string{"@ringcentral/sdk"},
			composer: []string{"ringcentral/ringcentral-php"},
		},
		{
			recipe:   "rollbar",
			npm:      []string{"rollbar"},
			composer: []string{"rollbar/rollbar", "rollbar/rollbar-laravel"},
			gomod:    []string{"github.com/rollbar/rollbar-go"},
		},
		{
			recipe:   "salesforce",
			npm:      []string{"jsforce"},
			composer: []string{"omniphx/forrest"},
			gomod:    []string{"github.com/simpleforce/simpleforce"},
		},
		{
			recipe: "sanity",
			npm:    []string{"@sanity/client"},
		},
		{
			recipe: "secretsmanager",
			npm:    []string{"@aws-sdk/client-secrets-manager"},
			gomod:  []string{"github.com/aws/aws-sdk-go-v2/service/secretsmanager"},
		},
		{
			recipe:   "segment",
			npm:      []string{"@segment/analytics-node", "analytics-node"},
			composer: []string{"segmentio/analytics-php"},
			gomod:    []string{"github.com/segmentio/analytics-go"},
		},
		{
			recipe:   "sentry",
			npm:      []string{"@sentry/node", "@sentry/browser", "@sentry/nextjs"},
			composer: []string{"sentry/sentry", "sentry/sentry-laravel", "sentry/sentry-symfony"},
			gomod:    []string{"github.com/getsentry/sentry-go"},
		},
		{
			recipe:   "ses",
			npm:      []string{"@aws-sdk/client-ses", "@aws-sdk/client-sesv2"},
			composer: []string{"aws/aws-sdk-php"},
			gomod:    []string{"github.com/aws/aws-sdk-go-v2/service/sesv2"},
		},
		{
			recipe: "shippo",
			npm:    []string{"shippo"},
		},
		{
			recipe: "snyk",
			npm:    []string{"snyk"},
		},
		{
			recipe: "sqs",
			npm:    []string{"@aws-sdk/client-sqs", "sqs-consumer"},
			gomod:  []string{"github.com/aws/aws-sdk-go-v2/service/sqs"},
		},
		{
			recipe:   "square",
			npm:      []string{"square"},
			composer: []string{"square/square"},
			gomod:    []string{"github.com/square/square-go-sdk"},
		},
		{
			recipe: "statuspage",
			npm:    []string{"statuspage-api"},
		},
		{
			recipe: "stytch",
			npm:    []string{"stytch"},
			gomod:  []string{"github.com/stytchauth/stytch-go"},
		},
		{
			recipe:   "taxjar",
			npm:      []string{"taxjar"},
			composer: []string{"taxjar/taxjar-php"},
		},
		{
			recipe:   "telnyx",
			npm:      []string{"telnyx"},
			composer: []string{"telnyx/telnyx-php"},
		},
		{
			recipe: "trello",
			npm:    []string{"trello.js", "node-trello"},
		},
		{
			recipe: "typeform",
			npm:    []string{"@typeform/api-client"},
		},
		{
			recipe: "vercel",
			npm:    []string{"@vercel/client", "@vercel/sdk"},
		},
		{
			recipe:   "vonage",
			npm:      []string{"@vonage/server-sdk"},
			composer: []string{"vonage/client"},
		},
		{
			recipe: "webflow",
			npm:    []string{"webflow-api"},
		},
		{
			recipe:   "workos",
			npm:      []string{"@workos-inc/node"},
			composer: []string{"workos/workos-php"},
			gomod:    []string{"github.com/workos/workos-go"},
		},
		{
			recipe:   "zendesk",
			npm:      []string{"node-zendesk"},
			composer: []string{"zendesk/zendesk_api_client_php"},
			gomod:    []string{"github.com/nukosuke/go-zendesk"},
		},
		{
			recipe: "zoom",
			npm:    []string{"@zoom/appssdk"},
		},
		{
			recipe: "linear",
			npm:    []string{"@linear/sdk"},
		},
		{
			recipe: "supabase",
			npm:    []string{"@supabase/supabase-js"},
		},
		{
			recipe: "neon",
			npm:    []string{"@neondatabase/serverless", "@neondatabase/api-client"},
		},
		{
			recipe: "svix",
			npm:    []string{"svix"},
			gomod:  []string{"github.com/svix/svix-webhooks"},
		},
		{
			recipe: "typesense",
			npm:    []string{"typesense"},
			gomod:  []string{"github.com/typesense/typesense-go"},
		},
		{
			recipe:   "meilisearch",
			composer: []string{"meilisearch/meilisearch-php"},
			npm:      []string{"meilisearch"},
			gomod:    []string{"github.com/meilisearch/meilisearch-go"},
		},
		{
			recipe: "qdrant",
			npm:    []string{"@qdrant/js-client-rest"},
			gomod:  []string{"github.com/qdrant/go-client"},
		},
		{
			recipe: "weaviate",
			npm:    []string{"weaviate-client", "weaviate-ts-client"},
			gomod:  []string{"github.com/weaviate/weaviate-go-client"},
		},
		{
			recipe: "openrouter",
			npm:    []string{"@openrouter/ai-sdk-provider"},
		},
		{
			recipe: "replicate",
			npm:    []string{"replicate"},
			gomod:  []string{"github.com/replicate/replicate-go"},
		},
		{
			recipe: "cohere",
			npm:    []string{"cohere-ai"},
			gomod:  []string{"github.com/cohere-ai/cohere-go"},
		},
		{
			recipe: "langfuse",
			npm:    []string{"langfuse", "langfuse-langchain"},
		},
		{
			recipe: "monday",
			npm:    []string{"@mondaydotcomorg/api", "monday-sdk-js"},
		},
		{
			recipe: "hetzner",
			gomod:  []string{"github.com/hetznercloud/hcloud-go"},
		},
		{
			recipe: "heroku",
			npm:    []string{"heroku-client"},
			gomod:  []string{"github.com/heroku/heroku-go"},
		},
		{
			recipe:   "woocommerce",
			composer: []string{"automattic/woocommerce"},
			npm:      []string{"@woocommerce/woocommerce-rest-api"},
		},
		{
			recipe: "dockerhub",
			gomod:  []string{"github.com/docker/docker"},
		},
		{
			recipe: "npmregistry",
			npm:    []string{"pacote", "libnpmsearch"},
		},
		{
			recipe: "increase",
			npm:    []string{"increase"},
			gomod:  []string{"github.com/Increase/increase-go"},
		},
		{
			recipe: "orb",
			npm:    []string{"orb-billing"},
			gomod:  []string{"github.com/orbcorp/orb-go"},
		},
		{
			recipe: "knock",
			npm:    []string{"@knocklabs/node"},
			gomod:  []string{"github.com/knocklabs/knock-go"},
		},
		{
			recipe: "novu",
			npm:    []string{"@novu/node", "@novu/api"},
			gomod:  []string{"github.com/novuhq/novu-go"},
		},
		{
			recipe: "merge",
			npm:    []string{"@mergeapi/merge-node-client"},
		},
		{
			recipe: "finch",
			npm:    []string{"@tryfinch/finch-api"},
			gomod:  []string{"github.com/Finch-API/finch-api-go"},
		},
		{
			recipe: "gusto",
			npm:    []string{"@gusto/embedded-api"},
		},
		{
			recipe: "alpaca",
			npm:    []string{"@alpacahq/alpaca-trade-api"},
			gomod:  []string{"github.com/alpacahq/alpaca-trade-api-go"},
		},
		{
			recipe: "polygon",
			npm:    []string{"@polygon.io/client-js"},
			gomod:  []string{"github.com/polygon-io/client-go"},
		},
		{
			recipe: "apify",
			npm:    []string{"apify-client", "apify"},
		},
		{
			recipe: "uploadcare",
			npm:    []string{"@uploadcare/upload-client", "@uploadcare/rest-client"},
		},
		{
			recipe: "backblaze",
			npm:    []string{"backblaze-b2"},
			gomod:  []string{"github.com/Backblaze/blazer"},
		},
		{
			recipe: "fly",
			gomod:  []string{"github.com/superfly/fly-go"},
		},
		{
			recipe: "xendit",
			npm:    []string{"xendit-node"},
			gomod:  []string{"github.com/xendit/xendit-go"},
		},
		{
			recipe: "dwolla",
			npm:    []string{"dwolla-v2"},
		},
		{
			recipe: "razorpay",
			npm:    []string{"razorpay"},
		},
		{
			recipe: "snowflake",
			npm:    []string{"snowflake-sdk"},
			gomod:  []string{"github.com/snowflakedb/gosnowflake"},
		},
		{
			recipe: "truelayer",
			npm:    []string{"truelayer-client"},
		},
		{
			recipe: "midtrans",
			npm:    []string{"midtrans-client"},
		},
		{
			recipe: "shortcut",
			npm:    []string{"@shortcut/client"},
		},
		{
			recipe: "onfido",
			npm:    []string{"@onfido/api"},
		},
		{
			recipe: "courier",
			npm:    []string{"@trycourier/courier"},
		},
		{
			recipe: "kratos",
			npm:    []string{"@ory/kratos-client"},
			gomod:  []string{"github.com/ory/kratos-client-go"},
		},
		{
			recipe: "hydra",
			npm:    []string{"@ory/hydra-client"},
			gomod:  []string{"github.com/ory/hydra-client-go"},
		},
		{
			recipe: "braze",
			npm:    []string{"braze-api"},
		},
		{
			recipe: "bandwidth",
			npm:    []string{"@bandwidth/messaging"},
			gomod:  []string{"github.com/Bandwidth/go-sdk"},
		},
		{
			recipe: "finnhub",
			npm:    []string{"finnhub"},
			gomod:  []string{"github.com/Finnhub-Stock-API/finnhub-go"},
		},
		{
			recipe: "buildkite",
			npm:    []string{"buildkite"},
			gomod:  []string{"github.com/buildkite/go-buildkite"},
		},
		{
			recipe: "discourse",
			npm:    []string{"node-discourse-api"},
		},
		{
			recipe: "airwallex",
			npm:    []string{"@airwallex/node-sdk"},
		},
		{
			recipe: "documenso",
			npm:    []string{"@documenso/sdk-typescript"},
		},
		{
			// Nordigen is the same API under its former name, and plenty of
			// manifests still carry it.
			recipe: "gocardlessbank",
			npm:    []string{"nordigen-node", "gocardless-nodejs"},
		},
		{
			recipe: "googleaddress",
			npm:    []string{"@googlemaps/addressvalidation"},
		},
		{
			recipe: "moderntreasury",
			npm:    []string{"modern-treasury"},
			gomod:  []string{"github.com/Modern-Treasury/modern-treasury-go"},
		},
		{
			recipe: "opsgenie",
			npm:    []string{"opsgenie-sdk"},
		},
		{
			recipe: "wordpress",
			npm:    []string{"wpapi", "@wordpress/api-fetch"},
		},
		{
			// Segment's client rather than Help Scout's own, and still
			// unambiguously a Help Scout client, which is the question
			// detection asks.
			recipe: "helpscout",
			npm:    []string{"helpscout"},
		},
		{
			recipe: "calcom",
			npm:    []string{"@calcom/atoms"},
		},
	}
}

// unshipped names the Recipes a mapping points at deliberately even though no
// Recipe ships for them. Detection reports those by name so a developer is told
// which calls will still reach the real network, which is more useful than the
// "looks like an API client" heuristic that would otherwise catch them.
//
// It exists so that a typo cannot pass for an intention. Without it, a mapping
// naming "postmarkk" would quietly become a warning about a provider nobody has
// heard of, and nothing would fail.
var unshipped = map[string]string{
	"anthropic": "worth a Recipe on its own merits; until then the warning is the value",
}

// suspectedAPIClient matches dependency names that look like third-party API
// clients we have no Recipe for. These are reported rather than ignored so the
// developer knows exactly which calls will still hit the real network.
var suspectedAPIClientHints = []string{
	"-sdk", "sdk-", "-api", "api-client", "-client",
}

// recipeForComposer resolves a Composer package to a Recipe name.
func recipeForComposer(pkg string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.composer {
			if strings.EqualFold(candidate, pkg) {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// recipeForNPM resolves an npm package to a Recipe name.
func recipeForNPM(pkg string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.npm {
			if strings.EqualFold(candidate, pkg) {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// recipeForGoModule resolves a Go module path to a Recipe name. Module paths
// carry major-version suffixes (/v76), so this matches on prefix.
func recipeForGoModule(path string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.gomod {
			if path == candidate || strings.HasPrefix(path, candidate+"/") {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// looksLikeAPIClient is a heuristic used only to populate Unmatched. It never
// causes a fake to boot.
func looksLikeAPIClient(pkg string) bool {
	lower := strings.ToLower(pkg)

	for _, hint := range suspectedAPIClientHints {
		if strings.Contains(lower, hint) {
			return true
		}
	}

	return false
}
