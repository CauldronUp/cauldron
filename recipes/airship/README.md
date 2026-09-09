# airship

Emulates the Airship push and channels API for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-09.**

Written against Airship's REST API reference at `www.airship.com/docs/developer/rest-api/ua` and struck live against `go.urbanairship.com` on 2026-09-09 with no credential and with a deliberately invalid one.

## What this Recipe found

**The challenge header advertises whichever scheme you already tried.**

```
GET /api/channels/     (no credential)
401  WWW-Authenticate: Basic realm="API"

GET /api/channels/     Authorization: Bearer notreal
401  WWW-Authenticate: Bearer realm="API"
```

This endpoint accepts three schemes — its own reference lists "Basic Auth (Master), Bearer Token, OAuth 2.0", and the tags endpoint beside it lists four. RFC 9110 lets a 401 carry several challenges precisely so that a server can say which ones it takes.

Airship sends exactly one, and picks it by **echoing what the caller already sent**. So the field that exists to tell you which schemes are available will only tell you about Bearer if you had already guessed Bearer. A caller holding an OAuth token and sending nothing is told to use Basic.

It is the second challenge-header fault in this batch, and the exact opposite of [cloudamqp](../cloudamqp): that one sends *no* challenge to a caller with no credential and a correct one to a caller with a wrong credential. One withholds the header from the caller who needs it; this one sends it and fills it in wrong.

**The failure is JSON labelled `text/plain`, and the reference says so.**

```
401  Content-Type: text/plain
{"ok":false,"error":"Unauthorized","error_code":40101}
```

The documentation's own response table for this operation gives the 401 a body of `Content-Type: text/plain`. So the mislabel is written down rather than accidental — a client that checks the type before parsing skips a parse it could have done, exactly as documented.

It is the fifth content-type fault in this collection, after Statsig sending none at all, Wrike serving text over JSON, Maxio serving JSON over text, and Affinity serving HTML over text.

**`error_code` is the HTTP status with two digits stuck on the end.** `40101` is 401 followed by 01. The machine-readable code is not independent of the status; it *contains* it, so a client switching on `error_code` is switching on a number whose leading digits it already had from the status line.

**Three content types across three responses, and the two JSON ones disagree.** The reference gives:

| Response | Documented `Content-Type` |
| --- | --- |
| 200 | `application/vnd.urbanairship+json; version=3` |
| 404 | `application/vnd.urbanairship+json` |
| 401 | `text/plain` |

The version parameter is part of the contract on success and absent on failure.

**And the vendor media type names a company that was renamed.** `urbanairship` is Urban Airship, which became Airship; the host in every example is still `go.urbanairship.com`.

That is the fourth rename survivor here, after [bird](../bird) authenticating with MessageBird's `AccessKey` scheme, [maxio](../maxio) served from `chargify.com`, and [terraformcloud](../terraformcloud) at `app.terraform.io`. It is the first where the old name is baked into a **media type**, which is the hardest of the four to ever change: altering it breaks every client that negotiates on it.

**A custom header names the field to read.** `Data-Attribute: channels` on the listing, `Data-Attribute: channel` on a single lookup. The response says, out of band, which key in the body holds the payload. That is a genuine answer to the envelope problem and nothing else in this collection does it.

**The next page arrives twice, and the reference contradicts itself about one of them.** `next_page` is a full URL in the body, and a `Link` header carries the same address. The response-header table describes `Link` as "Provides the URL to the **current** page of results"; the worked example three paragraphs below shows `rel=next`. Both cannot be right.

**Two examples on one page name the same things differently.** The single-channel lookup shows `"type": "open"` and `"address": "example@example.com"`. The listing shows `"device_type": "android"` and `"push_address": "FE66..."`. Platform and address, under two names each, in two examples of the same resource on one documentation page.

**Three booleans decide whether a device can be reached.** `opt_in`, `installed` and `background` — and a client counting reachable devices has to decide for itself which combination counts.

**The timestamps carry no zone.** `2020-03-06T18:52:59` in the reference's own examples: a `T` in the middle and nothing after the seconds, so they are not RFC 3339.

## Modelling limits

- **One route.** Listing channels. Named users, tags, attributes, subscription lists, push, schedules, segments and the whole automation surface each want their own evidence.
- **The vendor media type is not served on success.** The emulator answers `application/json`; the documented `application/vnd.urbanairship+json; version=3` is recorded above rather than reproduced, because the finding is what the reference declares across its three responses rather than one header value.
- **`Data-Attribute` is recorded, not served.** It is a real answer to the envelope problem and this format has no per-route header for naming the payload key.
- **This Recipe was blocked on documentation until 2026-09-09.** `docs.airship.com` is now a redirect stub whose real content lives at `www.airship.com/docs/developer/rest-api/ua`, which is server-rendered.
- **Nothing is mapped in detection.** Airship's published libraries are named for Urban Airship and target the push send surface rather than this one; a project holding one is sending notifications, not enumerating channels.
