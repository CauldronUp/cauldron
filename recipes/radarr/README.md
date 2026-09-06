# radarr

Emulates the Radarr API (v3), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Radarr's own generated OpenAPI document, published in its own repository — `Radarr/Radarr`, `src/Radarr.Api.V3/openapi.json`, 164 paths, `info.version` 3.0.0 — read on 2026-09-06. Radarr is self-hosted, so there is no public instance to strike.

**This Recipe is one of five siblings** — [sonarr](../sonarr), [radarr](../radarr), [prowlarr](../prowlarr), [lidarr](../lidarr), [readarr](../readarr) — and the three added later show the same enum diverging five ways rather than two. It exists in a pair with [sonarr](../sonarr), and the pair is the point.**

## What this Recipe found

**Radarr and Sonarr are the same API until they are not.** Same fork, same framework, same `/api/v3` prefix, same `info.version` string of `3.0.0`, same two credential schemes with the same names. **134 paths appear in both documents.** A client written for one and pointed at the other connects, authenticates, and answers 200.

Thirteen of those shared paths take different parameters. `GET /api/v3/history` takes `movieIds` and `includeMovie` here, and `seriesIds` and `includeSeries` on [sonarr](../sonarr). `GET /api/v3/calendar` differs by four parameters. The framework is ASP.NET Core, whose model binding **drops a query key with no matching property rather than refusing it** — so the wrong parameter is not an error, it is a filter that does not filter, and the response is every record instead of none.

**And the same integer means different things on the two.** `eventType` is a filter taking an array of integers, and a response field holding a string. The two documents list their members in this order:

| position | Radarr | Sonarr |
|---|---|---|
| 0 | `unknown` | `unknown` |
| 1 | `grabbed` | `grabbed` |
| 2 | `downloadFolderImported` | `seriesFolderImported` |
| 3 | `downloadFailed` | `downloadFolderImported` |
| 4 | `movieFileDeleted` | `downloadFailed` |
| 5 | `movieFolderImported` | `episodeFileDeleted` |
| 6 | `movieFileRenamed` | `episodeFileRenamed` |
| 7 | `downloadIgnored` | `downloadIgnored` |

Only 0, 1 and 7 agree. Neither document pins explicit values, so the integers a caller sends are the members' declaration positions — which means **`?eventType=3` asks for failed downloads here and successful imports on Sonarr**. A monitoring job counting failures against the wrong one of the pair counts successes, and nothing anywhere returns an error.

What tells them apart is `appName` on `GET /api/v3/system/status`, and nothing else: not the path, not the version, not the credential, not the shape of a listing. `instanceName` looks like a second discriminator and is not — it is settable, and an operator running both may rename them to the same thing.

**The same field is an integer going out and a string coming back.** You filter with `eventType=3` and read back `"eventType": "downloadFailed"`. One name, one endpoint, two types, and no client type round-trips it.

**The credential is documented as a query parameter as well as a header.** Both schemes are declared: `X-Api-Key` in a header, and `apikey` in the query string, described in the document as "Apikey passed as query parameter". A credential in a URL reaches access logs, proxy logs, browser history and `Referer` headers — and these are self-hosted boxes commonly sat behind a reverse proxy that logs full URLs by default.

**The status endpoint is a host inventory.** `SystemResource` carries `startupPath`, `appData`, `osName`, `osVersion`, `isDocker`, `sqliteVersion`, `migrationVersion`, `runtimeVersion`, `branch` and `packageUpdateMechanism`. Absolute filesystem paths and host details, behind a credential that the same document says may travel in the query string.

**One more difference the pair makes visible:** this document declares OpenAPI `3.0.4` and Sonarr's declares `3.0.1`, from the same generator on the same fork. A toolchain pinned to one patch of the specification reads one of them and refuses the other, and neither project's release notes mentions it.

**A history record carries one identifier where Sonarr carries two.** A series has episodes; a movie is the whole thing. So a record here cannot be narrowed below the title, and a client written against Sonarr's `episodeId` finds nothing to read.

## Modelling limits

- **Nothing here is verified against a live API.** Radarr is self-hosted and there is no public instance; every claim is the generated document's.
- **History and system status only.** 164 paths is a media manager: movies, files, queue, indexers, download clients, import lists, notifications, quality profiles, collections and the command endpoint each want their own evidence.
- **The `eventType` filter is served over the string, not the integers.** Encoding the integer mapping would mean asserting the enum's declaration positions are its wire values, which neither document states — it is inferred from the enum carrying no explicit values, and the inference is written down rather than encoded.
- **The credential is checked in the header.** The query-string scheme is recorded and not served, because serving it would make a credential in a URL the easy path in local development.
