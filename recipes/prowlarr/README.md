# prowlarr

Emulates the Prowlarr API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Prowlarr's own generated OpenAPI document, published in its own repository — `Prowlarr/Prowlarr`, `src/Prowlarr.Api.V1/openapi.json`, 93 paths, `info.version` 1.0.0, OpenAPI 3.0.4 — read on 2026-09-06. Prowlarr is self-hosted and there is no public instance.

**This Recipe is one of five siblings** — [sonarr](../sonarr), [radarr](../radarr), [prowlarr](../prowlarr), [lidarr](../lidarr), [readarr](../readarr). They are forks of one codebase, commonly installed on the same machine, sharing a credential header, a paged envelope and an API shape. **The set is the finding.**

## What this Recipe found

**The same integer means five different things.** `eventType` is a filter taking integers and a response field holding a string. None of the five documents pins explicit values, so the integers a caller sends are the members' declaration positions:

| pos | sonarr | radarr | prowlarr | lidarr | readarr |
|---|---|---|---|---|---|
| 0 | `unknown` | `unknown` | `unknown` | `unknown` | `unknown` |
| 1 | `grabbed` | `grabbed` | **`releaseGrabbed`** | `grabbed` | `grabbed` |
| 2 | `seriesFolderImported` | `downloadFolderImported` | `indexerQuery` | `artistFolderImported` | `bookFileImported` |
| 3 | `downloadFolderImported` | `downloadFailed` | `indexerRss` | `trackFileImported` | `downloadFailed` |
| 4 | `downloadFailed` | `movieFileDeleted` | `indexerAuth` | `downloadFailed` | `bookFileDeleted` |
| 5 | `episodeFileDeleted` | `movieFolderImported` | `indexerInfo` | `trackFileDeleted` | `bookFileRenamed` |
| 6 | `episodeFileRenamed` | `movieFileRenamed` | — | `trackFileRenamed` | `bookImportIncomplete` |
| 7 | `downloadIgnored` | `downloadIgnored` | — | `albumImportIncomplete` | `downloadImported` |
| 8 | — | — | — | `downloadImported` | `bookFileRetagged` |
| 9 | — | — | — | `trackFileRetagged` | `downloadIgnored` |
| 10 | — | — | — | `downloadIgnored` | — |

**Only position 0 agrees across all five.**

`?eventType=3` asks for a **successful import** on Sonarr and Lidarr, a **failed download** on Radarr and Readarr, and an **RSS poll** on Prowlarr. There is no integer that means "download failed" across the family: it is 4, 3, absent, 4 and 3.

A monitoring job pointed at the wrong one of five counts the wrong thing, and nothing anywhere returns an error.

**Nothing in a response tells the five apart except `appName`.** Same `X-Api-Key`, same paged envelope, same generated shape. `instanceName` looks like a second discriminator and is not — it is settable, and an operator running several may rename them to anything.

**The path prefix is not shared either, and that is one more thing that looks it.** Sonarr and Radarr serve `/api/v3`; Prowlarr, Lidarr and Readarr serve `/api/v1`. The number bears no relation to `info.version`, which is `"3.0.0"` for the two on v3 and `"1.0.0"` for the three on v1 — so a client deriving one from the other is right for the wrong reason on two forks and wrong on three.

The five documents also disagree about which OpenAPI patch they were generated against: 3.0.1, 3.0.4, 3.0.4, 3.0.4, 3.0.1.

**The credential is documented as a query parameter as well as a header**, on all five — `apikey`, "Apikey passed as query parameter", beside `X-Api-Key`. Recorded and not served, because serving it would make a credential in a URL the easy path in local development.

**The `eventType` filter is served over the string, not the integers.** Encoding the integer mapping would mean asserting the enum's declaration positions are its wire values, which none of the five documents states.

**Prowlarr is the one that is not a download manager at all.** It manages indexers. Its history is a log of indexer activity — queries, RSS polls, authentications — not of downloads, so four of its six event types have no counterpart in any sibling.

Its `HistoryResource` has **no `sourceTitle`, no `quality`, no `qualityCutoffNotMet`**; it has `indexerId` and a boolean `successful` that none of the other four carry. A client written against the shared shape reads `sourceTitle` and finds `undefined` on every Prowlarr record; one written against Prowlarr reads `successful` and finds `undefined` on the other four.

**And it renames the one event they all have.** Position 1 is `grabbed` on four siblings and `releaseGrabbed` here — same integer, same idea, different string. A client switching on the name handles four and falls through on the fifth.

## Modelling limits

- **Nothing here is verified against a live API.** Prowlarr is self-hosted and there is no public instance.
- **History and system status.** 93 paths is a media manager: indexers, applications, download clients, notifications, the search proxy and the indexer-stats surface each want their own evidence.
- **The credential is checked in the header.** The query-string scheme is recorded and not served.
