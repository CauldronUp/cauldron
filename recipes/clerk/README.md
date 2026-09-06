# Clerk

Emulates the Clerk API (2024-10-01), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Clerk's timestamps are Unix milliseconds -- code that divides by a thousand when it shouldn't, or fails to when it should, lands in 1970 or the year 55000, and a seconds-only fake makes both mistakes look correct. Failures also arrive as an array with a trace id sitting beside it rather than inside it: the trace id belongs to the response, and the errors belong to the individual fields that caused them, which is how a signup form maps a failure back to the right input.

## Sources

- Documentation: https://clerk.com/docs/reference/backend-api
- Machine-readable description: https://clerk.com/openapi.json, last checked 2026-08-31
  `cauldron drift clerk` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clerk     # run it
cauldron verify clerk -v # check every claim
```
