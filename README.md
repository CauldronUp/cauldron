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
| Webhooks (lifecycle events, signing, delivery) | Working |
| Fault injection, clock control, request log | Working |
| `serve`, `status`, `requests`, `seed`, `reset`, `fault`, `emit`, `clock` | Working |
| `cauldron up` (container orchestration) | **Not built.** `up` boots the emulated providers and says what it skipped |
| `snapshot` save/restore | Not built |
| Conformance suites (measured fidelity) | Not built |
| Scoped multi-segment paths (`/repos/{owner}/{repo}/…`) | Working |
| Recipes shipped | Stripe, GitHub |

## Try it

```bash
go build -o cauldron ./cmd/cauldron

./cauldron detect /path/to/your/project
./cauldron serve --fixture small-shop stripe
```

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
  -d '{"error":"rate_limit","count":1}'

# Age everything by a month, so a subscription falls into dunning
curl -X POST http://127.0.0.1:4600/_cauldron/clock/advance \
  -d '{"duration":"30d"}'

# Fire a signed webhook at your application
curl -X POST http://127.0.0.1:4600/_cauldron/stripe/emit \
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

## Recipes

A **Recipe** defines how Cauldron emulates an external dependency. It is not a mock, because anyone can return `200 OK`. A Recipe models behaviour, and carries the tests that prove it still matches the real API:

| Part | Covers |
|---|---|
| Auth | OAuth flows, API keys, HMAC signing, token refresh and expiry |
| Resources | Schemas, ID formats, relationships between objects |
| Routes | Paths, pagination and cursors, filtering, sorting |
| Behaviour | State transitions and side effects, so an order decrements inventory |
| Webhooks | Event catalogue, signing, delivery order, retries, duplicates |
| Errors | The provider's real error taxonomy and rate-limit shape |
| Fixtures | Named seed data: `empty`, `small-shop`, `enterprise` |
| Conformance | The suite that proves the Recipe still matches the real API |

## Design commitments

These are the decisions the project intends to be held to.

**Detection never guesses.** Package-to-Recipe mapping is an explicit table. A wrong guess is worse than no guess, because booting the wrong fake sends someone chasing a bug that doesn't exist. Anything unrecognised is reported, never silently faked.

**Determinism is enforced at the boundary, not requested politely.** Time comes from a movable sandbox clock and identifiers from a seeded generator. A Recipe has no access to the wall clock or to unseeded randomness. The same fixture and seed produce the same IDs on every machine, and a reset sandbox is indistinguishable from a fresh one.

**Nothing is held back to make a paid tier viable.** Everything that runs on your machine is Apache-2.0 and stays that way.

**Fidelity will be measured, not claimed.** The intent is that every Recipe carries a conformance suite run against the real provider, and that divergence is published rather than hidden. This is not built yet. The Recipes here are hand-modelled and unverified, and the status table says so.

## Contributing

The most valuable contribution is a Recipe for an API you actually use, especially one whose sandbox is painful to get. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licence

[Apache-2.0](LICENSE).

A [Brilliance Digital](https://github.com/BrillianceDigital) project.
