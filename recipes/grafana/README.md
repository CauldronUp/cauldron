# Grafana

Emulates the Grafana API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A dashboard save returns two identifiers, and the one called `id` is the one you can't use. Both are required in the response, but there is exactly one path that fetches a dashboard back -- `/api/dashboards/uid/{uid}` -- so the integer `id` is a row number specific to one Grafana instance and meaningless in the next. A deploy that stores `id` instead of `uid` has stored the identifier that won't survive a migration; the document even deprecates the equivalent field on folders (`folderId`) in favor of `folderUid`, in as many words.

A field called `title` in the save response actually carries the slug -- documented, verbatim, as "The slug of the dashboard" with example `"my-dashboard"` -- so the obvious way to show a user what they just saved shows them a URL fragment instead. And the version number needed to save a dashboard again isn't on the dashboard itself: a fetch returns `{dashboard, meta}`, the dashboard body is free-form JSON, and `version` lives in `meta`; saving with a stale one is a 412, not a 409 or 422.

## Sources

- Documentation: https://grafana.com/docs/grafana/latest/developers/http_api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve grafana     # run it
cauldron verify grafana -v # check every claim
```
