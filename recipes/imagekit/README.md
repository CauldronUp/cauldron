# ImageKit

Emulates the ImageKit API (v1), for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

It **has no state field because there is no state**:
transformation happens at request time on the edge. Its one asynchronous job
lives on another path and does have a handle.

## Sources

- Documentation: https://imagekit.io/docs/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve imagekit     # run it
cauldron verify imagekit -v # check every claim
```
