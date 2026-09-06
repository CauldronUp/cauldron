# overseerr

Emulates the Overseerr media-request API, for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Overseerr's own OpenAPI 3.0.2 document, published in its own repository — `sct/overseerr`, `overseerr-api.yml`, 134 paths, version 1.0.0 — read on 2026-09-06. Overseerr is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**Two fields called `status`, both integers, overlapping and meaning different things.** A request carries one, and the media it points at carries another:

| field | values |
|---|---|
| `MediaRequest.status` | 1 = PENDING APPROVAL, 2 = APPROVED, 3 = DECLINED |
| `MediaInfo.status` | 1 = UNKNOWN, 2 = PENDING, 3 = PROCESSING, 4 = PARTIALLY_AVAILABLE, 5 = AVAILABLE, 6 = DELETED |

The second is nested *inside* the first: `request.media.status`. So a request with `status: 2` — approved — routinely contains `media.status: 2`, which is pending. Same name, same type, same small integers, one inside the other, and a client that reads the wrong one gets a plausible answer rather than an error.

`status: 3` is **DECLINED** on the request and **PROCESSING** on the media. One is a refusal; the other is work in progress.

**And both examples are a value neither documents.** `MediaRequest.status` is declared `example: 0` against a description listing 1, 2 and 3. `MediaInfo.status` is declared `example: 0` against a description listing 1 through 6. Zero is in neither set.

So the examples in the published document — the values a mock server, a documentation page or a generated fixture will use — put every record in a state that does not exist. And zero is falsy, so `if (request.status)` treats the example as "no status" rather than as "impossible".

**The filter is a string vocabulary that matches neither integer set.** `GET /request` takes `filter`: `all`, `approved`, `available`, `pending`, `processing`, `unavailable`, `failed`, `deleted`, `completed`.

`approved` is a request status. `pending` is *both* a request status and a media status, meaning different things. `available`, `processing` and `deleted` are media statuses. And `failed`, `completed` and `unavailable` appear in neither list — so three of the words you can filter by name nothing the records can hold.

**The same word is the array and the count of it:**

```json
{"pageInfo": {"page": 1, "pages": 10, "results": 100},
 "results": [ … ]}
```

`body.results` is the records; `body.pageInfo.results` is how many there are. One word, two meanings, four characters apart in the same response.

**Paging is `take` and `skip`** — TypeORM's names straight through to the wire, not `limit`/`offset` or `page`/`per_page`. Both are declared **nullable**, so the document says `null` is a legal value for how many records you want.

**A 4K request is a separate record, not a flag.** The same film requested in 4K is a different row with its own id, status and root folder. A client that keys requests by `tmdbId` collapses them and loses one.

**Three identifiers in one response** — the request's own, Overseerr's internal media id, and TMDB's — and only the first addresses the endpoint. `tvdbId` is null on a film, so the field's absence is a genre rather than a gap.

**`serverId: 0` is a real server.** Overseerr numbers its Radarr and Sonarr instances from nought, so the falsy value is the first configured server rather than none.

**A session cookie is a documented API credential.** `connect.sid`, `in: cookie`, beside `X-Api-Key` — the second in this collection after [immich](../immich), with the same consequence: a browser sends it automatically, and the document mentions no CSRF token across 134 paths.

## Modelling limits

- **Nothing here is verified against a live API.** Overseerr is self-hosted and there is no public instance.
- **Requests, listed and fetched.** 134 paths is a request manager: media, users, settings, discovery, the Plex/Radarr/Sonarr integrations and the whole search surface each want their own evidence.
- **The cookie scheme is recorded and not served.** Accepting an ambient credential in local development teaches exactly the habit the finding warns about.
- **`example: 0` is recorded and not served.** Serving a state the provider documents as impossible would be inventing behaviour; the fixtures hold values from the documented sets.
