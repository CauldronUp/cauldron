# e2b

Emulates the E2B sandbox API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against E2B's reference at `e2b.dev/docs` and struck live against `api.e2b.dev` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The best malformed-credential message in this collection.**

```
GET /sandboxes    X-API-KEY: notreal
401 {"code":401,
     "message":"API key is malformed: expected the \"e2b_\" prefix,
                visit https://docs.e2b.dev/api-key for more information"}
```

It names the failure (**malformed**, not unrecognised), names the expected prefix (`e2b_`, in quotes), and links a page about API keys specifically. Three useful facts in one sentence — a caller who pasted a project id, an access token or a truncated key knows immediately which of those they did.

Nothing else here manages all three:

| | Names the failure | Names the format | Links the docs |
|---|---|---|---|
| [nango](../nango) | yes | yes (UUID v4) | no |
| [browserbase](../browserbase) | yes | no (names the *header*) | no |
| [saltedge](../saltedge) | yes | no | yes (exact anchor) |
| **e2b** | **yes** | **yes** | **yes** |

**And the missing case is accurate too** — `authorization header is missing` for a request that carried none, where five providers in this collection answer "invalid key" to exactly that.

So E2B gets right the two things this collection has watched most providers get wrong, in nine words each.

**The routing failure is a request validator.** `validation error: no matching operation was found` — "operation" is OpenAPI's word for a path-and-method pair, so E2B validates every request against its own document before routing it. **The document is the router**, and an unknown path is a specification mismatch rather than a missing resource.

**Two sandboxes on one team can cost different amounts** — four times the memory and twice the CPU, in one listing, with nothing in the envelope totalling it.

**`clientID` identifies the orchestrator node**, not the team or the key — so it changes when E2B reschedules, and a client keying on it keys on infrastructure.

## Modelling limits

- **One route.** The sandbox listing. Templates, filesystem operations, process execution, the CDP bridge and the whole metrics surface each want their own evidence.
- **Nothing is mapped in detection.** E2B's SDKs are named for the company and drive sandboxes through a websocket rather than this REST surface.
- **No `spec:`.** E2B's OpenAPI document lives in its repository rather than at a stable published URL — and the API validates against it without serving it, which is the finding above.
