# Coveralls

Emulates the Coveralls API (v1), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**Validation runs before the credential because there
is no credential at the HTTP layer** -- the repo token lives inside the uploaded
file.

## Sources

- Documentation: https://docs.coveralls.io/api-introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve coveralls     # run it
cauldron verify coveralls -v # check every claim
```
