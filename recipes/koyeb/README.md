# Koyeb

Emulates the Koyeb API (1.0.0), for local development and tests.

**14 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

Its create **points at three deployments at once**:
active, latest and last-provisioned, two of which routinely name different
records with nothing saying which to poll.

## Sources

- Documentation: https://www.koyeb.com/docs/reference/api
- Machine-readable description: https://api.prod.koyeb.com/public.swagger.json, last checked 2026-09-01
  `cauldron drift koyeb` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve koyeb     # run it
cauldron verify koyeb -v # check every claim
```
