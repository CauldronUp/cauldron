# Airtable

Emulates the Airtable API (v0), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

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
