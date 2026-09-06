# Duffel

Emulates the Duffel API (v2), for local development and tests.

**14 conformance cases, 1 checked against the live API on 2026-08-31.**

## What this Recipe found

It **refuses you for not naming a version before it checks
who you are.** Its remaining seven cases carry no verified date on purpose: the
runtime checks credentials before routing, so demonstrating a version failure
here needs a credential the real probe did not.

## Sources

- Documentation: https://duffel.com/docs/api/overview/welcome
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve duffel     # run it
cauldron verify duffel -v # check every claim
```
