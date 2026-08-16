<div align="center">

# Cauldron

**Your code. One command. Every dependency.**

The open-source environment compiler — it boots your application *and the third-party APIs it depends on*, locally, from one command.

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

Cauldron reads the manifests already in your repository, works out what the project needs, and boots it — including working emulations of the providers it talks to.

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

No recipe yet — these will still reach the real network:
  !  acme/weather-api-client  (composer.json)
```

That last section is deliberate. Falling back to the real network *silently* is how a test suite starts lying.

## Status

**Early. Not yet usable for real work.**

| Area | State |
|---|---|
| Detection engine — Composer, npm, Go modules | Working, tested |
| `cauldron detect` | Working |
| `cauldron up` | Prints the plan; orchestration not built |
| Recipe format and runtime | In progress |
| Container orchestration | Not started |

`cauldron up` currently tells you what it *would* do and says so plainly. It will not print a convincing boot log for work that didn't happen.

## Try the detector

```bash
go build -o cauldron ./cmd/cauldron
./cauldron detect /path/to/your/project
```

## Recipes

A **Recipe** defines how Cauldron emulates an external dependency. It is not a mock — anyone can return `200 OK`. A Recipe models behaviour, and carries the tests that prove it still matches the real API:

| Part | Covers |
|---|---|
| Auth | OAuth flows, API keys, HMAC signing, token refresh and expiry |
| Resources | Schemas, ID formats, relationships between objects |
| Routes | Paths, pagination and cursors, filtering, sorting |
| Behaviour | State transitions and side effects — an order decrements inventory |
| Webhooks | Event catalogue, signing, delivery order, retries, duplicates |
| Errors | The provider's real error taxonomy and rate-limit shape |
| Fixtures | Named seed data — `empty`, `small-shop`, `enterprise` |
| Conformance | The suite that proves the Recipe still matches the real API |

## Design commitments

These are the decisions the project intends to be held to.

**Detection never guesses.** Package-to-Recipe mapping is an explicit table. A wrong guess is worse than no guess — booting the wrong fake sends someone chasing a bug that doesn't exist. Anything unrecognised is reported, never silently faked.

**Determinism is enforced at the boundary, not requested politely.** Recipes get a seeded `ctx.now()` and `ctx.random()`. No ambient clock, no ambient randomness, no network. The same fixture and seed produce the same IDs on every machine.

**Nothing is held back to make a paid tier viable.** Everything that runs on your machine is Apache-2.0 and stays that way.

**Fidelity is measured, not claimed.** Recipes ship with conformance suites, and where a Recipe diverges from the real API that gets published rather than hidden.

## Contributing

The most valuable contribution is a Recipe for an API you actually use — especially one whose sandbox is painful to get. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licence

[Apache-2.0](LICENSE).

A [Brilliance Digital](https://github.com/BrillianceDigital) project.
