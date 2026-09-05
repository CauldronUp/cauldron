# Vercel

Emulates the Vercel API (v13), for local development and tests.

**14 conformance cases, 3 checked against the live API on 2026-09-05.**

Most of this Recipe still cites documentation, since deployments and domains need a real account. The credential and routing shapes needed no account at all, and checking them live found this Recipe's own claim wrong.

## What this Recipe found

The error envelope is thin: most failures are nested under `error` with only a `code` and a `message`, no `type` and nothing else in the body, so code written for a richer error shape (Stripe's, say) reads `error.type` and finds nothing. Deployment state is a string enum where `READY` is not the only terminal value -- polling until state stops being `BUILDING` and assuming success ships a broken build, because there are failure states on the other side of that check too.

## What checking it live found

Credential failures are the exception to "nothing else in the body": no token at all answers `"The request is missing an authentication token"` with `missingToken: true`, and a present, wrong token answers `"Not authorized"` with `invalidToken: true` instead -- two different sentences and two different flags, not the one sentence this Recipe had for both. An existing case had sent no token at all while calling it "a bad token"; fixed now. A path nothing declares is a 404, resolved before either credential question is asked.

## Sources

- Documentation: https://vercel.com/docs/rest-api
- Machine-readable description: https://vercel.com/openapi.json, last checked 2026-08-31
  `cauldron drift vercel` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vercel     # run it
cauldron verify vercel -v # check every claim
```
