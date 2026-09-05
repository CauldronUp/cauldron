# Printful

Emulates the Printful API (v1), for local development and tests.

**16 conformance cases, 4 checked against the live API.**

Everything past the credential and routing checks still cites documentation rather than an observation, because reaching it needs a real store. Those checks were verified directly against api.printful.com, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: an absent credential and a fictitious one are different sentences under the same code -- `{"code":401,"result":"Unauthorized","error":{"reason":"Unauthorized","message":"Unauthorized"}}` versus the same shape naming the token as the problem, `"The access token provided is invalid."` This file had modelled one message, "Invalid or missing API key.", for both. Routing runs ahead of the credential too, and is not one shape either: an unrouted path answers "Page not found." while a wrong method on a real path answers a genuine 404 (not 405) reading just "Not found" -- neither needs a credential at all.

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
