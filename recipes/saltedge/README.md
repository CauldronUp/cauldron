# saltedge

Emulates the Salt Edge account-information API for local development and tests.

**10 conformance cases, 6 checked against the live API on 2026-09-07.**

Written against Salt Edge's reference at `docs.saltedge.com` and struck live against `www.saltedge.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The documentation link points at the anchor for that exact error.**

```
GET /api/v5/customers        (no credential)
400 {"error":{"class":"AppIdNotProvided",
              "message":"'App-Id' was not provided in request headers.",
              "documentation_url":"https://docs.saltedge.com/account_information/v5/#errors-app_id_not_provided",
              "request_id":"…"}}
```

Not the API's front page. Not a generic errors page. The fragment for `app_id_not_provided` — the failure that just happened.

**This is the best `documentation_url` in the collection, and it is worth saying so.** Several providers here send one — GitHub's, [statsig](../statsig)'s, [docspring](../docspring)'s — and every one points at a page a reader still has to search. Salt Edge deep-links to the paragraph. Somebody reading a failure in a log can click once and be looking at the explanation of that failure.

**`class` is the exception class.** `AppIdNotProvided`, `ApiKeyNotFound` — Ruby class names, PascalCase, in the field a client branches on. A stable machine-readable code, which is more than most providers here manage, and also a name from the codebase: renaming the class renames the wire contract.

**A credential failure is 400.** Not 401, not 403. So the status a client retries on, the status it re-authenticates on, and the status it reads as "my payload is wrong" are one number — and the third reading is the one that sends somebody looking at their request body.

**The two halves of the credential fail differently.** No `App-Id` is `AppIdNotProvided`; a wrong `App-Id`/`Secret` pair is `ApiKeyNotFound`. Rarer here than it should be, and the reason the two failures have separate documentation anchors.

**An unknown path answers the credential failure**, so routing is judged after authentication and a path cannot be shown not to exist.

**The customer secret is in the listing.** It is Salt Edge's own per-customer token rather than a bank credential — but it is still a secret arriving in a collection response that nobody has to ask for.

**A blocked customer is marked by a date rather than a flag.** `blocked_at` present versus absent, so the boolean a client wants has to be derived.

**The paging block is nulls rather than absences.** `next_id` and `next_page` are both present and both `null` on the last page, so testing for the key finds it and the value has to be tested instead.

## Modelling limits

- **One route.** Customers. Connections, accounts, transactions, consents, the refresh lifecycle and the whole callback surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Salt Edge publishes a rendered reference site.
