# ZeroBounce

Emulates the ZeroBounce API (v2), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

It **asks this group's question of itself** in its own
documentation, and answers it across two fields.

## Sources

- Documentation: https://www.zerobounce.net/docs/email-validation-api-quickstart/v2-validate-emails
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zerobounce     # run it
cauldron verify zerobounce -v # check every claim
```
