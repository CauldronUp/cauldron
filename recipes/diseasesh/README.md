# disease.sh

Emulates the disease.sh API (diseasesh), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Six fields end in `PerOneMillion` and three of
them are integers. `casesPerOneMillion`, `deathsPerOneMillion` and
`testsPerOneMillion` are whole numbers; `activePerOneMillion`,
`recoveredPerOneMillion` and `criticalPerOneMillion` carry two decimal places.
One suffix, one obvious meaning, and two JSON types split by which statistic it
happens to be -- so a client that types the family as an integer truncates three
of them, and nothing in the response distinguishes the two groups.

**And `countryInfo` carries `_id`.** A field named the way a database names its
primary key, underscore and all, on a public API -- holding 124, which is the UN
numeric country code and not an internal row id at all. `updated` is a
JavaScript millisecond timestamp as a bare integer, and it changes on every
request while the figures beside it do not. The coordinates are whole degrees on
a record whose derived statistics carry two decimals, so the least precise
numbers here are the ones describing a place -- and one of them is called `long`
rather than `longitude`. A country that does not exist answers JSON and a path
that does not exist answers an Express HTML page, both with 404.

## Sources

- Documentation: https://disease.sh/docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve diseasesh     # run it
cauldron verify diseasesh -v # check every claim
```
