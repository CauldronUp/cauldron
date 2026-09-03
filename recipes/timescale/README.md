# Timescale

Emulates the Timescale API (v1), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Found after the search for it. `console.timescale.com` is
dead DNS and `api.timescale.com` is **a load balancer that refuses TLS** -- worse
than dead, because a client will keep retrying it.

## Sources

- Documentation: https://www.tigerdata.com/docs/reference/tiger-cloud-rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve timescale     # run it
cauldron verify timescale -v # check every claim
```
