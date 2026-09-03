# Scaleway

Emulates the Scaleway API (v1), for local development and tests.

**12 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**A path that does not exist answers 401.** The
credential gate runs before routing, so a typo anywhere on the authenticated
surface is indistinguishable from a permissions problem. Its catalogue shape
does *not* vary by zone, which is the answer to the question that Recipe was
started for: what varies is membership, and the smaller zone is a strict subset.

## Sources

- Documentation: https://www.scaleway.com/en/developers/api/instance/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve scaleway     # run it
cauldron verify scaleway -v # check every claim
```
