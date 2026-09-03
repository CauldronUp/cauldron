# ScrapingBee

Emulates the ScrapingBee API (v1), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**The target's status survives in a header until
you ask for it not to**: `Spb-initial-status-code` carries it, and the
`transparent_status_code` parameter makes the proxy's own status *become* the
target's, destroying the only signal that told them apart.

## Sources

- Documentation: https://www.scrapingbee.com/documentation/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve scrapingbee     # run it
cauldron verify scrapingbee -v # check every claim
```
