# Exa

Emulates the Exa API (unversioned), for local development and tests.

**12 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

**The one search API that labels a result cached**, per
document, where its two neighbours label nothing.

## Sources

- Documentation: https://docs.exa.ai/reference/search
- Machine-readable description: https://api.exa.ai/openapi.json, last checked 2026-09-13
  `cauldron drift exa` compares it against what this Recipe claims. It moved on
  2026-09-13: `/search`'s 200 became a `oneOf` of two objects differing by one
  required field with no discriminator, grew a `text/event-stream` alternative
  whose schema is a `oneOf` of six chunk kinds, and marked `resolvedSearchType`
  and `context` deprecated in place. The Recipe's header records what changed;
  nothing it serves moved with it.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve exa     # run it
cauldron verify exa -v # check every claim
```
