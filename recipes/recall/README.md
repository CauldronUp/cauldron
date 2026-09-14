# recall

Emulates the Recall.ai bot listing for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from the API reference at [`docs.recall.ai`](https://docs.recall.ai/reference/bot_list), and struck live on 2026-09-13 against two regions with no credential, with a wrong token, with the wrong scheme, with a scheme and no value, with a value and no scheme, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**`Allow` rides on a 401 that is not about methods.** Every refusal carries both of these:

```
WWW-Authenticate: Token
Allow: GET, POST, HEAD, OPTIONS
```

RFC 9110 puts `Allow` on a 405 and nowhere else. Here it is on the credential failure — and `PUT /api/v1/bot/`, a method that is *not* in that list, is refused for the credential rather than for the verb. So the header naming the permitted methods is served beside an answer that never looked at the method.

**The challenge names a scheme that is not registered.** `WWW-Authenticate: Token`, with no realm and no parameters. `Token` is not in the IANA HTTP Authentication Scheme Registry, so a client that implements the challenge has nothing to implement.

**A wrong token is answered with a list of the vendor's regions.**

> "Invalid API token. The API token provided might be for another Recall region (such as us-east-1, us-west-2, eu-central-1, or ap-northeast-1). Verify token is valid and belongs to the correct region."

Four deployment regions, enumerated to an unauthenticated caller, because the service cannot distinguish a token that is invalid from one that is valid somewhere else.

**And the two sentences about missing credentials fire on requests that carry them.** Struck live, all on `GET /api/v1/bot/`:

| `Authorization` header | code | detail |
|---|---|---|
| *(absent)* | `not_authenticated` | Authentication credentials were not provided. |
| `Token` | `authentication_failed` | *the region sentence* |
| `notarealtoken` | `authentication_failed` | *the region sentence* |
| `Bearer notarealtoken` | `authentication_failed` | *the region sentence* |
| `Banana notarealtoken` | `authentication_failed` | Invalid token header. No credentials provided. |

A header holding the word `Token` and nothing else is told its token is invalid. A header holding a scheme *and* a value is told no credentials were provided. And the scheme word is otherwise not read: `Bearer`, a bare value and `Token` all reach the same refusal.

**An unrouted path is a Django HTML page, with a marketing CSP on it.** 404, `text/html`, `<h1>Not Found</h1><p>The requested resource was not found on this server.</p>` — from an API that is JSON everywhere else. The response also carries a 1,700-byte `Content-Security-Policy` header naming several dozen third-party hosts: analytics, session recording, a chat widget, a payment processor, a data-enrichment vendor, a tag manager. On an API's 404, to a caller that is not a browser.

**A parameter's name is a Django ORM lookup.** `metadata__<key>` — double underscore — "Replace `<key>` with the metadata key you want to filter on." The framework's query syntax, exposed as the public filter interface.

**The paging flag is a boolean typed as a string.** `use_cursor` is declared `string` and documented "When present, uses cursor-based pagination instead of offset pagination", so its value is never read — only whether it is there. And the paging it replaces is worked by `page`, "A page number within the paginated result set", which is page-number paging, not the offset paging the flag's own description calls it. Neither mode has a page-size parameter: the reference lists none.

**A field the caller sets is deleted from the record.** `meeting_url` is "The url of the meeting… This field will be cleared a few days after the bot has joined a call." The record stops saying which meeting it recorded — and `meeting_url` is also a listing filter, so the same query stops matching the same bot.

Also pinned: the documented base `api.recall.ai` and the region host `us-west-2.recall.ai` answer identically, and so does `eu-central-1`; a trailing slash makes no difference to the refusal; `bot_name` defaults to "Meeting Notetaker"; and every entry in `status_changes` declares `code`, `sub_code`, `message` and `created_at` as required and read-only, with two of the four null on an ordinary lifecycle event.

## Sources

- [Recall.ai bot list reference](https://docs.recall.ai/reference/bot_list) — the response envelope, the `BotGet` fields, the query parameters, and the 60-requests-per-minute-per-workspace rate limit.
- Live: `us-west-2.recall.ai`, `api.recall.ai` and `eu-central-1.recall.ai`, struck 2026-09-13 with six different `Authorization` headers, an unrouted path, a wrong method, and with and without a trailing slash.

## Modelling limits

- **No description is published.** Recall.ai serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint and `cauldron drift` has nothing to compare.
- **This Recipe reads the scheme word and the provider does not.** Live, `Bearer x`, a bare `x` and `Token x` are judged identically. Here a header without the `Token ` prefix reaches the malformed verdict and gets "Invalid token header. No credentials provided." — which is the right sentence for `Banana x` and the wrong one for `Bearer x`. All three sentences are served; which request draws which is faithful for two of the five headers above.
- **The page size is not documented and not settable.** `limit_param` is declared `-`, which is this format's way of saying the provider accepts no name for it. The default page this Recipe serves is 100 records, which is a choice, not an observation.
- **Nine fields of a much larger record.** `BotGet` also carries `output_media`, `automatic_video_output`, `automatic_audio_output`, `chat`, `automatic_leave`, `zoom`, `google_meet`, `slack_team`, `webex` and `breakout_room`, each its own nested schema.
- **Two routes of many.** The bot listing and one bot. Create, delete, leave-call, the recording and transcript artifacts, calendars and webhooks are the rest.
- **The CSP header is described, not reproduced.** It is 1,700 bytes of third-party hostnames; the finding is that it is there at all.
- **Nothing is mapped in detection.** Recall.ai is reached through a plain HTTP call carrying `Authorization: Token`, which does not resolve to this host through a dependency file. Checked 2026-09-13.
