# Tomorrow.io

Emulates the Tomorrow.io API (v4), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-09-02.**

## What this Recipe found

It **has one clock and it is not the right one** -- a
clean UTC observation time, and nothing at all about when the forecast was
computed.

## Sources

- Documentation: https://docs.tomorrow.io/reference/realtime-weather
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tomorrowio     # run it
cauldron verify tomorrowio -v # check every claim
```
