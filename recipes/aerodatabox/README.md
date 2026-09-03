# AeroDataBox

Emulates the AeroDataBox API (1.15.3.0), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**Rejected and unentitled collapse at the gateway**
-- once the host header names a real product, a garbage key and a wrong-plan key
answer alike.

## Sources

- Documentation: https://doc.aerodatabox.com/rapidapi.html
- Machine-readable description: https://doc.aerodatabox.com/docs/openapi-rapidapi-v1.json, last checked 2026-09-03
  `cauldron drift aerodatabox` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve aerodatabox     # run it
cauldron verify aerodatabox -v # check every claim
```
