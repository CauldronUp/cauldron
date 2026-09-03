# Pipedream

Emulates the Pipedream API (v1), for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

**Sending nothing is a worse diagnosis than sending
rubbish**: no credential answers `404 {"error": "record not found"}` while a
junk token answers a proper `401`. The record supposedly not found is the
caller's own account.

## Sources

- Documentation: https://pipedream.com/docs/rest-api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pipedream     # run it
cauldron verify pipedream -v # check every claim
```
