# valtown

Emulates the Val Town API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Val Town serves without a credential at [`api.val.town/openapi.json`](https://api.val.town/openapi.json), and struck live on 2026-09-13 with no credential, with a deliberately invalid one, on a path that does not exist, and on a val id that does not exist.

## What this Recipe found

**Two required parameters, both with defaults.** `GET /v2/me/vals` declares `offset` and `limit` as `required: true` and gives each a `default` — 0 and 20. A required parameter can never be absent, and a default only applies when one is. So neither default can ever take effect, and a client generated from this document is obliged to send both on every call.

**The endpoint that fetches a val declares no security at all.**

| operation | declared security |
| --- | --- |
| `GET /v2/me/vals` | `[{bearerAuth: []}]` |
| `GET /v2/vals/{val_id}` | `null` |

That is on a resource whose own `privacy` enum has `private` in it. Struck live, an anonymous request for an id that does not exist answers 404 rather than 401 — so the lookup really does run before any credential check.

**And the val it cannot find is called a Project.**

```
/v2/vals/00000000-0000-0000-0000-000000000000   404 {"error":"Project could not be found"}
```

The noun in the error is the old one. The document still carries a `Project` schema beside `Val`, so the rename reached the types and not the sentences.

**There are two error shapes, from two layers.**

```
/v2/me/vals         401 {"error":"Unauthorized"}
/v1/cauldron-nope   404 {"message":"Route GET:/v1/cauldron-nope not found","error":"Not Found","statusCode":404}
```

The application answers one key. Fastify underneath it answers three, and reuses `error` for a reason phrase where the application uses it for the whole sentence. So `body.error` is `Unauthorized` on one failure and `Not Found` on the other, with the actual explanation moved into a key the application never sends.

**Neither 401 is declared.** The listing operation lists exactly one response — 200 — while declaring a bearer credential, so the first thing a caller meets is absent from the description. `cauldron drift` says so directly:

```
valtown   not backed: no operation the Recipe routes to answers 401, which unauthorized declares
          not backed: no operation the Recipe routes to answers 404, which unknown_route declares
          not backed: no operation the Recipe routes to answers 404, which val_not_found declares
```

A missing credential and a wrong one answer the same word, so they cannot be told apart either.

**`username` is required and nullable.** Every val carries an `author`, every author must carry a `username`, and that username is `anyOf: [string, null]` — the field naming who wrote it is guaranteed present and permitted to say nothing. `imageUrl` and `description` are required and nullable too, so three of the record's eight required fields guarantee a key rather than a value.

**`author.type` is two single-value enums instead of one enum.** The schema is `anyOf: [{enum: ["user"]}, {enum: ["org"]}]` where `enum: ["user", "org"]` says the same thing — so a generator emits a union of two string literal types for a field with two values.

**And `unlisted` is documented as security by link.** From the privacy enum's own description: "Unlisted resources do not appear on profile pages or elsewhere, but you can link to them."

**The document's version is `1` and it describes two.** `info.version` is the string `1`, over paths under both `/v1` and `/v2` — including two different shapes of the same val.

## Sources

- Live: `api.val.town`, struck 2026-09-13.
- [`openapi.json`](https://api.val.town/openapi.json) — served without a credential.

## What it added to Cauldron

Val Town's records carry a `links` object and so do its pages, and the drift gap report used one matching name as proof it had opened the envelope. That collision made it report `name`, `privacy` and `imageUrl` as fields the description does not declare — three sentences about a description that declares all three. The report now looks through a listing's envelope to the record before deciding what counts as missing.

## Modelling limits

- **Two routes.** Listing your own vals, and fetching one. Blobs, SQLite, branches, files, environment variables, orgs, memberships, email and the telemetry endpoints are 32 more paths.
- **`/v1` is not modelled.** The document describes two API versions with two shapes of a val in them; this Recipe serves the `/v2` one, and the `/v1` path appears only as the unrouted-404 case.
- **The unrouted 404 carries the path this Recipe struck.** Live it names the method and path that were asked for.
- **Nothing is mapped in detection.** Val Town is reached through its own CLI, its web editor and generic HTTP clients holding a token — no package name on npm or elsewhere says a project calls this API. Checked 2026-09-13.
