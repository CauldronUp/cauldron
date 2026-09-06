# Kickbox

Emulates the Kickbox API (v2), for local development and tests.

**16 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **sandbox is a tier rather than a parameter** -- the
key's own prefix decides, and reserved local-parts force each documented result.

## Sources

- Documentation: https://docs.kickbox.com/docs/single-verification-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kickbox     # run it
cauldron verify kickbox -v # check every claim
```
