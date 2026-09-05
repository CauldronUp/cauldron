# Freshdesk

Emulates the Freshdesk API (v2), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real help desk. Freshdesk is multi-tenant with no shared host, so the refusal cases were struck live, unauthenticated, against support.freshdesk.com -- Freshdesk's own help desk, running on its own product.

## What this Recipe found

Status and priority are integer enums, not words -- a ticket's status is `2`, not `"open"`, with `4` meaning resolved and `5` meaning closed. Code that compares against a string never matches, and code that treats "not 5" as still open counts every resolved ticket as outstanding too.

The live probe found that not every Freshdesk failure carries its sentence under `description`, the field this file's other errors use: a missing or wrong credential, and a wrong method on a real path, both put it under a plain `message` instead. An unrouted path answers before authentication is ever consulted, while a wrong method on a real path is a genuine 405 with a computed `Allow` header that matched exactly what came back live.

Collections are bare arrays with paging carried in a `Link` header rather than in the body. Freshdesk also has no sandbox on most plans, so a test that creates a ticket against a real help desk can fire a live automation and email an actual customer -- there's no safe environment to point tests at short of a separate account that inevitably drifts.

## Sources

- Documentation: https://developers.freshdesk.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve freshdesk     # run it
cauldron verify freshdesk -v # check every claim
```
