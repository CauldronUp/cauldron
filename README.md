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
| Detection engine (Composer, npm, Go modules) | Working. 271 of 298 Recipes are reachable from a dependency |
| Recipe format and validator | Working |
| Recipe runtime (routing, state, auth, pagination) | Working |
| Webhooks (lifecycle events, signing, delivery) | Working. Payload envelopes are declarable per Recipe. Of the 102 Recipes that emit events, 44 declare one and 58 fall back to a default modelled on Stripe's and not equal to it -- Stripe's own Recipe now declares the fuller real thing -- so the fallback is a convention rather than a claim about any provider |
| Fault injection, clock control, request log | Working |
| Network conditions (latency, throttling, timeouts, resets) | Working. Toxiproxy's vocabulary, applied to the emulated providers |
| `serve`, `status`, `requests`, `seed`, `reset`, `fault`, `network`, `emit`, `clock` | Working |
| `doctor`, `logs`, `open` | Working |
| `cauldron up` / `down` (container orchestration) | Working for backing services |
| `snapshot` save/restore | Working |
| Conformance suites (`cauldron verify`) | Working. 2985 cases, 670 of them checked against a live API |
| Spec drift (`cauldron drift`) | Working. 11 Recipes are checked against their provider's own OpenAPI document, 1 provider publishes one this cannot read, and 286 publish none |
| Scoped multi-segment paths (`/repos/{owner}/{repo}/…`) | Working |
| Headless mode (`--headless`, `--host`) | Working. Providers only, one line of JSON, no containers |
| Application runtimes in containers | Not built. Run your app as you normally do |
| Recipes shipped | 298: Stripe, GitHub, GitLab, Bitbucket, Shopify, Shopify GraphQL, BigCommerce, Etsy, Magento, eBay, Amazon SP-API, EasyPost, AfterShip, ShipStation, ShipEngine, DHL, Recharge, Lemon Squeezy, Toast, Avalara, FedEx, Easyship, Clover, Lightspeed, Printful, Printify, Medusa, Shopware, commercetools, VTEX, Saleor, Allegro, Akeneo, Voucherify, Apideck, Royal Mail, Lago, Polar, Gumroad, Ecwid, Squarespace, Wix, Metronome, Twilio, Slack, HubSpot, SendGrid, Airtable, Notion, Zendesk, Postmark, Plaid, Clerk, Intercom, Discord, Square, Mailchimp, Cloudflare, Vercel, Xero, QuickBooks, PagerDuty, Asana, Algolia, Sentry, Box, Calendly, Datadog, Front, Typeform, Miro, Contentful, Sanity, Klaviyo, Webflow, Zoom, Pipedrive, Freshdesk, Mailgun, Okta, Shippo, Dropbox, Google Drive, Google Cloud Storage, Auth0, Recurly, Trello, Paddle, CircleCI, Snyk, Statuspage, Buildkite, ClickUp, Basecamp, Shortcut, Chargebee, Vonage, Rollbar, Docusign, Lob, Segment, Greenhouse, Adyen, Salesforce, Help Scout, DigitalOcean, WooCommerce, WordPress, AWS SQS, Google Pub/Sub, AWS DynamoDB, AWS Secrets Manager, AWS SES, Resend, Cloudinary, LaunchDarkly, Vimeo, Mux, Pusher, Ably, Ghost, WorkOS, TaxJar, Stytch, Keycloak, FusionAuth, AssemblyAI, Documenso, Dropbox Sign, Deepgram, Gusto, OneSignal, Telnyx, Bill.com, Ramp, Mercury, Brex, Deel, Google Calendar, Jira, Confluence, Jira Service Management, OpenAI, ElevenLabs, PayPal, RingCentral, Mollie, Gmail, Increase, Svix, Docker Hub, Onfido, Netlify, Firecrawl, Modern Treasury, npm registry, Bugsnag, Persona, Orb, PostHog, Column, marqeta, braze, bandwidth, alpaca, meilisearch, knock, merge, truelayer, razorpay, gocardlessbank, finch, replicate, cohere, hightouch, dwolla, airwallex, snowflake, incidentio, fivetran, novu, uploadcare, midtrans, courier, opsgenie, xendit, gorgias, fastspring, kustomer, calcom, twilioverify, qdrant, revenuecat, backblaze, kratos, hetzner, heroku, polygon, fly, apify, finnhub, tradier, typesense, weaviate, hydra, openrouter, Gemini, Google Address Validation, ShipHero, Linear, Monday, Attio, Neon, PlanetScale, MongoDB Atlas, Supabase, Discourse, Langfuse, YouTube, HashiCorp Vault, Unleash, Pinecone, Hookdeck, Grafana, ConfigCat, SonarCloud, GrowthBook, PyPI, crates.io, RubyGems, Repology, Homebrew, Packagist, Hex.pm, NuGet, Go proxy, Maven Central, OSV.dev, deps.dev, Open-Meteo, USGS, Frankfurter, Nominatim, Open Library, Wikipedia, Hacker News, PokeAPI, MusicBrainz, Open Food Facts, Crossref, Zippopotam, TVmaze, Sunrise-Sunset, NHTSA vPIC, Wikidata, GBIF, Nager.Date, Wayback Machine, National Weather Service, OpenTDB, Open Brewery DB, PoetryDB, Datamuse, Art Institute of Chicago, Postcodes.io, Met Museum, Rick and Morty, Deck of Cards, TheMealDB, World Bank, disease.sh, Dog CEO, Advice Slip, SWAPI, Chuck Norris, Open Notify, Agify, Bible API, CoinGecko, iTunes Search, EPSS, Chess.com, ip-api, Where the ISS at, TfL, openFDA, OpenAlex, MBTA, Carbon Intensity, UK Police, FHRS, BoC Valet, FBI Wanted, Zenodo, UniProt, Spaceflight News, ClinicalTrials.gov, DataCite |

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
leaving something unwired. Of 493 events declared across these Recipes, 271 are
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
  10 of 10 cases passed
  10 from documentation only, none checked against the real API
