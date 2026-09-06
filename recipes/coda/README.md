# Coda

Emulates the Coda API (v1), for local development and tests.

**21 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose failure states one fact four times** -- the status
line, `statusCode`, `statusMessage` and `message`, the last two identical by
Coda's own admission. Its published description also calls itself the Superhuman
Docs API and names a server the live host is not.

## Sources

- Documentation: https://coda.io/developers/apis/v1
- Machine-readable description: https://coda.io/apis/v1/openapi.json, last checked 2026-09-01
  `cauldron drift coda` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve coda     # run it
cauldron verify coda -v # check every claim
```
