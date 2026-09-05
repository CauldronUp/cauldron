# Make

Emulates the Make API (v2), for local development and tests.

**6 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**A token from the wrong region looks like a token that
never existed** -- Make is partitioned by hostname, and using a real credential
on the wrong one produces the same "Not authorized." that a fabricated one does.

## Sources

- Documentation: https://developers.make.com/api-documentation/api-reference/scenarios
- Machine-readable description: https://eu1.make.com/api/v1/openapi.json, last checked 2026-09-05
  `cauldron drift make` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve make     # run it
cauldron verify make -v # check every claim
```
