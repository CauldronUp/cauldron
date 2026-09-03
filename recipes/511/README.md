# 511

Emulates the 511 API (1), for local development and tests.

**6 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

**A dot in the path skips the credential entirely**: any
segment carrying a file extension is served by the web server before the
application is reached.

## Sources

- Documentation: https://511.org/open-data/transit
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve 511     # run it
cauldron verify 511 -v # check every claim
```
