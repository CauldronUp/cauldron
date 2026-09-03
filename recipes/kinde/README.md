# Kinde

Emulates the Kinde API (1), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **management API is rich and whose token endpoint is ten
bytes**. Three distinguishable credential sentences on one, and `not_found` in
plain text on the other.

## Sources

- Documentation: https://docs.kinde.com/developer-tools/kinde-api/
- Machine-readable description: https://api-spec.kinde.com/kinde-management-api-spec.yaml, last checked 2026-09-01
  `cauldron drift kinde` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kinde     # run it
cauldron verify kinde -v # check every claim
```
