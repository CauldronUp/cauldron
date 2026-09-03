# Edamam

Emulates the Edamam API (v1.4), for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

It gives **every nutrient a unit and the total no basis**:
a whole-recipe figure with nothing saying it needs dividing.

## Sources

- Documentation: https://developer.edamam.com/edamam-docs-nutrition-api
- Machine-readable description: https://api.edamam.com/doc/open-api/nutrition-analysis-v1.yaml, last checked 2026-09-03
  `cauldron drift edamam` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve edamam     # run it
cauldron verify edamam -v # check every claim
```
