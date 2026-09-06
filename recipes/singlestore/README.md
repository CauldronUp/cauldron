# SingleStore

Emulates the SingleStore API (v2), for local development and tests.

**17 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose create hands back a password and nothing else
useful** -- `adminPassword`, `clusterID`, `groupID`, no state and no operation
to poll. A wrong method answers a real 405 with an `Allow` header, the direct
opposite of Turso, where no 405 exists at all.

## Sources

- Documentation: https://docs.singlestore.com/cloud/reference/management-api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve singlestore     # run it
cauldron verify singlestore -v # check every claim
```
