# argocd

Emulates the Argo CD API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Argo CD's own generated Swagger 2.0 document, published in its own repository — `argoproj/argo-cd`, `assets/swagger.json`, 82 paths — read on 2026-09-06. Argo CD is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**The published document declares no authentication at all.** `securityDefinitions` is absent. So is a top-level `security`. So is any per-operation `security`. The machine-readable description of the most widely deployed GitOps controller — a system whose entire job is deciding who may change what runs in a cluster — says nothing about credentials.

Argo CD does authenticate: a JWT from `POST /api/v1/session`, sent as `Authorization: Bearer <token>`. None of that is in the document, so a client generated from it has no credential support and every request it makes is anonymous.

**And it ships with the generator's placeholders still in it:**

```json
"title":       "Consolidate Services"
"version":     "version not set"
"description": "Description of all APIs"
```

All three are defaults nobody replaced. A tool that names a generated package from `info.title` produces one called `consolidate-services`; a tool that pins a client to `info.version` pins it to the string `"version not set"`.

**Six path parameters have dots in their names:** `{application.metadata.name}`, `{project.metadata.name}`, `{creds.url}`, `{repo.repo}`, `{id.value}`, `{source.repoURL}`. These are protobuf field paths leaking through gRPC-gateway into the URL template. A generator turns `{application.metadata.name}` into a parameter named `application.metadata.name`, which is not a legal identifier in most languages — so it becomes `application_metadata_name`, or `applicationMetadataName`, or is dropped, depending on which generator you ran.

**And one resource is addressed by three different names:**

```
GET /api/v1/applications/{appName}/server-side-diff
GET /api/v1/applications/{applicationName}/resource-tree
PUT /api/v1/applications/{application.metadata.name}
```

Same resource, same API, three spellings of "which application". A generated client has three unrelated parameters for one value.

**`code` is a gRPC status code, not an HTTP one.** The only error shape is `runtimeError`: `{code, error, message, details}`. `code` is gRPC's — 5 is `NOT_FOUND`, 7 is `PERMISSION_DENIED`, 16 is `UNAUTHENTICATED` — travelling in a body delivered over HTTP, beside an HTTP status that is a different number for the same event. A client reading `body.code` as a status reads 5 where the response said 404.

**`error` and `message` are both strings, both for the reason**, and the document says nothing about which is populated when.

**Every failure is `default`.** Operations declare a 200 and a `default` and nothing else — no 401, no 403, no 404. The status a failure arrives with is not constrained anywhere in the document.

**`project` and `projects` are the same parameter.** Both live on the applications listing. The document's own description of `project`: "the project names to restrict returned list applications (legacy name for backwards-compatibility)". Two parameters, one meaning, distinguished by a parenthetical. Sending both is undefined.

**The listing is a Kubernetes list, cursor and all.** `{items, metadata}`, where `metadata.continue` is an opaque cursor whose own description warns that continuing "may not be possible if the server configuration has changed or more than a few minutes have passed" — an expiry measured in minutes, with no field saying when it expired. The parameter is spelled `continue`, a reserved word in every C-family language, so generators rename it without saying so.

**The empty sync status is a real value.** `Synced`, `OutOfSync`, and `""` meaning the comparison has not run. A dashboard treating falsy as "no data" and `OutOfSync` as "a problem" shows the third as neither. Health is a separate six-valued axis where `Missing` is not `Unknown`.

## Modelling limits

- **Nothing here is verified against a live API.** Argo CD is self-hosted and there is no public instance.
- **Applications, listed and fetched.** 82 paths is a GitOps controller: sessions, repositories, clusters, projects, accounts, certificates, GPG keys and the sync and rollback operations each want their own evidence.
- **The credential this Recipe checks is the bearer token Argo CD actually uses,** which the document does not mention. Serving no credential at all — what the document describes — would teach a caller that anonymous requests work, and they do not.
- **`{application.metadata.name}` is recorded and not routed.** A dotted placeholder would collide with the dotted field paths used everywhere else in a Recipe, so the route uses the plain name and the finding is in words.
- **`code` is served as the gRPC code the document types it as.** It is not made to agree with the HTTP status, because upstream it does not.
