# Hunter.io

Emulates the Hunter.io API (v2), for local development and tests.

**19 conformance cases, 3 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **two limits are the wrong way round** -- 403 is the
rate limit and 429 is the plan quota.


**And Hunter's other published document disagrees with that one.**

`cauldron drift` reported this Recipe's fingerprint moved on 2026-09-09 and named the claim that moved with it: *no operation the Recipe routes to answers 403, which rate_limit declares*.

The OpenAPI document at `hunter.io/openapi.json` no longer declares 403 on any of the three routes modelled here. Domain Search now declares 200, 400, 401, 429 and a default; Email Finder adds 404; Email Verifier adds 202, 222 and 451. **403 is gone from all three and 429 has taken its place.**

The prose documentation at `hunter.io/api-documentation/v2`, checked the same day, still says word for word: "403 - Forbidden You have reached the rate limit." and "429 - Too many requests You have reached your usage limit. Upgrade your plan if necessary."

So one provider publishes two descriptions of its own API that disagree about which status a rate limit is, and a developer's answer depends on which document they happened to open. This Recipe serves the split unchanged, because the table it was quoted from still says it.

The OpenAPI's half of the change is not purely an improvement either. Every failure it declares on these routes — 400, 401, 429 and the default alike — now carries the identical description "The request could not be completed.", and neither specific sentence survives anywhere in the document. The status got more correct and the description got less useful in the same republish.

## Sources

- Documentation: https://hunter.io/api-documentation/v2
- Machine-readable description: https://hunter.io/openapi.json, last checked 2026-09-05
  `cauldron drift hunterio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hunterio     # run it
cauldron verify hunterio -v # check every claim
```
