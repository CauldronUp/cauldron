# beehiiv

Emulates the beehiiv API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against beehiiv's reference at `developers.beehiiv.com` and struck live against `api.beehiiv.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**"The api key is not valid" is the answer to sending no api key.**

Struck live, both ways, byte-identical:

```
GET /v2/publications    (no credential)
GET /v2/publications    Authorization: Bearer notreal
401 {"status":401,"statusText":"unauthorized",
     "errors":[{"message":"The api key is not valid","code":"INVALID_API_KEY"}]}
```

The **third** provider in this collection to describe the wrong failure in that exact way, after [loops](../loops) and [helicone](../helicone). It is the commonest mistake there is in an error message: the sentence gets written for the case somebody was thinking about and then reused for the case they were not.

The `code` is `INVALID_API_KEY`, which repeats the mistake in the field a client branches on. Both halves of the response agree, and both are wrong about half the requests that produce them.

**The status appears three times** — as the HTTP status, again as a number in `status`, and as the reason phrase lower-cased in `statusText`. Two of the three are derivable from the first, and the lower-casing is the same choice [terraformcloud](../terraformcloud) makes with its `title`.

**The `errors` array is the good part.** Objects rather than strings, each with a `message` and a `code`, so a request that fails validation several ways gets an entry per way and a client can show all of them. That is the shape [gong](../gong) reaches for and misses by putting strings in it, and the one [honeybadger](../honeybadger) misses by putting a string where the array goes. Three providers, one idea, and beehiiv is the one that got it right.

**An unknown path is an empty HTML response.** 404, `text/html`, zero bytes — so the failure a client meets while getting a URL wrong has no body at all.

**A listing reports four numbers about itself** — `limit`, `page`, `total_results`, `total_pages` — all at the top level beside `data` rather than in a `meta` object, so a record field called `page` would collide with the envelope.

**The creation time is seconds in a field called `created`.** Ten digits, and neither the name nor the value says which unit.

## Modelling limits

- **One route.** Publications. Posts, subscriptions, segments, automations, custom fields and the whole webhook surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** beehiiv publishes a rendered reference site.
