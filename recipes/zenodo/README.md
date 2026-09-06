# Zenodo

Emulates the Zenodo API (zenodo), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**The same identifier appears four times in four
forms.** `id` is the number `22173677`, `recid` is `"22173677"` as a string,
`doi` is `"10.5281/zenodo.22173677"`, and `doi_url` is that with a resolver in
front. A fifth, `conceptrecid`, is `"22173676"` -- one less, as a string -- and
`conceptdoi` is that one prefixed, so eight of the record's twenty top-level
keys are the same two numbers written out differently.

**And the record says it is finished three times, in three vocabularies:**
`"status": "published"`, `"state": "done"`, `"submitted": true`, with nothing
saying which is authoritative. `modified` and `updated` are byte-identical to
the microsecond. The array is called `hits` and lives inside an object called
`hits`, so reaching the records is `response.hits.hits`. And two failures use
different words for the same status -- a record that does not exist is about a
persistent identifier, and a path that does not exist is the web framework
apologising in case you typed the URL by hand.

## Sources

- Documentation: https://developers.zenodo.org/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zenodo     # run it
cauldron verify zenodo -v # check every claim
```
