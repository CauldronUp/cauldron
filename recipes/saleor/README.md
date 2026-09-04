# Saleor

Emulates the Saleor API (graphql), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Saleor's Money type carries the same amount twice and the currency once: an `amount` field that is an ordinary GraphQL float, and a `fractionalAmount` integer in minor units beside it. Three line items at 19.99 sum to 59.969999999999999 through the field with the obvious name, and to 5997 exactly through the one nobody reaches for. The disagreement shows up on a basket total or a monthly report, not on any figure a reviewer checks by hand.

A product also has no price of its own -- it has `pricing.priceRange.start.gross.amount`, five levels down and nullable at every level, because variants can differ. And five schema fields are deprecated with the reason "Always returns `null`", so a client can never tell whether a value is genuinely absent or was never implemented.

## Sources

- Documentation: https://docs.saleor.io/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve saleor     # run it
cauldron verify saleor -v # check every claim
```
