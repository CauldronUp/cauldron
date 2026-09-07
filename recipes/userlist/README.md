# userlist

Emulates the Userlist API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Userlist's reference at `userlist.com/docs` and struck live against `api.userlist.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**A credential failure is reported as a missing resource.**

```
GET /users    (no credential)
GET /users    Authorization: Push notreal
404 {"status":404,"code":"not_found",
     "errors":["The requested resource could not be found on the server."]}
```

Not 401, not 403 — **404**, with a sentence about a resource. The route exists; the caller's permission to see it does not; and the API answers by claiming the thing is not there.

It is the second provider here to do this, after [eightbyeight](../eightbyeight), and the two arrive at it differently. 8x8's 404 is Kong refusing to route a request without a matching consumer, so the gateway genuinely has no route. **Userlist's application answers with its own envelope and its own code, so this is a decision rather than a side effect.**

And it is arguably the right decision. Answering 401 to an unauthenticated request confirms that `/users` exists; answering 404 does not. Several providers in this collection worry about disclosure and none of the others follow through this far. The cost is that a developer with a working key and a typo, and a developer with a wrong key and a correct path, are told the same thing.

**The credential scheme is `Push`.** `Authorization: Push <token>` — a scheme name Userlist invented, which no HTTP library has a helper for. It is at least unambiguous: a token sent as `Bearer` is not silently half-accepted.

**`errors` is an array of bare strings** — the third shape for that field name in this collection:

| | `errors` is |
|---|---|
| [honeybadger](../honeybadger) | a single string |
| [gong](../gong) | an array of strings |
| **userlist** | an array of strings |
| [beehiiv](../beehiiv) | an array of **objects** with a code each |

Only the last can carry a machine-readable cause per failure.

**`status` and `code` say the same thing twice**, both derivable from the status line, beside the array carrying the only sentence.

**Properties are an untyped object that can be empty** — `{}` rather than `null` when nothing was set, so a client cannot tell "no properties" from "never set any".

## Modelling limits

- **One route.** Users. Companies, relationships, events, campaigns and the whole message surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Userlist publishes a rendered reference site.
