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
