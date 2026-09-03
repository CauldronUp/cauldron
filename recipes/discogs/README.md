# Discogs

Emulates the Discogs API (v2), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**A merged release leaves no trace** -- no redirect, no
tombstone, and a sentence that declines to choose between never and no longer.

## Sources

- Documentation: https://www.discogs.com/developers
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve discogs     # run it
cauldron verify discogs -v # check every claim
```
