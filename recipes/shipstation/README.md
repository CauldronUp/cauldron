# ShipStation

Emulates the ShipStation API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every date ShipStation sends is missing its timezone -- not UTC without a `Z`, just missing entirely -- and the instant it refers to is the account's own local time, which the API never states. Every date library in wide use parses that as UTC or local time depending on which one it is, and both are guesses; a shipment that left at nine in the morning becomes one that left at one in the afternoon, and nothing errors because the numbers are close enough that nobody checks.

Creating an order is also an upsert: `POST /orders/createorder` with an `orderKey` that already exists updates that order instead of creating a duplicate, which makes a retry after a timeout safe and a genuine duplicate key silently overwrite the wrong order.

## Sources

- Documentation: https://www.shipstation.com/docs/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shipstation     # run it
cauldron verify shipstation -v # check every claim
```
