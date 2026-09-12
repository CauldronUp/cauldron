# instantly

Emulates the Instantly v2 API for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-09-11.**

Read from Instantly's published OpenAPI document — 4.2MB of it, served without a key — and struck live against `api.instantly.ai` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The three keys of a failure arrive in one order, or the other.**

```
(no header)      401 {"statusCode":401,"error":"Unauthorized","message":"Missing authorization header"}
Bearer notreal   401 {"statusCode":401,"error":"Unauthorized","message":"Invalid API key"}
/cauldron-nope   404 {"message":"Route GET:/api/v2/cauldron-nope not found","error":"Not Found","statusCode":404}
```

Same three keys, exactly reversed. The 401s come from the application and the 404 from Fastify underneath it, and the two serialisers disagree about which key leads.

Nothing reads JSON positionally, so nothing breaks — but a fixture recorded from one path does not match a body from the other, and every byte-comparing test written against this API is pinned to whichever failure it happened to capture.

**The 404 reads the request back.** `Route GET:/api/v2/cauldron-nope not found` names the method and the path. The published document's example for a 404 is `Resource not found`, so the message a client actually sees for a mistyped path is not the one the document shows.

**An account has no identifier.** The `Account` schema has no `id` property at all. The record is keyed by `email` — which is also the field most likely to change when a person is renamed, so the primary key of an email account is an address that can move.

**The cursor's separator is the query string's separator.** `starting_after` is documented as `timestamp_created&email` format, and the document's own example value is:

```
2026-01-01T00:00:00.000Z&jon@doe.com
```

Sent as written, the `&` ends the parameter and `jon@doe.com` becomes a parameter of its own with no value. It has to be percent-encoded first, and the example is the form that does not work.

**And the field carrying it back has three possible types.** The document's description of `next_starting_after`, verbatim:

> this could either be a UUID, a timestamp, on an email depending on the specific API

One string field, three shapes, chosen by which endpoint you called — typo included. This listing's example is a UUID; the parameter it feeds is documented as the timestamp-and-email pair.

**Two enums carry their meaning in the sign.**

| field | values | what is below zero |
| --- | --- | --- |
| `status` | 1, 2, 3, -1, -2, -3 | connection error, soft bounce error, sending error |
| `warmup_status` | 0, 1, -1, -2, -3 | banned, spam folder unknown, permanent suspension |

The names live in `x-enumDescriptions`, a vendor extension. A generated client gets unnamed integers, and the reader has to know that negative means broken — in two enums whose zero means different things, since `warmup_status: 0` is *Paused* and `status` has no zero at all.

**`autofix_failed` is a boolean with three states.** Its own description: `null = in progress, true = failed, false = succeeded`. So null is not *unknown*, it is *still working* — and a client reading absent-or-null as "fine" reports a reconnection that has not finished as one that worked.

**`status_message` is the mail server's own error, passed through.** `code`, `command`, `response` and `responseCode` come straight off the SMTP conversation — `EENVELOPE`, `DATA`, `550-5.4.5 Daily user sending limit exceeded` — with `additionalProperties: true`. The JSON API's error detail for an email account is an uninterpreted slice of a different protocol.

**And the schema guarantees the wrong things.** `setup_pending` and `is_managed_account` are required; `status` — whether the account can send at all — is not.

## Sources

- Live: `api.instantly.ai`, struck 2026-09-11.
- [`api_v2.json`](https://api.instantly.ai/openapi/api_v2.json) — the OpenAPI document, named inside a fenced block in the documentation rather than linked from it.
- [List account](https://developer.instantly.ai/api-reference/account/list-account) — the page that names it.

## Modelling limits

- **One route.** Accounts. Campaigns, leads, lead lists, emails, the Unibox, inbox-placement tests, webhooks and the several AI agents each want their own evidence — the document describes well over two hundred operations.
- **The emulator emits one key order.** The live API answers two, depending on which layer raised the failure. A Recipe declares one shape per error, so the 404 here arrives in the same order as the 401 rather than reversed; the reversal is the finding, not something the sandbox can reproduce.
- **The 404's message is the recorded one.** Live it names the method and path that were asked for. Here it carries the path this Recipe struck.
- **429 is not modelled.** The document declares it with a `Rate limit exceeded` message; nothing in this Recipe provokes one, and the live API was not pushed into answering one.
- **Nothing is mapped in detection.** Instantly publishes no first-party client library, and the community packages on npm and PyPI are unofficial wrappers rather than a dependency that names the host. Checked 2026-09-11.
