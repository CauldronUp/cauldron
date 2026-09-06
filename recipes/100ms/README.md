# 100ms

Emulates the 100ms API (v2), for local development and tests.

**17 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **two deepest credential failures differ only in
latency** -- 221 milliseconds for a well-shaped fake key against under two for
everything shallower, which is a database lookup visible from outside.

## Sources

- Documentation: https://www.100ms.live/docs/server-side/v2/introduction/basics
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve 100ms     # run it
cauldron verify 100ms -v # check every claim
```
