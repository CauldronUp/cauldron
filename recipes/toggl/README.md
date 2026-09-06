# Toggl

Emulates the Toggl API (v9), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-03.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

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
