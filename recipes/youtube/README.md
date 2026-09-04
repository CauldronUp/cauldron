# YouTube

Emulates the YouTube API (v3), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A search result's identifier is an object, and none of its fields is called `id`. A video's own `id` is a plain string, but a search result's `id` is a `ResourceId` -- `{kind, videoId, channelId, playlistId}` -- where only the field matching `kind` is actually populated. `result.id` gets an object, `result.id.videoId` gets `undefined` for every channel in the results, and neither one errors.

Search also returns five results by default, not twenty or fifty like most page sizes elsewhere -- small enough that a developer eyeballing a live response can conclude the search matched five things total. And the response's own description of its `items` field, generated straight into the client-library source, calls it "pagination information," a copy-paste error sitting in the exact document real client libraries are built from.

## Sources

- Documentation: https://developers.google.com/youtube/v3/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve youtube     # run it
cauldron verify youtube -v # check every claim
```
