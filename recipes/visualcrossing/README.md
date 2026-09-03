# Visual Crossing

Emulates the Visual Crossing API (rest), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

Crossing's, **the one that says what the call cost**, and the
only true 405 in its group. It advertises POST in an Allow header and crashes on
it.

## Sources

- Documentation: https://www.visualcrossing.com/resources/documentation/weather-api/timeline-weather-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve visualcrossing     # run it
cauldron verify visualcrossing -v # check every claim
```
