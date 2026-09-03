# Baserow

Emulates the Baserow API (2.3.3), for local development and tests.

**15 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

It **sends field_1234 unless you ask for names**, which is
the answer Airtable's Recipe here records the opposite of.

## Sources

- Documentation: https://baserow.io/docs/apis%2Frest-api
- Machine-readable description: https://api.baserow.io/api/schema.json, last checked 2026-09-01
  `cauldron drift baserow` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve baserow     # run it
cauldron verify baserow -v # check every claim
```
