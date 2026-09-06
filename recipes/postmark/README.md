# Postmark

Emulates the Postmark API (1.0), for local development and tests.

**11 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real server token. The credential and routing checks were verified directly against api.postmarkapp.com on 2026-09-05.

## What this Recipe found

Checked live: no Server-Token header at all and a fictitious one answer the identical `{"ErrorCode":10,"Message":"Request does not contain a valid Server token."}` -- this file had assumed different wording for the same code. Routing also runs ahead of the credential: a wrong method on a real path answers a genuinely empty 405, zero bytes, with no token sent at all, while the same path with a real method and no token still gets the 401 above.

A successful Postmark send carries `ErrorCode: 0` and `Message: "OK"`, so a client cannot tell success from failure by whether an error field is present -- it has to actually read the code. Plenty of integrations check for the absence of an error and get it right only by accident, since the field is always there either way.

## Sources

- Documentation: https://postmarkapp.com/developer/api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve postmark     # run it
cauldron verify postmark -v # check every claim
```
