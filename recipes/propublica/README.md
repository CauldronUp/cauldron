# ProPublica

Emulates the ProPublica API (v2), for local development and tests.

**13 conformance cases, 12 checked against the live API on 2026-09-02.**

## What this Recipe found

**Two sentinel identifiers answer with a
placeholder**: zero and all-nines return a real Unknown Organization record
where every other unassigned number 404s.

## Sources

- Documentation: https://projects.propublica.org/nonprofits/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve propublica     # run it
cauldron verify propublica -v # check every claim
```
