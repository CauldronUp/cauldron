# flyio

Emulates the Fly.io Machines API app listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-14.**

Read from the OpenAPI document Fly serves without a credential at [`docs.machines.dev/openapi.json`](https://docs.machines.dev/openapi.json), and struck live on 2026-09-14 with no token, with a wrong token, without the required query parameter, on a path that does not exist, with a method the path does not take, and on the root.

## What this Recipe found

**The challenge asks for a scheme the API does not take.** Every refusal carries:

```
www-authenticate: Basic realm="api.machines.dev"
```

On an API authenticated with `Authorization: Bearer <fly token>`. A browser opening one of these gets a username-and-password box. A client implementing the challenge implements the wrong scheme. And `realm` is the hostname — the one piece of information the client already had.

**The document declares no authentication at all.** `components.securitySchemes` is `{}`, and no operation carries a `security` block, in a 68-path document for an API that refuses every request without a token. A client generated from it is generated without a way to send one.

**And the listing declares no failures.** `GET /v1/apps` has exactly one response in the document: `200`. Every refusal below was struck live and none of them is described, so `cauldron drift` reports each declared error as not backed:

```
not backed: no operation the Recipe routes to answers 401, which token_validation_error declares
not backed: no operation the Recipe routes to answers 404, which unknown_route declares
```

Correctly — there is nothing in the document to back them with.

**A missing required query parameter is reported as a missing page.** `org_slug` is `required: true` on the listing.

```
GET /v1/apps?org_slug=personal   (no token)      401  application/json
GET /v1/apps                     (wrong token)   404  text/plain  404 page not found
GET /v1/cauldron-nope                            404  text/plain  404 page not found
PUT /v1/apps                                     404  text/plain  404 page not found
```

Forgetting the parameter, mistyping the path and mistyping the verb are one answer — Go's default not-found handler, in plain text, from an API that is JSON everywhere else — and all three are answered before the credential is looked at.

**A missing token and a wrong one are the same sentence, and the sentence is a Go error chain.** `{"error":"Authenticate: token validation error"}` — the wrapping verb still on the front, and "token validation error" covering both a token that failed validation and a token that was never sent.

**A record carries three identifiers, and the one called `id` is not the one you address it by.** The path is `/v1/apps/{app_name}`, so an app is addressed by `name`. The record also carries `id`, a string, and `internal_numeric_id`, an integer — a field whose name says it is internal, in a public response.

Also pinned: the app record has no `required` array and not one field description in the whole schema, so nothing in the document says which of the nine fields is guaranteed or what any of them mean; `total_apps` rides beside `apps` as the only paging-shaped field, with no page, limit or cursor anywhere in the operation; and the API root answers `301` with `content-type: application/json` and no body.

## Sources

- [`docs.machines.dev/openapi.json`](https://docs.machines.dev/openapi.json) — OpenAPI 3.0.1, `Machines API 1.0`, 68 paths, 178 schemas, served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-14`. `docs.machines.dev/spec/openapi3.json`, the path Fly's own docs link, 301s to it.
- [Fly Machines API documentation](https://fly.io/docs/machines/api/).
- Live: `api.machines.dev`, struck 2026-09-14 with no token, a wrong token, a missing required parameter, an unrouted path, a wrong method, and the root.

## Modelling limits

- **The missing-parameter 404 is described, not served.** Live, omitting `org_slug` answers Go's plain-text 404 before the credential is read. This format has no way to make a query parameter required on a listing, so the Recipe serves the listing whether or not it is sent, and the three live probes above are the record of what really happens.
- **The not-found body for an app is not observed.** Reaching it needs a real token; the 404 served here is the provider's own envelope with the sentence its CLI surfaces, marked documentation-only.
- **The success side is document-derived.** Listing apps needs a real token; the records here are the listing's own schema with values of the declared types, and every case reading them is marked documentation-only.
- **Two routes of sixty-eight.** The app listing and one app. Machines, volumes, certificates, IP assignments, leases, metadata, exec, secrets and the rest of the Machines surface are the others.
- **Nothing is mapped in detection.** Fly is reached through `flyctl` or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
