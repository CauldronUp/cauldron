# pulumi

Emulates the Pulumi Cloud REST API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Pulumi's reference at `pulumi.com/docs` and struck live against `api.pulumi.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The message says "or".**

```
GET /api/user/stacks    (no credential)
GET /api/user/stacks    Authorization: token pul-notreal
401 {"code":401,"message":"Unauthorized: No credentials provided or are invalid."}
```

Byte-identical both ways — and the sentence covers both cases explicitly, joined by a conjunction.

That is honest. The server genuinely is not telling you which, and it says so. Most providers in this collection collapse the same two failures into one message and quietly pick a side, usually the wrong one: [loops](../loops) answers "Invalid API key" to a request carrying no key at all, and [helicone](../helicone) answers "No API key found" to one that carries a key. Pulumi names both.

It is better, and it is still not actionable. A client cannot branch on a conjunction.

The status also appears three times in one exchange: as the HTTP status, as the numeric `code`, and as the word "Unauthorized" before the colon.

**An unknown path is `text/plain`, and it is Go's.**

```
GET /api/cauldron-nope
404  Content-Type: text/plain
404 page not found
```

`http.NotFound`'s exact output from the standard library, left unreplaced. So the JSON envelope above is produced only by routes that exist, and the failure a client meets while getting a URL wrong is the one `.json()` throws on.

**The credential scheme is `token`, not `Bearer`.** GitHub popularised it and almost nothing else uses it, so every HTTP library's built-in bearer helper produces a header this API refuses — with a sentence that does not mention the scheme.

**A stack has no identifier.** It is addressed by the triple organization / project / stack, none of which is unique alone: two organizations can each have a `platform/prod`.

**A stack that has never run omits `lastUpdate` entirely** rather than sending zero — which is the right call, and it means a client sorting by last update has to decide where an absent value goes rather than finding it silently placed in 1970.

**The update time is seconds, not milliseconds.** Ten digits where most of this collection sends thirteen.

## Modelling limits

- **One route.** Stacks. Deployments, updates, previews, policy packs, ESC environments and the whole state-history surface each want their own evidence.
- **Nothing is mapped in detection.** The Pulumi CLI and the language SDKs build infrastructure rather than calling this REST surface, and nothing on npm, Packagist or the Go module proxy calls `api.pulumi.com` under an obvious name.
- **No `spec:`.** Pulumi publishes a rendered reference site.
