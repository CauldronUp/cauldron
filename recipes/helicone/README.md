# helicone

Emulates the Helicone API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Helicone's reference at `docs.helicone.ai` and struck live against `api.helicone.ai` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**`trace` is the name of the function that rejected you.**

```
POST /v1/request/query    (no credential)
401 {"error":"No authorization header","trace":"isAuthenticated.error"}

POST /v1/request/query    Authorization: Bearer sk-helicone-notreal
401 {"error":"No API key found","trace":"isAuthenticated.error"}

POST /v1/cauldron-nope    Authorization: Bearer sk-helicone-notreal
401 {"error":"No API key found","trace":"isAuthenticated.error"}
```

`trace` is a constant. Not a request id, not a span id, not a timestamp — the string `isAuthenticated.error`, identical on every failure, which is a **source-code identifier**: a function called `isAuthenticated` and the branch inside it. So the field named for correlating one request against a server log correlates nothing, and what it does disclose is a name from the codebase.

It is also the most useful thing in the response for whoever is debugging, which is the tension worth naming. Helicone is open source: a caller who greps the repository for `isAuthenticated` finds the check that refused them. That is arguably a deliberate affordance rather than a leak. It is still not a trace, and a client storing it against a support ticket has stored a constant.

**"No API key found" is the answer to sending an API key.** The second request carried a bearer token; the sentence says none was found. So the two failures are separated only by their prose, and the prose describes the wrong one — the same shape [loops](../loops) has, except that here the field which could have carried the distinction is occupied by a constant.

**An unknown path answers 401, not 404.** Routing happens after authentication, so a caller cannot discover whether an endpoint exists — and the sentence they get is about a key they did send.

**A read is a POST because the filter is a document.** `/v1/request/query` takes a filter tree in the body and answers a listing, and **the paging controls are in the body too**: `limit` and `offset` beside the filter. So a page cannot be bookmarked, does not appear in an access log, and a caching proxy sees one URL for every page of every query.

**Field names mix two conventions on one record.** `request_created_at` is snake case; `costUSD` is camel case with an acronym. A generated struct needs two tag styles for one object.

**Cost is a JSON number, not a decimal string.** `0.00021705` — so a client summing a month of these accumulates binary floating-point error on a currency value, and the field is named for the currency rather than typed for it.

**A failed upstream call still cost money.** A 500 from the model provider, no completion tokens, and 1204 prompt tokens already spent. A dashboard filtering to successful requests hides exactly the spend it most wants to explain.

## Modelling limits

- **One route.** The request query. Sessions, properties, prompts, datasets, evaluators and the whole caching surface each want their own evidence.
- **Nothing is mapped in detection.** Helicone is used as a proxy base URL rather than through a client library, so a project that depends on it usually depends on the OpenAI or Anthropic SDK with a changed host — which no dependency name can see.
- **No `spec:`.** Helicone publishes a rendered documentation site.
