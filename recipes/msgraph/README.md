# msgraph

Emulates the Microsoft Graph API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Microsoft's own generated OpenAPI 3.0.4 document -- `microsoftgraph/msgraph-metadata`, `openapi/v1.0/openapi.yaml`, 44,146,442 bytes and struck live against `graph.microsoft.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**11,493 paths and 17,777 operations, in 44MB.**

The largest description in this collection by a wide margin. [infisical](../infisical)'s is 13MB and 2,317 operations; this is three times the bytes and nearly eight times the operations, in 982,167 lines.

The title is *"OData Service for namespace microsoft.graph"* — what a converter writes when it turns CSDL into OpenAPI. The document is a translation rather than a description anybody wrote.

**10,101 of the 11,493 path keys have to be quoted** — 88% — because OData puts `$`, `(` and `)` in URLs. A YAML reader handles that; a developer grepping the file for a path does not, because the quoting is not what they will type.

**The longest path is 340 characters:**

```
/identityGovernance/entitlementManagement/resourceEnvironments/
{accessPackageResourceEnvironment-id}/resources/{accessPackageResource-id}/
roles/{accessPackageResourceRole-id}/resource/scopes/
{accessPackageResourceScope-id}/resource/uploadSessions/
{customDataProvidedResourceUploadSession-id}/files/
{customDataProvidedResourceFile-id}/$value
```

Five identifiers to name one file, a placeholder called `{customDataProvidedResourceUploadSession-id}`, and `$value` on the end because OData addresses a raw property by suffix rather than by header.

**A wrong token is answered with a library error code and a lesson about JWT serialisation.**

```
401 {"error":{"code":"InvalidAuthenticationToken",
              "message":"IDX14100: JWT is not well formed, there are no dots (.).\nThe token needs to be in JWS or JWE Compact Serialization Format. (JWS): 'EncodedHeader.EncodedPayload.EncodedSignature'. …"}}
```

`IDX14100` is `Microsoft.IdentityModel`'s own error identifier, and everything after it is that library's help text — two format templates and a literal `\n`. It is the most useful refusal in this batch **and** a library internal on the wire: the code is not documented as part of Graph, and it changes when the library does.

The `code` beside it is `InvalidAuthenticationToken` for both a missing token and an unparseable one, so the machine-readable field cannot separate them and the sentence, which can, belongs to somebody else's library.

**`request-id` and `client-request-id` are the same value.** `client-request-id` is meant to echo what the caller sent in the header of that name. A caller who sends none gets Graph's own id in the field named for theirs — two fields, one value, and the one named for the client is the server's.

**The timestamp inside the failure has no timezone.** `2026-09-07T21:07:01` — no `Z`, no offset.

**Three naming conventions in one small object.** `code`, `message` and `innerError` are camelCase; `request-id` and `client-request-id` are hyphenated; `date` is a bare word. All five live inside `error`.

**An unknown path answers 401 rather than 404**, so route existence is not discoverable without a token — which, for a document with 11,493 paths, is the one thing a caller might reasonably want to check.

**A user without a mailbox has two addresses that disagree.** `userPrincipalName` always looks like an email and is not one; `mail` is the address and can be `null`. A client showing "the user's email" has to know which, and only one of them is named for it.

## Modelling limits

- **One route of 17,777.** Users. Mail, calendars, drives, groups, teams, devices, identity governance and the rest of Graph each want their own evidence, and the document describing them is 44MB.
- **Nothing is mapped in detection.** The Microsoft Graph SDKs are generated per language and the .NET and JavaScript ones are named for the platform rather than the API, so no dependency name identifies this surface unambiguously.
- **A `spec:` is declared and fingerprinted.** `cauldron drift` reads it, which for a 44MB YAML document is the slowest check in the suite and the reason it runs on a schedule rather than on every push.
