# maxio

Emulates the Maxio Advanced Billing API — formerly Chargify — for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-07.**

Written against Maxio's reference at `developers.maxio.com` and struck live against a Chargify subdomain on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The header says JSON and the body is a sentence.**

```
GET /subscriptions.json          (no credential)
401  Content-Type: application/json
HTTP Basic: Access denied.
```

Twenty-six bytes of prose, declared `application/json`. Not a quoted string, which would at least parse — a bare sentence with a colon and a full stop in it. **Every client that does the correct thing, checking the type before parsing, parses it and throws.**

This collection now has three providers getting one header wrong in three different ways, and they are not equally bad:

| | Header says | Body is | A client trusting the header |
|---|---|---|---|
| [statsig](../statsig) | *nothing at all* | text | cannot parse, and is told so |
| [wrike](../wrike) | `text/plain` | JSON | skips a parse it could have done |
| **maxio** | `application/json` | text | **attempts a parse it cannot** |

Statsig's refusal is honest about being unparseable. Wrike's is over-cautious. Maxio's is the only one where believing the header is what breaks you.

The sentence is also Rack's, not Maxio's — `HTTP Basic: Access denied.` is what `Rack::Auth::Basic` writes when a challenge fails. So the middleware answers above the application, and Maxio's own JSON error envelope, which every documented failure uses, is never reached on this path.

**The product was renamed and the hostname was not.** Maxio bought Chargify; the documentation says Maxio Advanced Billing; the API lives at `<subdomain>.chargify.com`. The name a developer searches for and the name in the URL have not matched since 2022.

**Paths carry their format as an extension.** `/subscriptions.json`, not `/subscriptions` with an `Accept` header — Rails' `respond_to` convention on the wire, so the content type is chosen by the URL and a client cannot ask for anything else without changing it.

**Every record is wrapped under its own singular name.** A subscriptions listing is a bare array of `{"subscription": {...}}`, so the id is two levels down and a client indexing straight into the item finds nothing.

**A cancelling subscription still reads `active`.** `state: "active"` beside `cancel_at_end_of_period: true`. Two fields have to agree before a client can answer "is this person a customer", and only one of them is called `state`.

**Timestamps carry the site's own fixed offset**, `-04:00`, rather than UTC. Two Maxio sites answer different offsets for the same instant, so a client sorting on the string rather than the instant gets the order wrong.

## Detection

`chargely/chargify-sdk-php` names `chargify.com` 27 times and is mapped; the npm package `chargify` names it three times and has been a release candidate since 2021. Both are named for the old product, which is what the hostname still says.

## Modelling limits

- **Maxio's own error envelope is not served.** Every failure a public probe can reach comes from Rack, above the application. The documented `errors` array sits behind a working credential and nothing here could check it, so it is recorded rather than invented.
- **One route.** Subscriptions. Customers, products, components, invoices, allocations and the whole proration surface each want their own evidence.
- **No `spec:`.** Maxio publishes a rendered reference site.
