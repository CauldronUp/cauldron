# Kit

Emulates the Kit API (v4), for local development and tests.

**17 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **current API still says ConvertKit in a header** it
expects a client to parse, and whose two API versions take credentials that are
blind to each other.

## Sources

- Documentation: https://developers.kit.com/api-reference/subscribers/get-a-subscriber
- Machine-readable description: https://developers.kit.com/api-reference/v4.json, last checked 2026-09-01
  `cauldron drift kit` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kit     # run it
cauldron verify kit -v # check every claim
```
