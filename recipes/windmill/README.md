# Windmill

Emulates the Windmill API (1.798.1), for local development and tests.

**12 conformance cases, 11 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **one unauthenticated route leaks a source
location**: a job lookup needing no credential answers with a Rust file and line
number.

## Sources

- Documentation: https://www.windmill.dev/docs/core_concepts/webhooks
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve windmill     # run it
cauldron verify windmill -v # check every claim
```
