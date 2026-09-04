# Okta

Emulates the Okta API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A user's status in Okta is a lifecycle with seven states -- STAGED, PROVISIONED, ACTIVE, LOCKED_OUT, PASSWORD_EXPIRED, SUSPENDED, DEPROVISIONED -- and only one of them means the person can actually sign in. Code that just checks for the absence of DEPROVISIONED lets a suspended or locked-out user straight through.

A user's own attributes, like name and email, live under `profile`, and credential material is under a separate `credentials` key, so code that reads `user.email` directly off the object finds nothing. Okta also pages with a `Link` header, which Cauldron does not model, so cursor-based paging is not exercised here -- the states worth reproducing are the ones a live developer org cannot be made to sit in on demand anyway, like LOCKED_OUT or PASSWORD_EXPIRED.

## Sources

- Documentation: https://developer.okta.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve okta     # run it
cauldron verify okta -v # check every claim
```
