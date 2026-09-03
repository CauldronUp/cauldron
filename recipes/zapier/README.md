# Zapier

Emulates the Zapier API (v1), for local development and tests.

**12 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **catch hook accepts an id nobody was ever issued** --
200 success for an invented account and token, with no endpoint anywhere to check
the claim against.

## Sources

- Documentation: https://platform.zapier.com/reference/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zapier     # run it
cauldron verify zapier -v # check every claim
```
