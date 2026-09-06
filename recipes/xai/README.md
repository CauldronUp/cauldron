# xAI

Emulates the xAI API (v1), for local development and tests.

**12 conformance cases, 4 checked against the live API on 2026-08-31.**

## What this Recipe found

Missing a credential and sending the wrong one aren't the same failure, and they're not even the same shape: no credentials gets a 401 with one `code` field, a wrong key gets a 400 with a different `code` field and message, and a path nothing maps to answers with plain text instead of JSON entirely -- a client that calls `.json()` on every response because every other route answers JSON throws on that one. None of these three shapes appears anywhere in xAI's own published OpenAPI description, which documents only a single generic 400 for the whole surface.

Pricing also lives on the models endpoint itself, in a unit invented for this API alone -- USD cents per hundred million tokens -- and the running cost of a completion rides inside the response body as `usage.cost_in_usd_ticks`, an integer at ten billion ticks to the dollar.

## Sources

- Documentation: https://docs.x.ai/docs/api-reference
- Machine-readable description: https://docs.x.ai/openapi.json, last checked 2026-08-31
  `cauldron drift xai` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xai     # run it
cauldron verify xai -v # check every claim
```
