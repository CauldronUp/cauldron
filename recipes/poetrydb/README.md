# PoetryDB

Emulates the PoetryDB API (poetrydb), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-29.**

## What this Recipe found

`Status` is the number `404` and the string
`"405"`. The same key, in the same API, holding a JSON number on one failure and
a JSON string on the other: `body.status === 404` catches the first and misses
the second, `body.status === "405"` does the reverse, and `body.status >= 400`
is true for one and false for the other, because `"405"` is not a number.

**And neither of them is an HTTP failure.** Both come back 200, with
`application/json`, so `response.ok` is true and every retry, error boundary and
status check sees a success -- the only thing saying otherwise is that field,
whose type changes between the two ways of saying it. `linecount` is the string
`"14"`, so the API has a number that is a string beside a string that is a
number. `poemcount` is a valid input field and an invalid output field: one
refusal's list of allowed words contains it and the other's does not, and both
are 405. A path segment that is not a field is reported as Method Not Allowed on
a GET that is the only method the endpoint has. And `/author` answers
`{"authors": [...]}` while `/author/{name}` answers a bare array, so two paths
one segment apart cannot share a parser.

## Sources

- Documentation: https://poetrydb.org/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve poetrydb     # run it
cauldron verify poetrydb -v # check every claim
```
