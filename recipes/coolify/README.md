# coolify

Emulates the Coolify self-hosted-platform API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Written against Coolify's own OpenAPI 3.1 document, published in its own repository — `coollabsio/coolify`, `openapi.json`, 192 paths, `info.version` "0.1" — read on 2026-09-06. Coolify Cloud exists but every endpoint needs an account and there is no anonymous surface, so every claim here is the provider's own document.

## What this Recipe found

**The document uses `writeOnly` nowhere, and declares secrets on responses.**

`writeOnly` is OpenAPI's way of saying "this may be sent and is never returned". It appears **zero times** in this 192-path, one-megabyte document.

Meanwhile the `Application` schema — the one a listing answers with — declares:

| field | the document's description |
|---|---|
| `http_basic_auth_password` | "Password for HTTP Basic Authentication" |
| `manual_webhook_secret_github` | "Manual webhook secret for GitHub." |
| `manual_webhook_secret_gitlab` | — |
| `manual_webhook_secret_bitbucket` | — |
| `manual_webhook_secret_gitea` | — |

and the `PrivateKey` schema declares `private_key`, `format: private-key` — the SSH key itself.

So as published, nothing in this document distinguishes a secret you *send* from a secret you *get back*. A generated client types all of them onto the response; a generated TypeScript interface for `Application` has a password on it. Any code that logs a response body, serialises one into an error report, or caches one is handling five secrets per application and does not know it.

**This is a claim about the document, not about the running server.** Coolify may well redact these at runtime; there is no anonymous endpoint to check against, so this Recipe does not claim it either way. What is checkable is what the document says, and the document says they are returned. The fixture uses obvious placeholders, never plausible values.

**There is no `POST /applications`.** Creating one is five separate endpoints, chosen by where the code comes from:

```
POST /applications/public
POST /applications/private-github-app
POST /applications/private-deploy-key
POST /applications/dockerfile
POST /applications/dockerimage
```

`/applications` itself has a GET and nothing else. All five answer the same `Application` shape, so a client has one type for the result and five paths for the request.

**A bad credential is 400, and no credential is 401.** The document's own reusable responses: `400 "Invalid token."`, `401 "Unauthenticated."` — the opposite way round from most providers here. A client whose refresh-and-retry path fires on 401 never fires on an expired token, because an expired token is a 400, and 400 is the status every client treats as "my request was malformed, do not retry".

**Two identifiers, and the one called `uuid` is not one.** `id` is "the application identifier in the database", an integer, exposed and addressing nothing. `uuid` is what every path takes — and it is a 24-character lowercase alphanumeric string like `k4o8sw0kow4c8c0soo4k0go8`, with no hyphens and no version nibble. A client validating it against a UUID pattern because of what it is called rejects every real identifier, and a database column typed `uuid` cannot hold one.

**Ports are a string.** `ports_exposes` is `type: string` with no format and no pattern, holding `"3000,9229"`. Reading a port means splitting; writing one means joining; and the separator is a convention the schema never states.

**The status is a compound string.** `"running:healthy"` — two facts joined by a colon inside one untyped field, with no enum for either half, so a client comparing to `"running"` never matches. (`build_pack` next to it *is* a proper enum.)

**Eighty-six fields on one record**, returned in full by a listing that takes one parameter — `tag` — and no limit, no page, no cursor and no total.

**The 429 declares `Retry-After`,** which most providers modelled here do not.

## Modelling limits

- **Nothing here is verified against a live API.** Coolify Cloud requires an account and has no anonymous surface.
- **Applications, listed, fetched and created.** 192 paths is a platform: services, databases, servers, projects, teams, deployments, private keys and the whole scheduled-task surface each want their own evidence.
- **The secret-shaped fields are served with obvious placeholders.** What is reproduced is that the keys are present on a response — the finding — not a credential worth having.
- **Twelve of the eighty-six fields are modelled.** The rest are real and omitted, so this Recipe under-states how much a listing carries.
