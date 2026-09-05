# Clearbit

Emulates the Clearbit API (2019-12-19), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

It **is not retired and cannot be read**: the hosts still
answer and still stamp a 2019 version header, and the documentation now
redirects to a login.

## Sources

- Documentation: https://dashboard.clearbit.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clearbit     # run it
cauldron verify clearbit -v # check every claim
```
