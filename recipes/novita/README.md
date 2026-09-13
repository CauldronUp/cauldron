# novita

Emulates the Novita models API for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-13.**

Struck live against `api.novita.ai` on 2026-09-13 with no credential, with a deliberately invalid one, with a malformed one, and on a path that does not exist. The record shape is the live response, which anyone can fetch.

## What this Recipe found

**A wrong credential is not refused. It is ignored.** All three answer 200 and the whole catalogue:

```
(no header)                    200  117 models
Authorization: Bearer notreal  200  117 models
Authorization: notabearer      200  117 models
```

The endpoint never reads the header at all. So a client that has misconfigured its key gets a perfectly successful response here and a failure on its first real call — and the endpoint people reach for to check their setup is the one that cannot tell them anything about it.

**The price list is public.** 117 models with per-token pricing — input and output, cached and uncached — to a caller who has presented nothing.

**Every price is given twice, in two units, and one of them is a string.**

```json
"price_per_m": 1500,
"price_per_m_decimal": "0.15"
```

An integer beside a decimal string that is the same money at a different scale. 1500 is 0.15 at a scale of ten thousand, and nothing on the record says so: no currency, no exponent, no unit field. A client that reads the integer and formats it as cents is out by a factor of a hundred.

**And it is given twice again.** `input_token_price_per_m` at the top level and `pricing.prompt.price_per_m` inside are the same number, as are `output_token_price_per_m` and `pricing.completion.price_per_m`. Four fields, two numbers, and nothing saying which pair a client should believe if they ever disagree.

**`origin_price_per_m` sits beside `price_per_m`** — a list price and a charged price — so the record carries a discount structure with no discount in it and no field naming the difference.

**Three OpenAI-compatibility fields carry nothing.** `permission` is null, and `root` and `parent` are the empty string, where OpenAI sets `root` to the model's own id. The compatibility is in the key names rather than in the values, and a client written against OpenAI's shape reads `""` where it expects an identifier.

**The envelope is not OpenAI's either.** `{"data": […]}` with no `"object": "list"` beside it — on a path spelled `/v3/openai/models`.

**A model is named three times.** `id` and `title` are the same string, `zai-org/glm-5.3-flash`, and `display_name` is a third spelling of it, `GLM 5.3 Flash`. One of the three is what goes back in a request and nothing on the record says which.

**`status` is an unnamed integer.** `1`, with no enum, no sibling string, and nothing in the response to compare it against.

**And an unrouted path is not JSON.** `404 page not found` in plain text — Go's own `http.NotFound` string, the same one [pirsch](../pirsch) answers, from an API whose successes are JSON.

## Sources

- Live: `api.novita.ai`, struck 2026-09-13. Every case here is a live one; the record shape is the response itself.

## Modelling limits

- **One route.** Listing models. Chat completions, images, video, audio, the GPU-instance surface and the OpenAI- and Anthropic-compatible endpoints each want their own evidence.
- **Nothing here refuses a credential.** The route is declared public because that is what it is: three differently-wrong credentials were sent and all three were answered. There is no state in which this endpoint rejects a key, so there is no rejection to model.
- **One model of the 117.** The catalogue is large and the shape is uniform; the fixture carries the record the live response leads with, and a second one would say nothing the first does not.
- **No `spec:`.** Novita documents this API as prose. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** Novita is reached through an OpenAI client with the base URL changed, which names no Novita package at all. Checked 2026-09-13.
