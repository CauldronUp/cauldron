# IPinfo

Emulates the IPinfo API (core), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**No credential is not a failure**: an absent key is a
200 with a reduced body and a wrong one is a 403.

## Sources

- Documentation: https://ipinfo.io/developers
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ipinfo     # run it
cauldron verify ipinfo -v # check every claim
```
