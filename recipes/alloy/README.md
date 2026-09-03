# Alloy

Emulates the Alloy API (v1), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

**Whose refusals split across two statuses**: unreadable is a
400 and readable-but-wrong is a 401.

## Sources

- Documentation: https://developer.alloy.com/reference/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve alloy     # run it
cauldron verify alloy -v # check every claim
```
