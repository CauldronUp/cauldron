# Airtable

Emulates the Airtable API (v0), for local development and tests.

**9 conformance cases, 1 checked against the live API.**

One case was struck live against api.airtable.com on 2026-09-05, and it corrected a mistake: this file had claimed a missing token gets the message "You should provide valid api key to perform this operation", which the live host does not send -- it sends "Authentication required".

A second live finding is modelled but cannot be shipped as a case. Airtable reads a token's shape before its value: a wrong token of the right shape answers type UNAUTHORIZED and "Invalid authentication token", while a value of any other length answers the same AUTHENTICATION_REQUIRED as no header at all. A conformance case exercising that first branch would have to contain a literal of Airtable's token shape, and GitHub's push protection refuses one -- including a literal whose every character is a zero. So the Recipe serves the wrong-real-key answer, which is the one an integration actually meets, and [`recipe.yaml`](recipe.yaml) states the limit: a deliberately short garbage token gets the more informative answer here than it would from Airtable.

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
