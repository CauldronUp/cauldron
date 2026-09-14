# builderio

Emulates the Builder.io Content API listing for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-13.**

Struck live on 2026-09-13 with no key, with a wrong key, on a path that does not exist, with a method the path does not take, on a model that does not exist, on three API versions, and — for the record shapes — with the public demo key Builder publishes in its own quickstart.

## What this Recipe found

**One API, three error envelopes.**

```
no key, or a wrong one
  401 {"status":401,"message":"Authorization required","detail":"For help please contact support@builder.io"}

a model that does not exist
  404 {"status":404,"message":"Model not found"}

a path that does not exist, or a method it does not take
  404 {"errors":[{"status":"404","title":"Route Does Not Exist","detail":"Cannot GET `/api/v3/cauldron-nope`"}]}
```

Three shapes: a bare object with `message` and `detail`, the same object without `detail`, and an array under `errors` whose entries use `title` instead of `message`. And `status` is the **number** `401` in the first two and the **string** `"404"` in the third — so code comparing that field has to know which failure it is looking at before it can read the field that says which failure it is.

**Every mistake made without a key is the same 401.** No key, a wrong key, an unrouted path, a wrong method and an API version that serves something else entirely all answer "Authorization required", with a support email address in the machine-readable `detail`. The other two envelopes appear only once a valid key is in the URL — so a developer debugging a typo'd path sees a credential error until the credential is right.

**`/api/v1/content/{model}` answers 200 with Builder's marketing homepage.** With a valid key:

```
GET /api/v1/content/page   200  text/html          <title>Builder.io: Collaborative Development Platform</title>
GET /api/v2/content/page   200  application/json
GET /api/v3/content/page   200  application/json
```

A client pinned to v1 gets a web page with a success status, and `.json()` throws on it.

**Two fields that read as booleans are strings.** `published` is the string `"published"`. `data.hidden` is the string `"false"`. So `if (item.published)` is true whatever the state, and `if (item.data.hidden)` is true for a page that is not hidden.

**The envelope has no count, no cursor and no end.** `{"results":[…]}` and nothing else. Paging is `limit` and `offset`; an offset past the end answers `{"results":[]}`; and nothing in a full page says whether another one exists.

**`rev` is the same on unrelated records.** Two pages created eleven months apart, by different users, in different states, both carry `"rev":"taqizx9mb0d"` — a field named for a revision that does not vary with the record it is on.

**A publicly readable record carries a preview URL with a permission grant in it.** `meta.lastPreviewUrl` is a link whose query string includes `builder.user.permissions=read%2Ccreate%2Cpublish%2CeditCode%2CeditDesigns%2Cadmin%2C…` — percent-encoded, on content that anyone holding the public read key can fetch.

**The same asset is addressed two ways in one response.** `screenshot` points at `…/api/v1/image/assets%2F<key>%2F<id>` while an image inside the page body points at `…/api/v1/image/assets/<key>/<id>` — the same path shape, the slashes encoded in one and bare in the other.

Also pinned: `createdBy` and `lastUpdatedBy` are 28-character Firebase auth uids, served to anyone with the public read key; `query` carries the record's own targeting rules as objects tagged `"@type": "@builder.io/core:Query"`, and the blocks inside `data` are tagged `"@type": "@builder.io/sdk:Element"` with an `@version` — which look like JSON-LD and are not; `testRatio` and `variations` put A/B state on every record whether or not it has a test; and the credential is a query parameter by design, because it is meant to ship inside a browser.

## Sources

- Live: `cdn.builder.io`, struck 2026-09-13 with no key, a wrong key, an unrouted path, a wrong method, a model that does not exist, and `/api/v1`, `/api/v2` and `/api/v3` of the content path.
- The record shapes were read from `cdn.builder.io/api/v3/content/page` using `YJIGb4i01jvw0SRdL5Bt`, the read-only key Builder publishes in its own quickstart. It is a public key by design: Builder's client SDKs ship it in browser JavaScript.
- [Builder.io Content API documentation](https://www.builder.io/c/docs/content-api).

## Modelling limits

- **No description is published.** Builder.io serves no OpenAPI document for the content API, so there is no `spec` to fingerprint.
- **The fixture values are synthetic; the shapes are not.** Every field, type and quirk here was observed on a live response — the string `"published"`, the string `"false"`, the shared `rev`, the permissions in `meta.lastPreviewUrl`, the encoded `screenshot` path. The titles, urls and ids in the fixture are this Recipe's own, so the cases reading them are marked documentation-only rather than claiming a live match. The Firebase uids and the `rev` value are the ones observed, because their shape is the finding.
- **The v1 page is one line of a real one.** Live it is Builder's full marketing homepage; this Recipe serves a minimal document carrying the same `<title>`, the same status and the same content type.
- **One model of many.** `page` is served; every other model name answers `Model not found`, which is what the live API does for a model that is not in the space.
- **No single-record route.** `/api/v3/content/{model}/{id}` was not struck and is not served.
- **Nothing is mapped in detection.** Builder.io is reached through `@builder.io/sdk` or a plain URL carrying an `apiKey` parameter, and neither resolves to this host through a dependency file. Checked 2026-09-13.
