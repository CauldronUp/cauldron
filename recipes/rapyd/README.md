# Rapyd

Emulates the Rapyd API (v1), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **third failure tier is not about the credential at
all** -- a routing miss answers the same 401 whether credentials were absent,
garbage or well-formed.

## Sources

- Documentation: https://docs.rapyd.net/en/request-signatures.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rapyd     # run it
cauldron verify rapyd -v # check every claim
```