```

That second line is the honest one. Of every Recipe: 2985 cases, 670 run against
a live account and 2315 not. Documentation-derived cases are worth having,
and they are not the same as watching the provider do it. Adding a `verified:`
date to a case is a claim that someone did.

The six hundred and seventy are the cases whose provider can be asked without a key. Six are
OpenRouter's model-catalogue cases, whose numbers were read from the provider
rather than inferred; its completion cases carry no date, because calling that
endpoint costs money. Eight are the npm registry's, where what was checked is
the shape each case turns on rather than the values: that a per-version
document carries no time, no versions and no dist-tags, that `deprecated` is a
string and is absent rather than false when it does not apply, that
`integrity` sits inside `dist`, and that a missing package answers 404 with
`{"error":"Not found"}`. The Recipe records which real packages those were
observed on, so the dates can be audited rather than taken on trust.

Two of those seven were settled by the packages nobody wants to be famous.
`time` carries an entry for a version that `versions` does not, left behind by
an unpublish -- and the versions with one are `event-stream` 3.3.6, `rc`
1.2.9, 1.3.9 and 2.3.9, `flatmap-stream` 11.1.1 and `eslint-scope` 3.7.2,
every one pulled after a compromise. `eslint-scope` has thirty versions and
thirty-three `time` entries: thirty, `created`, `modified`, and the tombstone.
Fetching that version answers 404. So a release history rebuilt from `time` is
not merely wrong, it is disproportionately wrong about the releases somebody
removed on purpose.

The eighth is two shapes for one status. A package that does not exist answers
`{"error":"Not found"}`; a version that does not exist on a package that does
answers the bare JSON string `"version not found: 99.99.99"`. Code reading
`body.error` off the second finds `undefined` rather than failing, so a client
reporting `body.error` as the reason reports "undefined" -- quieter and worse
than throwing. Plain text would be the wrong way to model it, because a client
calling `.json()` on text throws and on this it succeeds and hands back a
string, so the case asserts the raw bytes: the quotes are the distinction.

Seven are Docker Hub's, which answers for a public repository without a token:
a repository nobody can see is a 404 carrying `{"message": "object not
found"}` and no code, `count` is of everything rather than of the page, and
one image wears several tags -- on one page of `library/nginx`, seventeen
digests were shared.

Three more of its cases were settled the same day. `latest` is an ordinary tag
and is routinely the older one: `library/mongo` last updated `latest` on
2026-07-23 under digest `sha256:e0ce8c35124d...`, while its `8.3.8` tag was
updated on 2026-08-18 under a different digest. Sixty-one of the hundred tags
on that page are newer than `latest`, and the same holds on `library/python`,
`library/nginx` and `library/redis`, so it is the ordinary condition rather
than a repository somebody forgot. And a tag can be marked `inactive` and
still be answered like any other: `library/registry` returns its `2.5.2` with
`tag_status: inactive`, one of fifteen in that state out of seventy, `2.5` and
`2.7.0` among them. Plain version tags, which is what makes it worth having --
somebody pinned to `registry:2.5.2` is pinned to something the registry has
already decided is a candidate for removal, and it says so in a field nobody
reads while answering 200. Fifteen are GitHub's, which answers for a public repository
without a token: its errors repeat the HTTP status in the body as a string,
an issue carries no owner and no repo, and an issue is addressed by its
number rather than by its id -- `golang/go` issue 81026 has id 5222669952,
and only one of those two works in a path. And the fourth corrected this
Recipe: it claimed the last page of a listing carries no `Link` header, and
GitHub sends one holding `rel="prev"`. A single page carries none, which is
where the wrong claim came from -- the first page and the last page are the
same page, and a single page is not a last page.

The difference is the whole reason a client reads that header. One that stops
when the header is missing never stops against GitHub, because the header is
there; it is `rel="next"` that is not. The case was teaching exactly the loop
its own comment said the header exists to prevent.

A `rel="prev"` is opt in rather than implied, because providers disagree:
Basecamp's own README describes `rel="next"` alone, so its last page really
does carry no header. Only offset and page numbering can have one at all, a
cursor being a position the caller was handed rather than a number to count
backwards from.

Asking GitHub for the rest of that Recipe's claims found one more thing it had
never modelled: `GET /issues` returns pull requests as well as issues, because
every pull request is an issue in GitHub's data model, and the only thing
telling them apart is whether a `pull_request` key is present. Not null --
absent. Twelve of the hundred open "issues" on one page of `golang/go` are
pull requests, so an open-issue count taken from that endpoint is wrong by
that much, a sync mirroring issues into a tracker imports pull requests, and
neither errors. The Recipe returned issues only, which meant code filtering
pull requests out had nothing to filter and code forgetting to filter them
looked right.

Two of that Recipe's sixteen cases still carry no date, and both are refusals
rather than oversights. Checking that a bad credential is rejected means
sending a credential-shaped header to somebody else's authentication endpoint
to watch it fail. Checking the rate-limit response means exhausting a rate
limit on purpose, which spends capacity belonging to everyone sharing this
address. The other fourteen needed nothing but a GET anybody can make.

Five are GitLab's, which answers for a public project. A merge request is
addressed by its per-project iid and not by its global id, and putting one
where the other belongs finds nothing -- or, on a busy instance, finds a
different merge request. A missing project answers one key and one only,
`{"message":"404 Project Not Found"}`, with the status repeated inside the
sentence and nothing called `error` or `errors`. A refused credential answers
the same way, `{"message":"401 Unauthorized"}`, so the status appears twice
and neither copy is a number.

The other two corrected this Recipe rather than confirming it. It had a closed
merge request answering with no `merged_at` and a merged one with no
`closed_at`. GitLab sends both dates on every merge request and nulls the one
that did not happen, so the key is there and its value is not. That is the
difference between `"merged_at" in mr` being false and being true: code asking
whether the key exists got one answer here and the other from GitLab, passed
locally, and read every merge request as merged.

Seven are Discourse's, which is the first provider added to this list since the
list was written. Any Discourse forum answers `/latest.json` and `/t/{id}.json`
without a key, so meta.discourse.org settles all five: the topics are two
levels down and there is no `topics` or `per_page` at the top level at all;
`topic_list.per_page` is thirty, the number this Recipe declares; a topic
carries `last_poster_username` as a string and neither a `last_poster` nor a
`user` object, while the names and avatars sit in a `users` array beside the
topics; `title` and `fancy_title` are the same words with the apostrophe as
`&rsquo;`, on a topic nobody wrote for the purpose; and `/t/{id}.json` answers
sixty-four keys of which `bumped_at` and `last_poster_username` are not two.

That last one is the case a description could not have settled. A description
lists what a response may carry. It does not say what a response leaves out,
and the whole point of that case is a field the listing has and the topic does
not.

Six are WordPress's, which serves any public site's posts to anybody:
`title` is an object whose only key is `rendered` and `id` beside it is a JSON
number; `categories` and `tags` are arrays of term ids with no names anywhere
in the post; a missing post answers exactly three keys -- `code`, `message`,
`data` -- with the status nested at `data.status` and nothing called `status`,
`error` or `errors` at the top; and asking for `page=0` is refused with "page
must be greater than or equal to 1", which settles the counting-from-one half
outright rather than by inference.

Two of those six were left without a date until a site could answer them, and
what the notes asked for turned out to exist. `date` and `date_gmt` are both
there with no zone marker on either, which is the shape -- and on
wordpress.org/news they are identical, because that site runs UTC, so nothing
of the trap was visible. wptavern.com runs on another timezone: its post
185079 carries `2025-01-08T22:48:02` and `2025-01-09T03:48:02`, five hours
apart and on different days, so a client parsing the first as UTC files that
post under the wrong date. Likewise, wordpress.org/news publishes two pages
and neither has a parent, which cannot show a child naming one; ma.tt's page
2545 carries `parent: 2536`, and 2536 carries `parent: 0`, so the id is an
integer pointing at another page and zero is what a top-level page says
rather than the field being absent.

One is still without a date. Eight of twenty posts carry `featured_media: 0`,
so the zero is real and common, but the case claims a sticky post as well --
and wordpress.org/news, wptavern.com, ma.tt and techcrunch.com were each
asked for `?sticky=true` and every one answered with nothing. Half a case
watched is not a case watched.

Four are Bitbucket's, whose public workspaces answer without a token. Its
listing envelope is exactly `values`, `size`, `pagelen`, `page` and `next`,
with no `data` and no `items` -- the two names a client reaches for first,
both finding undefined rather than an error -- and `size` was 407 against a
`pagelen` of 2, so it counts everything rather than what arrived. A repository
carries `is_private` as a boolean and no `visibility` key at all. A missing
repository answers two keys, `type` and `error`, and `error` holds nothing but
`message`.

The fourth is the one worth reading twice. `?pagelen=1` answers one value.
`?limit=1` answers ten, with `pagelen: 10` in the body: the wrong name is not
refused, it is ignored, and a full page comes back looking like a successful
request for one. That is the failure this README keeps describing, watched
happening.

Checking Bitbucket also corrected a message this project had invented. Its 404
had a generic "Resource not found", and Bitbucket's own words are that the
repository may not exist *or* may not be visible to you, and it declines to
say which. Code treating a 404 as a deletion -- dropping a row, ending a sync,
marking a connection dead -- is acting on a message that explicitly refuses to
support it, and the case it gets wrong is a repository somebody made private
this morning.

Thirteen are PyPI's, which publishes no OpenAPI for its JSON API and needs no
credential to answer, so the whole Recipe was written from live responses and
every case carries the date it was checked. The headline is three download
counters that always say minus one -- `"downloads": {"last_day": -1,
"last_month": -1, "last_week": -1}` -- kept after Warehouse stopped serving
download statistics through this API, with a fourth on every file. The rest
were settled the same way: `releases` is on `/pypi/{project}/json` and absent
from `/pypi/{project}/{version}/json`, so the key a client walks to enumerate
versions is missing from the more specific endpoint; three of requests'
releases map to empty arrays, so counting them counts versions nobody can
install; `core-metadata` is the one hyphenated key among snake-cased
neighbours; `md5_digest` and `digests.md5` carry the same hash and
`upload_time` and `upload_time_iso_8601` the same instant at two precisions;
and 1.26.20 is a real release of `urllib3` that answers 404 under `requests`,
which is what makes the project segment load-bearing rather than decorative.

Eleven are crates.io's, written the same way and for the same reason: no
OpenAPI, no credential, so every claim was checked live. Four fields say which
version is the latest -- `max_version`, `newest_version`, `max_stable_version`
and `default_version` -- and on `serde` all four read 1.0.229, which is why
nobody notices. On `rand` the field called `newest_version` is 0.8.8 while
`max_version` is 0.10.2: newest means most recently published and max means
highest by semver, and a patch to an old line went up on 25 August against a
0.10.2 from 2 July. A request with no `User-Agent` is a 403 that has nothing to
do with permission -- "we ask that your user agent actually identify your bot"
-- on a registry that needs no credential at all. And `id` is the crate's name
on a crate and an integer on a version, in one response.

Nine are RubyGems', and the same reason again: no OpenAPI, no credential for
reads, every claim checked live. The flag that decides whether publishing needs
two-factor authentication is the string `"true"`, because a gemspec's metadata
is a map of strings to strings -- so `if (metadata.rubygems_mfa_required)` is
true for `"true"` and also for `"false"`. A gem nobody has published answers a
404 whose entire body is `This rubygem could not be found.`, in plain text, from
a path ending in `.json`. `authors` is prose and `licenses` is a list, so
nokogiri's four authors arrive as one comma-joined string. A dependency
requirement is a sentence -- `{"name": "actioncable", "requirements": "=
8.1.3.1"}`. And the gemspec's URIs are promoted to the top level and left in
place, agreeing three times out of four: `homepage_uri` is at the top and absent
from `metadata`, so neither copy is authoritative.

Nine are Repology's, and it is the one registry here that does not answer a
missing thing with a 404: asking for a project nobody packages is a **200 with
an empty array**, so `if (!response.ok) throw` never fires and a typo is
indistinguishable from a project packaged nowhere. One project's response is a
bare array with an entry per packaging -- 806 of them for `curl` -- and only
five of the fourteen fields are on every entry. `srcname` is missing from
twelve, `vulnerable` from a hundred and seventy-four, and `binname` and
`binnames` are the same idea singular and plural with no entry carrying both.
Six statuses appear at once: 377 outdated, 255 legacy, 163 newest, 7 rolling, 3
devel and 1 `incorrect`, which is the index's own judgement that a version
string could not be read as a version, sitting in an ordinary entry beside the
rest.

Eight are Homebrew's. An object called `versions` holds two strings and a
boolean -- `{"stable": "8.21.0", "head": "HEAD", "bottle": true}` -- so one of
its three values is a version, `head` is a git ref spelled the same way on every
formula that has one, and `bottle` says whether a prebuilt binary exists. Four
more fields describe a computer that is not involved: `installed` is `[]`,
`pinned` and `outdated` are `false` and `linked_keg` is `null` on everything the
API serves, because it is a static document on a CDN and the same schema is what
`brew info --json` prints locally. A missing formula answers a **full HTML
page** from a path ending in `.json`, which is the third way this collection has
seen that refused. And Ruby symbols arrive as strings carrying their colon:
`keg_only_reason.reason` is `":provided_by_macos"` and every bottle's `cellar`
is `":any"`.

Eleven are Packagist's, and it is the sharpest of the registry findings: the
versions array is **a chain of diffs**. `GET /p2/monolog/monolog.json` returns
eighty-seven entries where the first has twenty-one keys and the second has
eight, because each entry after the first carries only the fields that differ
from the one before it. So `packages["monolog/monolog"][1].license` is undefined
-- not because 3.9.0 has no licence, but because it has the same one as 3.10.0.
The chain runs **newest first**, so the deltas apply backwards in version order.
A field that goes away is the string `"__unset"`. And the only signal that any
of it is happening is a sibling key at the top of the document, `"minified":
"composer/2.0"` -- nothing inside an entry says it is a patch.

Eight are Hex.pm's, where the API hands you the dependency line to paste and the
three it hands you disagree. For `phoenix` 1.8.13 the `configs` object is
`mix.exs: {:phoenix, "~> 1.8"}`, `rebar.config: {phoenix, "1.8.13"}` and
`erlang.mk: dep_phoenix = hex 1.8.13` -- the first a range that will pick up
1.8.14 on the next resolve, the other two exact pins. So the same request tells
an Elixir project to float within a minor line and an Erlang project to freeze
on a patch, and the only thing choosing between those policies is which key you
copied. Also: four download counters where `all`, `day` and `week` say what they
measure and `recent` does not, with no `month` beside it to infer from; a
`meta.maintainers` that is always empty while `owners` holds the actual people;
and an advisory carrying three identifiers for one finding.

Twelve are NuGet's, and it is the one where the same package answers two
different ways at once. The service index advertises the registration resource
under four versions of its type pointing at three base URLs, and a client picks
by asking for a version of the *type*: resolve the unversioned name and you get
`registration5-semver1`, which for `Microsoft.Extensions.Logging` is 98
versions. The `semver2` base URL has 175. The seventy-seven missing are the
SemVer 2.0.0 ones -- every preview of the current major and the one before it --
because SemVer 1 cannot express a dotted prerelease identifier, and nothing in
the response says anything was left out. Whether the index arrives whole or in
pieces is decided by size rather than by endpoint: 84 versions come back
inlined with fragment links, 600 come back as pages with real URLs, and
`Microsoft.Extensions.Logging` is on both sides of that line at once. And asking
for a package that does not exist reports that a *blob* does not exist, in
Azure Storage's XML, from a JSON API.

Nine are the Go proxy's, which refuses the module path in the URL and
then hands it back in the body. A module path is case-sensitive and the proxy
protocol will not carry capitals, so `github.com/BurntSushi/toml` -- the
spelling in every `go.mod` -- is fetched as `github.com/!burnt!sushi/toml`, and
asking for the author's spelling answers **404 with the sentence "bad request"**.
The escaped form then returns `Origin.URL: "https://github.com/BurntSushi/toml"`,
capitals restored, from the server that would not take them. The JSON keys are
Go struct field names -- `Version`, `Time`, `Origin`, `VCS`, `URL`, `Ref`,
`Hash` -- because the type is marshalled straight out of `cmd/go` with no tags.
And a module nobody has published answers with the proxy's own shell failure:
the `git ls-remote` it ran, a path inside its cache on its own disk, `exit
status 128`, git's "terminal prompts disabled", and two lines of advice.

Nine are Maven Central's, where the search API is Apache Solr's response handed
over unedited. A successful search reports `"status": 0` -- Solr's convention,
not HTTP's -- inside a body that arrived with a 200, beside `QTime` and a
`params` object echoing the caller's query along with the internal field list
the service asked Solr for and the sort expression it used. The results
themselves are single letters: `g`, `a`, `v`, `p` and `ec` for groupId,
artifactId, version, packaging and the file extensions that exist, and nothing
in the response says so. `start` appears twice in one document in two types --
`0` where the service computed it, `""` where the caller's value is read back as
text. A search matching nothing is a 200 that says it succeeded, so a mistyped
query and a missing artifact are the same response. And a query Solr will not
parse is plain text that stops mid-sentence: `Solr returned 400, msg:` with
nothing after the colon, an error relayed by a service that did not read it.

Thirteen are OSV.dev's, where one vulnerability is three records that disagree.
The SQL injection fixed in Django 2.2.28, 3.2.13 and 4.0.4 is in the database
under three ids from three sources, each naming the others in `aliases`:
`GHSA-2gwj-7jmv-h26r` has a summary, a severity and three `affected` entries;
`PYSEC-2022-190` has no summary and no severity; `CVE-2022-28346` has neither
and **no `package` at all** -- no name, no ecosystem, no purl, and a `GIT` range
whose events are commit hashes. So the shape of the answer depends on which name
you looked it up by, and the CVE, the name a scanner is most likely to hold, is
the emptiest of the three. They are not even on one schema version: the CVE
record says `1.9.0` and the other two say `1.7.3`.

**And the records disagree about which version fixed it.** GitHub's splits the
three release branches into three entries with two events each, ascending;
PyPA's packs all three into one entry with six events, descending. So
`affected[0].ranges[0].events[1].fixed` is `2.2.28` on one and `4.0.4` on the
other, for the same vulnerability in the same ecosystem. A query counts each
advisory once per source, too: django 3.2.0 answers 63 vulns of which 26 are
pairs of the same advisory, so there are 37, and the array is sorted by id so
the halves of a pair sit forty entries apart. `versions` is sorted as text --
`2.2.10` before `2.2.2` -- so its last entry is `2.2.9` while the highest
affected release is `2.2.27`. No vulnerabilities is a 200 whose body is `{}`
rather than `{"vulns": []}`. And `code` is a gRPC status on one failure and an
HTTP status on another: `5` beside a 404, `3` beside a 400, and `400` beside a
400 when the body will not parse.

Twelve are deps.dev's, where capitalising a package name gets you a different
package and a 200. On npm the name is case-sensitive, so
`/v3/systems/npm/packages/express` answers with 288 versions and
`/v3/systems/npm/packages/Express` answers with three -- a real package, three
releases in 2016, all deprecated, whose deprecation text is the only thing in
the response that says anything is wrong. `packageKey.name` echoes `Express`
back as though it were right and `isDefault` is set on one of the three, so a
client that title-cased a name reads a nine-year-stale version history and gets
no error at all. **The other ecosystems do not agree**: on PyPI the same request
normalises instead, and `/v3/systems/pypi/packages/Django` answers with
`packageKey.name: "django"` -- the name you asked for silently replaced with the
one it resolved to. One API, one path, two opposite readings of the same
mistake.

The rest of it is the same kind of quiet. The system is echoed in upper case
(`NPM`, `PYPI`) where the path took lower. A version that is not deprecated
carries `deprecatedReason: ""` -- present and empty, the opposite of the npm
registry above, whose `deprecated` is absent rather than false. The dependency
graph is nodes and edges with integer indices rather than a tree, the nodes are
in **alphabetical** order rather than dependency order, and the package itself
is node 0 with relation `SELF`, so the array is one longer than the dependency
count and filtering it makes every edge point somewhere else. A successful graph
says so with `error: ""`. And the failures are bare plain text: seventeen bytes
reading `package not found`, `text/plain`, no trailing newline, no JSON --
which is also what an ecosystem that does not exist gets, naming the wrong one
of the two things in the URL.

Eleven are Open-Meteo's, where the timestamps carry no offset and the same
string means two different hours. Ask for Toronto with a timezone and the hourly
times read `2026-08-28T00:00`; ask for the same coordinates without one and they
read `2026-08-28T00:00` again, four hours earlier, because the default timezone
is GMT. The first is 17.5 °C and the second is 21.7 °C. There is no `Z`, no
`+00:00` and no seconds, so `new Date()` on either parses as local time in
whatever zone the reader is in -- a third answer -- and the only way to be right
is to read `utc_offset_seconds` from the other side of the document and apply it
by hand. **And the coordinates come back moved**: a request for 43.65, -79.38
answers with 43.646603, -79.38272, the nearest cell of the weather model,
silently substituted.

The rest is the same quiet. The hourly data is **parallel arrays** -- `time`,
`temperature_2m` and `precipitation` as three lists whose only relationship is
the index -- with the units in a different object, where `hourly_units.time` is
`"iso8601"`, a serialisation format described as a unit for a field that has
none. Adding a second coordinate changes the top-level type from an object to an
array, and `location_id` is then absent on the first element rather than 0.
`timezone_abbreviation` is `"GMT-4"` rather than `EDT`. And a failure names a
field `error` that is only ever `true`, beside a `reason` that is a Swift
decoder's stringified error, complete with the service's internal generic type
signature -- while a forgotten `latitude` is reported as `Parameter 'latitude'
and 'longitude' must have the same number of elements`.

Eleven are USGS's, where asking for a page removes the count. The earthquake
catalogue's `metadata` block changes shape with the request: ask without a limit
and it carries `count`; ask with one and `count` is gone, replaced by `limit`
and `offset`. So a client that pages and reads `metadata.count` to know when to
stop finds `undefined` -- on exactly the requests where it is paging, and only
there. **And `offset` is one-based**: a request that sent no offset at all comes
back saying `offset: 1`, so anything treating it as a zero-based cursor and
adding the page size skips a record on every page after the first.

The features themselves are GeoJSON, with GeoJSON's traps. `coordinates` is
`[longitude, latitude, depth]` -- longitude first, and the third element is
depth in kilometres -- so `[-103.4604, 8.4191, 10]` read as a latitude-longitude
pair puts the epicentre in the wrong hemisphere without erroring. `ids`,
`sources` and `types` are comma-delimited strings that **begin and end with a
comma**, so a naive split gives five elements for three identifiers. The key
`type` appears four times in one document with four vocabularies:
`FeatureCollection`, `Feature`, `Point`, `earthquake`. `mag` is a JSON number,
so a magnitude of 5.0 arrives as `5` while the title the service built beside it
says `M 5.0`. `tz`, `felt`, `cdi`, `mmi` and `alert` are present and null rather
than absent. And a request that said `format=geojson` gets its failures as
**plain text** -- a multi-line human-readable report, `text/plain;charset=UTF-8`
with no space after the semicolon -- which reads the request back with its
ampersand written as `&amp;`, an HTML entity in a plain text body.

Nine are Frankfurter's, where asking for a Sunday gets you Friday with a 200 and
no word about it. `GET /v1/2026-08-23` answers `{"date": "2026-08-21", ...}`,
because the European Central Bank publishes on working days and this service
falls back to the most recent fixing. The only thing that says the date moved is
the `date` field, which a client that already knows what it asked for has no
reason to read -- so a two-day-old rate arrives labelled as the rate for the day
you wanted, and over Easter that is four days. **And the base is never among the
rates**: with `base=USD` the object holds EUR, GBP and JPY and no USD, so
`rates[base]` is undefined, and `symbols=USD,EUR` answers with EUR alone -- the
currency removed from your own list without a word.

A range changes the shape under the same key names: `date` becomes `start_date`
and `end_date`, and `rates` goes one level deeper, so `rates.EUR` works on one
and is undefined on the other. The range you get is the range that has data --
asking to the 31st of December answers to the 27th of August, silently clamped
-- and the days between are simply absent, so iterating dates and indexing
`rates[date]` finds nothing on two days in seven. A currency that does not exist
and a date before the series began are the same 404 with the same
`{"message": "not found"}`, naming neither. But a **path** that does not exist
answers `{"status": 404, "message": "not found"}` -- the same words with a status
added, so code branching on `body.status` finds a number when it mistyped the
URL and undefined when it mistyped the currency.

Nine are Nominatim's, where nothing found is an empty array on one path and an
error object on the other, and both are 200. `GET /search?q=...` with no matches
answers `[]`; `GET /reverse?lat=0&lon=0` with no match answers
`{"error": "Unable to geocode"}`. Same service, same `format=jsonv2`, same
status, and the reverse one is the expensive shape: an object with an error key
where an object of results should be, so `response.lat` is undefined and
`response.ok` is true. **And search answers with an array while reverse answers
with a bare object** -- not an array of one, an object -- so the two endpoints of
one geocoder cannot share a parser.

The rest is what OpenStreetMap's data model does to a JSON API. `lat` and `lon`
are **strings**, so arithmetic concatenates and comparison sorts lexically.
`boundingbox` is four strings in the order south, north, west, east -- neither
GeoJSON's `[west, south, east, north]` nor a pair of coordinates in sequence.
The keys of `address` depend on the place, and one of them carries an
administrative level inside its own name: `ISO3166-2-lvl4`, where the 4 is the
depth at which that country keeps its subdivisions and is 3 or 6 elsewhere.
`importance` arrives in scientific notation. `place_id` is local to the
installation and changes when the database is rebuilt, so the stable identity is
`osm_type` and `osm_id`, two fields. And a request without a `User-Agent` is
refused in **plain text** -- 403 and one line pointing at the usage policy --
from an endpoint that was asked for JSON.

Eight are Open Library's, where the answer is keyed by the string you sent and a
miss is that key's absence. `GET /api/books?bibkeys=ISBN:0451526538` answers
`{"ISBN:0451526538": {...}}`, so reading the response means rebuilding the query
string -- prefix, colon and all -- to index it; and a lookup that matched nothing
answers `{}` with a 200, not a null under that key and not a 404. So the
difference between found and not found is whether a key you constructed yourself
is there, and `data["ISBN:" + isbn].title` throws a TypeError on the miss without
saying why. **And the same field name carries a reference on one endpoint and a
resolved object on the other**: the canonical document answers
`"authors": [{"key": "/authors/OL18319A"}]` while `/api/books`, for the same
book, answers `[{"url": ..., "name": "Mark Twain"}]`. One is a pointer that costs
another request and nothing says which you are holding.

The rest is the data model of a wiki on a JSON surface. `/isbn/{isbn}.json` answers a
**302 when it finds the book** and a 404 when it does not, so the redirect is the
happy path -- and that 404 is a full HTML page served as `text/html` from a path
ending in `.json`. Timestamps are typed objects,
`{"type": "/type/datetime", "value": "2008-04-01T03:28:50.625462"}`, with
microseconds and no timezone at all. `type` is a reference too. `revision` and
`latest_revision` are two fields carrying one number. Every identifier is an
array, including `"openlibrary": ["OL1017798M"]`, which cannot repeat. And there
are **two schemes in one document**: `url` and `authors[].url` are `http://`
while `subjects[].url` and every cover URL are `https://` -- the record links are
the insecure ones.

Eight are Wikipedia's, where a page that does not exist is reported as an
internal error. `GET /api/rest_v1/page/summary/Zzzznotarealpagexyz` answers
`{"status": 404, "type": "Internal error"}` and that is the whole body: nothing
names the page, and the `type` says the same thing for a title nobody created as
for a genuine failure inside the service, so a client branching on it cannot tell
a 404 it should expect from a 500 it should alert on. **And the title arrives
five times in four fields, two of them HTML** -- `title`, `displaytitle`, and
`titles` holding `canonical`, `normalized` and `display`, where `displaytitle`
and `titles.display` are the same markup twice and nothing says which a heading
should use.

A disambiguation page is a 200 whose `extract` is a flattened list:
`"Mercury most commonly refers to:Mercury (planet), the closest planet to the
Sun\nMercury (element), a chemical element"` -- no space after the colon,
newlines between the entries, neither a sentence nor a list. The spelling you
asked for is nowhere in the answer, so asking for `toronto` gives `Toronto` in
every field that could have recorded the correction. `originalimage.source` ends
`3840px-...jpg` while `originalimage.width` says 6632, so laying out from the
declared size and loading the URL gets a picture forty per cent smaller -- and
those image URLs carry `utm_source`, `utm_campaign` and `utm_content`
parameters. `content_urls` has `desktop` and `mobile` with the same four keys,
three of them byte-identical. `revision` is the string `"1370998502"` where
`pageid` is the number 64646. And the main namespace's `text` is the empty
string.

Nine are Hacker News's, where a thing that is not there is the four bytes
`null`, with a 200. Not an empty object, not a 404, not an error field: the JSON
literal, sent as `application/json`, from both the item endpoint and the user
one. So `response.ok` is true, `.json()` succeeds, and what comes back is the
value that means "no value" in the language you are writing in -- and
`(await res.json()).title` throws a TypeError with no status anywhere to have
branched on first. **And two of the endpoints have no envelope at all**:
`/v0/maxitem.json` answers a bare integer and `/v0/topstories.json` answers a
bare array of five hundred integers, so rendering a front page is one request for
the list and five hundred for the items, with no batch endpoint anywhere.

The rest is one resource doing five jobs. `id` is a number on an item and the
username on a user. `text` is HTML inside a JSON string with no plain-text
companion, so anything rendering it strips tags and anything not rendering it
shows them. `time` is epoch **seconds**, ten digits where most of this collection
is thirteen. A story has `title`, `url`, `score` and `descendants`; a comment has
`parent` and `text` and none of those, and nothing but `type` says which arrived.
And `kids` is the direct replies while `descendants` is the total, so neither is
the length of the other.

Nine are PokeAPI's, where the description is still formatted for a Game Boy.
`flavor_text` comes back as
`"A strange seed was\nplanted on its\nback at birth.\fThe plant sprouts\nand
grows with\nthis POKéMON."` -- three things in one string. The newlines are hard
wraps at the original 1996 screen's column width, so re-flowing them in a modern
layout breaks mid-phrase. The `\f` is a **form feed**, U+000C, the page break
from the handheld's text box, sitting in a JSON string today. And the name is
spelled with a lowercase é between two capitals, because the font in those games
had no capital É. None of it is escaped, flagged or duplicated in a cleaned-up
field.

**And there are twenty-eight English descriptions, not one.**
`flavor_text_entries` is a flat array of a hundred and two, one per game version
per language, each carrying `language` and `version` as references -- so "the
English description" is a filter that returns twenty-eight answers and picking
the first is picking Red from 1996. Around that: a miss is the bare plain text
`Not Found` with a 404, from an API whose every success is JSON; a listing row is
a `name` and a `url` and **no identifier at all**, so the id exists only as a
path segment inside that URL and anything more than the name costs a request per
row; `previous` and `next` are both always present and each is null at its own
end, so `"next" in body` is true on every page; the lookup is case-insensitive
and says nothing about it; a type is `{slot, type: {name, url}}`, so the name is
two levels down and the array index is not the slot; and four of the ten
top-level `sprites` keys are present and null.

Eight are MusicBrainz's, which answers XML unless you ask for JSON and explains
that in XML. The default on `/ws/2/artist/{mbid}` is
`<?xml version="1.0" ...>`; `?fmt=json` gets a JSON document; and `?fmt=yaml`
gets **406 in application/xml**, carrying the message that the recognised types
are application/json and application/xml. The one answer that would tell you how
to ask for JSON is the one you cannot parse without already having solved the
problem. **And an identifier that does not exist is a 400 about the identifier's
format**: a perfectly well-formed UUID that nothing in the database uses answers
`{"error": "Invalid mbid."}`, exactly as a string that is not a UUID at all does,
so "no such artist" and "you sent nonsense" are one answer and there is no 404 in
the shape. The two bodies arrive with their keys in opposite orders.

A request without a `User-Agent` is **403 Forbidden with a message about being
throttled** -- the status says you are not allowed, the text says you are going
too fast, and the cause is neither -- and that body is `application/json` with no
charset where every other JSON answer declares one. The field names are
hyphenated (`sort-name`, `type-id`, `begin-area`, `iso-3166-1-codes`), so in
JavaScript every one has to be reached with brackets. `life-span` carries
`begin: "1987"` beside `end: "1994-04-05"` -- two precisions, both strings,
neither saying which. And there are four spellings of "not applicable" in one
document: `gender` is null, `end-area` is null, `ipis` is `[]`, and inside `area`
the `disambiguation` is `""` while `type` is null.

Seven are Open Food Facts's, where a product that is not there is a 200 with
`status: 0`. The HTTP status says it worked and the truth is a number in the
body, so `response.ok` is true and there is **no `product` key at all** --
`body.product.product_name` throws a TypeError on the answer a barcode scan is
most likely to get. And `status: 0` has more than one cause: `product not found`
and `no code or invalid code` share the number, so a client branching on it
cannot tell a barcode that does not exist from one it should not have sent.
**And the code you sent is not always the code that comes back**: asking for
`0000000000000`, thirteen zeros, answers `"code": "00000000"` -- eight -- with no
field anywhere recording what was asked for.

Then the nutrients. `energy` is **kilojoules**: the same object carries
`energy: 2252` and `energy-kcal: 539`, and the unqualified name is the kJ figure,
so reading `energy` and labelling it calories is out by a factor of 4.18. Every
nutriment is repeated -- `fat`, `fat_100g`, `fat_value` and `fat_unit` are four
keys for one number and its unit, and energy has twelve. The key names mix
separators inside one key: `added-sugars_100g`. `brands` is a comma-joined string
rather than a list, crowd-sourced, and one of the three on the reference product
is not a brand. And `nutriments_estimated` is a parallel object of the same
nutrients, computed rather than declared, keyed identically, with nothing inside
either saying which is which.

Nine are Crossref's, where a date is an array inside an array and it is not
always the same length. `"issued": {"date-parts": [[2013, 7, 31]]}` on one work
and `[[1970, 6]]` on another: the outer array is there because the field can hold
a range, and the inner one stops wherever the metadata stopped. Nothing says
which precision arrived, so `date-parts[0][2]` is a number on one work and
`undefined` on the next, and building a Date from three positional arguments
silently reads the month as the day. **And one work has three publication dates
at two precisions** -- `published` and `published-online` are `[[2013, 7, 31]]`
while `published-print` is `[[2013, 8]]`, so "when was this published" has three
answers, one of them a month, and the print date is after the online one.

`title` is an array: a paper has one title and it arrives in a list of one, as do
`container-title`, `short-container-title` and `ISSN`. `created` carries the same
instant three ways -- `date-parts`, `date-time` and a `timestamp` in epoch
milliseconds -- where `published` beside it carries only the first. The envelope
is `status`, `message-type`, `message-version` and `message`, with a version of
the envelope rather than the API that has read `1.0.0` for years and a `status`
of the string `"ok"` beside an HTTP 200. A single-work lookup carries a
relevance `score` for a query it did not make. And a DOI that does not resolve is
the bare plain text `Resource not found.` with a 404 and no charset -- no
envelope, no status field, from an API whose every success is that four-field
object.

Nine are Zippopotam's, where the keys have spaces in them. `"country
abbreviation"`, `"post code"`, `"place name"`, `"state abbreviation"` -- not
camelCase, not snake_case, not kebab-case: a space, in the key. Nothing in the
response can be destructured, nothing can be reached with a dot, and every access
is a bracket with a quoted string in it. **And `"CA"` is two different things in
one API**: the country abbreviation for Canada and the state abbreviation for
California, in fields whose names differ by one word.

`latitude` and `longitude` are strings, and with a minus sign in front the
lexical order is not close to the numeric one. A postcode that does not exist is
`{}` with a 404 -- an empty object, so `.json()` succeeds and every field is
undefined -- and a country code that does not exist answers exactly the same way,
as does a real postcode under the wrong country. `places` holds objects of **two
different shapes** depending on which endpoint answered: looking up a postcode
gives entries with `state` and `state abbreviation`, looking up a place gives
entries with `post code` instead, under the same key with nothing saying which
arrived. And the reverse lookup repeats `"place name"` at both levels, on the
request that named it in the first place.

Eight are TVmaze's, where the failure puts the reason in `name` and leaves
`message` empty. A missing show answers
`{"name": "Not Found", "message": "", "code": 0, "status": 404}` -- four fields,
and the one a client would read is blank. `name` carries the reason phrase, and
on every successful response `name` is the title of the show, so the same key
holds "Under the Dome" on a 200 and "Not Found" on a 404. `code` is 0, which is
neither an HTTP status nor an error code, and `status` is the HTTP status
restated inside the body.

**And "where does it air" is two mutually exclusive fields.** Every show carries
`network` and `webChannel` and exactly one is an object, so `show.network.name`
throws on half the catalogue and the check has to be written both ways round --
and inside them the shapes disagree, the network carrying a country object where
the web channel carries `country: null`. Around that: `runtime` and
`averageRuntime` are two fields where one is often null and nothing says which to
prefer; `schedule.time` is the empty string on a show that has a broadcast day;
`externals` holds three identifiers in three types, `null`, a number and the
string `"tt9140554"`; `summary` is HTML in a JSON string with no plain-text
companion; and `updated` is epoch seconds.

Eight are Sunrise-Sunset's, where the sunset is before the sunrise. A request
for Toronto answers `"sunrise": "10:35:37 AM"` and `"sunset": "12:01:46 AM"`:
both UTC, both formatted as a twelve-hour wall clock with no date and no zone.
Toronto's sunset that evening is a minute past midnight UTC -- the next calendar
day -- so parsed against the date that was asked for, the sunset lands sixteen
hours **before** the sunrise, and nothing in the response says the day rolled
over. **And the times are UTC while looking like local time**: `"10:35:37 AM"`
reads as half past ten in the morning, Toronto's sunrise was 06:35, and the only
thing that says otherwise is `tzid: "UTC"`, a sibling of `results` rather than a
field inside it.

A latitude of 999 answers `status: "OK"` with a full set of times -- there is no
validation of the coordinate at all. A failure sets `results` to the **empty
string**, where a success has it as an object, so `typeof results` differs
between them. `day_length` is `"13:26:09"`, a duration in the same
colon-separated shape as the times beside it and the only field in the object
that is not a clock reading. And at 78 degrees north on the solstice every time
is `"12:00:01 AM"` and `day_length` is `"00:00:00"`, with `status: "OK"` -- the
sun does not set there in June, and the API reports a zero-length day.

Eight are NHTSA's vPIC, where four errors arrive in one string and one of them
is 400 on a 200. `ErrorCode` is `"6,7,11,400"` -- a comma-separated list of
numbers inside a string, four codes for one VIN, and `400` among them looks like
an HTTP status and is not, because the response is a 200 and always is.
`ErrorText` is those four joined with `"; "`, and one of them contains its own
semicolon, so splitting on the separator gives five pieces for four errors.
**And the failure is a success with 148 empty strings in it**: `Count` is 1, the
`Results` array holds one object, and 148 of that object's 154 fields are `""` --
nothing absent, nothing null, so a VIN the service could not read looks exactly
like one it read and found nothing about.

Every value is a string: `ModelYear` is `"2003"`, `EngineCylinders` is `"6"`, and
`DisplacementL` is `"2.998832712"`, a three-litre engine to nine decimal places
as text. `Message` is a 250-character disclaimer on every successful response,
explaining that a missing value does not mean the feature is absent.
`SearchCriteria` is the request as prose -- `"VIN(s): 1HGCM82633A004352"`. The
error fields are populated on success, where `ErrorCode` is the string `"0"`. And
`SuggestedVIN` on a bad VIN is `"N!TAV!N"`: the input with its invalid characters
replaced by exclamation marks, which is not a VIN and cannot be submitted
anywhere. The one failure status on the whole API is for a path that does not
exist, and its message quotes the URI as **the backend saw it** --
`backend-vpic-api.nhtsa.dot.gov`, a host the caller never named -- under a
`message` key spelled in lower case where every success spells it `Message`.

Eight are Wikidata's, where the English label is not there and the name is under
a code that is not a language. Douglas Adams's item carries 75 labels and `en` is
not one of them -- nor `de`, `fr`, `es`, `it`, `nl`, `pt`, `pl` or `sv`, because
every Latin-script language was consolidated into a single entry under `mul`,
which is ISO 639-3 for "multiple languages". Asking the label endpoint directly
confirms it: `/labels/en` is a 404 and `/labels/mul` is `"Douglas Adams"`. So
`item.labels[userLang] ?? item.labels.en`, which is what every client writes,
shows nothing at all -- while `item.descriptions.en` beside it works, because the
descriptions were not consolidated. Two sibling objects, keyed the same way,
disagreeing about which languages exist.

**And every key in `statements` is an opaque number** -- 312 of them on this
item, and nothing in the response says what any means. A value is
`{"type": "value", "content": ...}` where `content` is a string for an item
reference and an object for a date; a date is `time` with a leading `+` on the
year, an integer `precision` nothing decodes, and a `calendarmodel` that is a URL
you must resolve to learn the calendar is Gregorian. A statement identifier is
`Q42$F078E5B3-...` -- the item, a dollar sign and a UUID -- and the item's own
case is not consistent across them: 285 begin `Q42$` and 27 begin `q42$`. A
missing item and a missing label are the same body with one nested field
different.

Eight are GBIF's, where a name it could not find comes back at a hundred per cent
confidence. `?name=Zzzznotaspeciesxyz` answers
`{"confidence": 100, "matchType": "NONE", "synonym": false}` with a 200 --
`confidence` is the highest the scale goes, on the one answer that found nothing,
so code ranking results by confidence puts the failures first. Every real match
scores lower: an exact one is 99, a fuzzy one 95, a genus standing in for a
species 94. **And `synonym` is present only when nothing was found**: the
no-match carries `synonym: false`, and the response whose `status` actually is
`SYNONYM` does not carry the field at all.

A deprecated name matches **EXACTly**. `Felis concolor` is what cougars were
called until 1993, and asking for it answers `matchType: "EXACT"` at confidence
98 with `status: "SYNONYM"` -- exact because the string exactly matched a name
nobody should use. That response carries two keys for two different taxa:
`usageKey` is the synonym and `speciesKey` and `acceptedUsageKey` are the
accepted species, so the field with the plainest name is the one not to store. A
misspelled species silently becomes a genus: `Puma notaspecies` answers
`matchType: "HIGHERRANK"` at 94 with `scientificName: "Puma Jardine, 1834"`,
which a client will print as though it were the species asked for. A typo is
corrected without saying what changed. And `/v1/species/notarealendpoint` answers
400 with the bare text `For input string: "notarealendpoint"` -- Java's
`NumberFormatException`, verbatim, because `/v1/species/{key}` takes an integer
and the segment was handed to a number parser.

Eight are Nager.Date's, where `fixed` is false on Christmas Day. The
twenty-fifth of December is as fixed as a date gets, and the field that says
whether a holiday falls on the same date every year says no -- as it does for New
Year's Day, and for every other holiday: across Canada, the United States, the
United Kingdom and Germany, eighty-one holidays in 2026, `fixed` is false
eighty-one times. `launchYear` is null on all eighty-one too. Both are fields
with one value, so branching on either is branching on a constant.

**And one date carries six holidays, none of them national.** The third of August
2026 in Canada is the Civic Holiday, British Columbia Day, Heritage Day, New
Brunswick Day, Natal Day and Saskatchewan Day -- six entries, one date, thirteen
provinces and territories between them, every one `global: false`. So
`holidays.find(h => h.date === today)` returns whichever the array ordered first,
and calling it "the holiday" is wrong in five provinces out of six. The field
naming those subdivisions is called `counties`, holds Canadian provinces and
German states, and is `null` rather than empty when the holiday is national. And
the two failures share an envelope and disagree about the rest of it: an unknown
country carries `detail`, an unsupported year carries `errors`, and the `title`
on the second is a framework's default validation sentence with the real message
two levels down inside an array.

Seven are the Wayback Machine's, where asking whether a page existed in 1990
answers yes, from 2002. `closest` means closest available, `available` is `true`,
`status` is `"200"`, and nothing in the response says the snapshot is twelve
years from the date that was asked for. The distance is computable, but only by
subtracting two timestamps in two different formats, and no field states it.

**And `timestamp` is in the document twice, meaning two different things.** The
one inside `closest` is when the snapshot was taken, fourteen digits. The one at
the top level is the date the caller asked for, eight digits, echoed back -- and
absent entirely when the caller sent none. So the same key is a request
parameter in one place and an answer in the other, at two lengths, one of them
conditional. `url` is doubled the same way. `status` is the string `"200"`, so
`=== 200` is false and `== 200` is true. A URL nothing has archived answers
`archived_snapshots: {}` with a 200, where the hits have `closest` -- and a
request with no `url` at all answers the plain text `Error: no url parameter`,
also with a 200, in `text/html`, from an endpoint whose successes are JSON.

Nine are the National Weather Service's, where the path says lat,lon and the
body says lon,lat. `GET /points/39.7456,-97.0892` answers
`"coordinates": [-97.0892, 39.7456]` -- the same pair, in the same exchange, in
the opposite order, because the path takes latitude first the way a person
writes a coordinate and the body puts longitude first the way GeoJSON requires.
Nothing in the response says the order changed, and near the equator or the
prime meridian reading `coordinates[0]` as what you sent gives a plausible place
rather than an error.

**And `type` is in that document three times, in three vocabularies.** The top
level says `Feature`, which is GeoJSON's; `properties.@type` says `wx:Point`,
which is the NWS ontology's; `properties.type` says `land`, which is the NWS's
own classification -- and only the middle one is namespaced. A number with a
unit is expressed four ways in one forecast: `elevation` is a `unitCode` object,
`temperature` is a bare number beside a separate `temperatureUnit`, `windSpeed`
is the prose `"10 to 15 mph"` and sometimes the scalar `"15 mph"`, and
`probabilityOfPrecipitation` is a `unitCode` object for a percentage.
`validTimes` is an interval whose right half is a duration --
`2026-08-28T16:00:00+00:00/P7DT9H` -- so `Date.parse` returns `NaN`. A request
with no `User-Agent` is refused by the CDN rather than the API, 403 in HTML,
with no `correlationId` to quote to anybody. And three failures are 404 and mean
three different things: the one whose title is accurate is the one carrying the
least, while a latitude of 999 says only "Not Found" and hides the reason in a
`parameterErrors` array RFC 9457 does not define -- an array whose message, for
an unknown forecast office, is 840 characters long because it inlines all 133
valid office codes into an English sentence.

Nine are OpenTDB's, where the success key is `results` and the rate limit's key
is `result`. One letter, on one of four answers. Code reading
`body.results.length` works against a hit, against no results and against an
invalid parameter, and finds `undefined` on the fourth -- and the fourth is the
one that arrives when the caller is going too fast, so the crash appears under
load and not before it.

**And three of the four are HTTP 200.** No results is 200. An invalid parameter
is 200. `response.ok` is true for every outcome except the rate limit, and the
only thing separating them is a number in the body that has no name anywhere in
the response. An unknown category is reported as an empty one -- `?category=999`
answers `response_code` 1, for no results, not 2, for a bad parameter -- so a
typo in a category id is indistinguishable from a category that happens to be
empty. The strings are HTML-encoded inside JSON: a question arrives as `Kraft
&quot;Cheez Whiz&quot; ... doesn&#039;t`, entities in a format that needs none,
so every string wants an HTML decode after it has already been JSON-decoded.
`incorrect_answers` is three entries on a multiple-choice question and one on a
boolean, with nothing but the sibling `type` to say which. And `encode=base64`
encodes the enum values too, so `type === "multiple"` is false until even the
fields a caller would never think to decode have been decoded.

Ten are Open Brewery DB's, where `brewery_type` says `closed`. The field naming
what kind of brewery this is also answers whether it is one: its fourteen values
include `micro`, `brewpub` and `regional`, which are kinds, alongside `closed`,
which is a status, `planning`, which is a stage, and `location`, which is
neither. 642 are typed `closed` and 639 `planning` -- 1,281 of 11,848 -- so
filtering `by_type=micro` silently drops every closed micro, and showing all
types puts businesses that do not exist beside businesses that do.

**And two pairs of fields are the same field twice.** `state` and
`state_province` carry the same string on every record, and so do `street` and
`address_1` -- not similar, identical, `"Oklahoma"` and `"Oklahoma"`. Two are
current and two are kept from an older shape, and nothing says which pair is
which. Asking for one brewery too many is a redirect: `?per_page=201` answers
302 to the homepage rather than 400 or a clamp, so a client that follows
redirects gets the landing page and a 200. A brewery that does not exist answers
404 as an HTML page. The documented `autocomplete` path is a 301 to `search`.
`by_state` has 214 keys, most of which are not states -- ACT is Australian,
Argyll Scottish, Auckland a New Zealand region, Aveiro a Portuguese district,
all sorted in among Alabama and Arizona. And `phone` is two formats in one
column: `"+49 9261 628000"` and `"4058160490"`.

Eight are PoetryDB's, where `status` is the number `404` and the string
`"405"`. The same key, in the same API, holding a JSON number on one failure and
a JSON string on the other: `body.status === 404` catches the first and misses
the second, `body.status === "405"` does the reverse, and `body.status >= 400`
is true for one and false for the other, because `"405"` is not a number.

**And neither of them is an HTTP failure.** Both come back 200, with
`application/json`, so `response.ok` is true and every retry, error boundary and
status check sees a success -- the only thing saying otherwise is that field,
whose type changes between the two ways of saying it. `linecount` is the string
`"14"`, so the API has a number that is a string beside a string that is a
number. `poemcount` is a valid input field and an invalid output field: one
refusal's list of allowed words contains it and the other's does not, and both
are 405. A path segment that is not a field is reported as Method Not Allowed on
a GET that is the only method the endpoint has. And `/author` answers
`{"authors": [...]}` while `/author/{name}` answers a bare array, so two paths
one segment apart cannot share a parser.

Ten are Datamuse's, where a definition is a part of speech, a tab character
and a sentence, in one string. `"defs": ["n\tA large amount. "]` packs two
fields into a JSON string with a control character between them, inside a format
that has arrays and objects for exactly this -- and every entry ends in a space,
so a client that splits on the tab and renders index 1 shows trailing whitespace
it never asked for.

**And `tags` is three vocabularies in one array of strings.** `["n", "v",
"pron:S L UW1 "]` holds bare part-of-speech codes beside a colon-prefixed
key-value pair, so `tags.includes("n")` works and getting the pronunciation
means scanning the same array for a prefix and splitting on a colon -- and the
pronunciation carries a trailing space too. Three queries on one endpoint answer
three different field sets, and `md=` explains only one of the differences:
`?ml=ocean` gets tags and no syllable count, `?rel_rhy=blue` gets a syllable
count nobody asked for and no tags at all, and adding `md=dpsr` to the rhyme
query gets all five. An unknown parameter answers `[]` with 200, so a typo in a
parameter name is indistinguishable from a query that matched nothing. And
`score` is an integer with no scale anywhere -- 40,041,792 for a synonym and
10,051 for a spelling suggestion, three orders of magnitude apart on one host,
with nothing saying what either means.

Eleven are the Art Institute of Chicago's, where one field of ninety-eight is
licensed differently from the rest and the response says so in prose. `info.
license_text` reads "The `description` field in this response is licensed under
a Creative Commons Attribution 4.0 Generic License (CC-By) ... All other data in
this response is licensed under a Creative Commons Zero (CC0) 1.0 designation",
and that English sentence in a sibling key is the only statement of which field
is which. Nothing structural marks it: `description` sits in `data` beside the
other ninety-seven and looks the same, so a client that stores the response and
republishes it under one licence is wrong about exactly one column.

**And an image URL is three pieces from two levels of the response.** The
artwork carries `image_id` and no URL at all; the base lives in a top-level
`config` key beside `data` rather than in it; and the rest of the path --
`/full/843,/0/default.jpg` -- is in neither. `date_display` and `date_end`
disagree, so Seurat's Grande Jatte reads "1884-86, border added 1888-89" while
the number a query sorts on is 1886. The same coordinate is in the response
three times and the third is rounded: `latitude` and `longitude` carry thirteen
decimal places and `latlon` joins them as a string with twelve.
`colorfulness` is 0 on a pointillist painting made of coloured dots and 43.0038
on a Franz Marc, so nothing distinguishes "not computed" from "computed, and
grey". And `has_not_been_viewed_much` is a negated boolean, which makes the
ordinary question a double negative.

Fifteen are Postcodes.io's, where four of the fields are named after
organisations that no longer exist. `nhs_ha` is a Health Authority, abolished in
2002. `primary_care_trust` is a PCT, abolished in 2013. `nuts` is an EU region
the UK left in 2020. And `ccg` is a Clinical Commissioning Group, abolished in
2022 -- sitting next to `icb`, the body that replaced it, holding almost the
same string, with nothing marking either as historical.

**And in Wales the two abolished English fields hold the name of a Welsh body.**
CF10 3NP answers "Cardiff and Vale University Health Board" for both `nhs_ha`
and `primary_care_trust` -- two fields named after two different defunct English
structures, filled with a body that is neither -- while `icb`, which is an
English structure, says "Wales". The unversioned `lsoa` is an alias, and only a
postcode whose boundaries moved can say which census it follows: in Westminster
`lsoa`, `lsoa11` and `lsoa21` all agree, and in Cardiff they do not, which is
what shows `lsoa` tracking 2021. A null field's code is a row of nines --
`admin_county` is null and `codes.admin_county` is `"E99999999"` -- and Cardiff
carries the same sentinel under `codes.icb` while `icb` itself says "Wales".
`ccg_id` is a different shape on every record: `"W2U3Z"` on the English one,
`"7A4"` on the Welsh, against a letter and eight digits everywhere else. And
`date_of_introduction` is `"198001"`, six characters with no separator.

Ten are the Met Museum's, where `artistGender` is `"Female"` or
the empty string. There is no `"Male"`. A man is recorded as the absence of a
value -- and the absence of a value is what twenty other fields on the same
record use to mean "this does not apply". `culture`, `period`, `dynasty`,
`reign`, `city`, `country`, `locus` and `excavation` are all `""` on the van
Gogh, so `artistGender`'s empty string sits in a crowd of empty strings that
mean something else entirely. Anything counting women in the collection can do
it; anything counting men cannot tell them from the unattributed, the anonymous,
and the records nobody has filled in.

**And absence is spelled `""` rather than `null` on twenty-one of fifty-seven
fields**, so `if (object.culture)` is falsy, `object.culture === null` is false,
and `"culture" in object` is true -- the three ordinary ways of asking whether a
field is set disagree with each other. Three year fields carry two types:
`accessionYear` is the string `"1993"`, `artistBeginDate` the string `"1853"`,
and `objectBeginDate` the number `1889`. The same `dimensions` field uses U+00D7
on one record and the letter `x` on another. A missing image is `""` where a URL
goes, so assigning it to an `img` src asks the page for its own address. And
`GalleryNumber` is the only capitalised field name on a record that also shouts
`artistWikidata_URL` and `AAT_URL` in the middle: fifty-seven fields, three
naming conventions.

Eleven are Rick and Morty's, where asking for a character by a name
instead of a number is a 500. `GET /api/character/abc` answers
`{"error": "Hey! you must provide an id"}` with a server-fault status, so the
retry wrapper retries, the alert fires and the dashboard records an outage --
for a request that was never going to work. The 404 beside it, for a number that
is simply not a character, is the correct status, so the API can produce the
right answer and does not for the wrong input.

**And `"unknown"` is the only lowercase value in the enum.** `status` is
`"Alive"`, `"Dead"` or `"unknown"`, so the one value meaning "we have no answer"
is also the one that breaks a switch written from the other two. One object says
"no value" two ways at once: on a character whose origin is not known,
`origin.name` is the string `"unknown"` and `origin.url` is the empty string --
a word in one field and an absence in its neighbour, inside the same two-key
object. Related records are URLs rather than identifiers, so joining anything
means parsing an integer back out of a string. And three failures carry three
sentences and two statuses, only one of which ends in a full stop.

Seven are Deck of Cards', where a failed draw empties the deck and
hands you the cards anyway. Asking a fifty-two card deck for 999 cards does not
refuse: it draws all fifty-two, sets `remaining` to 0, says `"success": false`,
returns them, and answers HTTP 200. There is no safe way to read that. A client
that checks `success` and discards the body has thrown away fifty-two real cards
and left the deck empty; a client that ignores `success` has fifty-two it did not
ask for the right number of. Either way the deck is spent, and the request that
spent it is the one the API called a failure.

**And the ten is a zero.** Card codes are two characters, so the ten of spades is
`"0S"` -- while its `value` field says the string `"10"`. Anything parsing a code
by taking its first character reads a nought where a ten should be. `value` is a
string for numbers and a word for faces, so sorting a hand by it puts `"10"`
before `"6"` and both before `"ACE"`. Every image is in the response twice, as
`image` and `images.png`, character for character. And a count that is not a
number answers 500 with Django's default error page in HTML, from a JSON API --
the framework's page, not anything this API wrote.

Nine are TheMealDB's, where `meals` is an array, or `null`, or a string. A
lookup that finds something answers `{"meals": [ {...} ]}`; one that finds
nothing answers `{"meals": null}`; and one with no identifier at all answers
`{"meals": "Invalid ID"}`. One key, three JSON types, and HTTP 200 on all three.
`body.meals.length` works on the first, throws on the second, and returns 10 on
the third -- because `"Invalid ID"` is ten characters long -- so a client
checking `if (body.meals)` treats the error message as a truthy result and
iterates it character by character.

**And forty of the record's fifty-four fields are one array flattened into
columns.** `strIngredient1` through `strIngredient20` and `strMeasure1` through
`strMeasure20`, so a nine-ingredient recipe carries eleven empty slots -- and
they are not empty the same way. Slots 10 to 15 are the empty string and slots
16 to 20 are `null`, in one numbered series, split at an index nothing explains.
`strTags` is a comma-joined string where an array belongs. `strInstructions`
carries CRLF inside a JSON string. And `strArea` is `"Japanese"` beside
`strCountry`'s `"Japan"`: the adjective and the noun, in two fields, on one
record.

Seven are the World Bank's, where the success is a two-element array and the
failure is a one-element array. A country lookup answers
`[{"page": 1, "pages": 1, "per_page": "50", "total": 1}, [ {...} ]]` and a bad
code answers `[{"message": [ {...} ]}]`, both with HTTP 200. `body[1]` is the
collection on success and `undefined` on failure, and nothing about the status,
the content type or the outer shape says which arrived except how many elements
the array has -- so `body[1].map(...)` throws a TypeError that names nothing
useful.

**And `per_page` is a string where its three neighbours are numbers.** The paging
object is four fields describing one page, three of them integers and one
quoted. `iso2Code` and `iso2code` are both on the record, capital-C at the top
and lowercase one level down. Three of those codes are not ISO codes at all --
North America is `"XU"`, High income is `"XD"`, Not classified is `"XX"`, all
from the range the standard reserves for private use, in a field named after it.
A missing value is an object of three empty strings rather than `null`, so
`if (country.adminregion)` is true and reading `.value` off it gives `""`. And
the coordinates are quoted on a record whose page numbers are not.

Eight are disease.sh's, where six fields end in `PerOneMillion` and three of
them are integers. `casesPerOneMillion`, `deathsPerOneMillion` and
`testsPerOneMillion` are whole numbers; `activePerOneMillion`,
`recoveredPerOneMillion` and `criticalPerOneMillion` carry two decimal places.
One suffix, one obvious meaning, and two JSON types split by which statistic it
happens to be -- so a client that types the family as an integer truncates three
of them, and nothing in the response distinguishes the two groups.

**And `countryInfo` carries `_id`.** A field named the way a database names its
primary key, underscore and all, on a public API -- holding 124, which is the UN
numeric country code and not an internal row id at all. `updated` is a
JavaScript millisecond timestamp as a bare integer, and it changes on every
request while the figures beside it do not. The coordinates are whole degrees on
a record whose derived statistics carry two decimals, so the least precise
numbers here are the ones describing a place -- and one of them is called `long`
rather than `longitude`. A country that does not exist answers JSON and a path
that does not exist answers an Express HTML page, both with 404.

Five are Dog CEO's, where `message` is a URL, or an array of strings, or an
error sentence. The random-image endpoint answers a URL under it, the sub-breed
list answers an array of names under it, and a breed that does not exist answers
a sentence under it -- and two of the three are strings, so a client cannot tell
the image from the error by type. The only way to know which arrived is to read
`status` first, or to notice that one of them starts with `https`. Code that
assigns `message` straight to an `img` src renders a broken image whose alt text
is the error.

**And the object's shape depends on whether it worked.** A success carries
`message` and `status`; a failure carries `message`, `status` and `code` -- so
the status line is restated twice, once as a word and once as a number, and only
on the way that failed. The sub-breeds are bare names rather than paths or
objects, so a client holding the list has to remember which breed it asked about
to build the next URL.

Six are Advice Slip's, where every response carries two `Cache-Control` headers
and they contradict. One says `max-age=3600`; the other says
`max-age=600, private, must-revalidate`. Both are on the wire, on the same 200,
on every request -- and a recipient that follows RFC 9110 joins repeated field
lines with commas, producing one directive list with two `max-age`s in it, which
is not a thing the grammar has an answer for. `fetch`'s `headers.get` returns
exactly that string, so a client reading its own cache policy has to decide
which half to believe.

**And three outcomes carry three different top-level keys.** A slip by id
answers `{"slip": {...}}`, a search answers `{"slips": [...]}`, and a failure
answers `{"message": {...}}` -- singular, plural, and neither, so there is no key
a client can read first to find out what happened. `total_results` is the string
`"1"` beside an `id` that is the number 1. The same slip has two field sets:
`/advice/1` answers `{id, advice}` and the search answers `{id, advice, date}`
for that slip, so a date exists that one route never shows. And a miss and an
empty search are both 200, distinguished only by one word -- `"error"` against
`"notice"` -- in a field called `type` inside an object called `message`.

Six are SWAPI's, where every number is a string and one of them has a unit in
it. `rotation_period` is `"23"`, `diameter` is `"10465"`, `population` is
`"200000"` -- and `gravity` is `"1 standard"`. Sorting planets by diameter puts
`"10465"` before `"9000"` because `"1"` comes before `"9"`, and `parseFloat` on
gravity gives 1 while the string it came from means one standard gravity, so the
field that looks most like a number is the one least safely read as one.

**And the absence of data is spelled two ways in one record.** The planet at
`/planets/28` is named `"unknown"`, and its six numeric fields split
three-three: `rotation_period`, `orbital_period` and `diameter` are `"0"`, while
`surface_water`, `population` and `gravity` are `"unknown"`. The same missing
information, two encodings, in the same object -- so a client treating `"0"` as
a real zero charts a planet with no diameter, and one treating `"unknown"` as
the sentinel misses the three not spelled that way. Related records are URLs
rather than identifiers, and a planet that does not exist answers the static
site's own HTML shell with a 200.

Nine are Chuck Norris's, where the success and the failure carry timestamps in
two different grammars. A joke's `created_at` is `"2020-01-05 13:42:19.324003"`
-- a space instead of a T, six digits of fraction, no timezone -- and an error's
`timestamp` is `"2026-08-30T06:21:50.575Z"`, proper ISO 8601. Same API, same
request, and which format arrives depends on whether it worked, so a client that
logs when something happened gets a date on its errors and `Invalid Date` on its
results.

**And the two identifier alphabets do not match either.** Most jokes carry
base64url -- `"xMUTjH93TyWHGR6MVf_vpA"`, mixed case with hyphens and underscores
-- and some carry `"rhvv9w42qka07narzt47da"`, lowercase letters and digits only.
Both twenty-two characters, two character sets, one collection. `updated_at` is
never different from `created_at`, so a field that exists to say when something
changed has never said it. A query parameter that is not a category answers 404
naming `"path": "/jokes/random"` -- the path that does exist and did answer. And
an empty search is `{"total": 0, "result": []}`, with the array called `result`,
singular, beside a total that is not.

Eight are Open Notify's, where **there is no HTTPS at all**. Port 443 does not
answer: `http://api.open-notify.org/iss-now.json` returns 200 and the same URL
over TLS times out, so no page served over HTTPS can call this API from a
browser. The only published client hardcodes `http://` as its default domain,
because there is no other one to reach.

**And every failure declares JSON and sends none.** A 404 and a 405 both answer
`Content-Type: application/json; charset=UTF-8` with a body of zero bytes, so a
client that trusts the content type calls `.json()` on an empty string and
throws -- having been told by the response itself that a document was coming.
The endpoint the documentation still describes fails differently again:
`/iss-pass.json` is gone and answers nginx's own HTML 404, `text/html`, 169
bytes, naming the version it runs. Two 404s in two content types on one API, and
which one arrives depends on how far in the request got. The position is two
strings beside a number -- `iss_position.longitude` and `iss_position.latitude`
are quoted while the `timestamp` next to them is an integer -- `message` reads
`"success"` on a 200, and `number` is `people.length` sent alongside the array,
on an endpoint where `craft` is not always `"ISS"` because Tiangong appears in
it too.

Seven are Agify's, where **three APIs share one budget**. api.agify.io,
api.genderize.io and api.nationalize.io are presented as separate services with
separate documentation, and interleaved calls decrement a single counter: 15,
14, 13, 12, 11 across the three hosts, one allowance of twenty-five a day. So
asking all three about a name -- the obvious thing to build -- costs three of
them, and a page that does it for eight visitors is finished until midnight.

**And `X-Rate-Limit-Reset` is seconds until midnight UTC, not a timestamp.** It
read `60485` while there were `60501` seconds left in the UTC day, so a client
passing it to `new Date(reset * 1000)` lands in January 1971. A name nobody has
heard of is a 200 with `"age": null`, not a 404. Present-but-empty and absent
are different: `?name=` answers 200, and omitting the parameter answers 422.
And the 422 is the only response carrying no rate-limit headers at all, on a
content type without the `charset=utf-8` every success declares -- the response
that says a client is doing something wrong is the one that will not say how
much budget is left.

Ten are Bible API's, where **the same verse comes back with different
whitespace depending on the translation, and not the same kind of whitespace.**
The World English Bible pads `John 3:16` with a leading newline and two trailing
ones; the King James pads it with ten leading spaces and one trailing newline.
One query parameter apart, on the same path, for the same verse -- so a client
that renders the text raw gets a blank line for one and an indent for the other,
and a client comparing two translations compares their padding as well as their
words.

**And the text is sent twice, both copies untrimmed:** the top-level `text` and
`verses[0].text` are byte-identical, padding included. The two endpoints share
no keys at all -- a reference answers `{reference, verses, text, translation_id,
translation_name, translation_note}` and `/data/web/random` answers
`{translation, random_verse}`, not one field name in common. The same concept
carries two names, `book_name` on one and `book` on the other, and the string
`"Public Domain"` arrives as `translation_note` in one place and `license` in
the other. There are two 404s and one of them is nginx's own HTML: an unknown
book is the application's `{"error": "not found"}`, and a chapter past the end
of a book never reaches the application at all.

Six are CoinGecko's, where **the map is keyed by currency and one of the keys
is not a currency.** `/simple/price` answers `{"bitcoin": {"usd": 78148, "eur":
67468, "last_updated_at": 1788075830}}`, so the field that says when the prices
were taken sits at the same depth, under the same kind of key, as the prices
themselves. A client that iterates the object to list them -- the only way to
read a map whose keys are whatever the caller asked for -- gets a currency
called `last_updated_at` worth 1,788,075,830.

**And anything the API does not recognise is dropped in silence.** An unknown
coin is not an error and not a null: `?ids=bitcoin,zzzznotacoin` answers
`{"bitcoin": {...}}` with a 200 and nothing to say one was rejected, and asking
only for unknown coins answers `{}`. One price is an integer and the next a
float in the same response, because the value is truncated by magnitude rather
than typed by field. There are two error shapes -- a flat `{"error": "coin not
found"}` and a nested `{"status": {"error_code": 429, ...}}` -- so one client
needs two parsers. And `/ping` answers `{"gecko_says": "(V3) To the Moon!"}`,
where the only statement of which API version replied is the `(V3)` in the
middle of a sentence.

Nine are iTunes', where **the JSON is served as JavaScript.**
Every response carries `Content-Type: text/javascript; charset=utf-8`, success
and failure alike. A browser sniffs its way to the right parser and
`fetch(...).json()` works, so nobody notices -- until the request goes through
anything that checks the content type first, which then refuses a body it could
have read, or hands it to a JavaScript evaluator because that is what the header
asked for.

**And there is no such thing as no results.** `?term=zzzzqqqxyznothing` answers
`resultCount: 60`, every one an audiobook by "ZZZ Sleep": the search decomposes
the term and finds something for the part of it that is pronounceable, so a page
reporting sixty results for a word nobody typed on purpose is showing sleep
recordings. A band is not the first result for its own name either --
`?term=radiohead` answers a feature film, `kind: "feature-movie"`, and Radiohead
appears only once the caller asks for `entity=musicArtist`. The array is
heterogeneous, a track and an artist sharing almost no keys and told apart only
by `wrapperType`; `amgArtistId` is a foreign key into All Media Guide, renamed
in 2013; and an id nobody has answers the same `{"resultCount": 0, "results":
[]}` as a search with no term at all, so a rejected request and an empty one are
one answer.

Six are EPSS's, where **one of the keys has a hyphen in it.** The envelope is
`{"status": "OK", "status-code": 200, ...}`, and `status-code` cannot be reached
with a dot in any language whose dot operator takes an identifier. In
JavaScript, `response["status-code"]` is 200 and `response.status-code` is
`response.status` minus a variable called `code` -- `NaN`, at runtime, with no
error and no warning. The field most likely to be read by code checking whether
the request worked is the one most likely to be read wrong.

**And the status is stated three times:** the HTTP status line, `status`, and
`status-code`, three sources for one fact in one response. The probability is a
string to nine decimal places -- `"0.999990000"`, quoted and padded, wrong for
every numeric comparison that does not parse it first. A CVE that is not a CVE
is a 200 with an empty array, exactly as an unknown but well-formed one is,
while omitting the parameter entirely is a 404: the malformed request succeeds
and the incomplete one fails. And failing makes the data private -- a 200 says
`"access": "public"` and a 404 says `"access": "private"`, so a field describing
what may be done with the data changes because the request did not work.

Eight are Chess.com's, where **the country is a URL.** A player's `country` is
`"https://api.chess.com/pub/country/US"` -- not `"US"` and not `"United
States"` -- so a client rendering `player.country` puts an API address on the
page, and getting the two letters costs a second request to a resource that is
`{"@id", "code", "name"}` and nothing else.

**And the player has two URLs that are not the same URL:** `"@id"` is
`https://api.chess.com/pub/player/hikaru` and `"url"` is
`https://www.chess.com/member/Hikaru`. Two hosts, two paths, two
capitalisations of one username, and nothing in either name saying which is
which -- and `@id` begins with a sigil, so it needs a subscript in most
languages. `twitch_url` is sent twice, once at the top level and once as
`streaming_platforms[0].channel_url`, byte-identical, so the field and the
structure meant to replace it both ship. `status` is a subscription tier,
`"premium"`, in a response whose HTTP status is 200. And the failure code is
zero: a player who does not exist answers `{"code": 0, ...}`, and so does a path
that does not exist, whose message is `Data provider not found for key
"/pub/player/hikaru/notanendpoint".` -- internal vocabulary, with the path
echoed back inside escaped quotes.

Eight are ip-api's, where **TLS is a paid feature and the message saying so is
delivered over TLS.** `https://ip-api.com/json/1.1.1.1` answers 403 with
`{"status": "fail", "message": "SSL unavailable for this endpoint, order a key
at https://members.ip-api.com/"}`. The connection is accepted, the certificate
is valid, the handshake completes, and the body says SSL is unavailable --
everything TLS is for worked, and what is missing is the entitlement.

**And every other failure is a 200.** A malformed address and a string that is
not an address at all both answer HTTP 200 with `{"status": "fail", "message":
"invalid query"}`, so the status line is not a signal and a client checking
`response.ok` gets true for both. A format the API does not recognise is the
marketing website: `/nosuchformat/` answers 200 with eleven kilobytes of HTML,
meta keywords and all. The rate limit is `X-Rl` and `X-Ttl`, two abbreviations
and neither standard. And the autonomous system is a number and a name in one
string, `"AS13335 Cloudflare, Inc."`, comma included, beside an `isp` with no
full stop and an `org` that is a different organisation again.

Nine are Where the ISS at's, where **the same field is a number on one endpoint
and a string on another.** `/v1/satellites/25544` sends `"latitude":
26.268487855719` and `/v1/coordinates/37.79,-122.39` sends `"latitude":
"37.795517"` -- one API, one field name, two types. `id` does it too: the
satellite carries `25544` and its own two-line element set carries `"25544"`.

**And the units are declared after the numbers they govern.** `altitude` and
`velocity` mean different things depending on a query parameter, and the only
thing that says which is a field at the end of the object: 419.69504466844
kilometres by default and 260.21728290923 miles with `?units=miles`, so a
parser reading fields in order has the numbers before it has the scale. Beside
them, `daynum` is a Julian date -- 2461282.8361574, counting from 4713 BC --
sharing an object with a Unix `timestamp` of 1788077044 and nothing to tell
them apart. The failure puts the status last, `{"error": "satellite not found",
"status": 404}`, and a path with no controller behind it says exactly that:
`"Invalid controller specified (nosuchthing)"`.

Eight are TfL's, where **the severity scale runs backwards.** `statusSeverity`
10 is `"Good Service"` and 6 is `"Severe Delays"`, with `"Part Suspended"` at 3:
the bigger number is the better railway. Every other severity scale a client has
met counts upward into trouble -- syslog, HTTP, CVSS, log levels -- so sorting
descending to put the worst lines first puts the healthiest lines first instead,
silently, on data that looks sorted.

**And every object carries its .NET type and assembly.** `"$type":
"Tfl.Api.Presentation.Entities.Line, Tfl.Api.Presentation.Entities"` --
namespace, class, comma, assembly, on every object at every depth, beginning
with a sigil so it needs a subscript in most languages. A `LineStatus` carries
`"created": "0001-01-01T00:00:00"`, which is .NET's `DateTime.MinValue`
serialised straight onto the wire because nothing set it, beside a `Line` whose
own `created` is a real 2026 date. A failure says what went wrong four times --
`exceptionType` as a .NET class name, the status as a number, the status as a
string, and a sentence -- with a timestamp to seven fractional digits. And a
path with nothing behind it answers `"Resource not found:
http://api:8001/NoSuchEndpoint"`, a service name and a port from inside the
deployment.

Seven are openFDA's, where **every value begins with its own field name.**
`"purpose"` is `["Purpose Pain reliever"]`, `"do_not_use"` is `["Do not use if
you are allergic to aspirin ..."]`, `"questions"` is `["Questions or comments?
Call 1-877-753-3935 ..."]`. These are scanned drug labels and the section
heading came along with the section, so the field name is both the key and the
first words of the value, and rendering "Purpose: Purpose Pain reliever" is the
natural thing to build.

**And two of them begin with a different field's name.**
`indications_and_usage` starts "Uses for the temporary relief" and
`storage_and_handling` starts "Other information store between", so stripping
the prefix is not a fix: the prefix is the heading printed on the box and the
key is what a database called it. Every field is an array of one, identifiers
included. Finding nothing is a 404 -- `{"code": "NOT_FOUND", "message": "No
matches found!"}`, with an exclamation mark -- and a search naming a field that
does not exist answers exactly that, byte for byte. A path that does not exist
is Express's own page in `text/html`. And the temperature is written with
U+00BA MASCULINE ORDINAL INDICATOR rather than U+00B0 DEGREE SIGN: a character
that looks right at a glance and matches no search for degrees.

Eight are OpenAlex's, where **the error message is the schema.** A filter that
does not exist answers 400 with `"nosuchfilter is not a valid field. Valid
fields are underscore or hyphenated versions of: abstract.search,
abstract.search.exact, ... updated_date, version"` -- 4,779 characters and 207
field names, comma-separated, inside one sentence in one string. The failure is
larger than most successful responses, the entire queryable surface of the API
arrives as prose, and finding out which field was wrong means splitting an
English sentence on a colon and then on commas.

**And a free API quotes a price:** every response carries `"cost_usd": 0.0001`
in its meta, a tenth of a hundredth of a cent for a request nobody is billed
for, beside an `x_query` that leaks the internal query language. Three failure
shapes share one host: Flask's default HTML page for a missing work, a lone
`message` for a missing entity, and `{error, message}` for a bad filter.
`title` and `display_name` are byte-identical. And the identifiers are URLs, so
getting the bare W-number or the bare DOI means string-splitting an address.

Seven are the MBTA's, where **`type` is a number and a string in the same
object.** `data.attributes.type` is `1` -- the GTFS route type, where 1 means
subway -- and `data.type` is `"route"`, the JSON:API resource type. Two levels
apart, one key name, and a client reading `type` without saying which one gets
whichever its path happened to reach. A third lives in every relationship, where
`relationships.agency.data.type` is `"agency"`.

**And the human-readable text moves keys between failures.** A resource that
does not exist answers `{"code": "not_found", "status": "404", "title":
"Resource Not Found"}` -- title, no detail -- and a filter that does not exist
answers `{"code": "bad_request", "status": "400", "detail": "Unsupported
filter(s): nosuchfilter"}` -- detail, no title. So a client rendering
`errors[0].detail` shows nothing when a route is missing. `status` is a string
in the body beside a number on the wire; the same error body arrives under
`application/vnd.api+json` from one path and `application/json` from another;
colours are six hex digits with no hash, so `style="color: DA291C"` renders
nothing; and `direction_names` and `direction_destinations` are two arrays
joined only by index.

Seven are Carbon Intensity's, where **a path that does not exist answers 200 and
the body says 400.** `GET /nosuchpath` returns HTTP 200 carrying `{"error":
{"code": "400 Bad Request", "message": "Please enter a valid path e.g.
/intensity/"}}` -- so `response.ok` is true for a URL the API does not have, and
the code is not a code but a status line as a string, with the number and the
reason phrase in one field. A date it cannot parse, meanwhile, answers a real
400 with the same body.

**And the same fuels are named two ways on two endpoints.** `/generation` calls
them `"gas"`, `"imports"`, `"nuclear"` -- lowercase, one word -- and
`/intensity/factors` calls them `"Gas (Combined Cycle)"`, `"Gas (Open Cycle)"`,
`"Dutch Imports"`, `"French Imports"`, `"Irish Imports"`: capitalised, spaced,
parenthesised, and split finer. Neither vocabulary is a key of the other, so
joining a fuel's share to its emissions factor needs a mapping the API does not
publish. The timestamps have no seconds -- `"2026-08-30T10:00Z"`, which is the
form its own error message asks for -- and every half hour carries a `forecast`,
an `actual` and an `index`, two numbers and a category derived from one of them
by a threshold nothing states.

Eight are UK Police's, where **`url` means three different things.** A force
sends `"url": "http://www.leics.police.uk/"`, a crime category sends `"url":
"all-crime"`, and the telephone engagement method sends `"url": ""`. A web
address, a slug, and nothing -- so a client rendering `url` as a link produces a
working link, a broken relative one, and an empty href, depending on which
endpoint it came from.

**And every link the force owns is plain HTTP.** Its website, its Facebook page,
its Twitter and its YouTube are all `http://`, in a payload delivered over TLS;
the one exception is the RSS feed, which is `https://`, so the API is not
consistently either. Three content types answer three kinds of request, and only
one is JSON: a force that does not exist is `text/plain` carrying the nine bytes
`Not Found`, and a path that does not exist is `text/html` carrying an HTML
fragment with no doctype around it. `type` and `title` are the same word on every
engagement method. And `stop-and-search` is a hyphenated JSON key, which a dot
cannot reach.

Eight are FHRS', the Food Hygiene Rating Scheme, where **the same digit means best
in one field and not in the other.** `"RatingValue": "5"` is the top of a scale
that counts up, and `"scores": {"ConfidenceInManagement": 0}` is the top of one
that counts points against you -- so zero is perfect there and five is not, in
the same object, and the first is a string while the second is a number.

**And absence is spelled three ways in one record:** `Phone` is `""`, `Distance`
is `null`, and `AddressLine2` is `""` while `AddressLine3` holds the town, so an
empty line is a gap in the middle of an address rather than the end of it. A
forgotten `x-api-version` header is reported as a path that does not exist --
`"The API 'Authorities' doesn't exist"` -- on a path that works the moment the
header is sent. That failure is a bare JSON string rather than an object, so
`JSON.parse` succeeds and `body.Message` is `undefined` rather than throwing,
and the name inside it is the internal handler: the version and the path joined
by a dot. The coordinates are strings where the scores are numbers, and
`RatingKey` packs the scheme, the rating and the locale into `"fhrs_5_en-gb"`.

Six are BoC Valet's, where **the response documents its own
key names inside itself.** `seriesDetail.FXUSDCAD.dimension.key` is `"d"` and
`.name` is `"Date"`, and the observations carry the date under `"d"` -- so a
client that wants it has to read the dimension first and subscript by whatever
it finds, or hardcode one letter and hope.

**And the number is two levels down, under a key the caller supplied.** A rate
is `obs["FXUSDCAD"].v`: the series name is a JSON key rather than a value, and
what sits under it is an object with one field, holding a string. The same key
differs by one letter between endpoints -- `/observations` returns
`seriesDetail`, a map keyed by series name, and `/series/{name}` returns
`seriesDetails`, flat, with the name as a field. The terms of use are the first
key of every successful body. And a bad parameter answers with a link to a
generated OpenAPI operation id, underscores and all, in a URL fragment.

Eight are FBI Wanted's, where **the reward is zero and the sentence beside it
says twenty-five thousand dollars.** `reward_min` and `reward_max` are both `0`,
and `reward_text` reads "The FBI is offering a reward of up to $25,000 for
information leading to the identification, arrest, and conviction ...". The
numbers are the machine-readable reward and the prose is the real one, so
anything sorting by `reward_max` puts a twenty-five-thousand-dollar case last
and anything filtering on `reward_max > 0` does not show it at all.

**And thirty-three of the record's fifty-four fields are null.** Every physical
descriptor -- `sex`, `hair`, `eyes`, `build`, `complexion`, `weight`, `race`,
`scars_and_marks`, `place_of_birth` -- because the subject has not been
identified, on a record that still carries a `warning_message` reading "SHOULD
BE CONSIDERED ARMED AND DANGEROUS". Absence is spelled two ways among fields of
the same kind: eight lists are `null` and one, `coordinates`, is `[]`. `detail`
is an array on one failure and a string on the other. Two timestamps in one
record disagree about time zones. And `description` packs three facts into one
string with carriage returns between them.

Seven are Zenodo's, where **the same identifier appears four times in four
forms.** `id` is the number `22173677`, `recid` is `"22173677"` as a string,
`doi` is `"10.5281/zenodo.22173677"`, and `doi_url` is that with a resolver in
front. A fifth, `conceptrecid`, is `"22173676"` -- one less, as a string -- and
`conceptdoi` is that one prefixed, so eight of the record's twenty top-level
keys are the same two numbers written out differently.

**And the record says it is finished three times, in three vocabularies:**
`"status": "published"`, `"state": "done"`, `"submitted": true`, with nothing
saying which is authoritative. `modified` and `updated` are byte-identical to
the microsecond. The array is called `hits` and lives inside an object called
`hits`, so reaching the records is `response.hits.hits`. And two failures use
different words for the same status -- a record that does not exist is about a
persistent identifier, and a path that does not exist is the web framework
apologising in case you typed the URL by hand.

Seven are UniProt's, where **the same acronym is capitalised two ways in two
top-level keys of one object.** `uniProtkbId` has a lowercase `kb` and
`uniProtKBCrossReferences` an uppercase `KB`, four fields apart -- so anything
generating types from the response gets two spellings of one name, and anything
guessing the other finds nothing.

**And the failure echoes your URL back with the scheme downgraded.** A request
to `https://rest.uniprot.org/uniprotkb/NOSUCHACC` answers `{"url":
"http://rest.uniprot.org/uniprotkb/NOSUCHACC", "messages": [...]}` -- the
request was over TLS and the URL in the reply is not, so a client that logs it
or retries it drops to plain HTTP without being told. That failure is a 400
rather than a 404, while a path that does not exist is nginx's own HTML page
naming the version it runs. `proteinExistence` packs a number and its label into
one string, `"1: Evidence at protein level"`. And `entryType` is a sentence with
a parenthetical that names the database twice.

Seven are Spaceflight News's, where **one record links the same site twice, once
securely and once not.** `url` is `https://spaceflightnow.com/...` and
`image_url` is `http://spaceflightnow.com/...`, on the same host, in the same
object -- so a page served over TLS that renders the picture gets a
mixed-content block and shows nothing, while the link beside it works.

**And the API's own 404 page carries a third-party analytics beacon**, a script
tag pointing at `static.cloudflareinsights.com/beacon.min.js` inside an error
response. Two timestamps in one record carry two precisions: `published_at` to
the second and `updated_at` to the microsecond. A `limit` that is not a number
is silently replaced with the default, so a client with a broken parameter never
learns. Absence is spelled two ways four fields apart, `socials: null` and
`events: []`. And the two failures do not agree on a format at all: an article
that does not exist answers JSON with the model name capitalised mid-sentence,
and a path that does not exist answers HTML.

Eight are ClinicalTrials.gov's, where **every failure is plain text and one of
them is nothing at all.** A malformed identifier answers 400 with `Parameter
`nctId` has incorrect format`; a parameter that will not convert answers 400
naming the integer width; and a path that does not exist answers 404 with no
body and no `Content-Type` header, so there is not even a wrong document to
parse. The two that do speak wrap their parameter names in Markdown backticks,
in responses declared as plain text, formatted for a renderer that is not there.

**And a date says whether it actually happened.** `startDateStruct` is
`{"date": "2019-01-24", "type": "ACTUAL"}` and `primaryCompletionDateStruct` is
`{"date": "2024-12-31", "type": "ESTIMATED"}` -- a fact and a guess, told apart
by one word a level down, in fields whose names end in a word about the code.
`overallStatus` is `"UNKNOWN"` beside a `lastKnownStatus` of `"RECRUITING"`, so
the study is both. `statusVerifiedDate` is `"2024-01"`, a year and a month with
no day, among dates that have all three. And `phases` is a list of one string,
and the string is `"NA"`.

Nine are DataCite's, where **the same key holds a list in one failure and an
object in another.** A missing DOI answers `{"errors": [{"status": "404",
"title": "..."}]}` and an unparseable query answers `{"errors": {"title":
"parse_exception ..."}}`. A client reading `errors[0].title` gets the sentence
on the first and undefined on the second, and one reading `errors.title` gets it
the other way round -- and the key is spelled the same, plural, in both. The
object one is 318 characters of the query parser's own Lucene grammar, newlines
and token names included, forwarded to whoever typed the URL.

**And the listing says there are 133,915,241 records and 10,000 pages of one.**
`totalPages` is not the total divided by the page size; it is ten thousand
records' worth, whatever the page size. A client paging until `page >
totalPages` reads 10,000 of 133 million and stops, told both numbers by the same
object. Asking for page 10,001 answers 200 with `"page": 10000`, a different
number from the one asked for, in the field named after it. A `page[size]` that
is not a number answers 200 with `"data": []` and `"totalPages": 0` beside that
unchanged total. Absence is spelled three ways in one record -- `container: {}`,
`formats: []`, `reason: null` -- the same year is sent twice in two types, and
the record's type is sent six times in six spellings, five saying dataset and
the sixth saying `misc`.

Every other provider needs an account, and a date nobody can reproduce is
worth less than an empty field that says so.

Cases run in process and need no credentials, so CI runs them on every push and
a Recipe edit that drifts from the provider fails there rather than in an
application months later.

Writing a Recipe is the format's stress test, and each new provider has found
something. Slack's failures arrive with HTTP 200 and its identifier lives in a
query parameter, so the router learned RPC shapes. HubSpot nests every business
attribute under `properties`. Notion refuses a request with no version header.
Plaid's error category is a separate field from its code. Intercom's paging
state is nested, which exposed declared constants silently overwriting computed
values. Discord's snowflakes are numeric strings long enough that minting small
integers would let a rounding bug through. Square's money is an object in minor
units. Cloudflare puts every payload under `result`, success or failure, so a
client branching on the HTTP status is checking the wrong thing. Xero answers a
request for one invoice with a list of one. PagerDuty authenticates with
`Token token=`, not `Bearer`. Sentry's entire error body is one string.

Datadog sends its errors as an array of bare strings, so a client reading
`.message` from each entry finds `undefined` on every one.

Miro nests coordinates under `position`, so code reading `item.x` gets nothing
and the item silently lands at the origin. Contentful keeps the identifier at
`sys.id`.

Webflow found a bug that had already shipped: every timestamp field was filled
in automatically, so a site that had never been published still carried a
`lastPublished`, and a Typeform response somebody abandoned still carried a
`submitted_at`. The emulator was claiming events that never happened, which no
test written against it could catch.

Mailgun found the worst one. Basic authentication only ever compared the
username, which is right for Twilio, where the account SID is the username, and
wrong for Mailgun, whose username is the constant `api` and whose key is the
password. A request with the correct username and a completely wrong key
returned 200, so the failure path a test most wants to exercise could not be
reached at all.

Dropbox found a limitation in the conformance checker rather than the emulator.
It names a field `.tag`, where the leading dot is part of the name, and the
checker split every path on dots, so there was no way to assert on it at all.
The emulator had been sending it correctly the whole time while every case that
mentioned it failed.

Trello is the clearest case for reproducing something rather than improving on
it. Its credentials travel in the query string, which is a bad idea, and its
failures are plain text rather than JSON. A fake that accepted a header would
hide the fact that the key ends up in access logs, and one that answered in
JSON would hide that calling `.json()` on a real error throws.

A pattern runs through the newer Recipes: the state that breaks an integration
is almost never a failure, it is a third state nobody branched on. A GitLab
merge request that is closed rather than merged. A CircleCI workflow waiting on
approval, which is neither a pass nor a failure. A Paddle subscription with a
cancellation scheduled but not yet applied. Each fixture carries one.

Of the thirty-five Recipes from Fivetran to Langfuse, twenty-five needed no
change to the format or the runtime. Writing a Recipe is now mostly research
rather than engineering: the work is reading the provider closely enough to
know which state nobody branches on, not teaching the format a new trick.

That window is named rather than described as "the last thirty-five", because
a sliding window is a number that goes stale without anybody editing it. It is
also the one figure here no test can hold: it is a claim about commits, and CI
clones without history. `git log --diff-filter=A` over `recipes/` and
`internal/` settles it for anyone who wants to check.

The suite has started finding faults in itself as well. Vonage reports a
successful send with the string `"0"` and Docusign counts with strings, and the
checker was comparing scalars by their rendered text, so `"0"` and `0` were
indistinguishable and both of those cases passed whichever the emulator sent.
Scalars now compare kind as well as value, which immediately turned up a real
inconsistency: the nested error style was sending PagerDuty's numeric codes as
text.

Not every provider fits, and one is left out deliberately rather than faked.
Temporal Cloud is gRPC rather than HTTP: this format describes HTTP surfaces,
and nothing built on it would be a Temporal client. No Recipe is better than a
misleading one.

That sentence used to name two others, and both of them ship now. Linear was
out for being GraphQL -- a single endpoint where the request body decides the
response shape -- until a route learned to match on the field a query names,
and Linear, Monday and ShipHero all shipped on it. Attio was out for the same
reason, recorded from an assessment that was simply wrong: Attio is REST and
publishes an OpenAPI description. Both entries sat unchallenged because
nothing rereads a list of things that cannot be done, and a reason that has
quietly expired reads exactly like one that still holds.

The earlier round found real bugs rather than confirming what the code did:
routes like Shopify's `/orders/{id}.json` matched nothing at all, so every
single-object route on Shopify and Twilio answered 404; errors came back in
Stripe's nested shape for providers that send a flat one; Twilio's identifier
is `sid`, not `id`; creates answered 200 where the provider answers 201; and
the list envelope carried a `next_cursor` field that Stripe does not send,
which is the worst kind of infidelity because code written against it works
locally and breaks in production.

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
  mbta                     a format this cannot read: this is Swagger 2.0, and only OpenAPI 3 is read

11 unchanged, 0 moved, 0 unreachable, 1 in a format this cannot read, 0 unrecorded, 286 with no description to check.
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
| `unsupported` | The provider publishes a description in a format this cannot read. Separate from `unreachable` because it is permanent -- api-v3.mbta.com serves Swagger 2.0, answers instantly, and will never once be unreachable |
| `unrecorded` | A description is named and has never been fingerprinted, so nothing was compared |
| `undeclared` | The provider publishes nothing to check. 286 of the 298 |

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

## Design commitments

These are the decisions the project intends to be held to.

**Detection never guesses.** Package-to-Recipe mapping is an explicit table. A wrong guess is worse than no guess, because booting the wrong fake sends someone chasing a bug that doesn't exist. Anything unrecognised is reported, never silently faked.

A project can also say so itself. `cauldron add mercury` writes a `cauldron.yaml` listing the providers it talks to, which is how a project reaches a Recipe no dependency maps to, or one it talks to over raw HTTP with no library at all. The first `add` copies whatever detection already found into the file, so starting one loses nothing, and from then on the file is the answer rather than the guess.

The table reaches 271 of the 298 Recipes. The other 27 ship and can be named
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

The most valuable contribution is a Recipe for an API you actually use, especially one whose sandbox is painful to get. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licence

[Apache-2.0](LICENSE).

A [Brilliance Digital](https://github.com/BrillianceDigital) project.
