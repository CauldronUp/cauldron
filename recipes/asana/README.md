# Asana

Emulates the Asana API (1.0), for local development and tests.

**24 conformance cases, 1 checked against the live API.**

One was struck live against app.asana.com on 2026-09-05, with no Authorization header at all and again with a made-up Bearer token: byte-identical 401. This file had claimed the message was "Not Authorized: authentication is required"; Asana sends plainly "Not Authorized", with no such phrase attached.

## What this Recipe found

Asana has no sandbox -- a personal access token reaches a real workspace, so a test that creates tasks creates them where people can see them and get notified. The identifier is gid, not id, and it's a numeric string large enough to lose precision if parsed as a number; code that reads task.id finds nothing. Every object also declares its own resource_type, and everything -- single object or collection -- is wrapped under "data", so reading the response body directly finds nothing either way.

## Sources

- Documentation: https://developers.asana.com/reference/rest-api-reference
- Machine-readable description: https://raw.githubusercontent.com/Asana/openapi/master/defs/asana_oas.yaml, last checked 2026-08-30
  `cauldron drift asana` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve asana     # run it
cauldron verify asana -v # check every claim
```
