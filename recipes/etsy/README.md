# Etsy

Emulates the Etsy API (v3), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real shop; the refusal cases were struck live, unauthenticated, against openapi.etsy.com.

## What this Recipe found

Etsy requires two credentials on every request, and the failure this creates is the reason this Recipe exists: a bearer token is refused if `x-api-key` isn't sent alongside it, and the refusal doesn't say which header it wanted. The obvious read is that the token is bad, and hours go into reissuing a credential that was fine all along. The live probe corrected the shape of that refusal: it is a 403, not the 401 this file guessed, and the sentence is about the key's required format, not "missing or invalid api key" -- and the identical body answers a request with no bearer token at all, because x-api-key's absence dominates regardless. What a valid key beside a wrong bearer token answers was not observable without a real one, and that case stays exactly as documented.

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
