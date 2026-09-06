# Miro

Emulates the Miro API (v2), for local development and tests.

**18 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.miro.com, no account and no key -- a missing credential and a wrong one both answer the identical body. This file declared code "1.0000" and "Invalid or missing access token"; the real answer is `tokenNotProvided` / "No authorization data was found on the request". Fixed below.

## What writing this Recipe changed

It nests coordinates under `position`, so code reading `item.x` gets nothing and
the item silently lands at the origin -- a failure that looks like a layout bug
rather than a parsing one.

## Sources

- Documentation: https://developers.miro.com/reference
- Machine-readable description: https://raw.githubusercontent.com/miroapp/api-clients/main/packages/generator/spec.json, last checked 2026-09-05
  `cauldron drift miro` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve miro     # run it
cauldron verify miro -v # check every claim
```
