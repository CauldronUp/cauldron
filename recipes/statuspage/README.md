# Statuspage

Emulates the Statuspage API (v1), for local development and tests.

**13 conformance cases, 2 checked against the live API on 2026-09-05.**

Statuspage has no sandbox, so the incident and component cases still cite documentation. The credential shape needed no page of its own, and checking it live found this Recipe's own error model wrong.

## What this Recipe found

Incident status and impact are independent fields: an incident can be resolved and have had critical impact, so deriving a page's health from status alone gets it wrong in both directions. Component status isn't a boolean either -- `degraded_performance` and `partial_outage` sit between operational and major outage, and a banner that treats anything other than fully operational as "down" overstates every partial problem.

## What checking it live found

The pairing reads backwards from the usual one: no `Authorization` header at all answers `"Invalid authentication token."` -- the more specific-sounding sentence -- and a present, wrong OAuth token answers the vaguer `"Could not authenticate"`. This Recipe had one message for both; `auth.absent_error` keeps them apart now. A path nothing declares and a wrong method both get the absent-credential sentence too, so the default check-before-routing order already matched and needed no change.

## Sources

- Documentation: https://developer.statuspage.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve statuspage     # run it
cauldron verify statuspage -v # check every claim
```
