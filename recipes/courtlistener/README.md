# CourtListener

Emulates the CourtListener API (v4), for local development and tests.

**22 conformance cases, 16 checked against the live API on 2026-09-02.**

## What this Recipe found

**A sealed record is undetectable by
design** -- its own documentation says the system has no way to know an item is
sealed.

## Sources

- Documentation: https://wiki.free.law/c/courtlistener/help/api/rest/v4/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve courtlistener     # run it
cauldron verify courtlistener -v # check every claim
```
