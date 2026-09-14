# payload

Emulates the Payload CMS collection listing for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-14.**

Struck live on 2026-09-14 against `payloadcms.com` — Payload's own website, whose REST API answers reads without a credential — on a collection that does not exist, with a method the path does not take, with a delete, on a document that does not exist, on a malformed identifier, and on a page past the end.

## What this Recipe found

**Nine paging fields, and three of them say the same thing.** The envelope is:

```
docs, totalDocs, limit, totalPages, page, pagingCounter,
hasPrevPage, hasNextPage, prevPage, nextPage
```

`hasNextPage: true`, `nextPage: 2` and `page < totalPages` are one fact written three ways. `pagingCounter` is the one-based index of the first record on the page — `(page - 1) * limit + 1`, computable from two fields already present, and sent anyway.

**And a page past the end is reported as if it existed.** Struck live with `?limit=1&page=9999` against a collection of 116:

```json
{"docs":[],"totalDocs":116,"limit":1,"totalPages":116,"page":9999,
 "pagingCounter":9999,"hasPrevPage":true,"hasNextPage":false,
 "prevPage":9998,"nextPage":null}
```

`page` is 9999. `pagingCounter` is 9999. `prevPage` is 9998 — a previous page that does not exist either, offered to a caller who has already walked off the end.

**Two error envelopes, on the same status, from the same API.**

```
GET /api/cauldron-nope          404  {"message":"Route not found \"/api/cauldron-nope\""}
PUT /api/posts                  404  {"message":"Route not found \"/api/posts\""}
GET /api/posts/000000000000…    404  {"errors":[{"message":"Not Found"}]}
```

A bare `message` on one, an `errors` array on the other. And a verb the path does not take is reported as the *path* not being found, with the path quoted inside the sentence — so the JSON carries escaped quotes around a value that could have been a field of its own.

**A malformed identifier is a 404, not a 400.** `GET /api/posts/not-an-objectid` answers the same "Not Found" as a well-formed id that is not there, so a client cannot tell a bug in its own code from a document that was deleted.

**An unauthenticated DELETE reaches body validation.** `DELETE /api/posts` with no credential answers `400 {"errors":[{"message":"Missing 'where' query of documents to delete."}]}` — a sentence explaining how to phrase a bulk delete, to a caller who has not identified themselves.

**A field's value is the name of another field.** `featuredMedia` is `"videoUrl"`, and the record carries both `image: null` and a populated `videoUrl`. The discriminator is a field name, in a string, resolved by the reader.

**The content is an editor's internal state.** `excerpt` and `content` are Lexical documents — `{"root":{"children":[{"detail":0,"format":0,"mode":"normal",…}]}}` — so the API's representation of prose is the shape one particular editor keeps in memory, `detail` and `format` integers included.

**And a content field contains a component tag.** `addToDocs` is the string `<Resource id="6aa2f41763d0ba1d0ceeb5ee" />` — JSX, as data, waiting for something downstream to parse and render it.

Also pinned: `id` is the *last* key of the record and `createdAt` and `updatedAt` are the first two, so the identifier arrives after every piece of content; `_status` carries an underscore among eighteen fields that do not; `guestSocials` is `{}` rather than null; and `/api/users` on this deployment answers 200 to anyone, names and photographs included — which is a configuration choice on one site, not a property of Payload.

## Sources

- Live: `payloadcms.com/api`, struck 2026-09-14 on `/posts`, `/users`, a collection that does not exist, a wrong method, a delete, a document that does not exist, a malformed id, and `?page=9999`.
- [Payload REST API documentation](https://payloadcms.com/docs/rest-api/overview).

## Modelling limits

- **No description is published.** Payload generates an OpenAPI document only through a plugin, and this deployment does not serve one, so there is no `spec` to fingerprint.
- **The fixture content is synthetic; the shapes are not.** Every field, type and quirk here was observed on a live response — the `featuredMedia` discriminator, the Lexical trees, the JSX string, the underscore on `_status`, the paging envelope. The titles, slugs and prose are this Recipe's own rather than Payload's blog, so the cases reading them are marked documentation-only.
- **`pagingCounter` and `hasPrevPage` are not served.** This format has names for eight of the nine paging fields; those two have none, and inventing constants for values that change per page would be worse than leaving them out. They are quoted above from the live response.
- **The delete is recorded, not served.** `DELETE /api/posts` is a write; this Recipe serves the two reads and states what the delete answers.
- **One collection of many.** `posts`. This deployment also serves `users`, `media`, `categories` and the rest of its own schema, and any Payload instance serves whatever its config defines.
- **Nothing is mapped in detection.** Payload is reached through its own REST API on whatever host the instance runs on, which does not resolve to any one host through a dependency file. Checked 2026-09-14.
