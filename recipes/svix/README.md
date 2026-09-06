# Svix

Emulates the Svix API (v1), for local development and tests.

**23 conformance cases, 3 checked against the live API on 2026-09-05.**

The delivery and backoff cases still cite documentation, since Svix's test environment has the same waiting problem the real one does. The credential and routing shapes needed no environment at all, and checking them live found this Recipe had collapsed two sentences into one.

## What this Recipe found

Accepted is not delivered. Creating a message answers 202 with an identifier and no delivery information, because none exists yet -- delivery is attempted for hours afterwards, with backoff between tries, and a message doesn't carry a status at all; the attempts do, one object per try. A failing endpoint gets disabled after repeated failure and stays disabled until someone re-enables it, and that's the state that breaks an integration silently: messages keep getting accepted, attempts simply stop, and nothing tells the sender.

## What checking it live found

No Authorization header at all names the header (`` `Authorization` header required ``) and a present, wrong token names the token instead (`Invalid token. Have you set the correct server URL?`) -- two sentences this Recipe had merged into one generic "Authentication failed". A path nothing declares answers empty, zero bytes and no `Content-Type`, rather than the `{code, detail}` shape an unrouted path borrows by default, and it is resolved before the credential rather than after.

## Sources

- Documentation: https://api.svix.com/docs
- Machine-readable description: https://api.svix.com/api/v1/openapi.json, last checked 2026-08-31
  `cauldron drift svix` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve svix     # run it
cauldron verify svix -v # check every claim
```
