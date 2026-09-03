# OpenTDB

Emulates the OpenTDB API (opentdb), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-29.**

## What this Recipe found

The success key is `results` and the rate limit's key
is `result`. One letter, on one of four answers. Code reading
`body.results.length` works against a hit, against no results and against an
invalid parameter, and finds `undefined` on the fourth -- and the fourth is the
one that arrives when the caller is going too fast, so the crash appears under
load and not before it.

**And three of the four are HTTP 200.** No results is 200. An invalid parameter
is 200. `response.ok` is true for every outcome except the rate limit, and the
only thing separating them is a number in the body that has no name anywhere in
the response. An unknown category is reported as an empty one -- `?category=999`
answers `response_code` 1, for no results, not 2, for a bad parameter -- so a
typo in a category id is indistinguishable from a category that happens to be
empty. The strings are HTML-encoded inside JSON: a question arrives as `Kraft
&quot;Cheez Whiz&quot; ... doesn&#039;t`, entities in a format that needs none,
so every string wants an HTML decode after it has already been JSON-decoded.
`incorrect_answers` is three entries on a multiple-choice question and one on a
boolean, with nothing but the sibling `type` to say which. And `encode=base64`
encodes the enum values too, so `type === "multiple"` is false until even the
fields a caller would never think to decode have been decoded.

## Sources

- Documentation: https://opentdb.com/api_config.php
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve opentdb     # run it
cauldron verify opentdb -v # check every claim
```
