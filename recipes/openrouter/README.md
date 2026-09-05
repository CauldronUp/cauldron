# openrouter

Emulates the openrouter API (v1), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-08-23.**

## What this Recipe found

Six of its cases are the model-catalogue ones, whose numbers were read from the
provider rather than inferred. Its completion cases carry no verification date,
because calling that endpoint costs money -- an honest reason for an unverified
claim, and one worth stating rather than leaving as a blank field.

## Sources

- Documentation: https://openrouter.ai/docs/api-reference/overview
- Machine-readable description: https://openrouter.ai/openapi.json, last checked 2026-09-05
  `cauldron drift openrouter` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openrouter     # run it
cauldron verify openrouter -v # check every claim
```
