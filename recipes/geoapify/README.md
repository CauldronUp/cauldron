# Geoapify

Emulates the Geoapify API (v1), for local development and tests.

**11 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

It **has no confidence to report** at any tier, and
whose credential check runs before routing on every path and verb tried.

## Sources

- Documentation: https://apidocs.geoapify.com/docs/ip-geolocation/
- Machine-readable description: https://raw.githubusercontent.com/geoapify/geoapify-openapi-specs/main/api-specs/ip-geolocation/ip-geolocation-api-openapi-specs.json, last checked 2026-09-02
  `cauldron drift geoapify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve geoapify     # run it
cauldron verify geoapify -v # check every claim
```
