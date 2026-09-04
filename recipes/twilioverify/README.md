# twilioverify

Emulates the twilioverify API (v2), for local development and tests.

**20 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A wrong verification code is a 200. You post the digits, Twilio answers 200 OK, and whether they were right is a word inside the body -- `status` says `approved` or not, and the naive `if (!response.ok) throw` treats every wrong code as success and lets the wrong person in. There are also two competing fields for the answer, `status` and the legacy `valid` boolean, and Twilio's own schema calls `valid` legacy and tells you to read `status` instead -- but both are present on every response.

A verification is also consumed: Twilio deletes the record once it's approved, expired, or out of attempts, so checking a second time is a 404, identical to checking one that never existed. And five wrong codes in a row switches the failure from 200 to 429, so retry-on-429 middleware retries the one response that will never succeed.

## Sources

- Documentation: https://www.twilio.com/docs/verify/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve twilioverify     # run it
cauldron verify twilioverify -v # check every claim
```
