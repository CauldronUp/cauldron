# outreach

Emulates the Outreach prospects API for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-09-11.**

Read from the OpenAPI document Outreach serves without a credential at `api.outreach.io/api/v2/schema/openapi.json`, and struck live against `api.outreach.io` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The API demands JSON:API and does not answer it.** Every request must carry `Content-Type: application/vnd.api+json` or take a 415, and the documented error object is JSON:API's:

```json5
{ errors: [ { id: "unsupportedMediaType", title: "…", details: "…" } ] }
```

The live failure is not that:

```
(no header)      401 application/json {"error":"No Authorization header.","description":"You must provide an Authorization header with a valid access token."}
Bearer notreal   401 application/json {"error":"Invalid JWT token.","description":"The JWT token could not be decoded."}
```

No `errors` array, no `id`, no `title`, and `application/json` rather than the media type the API insists its callers send. A client written against the documented shape reads `body.errors[0].title` and gets a TypeError on the first failure it meets.

**And the documented error object misuses JSON:API twice.** The specification defines `id` as "a unique identifier for this particular occurrence of the problem"; Outreach puts an error *code* there — `unsupportedMediaType`, the same string on every occurrence. The specification calls the human-readable explanation `detail`; Outreach calls it `details`.

**The failure names the credential's format.** `Invalid JWT token.` and `The JWT token could not be decoded.` tell a caller who has presented nothing valid that the token is a JWT — which is to say that it is base64 and readable without the key. Both sentences end in a full stop, so `error` is prose rather than a code, and anything switching on it is switching on English punctuation.

**The machine-readable document declares no parameters and no 401.** `GET /prospects` lists `parameters: []` and responses 200 and 422 — on an endpoint whose prose documentation spends four sections on `filter`, `sort`, `count` and `page`, and which answers 401 to every request arriving without a token. The definition omits every input and the output a caller meets first.

`cauldron drift` names it directly:

```
outreach   not backed: no operation the Recipe routes to answers 401, which invalid_jwt_token declares
           not backed: no operation the Recipe routes to answers 401, which no_authorization_header declares
```

**Two security schemes are defined identically under different names.** `bearerAuth` and `s2sAuthToken` are both `{type: http, scheme: bearer, bearerFormat: JWT}`, byte for byte, and only the first is referenced anywhere.

**The page-size parameter has two names.** Cursor pagination is documented as `page[size]`; offset pagination, three paragraphs later on the same page, as `page[limit]`. Same quantity, two spellings, and the `links` each style returns are spelled with its own.

**Nothing on a prospect is required.** The `prospect` schema has 230 attributes and no `required` array, and `prospectResponse` has none either — so `id` and `type`, which JSON:API requires on every resource object, are optional in Outreach's own description of one.

**And the field metadata is images.** Every attribute description in the OpenAPI begins with raw HTML:

```
<img src=https://developers.outreach.io/badges/filterable.svg>
<img src=https://developers.outreach.io/badges/sortable.svg>
<img src=https://developers.outreach.io/badges/readonly.svg>
```

Whether a field can be filtered, sorted or written is encoded as a hotlinked picture inside a free-text description. A generated client gets an unescaped `<img>` tag in its docstring, a reader offline sees nothing at all, and the only mechanical way to answer "is this filterable" is to match an image filename.

**Relationships come in two shapes on one record.** A to-one relationship carries `data` with a type and an id; a to-many carries `links.related` with a URL. Reading a relationship means testing which kind it is before you can read it.

## Sources

- Live: `api.outreach.io`, struck 2026-09-11.
- [`openapi.json`](https://api.outreach.io/api/v2/schema/openapi.json) — served without a credential, 1.4MB.
- [Making requests](https://developers.outreach.io/api/making-requests/) — content negotiation, filtering, both paginations, and the documented error envelope.

## Modelling limits

- **One route.** Listing prospects. Accounts, opportunities, sequences, mailings, tasks, users, the batch action surface and the custom-object schema endpoints are 146 more paths.
- **The 415 is not modelled.** Sending the wrong `Content-Type` is documented to produce the JSON:API error envelope, and nothing this Recipe does provokes one — a request without a token answers 401 first.
- **The documentation's examples are JSON5.** Unquoted keys, trailing commas and `//...` comments, in every response sample. The fixtures here are transcribed into JSON; the samples as published will not parse.
- **The deprecated offset pagination is not modelled.** `page[offset]` and `page[limit]` still work and are documented as deprecated, with a maximum offset of 10,000. This Recipe serves the cursor form the documentation recommends.
- **Nothing is mapped in detection.** Outreach publishes no first-party client library for this API — the SDK it does ship is for browser and web-widget extensions rather than for the REST surface — so a dependency list does not say this API is in use. Checked 2026-09-11.
