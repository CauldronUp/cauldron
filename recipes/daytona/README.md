# daytona

Emulates the Daytona sandbox listing for local development and tests.

**9 conformance cases, 1 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Daytona serves without a credential at [`api.daytona.io/openapi.json`](https://api.daytona.io/openapi.json), and struck live on 2026-09-13 against that host with no credential, with a wrong token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The published document's only server is localhost over cleartext.** `servers: [{"url": "http://localhost:3000"}]`, on a document served from `api.daytona.io`. Its OpenID Connect discovery URL is `http://localhost:3000/.well-known/openid-configuration` to match. Every client generated from this document is pointed at the developer's own machine, over plain HTTP.

**And the host that serves it answers nothing else.**

```
GET /api/sandbox         301 -> https://www.daytona.io/api/sandbox
GET /api/health          301 -> https://www.daytona.io/api/health
GET /api/cauldron-nope   301 -> https://www.daytona.io/api/cauldron-nope
PUT /api/sandbox         301 -> https://www.daytona.io/api/sandbox
```

Following any of them lands on a Framer marketing page, 200, with `<!-- Made in Framer · framer.com ✨ -->` in the first line of HTML. The document is the one thing `api.daytona.io` will answer with, and everything else on it is a permanent redirect to a website. `cauldron drift` puts it the other way round: it reports `GET /api/sandbox` as undeclared, because the document describes `/sandbox` under a localhost server and the live host serves `/api/`.

**A sandbox listing is required to carry each sandbox's environment variables.** `env` — "Environment variables for the sandbox" — is in the `required` array of the record the listing returns, so reading the list reads every sandbox's environment.

**Two path parameters are authentication tokens.** `/organizations/otel-config/by-sandbox-auth-token/{authToken}` and `/organizations/sandbox-identity/by-sandbox-auth-token/{authToken}` put the secret in the URL, where every access log and every `Referer` keeps it.

**Network policy is a comma-separated string.** `networkAllowList` is "Comma-separated list of allowed CIDR network addresses" and `domainAllowList` is the same for domains — two lists, encoded as text, on a JSON record, in the fields that decide what a sandbox may reach.

**Four timers and at least two conventions for "off".** `autoStopInterval` and `autoPauseInterval` say "0 means disabled"; `autoDeleteInterval` uses a negative value; `autoArchiveInterval` says nothing about being disabled at all.

**A field's description names a resource that is not this one.** `user` is "The user associated with the project", on a schema called `Sandbox`.

**And there are two listings, one of them named after paging.** `/sandbox` and `/sandbox/paginated`, beside `/sandbox/for-runner`, in the same public document.

**`state` and `desiredState` sit on one record**, so a sandbox says what it is and what it is meant to be; `spot` and `spotEvictedAt` carry the moment it was preempted; and `recoverable` is a boolean about whether the error beside it can be undone.

## Sources

- [`api.daytona.io/openapi.json`](https://api.daytona.io/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.daytona.io`, struck 2026-09-13 on four paths, every one of them a 301 to `www.daytona.io`.

## Modelling limits

- **No credential is modelled.** `api.daytona.io` answers no API path at all — a 301 to a marketing site is the whole of it — so the shape of a refusal here was never observed, and this Recipe declines to invent one.
- **One route of a hundred and thirty-five**, plus the redirect. Snapshots, volumes, toolbox, runners, organisations, API keys, audit and the Docker registry surface are the rest.
- **The `env` values in the fixture are illustrative.** The schema requires the field; what a real sandbox puts in it is whatever its owner set, and this Recipe serves a shape rather than anyone's secrets.
- **Nothing is mapped in detection.** Daytona is reached through `@daytonaio/sdk` on npm or `daytona` on PyPI, and neither resolves to this host through a dependency file. Checked 2026-09-13.
