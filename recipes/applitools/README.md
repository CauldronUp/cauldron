# applitools

Emulates the Applitools Eyes batch-results read for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from the Eyes server API reference at [`applitools.com/docs`](https://applitools.com/docs/eyes/reference/server-api/batches/list-batch-results), and struck live on 2026-09-13 with no key, with a wrong key in the header, with a wrong key in the query parameter, on paths that do not exist, and with methods those paths do not take.

## What this Recipe found

**The problem document's `type` is the status code's own specification.** Every 401:

```json
{"type":"https://tools.ietf.org/html/rfc9110#section-15.5.2",
 "title":"Unauthorized","status":401,
 "traceId":"00-13cbcc15090cfe6121fa1fa0bbfad0eb-5ace6f4db87f3461-00"}
```

RFC 9110 §15.5.2 *is* the definition of 401 Unauthorized. RFC 9457 asks a problem document's `type` URI to identify the problem; this one identifies the status — which is already in `status`, and in `title`, and in the status line. Four copies of one fact, and nothing about what went wrong.

**And the URI is dead.** `tools.ietf.org` was retired. The link 301s to `https://datatracker.ietf.org/doc/html/rfc9110`, and the redirect drops the `#section-15.5.2` fragment, so a reader who follows it lands at the top of a three-hundred-page RFC rather than the paragraph about 401.

**The support identifier points at a trace that was not recorded.** `traceId` is a W3C `traceparent`: version, trace id, span id, flags. The flags are `00` on every response struck — the sampled bit clear, meaning the trace was not sampled. The identifier handed to the caller is for a trace nobody kept. It is fresh on every request, so it genuinely identifies the request; it just does not lead anywhere.

**The problem document appears only where it is least needed.**

```
GET  /api/sessions/batches/{id}   (no key)   401  application/problem+json, 165 bytes
GET  /api/cauldron-nope                      404  no body, no Content-Type
PUT  /api/sessions/batches                   405  no body, no Content-Type
GET  /api/sessions                           405  no body, no Content-Type
```

The one failure whose cause the caller already knows is described in a standard format; the two that would tell them *which* path or *which* verb was wrong say nothing at all — not even a content type.

**The refusal is identical whatever the credential is, which makes the auth surface unobservable from outside.** No header, a wrong `X-Eyes-Api-Key`, and a wrong `?apiKey=` all draw the same body with a fresh trace id. Nothing in the response says whether the query parameter the SDKs have always sent is still read.

**A second host answers nothing.** `eyesapi.applitools.com` answers `404` with `Content-Length: 0` to every path tried, with and without a credential.

In the documented response shape: the batch's own identifier is at `statistics.id`, inside the object of counts, beside `statistics.pointerId` — so the record has no id of its own and the thing that identifies it is nested in the aggregate; one of the counts is called **`new`**, a reserved word in several of the languages the Eyes SDKs are written in; `runBy` and `assignedTo` are user email addresses on a test result; `env.app` is documented "Browser/application", under a key named for neither; and `steps` is present only when the request carried `steps=true`, so the shape of the record depends on a query parameter.

## Sources

- [List batch results](https://applitools.com/docs/eyes/reference/server-api/batches/list-batch-results) — the path, the `X-Eyes-Api-Key` header, and every field of the response.
- Live: `eyes.applitools.com` and `eyesapi.applitools.com`, struck 2026-09-13 on eight requests across six paths and four methods.

## Modelling limits

- **No description is published.** Applitools serves no OpenAPI document for the Eyes server API, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Reading a batch needs a real key and a real batch; the record here is the reference's own field list with values of the documented types, and every case reading it is marked documentation-only.
- **`traceId` is fixed here and fresh there.** Live, every response carries a new trace id; this Recipe serves one constant value, and the case that matters asserts the shape — `00-<32 hex>-<16 hex>-00` — rather than the value.
- **The 405 case is not live.** `PUT` and `DELETE` on the batch path are routed by the real API and answer 401, not 405; the genuine empty-bodied 405 was struck on `PUT /api/sessions/batches` and `GET /api/sessions`, which this Recipe does not route. The case here uses a method this Recipe does not route, and is marked documentation-only for that reason.
- **The batch id in the path is not matched.** Any batch id returns the one fixture batch, because the record carries no id field to match against — which is itself the finding.
- **One route of several.** The batch read. Batch properties, page coverage, the patch, and the session-level endpoints are the rest.
- **Nothing is mapped in detection.** Applitools is reached through `@applitools/eyes-*` on npm or the SDKs for other languages, and none of them resolves to this host through a dependency file. Checked 2026-09-13.
