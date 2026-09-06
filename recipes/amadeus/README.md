# Amadeus

Emulates the Amadeus API (v1 oauth2 / v2 flight-offers-search), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Neither documented Amadeus API host resolves at all -- test.api.amadeus.com and api.amadeus.com return no DNS record from any vantage point this Recipe could check, so every case here is read from Amadeus's own OpenAPI descriptions and developer guides on GitHub rather than observed.

The clearest trap: the OAuth token endpoint's errors are a flat four-key object ({error, error_description, code, title}), while every other Amadeus error on every other endpoint is a JSON:API-flavoured array ({errors:[{status, code, title, detail, source}]}). A client that learns Amadeus's error shape from the token endpoint -- the first call any integration makes -- learns the wrong shape for everything after it. Separately, a flight offer carries no expiry of its own: lastTicketingDate is a deadline for an offer already booked, not a statement of when an unbooked offer goes stale, and Amadeus's own published material never states what a stale or already-booked offer is refused with.

Amadeus's own docs also glue two different token failures to the identical code and title (38187, "Invalid parameters") for both a wrong grant type and wrong credentials -- distinguishable only by HTTP status or the `error` string, not by `code` alone.

## Sources

- Documentation: https://raw.githubusercontent.com/amadeus4dev/developer-guides/master/docs/API-Keys/authorization.md
- Machine-readable description: https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/Authorizaton_v1_swagger_specification.yaml, last checked 2026-09-02
  `cauldron drift amadeus` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve amadeus     # run it
cauldron verify amadeus -v # check every claim
```
