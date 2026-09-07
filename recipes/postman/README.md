# postman

Emulates the Postman API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Postman's reference at `learning.postman.com` and struck live against `api.getpostman.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The first provider in this collection to use `type` for what it is for.**

```
GET /cauldron-nope
404  Content-Type: application/problem+json
{"type":"https://api.postman.com/problems/not-found",
 "instance":"/cauldron-nope",
 "detail":"The requested resource could not be found.",
 "status":404,
 "title":"Not Found"}
```

All five of RFC 9457's members, and `type` is **a URI under Postman's own domain naming a Postman problem**. `instance` is the path that produced it — the other member almost nobody sends.

Three other providers here reach for the same specification and all three miss:

| | `type` |
|---|---|
| [triggerdev](../triggerdev) | MDN's page about the HTTP 401 status |
| [agicap](../agicap) | the RFC clause that *defines* 500 |
| [revai](../revai) | absent |
| **postman** | **`https://api.postman.com/problems/not-found`** |

It is worth recording that the specification is usable, because the other three make it look as though it is not.

**And the 401 is not a problem document at all.** `application/json`, a nested `error`, a class name and a sentence. So one API has a carefully-built problem document for routing failures and an ordinary envelope for credential ones, and a client needs both parsers.

**The key order flips between the two 401s.**

```
no credential   {"error":{"name":"AuthenticationError","message":"Invalid API Key. …"}}
wrong key       {"error":{"message":"Invalid API Key. …","name":"AuthenticationError"}}
```

Same fields, same values, opposite order. Two code paths building one object — the third instance in this collection after [zuora](../zuora) and [browserbase](../browserbase), and only a byte comparison notices.

**"Invalid API Key" is the answer to sending no API key** — the fourth provider here to write that sentence for a request that carried none, after [loops](../loops), [helicone](../helicone) and [beehiiv](../beehiiv).

**Every collection carries two identifiers.** `uid` is the owner's numeric id and the collection's UUID joined by a hyphen, and it is the one most endpoints take — so the field called `id` is not the handle, and building `uid` requires knowing who owns the thing.

## Modelling limits

- **One route.** Collections. Workspaces, environments, monitors, mocks, APIs and the whole SCIM surface each want their own evidence.
- **Nothing is mapped in detection.** Postman is used through its app and its CLI rather than through a client library, and nothing on any registry calls `api.getpostman.com` under an obvious name.
- **No `spec:`.** Postman publishes a rendered reference site, and its own public workspace holds a Postman collection rather than an OpenAPI document.
