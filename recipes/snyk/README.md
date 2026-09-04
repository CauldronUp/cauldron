# Snyk

Emulates the Snyk API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A vulnerability with no fix is the interesting one: severity says how bad it is, but `isUpgradable` and `isPatchable` say whether anything can be done about it today, and a gate that fails on severity alone blocks every build until an upstream maintainer ships a release. An ignored issue is also still returned -- Snyk records the ignore as a decision rather than removing the finding, so a naive count re-reports something a human already assessed.

## Sources

- Documentation: https://docs.snyk.io/snyk-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve snyk     # run it
cauldron verify snyk -v # check every claim
```
