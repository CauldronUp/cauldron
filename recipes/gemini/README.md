# Gemini

Emulates the Gemini API (v1beta), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A blocked prompt is a 200 with the candidates taken away. The response still carries a model version and usage metadata, but `candidates` -- the only part most integrations read -- is simply absent, with the explanation sitting in a sibling `promptFeedback.blockReason` field that almost nobody opens. Code that does `candidates[0].content.parts[0].text` throws on the very first index. It's worth contrasting with OpenAI's refusal shape, also in this collection: OpenAI hands back a message with `content: null` and the reason in `message.refusal`, so the read succeeds and yields nothing; Gemini removes the element entirely, so the read throws. Same event, two different failure textures, and neither shows up in the HTTP status.

An empty `finishReason` means the model hasn't stopped generating -- the same schema serves both the streaming and non-streaming calls, so an absent value means "still working" on one endpoint and "not applicable" on the other. And caching a prompt doesn't shrink the token count a cost model would read: `promptTokenCount` stays the full effective size even when `cached_content` is set, with `cachedContentTokenCount` the only field that actually says any of it was cheap.

## Sources

- Documentation: https://ai.google.dev/api/rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gemini     # run it
cauldron verify gemini -v # check every claim
```
