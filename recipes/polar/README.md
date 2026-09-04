# Polar

Emulates the Polar API (2026-10), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Paying for something and receiving it are two different Polar records, and the second one can fail permanently on its own. A `BenefitGrant` carries `is_granted`, `granted_at`, `is_revoked`, `revoked_at`, and a separate `error` object -- four real states rather than two: granted, on its way, revoked, or failed for good with the customer having paid for something they will never get. The order itself still shows as succeeded; nothing about it says the grant behind it died. The default listing view does not separate these either -- `is_granted` is a filter that is off by default, so the ordinary view of a customer's entitlements mixes what they actually have with what silently failed, and counting rows counts both.

There are four kinds of bearer token (personal access, organization access, customer session, member session), and all four look identical in the request header -- using the wrong one is a 401 with nothing in the response to say which kind was expected. The collection endpoint and the item endpoint also differ by exactly one character: `/v1/benefits/` lists, `/v1/benefits/{id}` fetches one.

Granting itself is not modelled, only the record of it -- nothing here actually calls Discord or GitHub, so the fixtures simply carry one of each of the four states, with the failed one being the point worth having.

## Sources

- Documentation: https://docs.polar.sh/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve polar     # run it
cauldron verify polar -v # check every claim
```
