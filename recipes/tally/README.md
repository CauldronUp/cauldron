# Tally

Emulates the Tally API (1.0.0), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

**Every failure is the same twelve bytes** -- absent,
wrong, unrouted, wrong method, all identical -- except an undocumented root path
that answers four bytes whatever the credential.

## Sources

- Documentation: https://developers.tally.so/api-reference/introduction
- Machine-readable description: https://developers.tally.so/api-reference/openapi.json, last checked 2026-09-03
  `cauldron drift tally` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tally     # run it
cauldron verify tally -v # check every claim
```
