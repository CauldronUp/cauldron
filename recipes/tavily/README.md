# Tavily

Emulates the Tavily API (unversioned), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

**Whose answer cites nothing**: no citation array, no result
index, established by walking the schema rather than the prose.

## Sources

- Documentation: https://docs.tavily.com/documentation/api-reference/endpoint/search
- Machine-readable description: https://tavily.com/.well-known/openapi.json, last checked 2026-09-05
  `cauldron drift tavily` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tavily     # run it
cauldron verify tavily -v # check every claim
```
