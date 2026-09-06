# Census

Emulates the Census API (v1), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whether you meet the credential check depends on the
verb**. A GET reaches the gate and answers 401 in plain text; a DELETE on the
same path finds no route for that verb and answers 404 in JSON.

## Sources

- Documentation: https://fivetran.com/docs/activations/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve census     # run it
cauldron verify census -v # check every claim
```
