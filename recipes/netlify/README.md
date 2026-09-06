# Netlify

Emulates the Netlify API (v1), for local development and tests.

**14 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.netlify.com, no account and no key -- a missing credential and a wrong one both answer the identical body. This file declared an invented sentence and no `code`; the real answer is `{"code":401,"message":"Access Denied"}`. Fixed below.

## What this Recipe found

A Netlify deploy has an identifier before it has anything else -- the create response carries an id and `state: "new"` with no `deploy_ssl_url` at all, so code that creates a deploy and immediately reads its URL gets nothing. `ready` also is not `published`: a deploy can build successfully and be entirely ready without being the one actually served, since publishing is a separate act, and only the site's own `published_deploy_id`, not a field on the deploy itself, says which one is live.

The failure state is called `error`, with the actual reason living in a prose `error_message` and no error code at all -- so a script deciding whether to retry has to match against a sentence. A building deploy and a failed deploy both have no URL either, so absence alone does not distinguish them and the state field is the only way to tell them apart. Site ids and deploy ids also have deliberately different shapes, a UUID versus a hex string, even though both appear in the same URL paths, so passing one where the other belongs gets a 404 that says only that nothing was found.

Publishing a deploy and deleting a site are deliberately not modelled here, the same call this collection makes for Mercury and RingCentral: an action that changes what the public sees with no confirmation step does not belong in an emulator.

## Sources

- Documentation: https://open-api.netlify.com/
- Machine-readable description: https://open-api.netlify.com/openapi.json, last checked 2026-09-05
  `cauldron drift netlify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve netlify     # run it
cauldron verify netlify -v # check every claim
```
