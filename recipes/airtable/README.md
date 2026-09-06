# Airtable

Emulates the Airtable API (v0), for local development and tests.

**12 conformance cases, 2 checked against the live API.**

Two cases were struck live against api.airtable.com on 2026-09-05, and the first corrected a mistake: this file had claimed a missing token gets the message "You should provide valid api key to perform this operation", which the live host does not send -- it sends "Authentication required".

Airtable also reads a token's shape before its value, and does not tell a wrongly shaped one from a missing one: both answer `AUTHENTICATION_REQUIRED`, while a well-formed token it does not know answers `UNAUTHORIZED` and "Invalid authentication token". Both branches are served. The one thing this Recipe cannot do is ship a well-formed token to prove the second with: GitHub's push protection refuses any literal of Airtable's token shape, including one whose every character is a zero, so [`recipe.yaml`](recipe.yaml) writes the real rule as its shape and exempts the single fixture key it publishes, in the open.

## What this Recipe found

A record's own fields live nested under "fields" -- only id and createdTime sit at the top level, so code that reads record.Name directly finds nothing. Base and table are path segments that genuinely partition the data, which is easy for a hand-rolled fake to get wrong silently by leaking records across tables.

The one fidelity gap worth knowing: Airtable's createdTime carries millisecond precision, Cauldron's only second precision. Cases here are read from Airtable's Web API reference rather than observed against a live base.

## Sources

- Documentation: https://airtable.com/developers/web/api/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve airtable     # run it
cauldron verify airtable -v # check every claim
```
