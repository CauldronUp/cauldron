# grafanacloud

Emulates the Grafana Cloud API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Grafana's reference at `grafana.com/docs` and struck live against `grafana.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**"Login Required" is a browser sentence, on an API.**

```
GET /api/orgs    (no credential)
401 {"code":"Unauthorized","message":"Login Required","requestId":"01efe1da-…"}

GET /api/orgs    Authorization: Bearer notreal
401 {"code":"InvalidCredentials","message":"Token could not be parsed","requestId":"28dcdfe5-…"}
```

Nothing on this surface has a login. A caller holds a service-account token or a cloud access policy token; there is no session and no form. "Login Required" is what the same handler says to a browser hitting `grafana.com` without a cookie, reaching an API client through shared middleware — **so the sentence is accurate about a completely different request**.

**The two failures do carry different codes, and that is the good half.** `Unauthorized` for nothing sent and `InvalidCredentials` for something wrong, with two different sentences. Most providers in this collection collapse these into one response; Grafana separates them in both the machine-readable field and the prose, and then puts the wrong prose in one of them.

**"Token could not be parsed" is precise about the right thing.** It says the string was unreadable rather than unrecognised — so a caller who pasted a truncated token knows to look at the string, and one whose token was revoked knows this is not their message.

**An unknown path answers the credential failure**, so routing is judged after authentication.

**`requestId` is a UUID and it is the only correlation id.** One field, in the body, on every failure, with no header twin — which after [gong](../gong) and [imgix](../imgix) in this collection is worth noticing as the simpler and better choice.

**A record carries two handles for one organisation.** The `slug` is what every URL uses and the numeric `id` is what the API returns, so a person reading a dashboard link and a client reading a record are holding different keys for the same thing.

## Modelling limits

- **One route.** Organisations. Stacks, instances, access policies, tokens and the Prometheus-shaped metrics surface each want their own evidence — and the last of those is a second API with its own conventions.
- **Nothing is mapped in detection.** Grafana's clients are Grafana clients rather than Grafana *Cloud* clients: they talk to a stack's own HTTP API, not to `grafana.com`.
- **No `spec:`.** Grafana publishes a rendered reference site for the Cloud API.
