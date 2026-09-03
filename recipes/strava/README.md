# Strava

Emulates the Strava API (v3), for local development and tests.

**10 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

It **has one answer for everything** -- four credential
failures collapse into one body, and a wrong method 404s exactly as an invented
path does.

## Sources

- Documentation: https://developers.strava.com/docs/reference/
- Machine-readable description: https://developers.strava.com/swagger/swagger.json, last checked 2026-09-02
  `cauldron drift strava` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve strava     # run it
cauldron verify strava -v # check every claim
```
