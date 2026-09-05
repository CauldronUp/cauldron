# Increase

Emulates the Increase API (v1), for local development and tests.

**14 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against api.increase.com, no account and no key. This file declared the auth failure as code `authentication_error`, "Your API key is invalid."; the real one is `invalid_api_key_error`, "Invalid API key.", with a nested reason and detail that differ between no credential and a wrong one. Also confirmed: the credential is checked before routing, so a path this file does not even declare answers the identical auth failure rather than a 404. Fixed below.

## What this Recipe found

A return is a brand new object, not a correction to the original transfer. The transfer that got returned still says it was submitted and still carries the amount that left -- nothing about it is marked wrong. The return sits beside it with its own id and its own date, often more than a week later, and code that reconciles by reading transfers alone never sees the money come back. The return's reason is also the only thing that says whether retrying is even legal: `insufficient_funds` may be retried, `no_account` may not, and `authorization_revoked` must never be retried, because doing so is a regulatory problem rather than a failed request -- yet all three arrive in the identical shape.

Like Mercury, Bill.com, Gusto and Deel in this collection, this Recipe deliberately models reading and not moving money: there is no endpoint here that creates a transfer, stated in the header rather than left for someone to discover the hard way. Amounts are integer minor units where the sign carries direction -- a credit is positive, a debit negative, on the same field -- so summing a column without reading the sign produces a number that's neither the money in nor the money out.

## Sources

- Documentation: https://increase.com/documentation/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve increase     # run it
cauldron verify increase -v # check every claim
```
