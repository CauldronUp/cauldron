# Brave Search

Emulates the Brave Search API (v1), for local development and tests.

**10 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Search's, which **calls a missing credential a validation
failure** -- 422 with `code: "VALIDATION"`, not a 401, so the credential is
treated as a malformed parameter and a client branching on 401 never sees it.
Brave's documentation says its rate-limit headers appear on "every API
response"; they are absent from every response reachable without a key.

## Sources

- Documentation: https://api-dashboard.search.brave.com/app/documentation/web-search/get-started
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bravesearch     # run it
cauldron verify bravesearch -v # check every claim
```
