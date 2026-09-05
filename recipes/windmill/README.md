# Windmill

Emulates the Windmill API (1.798.1), for local development and tests.

**12 conformance cases, 11 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **one unauthenticated route leaks a source
location**: a job lookup needing no credential answers with a Rust file and line
number.

## Sources

- Documentation: https://www.windmill.dev/docs/core_concepts/webhooks
- Machine-readable description: https://app.windmill.dev/api/openapi.json, last checked 2026-09-05
  `cauldron drift windmill` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve windmill     # run it
cauldron verify windmill -v # check every claim
```
