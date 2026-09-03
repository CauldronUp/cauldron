# Trakt

Emulates the Trakt API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**The credential is the only thing that ever speaks**:
with no key header the edge answers instead of the application, and the version
header never distinguishes anything.

## Sources

- Documentation: https://docs.trakt.tv
- Machine-readable description: https://docs.trakt.tv/openapi/openapi.json, last checked 2026-09-03
  `cauldron drift trakt` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve trakt     # run it
cauldron verify trakt -v # check every claim
```
