# waypoint

Emulates the HashiCorp Waypoint server API, for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Waypoint's own generated Swagger 2.0 document, published in its own repository — `hashicorp/waypoint`, `pkg/server/gen/server.swagger.json`, 148 paths — read on 2026-09-06. Waypoint is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**A path parameter called `application.application.application`.** One real path from this document:

```
/project/{application.application.project}/application/{application.application.application}/instances
```

The word "application" appears **five times** in one URL template: once as a literal segment and four times inside two parameter names. The second parameter is the same word three times, separated by dots.

These are protobuf field paths leaking through gRPC-gateway. **38** of this document's path parameters contain a dot, and the deepest are four segments long — `deployment.sequence.application.project`, `pipeline.owner.project.project`, and `targetRunner.id.id`, which is a field called `id` inside a field called `id`.

A generator turns `{application.application.application}` into a parameter of that name, which is not a legal identifier in any language it targets. [argocd](../argocd) has six such parameters from the same toolchain; this one has 38.

**The paging parameters are dotted too** — `pagination.page_size` and `pagination.next_page_token`, with dots, in a query string. Every HTTP library sends them verbatim; every struct generator has to rename them.

**A listing returns references, not records.**

```
ListProjectsResponse:
  projects: array of Ref.Project

Ref.Project:
  project: string
```

`projects` is an array of one-field objects wrapping a name. `projects[0].name` is `undefined`; `projects[0].project` is the name. And a client that wants anything about a project fetches each one, so listing N projects and reading them is N+1 requests.

**Paging tokens go empty rather than null.** From the document's own description of both: "The value will become empty when there are no more pages." So the end of a listing is `next_page_token === ""` — not null, not absent. Empty is falsy, so `while (token)` terminates by accident and `token !== null` never does.

**The count is a string.** `total_count` is `{"type": "string", "format": "uint64"}` — protobuf's JSON mapping, the same trade [etcd](../etcd) makes. `total_count === 0` is never true on an empty listing.

**Two paths per resource, one per identifier:**

```
/add-on-definition/by-id/{add_on_definition.id}
/add-on-definition/by-name/{add_on_definition.name}
```

Five `by-id` paths and three `by-name` ones. **This is the good half of the document.** Rather than one path accepting either — which is [harbor](../harbor)'s `X-Is-Resource-Name` header deciding what a segment means, and [gitea](../gitea)'s number-or-slug — Waypoint puts the choice in the URL, so a request says which key it used and a name made of digits is unambiguous.

**A project's applications are a cache the document admits may be stale:** "the set of applications that are known about this project. Note that this may not exactly represent the project configuration if a user hasn't run `waypoint init` yet." So an empty list means either "no applications" or "nobody has initialised this" — two facts with one representation, and the document says so rather than letting a reader assume the field is authoritative.

**And the placeholders again.** `title` is `pkg/server/proto/server.proto` — a filename — and `version` is `"version not set"`. That is the **third** document in this collection with that exact version string, after [argocd](../argocd) and [etcd](../etcd), all three from grpc-gateway. No `securityDefinitions` either, which makes five.

## Modelling limits

- **Nothing here is verified against a live API.** Waypoint is self-hosted and there is no public instance.
- **Projects, listed and fetched.** 148 paths is a deployment platform: applications, deployments, releases, builds, jobs, runners, pipelines, add-ons and the whole log surface each want their own evidence.
- **The dotted path parameters are recorded and not routed.** A dotted placeholder would collide with the dotted field paths this format uses everywhere else — the same limit [argocd](../argocd) records.
- **The credential is checked as a bearer token.** The document declares none and Waypoint uses one; serving no credential would teach a caller that anonymous requests work.
