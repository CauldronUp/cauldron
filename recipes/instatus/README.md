# Instatus

Emulates the Instatus API (v1), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

**A monitor can be up and paused at once** -- two
independent fields, and no last-checked time on the record at all.

## Sources

- Documentation: https://instatus.com/help/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve instatus     # run it
cauldron verify instatus -v # check every claim
```
