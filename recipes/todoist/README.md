# todoist

Emulates the Todoist API for local development and tests.

**15 conformance cases, 8 checked against the live API on 2026-09-11.**

Written against Todoist's published OpenAPI document at `developer.todoist.com` and struck live against `api.todoist.com` on 2026-09-11 with no credential, with a deliberately invalid one, on a path that does not exist, and on two retired ones.

## What this Recipe found

**The retired endpoint answers 410 with prose, a wink, and the wrong version number.**

```
GET /rest/v2/projects
410 Gone   Content-Type: text/plain

This endpoint is deprecated.

If you're reading this on a browser, there's a good chance you can change
the v9 part on the URL to v1 and get away with it. ;)

If you're using the API directly, please update your use case to rely
on the new API endpoints, available under /api/v1/ prefixes.

For more details, please see documentation at https://developer.todoist.com/
```

**410 rather than 404 is exactly right** — the path was real, the endpoint is gone, and retrying will not help. Almost nothing else in this collection sends one.

And then the message tells the reader to change "the v9 part on the URL", which this URL does not contain: it was `/rest/v2/`. One handler serves every retired version with a sentence hardcoded to name one of them.

Struck live, `/rest/v9/projects` answers **404**. So the version the advice names is the version that no longer routes, and the version that *does* route is told to go looking for it.

It is also addressed to a person — "if you're reading this on a browser", and a winking emoticon — in a `text/plain` body, delivered to whatever HTTP client asked.

**The current API says the status three times, in three vocabularies.**

```json
401 {"error":"Unauthorized",
     "error_code":477,
     "error_extra":{"event_id":"0e17262e…","retry_after":3},
     "error_tag":"UNAUTHORIZED",
     "http_code":401}
```

`http_code` repeats the status line. `error_tag` is the reason phrase in screaming snake. `error_code` is **477** — which is not an HTTP status at all, and is the only one of the three carrying information the status line did not.

**And `retry_after` escalates on a failure that retrying cannot fix.** Four consecutive unauthenticated calls answered **18, 35, 67 and 128** — roughly doubling.

That is a backoff schedule, attached to a 401. No amount of waiting turns a missing credential into a present one, and a client that honours the hint sleeps two minutes before failing in exactly the same way. A 404 carries one too.

Probing stopped there. A number that climbs per attempt is the provider asking to be left alone, so this Recipe records the escalation rather than mapping its ceiling.

**The challenge header is sent only to callers who already sent something.** A missing credential gets no `WWW-Authenticate` at all. A wrong one gets:

```
WWW-Authenticate: Bearer error="invalid_token",
  error_description="Unauthorized",
  resource_metadata="https://api.todoist.com/.well-known/oauth-protected-resource"
```

RFC 6750's error parameters and RFC 9728's metadata pointer together — one of the two best challenge headers in the collection, withheld from the one caller who needs it. [cloudamqp](../cloudamqp) has the same fault with a plainer header; [fastmail](../fastmail) sends the good one to everybody.

**The published document declares no authentication whatsoever.** `components.securitySchemes` is empty and `security` is null, on an API that answers 401 to every request without a credential.

That is not a scheme defined and left unreferenced, as [close](../close) and [middesk](../middesk) do, and not a hand-rolled header parameter either. It is nothing at all.

**A project carries four ordering fields.** `child_order` and `default_order` are integers; `order_key` and `default_order_key` are fractional-indexing strings. Two orderings, each expressed twice in two different schemes, all four required by the schema.

**And `is_deleted` is required on every project in the listing**, so a list of projects contains deleted projects, and a client that does not filter shows them.

**The timestamps are required and nullable.** `created_at` and `updated_at` are `anyOf: [string, null]` and both appear in `required` — the key is guaranteed and the value is not. That is the third shape of that trade this week, after [close](../close) and [insightly](../insightly).

**Eight booleans describe one project**, two of them about what the caller may do rather than about the project: `can_assign_tasks` and `can_comment` sit beside `is_archived`, `is_deleted`, `is_favorite`, `is_shared`, `is_collapsed` and `is_frozen`.

## A note on how this Recipe was written

Todoist's description could not be read at all when this started. `cauldron drift` reported the provider unreadable, and the document turned out to be fine — the reader was re-serialising it through YAML text on its JSON recovery path, and yaml.v3 cannot parse the block scalars it writes for strings that begin with a newline. Nine of Todoist's code samples begin with a newline.

That is fixed separately, and the 357-byte reduction of this document is the regression test.

## Modelling limits

- **One listing and one gravestone.** Projects, and the retired `/rest/v2/` surface. Tasks, sections, comments, labels, filters, activity and the whole sync API each want their own evidence — the document has 76 paths.
- **`retry_after` is served as a constant.** Live it climbs per failed attempt; reproducing an escalation would mean modelling how many times a caller has already failed, and the finding is that the hint is there at all rather than what it counts to.
- **Nothing is mapped in detection.** Todoist is reached through its own apps and through generic HTTP clients holding a personal token, and no client library on npm, Packagist or the Go module proxy calls this API under an obvious name, checked 2026-09-11.
