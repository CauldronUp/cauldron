# Render

Emulates the Render API (v1), for local development and tests.

**11 conformance cases, 2 checked against the live API on 2026-08-31.**

## What this Recipe found

One is Upstash's. Those are the live-checked halves of two Recipes
otherwise drafted from their descriptions. Render does not distinguish a missing
credential from a wrong one; Upstash answers three different shapes depending on
how you fail to authenticate, only one of which matches its documented scheme.

## Sources

- Documentation: https://render.com/docs/api
- Machine-readable description: https://api-docs.render.com/v1.0/openapi/render-public-api-1.json, last checked 2026-08-31
  `cauldron drift render` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve render     # run it
cauldron verify render -v # check every claim
```
