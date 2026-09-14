# voyageai

Emulates the Voyage AI embeddings endpoint for local development and tests.

**7 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from Voyage's own published Python client ([`voyage-ai/voyageai-python`](https://github.com/voyage-ai/voyageai-python)), and struck live against `api.voyageai.com` on 2026-09-13 with no credential, with a malformed body, with an unknown argument, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The body is checked before the credential.** With no `Authorization` header on any of them:

```
{}                                       400  …missing one or more required arguments
not json                                 400  …Error for argument '0': JSON decode error
{"input":["hi"],"model":"…","zzz":1}     400  …Argument 'zzz' is not supported by our API
{"input":["hello"],"model":"voyage-3"}   401  {"detail":"Unauthorized"}
```

An anonymous caller can enumerate the request contract one argument at a time — the validator names the unsupported argument and says the required ones are missing — and only a well-formed request gets as far as being refused.

**And the refusal is one word.** `{"detail":"Unauthorized"}`: no code, no request id in the body, nothing to branch on.

**The validation message says the same thing twice, with a `.;` between.** Sending `{}`:

```
The request body is not valid JSON, or some arguments were not specified properly.
In particular, Your request body was missing one or more required arguments.; Your
request body was missing one or more required arguments.
```

A fixed preamble that offers two possibilities, then "In particular," and then one clause per problem — each a complete sentence starting with a capital, joined by a full stop and a semicolon. A body that is not JSON at all is blamed on `argument '0'`: a positional argument, in an API whose arguments are all named.

**The client reads the request id from the wrong header.** `error.py` does `self.headers.get("request-id", None)`, and the response carries `x-request-id`. So `VoyageError.__str__`, which prints `Request {id}: {message}` when it finds one, never finds one.

**And it parses an error shape the API does not send.** `construct_error_object` returns nothing unless the body has an `error` key holding a dict. Every failure above is `{"detail": "…"}`. The published client reads two things off a failure and gets neither.

**The response wrapper refuses to hold an empty string.** `VoyageResponse.__setitem__` raises `ValueError` on `""` with the message "We interpret empty strings as None in requests" — so any response field that arrives as two quotes throws inside the client, and the reason given is about requests.

**And the client keeps two things and discards the rest.** `EmbeddingsObject.update` appends `d.embedding` for each item and adds `response.usage.total_tokens`; nothing else survives. Its `embeddings` is typed `Union[List[List[float]], List[List[int]]]`, so one field holds floats or integers depending on what was asked for — and `update` accumulates across calls rather than replacing.

**The server names itself `uvicorn`** in a response header on every request.

## Sources

- [`voyage-ai/voyageai-python`](https://github.com/voyage-ai/voyageai-python) — `error.py`, `api_resources/response.py` and `object/embeddings.py`.
- Live: `api.voyageai.com`, struck 2026-09-13 with no credential, a malformed body, an unknown argument, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The embeddings endpoint. Reranking, multimodal and contextualized embeddings are the rest.
- **The credential is checked before the body here, and live it is the other way round.** That ordering is the finding above; a declarative Recipe authenticates first, so the case that provokes the validation failure presents a good credential, and the live transcript for the anonymous forms is recorded rather than emulated.
- **The malformed-JSON failure is recorded, not served.** A conformance case sends JSON by construction and cannot send something that is not.
- **The response carries only what the client reads.** Voyage publishes no description of its responses, so the fixture holds `data[].embedding` and `usage.total_tokens` — the two members `EmbeddingsObject` keeps — and does not invent fields the API may add beside them.
- **Nothing is mapped in detection.** Voyage is reached through `voyageai` on PyPI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
