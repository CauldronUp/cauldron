# FireHydrant

Emulates the FireHydrant API (v1), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

It **says which credential problem you have**, in two
distinct sentences. Writing it is what finally closed the gap: a dozen Recipes
before it had discarded half of their provider's answer because the credential
check returned one bool. It now returns which way the credential failed --
absent, malformed, or presented and refused -- and a Recipe names an error per
verdict, so FireHydrant serves both sentences with nothing armed.

## Sources

- Documentation: https://docs.firehydrant.com/reference/firehydrant-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve firehydrant     # run it
cauldron verify firehydrant -v # check every claim
```
