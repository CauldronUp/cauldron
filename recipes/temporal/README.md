# Temporal

Emulates the Temporal API (v0.19.1), for local development and tests.

**15 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**A create never returns the thing it created**.
Making a namespace answers 200 -- not 201 -- with an operation in
`STATE_PENDING`, and the outcome arrives only by polling somewhere else. Nothing
here advances a fixture between requests, so a pending operation stays pending
and a client that polls until it changes never stops. No other Recipe
states that limit as sharply, though four record milder versions of it.

## Sources

- Documentation: https://docs.temporal.io/ops
- Machine-readable description: https://saas-api.tmprl.cloud/spec.json, last checked 2026-09-01
  `cauldron drift temporal` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve temporal     # run it
cauldron verify temporal -v # check every claim
```
