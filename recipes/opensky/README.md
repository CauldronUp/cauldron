# OpenSky

Emulates the OpenSky API (1.4.0), for local development and tests.

**11 conformance cases, 10 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **receiver field is always null for anyone anonymous**
-- the one provider in its group with a field for provenance, reserved.

## Sources

- Documentation: https://openskynetwork.github.io/opensky-api/rest.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve opensky     # run it
cauldron verify opensky -v # check every claim
```
