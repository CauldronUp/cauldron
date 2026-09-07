# tinybird

Emulates the Tinybird API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Tinybird's reference at `docs.tinybird.co` and struck live against `api.tinybird.co` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The refusal echoes the token as a Python bytes literal.**

```
GET /v0/pipes    Authorization: Bearer notreal
403 {"error": "invalid authentication token. Invalid token b'notreal': ",
     "documentation": "https://docs.tinybird.co/api-reference/overview#authentication"}
```

`b'notreal'` is Python's `repr` of a bytes object — the prefix, the quotes and all. **The string a caller sent comes back wrapped in the syntax of the language that rejected it**, and a real token would come back the same way, in a body a client is likely to log.

The sentence also ends `": "` — a colon, a space, and nothing. Something was meant to follow: an exception's own message, empty because the exception carried none. Two failures concatenated into one sentence, and the second half is missing.

**Every failure carries a documentation link, and it is a deep one.** The credential failures point at the authentication anchor; the routing failure points at the API reference index. Two links, each aimed at the page for its own failure.

That is the **second-best** `documentation_url` in this collection. [saltedge](../saltedge) deep-links to the anchor for the *exact error class*, which is better still; Tinybird links to the section, which is more than GitHub, [docspring](../docspring) or [statsig](../statsig) manage.

**A credential failure is 403**, not 401, for a request carrying nothing at all — so the status that tells a client to authenticate never arrives.

**The JSON is pretty-printed with a space after each colon** — Python's `json.dumps` default. The same class of tell as [gong](../gong)'s space *before* the colon, which is Jackson's.

**A copy pipe has no endpoint to call.** `type: "copy"` writes to a datasource on a schedule and has no URL, so `endpoint` is `null` — and a client building a query URL from a pipe listing builds a broken one for every pipe of that kind.

**An endpoint pipe names itself as its endpoint.** `endpoint` holds the pipe's own id rather than a URL, so the field named for an endpoint is not one.

## Modelling limits

- **One route.** Pipes. Datasources, tokens, jobs, the query endpoint and the whole ingestion surface each want their own evidence.
- **Nothing is mapped in detection.** Tinybird is used through its CLI and through a query URL rather than a client library, and nothing on any registry calls `api.tinybird.co` under an obvious name.
- **No `spec:`.** Tinybird publishes a rendered reference site.
