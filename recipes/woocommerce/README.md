# WooCommerce

Emulates the WooCommerce API (wc/v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

`processing` is the paid state; `completed` means fulfilled and shipped by a human clicking a button in wp-admin. An integration that waits for `completed` before granting access to a digital product waits for something that, on a digital product, never happens. Every monetary value is also a decimal string -- `total` is `"49.99"`, not a number -- so naive addition in JavaScript concatenates two totals into `"49.9912.00"` rather than throwing, and it gets stringified right back out by anything downstream that doesn't parse it first.

`meta_data` is an array of `{id, key, value}` objects rather than a map, so `order.meta_data.some_key` finds nothing and the value has to be found by scanning the array.

## Sources

- Documentation: https://woocommerce.github.io/woocommerce-rest-api-docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve woocommerce     # run it
cauldron verify woocommerce -v # check every claim
```
