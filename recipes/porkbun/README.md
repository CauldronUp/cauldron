# Porkbun

Emulates the Porkbun API (v3), for local development and tests.

**11 conformance cases, 10 checked against the live API on 2026-09-03.**

## What this Recipe found

It **takes its credential in a header after all**, against
every example it publishes, and answers 400 to every credential failure.

## Sources

- Documentation: https://porkbun.com/api/json/v3/documentation
- Machine-readable description: https://porkbun.com/api/json/v3/spec, last checked 2026-09-03
  `cauldron drift porkbun` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve porkbun     # run it
cauldron verify porkbun -v # check every claim
```
