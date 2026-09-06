# Rootly

Emulates the Rootly API (v1), for local development and tests.

**24 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose severity is a second JSON:API document inside a
field** -- a complete `data`/`id`/`type`/`attributes` object as an attribute's
value, rather than a reference with the body in `included`.

## Sources

- Documentation: https://docs.rootly.com/api-reference/overview
- Machine-readable description: https://rootly.com/swagger/v1/swagger.json, last checked 2026-09-05
  `cauldron drift rootly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rootly     # run it
cauldron verify rootly -v # check every claim
```
