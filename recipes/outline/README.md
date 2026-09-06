# outline

Emulates the Outline wiki API, for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-06.**

The live cases were struck against `app.getoutline.com` with no account: every one of them is a refusal, and being refused is the thing being recorded. Written against Outline's own OpenAPI 3.0 document — `outline/openapi`, `spec3.yml`, 154 paths.

## What this Recipe found

**154 paths. 154 POSTs. No GET anywhere.** Every operation is a POST, including `documents.list`, `documents.info`, `documents.search` and `apiKeys.list`. The paths are RPC method names — `/documents.info`, not `/documents/{id}` — so the resource and the verb share one path segment, separated by a dot.

What that costs, none of which any single response reveals:

- **Nothing is cacheable.** A read is a POST, so no HTTP cache, no CDN, no conditional request, no ETag flow, no browser cache.
- **Nothing is idempotent by method.** A client retrying after a timeout cannot tell from the method whether a repeat is safe; `documents.list` and `documents.delete` are the same verb.
- **Every cross-origin call is preflighted**, because a POST with a JSON content type is never a simple request.
- **A read's parameters are in a JSON body**, so a query cannot be a link, a bookmark, or a line in an access log.

**And the one mistake this design invites is reported as a missing endpoint.** Struck live:

```
GET  /api/documents.list   -> 404
     {"ok":false,"error":"not_found","status":404,"message":"Resource not found"}

POST /api/nosuch.method    -> 404
     {"ok":false,"error":"not_found","status":404,"message":"Resource not found"}
```

Byte-identical, with no `Allow` header on either. So in an API where every endpoint is a POST, a developer who reaches for GET on a listing — the single likeliest mistake there is — is told the endpoint does not exist. The URL is right. The response says it is not.

**The envelope carries the status twice and a boolean besides.** Struck live: `{"ok": false, "error": "authentication_required", "status": 401, "message": "Authentication required"}`. `ok` is a boolean, `status` repeats the HTTP status, `error` is a machine code, `message` is prose. Four ways to ask the same question.

**And the machine code cannot tell two failures apart.** Struck live, one request after the other:

| what was sent | `error` | `message` |
|---|---|---|
| no credential | `authentication_required` | "Authentication required" |
| malformed token | `authentication_required` | "Unable to decode token" |

Same code, same status. The only thing distinguishing "you sent nothing" from "your token is broken" is the prose — the part that changes between releases and the part nobody should branch on.

**Paging travels in the body.** `offset` and `limit` are properties of the POST body, not query parameters, so a listing's position never appears in a URL, a log line, or a proxy's record of the request.

**Permissions travel beside the records.** A listing answers `{data: [...], policies: [...]}`, where a `Policy` carries an `abilities` object declared with `additionalProperties` — so what a caller may do is untyped data, keyed by ability name, arriving alongside the records rather than being implied by them, and a client has to join it back itself.

**The new filter is exclusive with four deprecated ones.** `documents.list` takes `filters`, whose own description says it "cannot be combined with the deprecated `collectionId`, `userId`, `parentDocumentId` or `statusFilter` parameters". Five ways to narrow a listing, four deprecated, and mixing them is an error rather than a merge.

**Three nulls on one record mean three different things.** `publishedAt`, `archivedAt` and `deletedAt` are separate nullable timestamps; a never-published draft has all three null, and only their names say which state it is in.

**No rate-limit headers on anything struck** — no `RateLimit`, no `Retry-After`, on a 401 or a 404.

## Modelling limits

- **Documents, listed and fetched.** 154 paths is a wiki: collections, users, groups, shares, comments, revisions, attachments, search and the whole OAuth surface each want their own evidence.
- **The GET-is-404 behaviour is served, because it is the finding.** A GET route declared here has to answer something, and what Outline answers is the not-found body — so this Recipe answers it too.
- **`policies` is served as a fixed array beside the data.** Its abilities are the fixture's, not computed from the credential, because this Recipe has one caller.
- **The deprecated filters are recorded and not served.** Serving them would mean deciding what happens when they are combined with `filters`, and the document says only that they cannot be.
