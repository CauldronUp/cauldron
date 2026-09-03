# Toggl

Emulates the Toggl API (v9), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**The negative duration nobody can confirm**: the sign is
still documented and the formula that made it mean something is gone.

## Sources

- Documentation: https://engineering.toggl.com/docs/track/api/time_entries
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve toggl     # run it
cauldron verify toggl -v # check every claim
```
