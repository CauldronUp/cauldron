# Perplexity

Emulates the Perplexity API (v1), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-08-31.**

## What this Recipe found

**One route is versioned and the other is not.**
`POST /chat/completions` is a 401 and `POST /v1/chat/completions` is a 404 with
zero bytes, while `/v1/models` is versioned and works -- and that model
catalogue is documented as backing a different product entirely, so a client
listing models to pick one for chat is reading the wrong list.

## Sources

- Documentation: https://docs.perplexity.ai/api-reference
- Machine-readable description: https://docs.perplexity.ai/openapi.json, last checked 2026-08-31
  `cauldron drift perplexity` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve perplexity     # run it
cauldron verify perplexity -v # check every claim
```
