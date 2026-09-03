# Smartsheet

Emulates the Smartsheet API (2.0), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**Four routes disagree about what a missing
credential is.** No credential on a collection is 403 `errorCode 1004`; a wrong
token is 401 `errorCode 1002`; and `/sheets/{id}` answers a missing credential
with an entirely empty 401 -- zero bytes, no Content-Type, no code.

## Sources

- Documentation: https://developers.smartsheet.com
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve smartsheet     # run it
cauldron verify smartsheet -v # check every claim
```
