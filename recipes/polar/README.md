# Polar

Emulates the Polar API (2026-10), for local development and tests.

**16 conformance cases, 4 checked against the live API.**

Everything about grants, benefits, and pagination still cites documentation rather than an observation, because reaching it needs an organization this Recipe cannot fabricate. The credential and routing behaviour was verified directly against api.polar.sh, unauthenticated, on 2026-09-05.

## What this Recipe found

Routing is checked before the credential, checked live: an unrouted path answers 404 and a wrong method answers 405, neither needing anything in the Authorization header -- the opposite arrangement from most bearer APIs in this collection. And an absent credential is not the same failure as a wrong one: no header at all answers `{"error":"Unauthorized","detail":"Unauthorized"}`, while a syntactically fine but fictitious bearer token answers an OAuth2 introspection shape instead, `{"error":"invalid_token","error_description":"..."}`, with no `detail` field at all. This file had modelled one message, "Invalid token.", for both.

Paying for something and receiving it are two different Polar records, and the second one can fail permanently on its own. A `BenefitGrant` carries `is_granted`, `granted_at`, `is_revoked`, `revoked_at`, and a separate `error` object -- four real states rather than two: granted, on its way, revoked, or failed for good with the customer having paid for something they will never get. The order itself still shows as succeeded; nothing about it says the grant behind it died. The default listing view does not separate these either -- `is_granted` is a filter that is off by default, so the ordinary view of a customer's entitlements mixes what they actually have with what silently failed, and counting rows counts both.

There are four kinds of bearer token (personal access, organization access, customer session, member session), and all four look identical in the request header -- using the wrong one is a 401 with nothing in the response to say which kind was expected. The collection endpoint and the item endpoint also differ by exactly one character: `/v1/benefits/` lists, `/v1/benefits/{id}` fetches one.

Granting itself is not modelled, only the record of it -- nothing here actually calls Discord or GitHub, so the fixtures simply carry one of each of the four states, with the failed one being the point worth having.

## Sources

- Documentation: https://docs.polar.sh/api-reference
- Machine-readable description: https://polar.sh/docs/openapi.json, last checked 2026-09-05
  `cauldron drift polar` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve polar     # run it
cauldron verify polar -v # check every claim
```
