# YouTube

Emulates the YouTube API (v3), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

The search-result and video cases still cite the discovery document. The credential shape needed no project at all, and checking it live found this Recipe's status and sentence both wrong.

## What this Recipe found

A search result's identifier is an object, and none of its fields is called `id`. A video's own `id` is a plain string, but a search result's `id` is a `ResourceId` -- `{kind, videoId, channelId, playlistId}` -- where only the field matching `kind` is actually populated. `result.id` gets an object, `result.id.videoId` gets `undefined` for every channel in the results, and neither one errors.

Search also returns five results by default, not twenty or fifty like most page sizes elsewhere -- small enough that a developer eyeballing a live response can conclude the search matched five things total. And the response's own description of its `items` field, generated straight into the client-library source, calls it "pagination information," a copy-paste error sitting in the exact document real client libraries are built from.

## What checking it live found

No key at all is 403 `PERMISSION_DENIED`, `"Method doesn't allow unregistered callers..."`; a present, wrong key -- including an existing case's own fixture key -- is 400 `INVALID_ARGUMENT`, `"API key not valid."`, never the 401 `UNAUTHENTICATED` this Recipe had claimed for both. Google's older `error.errors[0].reason` field carries a third value again (`forbidden` / `badRequest`) beside the newer `status`, and several client libraries switch on it. A path nothing declares is empty and `text/html`, resolved before either credential question is asked.

## Sources

- Documentation: https://developers.google.com/youtube/v3/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve youtube     # run it
cauldron verify youtube -v # check every claim
```
