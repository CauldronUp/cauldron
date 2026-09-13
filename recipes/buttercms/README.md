# buttercms

Emulates the ButterCMS blog posts API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document ButterCMS serves at [`read_api.yaml`](https://buttercms.com/docs/api/openapi/read_api.yaml), and struck live against `api.buttercms.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**Three failures, one key, and the bytes differ.**

```
(no token)       401 {"detail":"Incorrect authentication credentials."}
?auth_token=x    401 {"detail":"Invalid token."}
/cauldron-nope/  404 {"detail": "Not found."}
```

Read the third one carefully: there is a space after the colon, and there is not one on the first two. Checked at the byte level, twice.

Two different JSON serialisers are answering from the same API, and the only way to see it is to look at the bytes rather than at the parsed object — which means a diff-based or hash-based test over recorded responses breaks on it and nothing else ever will.

**A request with no credential is told its credentials are incorrect.** "Incorrect authentication credentials." answers a request that supplied none; "Invalid token." answers one that supplied a wrong one. The vaguer sentence goes to the case that could have been named exactly, and the precise one goes to the case a caller can already guess.

**Four pagination parameters on one operation, mutually exclusive in pairs.** `page` and `page_size`, `limit` and `offset`. Each carries the same warning in its own description — "**Mutually exclusive with** the other pair" — which is the only place the rule exists. Nothing in the schema expresses it, so a generated client offers all four on one call.

**And the envelope changes shape depending on which pair you sent.** The document declares two response types for this one operation:

| response type | `meta` |
| --- | --- |
| `PageBasedPostsResponse` | `{previous_page, next_page, count}` |
| `OffsetBasedPostsResponse` | `{count, next_offset, previous_offset}` |

A client has to remember how it asked in order to know how to read the reply, and `meta.next_page` is present or absent according to a decision made in the request.

**An out-of-range page size is silently corrected.** From `limit`'s own description:

> Values above 100 are capped at 100; values below 1 or invalid fall back to the default of 10 (the request is not rejected).

A caller who asks for 500 records gets 100, a caller who asks for 0 gets 10, and neither is told.

**The token travels in the query string or in a header, and the document declares both.** Two security schemes for one credential — and the query form is the one every example in the documentation uses, so the secret is in the URL by default.

**Publication state is in three fields.** `status` is an enum of draft, published and scheduled, with `published` and `scheduled` as nullable timestamps beside it. Nothing says what `status: scheduled` means when `scheduled` is null.

**And nothing on a post is required.** Seventeen properties, no `required` array — not `slug`, not `title`, not `url`.

## Sources

- Live: `api.buttercms.com`, struck 2026-09-13.
- [`read_api.yaml`](https://buttercms.com/docs/api/openapi/read_api.yaml) — 215KB, served without a credential.
- [List All Posts](https://buttercms.com/docs/api-reference/blog-posts/list-all-posts) — the page, with the same document embedded in it.

## Modelling limits

- **One route.** Listing blog posts. Pages, collections, authors, categories, tags, search, the feeds and the Write API each want their own evidence.
- **The whitespace difference is recorded, not served.** Cauldron serialises every body the same way, so the 404 here is byte-comparable with the 401s and live it is not. The difference is the finding; reproducing it would mean declaring a serialiser per error, which no Recipe can do and no client should depend on.
- **`meta.next_page` is a string here and an integer in the document.** Cauldron renders a cursor as a string, and this provider's cursor is a page number. The case asserts its shape rather than its type.
- **Offset pagination is not modelled.** The other pair of parameters returns the other `meta` shape, and Cauldron declares one envelope per route.
- **Nothing is mapped in detection.** The `buttercms` packages on npm and Packagist are real clients of this API, but a project holding one is pointed at whichever ButterCMS account its read token belongs to — and the same package also talks to the Write API on a different host. Checked 2026-09-13.
