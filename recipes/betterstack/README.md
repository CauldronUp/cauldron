# Better Stack

Emulates the Better Stack API (v2), for local development and tests.

**19 conformance cases, 7 checked against the live API on 2026-09-03.**

## What this Recipe found

Stack's, whose **404 always says GET** whatever method was
sent, verified across six verbs.

## Sources

- Documentation: https://betterstack.com/docs/uptime/api/getting-started-with-uptime-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve betterstack     # run it
cauldron verify betterstack -v # check every claim
```
