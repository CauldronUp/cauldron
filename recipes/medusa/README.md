# Medusa

Emulates the Medusa API (v2), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The sharpest thing here is one line in Medusa's own spec about the `fields` query parameter: without a prefix it replaces the default field set entirely, with `+` it adds to it. So `?fields=id` does not give you an order with its id emphasised, it gives you an order with nothing else at all -- no total, no status, no email -- and it is not an error. A developer trying to trim a response to one extra field instead gets a response with only that field, and the code reading everything else starts seeing `undefined` for a reason documented in a sentence nobody reads twice.

Two identifiers exist on every order for different purposes: `id` is a string and `display_id` is a number, and `display_id` is the one a customer actually quotes -- joining on the wrong one finds nothing. `x-publishable-api-key` is required on every store route and is not the credential; it is checked separately from authentication, so a valid token without the key is refused for a reason that has nothing to do with the token. An order also carries seven separate totals (total, subtotal, item_total, original_total, tax_total, discount_total, shipping_total), and which one a report should use depends on the question, not the name.

Only the Storefront API is modelled; Medusa's Admin API has its own credential and shapes entirely, and carts, checkout, payment collections and customer accounts are not covered.

## Sources

- Documentation: https://docs.medusajs.com/api/store
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve medusa     # run it
cauldron verify medusa -v # check every claim
```
