# NeverBounce

Emulates the NeverBounce API (v4.2), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

It **calls its ambiguous answer unverifiable**, and
whose documented bearer header is silently never read.

## Sources

- Documentation: https://developers.neverbounce.com/reference/single-check
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve neverbounce     # run it
cauldron verify neverbounce -v # check every claim
```
