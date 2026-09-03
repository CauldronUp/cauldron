# Namecheap

Emulates the Namecheap API (unversioned), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **credential includes the caller's own address**, and
where a valid address that is not the whitelisted one reads as a wrong key.

## Sources

- Documentation: https://www.namecheap.com/support/api/intro/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve namecheap     # run it
cauldron verify namecheap -v # check every claim
```
