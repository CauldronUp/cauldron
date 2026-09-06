# Bunny

Emulates the Bunny API (v1), for local development and tests.

**23 conformance cases, 11 checked against the live API on 2026-09-01.**

## What this Recipe found

**Sending no body at all is a 500** -- and attaching
any body, still with no credential, reverts to an ordinary 401.

## Sources

- Documentation: https://docs.bunny.net/reference/bunnynet-api-overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bunny     # run it
cauldron verify bunny -v # check every claim
```
