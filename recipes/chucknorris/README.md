# Chuck Norris

Emulates the Chuck Norris API (chucknorris), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Norris's, where the success and the failure carry timestamps in
two different grammars. A joke's `created_at` is `"2020-01-05 13:42:19.324003"`
-- a space instead of a T, six digits of fraction, no timezone -- and an error's
`timestamp` is `"2026-08-30T06:21:50.575Z"`, proper ISO 8601. Same API, same
request, and which format arrives depends on whether it worked, so a client that
logs when something happened gets a date on its errors and `Invalid Date` on its
results.

**And the two identifier alphabets do not match either.** Most jokes carry
base64url -- `"xMUTjH93TyWHGR6MVf_vpA"`, mixed case with hyphens and underscores
-- and some carry `"rhvv9w42qka07narzt47da"`, lowercase letters and digits only.
Both twenty-two characters, two character sets, one collection. `updated_at` is
never different from `created_at`, so a field that exists to say when something
changed has never said it. A query parameter that is not a category answers 404
naming `"path": "/jokes/random"` -- the path that does exist and did answer. And
an empty search is `{"total": 0, "result": []}`, with the array called `result`,
singular, beside a total that is not.

## Sources

- Documentation: https://api.chucknorris.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve chucknorris     # run it
cauldron verify chucknorris -v # check every claim
```
