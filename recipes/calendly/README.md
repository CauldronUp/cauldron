# Calendly

Emulates the Calendly API (v2), for local development and tests.

**9 conformance cases, 1 checked against the live API.**

Struck live against api.calendly.com on 2026-09-05, both with no Authorization header at all and with a made-up Bearer token: byte-identical 401, and this file's claim held exactly -- title "Unauthenticated", message "The access token is invalid".

## What this Recipe found

Calendly addresses everything by URI rather than a bare identifier -- the field is literally called "uri" and holds a full URL, including in relationships, so code expecting to build a path from an id has nothing to build from. Paging is nested under "pagination" with a cursor token rather than an id, so a client that pages by remembering the last id it saw gets stuck on the first page forever.

Calendly has no sandbox, so testing the states that actually matter -- a cancelled event, a no-show, a rescheduled invitee -- means creating real scheduled events that email real confirmation and calendar invites to whoever is named as the invitee, which is why those states are rarely exercised at all.

## Sources

- Documentation: https://developer.calendly.com/api-docs
- Machine-readable description: https://developer.calendly.com/openapi/calendly-api.json, last checked 2026-09-05
  `cauldron drift calendly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve calendly     # run it
cauldron verify calendly -v # check every claim
```
