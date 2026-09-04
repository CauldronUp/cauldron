# Vercel

Emulates the Vercel API (v13), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The error envelope is thin: failures are nested under `error` with only a `code` and a `message`, no `type` and nothing else in the body, so code written for a richer error shape (Stripe's, say) reads `error.type` and finds nothing. Deployment state is a string enum where `READY` is not the only terminal value -- polling until state stops being `BUILDING` and assuming success ships a broken build, because there are failure states on the other side of that check too.

## Sources

- Documentation: https://vercel.com/docs/rest-api
- Machine-readable description: https://vercel.com/openapi.json, last checked 2026-08-31
  `cauldron drift vercel` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vercel     # run it
cauldron verify vercel -v # check every claim
```
