# Anrok

Emulates the Anrok API (v1), for local development and tests.

**20 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its two calculate endpoints **differ by one absent field**:
a stored transaction carries `"version": 1` and an ephemeral quote does not.

## Sources

- Documentation: https://apidocs.anrok.com/
- Machine-readable description: https://apidocs.anrok.com/openapi.yaml, last checked 2026-09-01
  `cauldron drift anrok` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve anrok     # run it
cauldron verify anrok -v # check every claim
```
