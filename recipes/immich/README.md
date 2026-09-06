# immich

Emulates the Immich photo-library API, for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Immich's own generated OpenAPI 3.0.0 document, published in its own repository — `immich-app/immich`, `open-api/immich-openapi-specs.json`, 191 paths, version 3.2.0-rc.0 — read on 2026-09-06. Immich is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**The schema announces a null it has not sent yet.** Immich ships a vendor extension, `x-immich-history`, giving **141 fields** a per-version lifecycle. On an album's `description`:

```json
[{"version": "v1", "state": "Added"},
 {"version": "v3", "state": "Updated",
  "description": "An empty string is returned instead of null for
                  backwards compatibility; null will be returned in v4."}]
```

So today `description` is `""` where it means absent, and in v4 it will be `null`. Both are falsy, so `if (album.description)` behaves identically across that change and `album.description.length` throws on one side of it. The warning is machine-readable, sits in the shipped document — and no standard OpenAPI tool reads it, because `x-` extensions are ignored by specification.

**128 fields are deprecated, and four of them only in that extension.** 124 carry OpenAPI's own `deprecated: true` beside `x-immich-state: Deprecated`, so a generator warns. Four do not:

| field | deprecated at |
|---|---|
| `ApiKeyCreateResponseDto.apiKey` | v3.2.0 |
| `AssetResponseDto.libraryId` | v1 — the version it was added in |
| `AssetResponseDto.resized` | v1.113.0 |
| `SearchAssetResponseDto.total` | v3.0.0 |

The first has a mechanical excuse: it is a `$ref` with siblings, and OpenAPI 3.0 ignores keys beside a `$ref`, so `deprecated` there would have done nothing either. The other three are inline schemas that simply lack it. All four are invisible to every tool that reads the standard keyword — and the first is the field carrying a **newly created API key**.

`libraryId` was added and deprecated in the same version. `SearchAssetResponseDto.total` has a `Deprecated` entry and no `Added` entry at all, so the document does not say when the total on a search response began existing, only that it should not be used.

**An array whose index means something, and whose second element depends on who is asking.** `albumUsers`, in the document's own words:

> First entry is always the album owner. Second entry is the auth user, if it differs from the owner. The rest are ordered alphabetically.

Position 0 is the owner. Position 1 is *you* — unless you own the album, in which case position 1 is the alphabetically first of everyone else. The ordering carries information that exists nowhere else in the response and changes with the caller. Any client that sorts, filters, or renders these through a `Set` destroys it, with no field to recover it from.

**A cookie is a first-class documented credential.** Three schemes on every authenticated operation: `bearer` (a JWT), `api_key` (`x-api-key`), and `cookie` — `immich_access_token`, `in: cookie`. Browsers send cookies automatically, and the document mentions no CSRF token anywhere. Every state-changing endpoint in a 191-path API is documented as reachable with an ambient credential.

**A photo count is bounded by JavaScript.** `assetCount` and `SearchAssetResponseDto.total` both declare `maximum: 9007199254740991` — `Number.MAX_SAFE_INTEGER`, the ceiling of the language the server is written in, published as a property of the data.

**A listing of albums takes no parameter that would bound it.** `GET /albums` takes `assetId`, `id`, `isOwned`, `isShared` and `name`. No limit, no page, no cursor, no total. An account with five thousand albums answers five thousand albums, every time.

**A missing album is a 400, and so is one you may not read** — one status and one message, `"Not found or no album.read access"`, covering both. A client cannot tell a permission problem from a typo.

**A thumbnail is not implied by a non-zero count.** `albumThumbnailAssetId` is nullable and required, so it is always present and may be null on an album holding assets.

**The one public endpoint answers under a key called `res`.** `GET /server/ping` is the only operation in 191 paths declaring no `security` at all, and it answers `{"res": "pong"}`.

## Modelling limits

- **Nothing here is verified against a live API.** Immich is self-hosted and there is no public instance.
- **Albums and the ping.** 191 paths is a photo library: assets, people, faces, search, timeline, sharing, libraries, jobs and the upload surface each want their own evidence.
- **`description` is served as the empty string**, which is what v3 sends. The v4 null is recorded and not served — serving a value the provider does not send yet would be inventing a future.
- **The cookie scheme is recorded and not served.** Accepting an ambient cookie in local development would teach exactly the habit the finding warns about.
- **`albumUsers` is served in the documented order for the fixture's own viewer.** The Recipe cannot vary it by caller, because it does not know who is calling beyond the credential — which is the finding restated.
