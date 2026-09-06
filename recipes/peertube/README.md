# peertube

Emulates the PeerTube federated-video API (v1), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-06.**

Struck against `peertube2.cpy.re` — a host named in PeerTube's own `servers` list as a live test server — with no credential, because the video surface is public. Written against PeerTube's own OpenAPI 3.0 document: `Chocobozzz/PeerTube`, `support/doc/api/openapi.yaml`, 249 paths, version 8.1.0.

## What this Recipe found

**Three identifiers per video, and the embed URL uses the third.** Struck live:

```
id         951957
uuid       "c40f992b-4ba3-4158-8e38-7c710db01707"
shortUUID  "qdcVhVSdcZMtpqeqnfVsJk"
embedPath  "/videos/embed/qdcVhVSdcZMtpqeqnfVsJk"
```

All three address the same video and `/videos/{id}` accepts any of them. The integer is local to the instance and meaningless anywhere else; the uuid is global; the shortUUID is the uuid re-encoded shorter, and it is the one the embed path uses. A client that stores `id` cannot build an embed URL; one that stores `uuid` has to convert.

**`{id, label}` on four fields, with three different types of id.** Struck live, on one video:

| field | value |
|---|---|
| `privacy` | `{"id": 1, "label": "Public"}` |
| `licence` | `{"id": 6, "label": "Attribution - Non Commercial - No Derivatives"}` |
| `language` | `{"id": "fr", "label": "French"}` |
| `category` | `{"id": null, "label": "Unknown"}` |

Same shape four times; the id inside is a number, a number, a **string** and a **null**. So `video.privacy === "Public"` is false, `video.language.id` needs a string comparison where `video.licence.id` needs a numeric one, and **`video.category` is truthy on a video with no category** — because "no category" is an object whose id is null and whose label is the word "Unknown". A client filtering uncategorised videos with `if (!video.category)` filters out none of them.

`licence` is spelled the British way, so a client mapping `license` reads `undefined` on every video.

**Errors are RFC 7807, and the content type says so.** Struck live: `Content-Type: application/problem+json; charset=utf-8` — not `application/json`. A client that checks the content type before parsing, or a framework configured to decode only `application/json`, skips the body on exactly the responses that explain themselves.

**And the problem documents are not the same shape as each other.** Struck live, one after the other:

```
404  {"type":"about:blank","title":"Not Found","detail":"Video not found",
      "status":404,"docs":"https://docs.joinpeertube.org/…#operation/getVideo"}

401  {"type":"https://docs.joinpeertube.org/…#section/Errors/unauthorized_request",
      "detail":"Token is invalid","status":401,"code":"unauthorized_request"}
```

The 404 has `title` and `docs` and no `code`. The 401 has `code` and no `title` and no `docs`, and its `type` is a real URL where the 404's is `about:blank`. Two RFC 7807 responses from one host sharing two members out of five: `problem.title` is undefined on the 401, `problem.code` is undefined on the 404.

**A 400 carries a hyphenated key.** Struck live: `"invalid-params":{"count":{"type":"field","value":"999","msg":"…"}}`. RFC 7807 allows extension members and does not constrain their names, so this is legal — and `problem.invalid-params` is a *subtraction* in JavaScript. It has to be `problem["invalid-params"]`, and a struct generator does whatever it does with a hyphen. The 400 also carries `instance`, which the other two do not.

**An oversized page is refused rather than trimmed.** `count=999` answers 400. Most providers modelled here clamp silently; this is the louder and better half of that choice.

**Three fields answer one question about adult content.** `nsfw` is a boolean, `nsfwFlags` is an integer bitfield (struck live as `0`), `nsfwSummary` is a nullable string. Two of the three "nothing" values are falsy, so they are indistinguishable from absent.

**Rate-limit headers on an unauthenticated request** — `X-RateLimit-Limit: 500`, `X-RateLimit-Remaining: 499`, on a public listing with no credential.

**A live stream has a duration of zero,** which is not the same as a zero-length video — and `isLive` is the only field that says which. A listing mixes local and federated records, and `isLocal` is the only field that separates them.

**Paging is `count` and `start`,** not `limit` and `offset`.

## Modelling limits

- **Videos, listed and fetched.** 249 paths is a federated video platform: accounts, channels, playlists, comments, live streams, subscriptions, abuse reports, the ActivityPub surface and the upload flow each want their own evidence.
- **The three identifiers are served on the record; the fetch route takes the uuid alone.** Accepting all three would mean routing one path three ways, and what the finding needs is that all three are in the response.
- **`state` is served as null.** It is non-null only while a video is being transcoded, which a fixture cannot represent honestly.
