# Bitly

Emulates the Bitly API (v4), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

It **calls retargeting redirecting and bills for it** --
a first-class operation charged as an encode, with no field anywhere holding the
destination a link used to have.

## Sources

- Documentation: https://dev.bitly.com/api-reference/
- Machine-readable description: https://dev.bitly.com/v4/v4.json, last checked 2026-09-03
  `cauldron drift bitly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bitly     # run it
cauldron verify bitly -v # check every claim
```
