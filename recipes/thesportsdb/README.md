# TheSportsDB

Emulates the TheSportsDB API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**A finished match says nothing about being
over**: a 2019 fixture carries a null status and only its two score fields
prove it happened.

## Sources

- Documentation: https://www.thesportsdb.com/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve thesportsdb     # run it
cauldron verify thesportsdb -v # check every claim
```
