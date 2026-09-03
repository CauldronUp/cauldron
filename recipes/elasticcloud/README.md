# Elastic Cloud

Emulates the Elastic Cloud API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Cloud's, where **attaching any credential hides the routing**:
with none, a 404 and a 405 answer properly; with garbage, both collapse into one
undocumented shape.

## Sources

- Documentation: https://www.elastic.co/docs/api/doc/cloud/
- Machine-readable description: https://www.elastic.co/docs/api/doc/cloud.json, last checked 2026-09-01
  `cauldron drift elasticcloud` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve elasticcloud     # run it
cauldron verify elasticcloud -v # check every claim
```
