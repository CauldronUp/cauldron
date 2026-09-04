# ShipEngine

Emulates the ShipEngine API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A rate request can half-fail with a 200. `rate_response.rates` holds the carriers that answered and `rate_response.errors` holds the ones that refused, side by side in the same body, and `rate_response.status` says "completed" regardless -- it describes whether ShipEngine finished asking, not whether the answers are complete. A client that checks the status code and reads `rates` sees a shorter list than it asked for, with no sign anything went wrong: a customer shown two options instead of three, missing the cheapest one, and nobody notices until someone checks a carrier's own site.

## Sources

- Documentation: https://www.shipengine.com/docs/rates/
- Machine-readable description: https://raw.githubusercontent.com/ShipEngine/shipengine-openapi/master/openapi.yaml, last checked 2026-08-30
  `cauldron drift shipengine` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shipengine     # run it
cauldron verify shipengine -v # check every claim
```
