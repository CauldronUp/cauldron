# gumlet

Emulates the Gumlet video asset listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Gumlet publishes at [`docs.gumlet.com/reference/openapi.json`](https://docs.gumlet.com/reference/openapi.json), and struck live against `api.gumlet.com` on 2026-09-13 with no header, with a wrong bearer, with a Basic header, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**"Invalid API Bearer token" is what you are told for sending none.** With no `Authorization` header at all:

```json
{"error":{"code":"invalid_authorization_header","message":"Invalid API Bearer token. Authorization header must be of format 'Bearer YOUR_API_KEY'"}}
```

A token is called invalid when there was no token. A wrong one gets a different code — `invalid_api_key`, "API key supplied with request is invalid" — so the API does distinguish the two cases, and the sentence for the first one does not.

**And that is the answer to everything.** An unrouted path, a method the route does not take, and a Basic header all receive the same 401, because the credential is checked before the path is. There is no way to find out whether a path exists.

**The published example response is somebody's real workspace.** The 200 example for this operation carries `total_asset_count: 2159` and a `distinct_tags` array of that account's whole tag vocabulary — `"aasas"`, `"asasaasas"`, `"asasas"`, `"asasasasa"`, `"asasass"`, `"asasdfdfdd"`, `"asqs"`, `"dsddsd"`, `"rewew"`, `"sdsdsdsd"` — alongside working labels and two personal names. Keyboard mash from a live testing session, shipped as the worked example in the public reference.

**And `distinct_tags` is on the listing, not the record.** Every tag in the workspace arrives with every page of assets, so paging through two thousand assets fetches the same vocabulary each time.

**The example's status is not one the filter can ask for.** The `status` query parameter is enumerated `queued`, `processing`, `ready`, `errored`, `deleted`. The asset in the document's own example has `"status": "optimizing"` — a sixth value, which arrives in responses and cannot be filtered on.

**One path parameter, two spellings.** `/video/assets/{asset_id}` is a GET and a DELETE; `/video/assets/{asset_ID}/thumbnail`, `/video/assets/{asset_ID}/subtitle/upload` and `/video/assets/{asset_ID}/audio/upload` are POSTs. The same identifier, cased two ways in one document, so a generated client has two variable names for it.

**There is a path in the document that does nothing.** `/video/assets/{asset_ID}/audio/upload-1` declares no operations at all — an empty path item with a de-duplication suffix in its name, published as part of the API.

**A record carries a URL back to itself.** `output.status_url` is `https://api.gumlet.com/v1/video/assets/6192269e0822a81d955d1a4b` — the address of the thing you are already holding.

**The security scheme is named `API_KEY` and declared `{"type": "http", "scheme": "bearer"}`.** `info.version` is `1.4` while the server URL says `/v1`; the listing's own description begins `[Deprecated]`; and `current_offset` is `1` on the first page, which is not what an offset is.

## Sources

- [`docs.gumlet.com/reference/openapi.json`](https://docs.gumlet.com/reference/openapi.json) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.gumlet.com`, struck 2026-09-13 with no header, a wrong bearer, a Basic header, an unrouted path, and a wrong method.

## Modelling limits

- **One route of fifty-six.** The deprecated asset listing. Images, live streams, playlists, profiles, workspaces, folders, insights, webhooks and the multipart upload surface are the rest.
- **`distinct_tags` is a short list here.** Live it is the account's whole vocabulary; the fixture carries five entries rather than reproducing the published example's, which are one workspace's real labels.
- **The record shape is the document's example, not a schema.** This operation declares its 200 as an `examples` block with no schema beside it, so what the fixture carries is the example's own field set.
- **Nothing is mapped in detection.** Gumlet is reached through a plain HTTP call carrying a bearer token or through its image URL rewriting, and neither resolves to this host through a dependency file. Checked 2026-09-13.
