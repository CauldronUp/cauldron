# National Weather Service

Emulates the National Weather Service API (nationalweatherservice), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

Weather Service's, where the path says lat,lon and the
body says lon,lat. `GET /points/39.7456,-97.0892` answers
`"coordinates": [-97.0892, 39.7456]` -- the same pair, in the same exchange, in
the opposite order, because the path takes latitude first the way a person
writes a coordinate and the body puts longitude first the way GeoJSON requires.
Nothing in the response says the order changed, and near the equator or the
prime meridian reading `coordinates[0]` as what you sent gives a plausible place
rather than an error.

**And `type` is in that document three times, in three vocabularies.** The top
level says `Feature`, which is GeoJSON's; `properties.@type` says `wx:Point`,
which is the NWS ontology's; `properties.type` says `land`, which is the NWS's
own classification -- and only the middle one is namespaced. A number with a
unit is expressed four ways in one forecast: `elevation` is a `unitCode` object,
`temperature` is a bare number beside a separate `temperatureUnit`, `windSpeed`
is the prose `"10 to 15 mph"` and sometimes the scalar `"15 mph"`, and
`probabilityOfPrecipitation` is a `unitCode` object for a percentage.
`validTimes` is an interval whose right half is a duration --
`2026-08-28T16:00:00+00:00/P7DT9H` -- so `Date.parse` returns `NaN`. A request
with no `User-Agent` is refused by the CDN rather than the API, 403 in HTML,
with no `correlationId` to quote to anybody. And three failures are 404 and mean
three different things: the one whose title is accurate is the one carrying the
least, while a latitude of 999 says only "Not Found" and hides the reason in a
`parameterErrors` array RFC 9457 does not define -- an array whose message, for
an unknown forecast office, is 840 characters long because it inlines all 133
valid office codes into an English sentence.

## Sources

- Documentation: https://www.weather.gov/documentation/services-web-api
- Machine-readable description: https://api.weather.gov/openapi.json, last checked 2026-08-31
  `cauldron drift nationalweatherservice` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nationalweatherservice     # run it
cauldron verify nationalweatherservice -v # check every claim
```
