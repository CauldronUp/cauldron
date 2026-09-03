# Heap

Emulates the Heap API (v1), for local development and tests.

**4 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**The application identifier authenticates nothing** -- one
never registered gets the identical success a real one would.

## Sources

- Documentation: https://developers.heap.io/reference/track-1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve heap     # run it
cauldron verify heap -v # check every claim
```
