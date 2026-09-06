# OpenAI

Emulates the OpenAI API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A chat completion's refusal does not show up where anyone looks for it. `content` and `refusal` are two separate nullable string fields on the same message, and when a model declines, `content` is null while the actual explanation sits in `refusal`, a field almost nobody reads. It is a 200, with `finish_reason: stop`, and nothing else anywhere signals that a refusal happened. Every integration written the obvious way logs "the model returned nothing," retries, gets the same nothing, and reports an outage, while the model answered both times in the field beside the one being read.

Token accounting has a similar trap: `completion_tokens` includes `rejected_prediction_tokens`, tokens generated but never actually delivered, which OpenAI's own docs say are "still counted in the total completion tokens for purposes of billing." A cost model built by tokenizing the visible response text undercounts every reasoning request. `max_completion_tokens` covers both visible output and invisible reasoning tokens together, so setting it too low against a reasoning model can consume the whole budget before a single visible word is emitted -- `finish_reason: length`, empty content, and a real bill.

Nothing is persisted unless the caller opts in: `store` defaults to false, so a completion just created answers 404 on a direct fetch and is absent from the listing too, for the same reason. Streaming is not modelled at all here, since the response shape genuinely differs, delta objects rather than a full message, and answering a stream request with a whole body would teach the wrong parsing entirely.

## Sources

- Documentation: https://platform.openai.com/docs/api-reference/chat
- Machine-readable description: https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml, last checked 2026-08-30
  `cauldron drift openai` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openai     # run it
cauldron verify openai -v # check every claim
```
