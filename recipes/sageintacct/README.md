# Sage Intacct

Emulates the Sage Intacct API (3.0), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Intacct's, whose gateway **cannot tell three mistakes apart** -- an
unregistered sender, mismatched XML, and a missing operation block all collapse to
one "Invalid request" -- while distinguishing three things a caller does not need.

## Sources

- Documentation: https://developer.intacct.com/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sageintacct     # run it
cauldron verify sageintacct -v # check every claim
```
