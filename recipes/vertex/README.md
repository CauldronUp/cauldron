# Vertex

Emulates the Vertex API (v2), for local development and tests.

**3 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developer.vertexinc.com/oseries
- Machine-readable description: https://dash.readme.com/api/v1/api-registry/1bto7bwmt0de2rl, last checked 2026-09-01
  `cauldron drift vertex` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vertex     # run it
cauldron verify vertex -v # check every claim
```
