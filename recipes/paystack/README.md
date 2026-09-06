# Paystack

Emulates the Paystack API (1.0), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-08-31.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

They are here for a **negative** result. Every response
is `{"status": bool, "message": string, "data": ...}`, so there are two
statements of the outcome in every reply -- exactly the shape that lets a body
disagree with its own HTTP status, as Slack's does. On the whole surface
reachable without a key, they agree: `true` with 200, `false` with 401, and no
200 ever carrying `false`. Worth recording about an API built to let them
diverge. Its 401 does carry `"code": "invalid_Key"`, with a capital K in the
middle of a snake_case value, so anything matching that string has to know. And
an unknown country is a 200 with an empty array -- identical to a real country
it has no banks for, with the success flag saying `true` for both.

## Sources

- Documentation: https://paystack.com/docs/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve paystack     # run it
cauldron verify paystack -v # check every claim
```
