# dokploy

Emulates the Dokploy project listing for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the OpenAPI document Dokploy serves at [`docs.dokploy.com/openapi.json`](https://docs.dokploy.com/openapi.json) and from Dokploy's own source, and struck live on 2026-09-14 against `app.dokploy.com` with no key, with a wrong key, on a path that does not exist, with a method the path does not take, and without a required parameter.

## What this Recipe found

**Every success in the document is declared to be empty.** All 604 operations — every one of them — give their `200` as:

```json
{"type": "object", "properties": {}, "additionalProperties": false}
```

An object with no properties and no others permitted. The document says every endpoint in this API returns exactly `{}`. `/project.all` says it. So do `/application.one`, `/server.all`, `/deployment.allByType` and the other six hundred. The implementation returns rows — `db.query.projects.findMany` with nested relations — so the schema does not merely omit the shape, it forbids it.

**The only reusable types in a 1.8 MB document are its errors.** `components.schemas` holds five entries: `error.BAD_REQUEST`, `error.UNAUTHORIZED`, `error.FORBIDDEN`, `error.NOT_FOUND`, `error.INTERNAL_SERVER_ERROR`. Each is fully specified — `message`, `code`, an `issues` array, examples, a title. Every failure is typed and no success is.

**And the live failure does not match the error schema either.** `error.UNAUTHORIZED` requires `message` **and** `code`, with the example message "Authorization not provided" and code `"UNAUTHORIZED"`. Struck live:

```json
{"message":"Unauthorized"}
```

No `code`, on a schema that requires it, and a different sentence from the one the document gives as its example.

**One answer for every mistake.** No key, a wrong key, a path that does not exist, a verb a path does not take, and a required parameter left out all answer that same 401. Nothing in the response distinguishes them.

**The document's only server is a hostname that does not exist.** `servers: [{"url": "https://your-dokploy-instance.com/api"}]`. Checked 2026-09-14: `your-dokploy-instance.com` does not resolve — NXDOMAIN. A client generated from this document points at an unregistered domain, and whoever registers it receives the requests, `x-api-key` header included.

**One authentication scheme is declared twice.** `securitySchemes` has `apiKey` and `x-api-key`, both `{"type": "apiKey", "in": "header", "name": "x-api-key"}`, differing only in their description. The global `security` is `[{"apiKey": []}, {"x-api-key": []}]` — which reads as two alternatives, and is one requirement written twice.

**A project's environment is encrypted at rest and fails open on read.** Dokploy's `encryptedText` column type decrypts in `fromDriver` and, on failure, returns the raw value with a log line:

> "Fail open so a key mismatch (e.g. restoring a backup under a different BETTER_AUTH_SECRET) degrades to showing ciphertext instead of breaking every query that touches the row."

So `env` can arrive as AES-256-GCM ciphertext in the same string field that usually holds the variables, with nothing on the record saying which it is.

Also pinned: paths are RPC names with dots — `/project.all`, `/project.one`, `/application.deployNginxQuickstart` — and the verbs are not consistent: `network.all` and `certificates.all` beside `gitProvider.getAll` and `ai.getAll`, with `deployment.allByCompose`, `deployment.allByServer`, `deployment.allCentralized` and `deployment.allByType` on one resource; `createdAt` is a `text` column holding an ISO string rather than a timestamp; and the listing nests one key per service type — `applications`, `libsql`, `mariadb`, `mongo`, `mysql`, `postgres` and the rest — so the response grows a top-level key every time Dokploy supports another engine.

## Sources

- [`docs.dokploy.com/openapi.json`](https://docs.dokploy.com/openapi.json) — OpenAPI 3.1.0, `Dokploy API 1.0.0`, 604 paths, 5 schemas; recorded by `cauldron drift`, `spec_seen: 2026-09-14`.
- [`Dokploy/dokploy`](https://github.com/Dokploy/dokploy) — `packages/server/src/db/schema/project.ts` for the project columns, `.../schema/utils.ts` for `encryptedText`, and `apps/dokploy/server/api/routers/project.ts` for what `project.all` actually returns.
- Live: `app.dokploy.com`, struck 2026-09-14 with no key, a wrong key, an unrouted path, a wrong method, and a missing required parameter.

## Modelling limits

- **The record comes from the source, not the description.** The document declares nothing about any success, so the fields here are Dokploy's own Drizzle columns and the resolver's `with:` clause. Every case reading them is marked documentation-only.
- **The not-found body is not observed.** `error.NOT_FOUND` is declared in the document; reaching it needs a real key, so the 404 served here is that schema's shape with a plausible sentence, marked documentation-only.
- **Two routes of six hundred and four.** `/project.all` and `/project.one`. Applications, compose stacks, databases of six kinds, servers, deployments, domains, certificates, registries, SSH keys, backups, notifications, git providers and the rest are the others.
- **The nested service keys are abridged.** The real listing nests six database types plus applications and compose; the fixture carries applications, postgres, and four empty arrays, which is enough to show the shape.
- **Nothing is mapped in detection.** Dokploy is reached through a plain HTTP call carrying an `x-api-key` header against a self-hosted instance, which does not resolve to any host through a dependency file. Checked 2026-09-14.
