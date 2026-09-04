# Ghost

Emulates the Ghost API (v5.0), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Publishing a post can send a real email, immediately and irreversibly. Setting `status: published` together with an `email_segment` doesn't just publish the post, it mails it to real members on the spot -- there's no sandbox to catch this because Ghost, like WordPress, is software rather than a hosted service, so testing means a second install that inevitably drifts from production.

A few shapes trip up an integration written by analogy with other APIs. Every collection is wrapped in a plural key named after the resource, even for a single item, so reading one post means `posts[0]`. The Admin API wants a JWT built from the API key, not the key itself, and sending the raw key looks exactly like a bad token rather than a malformed one. And content posted to the Admin API has to be Lexical (a JSON string), not HTML -- posting HTML directly succeeds and produces a silently empty post unless a `source` parameter says otherwise.

Status has four values, and `sent` is the one integrations don't expect: a `scheduled` post has a future `published_at` and isn't on the site yet, while a `sent` post is an email that was never on the site at all.

## Sources

- Documentation: https://ghost.org/docs/admin-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ghost     # run it
cauldron verify ghost -v # check every claim
```
