# Nobel Prize

Emulates the Nobel Prize API (v2.1), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Prize's, where **a laureate who does not exist is an array of
length one.** `/2.1/laureate/745` answers `[{"id": "745", ...}]` and
`/2.1/laureate/999999` answers `[{"meta": {...}}]` -- both arrays, both one
element, and the second holds the licence notice and nothing else. A client that
asks "did I get anything?" by reading `.length` is told yes, and finds out only
by reaching for a field. There is no 404 on that path at all.

**And the language maps have keys in them that are not languages.** `city` is
`{en, no, se}`; `cityNow` is `{en, no, se, sameAs, latitude, longitude}`;
`nameNow` is `{en}` where its sibling `name` has three; and the laureate's own
`knownName` has two, so `knownName.no` is undefined where `city.no` is not. Four
objects in one record and no two agree on what a language map is -- and a client
doing `Object.keys(cityNow)` to find translations gets an array of linked-data
URLs and a pair of coordinates among them. The coordinates are strings keeping a
trailing zero each, `"40.825930"`, a precision a JSON number cannot carry.
**The record also links to a different version of the API that served it**: from
`/2.1/`, its own `links[0].href` is `https://api.nobelprize.org/2/laureate/745`.
A share of a prize is a fraction in a string, `"1/3"`, while the money beside it
is an integer with no currency named anywhere.

## Sources

- Documentation: https://app.swaggerhub.com/apis/NobelMediaAB/NobelMasterData/2.1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nobelprize     # run it
cauldron verify nobelprize -v # check every claim
```
