# fly

Emulates the fly API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://docs.machines.dev/
- Machine-readable description: https://docs.machines.dev/openapi.json, last checked 2026-08-31
  `cauldron drift fly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fly     # run it
cauldron verify fly -v # check every claim
```
