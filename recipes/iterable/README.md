# iterable

Emulates the Iterable API for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Iterable's reference at `api.iterable.com/api/docs` and struck live against `api.iterable.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The refusal hands back the caller's IP address and one of Iterable's own.**

Struck live, twice, minutes apart from the same client:

```
GET /api/lists                  (no credential)
401 {"msg":"No API key found on request","code":"Unauthorized",
     "params":{"ip":"<the caller>,10.10.50.35","endpoint":"/api/lists"}}

GET /api/lists                  Api-Key: not-a-real-key
401 {"msg":"Invalid API key","code":"Unauthorized",
     "params":{"ip":"<the caller>,10.10.49.6","endpoint":"/api/lists"}}
```

`ip` is the `X-Forwarded-For` chain, joined with a comma and sent as one string. The first address is the caller's. **The second is inside Iterable** — `10.10.50.35` on one request and `10.10.49.6` on the next, from the same client minutes apart, which makes it the load-balancer node that happened to take the call.

So an unauthenticated 401 reports which of Iterable's internal hosts answered, and enough refusals map the range. No credential is needed for any of them.

It is also a header the caller controls. `X-Forwarded-For` is appended to by each hop, so the first element is whatever the client sent — and a client that sends one gets it echoed into a body it can read without authenticating.

**`params.endpoint` echoes the path**, which is useful in a log and means the 401 confirms which routes exist — except that it doesn't, because **an unknown path is HTML**. 404, `text/html`, a full application page. The JSON 401 above is reachable only on paths that do exist.

**The prose field is `msg` and the field named `code` holds a title.** Three letters with the vowel taken out, beside `"Unauthorized"` in the field a client would branch on — and that value is identical on both failures, so the machine-readable half carries no information and the difference lives in the prose.

**The credential header is `Api-Key`.** Not `Authorization`, not `X-Api-Key`. One capital, one hyphen. HTTP header names are case-insensitive so it makes no difference on the wire, and every code sample and every reader's memory carries the spelling anyway.

**Creation times are milliseconds as a bare Number.** Thirteen digits, no string, no timezone — a client dividing by 1000 lands in 1970 and the result still looks like a date.

**A `Dynamic` list sits beside a `Standard` one** in the same array. Membership of the first is recomputed and of the second is stored, so "remove this subscriber" means two different things depending on one word in the record.

## Detection

`@iterable/web-sdk` names `api.iterable.com` in its published archive and is mapped.

## Modelling limits

- **One route.** Lists. Users, campaigns, templates, catalogues, events, workflows and the whole export surface each want their own evidence.
- **The IP values in the fixture are placeholders.** `203.0.113.10` is the RFC 5737 documentation range; the internal addresses are the two that were actually observed. A fake echoing the caller's real address would be inventing behaviour the emulator has no business having.
- **No `spec:`.** The reference at `api.iterable.com/api/docs` is a rendered Swagger UI page whose served HTML carries no document.
