# Snipcart

Emulates the Snipcart API (v3), for local development and tests.

**14 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **cart 404s at the moment its order appears**: the
same token, refused on one endpoint and answered on the other, modelled by
putting it in one fixture and not the other.

## Sources

- Documentation: https://docs.snipcart.com/v3/api-reference/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve snipcart     # run it
cauldron verify snipcart -v # check every claim
```
