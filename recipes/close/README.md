# close

Emulates the Close CRM API for local development and tests.

**14 conformance cases, 5 checked against the live API on 2026-09-09.**

Written against Close's published OpenAPI document at `developer.close.com` and struck live against `api.close.com` on 2026-09-09 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**Three rate-limit vocabularies on one response, and they disagree.**

```
ratelimit:           limit=100; remaining=100; reset=1
ratelimit-limit:     60
ratelimit-remaining: 59
ratelimit-reset:     1
x-ratelimit-limit:   60, 60;w=1, 100;w=1
```

Three conventions at once. `ratelimit` is the IETF draft's single structured field. `ratelimit-limit` and its siblings are the same draft's earlier split-header form. `x-ratelimit-limit` is the legacy `X-` name carrying the draft's quota-policy syntax inside it.

And the numbers do not agree. The structured field says the limit is **100**. The split header beside it says **60**. The legacy one says both, in a comma-separated list.

So a client reading any single header gets a different answer about the same request — and the two a reasonable client is most likely to reach for, `ratelimit-limit` and `x-ratelimit-limit`, are the two that disagree with the field the draft actually standardises. A client parsing `x-ratelimit-limit` as an integer gets 60, correctly, by luck: `parseInt` stops at the comma.

**They are sent to a caller holding no credential**, so the quota is public, and a request that has already been refused still spends some of it and says so.

**The document declares two security schemes and applies neither.** `components.securitySchemes` defines `ApiKeyAuth` — HTTP Basic, "use your API key as the username and leave the password empty" — and `OAuth2`. Top-level `security` is absent. No operation declares one either.

What every operation *does* declare is a plain header parameter named `Authorization`, `required: true`, `type: string`.

So the machinery OpenAPI has for describing authentication is present, populated, and unreferenced, and the real requirement is expressed as a string the caller assembles itself. A generator reading this emits a client whose every method takes an `authorization` string and does no auth handling at all.

It is the opposite mistake to [adobesign](../adobesign), whose document marks `Authorization` `required: false` on an API that refuses every request without one. Close over-declares the header and under-declares the scheme; Adobe declares the scheme and lets the header be optional. Between them they cover both ways of getting one field wrong.

**`name` is required and nullable.** So are `url`, `description`, `created_by` and `updated_by` — all five are listed in `required` and typed `["string", "null"]`.

A lead is the record a CRM exists to hold, and the schema guarantees that a lead carries a `name` key while permitting that key to be null. That is the distinction between *present* and *has a value*, written into the contract. Code that checks `if ("name" in lead)` is satisfied by a lead with no name.

**And two fields have no type at all.** `primary_email` and `primary_phone` are declared with a description and no `type`, so the document says those fields exist and says nothing whatever about what they hold.

**The three declared failures are typed "Any type".** 400, 401 and 404 each carry a description and a schema whose entire content is `{"description": "Any type"}`.

Three providers in this collection now sit on that spectrum, and they cover it end to end: [attentive](../attentive) declares five failure statuses and gives none of them a body at all, Close declares three and types them as anything, and Middesk declares no failure whatsoever. Every one of them answers a real, specific JSON body on the wire.

**The 404 is Flask's, and it is addressed to a person.**

```json
{"error": "The requested URL was not found on the server. If you entered the URL manually please check your spelling and try again."}
```

The framework's default 404 text, wrapped in the API's own `error` field. It tells a machine to check its spelling. It is reached with a wrong credential in hand, so routing runs before the credential check and a typo is reported as a typo.

**The status is an identifier and a label side by side.** `status_id` is stable and `status_label` is the organisation's own wording, which anyone can rename in the pipeline settings — so a client switching on the label breaks the day somebody edits a dropdown, and the field that would not break is the opaque one.

**Paging is counting.** `_limit` and `_skip`, leading underscores to distinguish paging from the filter parameters sharing the query string, with `has_more` and no cursor anywhere. A caller keeps its own offset, and a lead created mid-scan shifts every page after it.

## Modelling limits

- **One route.** Listing leads. Contacts, opportunities, activities, tasks, custom fields, the whole reporting and bulk-action surface each want their own evidence.
- **The three rate-limit headers are served as constants.** They are what a real refusal carried on 2026-09-09; nothing here decrements them, because the finding is that the three disagree with each other rather than how any one of them counts down.
- **Nothing is mapped in detection.** No client for this API on npm, Packagist or the Go module proxy under an obvious name, checked 2026-09-09.
