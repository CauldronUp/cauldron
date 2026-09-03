# Apollo.io

Emulates the Apollo.io API (v1), for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-02.**

## What this Recipe found

**The absent credential is a validation failure**:
422 for nothing sent, 401 for something wrong.

## Sources

- Documentation: https://docs.apollo.io/reference/people-enrichment
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apolloio     # run it
cauldron verify apolloio -v # check every claim
```
