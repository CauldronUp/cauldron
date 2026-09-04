# Svix

Emulates the Svix API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Accepted is not delivered. Creating a message answers 202 with an identifier and no delivery information, because none exists yet -- delivery is attempted for hours afterwards, with backoff between tries, and a message doesn't carry a status at all; the attempts do, one object per try. A failing endpoint gets disabled after repeated failure and stays disabled until someone re-enables it, and that's the state that breaks an integration silently: messages keep getting accepted, attempts simply stop, and nothing tells the sender.

## Sources

- Documentation: https://api.svix.com/docs
- Machine-readable description: https://api.svix.com/api/v1/openapi.json, last checked 2026-08-31
  `cauldron drift svix` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve svix     # run it
cauldron verify svix -v # check every claim
```
