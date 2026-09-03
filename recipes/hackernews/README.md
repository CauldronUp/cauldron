# Hacker News

Emulates the Hacker News API (v0), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

News's, where a thing that is not there is the four bytes
`null`, with a 200. Not an empty object, not a 404, not an error field: the JSON
literal, sent as `application/json`, from both the item endpoint and the user
one. So `response.ok` is true, `.json()` succeeds, and what comes back is the
value that means "no value" in the language you are writing in -- and
`(await res.json()).title` throws a TypeError with no status anywhere to have
branched on first. **And two of the endpoints have no envelope at all**:
`/v0/maxitem.json` answers a bare integer and `/v0/topstories.json` answers a
bare array of five hundred integers, so rendering a front page is one request for
the list and five hundred for the items, with no batch endpoint anywhere.

The rest is one resource doing five jobs. `id` is a number on an item and the
username on a user. `text` is HTML inside a JSON string with no plain-text
companion, so anything rendering it strips tags and anything not rendering it
shows them. `time` is epoch **seconds**, ten digits where most of this collection
is thirteen. A story has `title`, `url`, `score` and `descendants`; a comment has
`parent` and `text` and none of those, and nothing but `type` says which arrived.
And `kids` is the direct replies while `descendants` is the total, so neither is
the length of the other.

## Sources

- Documentation: https://github.com/HackerNews/API
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hackernews     # run it
cauldron verify hackernews -v # check every claim
```
