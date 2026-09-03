# Akeyless

Emulates the Akeyless API (v2), for local development and tests.

**10 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

Its refusal **leaks a live session id and a console audit
URL** into the body handed to whoever sent a bad token.

## Sources

- Documentation: https://docs.akeyless.io/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve akeyless     # run it
cauldron verify akeyless -v # check every claim
```
