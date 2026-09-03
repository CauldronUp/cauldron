# PyPI

Emulates the PyPI API (json), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-08-25.**

## What this Recipe found

It publishes no OpenAPI for its JSON API and needs no
credential to answer, so the whole Recipe was written from live responses and
every case carries the date it was checked. The headline is three download
counters that always say minus one -- `"downloads": {"last_day": -1,
"last_month": -1, "last_week": -1}` -- kept after Warehouse stopped serving
download statistics through this API, with a fourth on every file. The rest
were settled the same way: `releases` is on `/pypi/{project}/json` and absent
from `/pypi/{project}/{version}/json`, so the key a client walks to enumerate
versions is missing from the more specific endpoint; three of requests'
releases map to empty arrays, so counting them counts versions nobody can
install; `core-metadata` is the one hyphenated key among snake-cased
neighbours; `md5_digest` and `digests.md5` carry the same hash and
`upload_time` and `upload_time_iso_8601` the same instant at two precisions;
and 1.26.20 is a real release of `urllib3` that answers 404 under `requests`,
which is what makes the project segment load-bearing rather than decorative.

## Sources

- Documentation: https://docs.pypi.org/api/json/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pypi     # run it
cauldron verify pypi -v # check every claim
```
