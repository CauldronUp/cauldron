# Postmark

Emulates the Postmark API (1.0), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A successful Postmark send carries `ErrorCode: 0` and `Message: "OK"`, so a client cannot tell success from failure by whether an error field is present -- it has to actually read the code. Plenty of integrations check for the absence of an error and get it right only by accident, since the field is always there either way.

## Sources

- Documentation: https://postmarkapp.com/developer/api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve postmark     # run it
cauldron verify postmark -v # check every claim
```
