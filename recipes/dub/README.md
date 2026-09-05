# Dub

Emulates the Dub API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**A bad key is not rejected until a workspace resolves**:
the same wrong token answers a missing-workspace 404 without one and an
invalid-key 401 with one.

## Sources

- Documentation: https://dub.co/docs/api-reference/introduction
- Machine-readable description: https://dub.co/openapi.json, last checked 2026-09-05
  `cauldron drift dub` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dub     # run it
cauldron verify dub -v # check every claim
```
