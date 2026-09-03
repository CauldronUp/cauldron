# Nuvei

Emulates the Nuvei API (REST 1.0), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**A wrong signature and a fake merchant are one event**,
because merchant identity is checked before the checksum ever is.

## Sources

- Documentation: https://docs.nuvei.com/documentation/accept-payment/server-to-server/rest-1-0/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nuvei     # run it
cauldron verify nuvei -v # check every claim
```
