# codat

Emulates the Codat platform API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Codat's reference at `docs.codat.io` and struck live against `api.codat.io` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The failure tells you whether to retry, and the answer is "Unknown".**

```
GET /companies         (no credential)
401 {"statusCode":401,"service":"PublicApi","error":"Unauthorized",
     "correlationId":"3a94fab4…","canBeRetried":"Unknown","detailedErrorCode":0}
```

`canBeRetried` is a **tri-state sent as a string**, and it is `"Unknown"` on every failure a public probe can reach. The field that exists to answer the one question a client has about a failure declines to answer it — in a type that has to be string-compared rather than tested as a boolean.

It is still better than not having it, and the three-valued design is right: yes, no, and "I do not know" are genuinely different answers, and most providers here offer none of them. The problem is that the third is the only one reachable without an account, so a client written against this API has an untested branch for the other two.

**`detailedErrorCode` is `0` on everything** — a second code beside `error`, reserved for detail, carrying none.

**`service` names the deployment that answered.** `"PublicApi"`, on a request with no account behind it. Mild, and the same class of thing [gong](../gong) does with a build number and [imgix](../imgix) with a version.

**The 404 keeps every field.** Six of them, one word different, and the correlation id is still there — which is exactly what [confluent](../confluent) gets wrong on the same failure, one Recipe over.

**There is no message field anywhere.** The only prose about a failure is the HTTP status phrase in PascalCase, in a field called `error`.

**Timestamps carry seven decimal places.** `2026-09-06T03:00:12.4410000Z` — .NET's round-trip format, hundred-nanosecond ticks. A parser expecting milliseconds either truncates or refuses, and both readings look like a working client.

**A company can exist before anything is connected to it.** An empty `platform`, no `dataConnections`, no `lastSync`, and a redirect link waiting — so a listing of companies is not a listing of data sources, and the count a dashboard shows is the wrong one.

## Modelling limits

- **One route.** Companies. Connections, the whole accounting/banking/commerce data model, sync status, webhooks and the assess surface each want their own evidence.
- **Nothing is mapped in detection.** Codat's SDKs are generated per data type and per language and none of them resolves under an obvious registry name; nothing found on 2026-09-07 calls `api.codat.io`.
- **No `spec:`.** Codat publishes a rendered reference site.
