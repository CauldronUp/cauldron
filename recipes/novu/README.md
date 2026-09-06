# novu

Emulates the novu API (v1), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.novu.co, no account and no key. This file's declared "API Key not found" matched a well-formed wrong key exactly -- but no Authorization header at all gets a different sentence, "Missing authorization header", which this file had never distinguished. Split below.

## What this Recipe found

Novu creates a subscriber implicitly the first time you trigger a workflow against an id it has not seen, which means a typo can never fail outright. Send to `usr_1O4` when you meant `usr_104` and Novu answers 201, acknowledges the trigger, creates a subscriber with no email and no phone, and delivers nothing -- everything reports success, and the only visible symptom is a subscriber list that quietly grows by one every time somebody fat-fingers an id.

`acknowledged: true` on a trigger response means received, not sent -- the workflow has not run yet, no message exists, and treating the 201 as delivery reports success for notifications nobody got. The `transactionId` that comes back is a correlation key for a later search, not a resource you can fetch by id, so code that stores it expecting a handle stores something with no endpoint behind it. A subscriber's channel credentials (push tokens and the like) live in a separate array from their email and phone, keyed by provider id, so having an email does not mean having a push token -- a workflow step targeting a channel the subscriber has no credential for is skipped rather than failed.

## Sources

- Documentation: https://docs.novu.co/api-reference/overview
- Machine-readable description: https://docs.novu.co/openapi.json, last checked 2026-09-05
  `cauldron drift novu` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve novu     # run it
cauldron verify novu -v # check every claim
```
