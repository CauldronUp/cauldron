# GNews

Emulates the GNews API (v4), for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

It **tells you whether your key is the right shape before
whether it is right.** Exactly 32 hexadecimal characters answers "Invalid API
Key"; 31, 33, or 32 non-hex answers "You did not provide an API key", which is
false. Truncate a key by one character and you are told you sent none.

## Sources

- Documentation: https://docs.gnews.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gnews     # run it
cauldron verify gnews -v # check every claim
```
