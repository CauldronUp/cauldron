# WeatherAPI

Emulates the WeatherAPI API (v1), for local development and tests.

**10 conformance cases, 9 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **own field table contradicts itself**: three fields
described as local time in Unix time, which is not a thing.

## Sources

- Documentation: https://www.weatherapi.com/docs/
- Machine-readable description: https://www.weatherapi.com/openapi.json, last checked 2026-09-05
  `cauldron drift weatherapi` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve weatherapi     # run it
cauldron verify weatherapi -v # check every claim
```
