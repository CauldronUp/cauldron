# infisical

Emulates the Infisical secrets-management API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-06.**

Written against Infisical's own OpenAPI 3.0.3 document, served by the API itself at `us.infisical.com/api/docs/json` — 1479 paths, 2317 operations, **12,974,420 bytes** — and struck live against `us.infisical.com` on 2026-09-06. The five live cases are the failure path: made with no credential, then with two deliberately invalid ones. The five documentation-derived cases are the listing, which needs an account.

## What this Recipe found

**The document contains zero `$ref`s.**

Not few. None. `components.schemas` is empty and the string `"$ref"` does not occur once in 12.9MB. Every schema is written out longhand at every use site.

```
response schemas declared   16,097
distinct shapes among them     869
copies of the most common    2,298
```

**94.6% of the schemas in the file are a copy of another schema in the file.** The most repeated one is the 400 error envelope, written out 2,298 times — once per operation — and the six error statuses together account for **13,788 inline copies of five fields**.

What that does to anything reading it: a generator produces sixteen thousand anonymous inline types where eight hundred named ones exist, and the same error envelope becomes 2,298 separate unrelated structs. Nothing can tell that the 401 from one endpoint is the 401 from another, because the document never says they are the same thing. The 13MB is almost entirely that.

**The declared error requires a field the API does not always send.**

Every one of those 13,788 envelopes is `required: [reqId, statusCode, message, error]` with `additionalProperties: false` — the strictest a schema can be. Struck live:

```
GET /api/v1/workspace                  (no credential)
401 {"reqId":"req-us-mlUwgIP7uYC7Za","statusCode":401,
     "message":"Token missing","error":"UnauthorizedError"}

GET /api/v1/nonexistent-cauldron
404 {"message":"Route GET:/api/v1/nonexistent-cauldron not found",
     "error":"Not Found","statusCode":404}
```

**The 404 has no `reqId`.** The document requires it and the API omits it, on the failure a client meets first and most often. Struck twice on two different unknown paths. The key order moves with it: the credential failures lead with `reqId`, this one leads with `message`.

That 404 is Fastify's own not-found handler, printing the method and the path back into the sentence — the framework answering from below whatever stamps `reqId` on everything else. Anything validating responses against this document rejects a real 404 from the real API.

**A credential that cannot be read at all is a 403.**

Three requests, struck live:

| Sent | Status | `error` | `message` |
|---|---|---|---|
| no `Authorization` header | **401** | `UnauthorizedError` | `Token missing` |
| `Bearer not-a-real-token` | **403** | `TokenError` | `The provided access token is malformed. Please use a valid token or generate a new one and try again.` |
| a real-shaped JWT, bad signature | **403** | `TokenError` | `invalid signature` |

403 means the caller was identified and is not allowed. Here it comes back for a string that was never parsed as a credential at all — so a client's "refresh my token" branch, which nearly always hangs off 401, does not fire on the failure a stale token actually produces.

**And the third message is not Infisical's.** "The provided access token is malformed. Please use a valid token or generate a new one and try again." is a written sentence with a suggestion in it. "invalid signature" is `jsonwebtoken`'s own string — lower case, no full stop, no advice — arriving unchanged under the same `error` class. Two failures one signature-check apart, one prose and one library internals.

**`reqId` names the region that answered.** `req-us-mlUwgIP7uYC7Za`, from `us.infisical.com`. The document lists `us.infisical.com`, `eu.infisical.com` and `localhost:8080` as servers and the request id carries which one you reached — genuinely useful, and also a field with no documented format that a user has to copy exactly into a support ticket.

**Not one of the 2317 operations has a `summary`.** Zero, in 12.9MB. And **300 of them have no `operationId`** either, so an eighth of this API's methods are unnameable and every generator invents a name from the path. 2294 carry tags, so the document was not written without care — it was written by a generator that had nowhere to get a sentence from.

**`version` is `0.0.1`**, on the document describing a production secrets manager serving two regions.

**The listing has no paging of any kind.** No limit, no offset, no cursor, no total: `GET /api/v1/projects/{projectId}/roles` answers `{"roles": [...]}` and `additionalProperties: false` leaves nowhere for a cursor to go even if one were sent. A project with a thousand roles sends a thousand roles.

**`description` is nullable and not required.** So a role with no description has two legal representations — the key absent, or the key present and `null` — and the schema permits both on an object that is otherwise strict about every field that may appear. A client testing `'description' in role` and one testing `role.description !== null` disagree about the same record.

## Detection

Four clients across three ecosystems:

| Package | Ecosystem | Evidence |
|---|---|---|
| `@infisical/sdk` | npm | names `app.infisical.com` |
| `infisical-node` | npm | names `app.infisical.com` |
| `github.com/infisical/go-sdk` | Go | names `infisical.com` in `client.go` and `constants.go` |
| `infisical/php-sdk` | Packagist | takes the host as a parameter — see below |

**`infisical/php-sdk` fails the host grep and is mapped anyway.** It hardcodes no host at all: `baseUrl` is a constructor argument, and the only `infisical.com` strings in the archive are documentation links. What it does contain is this API's own paths — `/api/v1/auth/universal-auth/login` and `/api/v3/secrets/raw` — which is the better evidence, and arguably the more correct design for a product that is meant to be self-hosted. The host test is a proxy for "does this call the API"; here the paths answer the real question directly.

## Modelling limits

- **One listing of 2317 operations.** Secrets, environments, folders, identities, tokens, audit logs, dynamic secrets, PKI and the whole approval-workflow surface each want their own evidence. What is here is the project-role listing and the failure path everything shares.
- **The success path is documentation-derived.** Reading anything from Infisical needs an account, so the five cases that answer 200 come from the document and the five that refuse come from the wire.
- **`reqId` rides the three failures that send it, rather than the envelope.** An envelope constant can be added per error but not taken away, and declaring `reqId` for every failure would put it on the 404 — which is the one place the real API leaves it out, and the finding.
- **The duplication is measured, not served.** A fake cannot reproduce "this schema appears 2,298 times"; the counts are in the header of the Recipe and here.
