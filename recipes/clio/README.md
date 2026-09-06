# Clio

Emulates the Clio API (v4), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **two 401s share no field names** -- `error` is an object
in one and a bare string in the other.

## Sources

- Documentation: https://docs.developers.clio.com/api-docs/
- Machine-readable description: https://docs.developers.clio.com/openapi.json, last checked 2026-09-01
  `cauldron drift clio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clio     # run it
cauldron verify clio -v # check every claim
```
