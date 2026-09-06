# Inngest

Emulates the Inngest API (v1), for local development and tests.

**14 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose failures declare `text/plain` and send JSON.** A
client that trusts the content type and treats the body as a string gets one; a
client that calls `.json()` succeeds despite being told not to.

## Sources

- Documentation: https://www.inngest.com/docs/events
- Machine-readable description: https://api-docs.inngest.com/api-specs/v1.json, last checked 2026-09-01
  `cauldron drift inngest` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve inngest     # run it
cauldron verify inngest -v # check every claim
```
