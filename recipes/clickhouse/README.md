# ClickHouse

Emulates the ClickHouse API (v1), for local development and tests.

**15 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **API root is its own description**: `GET /v1`
answers 200 unauthenticated with 1.2MB of OpenAPI, which is why this Recipe's base
URL and its recorded description are the same string.

## Sources

- Documentation: https://clickhouse.com/docs/cloud/manage/api/api-overview
- Machine-readable description: https://api.clickhouse.cloud/v1, last checked 2026-09-01
  `cauldron drift clickhouse` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clickhouse     # run it
cauldron verify clickhouse -v # check every claim
```
