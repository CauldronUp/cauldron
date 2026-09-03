# ClinicalTrials.gov

Emulates the ClinicalTrials.gov API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**Every failure is plain text and one of
them is nothing at all.** A malformed identifier answers 400 with `Parameter
`nctId` has incorrect format`; a parameter that will not convert answers 400
naming the integer width; and a path that does not exist answers 404 with no
body and no `Content-Type` header, so there is not even a wrong document to
parse. The two that do speak wrap their parameter names in Markdown backticks,
in responses declared as plain text, formatted for a renderer that is not there.

**And a date says whether it actually happened.** `startDateStruct` is
`{"date": "2019-01-24", "type": "ACTUAL"}` and `primaryCompletionDateStruct` is
`{"date": "2024-12-31", "type": "ESTIMATED"}` -- a fact and a guess, told apart
by one word a level down, in fields whose names end in a word about the code.
`overallStatus` is `"UNKNOWN"` beside a `lastKnownStatus` of `"RECRUITING"`, so
the study is both. `statusVerifiedDate` is `"2024-01"`, a year and a month with
no day, among dates that have all three. And `phases` is a list of one string,
and the string is `"NA"`.

## Sources

- Documentation: https://clinicaltrials.gov/data-api/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clinicaltrialsgov     # run it
cauldron verify clinicaltrialsgov -v # check every claim
```
