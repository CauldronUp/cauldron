# Wise

Emulates the Wise API (v1), for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-08-31.**

## What this Recipe found

**One endpoint answers without a credential and the rest
do not** -- `/v1/quotes` computes a real quote for anybody. It also changes the
order it checks routing and authentication depending on *which* credential
problem you have.

## Sources

- Documentation: https://docs.wise.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wise     # run it
cauldron verify wise -v # check every claim
```
