# attentive

Emulates the Attentive V2 API for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-07.**

Written against Attentive's published OpenAPI document for the V2 API and struck live against `api.attentivemobile.com` on 2026-09-07 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**A 401 with nothing in it at all.** Struck live, both with no `Authorization` header and with `Bearer notreal`:

```
HTTP/1.1 401 Unauthorized
x-content-type-options: nosniff

<no body>
```

No body. No `Content-Type`. And — the part that matters — **no `WWW-Authenticate`**.

RFC 9110 section 15.5.2 says a 401 response **MUST** include one. It is the only mandatory field on this status, and it is the only field this response does not have. So there is nothing anywhere in the response that distinguishes a request carrying no credential from one carrying a wrong credential, and nothing that names the scheme the API wants.

The one header it does send is `nosniff`: an instruction to a browser not to guess the media type of a body that is not there.

[buffer](../buffer) sends an empty 401 too, and puts the whole failure into `WWW-Authenticate` — malformed, single-quoted, and *there*. Attentive sends the empty half and not the other one. Between them they demonstrate that the header and the body are two independent decisions, and that a provider can get either wrong on its own.

**The 404 is the static-file handler.**

```json
{"timestamp":"2026-09-07T22:39:48.884+00:00","path":"/v2/cauldron-nope",
 "status":404,"error":"Not Found","requestId":"c645e37d-6576697",
 "message":"No static resource v2/cauldron-nope."}
```

Those are Spring Boot's default error attributes, and "No static resource" is the resource handler at the end of the chain reporting a missing **file** — because nothing before it claimed the path. A REST API's answer to an unknown endpoint is a filesystem lookup that failed.

**And the same body spells the path two ways.** `path` keeps the leading slash; the sentence in `message` has dropped it. One object, one path, two renderings.

**`requestId` is not a UUID.** `c645e37d-6576697` is eight hex characters, a hyphen, then seven — the first two groups of a UUID and nothing more. Anything validating it against a UUID pattern rejects the one identifier support would ask for.

**The document declares five failures and describes none of them.** `GET /v2/segments` lists 400, 401, 403, 429 and 500, and each is a `description` string with no `content` and no schema. A generator reading it produces a client that knows the statuses and nothing at all about the bodies — which, for the 401 above, happens to be exactly right, by accident.

**A segment has no identifier of its own.** `SegmentResponse` is `externalId`, `name`, `description`, `created` and `updated`. There is no `id`. The only identifier is the one the caller supplied, `/v2/segments/external/{externalId}` is the only lookup path, and **a segment created without an external id cannot be addressed at all**.

The schema also has no `required` array, so not even that field is guaranteed to be present.

**Two fields for one fact.** The listing carries `cursor`, documented as "null if no more results", and `hasMore` beside it. Either alone would answer the question. [affirm](../affirm) has the same pair with the cursor missing — a boolean saying more exists and nothing to ask for it with — so between the two providers the flag is redundant once and load-bearing once.

**Authentication is declared twice, in two vocabularies, and they disagree.** The document has a real `security` block, top-level and again per operation, *and* a vendor extension `x-requires-auth` on five of its seven paths. `/v2/me` and `/v2/user/attributes` carry the standard declaration and not the custom one. A reader following the extension concludes those two are public; a reader following `security` concludes they are not.

**A segment never edited has two identical timestamps.** `created` and `updated` are the same string until something changes, so `updated` cannot be used to find changed records without comparing it to `created` first.

## Modelling limits

- **One route.** Listing segments. Segment membership, bulk jobs, user attributes and the `/v2/me` identity endpoint each want their own evidence.
- **The v1 document is not modelled.** Attentive publishes two OpenAPI documents; this Recipe describes the seven-path V2 one, and the 284KB v1 surface is a separate piece of work.
- **Nothing is mapped in detection.** Attentive's integrations are platform apps and a tag rather than a client library, and nothing on npm, Packagist or the Go module proxy calls `api.attentivemobile.com` under an obvious name, checked 2026-09-07.
