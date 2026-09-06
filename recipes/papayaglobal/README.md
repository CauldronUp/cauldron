# Papaya Global

Emulates the Papaya Global API (v1), for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Global's, which is **not the product its name promises**: the
only reachable documented API is payment disbursement, and a beneficiary's entity
type describes the payee's own legal nature rather than who employs them.

## Sources

- Documentation: https://docs.papayaglobal.com/api-reference/beneficiaries-management
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve papayaglobal     # run it
cauldron verify papayaglobal -v # check every claim
```
