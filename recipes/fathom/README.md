# fathom

Emulates the Fathom Analytics sites API for local development and tests.

**8 conformance cases, 2 checked against the live API on 2026-09-13.**

Read from Fathom's own API documentation, and struck live against `api.usefathom.com` on 2026-09-13 with no credential, with a deliberately invalid one, and with and without an `Accept` header.

## What this Recipe found

**Without `Accept: application/json`, an API request is redirected into a login page.**

```
GET /v1/sites                            302 text/html
                                         Location: https://app.usefathom.com/login
GET /v1/sites  Accept: application/json  401 {"message":"Unauthenticated."}
```

Not 401, not JSON, not an error at all as far as the status line goes — a redirect to an interactive sign-in form. Every HTTP client follows a 302 by default, so what a caller actually receives is 200 and an HTML page. That page carries a `<meta http-equiv="refresh">` as well, so it redirects again if the first one is somehow not followed.

Fathom documents this, which makes it a decision rather than an accident:

> Always send `Accept: application/json` so that errors are returned as JSON rather than an HTML error page.

So the header that makes no difference to a successful response is the one that decides whether a failure is machine-readable at all.

**And the JSON failure is not the documented one either.** The error section says most failures "return a single `error` key with a human-readable message", and gives `{"error": "…"}` as its Unauthorized example. Live, the 401 is `{"message":"Unauthenticated."}` — Laravel's default sentence, in a key the documentation never mentions, with a trailing full stop. A missing token and a wrong one answer identically.

**The documentation says its own messages are not to be relied on.** Verbatim:

> Error messages are written for humans and may change over time. When handling errors programmatically, branch on the HTTP status code rather than matching on the message text.

Which is good advice, and leaves `400` covering seven unrelated conditions — invalid parameters, a failed validation, a token without permission for the action, an unsupported `entity`, an hourly `date_grouping` over a range longer than 7 days, and an account whose subscription has lapsed — with nothing but prose to tell them apart.

**Two failure shapes, chosen by the kind of failure.**

```json
{"error": "This token doesn't have permission to access this endpoint"}
{"errors": {"name": ["The name field is required."]}}
```

Singular and a string, or plural and an object of arrays. A client has to test which key arrived before it can read either.

**The cursor is not in the envelope.** The list object carries `object`, `url`, `has_more` and `data`, and the value to send as the next cursor is, in the documentation's own words, "the `id` of the first or last item in `data`".

So paging means reaching into the collection and pulling an identifier out of a record — and a page with no records has no cursor at all.

**And the direction parameter changes the sort order.** `starting_after` sorts chronologically; `ending_before` sorts in reverse. The same endpoint returns records in two different orders depending on which cursor a caller holds, and the two may not be combined.

**A site's timestamp is not a timestamp.** `created_at` is `2020-07-27 12:01:01`: a space where RFC 3339 wants a `T`, and no offset. The record carries a `timezone` field separately, so the zone the timestamp should be read in is a different key — and nothing says whether it applies to `created_at` or only to the reports.

## Sources

- Live: `api.usefathom.com`, struck 2026-09-13.
- [Fathom Analytics API v1](https://usefathom.com/api) — published as a single Markdown file for machine readers.

## Modelling limits

- **One route.** Listing sites. Account, tokens, events, milestones, aggregations, current visitors and the reporting surface each want their own evidence.
- **The 302 is recorded, not served.** Cauldron answers one shape per route, and Fathom's depends on the request's `Accept` header — a redirect to a login page without it, JSON with it. Every case here sends the header, and the redirect is quoted above because it is the finding.
- **The documented `error` shape is not modelled.** The 401 served here is the one that arrives; the `{"error": …}` and `{"errors": {…}}` shapes are described in the documentation and nothing this Recipe can do provokes either.
- **`ending_before` is not modelled.** Paging backwards reverses the sort order, and Cauldron declares one cursor parameter per route.
- **No `spec:`.** Fathom documents this API as prose with worked examples. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** A project using Fathom holds a script tag or a proxy route rather than a client of this API, and the packages named for it are tracker wrappers that post pageviews rather than read the reporting surface. Checked 2026-09-13.
