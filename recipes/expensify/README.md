# Expensify

Emulates the Expensify API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **HTTP status is always 200**. The real one is a body
field reusing HTTP numbers with none of their meaning: a genuine authentication
failure answers `responseCode: 404`.

## Sources

- Documentation: https://integrations.expensify.com/Integration-Server/doc/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve expensify     # run it
cauldron verify expensify -v # check every claim
```
