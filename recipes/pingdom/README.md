# Pingdom

Emulates the Pingdom API (3.1), for local development and tests.

**14 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

It **separates the last poll from the last outage**, the
only provider in its group to do so, and the only one with a state for
unconfirmed.

## Sources

- Documentation: https://docs.pingdom.com/api/
- Machine-readable description: https://docs.pingdom.com/API_3.1.yaml, last checked 2026-09-03
  `cauldron drift pingdom` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pingdom     # run it
cauldron verify pingdom -v # check every claim
```
