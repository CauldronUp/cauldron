# SurveyMonkey

Emulates the SurveyMonkey API (v3), for local development and tests.

**15 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **own example gets its own wire wrong**, quoting
a numeric error id as a string, and whose responses carry identifiers and never
a label.

## Sources

- Documentation: https://developer.surveymonkey.com/api/v3/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve surveymonkey     # run it
cauldron verify surveymonkey -v # check every claim
```
