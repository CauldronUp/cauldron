# Split

Emulates the Split API (2), for local development and tests.

**9 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **failure sentence echoes the last five characters of
your key**, unmasked, inside the message explaining the rejection -- while the
mask before them is a fixed width that says nothing about the real length.

## Sources

- Documentation: https://docs.split.io/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve split     # run it
cauldron verify split -v # check every claim
```
