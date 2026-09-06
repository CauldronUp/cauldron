# Livepeer

Emulates the Livepeer API (v1), for local development and tests.

**11 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its asset and task **disagree about being done** -- two
status vocabularies sharing two words, where the task's `completed` says nothing
about the asset's `ready`.

## Sources

- Documentation: https://docs.livepeer.org/api-reference/asset/overview
- Machine-readable description: https://raw.githubusercontent.com/livepeer/studio/master/packages/api/src/schema/api-schema.yaml, last checked 2026-09-01
  `cauldron drift livepeer` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve livepeer     # run it
cauldron verify livepeer -v # check every claim
```
