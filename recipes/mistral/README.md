# Mistral

Emulates the Mistral API (v1), for local development and tests.

**12 conformance cases, 4 checked against the live API on 2026-08-31.**

## What this Recipe found

Four are xAI's. Mistral and xAI each offer an
OpenAI-compatible surface, and each marks where that stops. Mistral's errors are flat `detail`
rather than OpenAI's nested `error{}`, its ids are `cmpl-` not `chatcmpl-`, and
`AssistantMessage` has no `refusal` field at all. xAI distinguishes a missing
credential (401) from a wrong one (400) where Mistral answers identically to
both -- and on one xAI route body validation runs *before* authentication,
which is the opposite order from its own `/v1/models`.

## Sources

- Documentation: https://docs.mistral.ai/api/
- Machine-readable description: https://docs.mistral.ai/openapi.yaml, last checked 2026-08-31
  `cauldron drift mistral` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mistral     # run it
cauldron verify mistral -v # check every claim
```
