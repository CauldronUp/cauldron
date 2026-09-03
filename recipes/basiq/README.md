# Basiq

Emulates the Basiq API (v3), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**The connection is fine and the account is not** -- two
accounts under one active connection, one of them two years stale, with nothing
on the connection saying so.

## Sources

- Documentation: https://api.basiq.io/reference/getconnection
- Machine-readable description: https://raw.githubusercontent.com/basiqio-oss/Basiq-docs/v3.0/reference/connect.json, last checked 2026-09-02
  `cauldron drift basiq` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve basiq     # run it
cauldron verify basiq -v # check every claim
```
