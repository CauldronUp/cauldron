# karakeep

Emulates the Karakeep bookmark-manager API (v1), for local development and tests.

**12 conformance cases, 2 checked against the live API on 2026-09-06.**

The live cases were struck against `try.karakeep.app` with no account — both are refusals. Written against Karakeep's own OpenAPI 3.0 document, published in its own repository: `karakeep-app/karakeep`, `packages/open-api/karakeep-openapi-spec.json`, 35 paths.

## What this Recipe found

**The only response that is not JSON is the one you most need to read.** Struck live:

```
GET /api/v1/bookmarks?limit=1        (no credential)
HTTP 401  Content-Type: text/plain;charset=UTF-8
Unauthorized
```

Twelve bytes of prose. Every success in this API is `application/json`, and the document declares an `Error` schema with `code` and `message` both required — which the 401 does not use. Its declared content type in the document is `text/plain` too, so this is intended rather than accidental.

A client that calls `.json()` on every response throws here, on the one response whose meaning it most needs, and what it reports is a JSON parse error at position 0 — not "your credential is wrong".

**And the two refusals are byte-identical.** Sending nothing and sending a token that will not decode produce the same status, the same content type and the same twelve bytes. Nothing distinguishes a configuration problem from an expiry.

**Null is a member of the enum, not the absence of one.** A bookmark carries three status fields — `taggingStatus`, `summarizationStatus`, `embeddingStatus` — each declared:

```json
"enum": ["success", "failure", "pending", null]
```

So `null` is a fourth declared value meaning *never attempted*, not a missing field. Four states, three independent copies of them on one record, and the fourth is the one every `if (status)` and `status ?? "…"` quietly collapses into the others. `source` is the same shape: eight members plus `null`, where null means the bookmark predates the field.

**A bookmark's payload is a `oneOf`, and one variant has no payload:**

| `content.type` | fields |
|---|---|
| `link` | 23 — url, title, description, favicon, htmlContent, crawlStatus, readerViewScore, author, publisher, datePublished… |
| `text` | 3 |
| `asset` | 7 |
| `unknown` | **1** — its own name |

`unknown` is a declared variant carrying nothing but `type`. An exhaustive switch over the four has to handle a case with no data in it, and `bookmark.content.url` is `undefined` on three of the four. (`assets[].assetType` has twelve members and one of *those* is also `unknown`.)

**Tags say whether a machine wrote them.** Each tag carries `attachedBy`, `"ai"` or `"human"`, both required. So the tags a model invented are marked — and only if you look. A client that renders `tags` uniformly presents guesses and decisions identically; one that syncs them into another system copies the guesses in as facts.

**Two creation timestamps, and only one is required.** `createdAt` is in the schema's `required` list, `firstCreatedAt` is not, and nothing on either says how they differ.

**The cursor is required *and* nullable,** so the key is always present and only its value moves — which makes `"nextCursor" in body` useless and `body.nextCursor === null` the only correct end-of-listing test.

**`favourited` carries the British spelling** — both the field and the query parameter — so a client written to the American spelling filters nothing and reads `undefined`.

## Modelling limits

- **Bookmarks, listed and fetched.** 35 paths is a bookmark manager: lists, tags, highlights, assets, feeds and the whole backup surface each want their own evidence.
- **`content` is served as the `link` variant.** The other three are recorded and not served — serving a `oneOf` would mean deciding per fixture which variant each record is, and what the finding needs is that the `unknown` variant exists at all.
- **The two refusals are served identically,** because that is what the live API does. This Recipe therefore cannot tell a caller which mistake they made either — which is the finding, not a gap.
