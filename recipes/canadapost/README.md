# Canada Post

Emulates the Canada Post API (unversioned), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

Post's, which **checks the parcel's shape before the
caller's**: twelve, thirteen or sixteen alphanumerics, validated ahead of the
credential and even against a wrong one.

## Sources

- Documentation: https://www.canadapost-postescanada.ca/track-reperage/en
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve canadapost     # run it
cauldron verify canadapost -v # check every claim
```
