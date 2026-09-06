# CockroachDB

Emulates the CockroachDB API (v1), for local development and tests.

**14 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **state enum has no word for deleting**, so a
cluster mid-teardown is either still created or already gone.

## Sources

- Documentation: https://www.cockroachlabs.com/docs/api/cloud/v1
- Machine-readable description: https://cockroachlabs.cloud/assets/docs/api/latest/openapi.json, last checked 2026-09-02
  `cauldron drift cockroachdb` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cockroachdb     # run it
cauldron verify cockroachdb -v # check every claim
```
