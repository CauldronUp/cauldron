# OpenAIRE

Emulates the OpenAIRE API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The JSON is XML wearing a costume.** Every text
value sits under `"$"`, every attribute under a key beginning `"@"`, and the
element names keep their XML namespaces: reading a title means
`result.metadata["oaf:entity"]["oaf:result"].title[0]["$"]`, because `oaf:entity`
is not an identifier in any language a client is written in, and neither is `$`.
It is not a JSON API that happens to look odd -- it is an XML document run
through a mechanical transliteration, and the shape of the original shows
through every field.

**And `total` is a number when it is not zero and a string when it is.** One
match answers `"total": {"$": 1}`; no match answers `"total": {"$": "0"}` with
`"results": null`, so the two paths through the response share neither type. A
page size out of range puts the useful sentence in the wrong field: `message` is
`"400 - Illegal argument exception."`, restating the status a client already
has, while `exception` holds the only text that says what happened. And a path
that does not exist answers **Apache Tomcat's own error report, naming its
version -- 7.0.68, from 2016 -- twice**, once in the title and once at the foot,
with a `message` field of `<u></u>`.

## Sources

- Documentation: https://graph.openaire.eu/docs/apis/search-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openaire     # run it
cauldron verify openaire -v # check every claim
```
