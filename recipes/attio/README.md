# Attio

Emulates the Attio API (2.0.0), for local development and tests.

**6 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

This Recipe was excluded as impossible, and the assessment was simply wrong.
Attio is REST and publishes an OpenAPI description; it had been recorded as
GraphQL and nobody reread the reason.

Both entries on that list sat unchallenged for the same reason: nothing rereads
a list of things that cannot be done, and an expired reason reads exactly like a
live one.

## Sources

- Documentation: https://api.attio.com/openapi/api
- Machine-readable description: https://attio.com/openapi.json, last checked 2026-09-02
  `cauldron drift attio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve attio     # run it
cauldron verify attio -v # check every claim
```
