# Flutterwave

Emulates the Flutterwave API (v3), for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-08-31.**

## What this Recipe found

They exist to be read against Paystack's. Both wrap
everything in `{status, message, data}`; Paystack's `status` is the boolean
`false` and Flutterwave's is the string `"error"`, so `if (!body.status)` treats
every Flutterwave failure as a success.

## Sources

- Documentation: https://developer.flutterwave.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve flutterwave     # run it
cauldron verify flutterwave -v # check every claim
```
