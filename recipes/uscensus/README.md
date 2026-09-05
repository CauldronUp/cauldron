# US Census

Emulates the US Census API (2021/acs/acs1), for local development and tests.

**14 conformance cases, all of them checked against the live API, most recently on 2026-09-05.**

## What this Recipe found

**A page size here is not ignored, it is refused.** The variables catalogue answers thirteen and a half megabytes to a bare request -- all 36,428 entries -- and `?limit=2` is a **400**: "error: the 'get' argument must be a comma separated list of variable names", which is what any unrecognised parameter gets. Most providers quietly ignore a parameter they do not have, so a client that tries to page this one fails loudly rather than silently reading the same page forever. Struck live on 2026-09-05.

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
