# Hotelbeds

Emulates the Hotelbeds API (1.0), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

**The signature's quality cannot be observed**: the
key splits cleanly and every signature shape collapses into one 403.

## Sources

- Documentation: https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/
- Machine-readable description: https://bitbucket.org/ApiPortalHotelbeds/apitude-openapi/raw/master/OpenAPI-Hotel-BookingAPI-3.0.yaml, last checked 2026-09-02
  `cauldron drift hotelbeds` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hotelbeds     # run it
cauldron verify hotelbeds -v # check every claim
```
