# smartlead

Emulates the Smartlead campaigns API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-11.**

Read from Smartlead's own documentation — which answers Markdown at every address with `.md` appended, and publishes a single-file `llms-full.txt` — and struck live against `server.smartlead.ai` on 2026-09-11 with no credential, with an empty one, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**One page describes the response two ways and contradicts itself twice.** The reference for this endpoint carries, in this order:

| where | what it says |
| --- | --- |
| Key Features bullet | "Returns direct array of campaign objects (not wrapped)" |
| Response section | `success` (boolean), `campaigns` (array of campaign objects) |
| PHP sample | `count($data['campaigns'])` |
| Response example | `[ { "id": 2710262, … } ]` — a bare array, labelled "Actual from API" |
| Note beneath it | "returns a direct array of campaigns, not wrapped in a success object" |

Four statements about one body, two of them wrong. The sample code that ships beside the example counts a key the response does not have.

**Every failure is a different shape depending on which page you read.**

```
the error-handling guide:  {"error": {"code": …, "message": …, "details": {…}}}
this page's 401 example:   {"message": "Invalid API Key"}
this page's 404 example:   {"error": "Resource not found"}
this page's 422 example:   {"error": "Invalid parameters provided"}
```

So `error` is an object on one page and a string on another, and the guide calls its version "a consistent JSON structure". The guide's own Python sample reads `error_body.get('error', {}).get('message', 'Unknown')` — which against the live API prints `Unknown` for every failure there is.

**Live, there is one shape and two sentences.**

```
(no api_key)          401 {"message":"API key is required."}
?api_key=             401 {"message":"API key is required."}
?api_key=notreal      401 {"message":"Invalid API Key"}
/cauldron-nope        401 {"message":"API key is required."}
```

`API key is required.` has a trailing full stop and a lower-case k. `Invalid API Key` has neither, and title-cases both words. One field, one status, one endpoint, two house styles — and an empty `api_key=` counts as absent rather than wrong.

**And an unrouted path is a credential failure.** Authentication runs before routing, so a mistyped path answers 401 carrying the sentence for whichever credential state the caller is in. The documented 404 is not reachable by getting the path wrong, and a client that retries its credentials on a 401 will retry forever against a typo.

**The status enum disagrees with itself between pages.**

| page | values |
| --- | --- |
| Understanding Campaigns | DRAFT, ACTIVE, PAUSED, STOPPED, ARCHIVED, COMPLETED |
| this endpoint's reference | ACTIVE, PAUSED, STOPPED, ARCHIVED, DRAFTED |

Six against five, `DRAFT` against `DRAFTED`, and `COMPLETED` in one place only — on the field whose whole job is to be compared against a literal.

**The page's own Response Codes section omits two statuses it gives examples for.** It lists 200, 401 and 500. The examples above it include a 404 and a 422.

**The credential is a query parameter**, and the documentation offers no header form — so the secret is in every URL, every access log, and every `Referer`.

**And one object in the record is camelCase.** Every field on a campaign is snake_case — `user_id`, `track_settings`, `min_time_btwn_emails` — except inside `scheduler_cron_value`, which holds `startHour` and `endHour` beside `tz` and `days`.

**`track_settings` is a list of negatives**, so an empty list means track everything and two entries mean track nothing. A client reading "no tracking settings" as "no tracking" has it exactly backwards.

## Sources

- Live: `server.smartlead.ai`, struck 2026-09-11.
- [`llms-full.txt`](https://api.smartlead.ai/llms-full.txt) — the record shape, the field tables and the response examples.
- [Understanding Campaigns](https://api.smartlead.ai/core/campaigns) — the other status enum.
- [Error Handling Guide](https://api.smartlead.ai/guides/error-handling) — the envelope the API does not send.

## Modelling limits

- **One route.** Listing campaigns. Leads, email accounts, sequences, webhooks, analytics and the client/agency surface each want their own evidence.
- **The documented 404 and 422 are not modelled.** Both are described with a shape (`{"error": "…"}`) that no live probe produced, and nothing this Recipe can do provokes one — a mistyped path answers 401 instead.
- **The envelope is the live one, not the documented one.** This serves the bare array, because that is what the page's own example and its Note say the API sends, and the field table naming `success` and `campaigns` is the part that is wrong.
- **No `spec:`.** Smartlead documents this API as prose and field tables. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** Smartlead publishes no first-party client library, and the packages named for it on npm and PyPI are unofficial wrappers a project may or may not hold. Checked 2026-09-11.
