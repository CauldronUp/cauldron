# saucelabs

Emulates the Sauce Labs job listing for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-14.**

Read from the Jobs API reference at [`docs.saucelabs.com`](https://docs.saucelabs.com/dev/api/jobs/), and struck live on 2026-09-14 against both regions with no credential, with a wrong basic credential, on a path that does not exist, with a method the path does not take, and on the root.

## What this Recipe found

**A path that does not exist asks for a credential; a path that does exist says it does not.**

```
GET /rest/v1/cauldron-nope         401  {"detail":"Authorization failed"}
GET /rest/v1/cauldronuser/jobs     404  Not found      text/plain
PUT /rest/v1/cauldronuser/jobs     404  "Not Found"    application/json
GET /rest/v1.1/jobs                401  {"detail":"Authorization failed"}
```

The inversion is the finding. The invented path is challenged for a credential; the real endpoint shape — with a username that is not an account — answers 404 before any credential is read. The username segment is resolved ahead of authentication.

**And "not found" is spelled two ways in two formats on one path.** With no credential, `GET` answers `Not found` as `text/plain`, lowercase f. `PUT` on the same path answers `"Not Found"` — a bare JSON string, quotes included, `application/json; charset=utf-8`, capital F. A client calling `.json()` gets a throw on one and the string `"Not Found"` on the other, where `body.detail` is `undefined` rather than missing.

**The account name is a path segment.** `/rest/v1/{username}/jobs`, so the account appears in every URL, every access log and every `Referer` — and it is also half of the Basic credential.

**A field named for a boolean carries a word.** `public` is a string, and the reference's own example value is `"team"`, so `if (job.public)` is true for a job visible only to a team. `"private"` is truthy too.

**A field is a past-tense verb.** `breakpointed`, nullable, beside `container` and `performance_enabled` — three flag-shaped fields in three different grammars.

**Four identifiers on one record**: `id`, `org_id`, `team_id` and `group_id` — and `url` sits among them as `null` on the reference's own example, a field named for a link with nothing in it.

**And the format of the response is a query parameter.** `format` takes `json` or `csv` and defaults to `json`, so the media type of a successful answer is decided by a value in the query string rather than by `Accept`.

Also pinned: `creation_time` is Unix seconds and `deletion_time` is null beside it, with no unit stated on either; the offset parameter is called `skip`; and both regional hosts — `us-west-1` and `eu-central-1` — answer identically, with the region in the hostname and nothing on the record saying which one served it.

## Sources

- [Sauce Labs Jobs API](https://docs.saucelabs.com/dev/api/jobs/) — the endpoint, the response example, the job fields, and the `limit` / `skip` / `from` / `to` / `format` parameters.
- Live: `api.us-west-1.saucelabs.com` and `api.eu-central-1.saucelabs.com`, struck 2026-09-14 with no credential, a wrong Basic credential, an unrouted path, a wrong method, `/rest/v1.1/jobs`, and the root.

## Modelling limits

- **No description is published.** Sauce Labs serves no OpenAPI document for the v1 REST API, so there is no `spec` to fingerprint.
- **The plain-text `Not found` is recorded, not served.** Live, a username that is not an account answers it before the credential is read. This Recipe serves the listing for whatever username is in the path, so nothing in it reaches that answer; the four live probes above are the record of it. No real usernames were tried — the observation is about the path *shape*, and this Recipe does not test whether the response differs for an account that exists.
- **The username is not a scope.** The path carries it and the record does not, so partitioning records by it would mean adding a field the real API does not send.
- **The success side is document-derived.** Listing jobs needs a real account; the records here are the reference's own example plus a second of the same shape, and every case reading them is marked documentation-only.
- **`format=csv` is not served.** The finding is that the parameter exists; serving a second representation of a document-derived fixture would be inventing a CSV nobody has seen.
- **One route of many.** The job listing. Job detail, assets, the storage API, tunnels, the Real Device surface and the v2 team and user endpoints are the rest.
- **Nothing is mapped in detection.** Sauce Labs is reached through `saucectl` or a plain Basic call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
