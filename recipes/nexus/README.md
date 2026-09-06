# Sonatype Nexus

Emulates the Sonatype Nexus API (3), for local development and tests.

**13 conformance cases, 10 checked against the live API on 2026-09-01.**

## What this Recipe found

Its `x-siesta-faultid` was **identical across three different
requests**, so the field an operator would grep for names a class of requests
rather than theirs.

## Sources

- Documentation: https://help.sonatype.com/en/rest-and-integration-api.html
- Machine-readable description: https://repo.eclipse.org/service/rest/swagger.json, last checked 2026-09-01
  `cauldron drift nexus` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nexus     # run it
cauldron verify nexus -v # check every claim
```
