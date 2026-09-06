# Treasury Prime

Emulates the Treasury Prime API (1), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Prime's, where **cancelling is a status you write**, not an
endpoint you call: a PATCH setting status to canceled, valid only while pending.

## Sources

- Documentation: https://docs.treasuryprime.com/docs/making-your-first-ach-transfer
- Machine-readable description: https://docs.treasuryprime.com/api-reference/openapi-payments.json, last checked 2026-09-01
  `cauldron drift treasuryprime` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve treasuryprime     # run it
cauldron verify treasuryprime -v # check every claim
```
