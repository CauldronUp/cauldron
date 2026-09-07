# browserbase

Emulates the Browserbase API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Browserbase's reference at `docs.browserbase.com` and struck live against `api.browserbase.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The refusal names the exact header you forgot.**

```
GET /v1/sessions    (no credential)
401 {"statusCode":401,"error":"Unauthorized","message":"Missing x-bb-api-key header"}
```

`x-bb-api-key` — the header, spelled out, in the message. That is the single most useful thing a refusal can say to somebody wiring up an integration for the first time, and almost nothing in this collection says it. [pulumi](../pulumi) refuses a `Bearer` prefix without ever mentioning that it wants `token`; [openphone](../openphone) names its header and is the only other one here that does.

It matters more than usual for this API, because **the header is not `Authorization`**. A developer reaching for the default of every HTTP library gets nothing back, and the refusal tells them why in five words.

**A wrong key is answered with the word `Unauthorized` twice.** `error` and `message` both. So the precision is only on the first failure: once a key is present and wrong, there is nothing to branch on but a string that is also in the other field.

**The 404 has the same three fields in the opposite order.**

```
401 {"statusCode":401,"error":"Unauthorized","message":"Missing x-bb-api-key header"}
404 {"message":"Route GET:/v1/cauldron-nope not found","error":"Not Found","statusCode":404}
```

Fastify's default not-found handler building the object one way and the application building it the other. Same keys, same shape, and only a byte comparison notices — which is every recorded fixture and every snapshot test. [zuora](../zuora) has the same instability between its two credential failures, and [infisical](../infisical)'s 404 is this same Fastify handler answering from below the application.

**The 404 echoes the method and the path**, which is Fastify again and is genuinely useful: a caller sees exactly what the router was given, verb and all.

**Success has no envelope.** A bare array, where every failure has three fields — so the shape of a response depends on whether it worked.

**A running session has no end time and is already billing.** `proxyBytes` counts on a session with no `endedAt`, so a cost computed over finished sessions is lower than the bill.

**The region is per session rather than per project.** Two sessions in one project can be on two continents, so latency and egress vary within one listing and nothing groups them.

## Modelling limits

- **One route.** Sessions. Contexts, extensions, uploads, logs, live URLs and the whole CDP surface each want their own evidence.
- **Nothing is mapped in detection.** Browserbase is driven through Playwright or Puppeteer pointed at its CDP endpoint, so a project that uses it depends on a browser-automation library rather than on a client of this API — which no dependency name can distinguish.
- **No `spec:`.** Browserbase publishes a rendered reference site.
