# iTunes Search

Emulates the iTunes Search API (itunes), for local development and tests.

**12 conformance cases, 9 checked against the live API on 2026-08-30.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**The JSON is served as JavaScript.**
Every response carries `Content-Type: text/javascript; charset=utf-8`, success
and failure alike. A browser sniffs its way to the right parser and
`fetch(...).json()` works, so nobody notices -- until the request goes through
anything that checks the content type first, which then refuses a body it could
have read, or hands it to a JavaScript evaluator because that is what the header
asked for.

**And there is no such thing as no results.** `?term=zzzzqqqxyznothing` answers
`resultCount: 60`, every one an audiobook by "ZZZ Sleep": the search decomposes
the term and finds something for the part of it that is pronounceable, so a page
reporting sixty results for a word nobody typed on purpose is showing sleep
recordings. A band is not the first result for its own name either --
`?term=radiohead` answers a feature film, `kind: "feature-movie"`, and Radiohead
appears only once the caller asks for `entity=musicArtist`. The array is
heterogeneous, a track and an artist sharing almost no keys and told apart only
by `wrapperType`; `amgArtistId` is a foreign key into All Media Guide, renamed
in 2013; and an id nobody has answers the same `{"resultCount": 0, "results":
[]}` as a search with no term at all, so a rejected request and an empty one are
one answer.

## Sources

- Documentation: https://performance-partners.apple.com/search-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve itunes     # run it
cauldron verify itunes -v # check every claim
```
