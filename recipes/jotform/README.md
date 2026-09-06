# Jotform

Emulates the Jotform API (1), for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **envelope and status can disagree**, and whose
not-found echoes the path parameter's name rather than the value it was given.

## Sources

- Documentation: https://api.jotform.com/docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve jotform     # run it
cauldron verify jotform -v # check every claim
```
