# flotiq

Emulates the Flotiq content-object listing for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the Flotiq API documentation, and struck live on 2026-09-14 with no key, with a wrong key in the header, with a wrong key in the query parameter, on a path that does not exist, with a method the path does not take, on the schema endpoint and on the root.

## What this Recipe found

**Every response carries a frame policy no browser has honoured in a decade.** On all seven requests:

```
x-frame-options: ALLOW-FROM https://editor.flotiq.com
```

`ALLOW-FROM` was dropped from the specification. Chrome never implemented it; Firefox removed it. A browser that does not understand the directive treats the whole header as invalid, so the effect is *no* framing protection rather than the narrow one it names. It rides on `application/json` responses, where framing is not a thing that happens, and it names the vendor's own editor host to every unauthenticated stranger who makes a request.

**One answer for every mistake.** No key, a wrong key in the header, a wrong key in the query parameter, a path that does not exist, a verb a path does not take, and the schema endpoint all answer:

```json
{"code":401,"message":"Unauthorized"}
```

with `code` repeating the status that is already on the status line.

**And the machine-readable description is behind the credential.** `/api/v1/open-api-schema` answers that same 401, so the document describing the API cannot be read by anyone deciding whether to use it.

**The identifier contains the type.** An object's `id` is `blogposts-456712` — the content type's name, a hyphen, then a number. The id is not opaque, moving an object between types would change it, and parsing it back out is the only way to learn the type from the id alone. Which is also `internal.contentType`, one level down.

**Two timestamps are the empty string and a version number is negative.** `deletedAt` is `""`. `publishedAt` is `""`. `publicVersion` is `-1`. Beside them, `status` is `"public"`. Four fields describing publication, in three conventions for "no", on an object the fifth calls public.

**The user's fields and the system's share one level.** `id` and `internal` sit beside `title` and `postContent` — whatever the content type defines. A type with a field called `id`, or `internal`, has nowhere to put it.

Also pinned: the envelope carries four counters — `total_count`, `total_pages`, `current_page`, `count` — and no cursor; the key may travel as `X-AUTH-TOKEN` or as `?auth_token=`, and the query form puts it in every access log and `Referer`; `x-generated-by: Flotiq` is on every response including the 401s; and the API host's root is a styled HTML page titled "Flotiq - page not found".

## Sources

- [Flotiq API documentation](https://flotiq.com/docs/API/) — the two credential carriers and the content-object endpoints.
- Live: `api.flotiq.com`, struck 2026-09-14 with no key, a wrong key in each carrier, an unrouted path, a wrong method, `/api/v1/open-api-schema`, and the root.

## Modelling limits

- **The description exists and cannot be read.** Flotiq serves an OpenAPI document at `/api/v1/open-api-schema` and requires a key for it, so there is no `spec` to fingerprint from outside.
- **The success side is document-derived.** Listing content needs a real key; the records here are the documented envelope and `internal` block with values of the documented shapes, and every case reading them is marked documentation-only.
- **One content type is served.** `blogposts` is the documentation's own example. Flotiq's paths are `/api/v1/content/{ctdName}`, so a real account has as many listings as it has types; this Recipe has one, and its user fields (`title`, `postContent`) are that type's.
- **The not-found body is not observed.** Reaching it needs a real key; the 404 served here is the provider's own envelope with the status it would carry, marked documentation-only.
- **The unrouted-path 401 is what live sends.** It is declared here as `unknown_route` so the shape is named rather than inherited, but no credential ever gets far enough for a routing failure to be visible.
- **Nothing is mapped in detection.** Flotiq is reached through `@flotiq/flotiq-api-sdk` or a plain header call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
