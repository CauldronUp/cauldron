# Stytch

Emulates the Stytch API (v1), for local development and tests.

**17 conformance cases, 3 checked against the live API on 2026-09-05.**

Stytch has a real test environment, so most cases here still cite documentation rather than the sandbox itself. The credential and routing failures needed no project at all, and checking them live found this Recipe's own auth model wrong on every axis.

## What this Recipe found

A user's identity has no single verified flag -- it's an array of factors (emails, phone numbers, providers), each independently verified, so a user can have a verified email and an unverified phone at the same time while `status` reads `active` regardless. A session is similarly a list of `authentication_factors` recording how it was obtained, so code requiring a particular login method has to read the array rather than trust the session object.

Deleting a user also doesn't end their sessions -- they remain valid until they expire on their own, so a revoked account keeps working for as long as its token has left.

## What checking it live found

Not a 401, and not the code or message this Recipe had claimed. No `Authorization` header at all is `invalid_authorization_header`, a 400; a well-formed Basic credential naming a project id Stytch does not have -- including the fixture's own project id with a wrong secret, which an existing case already sent -- is a different 400, `invalid_project_id_authentication`. A path nothing declares is its own 404, `route_not_found`, checked before the credential rather than after.

## Sources

- Documentation: https://stytch.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stytch     # run it
cauldron verify stytch -v # check every claim
```
