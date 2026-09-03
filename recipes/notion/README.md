# Notion

Emulates the Notion API (2022-06-28), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It refuses a request that carries no version header. Forgetting that header is
the classic Notion integration bug, and a fake that waved it through would let
code ship that fails on its first real call.

## Sources

- Documentation: https://developers.notion.com/reference/intro
- Machine-readable description: https://developers.notion.com/openapi.json, last checked 2026-08-31
  `cauldron drift notion` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve notion     # run it
cauldron verify notion -v # check every claim
```
