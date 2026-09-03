# Navitia

Emulates the Navitia API (v1), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**Three credential sentences and no workaround**, where an
invented coverage region answers exactly what a real one would.

## Sources

- Documentation: https://doc.navitia.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve navitia     # run it
cauldron verify navitia -v # check every claim
```
