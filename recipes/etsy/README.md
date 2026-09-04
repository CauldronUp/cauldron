# Etsy

Emulates the Etsy API (v3), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Etsy requires two credentials on every request, and the failure this creates is the reason this Recipe exists: a perfectly valid OAuth bearer token is refused with a plain 401 if `x-api-key` isn't sent alongside it, and the refusal doesn't say which header it wanted. The obvious read is that the token is bad, and hours go into reissuing a credential that was fine all along.

A few fields aren't what they look like either. Price is `{amount, divisor, currency_code}`, not a number -- reading it directly gives NaN. `count` on a listing is the total number matching, not the number of results in the current page, so a loop sized off it reads past the page it actually has. And a failure body is `{"error": "a sentence"}`: one string under a key shaped like an object, with no code to branch on beyond the HTTP status itself.

Two separate rate limits also run at once -- ten requests a second and ten thousand a day, reported across four headers -- so backing off the per-second budget does nothing for the daily one, which is the one that actually ends a bulk import mid-afternoon.

## Sources

- Documentation: https://developer.etsy.com/documentation/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve etsy     # run it
cauldron verify etsy -v # check every claim
```
