# Descope

Emulates the Descope API (v1), for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose two commonest error codes are undocumented.**
`E011007` for a missing bearer header and `E011008` for an unparseable one --
the two mistakes a caller meets first -- appear nowhere in Descope's own error
reference, which documents `E011001` through `E011004` and stops.

## Sources

- Documentation: https://docs.descope.com/api
- Machine-readable description: https://docs.descope.com/examples/Descope_API.yaml, last checked 2026-09-01
  `cauldron drift descope` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve descope     # run it
cauldron verify descope -v # check every claim
```
