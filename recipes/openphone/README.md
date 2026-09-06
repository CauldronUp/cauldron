# openphone

Emulates the OpenPhone API for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-06.**

Written against OpenPhone's API reference and struck live against `api.openphone.com` on 2026-09-06 with no credential and then with a deliberately invalid one.

## What this Recipe found

**A failure is JSON and a not-found is HTML.** Struck live:

```
GET /v1/phone-numbers                    (no credential)
401  application/json
{"error":{"message":"Missing authorization header","key":"Unauthorized",
          "trace":"5640373784117800080"}}

GET /v1/cauldron-nonexistent
404  text/html; charset=utf-8
<!DOCTYPE html><html lang="en">…
```

So the one failure a client is guaranteed to meet while getting a URL right — getting it wrong first — is the one that is not JSON. Every `.json()` in both the happy path and the error path throws there, and the content type says so a beat too late for code that does not check it. It is Express's default handler underneath, showing through the API.

**`trace` is a 64-bit integer sent as a string.** `"5640373784117800080"` — 19 digits, quoted. Nineteen digits do not survive a double, which is *why* it is quoted, and it means the correlation id cannot be compared numerically or stored as a number. A client that does either silently corrupts the value it would quote to support, which is the one job that field has.

**`key` is the machine-readable field and it says the same thing twice.** Struck live:

| Sent | `message` | `key` |
|---|---|---|
| nothing | `Missing authorization header` | `Unauthorized` |
| `Authorization: not-a-real-key` | `Unauthorized` | `Unauthorized` |

So the field a client is meant to branch on is `Unauthorized` for both, and the only thing separating "you sent nothing" from "you sent the wrong thing" is the prose — where the missing case is precise and the wrong case is the key repeated. A client switching on `key` cannot tell them apart. A client switching on `message` gets one useful sentence and one echo.

**The credential carries no scheme.** `Authorization: <key>`, bare — no `Bearer `. OpenPhone's documentation says so, and every HTTP library's built-in bearer helper gets it wrong by adding one. That is served here rather than described: a fake that accepted `Bearer <key>` would let code ship that the real API refuses.

**A number is sent twice in two formats.** `number: "+15550100"` and `formattedNumber: "(555) 010-0"`, as separate fields. Only one of them is safe to send back in a request, and nothing in either name says which.

## Modelling limits

- **One route.** Phone numbers. Messages, calls, contacts, conversations, webhooks and call transcripts each want their own evidence.
- **The 404 is served as an empty body rather than as HTML.** What the case pins is that it is not JSON; reproducing Express's default page byte for byte would be asserting a template rather than a contract.
- **Nothing is mapped in detection.** `@openphone/node` and `openphone` are both 404 on npm as of 2026-09-06, and no client turned up on Packagist or the Go module proxy.
- **No `spec:`.** OpenPhone publishes a rendered documentation site.
