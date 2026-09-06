# Spoonacular

Emulates the Spoonacular API (v1), for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-03.**

## What this Recipe found

It **says per-serving everywhere but on the wire** --
the clearest statement of the group, in the documentation rather than the body.

## Sources

- Documentation: https://spoonacular.com/food-api/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve spoonacular     # run it
cauldron verify spoonacular -v # check every claim
```
