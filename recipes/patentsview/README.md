# PatentsView

Emulates the PatentsView API (v1), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**For an API that no longer exists** -- every host
redirects, and the one the credential lived on is gone from DNS.

## Sources

- Documentation: https://data.uspto.gov/support/transition-guide/patentsview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve patentsview     # run it
cauldron verify patentsview -v # check every claim
```
