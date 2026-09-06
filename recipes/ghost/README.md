# Ghost

Emulates the Ghost API (v5.0), for local development and tests.

**20 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real member-facing site, because publishing with an email filter sends real mail. The refusal cases were struck live, unauthenticated, against demo.ghost.io, a real public Ghost install.

## What this Recipe found

Publishing a post can send a real email, immediately and irreversibly. Setting `status: published` together with an `email_segment` doesn't just publish the post, it mails it to real members on the spot -- there's no sandbox to catch this because Ghost, like WordPress, is software rather than a hosted service, so testing means a second install that inevitably drifts from production.

A few shapes trip up an integration written by analogy with other APIs. Every collection is wrapped in a plural key named after the resource, even for a single item, so reading one post means `posts[0]`. The Admin API wants a JWT built from the API key, not the key itself, and sending the raw key looks exactly like a bad token rather than a malformed one. And content posted to the Admin API has to be Lexical (a JSON string), not HTML -- posting HTML directly succeeds and produces a silently empty post unless a `source` parameter says otherwise.

The live probe found Ghost's authentication failures split three ways where this file had modelled one: no credential at all is a 403 naming no authenticated user, a present-but-malformed one is a 400 `INVALID_JWT`, and this file's own guess (401 `UnauthorizedError`) was neither. Two narrower shapes were struck and left unmodelled -- a raw key with no "Ghost " prefix, and a well-formed JWT with no kid claim -- because a single credential pattern cannot host four different answers. An unrouted path and a wrong method on a real path both answer identically before authentication is ever consulted.

Status has four values, and `sent` is the one integrations don't expect: a `scheduled` post has a future `published_at` and isn't on the site yet, while a `sent` post is an email that was never on the site at all.

## Sources

- Documentation: https://ghost.org/docs/admin-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ghost     # run it
cauldron verify ghost -v # check every claim
```
