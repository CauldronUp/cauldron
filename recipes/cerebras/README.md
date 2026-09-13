# cerebras

Emulates the Cerebras Inference models API for local development and tests.

**6 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from Cerebras's own API reference, and struck live against `api.cerebras.ai` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The two credential failures are the wrong way round.**

```
(no header)      403 {"detail":"Not authenticated"}
Bearer notreal   401 {"message":"Wrong API Key","type":"invalid_request_error",
                      "param":"api_key","code":"wrong_api_key"}
```

401 means "tell me who you are" and 403 means "I know, and you may not". Cerebras answers 403 to the caller who has said nothing, and 401 to the one whose credential it read and rejected.

Its own error table agrees with the RFC and not with the API — it lists 401 as `AuthenticationError` and 403 as `PermissionDeniedError`.

**And they come from two different frameworks.** `{"detail": "…"}` is FastAPI's default shape; `{"message","type","param","code"}` is OpenAI's. One layer answers the unauthenticated caller and another answers the authenticated-but-wrong one, so a client needs two readers — `detail` on one and `message` on the other — to print either.

**`param` names something that is not a parameter.** OpenAI's `param` field identifies the offending *request parameter*. Here it says `api_key`, which travels in a header. The field pointing at what to fix points at a parameter the request does not have.

**An unrouted path answers 404 with nothing in it.** Zero bytes, labelled `text/plain`. `.json()` throws, `.text` is an empty string, and the status is the only thing in the response.

**Every model was created at the epoch.** `created` is documented as "the Unix timestamp (in seconds) of when the model was created", and the reference's own example gives `0` for every model in it. The field exists because OpenAI's does, and it says 1 January 1970 for all of them.

**The listing has no envelope beyond a discriminator.** `{"object": "list", "data": […]}` — no count, no cursor, no `has_more`, and no parameters on the operation to ask for any. The whole catalogue, every time.

## Sources

- Live: `api.cerebras.ai`, struck 2026-09-13.
- [List models](https://inference-docs.cerebras.ai/api-reference/models/list-models) — the record shape and the example.
- [Error Codes](https://inference-docs.cerebras.ai/support/error) — the status table the API disagrees with.

## Modelling limits

- **One route.** Listing models. Chat completions, completions, the single-model lookup and the tool-calling surface each want their own evidence.
- **402, 413, 422 and 429 are not modelled.** The error table names them with SDK exception classes and no bodies, and nothing this Recipe does provokes one.
- **No `spec:`.** Cerebras documents this API as field tables and worked examples, each page also answering Markdown at the same address with `.md` appended. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** The `cerebras_cloud_sdk` packages are real clients, but they are OpenAI-compatible wrappers pointed at whichever base URL a project configures — and most projects reach this API through an OpenAI client with the base URL changed, which names no Cerebras package at all. Checked 2026-09-13.
