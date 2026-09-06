# Recharge

Emulates the Recharge API (2021-11), for local development and tests.

**19 conformance cases, 3 checked against the live API.**

Everything past the credential and routing checks still cites documentation rather than an observation, because reaching it needs a real store. Those checks were verified directly against api.rechargeapps.com, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: no access token at all and a fictitious one both answer `{"error":"bad authentication"}` -- singular `error`, not the plural `errors` this file had the message field named for, and different wording too. Routing runs ahead of the credential entirely: an unrouted path answers `{"error":"Not Found"}` and a wrong method on a real path answers `{"error":"method not allowed"}`, neither needing anything sent at all.

Recharge sits on top of Shopify, so one purchase effectively exists as two separate orders -- Recharge bills a `charge`, Shopify records an `order`, and each carries its own identifier for the same money, linked only by a quiet `shopify_order_id` field. Joining the two systems on the wrong identifier either silently loses every subscription order or double-counts every one of them, and it is an easy mistake because both objects are called an order somewhere and both have an `id`.

The API version travels as a header rather than the URL, and omitting it does not fail -- it silently falls back to whatever the store's own admin-configurable default is, so an integration can break without anyone deploying anything. Money is a quoted string (`"19.99"`), so naive concatenation replaces addition, and `next_charge_scheduled_at` is a date with no time and no timezone, the store's own local date, so a job running at midnight UTC can charge a day early or late depending on which side of the date line the store sits.

The most consequential finding is the last one: a subscription stays `ACTIVE` while its charges are failing. Recharge retries a declined card several times before giving up, and the subscription's own status says nothing about it, so a report counting active subscribers counts people whose money stopped arriving weeks earlier -- the only way to see it is to read the charges sitting beside the subscription, not the subscription's status field.

## Sources

- Documentation: https://developer.rechargepayments.com/2021-11
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve recharge     # run it
cauldron verify recharge -v # check every claim
```
