# Hardcover

Emulates the Hardcover API (v1), for local development and tests.

**8 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

It **separates the work from the edition and says so**,
in prose and structurally, with the ISBN fields only on editions.

## Sources

- Documentation: https://docs.hardcover.app/api/getting-started/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hardcover     # run it
cauldron verify hardcover -v # check every claim
```
