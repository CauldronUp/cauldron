# MBTA

Emulates the MBTA API (v3), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**`Type` is a number and a string in the same
object.** `data.attributes.type` is `1` -- the GTFS route type, where 1 means
subway -- and `data.type` is `"route"`, the JSON:API resource type. Two levels
apart, one key name, and a client reading `type` without saying which one gets
whichever its path happened to reach. A third lives in every relationship, where
`relationships.agency.data.type` is `"agency"`.

**And the human-readable text moves keys between failures.** A resource that
does not exist answers `{"code": "not_found", "status": "404", "title":
"Resource Not Found"}` -- title, no detail -- and a filter that does not exist
answers `{"code": "bad_request", "status": "400", "detail": "Unsupported
filter(s): nosuchfilter"}` -- detail, no title. So a client rendering
`errors[0].detail` shows nothing when a route is missing. `status` is a string
in the body beside a number on the wire; the same error body arrives under
`application/vnd.api+json` from one path and `application/json` from another;
colours are six hex digits with no hash, so `style="color: DA291C"` renders
nothing; and `direction_names` and `direction_destinations` are two arrays
joined only by index.

## Sources

- Documentation: https://api-v3.mbta.com/docs/swagger/index.html
- Machine-readable description: https://api-v3.mbta.com/docs/swagger/swagger.json, last checked 2026-08-31
  `cauldron drift mbta` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mbta     # run it
cauldron verify mbta -v # check every claim
```
