# World Bank

Emulates the World Bank API (worldbank), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Bank's, where the success is a two-element array and the
failure is a one-element array. A country lookup answers
`[{"page": 1, "pages": 1, "per_page": "50", "total": 1}, [ {...} ]]` and a bad
code answers `[{"message": [ {...} ]}]`, both with HTTP 200. `body[1]` is the
collection on success and `undefined` on failure, and nothing about the status,
the content type or the outer shape says which arrived except how many elements
the array has -- so `body[1].map(...)` throws a TypeError that names nothing
useful.

**And `per_page` is a string where its three neighbours are numbers.** The paging
object is four fields describing one page, three of them integers and one
quoted. `iso2Code` and `iso2code` are both on the record, capital-C at the top
and lowercase one level down. Three of those codes are not ISO codes at all --
North America is `"XU"`, High income is `"XD"`, Not classified is `"XX"`, all
from the range the standard reserves for private use, in a field named after it.
A missing value is an object of three empty strings rather than `null`, so
`if (country.adminregion)` is true and reading `.value` off it gives `""`. And
the coordinates are quoted on a record whose page numbers are not.

## Sources

- Documentation: https://datahelpdesk.worldbank.org/knowledgebase/articles/889392
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve worldbank     # run it
cauldron verify worldbank -v # check every claim
```
