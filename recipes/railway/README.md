# railway

Emulates the Railway public GraphQL API for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-07.**

Written against Railway's public API reference at `docs.railway.com` and struck live against `backboard.railway.app` on 2026-09-07 with no credential, with a deliberately invalid one, with a query naming a field that does not exist, and on a path that does not exist.

## What this Recipe found

**A refusal is HTTP 200 and a typo is HTTP 400.**

```
POST /graphql/v2   { projects { edges { node { id name } } } }      (no credential)
200 OK
{"errors":[{"message":"Not Authorized","locations":[{"line":1,"column":3}],
            "path":["projects"],"extensions":{"code":"INTERNAL_SERVER_ERROR"},
            "traceId":"909477571025499935"}],"data":null}

POST /graphql/v2   { nope }
400 Bad Request
{"errors":[{"message":"Cannot query field \"nope\" on type \"Query\".",
            "extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}, …}]}
```

The request that can never succeed until the caller finds a token comes back a success. The request with a typo in it comes back a failure. GraphQL-over-HTTP permits the 200 — a transport success carrying a field-level error is the specification's own model — but Railway does not apply it consistently, and the half it does not apply it to is the half a client's retry logic reads first.

**`INTERNAL_SERVER_ERROR` is the code for "Not Authorized".** That is Apollo Server's default when a resolver throws something carrying no code of its own, so this refusal is an unhandled exception wearing a response's clothes.

The consequence is not cosmetic. **Every retry policy ever written treats `INTERNAL_SERVER_ERROR` as transient and a missing credential as permanent**, so a client that branches on `extensions.code` will retry this forever, at 200, with backoff, and never learn why. The message beside it — "Not Authorized" — is correct, and it is the half no machine reads.

The validation failure four bytes away gets `GRAPHQL_VALIDATION_FAILED`, which is exact. So this API has a precise code for the mistake a developer makes once and a wrong one for the mistake every integration makes on its first request.

**`data` is null on the 200 and absent from the 400.** `"data" in body` is true for the failure that arrived as a success and false for the one that arrived as a failure, which is the opposite of how a client would guess.

**One service, two identifier names in two formats.** The GraphQL errors carry `traceId`, an eighteen-or-nineteen-digit number sent as a string — a snowflake, which a client that parses it to a number loses the tail of. The REST 404 on an unrouted path carries `requestId`, a UUID:

```
POST /graphql/cauldron-nope
404 {"message":"Not Found","requestId":"7b46b820-314e-4f7e-850b-f9bce2c9ff0f"}
```

Neither response mentions the other's field, so quoting an identifier to support means knowing which half of the service answered you.

**Introspection is open to anonymous callers.** `__type(name: "Query")` with no credential returns the whole schema. Among the field names it hands back: `adminVolumeInstancesForVolume`, `apiTokens`, `auditLogs`, `botScopeBindings`, `bucketS3Credentials`, `complianceAgreements`. So an unauthenticated request cannot read one field and can read the name of every field there is, including the administrative ones.

**The cursor is on the edge and not on the project.** Which is the whole reason a Relay edge exists as a separate object, and the reason a client cannot page from the records alone.

**A deleted project stays in the listing.** `deletedAt` is null until a trial project is reaped and a timestamp afterwards, and the record does not move — so a client counting projects counts the reaped ones unless its query asks for a field it has no reason to think of.

## Modelling limits

- **One route.** Projects. Services, deployments, environments, variables, volumes and the whole log surface each want their own evidence.
- **No GraphQL is parsed.** A document naming `projects` gets the projects fixture whatever fields it asked for, and paging is read from `variables.first` and `variables.after` rather than from arguments inlined in the document.
- **The emulator pages on the record's own identifier**, where Railway's `after` takes the edge's opaque base64 cursor. The cursor is served on every edge; positioning on it is not modelled.
- **`data` is null on every declared error here.** Railway omits the key entirely on the 400 and this format has one error envelope per Recipe, so the difference is recorded rather than served.
- **Nothing is mapped in detection.** Railway is driven by its CLI and by generic GraphQL clients, neither of which a dependency name distinguishes.
- **No `spec:`.** The API is GraphQL, so its description is a schema — and that schema is served to anybody who asks.
