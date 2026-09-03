# UptimeRobot

Emulates the UptimeRobot API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**A wrong method answers 200.**
`GET is not supported. Use POST instead.` arrives as `text/html` with an HTTP
200, so a client checking `response.ok` proceeds and then throws parsing the
sentence. `PUT` on the same path gets Express's 404 instead.

## Sources

- Documentation: https://uptimerobot.com/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve uptimerobot     # run it
cauldron verify uptimerobot -v # check every claim
```
