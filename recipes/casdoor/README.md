# casdoor

Emulates the Casdoor identity-provider API, for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Casdoor's own generated Swagger 2.0 document, published in its own repository — `casdoor/casdoor`, `swagger/swagger.json`, 234 paths, version 1.503.0 — read on 2026-09-06. Casdoor is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**The user schema declares a password, a salt and a TOTP secret — and the document uses `writeOnly` nowhere.**

`object.User` has **181 properties**. Among them:

```
password          accessSecret     totpSecret
passwordSalt      accessToken      originalToken
passwordType      accessKey        originalRefreshToken
```

`GET /api/get-user` declares `object.User` as its 200 response schema, and `writeOnly` — OpenAPI's way of saying "may be sent, is never returned" — appears **zero times in 350 kilobytes of JSON**.

So as published, the document says a user fetch returns that user's password, their password salt, their TOTP secret and three tokens. On an identity provider. A generated client's `User` type carries all nine; a logged response body carries them; a cache holds them.

**This is a claim about the document, not about the running server.** Casdoor may well strip these in the handler — most do — and there is no public instance to check against, so this Recipe does not claim it either way. What is checkable is what the document says, and the document says they are returned. The fixture uses obvious placeholders.

[coolify](../coolify) records the same shape on a platform API; [exoscale](../exoscale) records the opposite — one correct `writeOnly`, hiding firewall rules.

**Deletion is a POST to a path whose name says delete.**

```
POST /api/delete-user
POST /api/delete-application
POST /api/delete-cert
```

Twenty-eight `delete-*` paths, every one of them POST. The verb lives in the path and the method disagrees with it, so a reviewer grepping for `DELETE` finds nothing, a proxy rule on the method never fires, and an audit classifying requests by method files every deletion as a write. The `get-*` paths *are* GET — so the method is redundant with the path name on 104 endpoints and contradicts it on 28.

**The response envelope has three numbered payload slots.**

```
controllers.Response:
  status  string
  msg     string
  data    "support string, struct or []struct"
  data2   "support string, struct or []struct"
  data3   "support string, struct or []struct"
```

`data`, `data2` and `data3`, each with the same description and no type. Which one holds what depends on the endpoint, and the document does not say for any of them.

**And `status` is a string in the body** — `"ok"` or `"error"`, carried at HTTP 200. So a failure is a success as far as the transport is concerned: `response.ok` is true, `status < 300` passes, `raise_for_status()` does not raise.

**The document ships two definitions the generator produced by accident:**

```
"232967.<nil>.string"
"233025.string.string"
```

Schema names containing a line number, a Go nil and a type name. `<nil>` is not a legal identifier fragment anywhere, and a generator that turns schema names into types has to invent something for both.

**One Casdoor serves many issuers.** Beside `/.well-known/openid-configuration` sits `/.well-known/{application}/openid-configuration`, and the same for JWKS and WebFinger — so the OIDC discovery document depends on which application is asking, and a client fetching the unqualified one gets a different issuer from the one it was configured for.

**A user is identified by a pair.** Every record is scoped by organisation, so the real key is `(owner, name)` — and the scope travels as a query parameter rather than a path segment, so two calls differing only in a query parameter address different directories.

**`isAdmin` and `isForbidden` are separate axes.** A record can be neither, either or both, and nothing collapses them into a role.

## Modelling limits

- **Nothing here is verified against a live API.** Casdoor is self-hosted and there is no public instance.
- **Users, listed, fetched and deleted.** 234 paths is an identity provider: applications, organisations, providers, tokens, certs, enforcers, groups, roles, permissions, the OIDC surface and the Casbin adapter each want their own evidence.
- **The secret-shaped fields are served with obvious placeholders.** What is reproduced is that the keys are on a response — the finding — not a credential worth having.
- **Nine of the 181 user fields are modelled**, so this Recipe badly understates how much a user record carries.
