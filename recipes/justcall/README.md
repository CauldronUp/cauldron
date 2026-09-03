# JustCall

Emulates the JustCall API (2.1), for local development and tests.

**10 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its recording link has **three shapes and no lifetime**
across its own three documents.

## Sources

- Documentation: https://developer.justcall.io/reference/call_get_v21
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve justcall     # run it
cauldron verify justcall -v # check every claim
```
