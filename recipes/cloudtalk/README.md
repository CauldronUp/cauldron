# CloudTalk

Emulates the CloudTalk API (v1.7), for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **two hosts disagree about their own API** on both
credential granularity and routing order.

## Sources

- Documentation: https://developers.cloudtalk.io/api-reference/overview
- Machine-readable description: https://developers.cloudtalk.io/api-reference/openapi.json, last checked 2026-09-01
  `cauldron drift cloudtalk` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cloudtalk     # run it
cauldron verify cloudtalk -v # check every claim
```
