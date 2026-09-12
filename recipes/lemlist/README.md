# lemlist

Emulates the lemlist campaigns API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-11.**

Read from the OpenAPI document lemlist embeds in each documentation page, and struck live against `api.lemlist.com` on 2026-09-11 with no credential, with a deliberately invalid one, with a malformed one, and on a path that does not exist.

## What this Recipe found

**A path that does not exist answers 200 and sends the web application.**

```
GET /api/cauldron-nope   200 text/html   <!DOCTYPE html> … (the SPA)
```

Not a 404 — a success, carrying the lemlist single-page app: script tags, stylesheets, a tag-manager block, and the client bootstrap configuration with analytics identifiers, a publishable payment key, internal service hostnames and the deployed commit hash in it.

So a client that mistypes a path gets `response.ok === true`, a `.json()` that throws, and nothing anywhere in the answer saying the path was wrong. Every other Recipe in this collection describes what a provider does *wrong* on a 404; this one has to describe a provider that does not send one.

**And none of the real failures is JSON either.** The document declares all four as `text/plain`, one sentence each:

| status | sentence | what it actually means |
| --- | --- | --- |
| 400 | `no api key provided` | the header is absent, or holds no password |
| 401 | `The authentication you supplied is incorrect` | the key is wrong |
| 403 | `User linked to this API key is blocked` | the person is blocked |
| 404 | `No user found for this API key` | the key has nobody behind it |

Four statuses describing four states of one credential. A missing key is a **400**, so a client branching on 401 never sees it, and a **404 on a listing** means the key is orphaned rather than that anything is missing. The live 400 arrives with no `Content-Type` header at all.

**The instructions call a colon a semicolon.** Verbatim, from the authentication section:

> Don't forget to add the semicolon (`:`) before your API key in curl command

The character in the code span is right and the word in front of it is wrong. A reader who follows the prose types `;`, which sends the key as the *username* with an empty password — and lemlist reads only the password half.

Struck live, a Basic header carrying a key with no colon answers `400 no api key provided`, identically to sending no header at all. So the typo does not produce "wrong key". It produces "no key", which is the one message that will not send a reader to look at their credential.

**`X-RateLimit-Reset` is a JavaScript timestamp.** Struck live on the 401:

```
x-ratelimit-reset: Sat Sep 12 2026 04:07:49 GMT+0200 (Central European Summer Time)
```

That is `Date.prototype.toString()` in an HTTP header. Not an HTTP-date, not a Unix second count, and it ends in a parenthesised timezone *name* that changes with the season — so the one field telling a client when it may retry is the one field no date parser will read.

`Retry-After: 2` comes with it on a request that was not rate limited at all — `x-ratelimit-remaining: 19` of 20 — so the header meaning "wait" is sent to callers who need not.

**The version parameter is optional, defaults to `v2`, accepts only `v2`, and the page warns you to send it.** The schema says `default: v2` and `enum: [v2]`; the prose above it says "Don't forget to set the query parameter version to `version=v2`". A default nobody may override should never need a warning, and the warning is the evidence that it does not apply.

**Two pagination schemes sit on the same listing.** `offset` with `limit`, and `page`, declared side by side with nothing saying which wins.

**And three of the four campaigns in the success example are broken.** The happy-path sample for listing campaigns carries `hasError: true` on three records, with the vendor's own sentences inside — "Your campaign does not have sender.", "One of your sender has no email provider.", "Your campaign have an invalid sender mailbox". The fourth omits `hasError` rather than sending false, so an absent key is the only way this API says a campaign is fine.

## Sources

- Live: `api.lemlist.com`, struck 2026-09-11.
- [Get Many Campaigns](https://developer.lemlist.com/api-reference/endpoints/campaigns/get-many-campaigns) — the page, and the OpenAPI document embedded in it.

## Modelling limits

- **One route.** Campaigns. Leads, activities, unsubscribes, team, senders, credits, schedules, sequences and the enrichment surface each want their own evidence.
- **The 404 is described, not served.** `No user found for this API key` is a state of a credential rather than of a request, and nothing this Recipe can do to a fixture puts a key into it.
- **The unrouted 200 serves an opening tag, not the app.** Live it is the whole single-page application; here it is `<!DOCTYPE html>` under `text/html`, which is enough to reproduce what breaks — a truthy `response.ok` and a `.json()` that throws.
- **No `spec:`.** lemlist embeds a complete OpenAPI document per endpoint inside that endpoint's page, the same arrangement [mailtrap](../mailtrap) uses. There is no single document to fingerprint.
- **Nothing is mapped in detection.** The npm and PyPI packages named for lemlist are unofficial wrappers, and the supported integrations are Zapier, Make and native CRM connectors rather than a library a project would depend on. Checked 2026-09-11.
