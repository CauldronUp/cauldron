# Zippopotam

Emulates the Zippopotam API (zippopotam), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

The keys have spaces in them. `"country
abbreviation"`, `"post code"`, `"place name"`, `"state abbreviation"` -- not
camelCase, not snake_case, not kebab-case: a space, in the key. Nothing in the
response can be destructured, nothing can be reached with a dot, and every access
is a bracket with a quoted string in it. **And `"CA"` is two different things in
one API**: the country abbreviation for Canada and the state abbreviation for
California, in fields whose names differ by one word.

`latitude` and `longitude` are strings, and with a minus sign in front the
lexical order is not close to the numeric one. A postcode that does not exist is
`{}` with a 404 -- an empty object, so `.json()` succeeds and every field is
undefined -- and a country code that does not exist answers exactly the same way,
as does a real postcode under the wrong country. `places` holds objects of **two
different shapes** depending on which endpoint answered: looking up a postcode
gives entries with `state` and `state abbreviation`, looking up a place gives
entries with `post code` instead, under the same key with nothing saying which
arrived. And the reverse lookup repeats `"place name"` at both levels, on the
request that named it in the first place.

## Sources

- Documentation: https://api.zippopotam.us/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zippopotam     # run it
cauldron verify zippopotam -v # check every claim
```
