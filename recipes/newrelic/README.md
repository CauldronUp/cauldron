# New Relic

Emulates the New Relic API (v2), for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-08-31.**

## What this Recipe found

Relic's, **whose own description would generate a broken client.**
The spec declares the applications *list* with the same schema as the single
fetch, so it promises `{"application": {...}}` where the wire sends
`{"applications": [...]}`. Its wrong-key failure also echoes the rejected
credential back in the body.

## Sources

- Documentation: https://docs.newrelic.com/docs/apis/rest-api-v2/get-started/introduction-new-relic-rest-api-v2/
- Machine-readable description: https://api.newrelic.com/docs/swagger.yml, last checked 2026-08-31
  `cauldron drift newrelic` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve newrelic     # run it
cauldron verify newrelic -v # check every claim
```
