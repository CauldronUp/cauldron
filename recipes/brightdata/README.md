# Bright Data

Emulates the Bright Data API (v3), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Data's, which **found a bug in this project's credential check**.
It groups a scheme with no value alongside a header with no scheme, where the
check called it a rejected credential; the check's own comment had been on Bright
Data's side all along.

## Sources

- Documentation: https://docs.brightdata.com/api-reference/rest-api/scraper/asynchronous-requests
- Machine-readable description: https://docs.brightdata.com/api-reference/rest-api/scraper/scraper-rest-api.json, last checked 2026-09-01
  `cauldron drift brightdata` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve brightdata     # run it
cauldron verify brightdata -v # check every claim
```
