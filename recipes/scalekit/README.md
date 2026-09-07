# scalekit

Emulates the Scalekit API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Scalekit's reference at `docs.scalekit.com` and struck live against `api.scalekit.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The API refuses on the domain before anything else.**

```
GET /api/v1/organizations    (no credential)
GET /api/v1/organizations    Host: cauldron.scalekit.com
GET /api/v1/organizations    Authorization: Bearer notreal

404  Content-Type: text/plain
invalid domain in request
```

All three, identically. Not 401, not 403 — **404, in plain text**, saying nothing about the credential and everything about the hostname.

Scalekit issues every customer their own subdomain, and `api.scalekit.com` is not one of them: it is the apex the documentation names, and it serves nothing. **So the first request a developer makes, copied from the reference, is refused for a reason the reference does not mention** — and the refusal does not say what a valid domain would look like.

**A `Host` header does not help**, because the tenant is resolved from the TLS SNI and the certificate rather than from the header. Which is correct, and it means the only way to reach this API is to already know your own subdomain.

**It is 404 rather than 421.** RFC 9110 defines **421 Misdirected Request** for exactly this: a request sent to a server that cannot produce a response for the target URI's authority. Nothing in this collection sends a 421, and this is the one place it belonged.

**The whole observable API is four words.** No JSON, no code, no correlation id, no link. Lower case, no full stop.

**And the host root serves a React application.** `GET /` answers 200 with an HTML shell carrying `<meta name="sk-onprem">` — the self-hosted build of Scalekit's dashboard, on the API hostname. The apex serves a single-page app to a browser and four words to an API client.

**The region is per organisation**, so a listing can span jurisdictions and a client cannot assume one for all of them.

## Modelling limits

- **One route.** Organizations. Connections, directories, users, roles, sessions and the whole SSO and SCIM surface each want their own evidence — none of them reachable on the apex.
- **The subdomain is not modelled.** The emulator serves the routes; reproducing "this hostname is not yours" would mean modelling TLS.
- **Nothing is mapped in detection.** Scalekit's SDKs are named for the company, and a project depending on one is using SSO, directory sync or passwordless — three different surfaces.
- **No `spec:`.** Scalekit's API is gRPC-transcoded, so its surface is described as protobuf rather than OpenAPI.
