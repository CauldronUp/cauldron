# sonarr

Emulates the Sonarr API (v3), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Sonarr's own generated OpenAPI document, published in its own repository — `Sonarr/Sonarr`, `src/Sonarr.Api.V3/openapi.json`, 162 paths, `info.version` 3.0.0 — read on 2026-09-06. Sonarr is self-hosted, so there is no public instance to strike.

**This Recipe exists in a pair with [radarr](../radarr), and the pair is the point.**

## What this Recipe found

**Sonarr and Radarr are the same API until they are not.** Same fork, same framework, same `/api/v3` prefix, same `info.version` string of `3.0.0`, same two credential schemes with the same names. **134 paths appear in both documents.** A client written for one and pointed at the other connects, authenticates, and answers 200.

Thirteen of those shared paths take different parameters. `GET /api/v3/history` takes `seriesIds` and `includeSeries` here, and `movieIds` and `includeMovie` on [radarr](../radarr). `GET /api/v3/calendar` differs by four parameters. The framework is ASP.NET Core, whose model binding **drops a query key with no matching property rather than refusing it** — so the wrong parameter is not an error, it is a filter that does not filter, and the response is every record instead of none.

**And the same integer means different things on the two.** `eventType` is a filter taking an array of integers, and a response field holding a string. The two documents list their members in this order:

| position | Sonarr | Radarr |
|---|---|---|
| 0 | `unknown` | `unknown` |
| 1 | `grabbed` | `grabbed` |
| 2 | `seriesFolderImported` | `downloadFolderImported` |
| 3 | `downloadFolderImported` | `downloadFailed` |
| 4 | `downloadFailed` | `movieFileDeleted` |
| 5 | `episodeFileDeleted` | `movieFolderImported` |
| 6 | `episodeFileRenamed` | `movieFileRenamed` |
| 7 | `downloadIgnored` | `downloadIgnored` |

Only 0, 1 and 7 agree. Neither document pins explicit values, so the integers a caller sends are the members' declaration positions — which means **`?eventType=3` asks for successful imports here and failed downloads on Radarr**. A monitoring job counting failures against the wrong one of the pair counts successes, and nothing anywhere returns an error.

What tells them apart is `appName` on `GET /api/v3/system/status`, and nothing else: not the path, not the version, not the credential, not the shape of a listing. `instanceName` looks like a second discriminator and is not — it is settable, and an operator running both may rename them to the same thing.

**The same field is an integer going out and a string coming back.** You filter with `eventType=3` and read back `"eventType": "downloadFolderImported"`. One name, one endpoint, two types, and no client type round-trips it.

**The credential is documented as a query parameter as well as a header.** Both schemes are declared: `X-Api-Key` in a header, and `apikey` in the query string, described in the document as "Apikey passed as query parameter". A credential in a URL reaches access logs, proxy logs, browser history and `Referer` headers — and these are self-hosted boxes commonly sat behind a reverse proxy that logs full URLs by default.

**The status endpoint is a host inventory.** `SystemResource` carries `startupPath`, `appData`, `osName`, `osVersion`, `isDocker`, `sqliteVersion`, `migrationVersion`, `runtimeVersion`, `branch` and `packageUpdateMechanism`. Absolute filesystem paths and host details, behind a credential that the same document says may travel in the query string.

## Modelling limits

- **Nothing here is verified against a live API.** Sonarr is self-hosted and there is no public instance; every claim is the generated document's.
- **History and system status only.** 162 paths is a media manager: series, episodes, files, queue, indexers, download clients, import lists, notifications, quality profiles and the command endpoint each want their own evidence.
- **The `eventType` filter is served over the string, not the integers.** Encoding the integer mapping would mean asserting the enum's declaration positions are its wire values, which neither document states — it is inferred from the enum carrying no explicit values, and the inference is written down rather than encoded.
- **The credential is checked in the header.** The query-string scheme is recorded and not served, because serving it would make a credential in a URL the easy path in local development.
