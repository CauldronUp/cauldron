# UK Police

Emulates the UK Police API (ukpolice), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Police's, where **`url` means three different things.** A force
sends `"url": "http://www.leics.police.uk/"`, a crime category sends `"url":
"all-crime"`, and the telephone engagement method sends `"url": ""`. A web
address, a slug, and nothing -- so a client rendering `url` as a link produces a
working link, a broken relative one, and an empty href, depending on which
endpoint it came from.

**And every link the force owns is plain HTTP.** Its website, its Facebook page,
its Twitter and its YouTube are all `http://`, in a payload delivered over TLS;
the one exception is the RSS feed, which is `https://`, so the API is not
consistently either. Three content types answer three kinds of request, and only
one is JSON: a force that does not exist is `text/plain` carrying the nine bytes
`Not Found`, and a path that does not exist is `text/html` carrying an HTML
fragment with no doctype around it. `type` and `title` are the same word on every
engagement method. And `stop-and-search` is a hyphenated JSON key, which a dot
cannot reach.

## Sources

- Documentation: https://data.police.uk/docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ukpolice     # run it
cauldron verify ukpolice -v # check every claim
```
