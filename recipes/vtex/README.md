# VTEX

Emulates the VTEX API (oms), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

One purchase is not one order. VTEX's own documentation for its order-group path parameter says so directly: when an order is fulfilled by multiple sellers, each seller gets its own order id, and all of them share one order-group id. Code that stores "the order id" stores one seller's share of a purchase, a naive count of rows double-counts a marketplace basket, and a refund issued against one id leaves the other seller's order still paid.

The order listing also stopped returning line items in 2018, documented with the exact date, while the schema still lists `items` as a field -- a client reading `order.items` off a listing gets `undefined` for every order and has to fetch each one individually. And paging on the order listing stops at thirty pages, so a backfill walking it can never retrieve more than three thousand orders no matter how it filters.

## Sources

- Documentation: https://developers.vtex.com/docs/api-reference/orders-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vtex     # run it
cauldron verify vtex -v # check every claim
```
