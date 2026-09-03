# NHTSA vPIC

Emulates the NHTSA vPIC API (vpic), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

VPIC, where four errors arrive in one string and one of them
is 400 on a 200. `ErrorCode` is `"6,7,11,400"` -- a comma-separated list of
numbers inside a string, four codes for one VIN, and `400` among them looks like
an HTTP status and is not, because the response is a 200 and always is.
`ErrorText` is those four joined with `"; "`, and one of them contains its own
semicolon, so splitting on the separator gives five pieces for four errors.
**And the failure is a success with 148 empty strings in it**: `Count` is 1, the
`Results` array holds one object, and 148 of that object's 154 fields are `""` --
nothing absent, nothing null, so a VIN the service could not read looks exactly
like one it read and found nothing about.

Every value is a string: `ModelYear` is `"2003"`, `EngineCylinders` is `"6"`, and
`DisplacementL` is `"2.998832712"`, a three-litre engine to nine decimal places
as text. `Message` is a 250-character disclaimer on every successful response,
explaining that a missing value does not mean the feature is absent.
`SearchCriteria` is the request as prose -- `"VIN(s): 1HGCM82633A004352"`. The
error fields are populated on success, where `ErrorCode` is the string `"0"`. And
`SuggestedVIN` on a bad VIN is `"N!TAV!N"`: the input with its invalid characters
replaced by exclamation marks, which is not a VIN and cannot be submitted
anywhere. The one failure status on the whole API is for a path that does not
exist, and its message quotes the URI as **the backend saw it** --
`backend-vpic-api.nhtsa.dot.gov`, a host the caller never named -- under a
`message` key spelled in lower case where every success spells it `Message`.

## Sources

- Documentation: https://vpic.nhtsa.dot.gov/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nhtsa     # run it
cauldron verify nhtsa -v # check every claim
```
