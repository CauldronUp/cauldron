# GitBook

Emulates the GitBook API (v1), for local development and tests.

**11 conformance cases, 1 checked against the live API on 2026-09-01.**

## What this Recipe found

It **checks whether a route exists before who you are** --
the opposite of Coda, and the reason only one of its cases is verified: the
runtime checks credentials first, so GitBook's missing-credential 403 is
unreachable from here.

## Sources

- Documentation: https://gitbook.com/docs/developers/gitbook-api
- Machine-readable description: https://api.gitbook.com/openapi.json, last checked 2026-09-01
  `cauldron drift gitbook` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gitbook     # run it
cauldron verify gitbook -v # check every claim
```
