# Statuspage

Emulates the Statuspage API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Incident status and impact are independent fields: an incident can be resolved and have had critical impact, so deriving a page's health from status alone gets it wrong in both directions. Component status isn't a boolean either -- `degraded_performance` and `partial_outage` sit between operational and major outage, and a banner that treats anything other than fully operational as "down" overstates every partial problem.

## Sources

- Documentation: https://developer.statuspage.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve statuspage     # run it
cauldron verify statuspage -v # check every claim
```
