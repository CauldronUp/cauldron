# arXiv

Emulates the arXiv API (classic), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

It **answers Atom XML and nothing else** -- not with
a query parameter, not with an Accept header, not for failures. Its failure
semantics are inverted from the convention: a well-formed identifier naming
nothing is a `200` with an empty feed, while a malformed one is a `400`, so the
merely-wrong request succeeds and the merely-unlucky one fails. An unrecognised
search field is swallowed rather than refused -- `search_query=zzzzz:electron`
answers 200 with no results, and the echoed query shows the whole string became
a literal phrase under `all:`. Asking for too much is a `500` saying the server
is overloaded, not a `400`.

## Sources

- Documentation: https://info.arxiv.org/help/api/user-manual.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve arxiv     # run it
cauldron verify arxiv -v # check every claim
```
