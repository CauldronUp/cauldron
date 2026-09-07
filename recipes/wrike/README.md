# wrike

Emulates the Wrike API for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-09-07.**

Written against Wrike's API reference at `developers.wrike.com` and struck live against `www.wrike.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The same body arrives as `text/plain` on one path and `application/json` on another.**

Struck live, three requests:

```
GET /api/v4/tasks              (no credential)
401  Content-Type: text/plain
{"errorDescription":"Authorization method unknown","error":"not_authorized"}

GET /api/v4/tasks              Authorization: Bearer not-a-real-token
401  Content-Type: text/plain
{"errorDescription":"Access token is unknown or invalid","error":"not_authorized"}

GET /api/v4/cauldron-nope      Authorization: Bearer not-a-real-token
401  Content-Type: application/json
{"errorDescription":"Access token is unknown or invalid","error":"not_authorized"}
```

The second and third are byte-identical JSON. The declared type is not. **A route Wrike serves says `text/plain`; a route it does not says `application/json`** — so the header is right on the path that does not exist and wrong on the one that does.

The consequence is precise. Any client that checks the content type before parsing — Guzzle with a JSON middleware, Axios' default transform, Go's `encoding/json` behind a type switch — refuses to parse the refusal it will actually receive, and parses the one it gets only when the URL is wrong. The failure looks like a transport bug on the correct path and like a clean error on the incorrect one, exactly backwards.

**An unknown path is a 401, not a 404.** The credential is judged before the route, so a typo in a URL is reported as an authentication problem, and there is no way from outside to find out whether a path exists without a working token.

**The two credential failures differ, and usefully.** "Authorization method unknown" for a request with no header; "Access token is unknown or invalid" for one carrying a token nobody issued. Both under `error: "not_authorized"`, so the machine-readable field is identical and the prose is what separates them — the ordinary trade, and here the sentences are at least accurate about which mistake was made. ("Method" means the authentication method, not the HTTP verb, which is a word that sentence picks badly.)

**The prose field is `errorDescription`.** Not `message`, not `detail`, not `error_description`. One camel-case field beside a lower-case `error` carrying the code.

**A task's status is not the status on the board.** `status` is one of four fixed values; the account's own workflow lives in `customStatusId`, which needs a second request to resolve. Two tasks both reading `Completed` here can sit in different columns.

**The permalink is not built from the identifier.** `id` is `IEACJ7L6KQAAAAAA` and the link a person opens carries a decimal number that appears nowhere else on the record — so a client cannot construct one from the other in either direction.

**Paging is opt-in.** A `nextPageToken` comes back only once `pageSize` has been sent. A caller who never sends one is served the whole collection and never sees a cursor, so a complete answer and a first page look the same.

## Modelling limits

- **One route.** Tasks. Folders, projects, custom fields, timelogs, comments, dependencies and the whole webhook surface each want their own evidence.
- **The success path is documentation-derived.** Reading anything needs an account; the three refusals come from the wire.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Wrike publishes a rendered reference site.
