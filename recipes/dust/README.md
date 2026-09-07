# dust

Emulates the Dust API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Dust's reference at `docs.dust.tt` and struck live against `dust.tt` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Two typed failures, and one of the types ends in `_error`.**

```
(no credential)   401 {"error":{"type":"not_authenticated",
                                "message":"The request does not have valid authentication credentials."}}
Bearer notreal    401 {"error":{"type":"malformed_authorization_header_error",
                                "message":"Malformed Authorization header"}}
```

Two distinct types for two distinct failures, which is more than most providers in this collection manage.

And the second is `malformed_authorization_header_error` — a type name, inside a field called `type`, inside an object called `error`, ending in the word *error*. The suffix is redundant twice over. `not_authenticated` beside it has no suffix, so **one vocabulary is built to two rules**.

**The first message is gRPC's.** "The request does not have valid authentication credentials." is the standard text for `UNAUTHENTICATED` in Google's status code definitions, which is where a great many services copied it from. It is **accurate here**, which is worth noting after five providers in this collection that describe the wrong failure.

**The second is accurate too, and narrower.** "Malformed Authorization header" — not "invalid key". A caller who sent a token with the wrong prefix, no prefix, or a stray newline knows to look at the header rather than at the key inside it.

**An unknown path is `text/plain`** — `404 Not Found`, four words, from Next.js' own handler. So the JSON above belongs to routes that exist, and a client cannot parse the failure it meets while getting a URL wrong.

**A new conversation has no title.** A conversation exists before its first message is summarised, so `title` is `null` on anything recent and a list view has to fall back to something — usually the first message, which is a second fetch.

**`visibility: "unlisted"` is the more private of the two values**, which is the opposite of what the word suggests to anyone who has used a video site.

## Modelling limits

- **One route.** Conversations. Messages, agents, data sources, spaces and the whole run surface each want their own evidence.
- **Nothing is mapped in detection.** Dust's SDK is named for the company and covers the agent runtime as well as this API.
- **No `spec:`.** Dust publishes a rendered reference site.
