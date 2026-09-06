# OpenStates

Emulates the OpenStates API (v3), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

It **admits its own identifiers are not stable** -- a
field listing the identifiers it used to give the same bill.

## Sources

- Documentation: https://v3.openstates.org/docs
- Machine-readable description: https://v3.openstates.org/openapi.json, last checked 2026-09-02
  `cauldron drift openstates` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openstates     # run it
cauldron verify openstates -v # check every claim
```
