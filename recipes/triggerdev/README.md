# triggerdev

Emulates the Trigger.dev management API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Trigger.dev's reference at `trigger.dev/docs` and struck live against `api.trigger.dev` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**One endpoint, two error models, and which you get depends on which mistake you made.**

Struck live:

```
GET /api/v1/runs                (no Authorization header)
401  Content-Type: application/problem+json
{"title": "Unauthorized",
 "status": 401,
 "type": "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401",
 "detail": "No authorization header provided",
 "error": "No authorization header provided"}

GET /api/v1/runs                Authorization: Bearer tr_dev_notreal
401  Content-Type: application/json
{"error":"Invalid API key"}
```

The first is RFC 9457, properly typed and properly labelled. The second is a single field under a different media type. Same route, same status, two unrelated shapes — so **a client cannot write one parser for a 401 from this endpoint**. It has to branch on the content type, on a failure, which is the code path least likely to have been exercised.

**The problem `type` points at MDN.** RFC 9457 says the type URI "identifies the problem type" and that dereferencing it should yield documentation *of that problem*. This one yields Mozilla's page on the HTTP 401 status code — a correct, useful page about a status the response already carries in two other fields, and nothing whatsoever about Trigger.dev. So the field that exists to distinguish one provider's problems from another's is filled with something identical across every provider that reaches for it.

**`detail` and `error` carry the same string.** RFC 9457 defines `title`, `status`, `type`, `detail` and `instance`; `error` is not among them. It is the *other* failure's only field, welded onto this one so that `body.error` works against both shapes. That is a thoughtful thing to do, and it also makes the document five fields carrying four facts.

**An unknown path is HTML.** 404, `text/html`, a full application shell with a stale-asset recovery script in it — `api.trigger.dev` and the web app are one deployment. So the 404 a client meets while getting a URL wrong is the one `.json()` throws on, and what it throws on is a React page.

**Paging uses bracketed parameter names.** `page[size]` and `page[after]`, and the cursor comes back nested under a `pagination` object rather than at the top level. A client building a query string has to encode brackets that URL builders escape inconsistently.

**Test runs share the listing with production ones**, separated by an `isTest` boolean with no default filter. A dashboard counting runs counts both.

## Detection

`@trigger.dev/sdk` names `api.trigger.dev` 42 times in its published archive and is mapped.

## Modelling limits

- **One route.** Runs. Tasks, schedules, batches, waitpoints, environment variables and the whole realtime surface each want their own evidence.
- **The 404 is not served as HTML.** What the cases pin is the two 401 shapes; reproducing a React shell byte for byte would be asserting a build artifact rather than a contract.
- **No `spec:`.** Trigger.dev publishes a rendered documentation site.
