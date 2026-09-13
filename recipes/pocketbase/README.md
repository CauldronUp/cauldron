# pocketbase

Emulates the PocketBase collections API for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from PocketBase's own source and documentation — it is open — and struck live against `pocketbase.io` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The success key is on every failure, empty.**

```
(no header)        401 {"data":{},"message":"The request requires valid record authorization token.","status":401}
Authorization: x   401 {"data":{},"message":"The request requires valid record authorization token.","status":401}
/api/cauldron-nope 404 {"data":{},"message":"File not found.","status":404}
```

`data` is where a successful write puts its record, and it is present and empty on every error. So `body.data` is truthy on both a success and a failure, and a client testing for it learns nothing.

**An unrouted API path is a missing file.** `File not found.` — because one process serves the API and the static web assets, so a path matching no route falls through to the file server. A caller who mistypes an endpoint is told about the filesystem.

**And the two credential states are one sentence.** A missing header and a wrong token answer byte-identical bodies.

**`null` is the most restrictive value a rule can have and `""` is the most permissive.** Each collection carries five access rules — `listRule`, `viewRule`, `createRule`, `updateRule`, `deleteRule` — typed as nullable strings, with these meanings from PocketBase's own documentation:

| value | who may perform the action |
| --- | --- |
| `null` | only an authorised superuser (the default) |
| `""` | anyone: superusers, authorised users **and guests** |
| `"expression"` | only callers satisfying the filter |

The difference between "nobody but an administrator" and "the entire internet" is `null` against `""` in a JSON string field — two values most languages treat as equally falsy, that most form encoders cannot tell apart, and that differ by two characters in a diff.

**A listing that fails its rule answers 200 and no records.** From the same documentation:

> the API will return 200 empty items response in case a request doesn't satisfy a listRule, 400 for unsatisfied createRule

Authorisation failure on a read is shaped exactly like an empty collection, so a client cannot tell "you may not see these" from "there are none".

**The rules are also filters.** Again from the documentation: "PocketBase API Rules act also as records filter". The same expression that decides whether a caller may list also decides which rows come back, so a partial result and a complete one are indistinguishable.

**One record shape covers three kinds of collection.** `Collection` embeds base, auth and view options together, so a base collection's record carries the fields an auth collection would use and a view collection's carries both.

**And a collection's indexes are raw SQL.** `"CREATE UNIQUE INDEX idx_users_email ON users (email)"` — a DDL statement, as a string, in a JSON array on the record.

## Sources

- Live: `pocketbase.io`, struck 2026-09-13.
- [`core/collection_model.go`](https://github.com/pocketbase/pocketbase/blob/master/core/collection_model.go) — the record, and the nullable rules.
- [API rules and filters](https://pocketbase.io/docs/api-rules-and-filters/) — what null and the empty string mean, and the 200-empty behaviour.

## Modelling limits

- **One route.** Listing collections. Records, realtime, files, logs, backups, settings and the auth surface each want their own evidence.
- **The 200-empty behaviour is recorded, not served.** A listing whose `listRule` a caller does not satisfy answers 200 with no items; Cauldron does not evaluate rules, so the rules here are data on the record rather than policy the sandbox enforces.
- **The live host is PocketBase's own demo.** `pocketbase.io` serves the project's site and a PocketBase instance from one process, which is why the 404 says "File not found." Any self-hosted instance is the same binary.
- **No `spec:`.** PocketBase documents this API as prose and ships its own JS and Dart clients. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** The `pocketbase` package on npm is a real client, but every instance is self-hosted and its address lives in configuration rather than in the dependency. Checked 2026-09-13.
