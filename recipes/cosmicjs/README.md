# cosmicjs

Emulates the Cosmic v3 Objects API for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Cosmic serves without a credential at [`api.cosmicjs.com/v3/openapi.json`](https://api.cosmicjs.com/v3/openapi.json), and struck live on 2026-09-13 on the collection root and on a bucket that does not exist.

## What this Recipe found

**"Route not found" arrives with a 200.**

```
GET /v3/buckets   200 {"message":"Route not found",
                       "documentation":"https://api.cosmicjs.com/v3/openapi.json",
                       "hints":["/v3/agents/* and /v3/users/* are served by https://dapi.cosmicjs.com",
                                "POST /v3/buckets/{slug}/media is served by https://workers.cosmicjs.com"]}
```

The body is the document's own `Error` schema — which that schema describes, in its own words, as "the error shape returned by every **non-2xx** response". Here it is on a 2xx.

A client branching on `response.ok` treats this as a success and then reads `body.objects`, which is not there.

**It is a helpful failure, which is the other half of the problem.** The body carries a link to the API's own description and two hints naming *other hosts* — `dapi.cosmicjs.com` for agents and users, `workers.cosmicjs.com` for media uploads. One API is three hosts, and the only place that is written down is a body a client will never look at, because the status said everything worked.

**A bucket that does not exist is a 404 with the slug in it.** `{"status":404,"message":"bucket with slug: 'cauldron-nope' not found"}` — the caller's own input, quoted back inside the sentence, with no credential required to ask.

**The credential is optional depending on a setting you cannot see.** From the security scheme's own description:

> The Bucket read key, passed as a query parameter. … Buckets with no read key configured accept reads without one.

So a declared `security` may or may not apply, decided by a dashboard setting on the resource rather than by anything in the request — and the key travels in the query string, so where it does apply it is in every URL.

**A timestamp has three types.** `publish_at` is `["string", "number", "null"]`, described as a UNIX millisecond timestamp, and sits beside `published_at`, which is a date-time string. Two publication timestamps, two encodings, and one of them can arrive as either.

**`thumbnail` is one thing on write and another on read.** Its description: "Media `name` of the Object thumbnail. Returned as a URL on read." The field a client writes and the field it reads back are different kinds of value under one name.

**`content` is deprecated in v3 and still on the record**, marked `deprecated: true`, with its description naming the replacement.

**And nothing on an Object is required.** Twenty properties, no `required` array — not `id`, not `slug`, not `type`. The envelope around them requires exactly one thing, `objects`, so a conforming response is a list of objects that need not be objects.

## Sources

- Live: `api.cosmicjs.com`, struck 2026-09-13.
- [`openapi.json`](https://api.cosmicjs.com/v3/openapi.json) — served without a credential, and named by the 200 above.

## Modelling limits

- **Two routes.** Listing Objects, and the collection root that answers 200 and an error. Object types, revisions, media, batch writes and the two other hosts the hints name each want their own evidence.
- **`/v3/buckets` is not in the description.** `cauldron drift` says so — the path that answers the "Route not found" 200 is not one of the document's twelve. That is the finding rather than a gap to close.
- **The bucket-not-found sentence is recorded, not served.** Live it quotes the slug that was asked for; here an unrouted path answers the generic `Route not found`, because a Recipe declares one message per failure.
- **The 401 is documentation-only.** The document declares it, and whether a caller ever meets it depends on whether the bucket has a read key configured — which no probe from outside can determine.
- **Nothing is mapped in detection.** The `@cosmicjs/sdk` package is a real client, but a project holding one is pointed at whichever bucket slug and read key its configuration names, and neither is in the dependency. Checked 2026-09-13.
