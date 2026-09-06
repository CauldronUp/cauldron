# Chroma

Emulates the Chroma API (v2), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its add **returns an object with no properties in it**,
formally declared as `{"type": "object"}` in its own live description.

## Sources

- Documentation: https://docs.trychroma.com/docs/overview/getting-started
- Machine-readable description: https://api.trychroma.com/openapi.json, last checked 2026-09-01
  `cauldron drift chroma` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve chroma     # run it
cauldron verify chroma -v # check every claim
```
