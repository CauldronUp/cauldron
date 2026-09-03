# TriNet

Emulates the TriNet API (v1), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its gateway answers from **two different subsystems** for an
empty credential and a wrong one -- different error-code namespaces, and
different capitalisation of the sentence beside them.

## Sources

- Documentation: https://www.trinet.com/customer-resource-center/technology/integration-center
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve trinet     # run it
cauldron verify trinet -v # check every claim
```
