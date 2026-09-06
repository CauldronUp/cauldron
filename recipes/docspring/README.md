# docspring

Emulates the DocSpring PDF generation API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-06.**

Written against DocSpring's API guide at `docspring.com/docs` and struck live against `api.docspring.com` on 2026-09-06 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Every failure echoes something the caller sent back at them.** Struck live:

```
no credential
401 {"status":"error","error":"Missing Basic Auth: Please provide an API token
     via Basic Auth. See: https://docspring.com/docs/api-guide/authentication/"}

Basic not-a-token:not-a-secret
401 {"status":"error","error":"Could not find API token with ID: not-a-token"}

GET /api/v1/cauldron-nonexistent
404 {"status":"error","error":"Unrecognized request URL (GET:
     /api/v1/cauldron-nonexistent). Please see https://docspring.com/docs or
     contact support if you need any help."}
```

The second quotes the token id back. The third quotes the method and the path. Both are useful when the caller is you, and here neither costs anything — because **the thing quoted is the token id, which is the public half of a Basic pair**. The secret is the password and it is never repeated.

That distinction is the point. Several providers in this collection echo a credential into an error and echo the wrong half. DocSpring echoes the half that identifies and not the half that authorises, which is the version of this behaviour that is safe to have.

**`status` is the string `"error"`.** Not the HTTP status, not a boolean. The literal word, on all three failures. So the only structured field in the envelope has one value and carries no information, and everything a client needs is inside `error` — which is prose, with URLs in it.

**There is no code, no type and no field name.** Three different failures — no credential, a bad credential, a bad path — and nothing in any of them a program can switch on except by matching English.

**The 401 says how long it took.** `x-runtime: 0.015572` rides the refusal: the server's own processing time, in seconds, as a float, disclosed to an unauthenticated caller. `x-request-id` is 32 hex characters beside it. Rails, showing through.

**Expiry is a number and its unit is a separate string.** `expire_after: 1` with `expiration_interval: "hours"`, beside `expire_after: 7` with `"days"`. Comparing `expire_after` across two templates without reading the other field compares hours to days, and the field that would warn you is the one with the less obvious name.

**A template still processing reports zero pages and no error.** `page_count: 0` is a state a client polls through, not a failure — so code that treats the count as authoritative renders an empty document rather than waiting.

## Detection

`docspring` names `api.docspring.com` five times in its published archive, and `docspring/docspring` is the same client on Packagist. Both are mapped.

## Modelling limits

- **One route.** Templates. Submissions, PDF generation, data requests, folders and the whole async generation lifecycle each want their own evidence.
- **The success path is documentation-derived.** Listing templates needs an account; the three failures come from the wire.
- **No `spec:`.** DocSpring publishes a rendered API guide.
