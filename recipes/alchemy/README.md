# Alchemy

Emulates the Alchemy API (v2), for local development and tests.

**8 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

**The credential is checked before the route exists** --
the key is a path segment, so there is no such thing as an unrouted path.

## Sources

- Documentation: https://www.alchemy.com/docs/reference/eth-getblockbynumber
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve alchemy     # run it
cauldron verify alchemy -v # check every claim
```
