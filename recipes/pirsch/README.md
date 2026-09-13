# pirsch

Emulates the Pirsch Analytics domains API for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-13.**

Struck live against `api.pirsch.io` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist. Record shapes come from Pirsch's own Go SDK, which is open.

## What this Recipe found

**Three failures, three shapes, and the first one is a success.**

```
(no header)      200 application/json           null
Bearer notreal   401 text/plain; charset=utf-8  {"validation":null,"error":["token contains an invalid number of segments"]}
/cauldron-nope   404 text/plain; charset=utf-8  404 page not found
```

**A request with no credential at all answers 200 and the four bytes `null`.** Not an empty array — the JSON null literal, which is what a Go handler returning a nil slice sends.

`response.ok` is true. `.json()` succeeds and hands back null. The `for (const d of body)` on the next line throws, and a client that reaches `body.length` instead gets `undefined` and concludes the account has no domains.

**The 401 is JSON labelled `text/plain`.** A client that checks the content type before parsing refuses to read it; one that trusts the label gets a string. The body is JSON either way.

**And its message is a JWT library's internal error text.** `token contains an invalid number of segments` is golang-jwt's own string, passed through verbatim — so the failure tells a caller who has proved nothing that the credential is a JWT, and which stage of parsing it failed at.

**`validation` is `null` here and `{}` elsewhere.** The same envelope on a different path answers `{"validation":{},"error":["Domain not found."]}`. One field, two ways of being empty, and a client testing `if (body.validation)` gets a different answer depending on which failure arrived. `error` is an array while `validation` is an object, so reading a failure means indexing one field and keying the other.

**An unrouted path is not JSON at all.** `404 page not found` is Go's own `http.NotFound` string, in plain text, from an API whose other failures are JSON. A client has three parsers to write and only one of the three responses tells it which to use.

**A domain record carries a third party's account.** `google_user_id`, `google_user_email` and `gsc_domain` sit on the domain beside the site's own settings — an email address on a resource that is otherwise about a hostname.

**The timestamps are called `def_time` and `mod_time`.** Every record carries them, on a `BaseEntity` the whole API embeds, where `created_at` and `updated_at` would say the same thing in words a reader already knows.

**And two fields name a person without agreeing on a type.** `user_id` is a string and `new_owner` is a nullable integer, on the same record.

**Billing state travels with the website.** `subscription_active` sits on the domain, so the record describing a hostname also says whether it is paid for.

## Sources

- Live: `api.pirsch.io`, struck 2026-09-13.
- [`pkg/types.go`](https://github.com/pirsch-analytics/pirsch-go-sdk/blob/master/pkg/types.go) — `Domain` and `BaseEntity`, in Pirsch's own Go SDK.

## What it added to Cauldron

`responses.list.null_when_empty` — send the JSON `null` literal in place of an empty collection, which is what a Go handler returning a nil slice sends. [hatchet](../hatchet) answers `null` for a failure and Pirsch answers it for an empty success; the second occurrence is what named the field. It sits beside `omit_when_empty` and is not the same thing: a missing key reads as `undefined` and a null reads as `null`, and the two break different code.

## Modelling limits

- **One route.** Listing domains. Statistics, sessions, events, funnels, keywords, the live-visitor endpoints and the whole reporting surface each want their own evidence.
- **The wrong-credential body is served as JSON under a `text/plain` label**, which is what Pirsch does. A client written against this sandbox meets the same mislabel it will meet in production.
- **`{"validation":{}}` is not modelled.** The empty-object form appears on a different path, quoted above; this Recipe serves the `null` form its own route answers.
- **No pagination.** The endpoint answers every domain the credential can see, and nothing in the SDK or the documentation names a page parameter.
- **No `spec:`.** Pirsch documents this API as prose and ships SDKs. There is no OpenAPI document at any address this Recipe could find; the record shape comes from the Go SDK, which is the artefact a description would have produced.
- **Nothing is mapped in detection.** A project using Pirsch holds a script tag or one of its tracking SDKs, which post pageviews to a different surface — not a client of the dashboard API this Recipe describes. Checked 2026-09-13.
