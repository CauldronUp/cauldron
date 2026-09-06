# telegram

Emulates the Telegram Bot API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-06.**

Written against Telegram's own documentation at `core.telegram.org/bots/api` — 2850 documented fields, read on 2026-09-06 — and struck live against `api.telegram.org` on the same day. Every live request used a deliberately invalid token: strings of the right shape that identify no bot.

## What this Recipe found

**The credential is the route.**

```
GET https://api.telegram.org/bot<token>/getMe
```

The token is the first path segment, prefixed with the literal `bot`. It cannot be set once as a default header, cannot be attached by an interceptor, and lands in every access log, proxy log and browser history that records a URL. Every response also carries `Access-Control-Allow-Origin: *`, so a browser can call this API directly — with the bot's token in the address bar.

This is the **fifth** Recipe here whose credential is a path segment, after [alchemy](../alchemy), [pubnub](../pubnub), [thesportsdb](../thesportsdb) and Inngest. Each of those four wrote a paragraph explaining that none of the format's auth schemes could describe it, and worked around it — PubNub by "reproducing a wrong key as a literal route matching one exact bad value". This one is served by `scheme: path`, which those four are why the format now has.

**The shape of the token decides whether the URL exists.** Struck live:

| Token sent | Status | Body |
|---|---|---|
| `123456789:AAFakeToken…` | **401** | `{"ok":false,"error_code":401,"description":"Unauthorized"}` |
| `123456789:AAF` | **401** | same |
| `abc:def` | **404** | `{"ok":false,"error_code":404,"description":"Not Found"}` |
| `123456789` (no colon) | **404** | same |
| `nope` | **404** | same |

A token whose first part is not digits, or which has no colon, makes the path not exist. A well-formed token nobody issued is refused. So **the same wrong credential produces 404 or 401 depending on how it is misspelled**, and a developer reporting "that endpoint is gone" is reporting a typo in their own secret.

`bot0:0` answers **401 "Unauthorized: invalid token specified"** — a longer sentence, same status, for bot id zero alone. Recorded and not served: one special-cased value with no way to say so in this format.

**An unknown method with a bad token is 401, not 404.** `/bot<wrong>/thisMethodDoesNotExist` answers Unauthorized. The token is checked before the method, so nobody without a working token can discover whether a method exists. That also means the documented case-insensitivity of method names — "All methods in the Bot API are case-insensitive" — **could not be confirmed from outside**: `getMe`, `getme`, `GETME` and `gEtMe` were all tried and all four answered the credential failure first. It is recorded here as the documentation's claim, not as an observation.

**DELETE, PUT and PATCH answer 501 with an empty body.**

```
DELETE /bot<token>/getMe
HTTP/1.1 501 Not Implemented
Content-Length: 0
(no Content-Type)
```

Not 405 — **501**, which means the server implements the method for no resource at all. And the body is empty with no content type, in an API whose entire contract is *"The response contains a JSON object, which always has a Boolean field `ok`"*. These three are the only responses this API sends that a client cannot call `.json()` on, and nothing in the contract warns about them.

**Telegram disclaims its own error code.** Verbatim from the documentation:

> "An Integer `error_code` field is also returned, but its contents are subject to change in the future."

So the one machine-readable part of a failure is the part the provider tells you not to depend on, and what is left is `description` — English prose. Every Telegram bot library in the world switches on those strings.

**94 fields are typed `True`.** Not `Boolean`: `True`, a type with a single inhabitant, where the field is present when it holds and **absent** when it does not. The same document types 362 other fields `Boolean`. Both encodings describe the same kind of fact and only the type column says which — `is_premium` is `True` and `is_bot` is `Boolean`, on the same `User` object. A client reading `user.is_premium === false` never matches; the answer it wants is `!('is_premium' in user)`.

**100 fields are typed `Integer or String`** — a union declared as the type, so a client has to accept either for one field.

**Reading the second page destroys the first.**

`getUpdates` takes an `offset`, and Telegram's own description of it:

> "An update is considered confirmed as soon as `getUpdates` is called with an offset higher than its `update_id`."

> "The negative offset can be specified to retrieve updates starting from `-offset` update from the end of the updates queue. **All previous updates will be forgotten.**"

So this GET mutates, and destructively. Paging forward permanently discards everything behind the cursor; there is no way to re-read a page; and a *negative* offset — a legal value for a parameter called `offset` — forgets the queue. A client that pages ahead before it has finished processing has destroyed the data it skipped, and nothing about the request said so.

**Identifiers are sized to fit JavaScript.** From the documentation of the id fields: *"at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier."* Telegram constrains its own identifier space to what a double holds exactly, and says so. Almost every other provider in this collection sends an int64 as a string instead and leaves the reader to work it out — [etcd](../etcd), [waypoint](../waypoint), [firebaseauth](../firebaseauth), [gmail](../gmail).

**Chat ids are negative and user ids are not.** A supergroup is `-1001234567890` and a user is `55501`, so **the sign of an integer is what says which kind of thing it names** — on a field called `id`, inside objects that both have one.

## Detection

Five clients across three ecosystems, each checked by fetching its published archive and counting `api.telegram.org`:

| Package | Ecosystem | Occurrences |
|---|---|---|
| `grammy` | npm | 7 |
| `node-telegram-bot-api` | npm | 4 |
| `telegraf` | npm | 2 |
| `irazasyed/telegram-bot-sdk` | Packagist | 3 |
| `github.com/go-telegram-bot-api/telegram-bot-api/v5` | Go | 2 files |

## Modelling limits

- **`getUpdates` is served unpaged, deliberately, and this is the interesting decision here.** Telegram's `offset` is an `update_id` and it is **inclusive** — "Identifier of the first update to be returned" — while every cursor this format serves is exclusive, resuming *after* the record named. The two differ by exactly one record, on the endpoint a bot calls in a loop forever, and an off-by-one in a fake is the failure this collection exists to prevent. So the cursor is not served rather than served wrong. There is no inclusive-cursor option in the format today; one provider needing one is not enough to add one, and the case is recorded for the next. The other half — that sending the parameter deletes data — has no counterpart in this format and should not have one.
- **`bot0:0` is recorded and not served.** One credential value with its own sentence under the same status.
- **Case-insensitive method names are the documentation's claim.** The credential is checked first, so it cannot be observed from outside without an account.
- **Two methods of about a hundred.** `getMe` and `getUpdates`, and the failure path they share with everything else. `sendMessage`, the whole update taxonomy, webhooks, files, inline queries and payments each want their own evidence.
- **No `spec:`.** Telegram publishes this API as an HTML page and nothing else. Third-party OpenAPI translations exist and none is Telegram's, so fingerprinting one would record somebody else's reading of the documentation rather than the documentation.
