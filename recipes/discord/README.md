# Discord

Emulates the Discord API (v10), for local development and tests.

**13 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource shapes cite documentation rather than an observation on a real guild; the refusal cases were struck live, unauthenticated, against discord.com/api/v10.

## What writing this Recipe changed

Its snowflakes are numeric strings long enough that minting small integers
would have let a rounding bug through unnoticed.

The live probe found Discord's router resolves a request's path and method before it ever looks at authentication: an unmatched path and a wrong method on a declared route both answer before the credential is consulted, and only a route that genuinely matches falls through to the ordinary 401. The declared authentication error was already exactly right.

## Sources

- Documentation: https://discord.com/developers/docs/reference
- Machine-readable description: https://raw.githubusercontent.com/discord/discord-api-spec/main/specs/openapi.json, last checked 2026-09-05
  `cauldron drift discord` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve discord     # run it
cauldron verify discord -v # check every claim
```
