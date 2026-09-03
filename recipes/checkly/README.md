# Checkly

Emulates the Checkly API (v1), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**One run disagrees with itself by location** -- one
`checkRunId` covering several results that independently differ on timing and
whether anything failed, so "did the check pass" has no single answer.

## Sources

- Documentation: https://www.checklyhq.com/docs/api
- Machine-readable description: https://api.checklyhq.com/openapi.json, last checked 2026-09-02
  `cauldron drift checkly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve checkly     # run it
cauldron verify checkly -v # check every claim
```
