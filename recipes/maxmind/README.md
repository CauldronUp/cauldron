# MaxMind

Emulates the MaxMind API (v2.1), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **credential stands in front of the question** -- a
reserved address never reaches the reserved-address error its own documentation
declares.

## Sources

- Documentation: https://dev.maxmind.com/geoip/docs/web-services?lang=en
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve maxmind     # run it
cauldron verify maxmind -v # check every claim
```
