# liveblocks

Emulates the Liveblocks room listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the REST API reference at [`liveblocks.io/docs`](https://liveblocks.io/docs/api-reference/rest-api-endpoints), and struck live on 2026-09-14 with no key, with a wrong secret key, with a public key where the secret one belongs, on a path that does not exist, with a method the path does not take, on the previous version and on the root.

## What this Recipe found

**Inside `/v2`, a failure carries four fields and two of them are advice.**

```json
{"error":"MISSING_SECRET_KEY",
 "message":"Missing secret key in authentication header",
 "suggestion":"Please use a secret key, the secret key is available in the dashboard: https://liveblocks.io/dashboard/apikeys",
 "docs":"https://liveblocks.io/docs/api-reference/rest-api-endpoints"}
```

A code, a sentence, what to do about it, and where to read more. The unrouted path gets the same treatment: *"No such endpoint exists. Please look at the documentation for available endpoints. Maybe there was a typo?"* This is the most helpful error envelope in this collection.

**Outside it, the same key carries an English phrase instead of a code.**

```
PUT /v2/rooms   405  {"error":"Method Not Allowed"}
GET /v1/rooms   404  {"error":"Not Found"}
GET /           404  {"error":"Not Found"}
```

One field. No message, no suggestion, no docs — and `error`, which held `MISSING_SECRET_KEY` a moment ago, now holds a status code's reason phrase. Code switching on `body.error` has to know which half of the API it is talking to before it can tell a constant from a sentence.

**Three ways to get the credential wrong, and the statuses cross over.**

| request | status | code |
|---|---|---|
| no `Authorization` header | **401** | `MISSING_SECRET_KEY` |
| `Bearer sk_dev_<wrong>` | **403** | `INVALID_SECRET_KEY` |
| `Bearer pk_dev_<public>` | **401** | `WRONG_KEY_USED` |

A key that is *wrong* is a 403. A key of the wrong *kind* is a 401 — the same status as no key at all. And the third case exists because the public key, the one meant to ship in browser JavaScript, is a plausible thing to reach for, so the API has a name for that mistake and a sentence explaining it: "Public key instead of the secret key".

**Permissions are glob-shaped strings in three fields and two shapes.** `defaultAccesses` is an array — `["*:write"]` — while `groupsAccesses` and `usersAccesses` are maps of id to the same array. The thing that decides who may write is sometimes a list and sometimes a map of lists, and its values are a scope language of their own inside a JSON string.

**The cursor and the parameter that takes it have different names.** The response carries `nextCursor`; the request takes `startingAfter`. A caller has to know the mapping, and nothing in either name says it.

Also pinned: `metadata` values are documented as "value or array", so one key's type depends on what was written to it; `type: "room"` rides on every record returned from an endpoint called rooms; and `limit` defaults to 20 with a documented range of 1 to 100.

## Sources

- [Liveblocks REST API endpoints](https://liveblocks.io/docs/api-reference/rest-api-endpoints) — the room shape, the envelope, and the `limit` / `startingAfter` / `organizationId` / `query` / `userId` / `groupIds` parameters.
- Live: `api.liveblocks.io`, struck 2026-09-14 with no key, a wrong secret key, a public key, an unrouted path, a wrong method, `/v1/rooms`, and the root.

## Modelling limits

- **No description is published.** Liveblocks serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Listing rooms needs a real secret key; the records here are the reference's own room shape with values of the documented types, and every case reading them is marked documentation-only.
- **`query`, `userId` and `groupIds` are not modelled.** `query` is a filter language of its own over room ids and metadata; serving a subset of it would be inventing which subset.
- **One route of many.** The room listing. Room detail, storage, the Yjs document, threads, comments, notifications and the webhook surface are the rest.
- **Nothing is mapped in detection.** Liveblocks is reached through `@liveblocks/node` or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
