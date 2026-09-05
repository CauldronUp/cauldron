# Exa

Emulates the Exa API (unversioned), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

**The one search API that labels a result cached**, per
document, where its two neighbours label nothing.

## Sources

- Documentation: https://docs.exa.ai/reference/search
- Machine-readable description: https://api.exa.ai/openapi.json, last checked 2026-09-05
  `cauldron drift exa` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve exa     # run it
cauldron verify exa -v # check every claim
```
