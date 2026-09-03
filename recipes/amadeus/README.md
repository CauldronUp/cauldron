# Amadeus

Emulates the Amadeus API (v1 oauth2 / v2 flight-offers-search), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://raw.githubusercontent.com/amadeus4dev/developer-guides/master/docs/API-Keys/authorization.md
- Machine-readable description: https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/Authorizaton_v1_swagger_specification.yaml, last checked 2026-09-02
  `cauldron drift amadeus` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve amadeus     # run it
cauldron verify amadeus -v # check every claim
```
