# FHRS

Emulates the FHRS API (fhrs), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

The Food Hygiene Rating Scheme, where **the same digit means best
in one field and not in the other.** `"RatingValue": "5"` is the top of a scale
that counts up, and `"scores": {"ConfidenceInManagement": 0}` is the top of one
that counts points against you -- so zero is perfect there and five is not, in
the same object, and the first is a string while the second is a number.

**And absence is spelled three ways in one record:** `Phone` is `""`, `Distance`
is `null`, and `AddressLine2` is `""` while `AddressLine3` holds the town, so an
empty line is a gap in the middle of an address rather than the end of it. A
forgotten `x-api-version` header is reported as a path that does not exist --
`"The API 'Authorities' doesn't exist"` -- on a path that works the moment the
header is sent. That failure is a bare JSON string rather than an object, so
`JSON.parse` succeeds and `body.Message` is `undefined` rather than throwing,
and the name inside it is the internal handler: the version and the path joined
by a dot. The coordinates are strings where the scores are numbers, and
`RatingKey` packs the scheme, the rating and the locale into `"fhrs_5_en-gb"`.

## Sources

- Documentation: https://api.ratings.food.gov.uk/help
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fhrs     # run it
cauldron verify fhrs -v # check every claim
```
