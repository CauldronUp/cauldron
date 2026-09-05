# iNaturalist

Emulates the iNaturalist API (v1), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**One field carries two schemes and the only
broken one is the record you asked for.** A single taxon response holds
forty-five `wikipedia_url` values -- the taxon, its ancestors, its children.
Twenty-one begin `http://` and twenty-four begin `https://`, same host, same
shape, one document. And exactly one of the forty-five contains a raw space:
`https://en.wikipedia.org/wiki/Apis mellifera`, the top-level record, the only
one a client reading `results[0]` will ever touch. Every nested one is escaped
with underscores. An `<a href>` works, because browsers repair it, so the bug
survives every manual check and fails where nobody is looking.

**And every fetch is a search.** `/v1/taxa/47219` answers with a paging
envelope, and a thing that does not exist answers 200 with `total_results: 0`
rather than 404 -- so `res.results[0].name` throws on undefined instead of
reading a status. The two disagree about the page size, too: the fetch that
found something reports `per_page: 30` and the one that found nothing reports
`per_page: 1`. The one real 404 on the host is Express's own HTML page, `Cannot
GET /v1/nosuchthing` inside a `<pre>`. `wikipedia_summary` is HTML in a field
called summary, with an en dash in its range of species counts. The ancestry is
sent twice in two types, as `ancestor_ids` and as a slash-joined string. And
`conservation_status` is `null` beside a `conservation_statuses` holding
fourteen.

## Sources

- Documentation: https://api.inaturalist.org/v1/docs/
- Machine-readable description: https://api.inaturalist.org/swagger.json, last checked 2026-09-05
  `cauldron drift inaturalist` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve inaturalist     # run it
cauldron verify inaturalist -v # check every claim
```
