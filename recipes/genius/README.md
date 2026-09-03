# Genius

Emulates the Genius API (v1), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

**Two pieces of software refuse the caller** and only
one of them agrees with the other about whether routing matters.

## Sources

- Documentation: https://docs.genius.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve genius     # run it
cauldron verify genius -v # check every claim
```
