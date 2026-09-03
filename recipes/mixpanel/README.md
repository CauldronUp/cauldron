# Mixpanel

Emulates the Mixpanel API (1.0.0), for local development and tests.

**13 conformance cases, 12 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **silent acceptance is now asserted rather than
described**, and whose gate turns out to be the presence of three fields and
never their validity.

## Sources

- Documentation: https://developer.mixpanel.com/reference/ingestion-api-overview
- Machine-readable description: https://docs.mixpanel.com/openapi/ingestion.openapi.yaml, last checked 2026-09-03
  `cauldron drift mixpanel` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mixpanel     # run it
cauldron verify mixpanel -v # check every claim
```
