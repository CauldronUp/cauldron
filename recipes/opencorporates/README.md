# OpenCorporates

Emulates the OpenCorporates API (v0.4), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

**An aggregator that has closed**, whose own
documentation says 401 means incorrect credentials and whose 401 fires for
absent ones too.

## Sources

- Documentation: https://api.opencorporates.com/documentation/API-Reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve opencorporates     # run it
cauldron verify opencorporates -v # check every claim
```
