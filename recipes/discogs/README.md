# Discogs

Emulates the Discogs API (v2), for local development and tests.

**17 conformance cases, 10 checked against the live API on 2026-09-02.** The 3 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

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
