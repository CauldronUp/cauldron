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
| Detection engine (Composer, npm, Go modules) | Working. 202 of 224 Recipes are reachable from a dependency |
| Recipe format and validator | Working |
| Recipe runtime (routing, state, auth, pagination) | Working |
| Webhooks (lifecycle events, signing, delivery) | Working. Payload envelopes are declarable per Recipe. Of the 102 Recipes that emit events, 44 declare one and 58 fall back to a default modelled on Stripe's and not equal to it -- Stripe's own Recipe now declares the fuller real thing -- so the fallback is a convention rather than a claim about any provider |
| Fault injection, clock control, request log | Working |
| Network conditions (latency, throttling, timeouts, resets) | Working. Toxiproxy's vocabulary, applied to the emulated providers |
| `serve`, `status`, `requests`, `seed`, `reset`, `fault`, `network`, `emit`, `clock` | Working |
| `doctor`, `logs`, `open` | Working |
| `cauldron up` / `down` (container orchestration) | Working for backing services |
| `snapshot` save/restore | Working |
| Conformance suites (`cauldron verify`) | Working. 2347 cases, 58 of them checked against a live API |
| Scoped multi-segment paths (`/repos/{owner}/{repo}/…`) | Working |
| Headless mode (`--headless`, `--host`) | Working. Providers only, one line of JSON, no containers |
| Application runtimes in containers | Not built. Run your app as you normally do |
| Recipes shipped | 224: Stripe, GitHub, GitLab, Bitbucket, Shopify, Shopify GraphQL, BigCommerce, Etsy, Magento, eBay, Amazon SP-API, EasyPost, AfterShip, ShipStation, ShipEngine, DHL, Recharge, Lemon Squeezy, Toast, Avalara, FedEx, Easyship, Clover, Lightspeed, Printful, Printify, Medusa, Shopware, commercetools, VTEX, Saleor, Allegro, Akeneo, Voucherify, Apideck, Royal Mail, Lago, Polar, Gumroad, Ecwid, Squarespace, Wix, Metronome, Twilio, Slack, HubSpot, SendGrid, Airtable, Notion, Zendesk, Postmark, Plaid, Clerk, Intercom, Discord, Square, Mailchimp, Cloudflare, Vercel, Xero, QuickBooks, PagerDuty, Asana, Algolia, Sentry, Box, Calendly, Datadog, Front, Typeform, Miro, Contentful, Sanity, Klaviyo, Webflow, Zoom, Pipedrive, Freshdesk, Mailgun, Okta, Shippo, Dropbox, Google Drive, Google Cloud Storage, Auth0, Recurly, Trello, Paddle, CircleCI, Snyk, Statuspage, Buildkite, ClickUp, Basecamp, Shortcut, Chargebee, Vonage, Rollbar, Docusign, Lob, Segment, Greenhouse, Adyen, Salesforce, Help Scout, DigitalOcean, WooCommerce, WordPress, AWS SQS, Google Pub/Sub, AWS DynamoDB, AWS Secrets Manager, AWS SES, Resend, Cloudinary, LaunchDarkly, Vimeo, Mux, Pusher, Ably, Ghost, WorkOS, TaxJar, Stytch, Keycloak, FusionAuth, AssemblyAI, Documenso, Dropbox Sign, Deepgram, Gusto, OneSignal, Telnyx, Bill.com, Ramp, Mercury, Brex, Deel, Google Calendar, Jira, Confluence, Jira Service Management, OpenAI, ElevenLabs, PayPal, RingCentral, Mollie, Gmail, Increase, Svix, Docker Hub, Onfido, Netlify, Firecrawl, Modern Treasury, npm registry, Bugsnag, Persona, Orb, PostHog, Column, marqeta, braze, bandwidth, alpaca, meilisearch, knock, merge, truelayer, razorpay, gocardlessbank, finch, replicate, cohere, hightouch, dwolla, airwallex, snowflake, incidentio, fivetran, novu, uploadcare, midtrans, courier, opsgenie, xendit, gorgias, fastspring, kustomer, calcom, twilioverify, qdrant, revenuecat, backblaze, kratos, hetzner, heroku, polygon, fly, apify, finnhub, tradier, typesense, weaviate, hydra, openrouter, Gemini, Google Address Validation, ShipHero, Linear, Monday, Attio, Neon, PlanetScale, MongoDB Atlas, Supabase, Discourse, Langfuse, YouTube, HashiCorp Vault, Unleash, Pinecone, Hookdeck, Grafana |

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

That second line is the honest one. Of every Recipe: 2347 cases, 58 run against
a live account and 2289 not. Documentation-derived cases are worth having,
and they are not the same as watching the provider do it. Adding a `verified:`
date to a case is a claim that someone did.

The fifty-eight are the cases whose provider can be asked without a key. Six are
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

The table reaches 202 of the 224 Recipes. The other 22 ship and can be named
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
