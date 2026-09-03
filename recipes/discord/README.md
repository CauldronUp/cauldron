# Discord

Emulates the Discord API (v10), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

Its snowflakes are numeric strings long enough that minting small integers
would have let a rounding bug through unnoticed.

## Sources

- Documentation: https://discord.com/developers/docs/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve discord     # run it
cauldron verify discord -v # check every claim
```
