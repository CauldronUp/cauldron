# US Census

Emulates the US Census API (2021/acs/acs1), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

Census's, where **the key gate fires before anything is
validated.** Every request carrying a non-empty `get` answers `302` to
`missing_key.html` with `X-DataWebAPI-KeyError: 1` -- across eight datasets and
every malformed variant tried -- so a correct query and a typo are
indistinguishable. Omitting `get` entirely gets a `400` first, which means the
only way to learn anything about a query is to leave out the parameter that
makes it one. `NAME`, which appears in every documented example, is not a key
in `variables.json` at all; fetching it resolves with a `concept` string
byte-identical to `GEO_ID`'s.

## Sources

- Documentation: https://www.census.gov/data/developers/guidance/api-user-guide.Core_Concepts.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve uscensus     # run it
cauldron verify uscensus -v # check every claim
```
