# Upstash

Emulates the Upstash API (v2), for local development and tests.

**6 conformance cases, 1 checked against the live API on 2026-08-31.**

## Sources

- Documentation: https://upstash.com/docs/devops/developer-api/introduction
- Machine-readable description: https://upstash.com/openapi.json, last checked 2026-08-31
  `cauldron drift upstash` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve upstash     # run it
cauldron verify upstash -v # check every claim
```
