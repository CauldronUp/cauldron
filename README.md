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
| Detection engine (Composer, npm, Go modules) | Working |
| Recipe format and validator | Working |
| Recipe runtime (routing, state, auth, pagination) | Working |
| Webhooks (lifecycle events, signing, delivery) | Working. Payload envelopes are declarable per Recipe; Recipes that declare none fall back to Stripe's shape, which is a default rather than a claim about that provider |
| Fault injection, clock control, request log | Working |
| `serve`, `status`, `requests`, `seed`, `reset`, `fault`, `emit`, `clock` | Working |
| `doctor`, `logs`, `open` | Working |
| `cauldron up` / `down` (container orchestration) | Working for backing services |
| `snapshot` save/restore | Working |
| Conformance suites (`cauldron verify`) | Working. 860 cases, all from documentation so far |
| Scoped multi-segment paths (`/repos/{owner}/{repo}/…`) | Working |
| Headless mode (`--headless`, `--host`) | Working. Providers only, one line of JSON, no containers |
| Application runtimes in containers | Not built. Run your app as you normally do |
| Recipes shipped | 102: Stripe, GitHub, GitLab, Bitbucket, Shopify, Twilio, Slack, HubSpot, SendGrid, Airtable, Notion, Zendesk, Postmark, Plaid, Clerk, Intercom, Discord, Square, Mailchimp, Cloudflare, Vercel, Xero, QuickBooks, PagerDuty, Asana, Algolia, Sentry, Box, Calendly, Datadog, Front, Typeform, Miro, Contentful, Sanity, Klaviyo, Webflow, Zoom, Pipedrive, Freshdesk, Mailgun, Okta, Shippo, Dropbox, Auth0, Recurly, Trello, Paddle, CircleCI, Snyk, Statuspage, Buildkite, ClickUp, Basecamp, Shortcut, Chargebee, Vonage, Rollbar, Docusign, Lob, Segment, Greenhouse, Adyen, Salesforce, Help Scout, DigitalOcean, WooCommerce, WordPress, AWS SQS, Google Pub/Sub, AWS DynamoDB, AWS Secrets Manager, AWS SES, Resend, Cloudinary, LaunchDarkly, Mux, Pusher, Ably, Ghost, WorkOS, TaxJar, Stytch, AssemblyAI, Documenso, Deepgram, Gusto, OneSignal, Telnyx, Bill.com, Ramp, Mercury, Brex, Deel, Google Calendar, Jira, PayPal, RingCentral, Mollie, Gmail, Increase, Svix |

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

Or drive it from the CLI:

```bash
cauldron status                                   # what is running
cauldron seed stripe --fixture small-shop         # load seed data
cauldron fault stripe --error rate_limit --count 2 # break the next two calls
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
  8 of 8 cases passed
  8 from documentation only, none checked against the real API
```

That second line is the honest one, and today it reads the same for every
Recipe: 860 cases, none yet run against a live account. Documentation-derived
cases are worth having, and they are not the same as watching the provider do
it. Adding a `verified:` date to a case is a claim that someone did.

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

Twenty-four of the last thirty-five needed no format change at all. Writing a
Recipe is now mostly research rather than engineering: the work is reading the
provider closely enough to know which state nobody branches on, not teaching
the format a new trick.

The suite has started finding faults in itself as well. Vonage reports a
successful send with the string `"0"` and Docusign counts with strings, and the
checker was comparing scalars by their rendered text, so `"0"` and `0` were
indistinguishable and both of those cases passed whichever the emulator sent.
Scalars now compare kind as well as value, which immediately turned up a real
inconsistency: the nested error style was sending PagerDuty's numeric codes as
text.

Not every provider fits, and two are left out deliberately rather than faked.
Linear is GraphQL: a single endpoint where the request body decides the
response shape, so a Recipe that ignored the query and returned fixed data
would let you ship a query with wrong field names and still pass. Attio stores
every attribute as an array of historical values, and a client reads
`record.values.name[0].value`; modelling those as scalars would teach a shape
that does not exist. In both cases no Recipe is better than a misleading one.

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
roughly in order, along with the two that were assessed and deliberately left
out. A Recipe has to earn its place: the test is whether the provider has
behaviour worth reproducing, not whether adding it raises the count.

**Fidelity will be measured, not claimed.** The intent is that every Recipe carries a conformance suite run against the real provider, and that divergence is published rather than hidden. This is not built yet. The Recipes here are hand-modelled and unverified, and the status table says so.

## Contributing

The most valuable contribution is a Recipe for an API you actually use, especially one whose sandbox is painful to get. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licence

[Apache-2.0](LICENSE).

A [Brilliance Digital](https://github.com/BrillianceDigital) project.
