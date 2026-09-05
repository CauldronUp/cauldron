# Sentry

Emulates the Sentry API (0), for local development and tests.

**10 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real organisation. The credential check itself was verified directly against sentry.io on 2026-09-05.

## What writing this Recipe changed

Its entire error body is one string, with no object around it to read a field
from.

## What checking it live found

No Authorization header at all answers exactly what this file already had, "Authentication credentials were not provided.", but a syntactically fine but fictitious bearer answers a different sentence, "Invalid org token", which was not modelled before. A wrong method on a real path lands on the missing-credential sentence too, so the credential is checked first -- but a genuinely unrouted path does not: it answers a distinctive plain-text 404, "Route not found, did you forget a trailing slash?", with a caret pointing at the fix, needing no credential at all.

## Sources

- Documentation: https://docs.sentry.io/api
- Machine-readable description: https://raw.githubusercontent.com/getsentry/sentry-api-schema/main/openapi-derefed.json, last checked 2026-09-05
  `cauldron drift sentry` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sentry     # run it
cauldron verify sentry -v # check every claim
```
