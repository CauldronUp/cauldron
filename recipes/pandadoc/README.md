# PandaDoc

Emulates the PandaDoc API (v1), for local development and tests.

**16 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Its create **tells you in words to poll** -- status
`document.uploaded` and an `info_message` beside it, on a document that is not
usable yet.

## Sources

- Documentation: https://developers.pandadoc.com/reference/about
- Machine-readable description: https://raw.githubusercontent.com/PandaDoc/pandadoc-openapi-specification/main/openapi.yaml, last checked 2026-09-01
  `cauldron drift pandadoc` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pandadoc     # run it
cauldron verify pandadoc -v # check every claim
```
