# Freshdesk

Emulates the Freshdesk API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Status and priority are integer enums, not words -- a ticket's status is `2`, not `"open"`, with `4` meaning resolved and `5` meaning closed. Code that compares against a string never matches, and code that treats "not 5" as still open counts every resolved ticket as outstanding too.

Collections are bare arrays with paging carried in a `Link` header rather than in the body. Freshdesk also has no sandbox on most plans, so a test that creates a ticket against a real help desk can fire a live automation and email an actual customer -- there's no safe environment to point tests at short of a separate account that inevitably drifts.

## Sources

- Documentation: https://developers.freshdesk.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve freshdesk     # run it
cauldron verify freshdesk -v # check every claim
```
