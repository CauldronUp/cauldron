# stability

Emulates the Stability AI generation-result endpoint for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Stability serves without a credential at [`api.stability.ai/v2alpha/openapi`](https://api.stability.ai/v2alpha/openapi), and struck live on 2026-09-13 with no header, with a wrong key, on a v1 path and on a v2beta one.

## What this Recipe found

**The refusal quotes both ends of the key you sent.** Struck live against v1 with `Authorization: Bearer sk-notarealtokenzzzzzz`:

```json
{"id":"35e2a3c9…","message":"Incorrect API key provided: sk-nota***********zzzz. You can find your API key at https://stability.ai.","name":"unauthorized"}
```

Seven characters of the front, four of the back, and one asterisk for every character between — so the response body carries both ends of the secret and its exact length, into every log that records a response body.

**Two error envelopes on one host.** The same 401, from v1 and from v2beta:

```
/v1/engines/list          {"id", "message", "name"}
/v2beta/stable-image/…    {"id", "errors": [...], "name"}
```

One sends prose as a string under `message`; the other sends it as an array of strings under `errors` and has no `message` at all. A client that handles a Stability failure handles half of them. And the v2beta sentence says the word twice: `"authorization: authorization: invalid or missing header value"`.

**A path that does not exist is answered by the credential check.** `GET /v1/cauldron-nope` with no header answers 401 "missing authorization header", so a typo in the path is reported as a problem with a header.

**And no operation declares a 401 at all.** The document lists 200, 202, 400, 404 and 500 on the result endpoint and never the status every unauthenticated request receives — which is why `cauldron drift` reports all three of this Recipe's refusals as unbacked. The v1 path is not in the document either.

**The key is a parameter, not a security scheme.** Every operation declares `{"name": "authorization", "in": "header", "required": true, "type": "string"}` as an ordinary parameter, with the same paragraph about key reuse repeated on each one. There is no `securitySchemes` entry, so a generated client takes the raw header value as a positional argument and nothing marks it as a secret.

**The same two facts are headers or body fields, and one changes name in between.** From the document's own note:

> This header is absent on JSON encoded responses because it is present in the body as `finish_reason`.

So `finish-reason` and `seed` arrive as response headers when the caller accepts bytes, and `finish_reason` and `seed` arrive in the body when the caller accepts JSON — one of them renamed across the hyphen, and `seed` declared `type: string` in the header and `type: number` in the body. An `Accept` header decides which.

**`CONTENT_FILTERED` is a success.** The enum is `["SUCCESS", "CONTENT_FILTERED"]`, and the second is documented as "successful generation, however the output violated our content moderation policy and has been blurred as a result". A finished generation, a blurred image, and a field that says it worked.

**The in-progress status is an enum of one.** The 202 body is `{id, status}` with `status` enumerated `["in-progress"]` — a field that can only ever hold one value, on a response whose status code already said it.

**And `id` means three different things on one endpoint.** On the 202 it is the generation; on the 404 and every other failure it is "a unique identifier associated with this error"; on the 200 it is absent entirely.

**The document is titled "StabilityAI REST API v2beta" and is served at `/v2alpha/openapi`**, with both `/v2alpha/` and `/v2beta/` paths inside it. And the engine types the vendor's own documentation bundle still ships are `gooseai.EngineInfo`, with a `gooseai.EngineTokenizer` of `GPT2` or `PILE` — the message definitions of a text-generation company, compiled into the reference page of an image API.

## Sources

- [`api.stability.ai/v2alpha/openapi`](https://api.stability.ai/v2alpha/openapi) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- `platform.stability.ai/assets/index-*.js` — the documentation application, which carries the `gooseai` message definitions.
- Live: `api.stability.ai`, struck 2026-09-13 with no header and with a wrong key, against a v1 path and a v2beta one.

## Modelling limits

- **Two routes.** The JSON form of the result endpoint, and one v1 path kept for its failures alone.
- **The v1 route holds no key.** No source describes the JSON the engine listing answers with, so this Recipe declines to invent one: the route accepts no credential, every request to it is refused, and the two refusals are what it is there for.
- **The byte form is not served.** Without `Accept: application/json`, the live endpoint answers the image itself with `finish-reason` and `seed` as headers. This Recipe serves the JSON form and records the other; the name change across the two is the finding, not something the sandbox reproduces.
- **The 202 is not modelled.** A generation that is still running answers `{id, status: "in-progress"}`; a deterministic sandbox has no in-progress state to be in.
- **`id` is a constant in the failures.** Live it is a fresh 32-hex value per response; the live case asserts its shape with a regex.
- **Nothing is mapped in detection.** Stability is reached through `stability-sdk` on PyPI or a plain HTTP call carrying an `Authorization` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
