# hatchet

Emulates the Hatchet API for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-11.**

Read from the multi-file OpenAPI contract in Hatchet's own repository — it is open source — and struck live against `cloud.onhatchet.run` on 2026-09-11 with no credential, with a deliberately invalid one, on a path that does not exist, and on tenant paths carrying a malformed and an unknown identifier.

## What this Recipe found

**Nothing is ever a 401.**

```
(no header)      403 {"message":"Please provide valid credentials"}
Bearer notreal   403 {"message":"Please provide valid credentials"}
```

Same status, same sentence. 403 means the server knows who you are and you may not; 401 means tell me who you are. Hatchet only ever says the first — to callers who have said nothing.

Because 403 carries no challenge, there is no `WWW-Authenticate` anywhere to name the scheme. A client cannot tell a missing credential from a rejected one, and a retry-on-401 interceptor never fires.

**A path that does not exist answers 403 and the body `null`.** Four bytes, under `application/json`. It parses — `JSON.parse` returns null rather than throwing — so a client reaches `body.message` on it and gets a TypeError where an error message should be.

**And the document declares none of this.** `GET /api/v1/users/current` lists 400, 401 and 405, no 403 at all, and points all three at `APIErrors`:

```json
{"errors": [{"code": 1400, "field": "name", "description": "…",
             "docs_link": "github.com/hatchet-dev/hatchet"}]}
```

The live failure is `{"message": "…"}`. Not one field in common with the declared type, so a generated client's error struct parses nothing this API actually sends. The declared shape's own `docs_link` example is a URL with no scheme.

**The operation's declared security is a browser cookie.** The document's top-level `security` offers `bearerAuth` and `cookieAuth`; this operation overrides it to `cookieAuth` alone — an apiKey in a cookie named `hatchet`. So the endpoint answering "who am I" is documented as reachable only from a browser session, while the sentence it refuses a Bearer token with says to provide valid credentials.

**The unauthenticated caller gets a four-way oracle.** With no credential at all:

| request | answer |
| --- | --- |
| `/api/v1/users/current` | `403 {"message":"Please provide valid credentials"}` |
| `/api/v1/cauldron-nope` | `403 null` |
| `/api/v1/tenants/notauuid` | `400 {"message":"invalid tenant id"}` |
| `/api/v1/tenants/<a well-formed uuid>` | `404 {"message":"not found"}` |

Identifier parsing and existence lookup both run before authentication, so the response separates a real tenant from an absent one for a caller who has proved nothing. The identifiers are version-4 UUIDs, so this is not a practical way to find one — the finding is the order the checks run in, which is what generalises to anything else keyed the same way.

**`/api/v1/meta` answers 200 to anybody** and describes the deployment: which sign-in schemes exist, whether signup is open, whether tenants may be created, whether observability is on, and whether authentication is disabled entirely.

The schema for it also declares `authDisabledToken` — "the embedded worker API token, only set on authdisabled builds", example `eyJhbGciOiJFUzI1NiIs…`. So the same unauthenticated endpoint that reports authentication is off is, by design on those builds, the one that hands out the token. Hatchet Cloud reports `authDisabled: false` and omits the field.

**An id may be the empty string.** `APIResourceMeta.id` is described as "the id of this resource, in UUID format" and constrained `minLength: 0`, `maxLength: 36` — so a conforming record can identify itself with nothing, and the bound that would actually pin a UUID is 36 at both ends.

**And a user carries a hash of their own email.** `emailHash`, "for use with Pylon Support Chat" — the API's own user record includes a derived identifier whose only consumer is a third-party support widget.

## Sources

- Live: `cloud.onhatchet.run`, struck 2026-09-11.
- [`components/schemas/user.yaml`](https://github.com/hatchet-dev/hatchet/blob/main/api-contracts/openapi/components/schemas/user.yaml) — the `User` message.
- [`components/schemas/metadata.yaml`](https://github.com/hatchet-dev/hatchet/blob/main/api-contracts/openapi/components/schemas/metadata.yaml) — `APIResourceMeta`, `APIMeta`, `APIError`.
- [`paths/user/user.yaml`](https://github.com/hatchet-dev/hatchet/blob/main/api-contracts/openapi/paths/user/user.yaml) — the operation, and its cookie-only security.

## Modelling limits

- **Two routes.** The current user and the public metadata. Tenants, workers, workflows, task runs, events, rate limits, webhooks and the whole V1 stable surface each want their own evidence — the contract's root file alone lists well over a hundred paths.
- **The unrouted 403 sends no body; live it sends `null`.** A Recipe declares a body shape, and the bare JSON null literal is not one of them. The difference matters in one direction: `JSON.parse("null")` succeeds and yields null, while `JSON.parse("")` throws. So code written against this sandbox fails at the parse and code running against Hatchet fails one line later, on the property read.
- **The tenant oracle is recorded, not served.** The 400-then-404 ordering was struck live and is quoted above; modelling it would need a tenant route marked answerable without a credential, which is a claim the probes do not support — an unknown tenant 404s anonymously, and what a *real* one does anonymously was never tested.
- **No `spec:`.** Hatchet's OpenAPI contract is a tree of files stitched by relative `$ref`, so the document at the root address is a manifest rather than a description and there is nothing at one URL to fingerprint. The files are cited individually above.
- **Nothing is mapped in detection.** Hatchet's SDKs — `@hatchet-dev/typescript-sdk`, `hatchet-sdk` on PyPI, the Go package — are worker libraries that connect over gRPC to whichever engine a deployment configures, not clients of this REST API. A project holding one is running tasks, not reading the management surface. Checked 2026-09-11.
