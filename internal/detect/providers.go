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
			// No Recipe ships for OpenAI yet, and the mapping is here on
			// purpose. A dependency Cauldron recognises and cannot emulate is
			// reported by name, so the developer is told those calls will
			// still reach the real network. Without the mapping it falls
			// through to the "looks like an API client" heuristic and the
			// warning is vaguer than it needs to be.
			recipe:   "openai",
			composer: []string{"openai-php/client", "openai-php/laravel"},
			npm:      []string{"openai"},
			gomod:    []string{"github.com/sashabaranov/go-openai"},
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
	"openai": "worth a Recipe on its own merits; until then the warning is the value",
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
