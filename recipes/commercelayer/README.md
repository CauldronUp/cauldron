# Commerce Layer

Emulates the Commerce Layer API (core), for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Layer's, which **has no cart resource at all** -- "The shopping
cart is a draft order", in its own words -- and whose routing failures carry its
own live typo, "was not not found."

## Sources

- Documentation: https://docs.commercelayer.io/core-api-reference/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve commercelayer     # run it
cauldron verify commercelayer -v # check every claim
```
