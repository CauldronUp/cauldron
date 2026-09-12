# plausible

Emulates the Plausible Analytics Sites API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-11.**

Written against Plausible's Sites API reference at `plausible.io/docs` and struck live against `plausible.io` on 2026-09-11 with no credential, with a deliberately invalid one, with a header carrying no scheme, on a second endpoint, and on a path that does not exist.

## What this Recipe found

**The two credential sentences are the best in this collection.**

```
(no header)      401 {"error":"Missing API key. Please use a valid Plausible
                      API key as a Bearer Token."}

Bearer notreal   401 {"error":"Invalid API key. Please make sure you're using
                      a valid API key with access to the resource you've
                      requested."}
```

The first names the failure, the product **and the scheme** — a caller who sent nothing is told what to send and how to send it.

The second names the failure and then hedges, correctly, on scope: *"with access to the resource you've requested"* covers the case where the key is real and simply not entitled. Most providers here either answer that with a bare 403 or do not distinguish it at all.

After seven providers in this collection that answer "invalid key" to a request carrying none, this is the pair to hold them against.

**And a header with no scheme is still reported as missing.** `Authorization: notreal` — present, non-empty, missing only the word `Bearer` — answers "Missing API key". So the one caller who demonstrably sent something is told it sent nothing.

That is the third time this exact shape has turned up: [wise](../wise) and [todoist](../todoist) both do it too. Three providers, one mistake — a header that fails to parse is counted as a header that was never there.

**A second endpoint gives a different sentence, and it admits it cannot tell two failures apart.**

```
GET /api/v1/stats/aggregate?site_id=x   Bearer notreal
401 {"error":"Invalid API key or site ID. Please make sure you're using a valid
              API key with access to the site you've requested."}
```

*"Invalid API key **or site ID**"* is an endpoint saying plainly that it does not know which of the two the caller got wrong — because answering accurately would tell an unauthenticated stranger whether a given site exists.

That is a real constraint, honestly stated, in a sentence, on the wire. Nothing else in this collection does that: the usual move is to pick one of the two failures and report it confidently.

**And then the 404 is the dashboard.** An unknown path under `/api/v1/` answers the Plausible web application's own 404 page — `text/html`, an apple-touch-icon, a viewport meta tag. So the three failures on this host are two carefully written JSON sentences and a web page.

**A site has no identifier except its domain.** The record is `{domain, timezone}` and nothing else.

So the primary key is the thing a customer changes when they move hosts, and a site that is renamed is — to any client holding a reference — a different site that has appeared while the old one vanished. The teams listing on the same API *does* carry a UUID, so the two listings are keyed differently.

**`meta` carries both cursors and no total.** `{"after": null, "before": null, "limit": 100}` — both directions always present, both null on a single page. A client testing for the key's presence to decide whether to keep paging never stops, and there is no count anywhere to fall back on.

**And a team can exist and not be reachable.** `/api/v1/sites/teams` returns each team with `api_available`, false on the personal one — the API listing, in its own response, which of your teams it will refuse to serve.

## Modelling limits

- **Two routes.** Sites and teams. The Stats API is a separate surface with its own key role, and its aggregate, timeseries and breakdown endpoints each want their own evidence — the "API key or site ID" sentence above is quoted from it rather than served here.
- **Nothing is mapped in detection.** A project using Plausible holds a script tag or a proxy route, not a client of this API; the npm packages named for it are tracker wrappers that post events to `/api/event` rather than reading the Sites API. Checked 2026-09-11.
- **No `spec:`.** Plausible publishes a rendered reference site.
