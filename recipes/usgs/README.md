# USGS

Emulates the USGS API (fdsnws-event-1), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

Asking for a page removes the count. The earthquake
catalogue's `metadata` block changes shape with the request: ask without a limit
and it carries `count`; ask with one and `count` is gone, replaced by `limit`
and `offset`. So a client that pages and reads `metadata.count` to know when to
stop finds `undefined` -- on exactly the requests where it is paging, and only
there. **And `offset` is one-based**: a request that sent no offset at all comes
back saying `offset: 1`, so anything treating it as a zero-based cursor and
adding the page size skips a record on every page after the first.

The features themselves are GeoJSON, with GeoJSON's traps. `coordinates` is
`[longitude, latitude, depth]` -- longitude first, and the third element is
depth in kilometres -- so `[-103.4604, 8.4191, 10]` read as a latitude-longitude
pair puts the epicentre in the wrong hemisphere without erroring. `ids`,
`sources` and `types` are comma-delimited strings that **begin and end with a
comma**, so a naive split gives five elements for three identifiers. The key
`type` appears four times in one document with four vocabularies:
`FeatureCollection`, `Feature`, `Point`, `earthquake`. `mag` is a JSON number,
so a magnitude of 5.0 arrives as `5` while the title the service built beside it
says `M 5.0`. `tz`, `felt`, `cdi`, `mmi` and `alert` are present and null rather
than absent. And a request that said `format=geojson` gets its failures as
**plain text** -- a multi-line human-readable report, `text/plain;charset=UTF-8`
with no space after the semicolon -- which reads the request back with its
ampersand written as `&amp;`, an HTML entity in a plain text body.

## Sources

- Documentation: https://earthquake.usgs.gov/fdsnws/event/1/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve usgs     # run it
cauldron verify usgs -v # check every claim
```
