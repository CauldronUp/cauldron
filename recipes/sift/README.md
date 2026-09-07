# sift

Emulates the Sift risk-scoring API for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Sift's reference at `developers.sift.com` and struck live against `api.sift.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Everything is HTTP 400 and the real status is a small integer in the body.**

```
GET  /v205/events           400 {"status": 50, "error_message": "Invalid request: HTTP method GET is not supported by this URL", "time": 1788814232}
POST /v205/events    {}     400 {"status":51,"error_message":"Invalid API Key. …","time":1788814232,"request":"{}"}
GET  /v205/users/abc/score  400 {"status": 51, "error_message": "Invalid API Key. …", "time": 1788814232}
GET  /v205/cauldron-nope    400 {"status": 50, "error_message": "Invalid request: Not Found", "time": 1788814233}
```

A wrong verb, a missing credential, a wrong credential and an unknown path all answer **400**. The HTTP status carries no information at all; `status` in the body does, and it is `50` or `51` — Sift's own numbering, unrelated to HTTP and unrelated to anything else in this collection.

So every branch a client writes has to read the body first, and a proxy, a retry policy or an alerting rule keyed on the status line sees one value for four different problems. `50` covers both a verb the URL does not take and a URL that does not exist; `51` is the credential. Those two numbers are the whole vocabulary a public probe can reach.

**`time` is the server's clock on every response** — seconds since the epoch, as a number, on failures as well as successes. Which is genuinely useful and unique here: a client can measure its own clock skew from a rejection, which matters for an API whose events are timestamped by the caller.

**The POST echoed the request body back.** `"request":"{}"` — the body exactly as sent, as a *string*, inside the failure. An empty object was sent on purpose. Sift's events API takes the API key as a **body field** (`$api_key`), so a request that reaches this handler with a real key has that key echoed into a response the caller will very likely log. This Recipe models the score API instead, where the credential is a query parameter and is at least not repeated back.

**Two serialisers.** The GET responses are spaced — `{"status": 50, …}` — and the POST response is compact. Same API, same version, same failure vocabulary, two JSON writers.

**A user has one score per abuse type.** Clean for account takeover and not for payments, on one record — so there is no single number, and a dashboard showing "the" score has picked a side.

**The identifier is whatever the customer calls the user.** Sift mints none of its own, so the path segment is an unvalidated, unbounded customer string that has to be URL-safe by luck.

## Modelling limits

- **One route.** The score API. Events, labels, decisions, workflows and the whole webhook surface each want their own evidence — and the events endpoint is a POST whose credential is in the body, which this Recipe deliberately does not model for the reason above.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Sift publishes a rendered reference site.
