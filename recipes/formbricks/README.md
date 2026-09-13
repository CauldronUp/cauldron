# formbricks

Emulates the Formbricks management API for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from Formbricks's own source — it is open — and struck live against `app.formbricks.com` on 2026-09-13 with no credential, with a deliberately invalid one, on both API versions, and on a path that does not exist.

## What this Recipe found

**Two API versions on one host, two error shapes.**

```
/api/v1/management/me  401 {"code":"not_authenticated","message":"Not authenticated",
                            "details":{"x-Api-Key":"Header not provided or API Key invalid"}}
/api/v2/me             401 {"error":{"code":401,"message":"Unauthorized"}}
```

`body.code` is the string `not_authenticated` in one and the number `401` in the other. `body.message` is at the top level in one and one level down in the other.

A client that talks to both — the normal case while a migration is in progress — needs two readers for the same failure, and nothing but the path tells it which it is about to get.

**The details map is keyed by a header name spelled three ways.**

```json
"details": { "x-Api-Key": "Header not provided or API Key invalid" }
```

Lower-case `x`, capital `A`, capital `K`. HTTP header names are case-insensitive, so a caller who sent `x-api-key` — which is what the documentation tells them to send — and then looks that name up in this map finds nothing. The key is a label rather than an identifier.

**And the detail merges the two conditions it exists to separate.** "Header not provided **or** API Key invalid": the one field in the body whose job is to say which of the two happened says both. Struck live with no header and with a wrong key, the bodies are byte-identical.

**The v2 success envelope makes its pagination optional.** From the source:

```ts
interface ApiResponseWithMeta<T> extends ApiResponse<T> {
  meta?: { total?: number; limit?: number; offset?: number }
}
```

`meta` is optional and so is every field inside it. A paginated response may carry no pagination at all, and a client cannot tell "there are no more" from "the server did not say".

**The listing can be sorted by two fields and neither is a name.** `sortBy` is an enum of `createdAt` and `updatedAt`. A caller who wants surveys in alphabetical order, or grouped by status, sorts them itself.

**And the status enum mixes two casings.** `draft`, `inProgress`, `paused`, `completed` — three lower-case words and one camelCase, in one enum.

**An unrouted path under `/api` answers the web application.** A full Next.js HTML document, stylesheet links and all, from the path prefix reserved for the API.

## Sources

- Live: `app.formbricks.com`, struck 2026-09-13.
- [`api/v2/management/surveys/types/surveys.ts`](https://github.com/formbricks/formbricks/blob/main/apps/web/modules/api/v2/management/surveys/types/surveys.ts) — the filters, the enums and the identifier format.
- [`api/v2/types/api-success.ts`](https://github.com/formbricks/formbricks/blob/main/apps/web/modules/api/v2/types/api-success.ts) — the envelope, and its optional pagination.

## Modelling limits

- **One route, on v1.** The v1 management surveys listing, which is the one the live probes reached. Responses, contacts, attributes, actions, webhooks, display and the whole v2 surface each want their own evidence.
- **The v2 error shape is recorded, not served.** Both shapes are quoted above; a Recipe declares one error envelope, and this one serves v1's because that is the version the route belongs to.
- **The survey record is the part the API's own filters pin down.** A Formbricks survey carries questions, logic, styling, triggers and a dozen more structures defined in Zod; the fields here are the ones the v2 filter and sort enums name, which are the ones a listing can be queried by.
- **No `spec:`.** Formbricks defines its API in Zod schemas in the application itself and renders a reference from them at build time, so there is no document at an address to fingerprint.
- **Nothing is mapped in detection.** A project using Formbricks holds its embed script or one of the tracking SDKs, which post responses rather than read the management surface — and self-hosted instances put the host in configuration. Checked 2026-09-13.
