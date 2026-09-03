# Telnyx

Emulates the Telnyx API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developers.telnyx.com/api
- Machine-readable description: https://telnyx.com/openapi.json, last checked 2026-08-31
  `cauldron drift telnyx` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve telnyx     # run it
cauldron verify telnyx -v # check every claim
```
