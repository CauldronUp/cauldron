# Veriff

Emulates the Veriff API (v1), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

It **never looks at the signature until the id is real** --
three credential codes from the client identifier alone, each echoing what it was
given.

## Sources

- Documentation: https://devdocs.veriff.com/apidocs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve veriff     # run it
cauldron verify veriff -v # check every claim
```
