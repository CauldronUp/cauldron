# FlightAware

Emulates the FlightAware API (4.17.1), for local development and tests.

**5 conformance cases, 2 checked against the live API on 2026-08-31.**

## What this Recipe found

**Every mistake is the same 401** -- valid endpoint,
typo, malformed id and the bare root, byte-identical including the leading
newline and eight spaces of padding.

## Sources

- Documentation: https://www.flightaware.com/aeroapi/portal/documentation
- Machine-readable description: https://www.flightaware.com/commercial/aeroapi/resources/aeroapi-openapi.yml, last checked 2026-08-31
  `cauldron drift flightaware` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve flightaware     # run it
cauldron verify flightaware -v # check every claim
```
