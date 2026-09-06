# SEC EDGAR

Emulates the SEC EDGAR API (edgar), for local development and tests.

**19 conformance cases, 6 checked against the live API on 2026-08-31.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

EDGAR's, where **the same field is a string in one endpoint and an
integer in another.** `/submissions/CIK0000320193.json` answers
`"cik": "0000320193"` and `/api/xbrl/companyconcept/...` answers `"cik": 320193`
-- same company, same field name, same API. A client carrying the value from one
to the other sends `320193` where the path wants `0000320193` and gets a 404,
and the comparison that would have caught it never fires, because
`"0000320193" == 320193` is false in every language here. The padding is
load-bearing: `/submissions/CIK320193.json` is a 404 for a company whose CIK
genuinely is 320193, so the endpoint that hands you the identifier as an integer
has already destroyed the thing that makes it usable.

**And the failures are XML.** `data.sec.gov` is S3 underneath, so a missing
company answers `<Error><Code>NoSuchKey</Code>` -- from a path ending `.json`,
with `Content-Type: application/xml`, in S3's vocabulary rather than EDGAR's.
Nothing in it mentions a company or a filing. A client calling `.json()` on a
404 throws on the first character.

## Sources

- Documentation: https://www.sec.gov/search-filings/edgar-application-programming-interfaces
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve secedgar     # run it
cauldron verify secedgar -v # check every claim
```
