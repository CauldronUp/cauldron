# Steam

Emulates the Steam API (v2), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-08-31.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**A missing app is a 200 that says success is false.**
`appdetails` for a nonexistent id answers thirty bytes -- `{"99999999":
{"success": false}}` -- with `data` absent rather than null, keyed by the id you
asked for rendered as a string. Its 64-bit `steamid`s are quoted, which is a
relief rather than a bug: a bare number would exceed 2^53 eightfold and lose
precision in every JavaScript client. But `categories[].id` is the number `1`
while `genres[].id` is the string `"1"` in the same body, and
`weighted_vote_score` is a bare `0.5` on unvoted reviews and a quoted
sixteen-digit string on voted ones.

## Sources

- Documentation: https://partner.steamgames.com/doc/webapi_overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve steam     # run it
cauldron verify steam -v # check every claim
```
