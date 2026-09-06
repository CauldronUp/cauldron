<div align="center">

# Cauldron

**Your code. One command. Every dependency.**

The open-source environment compiler. It boots your application *and the third-party APIs it depends on*, locally, from one command.

[![Licence](https://img.shields.io/badge/licence-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8.svg)](go.mod)

</div>

---

## The idea

Postgres runs in a container. Redis runs in a container. But the moment your application talks to something you don't own, local development turns into sandbox accounts, partner approvals, shared test tenants that somebody else keeps mutating, and `if (isLocal())` branches nobody remembers to delete.

Cauldron's claim is simple: **a third-party API is just another local service.**

```bash
git clone <repo>
cd project
cauldron up
```

Cauldron reads the manifests already in your repository, works out what the project needs, and boots it, including working emulations of the providers it talks to.

```
Detected a Laravel project.

Runtime
  +  php 8.5           (require.php)
  +  laravel 13.0      (laravel/framework)
  +  node 24.0.0       (engines.node)

Services
  +  redis             (predis/predis)
  +  horizon           (laravel/horizon)

Recipes
  +  stripe 17.0       (stripe/stripe-php)
  +  shopify 5.4       (shopify/shopify-api)

No recipe yet. These will still reach the real network:
  !  acme/weather-api-client  (composer.json)
```

That last section is deliberate. Falling back to the real network *silently* is how a test suite starts lying.

## Status

**Early, but the core works.** The emulator is real: you can point an SDK at it today.

| Area | State |
|---|---|
| Detection engine (Composer, npm, Go modules) | Working. 484 of 627 Recipes are reachable from a dependency |
| Recipe format and validator | Working |
| Recipe runtime (routing, state, auth, pagination) | Working. A credential check reports whether nothing was sent, something unusable was, or a well-formed credential was refused, and a Recipe may answer each differently -- most providers do not, and for those one message still covers all three |
| Webhooks (lifecycle events, signing, delivery) | Working. Payload envelopes are declarable per Recipe. Of the 103 Recipes that emit events, 45 declare one and 58 fall back to a default modelled on Stripe's and not equal to it -- Stripe's own Recipe now declares the fuller real thing -- so the fallback is a convention rather than a claim about any provider |
| Fault injection, clock control, request log | Working |
| Network conditions (latency, throttling, timeouts, resets) | Working. Toxiproxy's vocabulary, applied to the emulated providers |
| `serve`, `status`, `requests`, `seed`, `reset`, `fault`, `network`, `emit`, `clock` | Working |
| `doctor`, `logs`, `open` | Working |
| `cauldron up` / `down` (container orchestration) | Working for backing services |
| `snapshot` save/restore | Working |
| Conformance suites (`cauldron verify`) | Working. 8433 cases, 3309 of them checked against a live API |
| Spec drift (`cauldron drift`) | Working. 161 Recipes are checked against their provider's own OpenAPI document; 466 name none yet |
| Finding descriptions (`cauldron discover`) | Working. Proposes a description only where the document declares a path the Recipe already models. Found 35 across the collection, of which 29 declare every route |
| Description-backed routes (`serve --with-spec`) | Working. Adds routes from a Recipe's declared OpenAPI description for paths it does not model. Asana goes from 9 routes to 230; Adyen from 4 to 28. The Recipe always wins |
| Scoped multi-segment paths (`/repos/{owner}/{repo}/…`) | Working |
| Headless mode (`--headless`, `--host`) | Working. Providers only, one line of JSON, no containers |
| Application runtimes in containers | Not built. Run your app as you normally do |
| Recipes shipped | 627: Stripe, GitHub, GitLab, Bitbucket, Shopify, Shopify GraphQL, BigCommerce, Etsy, Magento, eBay, Amazon SP-API, EasyPost, AfterShip, ShipStation, ShipEngine, DHL, Recharge, Lemon Squeezy, Toast, Avalara, FedEx, Easyship, Clover, Lightspeed, Printful, Printify, Medusa, Shopware, commercetools, VTEX, Saleor, Allegro, Akeneo, Voucherify, Apideck, Royal Mail, Lago, Polar, Gumroad, Ecwid, Squarespace, Wix, Metronome, Twilio, Slack, HubSpot, SendGrid, Airtable, Notion, Zendesk, Postmark, Plaid, Clerk, Intercom, Discord, Square, Mailchimp, Cloudflare, Vercel, Xero, QuickBooks, PagerDuty, Asana, Algolia, Sentry, Box, Calendly, Datadog, Front, Typeform, Miro, Contentful, Sanity, Klaviyo, Webflow, Zoom, Pipedrive, Freshdesk, Mailgun, Okta, Shippo, Dropbox, Google Drive, Google Cloud Storage, Auth0, Recurly, Trello, Paddle, CircleCI, Snyk, Statuspage, Buildkite, ClickUp, Basecamp, Shortcut, Chargebee, Vonage, Rollbar, Docusign, Lob, Segment, Greenhouse, Adyen, Salesforce, Help Scout, DigitalOcean, WooCommerce, WordPress, AWS SQS, Google Pub/Sub, AWS DynamoDB, AWS Secrets Manager, AWS SES, Resend, Cloudinary, LaunchDarkly, Vimeo, Mux, Pusher, Ably, Ghost, WorkOS, TaxJar, Stytch, Keycloak, FusionAuth, AssemblyAI, Documenso, Dropbox Sign, Deepgram, Gusto, OneSignal, Telnyx, Bill.com, Ramp, Mercury, Brex, Deel, Google Calendar, Jira, Confluence, Jira Service Management, OpenAI, ElevenLabs, PayPal, RingCentral, Mollie, Gmail, Increase, Svix, Docker Hub, Onfido, Netlify, Firecrawl, Modern Treasury, npm registry, Bugsnag, Persona, Orb, PostHog, Column, marqeta, braze, bandwidth, alpaca, meilisearch, knock, merge, truelayer, razorpay, gocardlessbank, finch, replicate, cohere, hightouch, dwolla, airwallex, snowflake, incidentio, fivetran, novu, uploadcare, midtrans, courier, opsgenie, xendit, gorgias, fastspring, kustomer, calcom, twilioverify, qdrant, revenuecat, backblaze, kratos, hetzner, heroku, polygon, fly, apify, finnhub, tradier, typesense, weaviate, hydra, openrouter, Gemini, Google Address Validation, ShipHero, Linear, Monday, Attio, Neon, PlanetScale, MongoDB Atlas, Supabase, Discourse, Langfuse, YouTube, HashiCorp Vault, Unleash, Pinecone, Hookdeck, Grafana, ConfigCat, SonarCloud, GrowthBook, PyPI, crates.io, RubyGems, Repology, Homebrew, Packagist, Hex.pm, NuGet, Go proxy, Maven Central, OSV.dev, deps.dev, Open-Meteo, USGS, Frankfurter, Nominatim, Open Library, Wikipedia, Hacker News, PokeAPI, MusicBrainz, Open Food Facts, Crossref, Zippopotam, TVmaze, Sunrise-Sunset, NHTSA vPIC, Wikidata, GBIF, Nager.Date, Wayback Machine, National Weather Service, OpenTDB, Open Brewery DB, PoetryDB, Datamuse, Art Institute of Chicago, Postcodes.io, Met Museum, Rick and Morty, Deck of Cards, TheMealDB, World Bank, disease.sh, Dog CEO, Advice Slip, SWAPI, Chuck Norris, Open Notify, Agify, Bible API, CoinGecko, iTunes Search, EPSS, Chess.com, ip-api, Where the ISS at, TfL, openFDA, OpenAlex, MBTA, Carbon Intensity, UK Police, FHRS, BoC Valet, FBI Wanted, Zenodo, UniProt, Spaceflight News, ClinicalTrials.gov, DataCite, Deezer, iNaturalist, RCSB, Europe PMC, OpenAIRE, ROR, Nobel Prize, PDBe, Jisho, Stack Exchange, Kraken, Scryfall, Lichess, Gutendex, Wikimedia, SEC EDGAR, Alpha Vantage, NCBI Entrez, Semantic Scholar, Mercado Libre, US Census, NASA, Steam, NOAA, Guardian, arXiv, Healthchecks, PubNub, Paystack, Hugging Face, Vultr, Linode, Smartsheet, UptimeRobot, Render, Upstash, Mistral, xAI, Scaleway, Status.io, Tink, OpenAQ, Groq, Perplexity, TMDB, RAWG, Twitch, Spotify, Duffel, FlightAware, Honeycomb, New Relic, Wise, Flutterwave, Giphy, Unsplash, NewsAPI, GNews, Coda, GitBook, Split, Flagsmith, Brevo, Mailjet, Lithic, Moov, Descope, FRED, Twelve Data, Mapbox, OpenRouteService, Plivo, Sinch, Cronitor, Checkly, ScrapingBee, Browserless, Brave Search, Serper, Temporal, Inngest, Infobip, Customer.io, Pipedream, Make, Turso, SingleStore, Airbyte, Census, Rootly, FireHydrant, Axiom, Mezmo, Raygun, AppSignal, Ory, SuperTokens, Logto, Wasabi, Filebase, HERE, TomTom, what3words, Frontegg, Kinde, PropelAuth, Convex, Xata, Fauna, Airbrake, Papertrail, Sumo Logic, PandaDoc, SignNow, Anvil, Together AI, Fireworks, Lever, Ashby, Workable, Kit, Buttondown, Zapier, n8n, Windmill, Nuvei, Worldpay, Rapyd, Baserow, NocoDB, Retool, Bunny, ImageKit, Transloadit, Chroma, Milvus, Rippling, BambooHR, Workday, Diffbot, Bright Data, Zyte, Koyeb, Northflank, Porter, NetSuite, Sage Intacct, Expensify, ClickHouse, Timescale, InfluxDB, Sendcloud, ShipBob, Onfleet, Dynatrace, Scout APM, FullStory, Formstack, PDFMonkey, Api2Pdf, Akeyless, JFrog Artifactory, Sonatype Nexus, Swell, Commerce Layer, Snipcart, Payoneer, Treasury Prime, Melio, Tipalti, Anrok, Vertex, TriNet, Remote.com, Zenvia, Kaleyra, Textline, Constant Contact, Drip, Omnisend, Duo, Castle, Sardine, MinIO, Nile, Tigris, Elastic Cloud, Optimizely, VWO, Flexport, Bringg, Radar, TalkJS, CometChat, Whereby, Clio, Aha, Adapty, JustCall, CloudTalk, Missive, Papaya Global, Postscript, Appsmith, Budibase, Coveralls, SavvyCal, 100ms, Ayrshare, AWS Cognito, Google Maps Platform, Livepeer, UPS, USPS, Canada Post, Reddit, Bluesky, Mastodon, Coinbase, Binance, Alchemy, WeatherAPI, Tomorrow.io, Visual Crossing, IPinfo, MaxMind, Geoapify, Clearbit, Hunter.io, Apollo.io, Daily, LiveKit, Agora, Discogs, Last.fm, Genius, CockroachDB, Redis Cloud, SurrealDB, Checkout.com, Klarna, Mangopay, Amadeus, Kiwi.com, Hotelbeds, Sumsub, Veriff, Alloy, Anthropic, Stability AI, Fal.ai, MX, Belvo, Basiq, DeepL, Lokalise, Crowdin, TheSportsDB, football-data, balldontlie, Strava, Oura, Fitbit, Listen Notes, Podcast Index, Taddy, CourtListener, ProPublica, OpenStates, Electricity Maps, WattTime, EIA, ZeroBounce, Kickbox, NeverBounce, RentCast, ATTOM, Estated, Canvas, Moodle, Clever, OpenSky, AviationStack, AeroDataBox, Nutritionix, Edamam, Spoonacular, Google Books, ISBNdb, Hardcover, Toggl, Harvest, Clockify, Mixpanel, Amplitude, Heap, Exa, Tavily, SerpApi, Fastly, KeyCDN, Akamai, Namecheap, Porkbun, DNSimple, OpenCorporates, Companies House, PatentsView, Cloudflare Stream, api.video, Bunny Stream, Bitly, Dub, Short.io, SurveyMonkey, Jotform, Tally, 1Password, Bitwarden, Doppler, Better Stack, Instatus, Pingdom, Transitland, Navitia, 511, OMDb, Trakt, Watchmode |

## Try it

```bash
go build -o cauldron ./cmd/cauldron

cd /path/to/your/project
./cauldron up --fixture small-shop
```

```
Detected a Laravel project.

Services
  +  meilisearch  (laravel/scout)

Recipes
  +  stripe 17.0  (stripe/stripe-php)

Starting services (the first run pulls images, which can take a minute)
  +  meilisearch  http://127.0.0.1:7700 (key: cauldron)
  +  mailpit      smtp 127.0.0.1:1025 · inbox http://127.0.0.1:8025

Cauldron is listening on http://127.0.0.1:4600

  stripe       http://127.0.0.1:4600/stripe
```

Containers are labelled and scoped to the project directory, so two checkouts
side by side get separate environments. `cauldron down` removes them, and
`--keep-data` preserves the volumes.

### Headless: just the providers

If you already have a way to run your application, you probably do not want a
second one. Headless mode emulates the providers and nothing else: no
containers, no plan describing an environment Cauldron is not going to set up,
and one line of JSON instead of a banner addressed to a person.

```bash
cauldron serve --headless stripe woocommerce
```

```json
{"address":"http://127.0.0.1:4600","bind":"127.0.0.1","port":4600,"control":"http://127.0.0.1:4600/_cauldron/status","mounted":[{"recipe":"stripe","url":"http://127.0.0.1:4600/stripe"},{"recipe":"woocommerce","url":"http://127.0.0.1:4600/woocommerce"}],"missing":[],"unseeded":[]}
```

One line, printed before anything else, so a script can read it and get on with
starting the application itself. `missing` is in there for the same reason it is
in the banner: a provider Cauldron cannot emulate still reaches the real
network, and that should not be something you find out from a charge.

`cauldron up --headless` does the same and skips the containers, for a project
whose services are already running.

**[docs/headless.md](docs/headless.md)** covers dropping this into an
environment that already exists: Docker Compose in both directions, CI, waiting
for readiness without polling, and pointing each SDK at it. Every example there
is exercised by the `headless` CI job on each push.

Cauldron binds loopback by default. An application in its own container cannot
reach the host's loopback, so `--host` changes that:

```bash
cauldron serve --headless --host 0.0.0.0 stripe
```

The reported `address` stays dialable from this machine even on a wildcard bind,
and the real bind travels beside it as `bind` and `port`, because `http://[::]:4600`
is not a URL anybody can use. Binding past loopback also puts the control plane
on the network, where it can seed, reset and fault the providers, so Cauldron
says so on stderr when you do it.

Then point an SDK at `http://127.0.0.1:4600/stripe` and use it as you would the real thing:

```bash
# Create a customer, form encoded exactly as Stripe's own SDKs send it
curl -X POST http://127.0.0.1:4600/stripe/v1/customers \
  -H "Authorization: Bearer sk_test_cauldron" \
  -d "email=ada@example.com&name=Ada+Lovelace&metadata[plan]=pro"
```

```json
{
  "created": 1767225600,
  "currency": "usd",
  "email": "ada@example.com",
  "id": "cus_rfBd56ti2SMtYv",
  "metadata": { "plan": "pro" },
  "name": "Ada Lovelace"
}
```

That identifier is not random. The same seed and fixture produce it on every machine, every run.

### The things a sandbox cannot do

```bash
# Rate limit the next request, with the provider's real Retry-After header
curl -X POST http://127.0.0.1:4600/_cauldron/stripe/fault \
  -H 'Content-Type: application/json' \
  -d '{"error":"rate_limit","count":1}'

# Age everything by a month, so a subscription falls into dunning
curl -X POST http://127.0.0.1:4600/_cauldron/clock/advance \
  -H 'Content-Type: application/json' \
  -d '{"duration":"30d"}'

# Fire a signed webhook at your application
curl -X POST http://127.0.0.1:4600/_cauldron/stripe/emit \
  -H 'Content-Type: application/json' \
  -d '{"event":"payment_intent.payment_failed","data":{"amount":2500}}'

# See what your code actually sent
curl http://127.0.0.1:4600/_cauldron/stripe/requests
```

Webhook deliveries are recorded even when nothing is subscribed, so "this event would have fired" is assertable without standing up a listener.

Some webhooks arrive without being asked for. A create, update or delete fires
the event its route names, so posting a customer to Stripe delivers
`customer.created` with no further instruction, and creating a ticket in
Freshdesk delivers `ticket_create` -- which is what Freshdesk calls it. A route
says so with `emits:`, because the name is the provider's rather than a
convention: Bitbucket sends `pullrequest:created`, Box `FILE.DELETED`, Trello
`createCard`, Xero `INVOICE.CREATE`. A Recipe that names an event its own
`webhooks.events` does not list is refused, since the failure mode of that typo
is a webhook that never arrives, which is indistinguishable from a provider
that does not send one.

Most declared events are not named for a change at all. `crawl.completed`,
`video.asset.ready`, `user.session.start` and `payment.failed` are things that
happen to a record later, or to something that is not a record. Those are what
`emit` is for, and a Recipe listing them is describing its provider rather than
leaving something unwired. Of 504 events declared across these Recipes, 282 are
in that category.

One of the 264 fires anyway, which is the useful exception: Firecrawl's
`crawl.started` is what creating a crawl does, so that route names it. Whether
an event follows from a change is a question about the provider and not about
the shape of its name.

Or drive it from the CLI:

```bash
cauldron status                                   # what is running
cauldron seed stripe --fixture small-shop         # load seed data
cauldron fault stripe --error rate_limit --count 2 # break the next two calls
cauldron network stripe --latency 800ms           # make the link slow
cauldron clock advance 30d                        # age everything a month
cauldron emit stripe payment_intent.payment_failed # fire a webhook
cauldron requests stripe                          # what your code actually sent
cauldron reset                                    # back to a clean sandbox
```

```
#  METHOD  PATH           STATUS  NOTE
1  GET     /v1/customers  429     fault: rate_limit
2  GET     /v1/customers  429     fault: rate_limit
3  GET     /v1/customers  200     list
```

Point the CLI at a different server with `--url` or `$CAULDRON_URL`.

### Breaking the network, not just the provider

`cauldron fault` covers what a provider does *deliberately* when something is
wrong: rate limits, validation errors, 5xx. `cauldron network` covers what the
network does to you regardless of what the provider intended.

```bash
cauldron network stripe --latency 800ms --jitter 200ms   # a slow link
cauldron network stripe --bandwidth 50                   # a thin one
cauldron network stripe --timeout 30s                    # accept, answer nothing, hang up
cauldron network stripe --reset --probability 0.1        # one call in ten dies mid-flight
cauldron network stripe --clear
```

The difference is the whole point. A rate limit arrives as a well-formed 429
your client library already understands and probably already retries. A
connection that hangs for ninety seconds and then dies arrives as nothing at
all: no status, no body, no error type your code has a branch for. That is the
one that pages someone, and it is the one no provider's sandbox will ever hand
you.

The vocabulary is [Toxiproxy's](https://github.com/Shopify/toxiproxy) on
purpose — latency, jitter, bandwidth, timeout, reset, slice, limit. If you have
used Toxiproxy to break the link to your database, you already know what these
flags do, and learning a second set of words for the same idea would buy
nothing. Cauldron applies them to the third-party APIs Toxiproxy cannot help
with, because those do not exist locally for it to sit in front of.

The two compose. Point Toxiproxy at Postgres and Cauldron at Stripe, and every
dependency your application has can be made to misbehave from a single test. In
PHP, [`mpge/toxiproxy-php`](https://github.com/mpge/toxiproxy-php) is the other
half: it manages the Toxiproxy server for you the way Cauldron manages the
provider fakes.

Conditions are applied before the provider's own behaviour, are reproducible
from the sandbox seed even below a probability of 1, and show up in `status` and
`requests` so an odd timing has a visible cause rather than looking like a
flake. `docs/network.md` has the full reference, including the two places an
HTTP emulator honestly cannot be a TCP proxy.

### When something is wrong

```bash
cauldron doctor
```

```
ok    project                Laravel, 2 service(s), 1 recipe(s)
ok    recipes                4 bundled, all parsing
ok    docker                 running
warn  ports                  6379 wanted by redis, held by other-project-redis-1
                             stop the other container, or start this project without that service
ok    sandbox                serving stripe
```

It checks the things that are obvious once you know and cost twenty minutes
when you do not: Docker missing, a port another checkout is already holding, a
Recipe that stopped parsing, dependencies with no Recipe that will still reach
the real network.

```bash
cauldron logs redis          # recent output from one service
cauldron open mailpit        # the inbox, in a browser
cauldron open --print        # just print the address
```

### How faithful is it, really

Every fake eventually has to answer "how do you know it behaves like the real
thing?", and for most of them the honest answer is "you don't". Cauldron's
answer is a conformance suite: each Recipe carries checkable claims about the
provider, every claim cites where it came from, and the report separates what
was observed against the real API from what was only read in the documentation.

```bash
cauldron verify            # every bundled Recipe
cauldron verify stripe -v  # one Recipe, listing each case
```

```
stripe 0.1.0
  18 of 18 cases passed
  4 checked against the real API, 14 from documentation only
  last checked against the real API on 2026-09-05
```

That second line is the honest one. Of every Recipe: 8433 cases, 3309 run against
a live account and 5124 not. Documentation-derived cases are worth having,
and they are not the same as watching the provider do it. Adding a `verified:`
date to a case is a claim that someone did.

All 3309 are the cases whose provider can be asked without a key. Every other
provider needs an account, and a date nobody can reproduce is worth less than an
empty field that says so.

Every provider's own findings live beside its Recipe, in
`recipes/<name>/README.md` -- what was probed, what it answered, and what was
deliberately not modelled. They were here, all three thousand lines of them, and
a file nobody reads to the end is a poor place for the one note a person wants
when they are about to use a particular fake.

Cases run in process and need no credentials, so CI runs them on every push and
a Recipe edit that drifts from the provider fails there rather than in an
application months later.

Writing a Recipe is the format's stress test, and every provider has found
something -- a failure that arrives with HTTP 200, an identifier in a query
parameter, money as an object in minor units, an error body that is one bare
string. Those findings are in each Recipe's own README, beside the Recipe that
found them, because the one a reader wants is the one for the provider they are
about to fake.

A pattern runs through them. The state that breaks an integration is almost
never a failure: it is a third state nobody branched on -- a request closed
rather than completed, a job waiting on approval, a subscription with a
cancellation scheduled and not yet applied. Every fixture carries one, and
finding it is most of the work.

That work is now research rather than engineering. The format has needed no new
mechanism for most of the Recipes added recently; what a new provider costs is
reading it closely enough to know which state nobody branches on.

The suite has started finding faults in itself as well as in the Recipes. A
provider reporting success with the string `"0"` is what made the conformance
checker compare a scalar's kind as well as its value, and that change
immediately turned up a real inconsistency in the runtime's own error rendering.
Several Recipes record a bug they found in the emulator rather than in their
provider, which is the more useful half of writing one.

Not every provider fits, and one is left out deliberately rather than faked.
Temporal Cloud is gRPC rather than HTTP: this format describes HTTP surfaces,
and nothing built on it would be a Temporal client. No Recipe is better than a
misleading one.

That list used to be longer, and two of its entries were wrong. One named a
provider the format had since learned to describe; the other named one that had
been assessed incorrectly in the first place. Both sat unchallenged because
nothing rereads a list of things that cannot be done, and a reason that has
quietly expired reads exactly like one that still holds. Each of those Recipes
now records what it was excluded for and why that stopped being true.

### Snapshots

A bug that only reproduces after eleven API calls is a bug nobody else can
reproduce. Capture the sandbox instead of describing it:

```bash
cauldron snapshot save before-refund-bug
cauldron snapshot load before-refund-bug
cauldron snapshot list
```

A snapshot holds every mounted recipe's state, the sandbox clock, and the
position of the identifier generator, so restoring it does not rewind the next
identifier into one you already used. It lands in `.cauldron/snapshots/` inside
the project, which is the point: commit it, or attach it to the issue, and the
next person runs the same failure rather than reading about it.

Restoring refuses a snapshot that names a recipe you are not running, rather
than half-applying it. A recipe whose version has moved on since the capture
loads with a warning.

### Fixtures per recipe

A project with two providers rarely has one fixture name that fits both.
`--fixture` takes either a single name or `recipe=fixture` pairs:

```bash
cauldron up --fixture small-shop                        # applied where it fits
cauldron up --fixture stripe=small-shop,github=small-repo
```

A bare name is a convenience, so recipes that do not ship it are listed and left
at their default rather than failing the run. Name a recipe explicitly and a
missing fixture is an error, because you asked for something specific.

### Finding the descriptions nobody recorded

Twelve Recipes named a provider's description, and the README reported that as
though 301 providers publish nothing. It was never a measurement. It was how
many somebody had happened to find by hand.

`cauldron discover` looks for the rest, and the hard part is not looking. It is
refusing to believe what you find.

```
$ cauldron discover
  clerk                    10 of 10 route(s) declared, in 185 path(s), OpenAPI 3.0.3
                             spec: https://clerk.com/openapi.json
  royalmail                4 of 4 route(s) declared, in 11 path(s), Swagger 2.0
                             spec: https://api.parcel.royalmail.com/swagger/v1/swagger.json
  monday                   1 of 1 route(s) declared, in 2 path(s), OpenAPI 3.1.0
                             spec: https://monday.com/openapi.json
                             (weak: one route in a 2-path document -- check this
                              is the API and not a docs site)

35 proposed, 266 searched and nothing found, 0 with no verified host to search from, 12 already recorded.
```

**A URL resolving says nothing, and a 200 says nothing.** Documentation
platforms serve a generic `openapi.json` at the vendor's marketing domain:
`ramp.com` and `circleci.com` both answer one, with two and five paths,
describing the docs site rather than the API. Matching by domain instead scored
98 hits across this collection, a dozen of which were Recipes whose
documentation lives on github.com being paired with GitHub's own description.
Matching on a public directory's declared hosts paired Bugsnag with ClickSend.

So a proposal requires the document to declare a path the Recipe already
models -- the same intersection the fingerprint takes. Because the test is on
the document rather than the address, the search can afford to be wide: hosts
come from conformance sources and then widen to the siblings a description is
commonly published on. Attio is verified at `api.attio.com` and publishes at
`attio.com`; neither it nor Monday is reachable without that.

**Nothing is written.** The output is lines to paste, the same bargain `drift
--record` makes and for a sharper reason: a wrong spec URL does not fail
loudly. It reports drift against a document that was never this provider's,
every morning, until somebody works out why.

Of the 29 proposals declaring every route, `cauldron check` found nothing to
disagree with in 18, and those are recorded. Five disagree and are deliberately
not recorded -- Resend's description declares `SendEmailResponse` as `{id}`
alone where the Recipe answers `POST /emails` with a whole email, which is a
question about the Recipe that wants a live call to settle. Recording a
description a Recipe contradicts files the argument as a fingerprint and loses
it.

**The number will not reach 627.** Most providers publish nothing, and this
search reaches only descriptions served at a conventional URL -- a provider
publishing in a GitHub repository, which is how eleven of the first twelve here
were found, still needs a person to say where it is.

## When the provider moves

A Recipe is written once and the provider keeps going. Its conformance suite
cannot notice: that suite asserts what the Recipe says, so renaming a field
upstream leaves every case green and every case wrong. The suite checks the
Recipe's internal consistency, and internal consistency is exactly what
survives the provider changing its mind.

`cauldron drift` fetches the provider's own OpenAPI document and compares it
with a fingerprint recorded in the Recipe.

```
$ cauldron drift
  adyen                    unchanged since 2026-08-30
  ...

89 unchanged, 0 moved, 0 unreachable, 0 in a format this cannot read, 0 unrecorded, 421 with no description to check.
A Recipe with no description this can read is not verified by this. It is unexamined.
```

**The fingerprint is not a checksum of the file, and that is the whole design.**
Providers republish these documents constantly -- a reworded summary, a new
example, an endpoint in a product the Recipe has never heard of -- and a
checksum calls every one of those drift. A scan that goes red on every publish
gets switched off, and then the change that mattered arrives unannounced. The
noisy check and no check at all fail the same way; the noisy one costs more on
the way there.

So the fingerprint covers the intersection and nothing else: the paths and
methods the Recipe declares routes for, the response codes those operations
answer with, and the types of the fields the Recipe itself names. A path
appearing that the Recipe says nothing about does not move it. A field the
Recipe claims changing type does.

Six states, and only one of them fails a build:

| State | Means |
| --- | --- |
| `unchanged` | The description still says what the Recipe claims |
| `moved` | It does not. **The only state that exits non-zero** |
| `unreachable` | The host did not answer, or answered with something that is not a description. Not drift: a docs host returning 503 has said nothing about whether the provider changed anything |
| `unsupported` | The provider publishes a description in a format this cannot read -- RAML, AsyncAPI, WSDL. Separate from `unreachable` because it is permanent: a host answering 503 is a Tuesday, a format with no reader stays true tomorrow. Swagger 2.0 was here until it was rewritten on the way in |
| `unrecorded` | A description is named and has never been fingerprinted, so nothing was compared |
| `undeclared` | No description is recorded for this Recipe. 501 of the 627, and falling -- see `cauldron discover` |

It runs on a schedule rather than on every pull request, in
`.github/workflows/drift.yml`, because a pull request must not go red because
somebody else's docs host was restarting.

**What it cannot do is worth saying plainly.** A description will tell you a
payment has a `status` of type string on the day the provider starts answering
`"approved"` for payments nobody was paid for, and the fingerprint will not
move, because nothing the document says has changed. That is what a Recipe is
for. Drift catches the mechanical half -- a field renamed, a path moved, a
status dropped -- which is precisely the half a Recipe's own suite is
structurally unable to catch. A Recipe whose fingerprint has not moved is
un-contradicted by the description, on the parts it claims. That is a smaller
sentence than "unchanged", and it is the true one.

### Borrowing breadth from a description

A Recipe models what somebody sat down and checked. That is a small number of
endpoints, and it always will be:

| Recipe | Paths the description declares | Routes the Recipe models |
| --- | --- | --- |
| Box | 186 | 7 |
| OpenAI | 182 | 5 |
| Asana | 175 | 11 |
| Twilio | 121 | 4 |
| Intercom | 79 | 9 |
| ShipEngine | 70 | 7 |
| Adyen | 26 | 4 |

Hand-writing the remainder is not going to happen. `serve --with-spec` fetches
the description a Recipe already names, drafts routes from it, and mounts them
for the paths the Recipe is silent about:

```
$ cauldron serve --with-spec asana
...
221 route(s) added from provider descriptions, 9 already modelled and left alone.
  asana                    +221 derived

A derived route is what the provider says it does.
A Recipe route is what it was seen doing.
```

**The Recipe always wins.** A description may add an endpoint the Recipe does
not have; it may never change one the Recipe does, nor an error the Recipe
declares. Where a description's resource name collides with a Recipe's, the
route is dropped rather than guessed at, and the reason is printed.

**A derived route is a smaller claim, and the output says so every time.** This
is not modesty. Asana's own description omits `completed` and `due_on` from the
task listing -- they are `opt_fields`, real and undeclared -- so a mock built
from that description serves tasks without them, and the Recipe that pins them
is right exactly where the description is wrong. Kraken's description says
`error` is an array of string; nothing in OpenAPI can say that the array is
empty on success and that an empty array is *true*. Deezer's declares
`200: Artist`; it cannot declare that every failure is also a 200.

Three rules keep the addition honest rather than merely large:

- **It cannot break a Recipe.** The merged Recipe is validated, and if it does
  not validate the written Recipe is served unchanged with the reason printed.
  Adyen found that one: its Recipe wraps listings without a key, so every
  resource owes a collection name and the drafted ones had none.
- **Derived routes keep the description's base path.** Adyen declares
  `/sessions` beside a server of `.../v71`, so without this every derived route
  mounts at a path no client would call -- reported as added, and missing.
  Serving a route at the wrong path is worse than not serving it.
- **A derived resource follows the Recipe's own envelope where the Recipe is
  unanimous.** Asana wraps every listing in `data` and all three written
  resources say so, so a derived resource says `data` too rather than naming
  itself after its path. Where the written resources disagree there is no
  convention to follow and the path stands.

It is off by default, and has to be: reaching for a description is a network
call in a tool whose whole value is that it works without one. A Recipe whose
provider publishes nothing -- 501 of the 627 here -- is served exactly as it
always was, and told so plainly.

### Which version, and which one it replaced

`upstream.api` has always held the provider's version label -- 95 Recipes say
`v1` and 50 say `v2`. What it never held is whose, so a provider's v1 and its
v2 were unrelated strings, and the knowledge that one replaced the other lived
in prose beside a detection entry where nothing could check it.

```yaml
upstream:
  api: v3
  provider: mbta
  supersedes:
    - version: v2
      host: realtime.mbta.com
      note: the retired realtime interface, whose clients must not be offered this Recipe
```

The host is the part that does work: a client is told to be a v2 client by what
it talks to, and that is the only thing detection can tell versions apart by.
Three Recipes here already knew such a fact and could only write it down as a
comment -- the MBTA's `realtime.mbta.com`, ClinicalTrials.gov's
`clinicaltrials.gov/ct2`, and DataCite's `mds.datacite.org`.

## Recipes

A **Recipe** defines how Cauldron emulates an external dependency. It is not a mock, because anyone can return `200 OK`. A Recipe models behaviour, and carries the tests that prove it still matches the real API:

| Part | Covers |
|---|---|
| Auth | OAuth flows, API keys, HMAC signing, token refresh and expiry |
| Resources | Schemas, ID formats, relationships between objects |
| Routes | Paths, pagination and cursors, filtering, sorting |
| Behaviour | State transitions and side effects, so an order decrements inventory |
| Webhooks | Event catalogue, payload envelope, signing, delivery order, retries, duplicates |
| Errors | The provider's real error taxonomy and rate-limit shape |
| Fixtures | Named seed data: `empty`, `small-shop`, `enterprise` |
| Conformance | The suite that proves the Recipe still matches the real API |

Each one lives in `recipes/<name>/` and carries three things a reader wants
separately: `recipe.yaml` is the Recipe itself, with the full notes in its
header -- what was probed, what it answered, what was deliberately not modelled
and why. `README.md` is the short version: what this provider turned out to do
that a caller would not expect, and how many of its claims were checked against
the live API rather than read in documentation. The conformance cases sit inside
the Recipe, beside the claims they check.

If you are about to fake one provider, its README is the page to read.

## Design commitments

These are the decisions the project intends to be held to.

**Detection never guesses.** Package-to-Recipe mapping is an explicit table. A wrong guess is worse than no guess, because booting the wrong fake sends someone chasing a bug that doesn't exist. Anything unrecognised is reported, never silently faked.

A project can also say so itself. `cauldron add mercury` writes a `cauldron.yaml` listing the providers it talks to, which is how a project reaches a Recipe no dependency maps to, or one it talks to over raw HTTP with no library at all. The first `add` copies whatever detection already found into the file, so starting one loses nothing, and from then on the file is the answer rather than the guess.

The table reaches 484 of the 627 Recipes. The other 143 ship and can be named
directly. Every one of them was looked for: some publish no client at all,
and some publish a package whose name matches the provider and whose contents
do not -- `basecamp` on npm is a set of Astro components for somebody's scout
group, `mercury` is a frontend framework, `persona-api` is a chatbot. Mapping
one of those would boot the wrong fake for a project that never talks to that
provider, which is worse than the gap.

It also holds one entry with no Recipe behind it. A dependency Cauldron recognises and cannot emulate is reported by name, which tells you more than "this looks like an API client" does. `go test ./internal/detect/` prints the current figure and names what is missing.

**Determinism is enforced at the boundary, not requested politely.** Time comes from a movable sandbox clock and identifiers from a seeded generator. A Recipe has no access to the wall clock or to unseeded randomness. The same fixture and seed produce the same IDs on every machine, and a reset sandbox is indistinguishable from a fresh one.

**A name the Recipe chooses must be asserted where the value exists.** A
Recipe declaring `cursor_field: next_page_token` and only ever asserting that
the cursor is *absent* on a last page has claimed nothing: the absence holds
whatever the field is called. Twenty-one of these shipped before anything
checked, across twenty Recipes, and two of them turned out to have no
successful list case at all. The validator now refuses a declared paging name
that no case pins down.

**A conformance case must claim something the emulator decides.** A case that
posts `{"name": "Ada"}` and asserts `{"name": "Ada"}` is testing that the
request survived the round trip, which every fake does by construction. It
reads as evidence and is not, and it passes no matter how wrong the Recipe is.
Four were written before anything checked, each found only by breaking the
Recipe on purpose and noticing the case still passed; one could not be salvaged
and was replaced. The validator now refuses a case whose every claim repeats a
value it sent, and one non-echoed claim clears it.

**Nothing is held back to make a paid tier viable.** Everything that runs on your machine is Apache-2.0 and stays that way.

**More providers are queued.** [docs/backlog.md](docs/backlog.md) lists them,
roughly in order, along with the one that was assessed and deliberately left
out. A Recipe has to earn its place: the test is whether the provider has
behaviour worth reproducing, not whether adding it raises the count.

**Fidelity will be measured, not claimed.** The intent is that every Recipe carries a conformance suite run against the real provider, and that divergence is published rather than hidden. This is not built yet. The Recipes here are hand-modelled and unverified, and the status table says so.

## Contributing

The most valuable contribution is a Recipe for an API you actually use, especially one whose sandbox is painful to get. See [CONTRIBUTING.md](CONTRIBUTING.md), and [docs/format.md](docs/format.md) for every key a Recipe may carry — a test fails if that reference drifts from the format in either direction.

## Licence

[Apache-2.0](LICENSE).

A [Brilliance Digital](https://github.com/BrillianceDigital) project.
