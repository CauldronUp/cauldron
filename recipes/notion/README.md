# Notion

Emulates the Notion API (2022-06-28), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.notion.com, no account and no key. This file's declared "API token is invalid." matched a well-formed wrong token exactly -- but no Authorization header at all gets a different sentence, "Authorization header must use the format \"Bearer <token>\".", which this file had never distinguished. Split below.

## What writing this Recipe changed

It refuses a request that carries no version header. Forgetting that header is
the classic Notion integration bug, and a fake that waved it through would let
code ship that fails on its first real call.

## Sources

- Documentation: https://developers.notion.com/reference/intro
- Machine-readable description: https://developers.notion.com/openapi.json, last checked 2026-08-31
  `cauldron drift notion` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve notion     # run it
cauldron verify notion -v # check every claim
```
