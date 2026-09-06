# Nutritionix

Emulates the Nutritionix API (v2), for local development and tests.

**14 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

**Half a credential is the same as none**, and
whose nutrient array carries a number and a value with no unit and no name.

## Sources

- Documentation: https://docx.syndigo.com/developers/docs/nutritionix-api-guide
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nutritionix     # run it
cauldron verify nutritionix -v # check every claim
```
