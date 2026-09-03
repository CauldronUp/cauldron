# Notion

Emulates the Notion API (2022-06-28), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developers.notion.com/reference/intro
- Machine-readable description: https://developers.notion.com/openapi.json, last checked 2026-08-31
  `cauldron drift notion` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve notion     # run it
cauldron verify notion -v # check every claim
```
