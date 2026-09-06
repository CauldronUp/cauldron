# Auth0

Emulates the Auth0 API (v2), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A user_id encodes which connection it came from -- auth0|abc is a database user, google-oauth2|123 came from Google -- so code that assumes the auth0 prefix breaks the moment somebody signs in with Google, and ids aren't interchangeable between connections. email_verified is a separate field from the email and stays false until someone clicks the link, and a tenant used only for testing tends to contain nothing but pre-verified users, so the branch handling an unverified one never runs until production.

One fidelity gap: Auth0 returns a bare array by default and switches to an object with totals when include_totals is set; Cauldron always returns the bare array.

## Sources

- Documentation: https://auth0.com/docs/api/management/v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve auth0     # run it
cauldron verify auth0 -v # check every claim
```
