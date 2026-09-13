# deepinfra

Emulates the DeepInfra models API for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-13.**

Struck live against `api.deepinfra.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist. The record shape is the live response, which anyone can fetch.

## What this Recipe found

**The price of three is 2.9999999999999996.**

```json
"pricing": {
  "input_tokens": 1.0,
  "output_tokens": 2.9999999999999996,
  "cache_read_tokens": 0.2
}
```

A binary floating-point number where a price belongs. It is three, arrived at by arithmetic rather than written down, and it renders as `$2.9999999999999996` in anything that formats it without rounding.

There is no currency field and no unit field beside it, so a client cannot tell what the `1.0` counts either.

**And a `discount` sits next to them.** `0.61`, a bare multiplier, with nothing saying whether the prices above are before it or after it.

**The catalogue is public and a wrong key is refused.** No credential answers 200 and 191 models; `Bearer notreal` answers `401 {"detail":"User is not authorized to access this resource"}`.

That is worth recording because the neighbouring provider does the opposite: [novita](../novita) serves its whole catalogue to a wrong key without reading it. Two OpenAI-compatible model listings, one rejecting a bad credential and one ignoring it.

**Every model was created at the epoch.** `created: 0` on all of them — the same zero [cerebras](../cerebras) sends, for the same reason: the field exists because OpenAI's does.

**Every model carries an image generator's settings.** `default_width`, `default_height` and `default_iterations` are on a text model's record, all null. One metadata shape for every kind of model, so a chat model's record has three fields about pixels.

**And the description is Markdown.** The text field carries `[MiMo-V2-Flash](https://github.com/XiaomiMiMo/MiMo-V2-Flash)` — a link, in a syntax a JSON client has no reason to expect, which renders as literal brackets anywhere that does not know to parse it.

**The failures are FastAPI's default.** `{"detail": "…"}` for both the rejected credential and the unrouted path — one key, no code, no type, and the 404's whole body is `{"detail":"Not Found"}`.

## Sources

- Live: `api.deepinfra.com`, struck 2026-09-13. Every case here is a live one; the record shape is the response itself.

## Modelling limits

- **One route.** Listing models. Chat completions, embeddings, image generation, the inference endpoints and the account surface each want their own evidence.
- **One model of the 191.** The catalogue is large and the shape is uniform; the fixture carries the record the live response leads with.
- **The description is abridged.** The live one is several sentences; the fixture keeps the Markdown link, which is the part that matters.
- **No `spec:`.** DeepInfra documents this API as prose and an OpenAI-compatible surface. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** DeepInfra is reached through an OpenAI client with the base URL changed, which names no DeepInfra package at all. Checked 2026-09-13.
