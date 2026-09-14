# replit

Emulates the Replit project listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Replit serves without a credential at [`api.replit.com/openapi.json`](https://api.replit.com/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong key, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**Two paths have swapped names.** `/audit-logs` carries the operation `listAuditLogsLegacy`. `/compliance/messages` carries `listAuditLogs`, and its 200 is described as "Audit log event list."

So the path named for audit logs is the legacy one, and the current audit-log endpoint is the one named for compliance messages. A client picking by path name picks the wrong one, and a client picking by operation name picks a path that does not say what it is.

**The credential is checked before the path is.** Struck live, all four of these answer the same body:

```
GET /v1/projects              401
GET /v1/compliance/messages   401
GET /v1/cauldron-nope         401
PUT /v1/projects              401

{"error":{"code":"unauthenticated","message":"No API key provided.","details":null}}
```

A typo in the path, a typo in the verb, and a forgotten key are one answer. Nothing in it distinguishes them, so a caller debugging a 401 cannot tell whether the key is the problem.

**`details` is always present and always null.** Every failure carries a third key with nothing in it — on both the missing-key and wrong-key refusals, and on the unrouted path and the wrong method too.

**Paging is a `oneOf` discriminated by a boolean constant.** `pagination` is one of:

```json
{"cursor": "<string, 1..2048>", "hasMore": true}
{"cursor": null,               "hasMore": false}
```

Both branches are `required: ["cursor", "hasMore"]` and `additionalProperties: false`, so the key set never changes. What changes is the type: `cursor` is `string` in one branch and `"type": "null"` in the other — typed null rather than nullable — and the branches are told apart by a `const` on a boolean rather than by a field a reader can look up.

**A field called creator is documented as the owner.** `creatorId` is "Current project owner user ID, or null when unavailable" — present tense, ownership, under a name about who made it. It is in the `required` array and it is `string | null`, so the contract requires a field that may be absent in value.

**And a project that has never been updated says it was.** `updatedAt` is "Last project update time, or createdAt when no update timestamp exists". An untouched project reports its creation time as its update time, and nothing on the record separates that from a project genuinely updated in the second it was made.

**The page cap lives on the response type.** `data` is `{"type": "array", "maxItems": 100}`, so the limit on a page is a property of the schema the listing returns rather than of the `limit` parameter that asks for it.

Also pinned: `bearerFormat` is `"Replit API key"`, a human phrase in the slot where a token format name belongs, with the `rpl_` prefix only in the prose beside it; every schema in the document is `additionalProperties: false`; and the workspace embedded in a project declares `^[A-Za-z0-9]+$` for its id and `^[0-9A-Za-z-]+$` for its slug, both anchored, with no pattern at all on the name between them.

## Sources

- [`api.replit.com/openapi.json`](https://api.replit.com/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- [Replit API documentation](https://docs.replit.com/replit-api).
- Live: `api.replit.com`, struck 2026-09-13 with no credential, a wrong key, an unrouted path, and a wrong method.

## Modelling limits

- **One route of twelve.** The project listing. Audit logs (under both names), budgets, deployments, groups, group users, members, usage and workspaces are the rest.
- **The success fixture is document-derived.** Listing projects needs a real `rpl_` key; the records here are `project`'s own required properties with values of the declared types and patterns, and every case reading them is marked documentation-only.
- **The rate-limit headers are not served.** The document declares `X-Request-Id` and a family of `X-RateLimit-*` headers on the 200; this Recipe serves the body and does not invent limit numbers.
- **Nothing is mapped in detection.** Replit is reached through a plain HTTP call carrying a bearer token beginning `rpl_`, which does not resolve to this host through a dependency file. Checked 2026-09-13.
