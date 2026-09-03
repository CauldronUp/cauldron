# Sendcloud

Emulates the Sendcloud API (v2), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-01.**

## Sources

- Documentation: https://sendcloud.dev/api/v2/parcels/index.md
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sendcloud     # run it
cauldron verify sendcloud -v # check every claim
```
