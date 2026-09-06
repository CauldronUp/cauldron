# Buttondown

Emulates the Buttondown API (v1), for local development and tests.

**17 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**Unsubscribing does not delete anybody**: the proof
is a create parameter whose default refuses a second create for an address that
already exists.

## Sources

- Documentation: https://docs.buttondown.com/api-subscribers-introduction
- Machine-readable description: https://docs.buttondown.com/openapi.json, last checked 2026-09-01
  `cauldron drift buttondown` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve buttondown     # run it
cauldron verify buttondown -v # check every claim
```
