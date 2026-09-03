# xAI

Emulates the xAI API (v1), for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-08-31.**

## Sources

- Documentation: https://docs.x.ai/docs/api-reference
- Machine-readable description: https://docs.x.ai/openapi.json, last checked 2026-08-31
  `cauldron drift xai` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xai     # run it
cauldron verify xai -v # check every claim
```
