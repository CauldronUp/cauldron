# Bluesky

Emulates the Bluesky API (v1), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**The handle is a pointer and the DID is the
identity**: the same profile fetched both ways comes back byte-identical.

## Sources

- Documentation: https://docs.bsky.app/docs/category/http-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bluesky     # run it
cauldron verify bluesky -v # check every claim
```
