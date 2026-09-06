# WooCommerce

Emulates the WooCommerce API (wc/v3), for local development and tests.

**24 conformance cases, 4 checked against the live API on 2026-09-05.**

WooCommerce has no shared sandbox, since every install is a merchant's own site -- except one. woo.com runs its own storefront on WooCommerce, so its `/wp-json/wc/v3` is a real production API with nothing special about its unauthenticated behaviour, and checking it live found two claims here that had never been run.

## What this Recipe found

`processing` is the paid state; `completed` means fulfilled and shipped by a human clicking a button in wp-admin. An integration that waits for `completed` before granting access to a digital product waits for something that, on a digital product, never happens. Every monetary value is also a decimal string -- `total` is `"49.99"`, not a number -- so naive addition in JavaScript concatenates two totals into `"49.9912.00"` rather than throwing, and it gets stringified right back out by anything downstream that doesn't parse it first.

`meta_data` is an array of `{id, key, value}` objects rather than a map, so `order.meta_data.some_key` finds nothing and the value has to be found by scanning the array.

## What checking it live found

A present Basic credential that is not a real consumer key/secret pair -- including an existing case's own fixture pair with a deliberately wrong secret -- is intercepted by WordPress's own core Application Passwords authentication before WooCommerce's REST layer ever runs, because a consumer key is never a valid WordPress username: it fails with WordPress's `"invalid_username"` sentence, not WooCommerce's `"cannot_view"` one. An absent credential is not intercepted the same way and does reach WooCommerce's own check, which this Recipe had modelled correctly all along. A path nothing declares is WordPress's own routing failure, resolved before any credential is judged.

## Sources

- Documentation: https://woocommerce.github.io/woocommerce-rest-api-docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve woocommerce     # run it
cauldron verify woocommerce -v # check every claim
```
