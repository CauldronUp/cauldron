# Clever

Emulates the Clever API (v3.0), for local development and tests.

**22 conformance cases, 7 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **Allow header is spelled one letter per header**:
nine headers named zero through eight, each holding one character.

## Sources

- Documentation: https://dev.clever.com/
- Machine-readable description: https://raw.githubusercontent.com/Clever/clever-go/master/swagger.yml, last checked 2026-09-03
  `cauldron drift clever` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clever     # run it
cauldron verify clever -v # check every claim
```
