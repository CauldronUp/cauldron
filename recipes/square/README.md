# Square

Emulates the Square API (2026-01-22), for local development and tests.

**17 conformance cases, 4 checked against the live API on 2026-09-05.**

Most of this Recipe still cites documentation, since Square's own record-level behaviour needs a real account. What an unauthenticated request gets back does not, and checking it live confirmed the credential shape and found a routing claim this Recipe had never made.

## What writing this Recipe changed

Its money is an object in minor units rather than a number, so a client that
treats an amount as a scalar is wrong twice over.

## What checking it live found

The credential shape held up: an absent bearer and a garbage one answer byte-identical 401 `AUTHENTICATION_ERROR`/`UNAUTHORIZED`, checked before `Square-Version` is ever asked for. What was missing: a path nothing declares and a method a real path does not support both answer the same 404 `INVALID_REQUEST_ERROR`/`NOT_FOUND`, `"Resource not found."` -- there is no 405 anywhere in this API -- with no credential sent on either request, so routing is resolved before the credential is judged. `after_routing: true` says so now. A pre-existing case named "a bad token is refused" had also sent no token at all; the name is fixed and a case with an actual garbage token sits beside it.

## Sources

- Documentation: https://developer.squareup.com/reference/square
- Machine-readable description: https://raw.githubusercontent.com/square/connect-api-specification/master/api.json, last checked 2026-09-05
  `cauldron drift square` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve square     # run it
cauldron verify square -v # check every claim
```
