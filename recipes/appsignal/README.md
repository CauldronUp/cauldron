# AppSignal

Emulates the AppSignal API (1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

It **draws no line between an occurrence and a group**
-- the abstraction Rollbar, Bugsnag and Raygun all treat as central is simply
absent from its public REST API.

## Sources

- Documentation: https://docs.appsignal.com/api.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve appsignal     # run it
cauldron verify appsignal -v # check every claim
```
