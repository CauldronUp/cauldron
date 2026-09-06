# Appsmith

Emulates the Appsmith API (v1), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-09-01.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

It **says no from three pieces of software**: Spring
Security's empty 401, Appsmith's own JSON 405, and a CSRF filter that stops any
mutating verb on any path before either.

## Sources

- Documentation: https://docs.appsmith.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve appsmith     # run it
cauldron verify appsmith -v # check every claim
```
