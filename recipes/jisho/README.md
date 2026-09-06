# Jisho

Emulates the Jisho API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**There is no failure at all.** A search that finds
nothing, and a search with no keyword in it -- the one parameter the endpoint
exists for -- are both answered 200 with an empty list, exactly as a search that
works is answered 200 with a full one. A client whose query string was built
wrong gets an empty result and no reason.

**And `meta.status` is a field that can only hold one value.** It is 200 on all
three, because the only non-200 on the host is an HTML page from a path that
never reaches this code -- so a client checking `body.meta.status` is checking a
constant, which is worse than not checking, because it looks like a guard.
**Sibling objects in one array do not have the same keys**: a word's three
senses are identical in shape except that the third has a `sentences` key and
the first two do not, so mapping over them and reading `.sentences.length`
throws on two of three. `attribution` holds two booleans and a URL under one
object. A `parts_of_speech` of `["Wikipedia definition"]` says where the
definition came from rather than what kind of word it is. `jlpt` is
`["jlpt-n5"]`, the field's own name repeated inside its value, beside a `tags`
of `["wanikani7"]` -- a different company's learning app, and its level seven,
in a dictionary's public tags.

## Sources

- Documentation: https://jisho.org/forum/54fefc1f6e73340b1f160000-is-there-any-kind-of-api-for-jisho
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve jisho     # run it
cauldron verify jisho -v # check every claim
```
