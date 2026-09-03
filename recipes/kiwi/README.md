# Kiwi.com

Emulates the Kiwi.com API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**An API that answers behind documentation that does not** --
alive and probeable, with a docs portal wanting an account at the vendor's own
domain.

## Sources

- Documentation: https://tequila.kiwi.com/portal/docs/tequila_api/search_api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kiwi     # run it
cauldron verify kiwi -v # check every claim
```
