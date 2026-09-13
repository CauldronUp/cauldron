# jina

Emulates the Jina AI models API for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-13.**

Struck live against `api.jina.ai` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist. The record shape is the live response, which anyone can fetch.

## What this Recipe found

**One model is named three ways and the vendor is spelled differently in each.**

```
"id":              "jina-ai/jina-ocr-v1"
"hugging_face_id": "jinaai/jina-ocr-v1"
"name":            "Jina AI: Jina OCR v1"
```

Hyphenated, unhyphenated, and spaced with a colon. Three identifiers for one model, and no two agree on how to write the company's own name — so a client matching a model across this API and Hugging Face has to know the mapping rather than derive it.

**The models list is public and ignores a wrong key; the endpoint next to it does not.** `/v1/models` answers 200 and the whole catalogue with no header and with `Bearer notarealtoken` alike, while `/v1/embeddings` answers 401 to the same caller. One host, an endpoint that never reads the credential beside one that does — and the first is the one people reach for to check their setup.

**Prices are strings, including the zeroes.** `"prompt": "0.0000005"`, `"image": "0"`, `"request": "0"`. Sending decimals as strings avoids binary floating point, which is right — and it means a client comparing a price to zero is comparing to the string `"0"`. There is no currency field and no unit field anywhere to say what the number counts.

**`quantization` is the empty string.** Not null, not absent: two quotes, where a value belongs.

**A model record carries geography.** `datacenters` is a list of `{country_code}` objects, so "where does this run" sits on the same object as the price and the context length.

**The credential failure is three sentences and a dashboard link.** "Authentication required. Provide your API key via the Authorization header: `'Authorization: Bearer <api-key>'`. Get your API key at https://jina.ai/api-dashboard/key-manager." — an instruction, an example and a URL, in the field a client would print.

**And an unrouted path is `{"detail":"Invalid endpoint"}`** with a request id beside it — the same key as the 401, two words, and no mention of what was asked for.

## Sources

- Live: `api.jina.ai`, struck 2026-09-13. Every case here is a live one; the record shape is the response itself.

## Modelling limits

- **Two routes.** The public model listing, and the embeddings endpoint next to it that exists here only to show it refuses the same caller.
- **One model of the 30.** The catalogue is uniform in shape; the fixture carries the record the live response leads with, and its description is abridged.
- **`/v1/embeddings` is declared public and always fails.** That is the anonymous case, which is the one this Recipe is about; what it does with a real key is a different endpoint's evidence.
- **No `spec:`.** Jina documents this API as prose and an OpenAI-compatible surface. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** Jina is reached through an OpenAI client with the base URL changed, or through `curl` against a documented endpoint — neither names a Jina package. Checked 2026-09-13.
