# PropelAuth

Emulates the PropelAuth API (backend/v1), for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **status code tracks whether you sent a header, not
what was in it** -- 404 for none, 401 for any, byte-identical empty bodies, on a
hostname invented for the probe and on its own documented one alike.

## Sources

- Documentation: https://docs.propelauth.com/reference/api/getting-started
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve propelauth     # run it
cauldron verify propelauth -v # check every claim
```
