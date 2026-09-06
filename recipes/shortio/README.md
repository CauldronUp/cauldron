# Short.io

Emulates the Short.io API (v1), for local development and tests.

**16 conformance cases, 10 checked against the live API on 2026-09-03.**

## What this Recipe found

**One fact in four envelopes** -- the same unauthorized
condition rendered four different ways depending only on which route was asked.

## Sources

- Documentation: https://developers.short.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shortio     # run it
cauldron verify shortio -v # check every claim
```
