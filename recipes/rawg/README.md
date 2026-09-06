# RAWG

Emulates the RAWG API (v1), for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-08-31.**

## What this Recipe found

Its **two auth failures cannot both be served.** It does say
which credential problem happened -- "The key parameter is not provided" against
"The API key is not found" -- and one authorisation gate can carry one of them.
Six Recipes have now hit that limit.

## Sources

- Documentation: https://api.rawg.io/docs/
- Machine-readable description: https://api.rawg.io/docs/?format=openapi, last checked 2026-08-31
  `cauldron drift rawg` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rawg     # run it
cauldron verify rawg -v # check every claim
```
