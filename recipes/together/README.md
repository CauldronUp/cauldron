# Together AI

Emulates the Together AI API (v1), for local development and tests.

**6 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **OpenAI-compatible model list is a bare array**
where OpenAI sends `{object, data}` -- so a client following Together's own
compatibility instructions reads undefined on its first call.

## Sources

- Documentation: https://docs.together.ai/docs/inference/openai-compatibility
- Machine-readable description: https://docs.together.ai/openapi.yaml, last checked 2026-09-01
  `cauldron drift together` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve together     # run it
cauldron verify together -v # check every claim
```
