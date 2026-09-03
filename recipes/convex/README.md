# Convex

Emulates the Convex API (http-api), for local development and tests.

**7 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**Nothing tells you the deployment is wrong**. The
hostname wildcard always resolves and a gateway validates the body's shape before
asking whether the deployment exists.

## Sources

- Documentation: https://docs.convex.dev/http-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve convex     # run it
cauldron verify convex -v # check every claim
```
