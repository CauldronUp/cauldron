# Hunter.io

Emulates the Hunter.io API (v2), for local development and tests.

**16 conformance cases, 3 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **two limits are the wrong way round** -- 403 is the
rate limit and 429 is the plan quota.

## Sources

- Documentation: https://hunter.io/api-documentation/v2
- Machine-readable description: https://hunter.io/openapi.json, last checked 2026-09-05
  `cauldron drift hunterio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hunterio     # run it
cauldron verify hunterio -v # check every claim
```
