# Customer.io

Emulates the Customer.io API (1.0.0), for local development and tests.

**10 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose API answers its marketing site's 404** -- a 3817
byte HTML page, byte-identical across both API hosts and both kinds of mistake,
so a typo in an API path returns a web page.

## Sources

- Documentation: https://docs.customer.io/api/track/
- Machine-readable description: https://docs.customer.io/files/journeys-app.json, last checked 2026-09-01
  `cauldron drift customerio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve customerio     # run it
cauldron verify customerio -v # check every claim
```
