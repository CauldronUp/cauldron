# metagraph

Emulates the Meta Graph API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Meta's reference at `developers.facebook.com` and struck live against `graph.facebook.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Every credential failure is HTTP 400.**

```
GET /v21.0/me                       (no token)
400 {"error":{"message":"An active access token must be used to query information about the current user.",
              "type":"OAuthException","code":2500,"fbtrace_id":"Ax_aO2TO0AuY7zAL3_XwPPv"}}

GET /v21.0/me?access_token=notreal
400 {"error":{"message":"Invalid OAuth access token - Cannot parse access token",
              "type":"OAuthException","code":190,"fbtrace_id":"A3Q2jXILH9XCUMt5561RGqd"}}
```

Not 401. Meta's whole Graph API answers 400 to authentication problems — so **the status a client re-authenticates on never arrives**, and the status it reads as "my request is malformed" is the one that means "your token expired". Every Facebook integration in existence has a branch treating 400 as an auth failure, and that branch also catches every genuine validation error.

**The codes are small integers with a large gap.** `2500` for a missing token, `190` for an unparseable one — Meta's own numbering, four digits and three, unrelated to HTTP. `190` is the one every developer knows by heart, because it is what an expired user token produces.

**`type` is `OAuthException` for both.** A class name, the same choice [saltedge](../saltedge) makes with `class` — and here it is a constant across every credential failure, so the field that looks like a category carries nothing and `code` is where the difference lives.

**An unknown path is answered with the token error**, not a 404 — routing is judged after the token, so a path cannot be shown not to exist and a typo is reported as a credential problem.

**`fbtrace_id` is 23 characters of base64url**, in the body, on every failure. One correlation id, no header twin, and it is the thing Meta's support actually asks for — so unlike most of the trace fields in this collection, it is doing its job.

**A listing of pages hands back a token per page.** `GET /me/accounts` returns each page's own access token beside the page it opens, so **the response a client logs contains as many credentials as the user has pages**.

**Permissions are a list of verbs rather than a role.** `["ANALYZE","ADVERTISE","MODERATE","CREATE_CONTENT","MANAGE"]` — so "can this person post" is an array membership test, and two pages with different sets have no role to compare.

## Modelling limits

- **One route.** The page listing. Posts, comments, insights, webhooks, the whole Instagram and WhatsApp surfaces and the Marketing API each want their own evidence.
- **The token is served in the query string**, which is what Meta's documentation leads with — so it lands in access logs and browser history. Serving it in a header would hide that.
- **Nothing is mapped in detection.** The Facebook SDKs are named for the platform rather than the API, and a project depending on one may be using Login, Marketing or Graph.
- **No `spec:`.** Meta publishes a rendered reference site and a per-version changelog rather than a machine-readable description.
