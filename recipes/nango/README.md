# nango

Emulates the Nango API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Nango's reference at `docs.nango.dev` and struck live against `api.nango.dev` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**An unknown API path answers 200 with an HTML page, and the page carries an engineering note about token leakage.**

```
GET /cauldron-nope
200  Content-Type: text/html
```

3043 bytes of application shell, with a **200** on it. Not a 404 — a success — so a client that checks the status before the body concludes its request worked, and one that calls `.json()` throws on a page HTTP told it was fine.

Inside that page, as an HTML comment:

> Auth urls carry tokens in the path. Cloud serves this from a CDN, so the server's Referrer-Policy header doesn't apply. Must stay above any request-initiating tag.

It is a note from Nango's own engineers explaining a real mitigation: their OAuth callback URLs carry tokens as path segments, a path segment lands in the `Referer` header of any subresource the page loads, the CDN in front of the app does not pass the server's `Referrer-Policy` through, and so a `<meta name="referrer">` tag has to come before anything that makes a request.

**The mitigation is correct and carefully reasoned.** The comment is also a precise description of where the soft spot is, served with a 200 to anybody who mistypes a path — and the same reasoning tells a reader that a tag which sneaks above it in any future edit undoes it silently.

**The refusal says what shape the credential should have.**

```
(no header)      401 {"error":{"message":"Authentication failed. The request is missing the Authorization header.","code":"missing_auth_header","payload":{}}}
Bearer notreal   401 {"error":{"message":"Authentication failed. The provided secret key is not a UUID v4.","code":"invalid_secret_key_format","payload":{}}}
```

Two codes, two accurate sentences, and the second names the format — **UUID v4**, specifically, not just "a UUID". That is the most precise credential-shape message in this collection: a caller who has pasted a public key, an old key or a truncated one is told which check failed.

It is also, notably, **accurate about the missing case** — where [loops](../loops), [helicone](../helicone) and [beehiiv](../beehiiv) all answer "invalid key" to a request that carried none.

**`payload` is always an empty object** — a third field on every failure, reserved for structured detail, carrying none. The same shape [confluent](../confluent)'s `source: {}` has.

**The path is singular and the envelope is plural.** `GET /connection` answers `{"connections":[…]}`, so the route and the response disagree about how many things this is.

**The identifier is the customer's own user reference.** Nango mints none of its own, so `connection_id` is whatever the integrating application called the user — unvalidated, and URL-safe by luck.

## Modelling limits

- **One route.** Connections. Syncs, actions, the proxy, webhooks, integrations and the whole records surface each want their own evidence.
- **The HTML page is served in outline rather than byte for byte.** What the case pins is the status, the content type and the comment; reproducing a build artifact exactly would be asserting a bundler's output.
- **Nothing is mapped in detection.** Nango is used through its own dashboard and a proxy base URL more often than through a client library, and nothing on any registry calls `api.nango.dev` under an obvious name.
- **No `spec:`.** Nango publishes a rendered reference site.
