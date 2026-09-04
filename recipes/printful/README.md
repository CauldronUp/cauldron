# Printful

Emulates the Printful API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Printful order carries two full sets of money with identical field names inside each -- `costs` (what Printful charges you) and `retail_costs` (what you charged your customer) both have their own `subtotal`, `shipping`, `tax`, and `total`, with nothing in either name saying which is which. Reading `costs.total` as revenue reports your supplier's invoice as income; reading `retail_costs.total` as cost of goods reports your own price as what you paid. They are not even guaranteed to share a currency -- Printful bills in the store's currency while `retail_costs` reflects whatever the customer actually paid, so one order can carry two totals in two currencies, each with its own easily-ignored currency field.

An order also has two identifiers, Printful's own id and the `external_id` you supplied, and endpoints accept either, prefixed with `@` to mean yours, so the same path segment means a different lookup depending on one character. A draft order is not an order anyone will actually make; `status` stays `draft` until confirmed, and shipping on a draft is only an estimate that gets recalculated at confirmation, so the number quoted to a customer in the basket can change by the time of the invoice. Confirming an order is not modelled here, so a draft stays a draft.

## Sources

- Documentation: https://developers.printful.com/docs/
- Machine-readable description: https://developers.printful.com/docs/openapi.json, last checked 2026-08-31
  `cauldron drift printful` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve printful     # run it
cauldron verify printful -v # check every claim
```
