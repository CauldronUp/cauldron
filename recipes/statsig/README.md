# statsig

Emulates the Statsig Console API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Statsig's Console API reference at `docs.statsig.com` and struck live against `statsigapi.net` on 2026-09-07 with no credential, then with a deliberately invalid one, then on a path that does not exist.

## What this Recipe found

**Three different mistakes, one response, and it is not JSON.**

```
GET /console/v1/gates                (no credential)
GET /console/v1/gates                STATSIG-API-KEY: not-a-real-key
GET /console/v1/cauldron-nope        STATSIG-API-KEY: not-a-real-key

401
(no Content-Type header at all)
Unauthorized
```

Eleven bytes. No object, no field, no quotes — and **no `Content-Type`**, so a client cannot even ask what it was handed. `response.json()` throws; `response.text()` returns a word; and the same word answers "I sent nothing", "I sent the wrong key" and "that endpoint does not exist".

From outside, this API is opaque. A caller without a working key cannot discover whether a route exists, cannot tell a typo from a permissions problem, and cannot parse anything.

**That is a defensible security posture, and it is worth being explicit that it is one.** Most findings in this collection are mistakes. This one is a choice, and the cost of the choice is symmetric: the design that gives an attacker nothing gives a developer nothing. Every Statsig client library ends up with a branch that reads the body as text on a 401 and as JSON on everything else, and the failing path is the one nobody exercises.

**The credential header is `STATSIG-API-KEY`** — screaming snake case with hyphens, a convention nothing else here uses. Header names are case-insensitive so it makes no difference on the wire, and every code sample carries the shouting.

**The key names its own surface.** A Console key begins `console-`, a server key `secret-`, a client key `client-`. Sending a server key to the Console API is refused with the same eleven bytes as sending nothing, so the prefix is the only thing that tells you which of three products a key belongs to.

**A gate's id is the name somebody typed.** `id` and `name` are the same string, so renaming a gate changes its identifier and every reference to it.

**A disabled gate is still being checked.** `isEnabled: false` beside `checksPerHour: 8.25` — a gate that is off still answers every evaluation, and the count includes them. Traffic on a disabled gate is the normal case, not a bug, and a number here does not mean anyone is being let through.

**The last editor is a display name and nothing else.** No id, no email. Two people with the same name are one value, and the field cannot be joined to a user record.

## Detection

`statsig-node`, `statsig/statsigsdk` and `github.com/statsig-io/go-sdk` name `statsigapi.net` and never `console/v1`: they are the evaluation SDKs, a different surface on the same host with a different key family — a `secret-` key is refused by the Console API exactly like no key at all. Mapped anyway, because the question detection answers is whether a project talks to Statsig. The same call the Circle wallets SDK gets, and the opposite of [imgix](../imgix), where the sibling product is on a different host entirely.

## Modelling limits

- **The missing `Content-Type` is recorded, not served.** The real refusal carries no such header; this format has no way to omit one and the emulator sends `text/plain`. The cases assert the status and the eleven bytes, which are the parts that can be reproduced honestly. One provider needing a way to suppress a header is not enough to add one.
- **One route.** Gates. Experiments, dynamic configs, segments, metrics, exposures and the whole holdout surface each want their own evidence.
- **No `spec:`.** Statsig publishes a rendered documentation site.
