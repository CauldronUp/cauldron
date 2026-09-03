# ShipEngine

Emulates the ShipEngine API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://www.shipengine.com/docs/rates/
- Machine-readable description: https://raw.githubusercontent.com/ShipEngine/shipengine-openapi/master/openapi.yaml, last checked 2026-08-30
  `cauldron drift shipengine` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shipengine     # run it
cauldron verify shipengine -v # check every claim
```
