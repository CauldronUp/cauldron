# Plaid

Emulates the Plaid API (2020-09-14), for local development and tests.

**12 conformance cases, 1 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs an item this Recipe cannot fabricate. The credential shape itself was checked directly against sandbox.plaid.com, unauthenticated, on 2026-09-05.

## What writing this Recipe changed

Its error category is a separate field from its code, so branching on the code
alone tells a caller less than the response does.

## What checking it live found

The live host reads PLAID-CLIENT-ID and PLAID-SECRET as headers, exactly as this Recipe already modelled from the OpenAPI document's security block, even though that same document also lists `client_id`/`secret` as optional body properties -- a documentation ambiguity this Recipe had to actually test rather than assume. Two well-formed but fake headers answer 400 INVALID_API_KEYS with the exact sentence already on file. Omitting a header entirely, though, answers a different code and category (MISSING_FIELDS/INVALID_REQUEST) naming the body field Plaid wants -- "client_id", not "PLAID-CLIENT-ID" -- and this project's required-header mechanism can only substitute the header's own name, so that exact text is described rather than asserted as verified. Separately, the live probe also confirmed something the spec merely implied: the `status` field Plaid's schema documents on an error object is webhook-only and never appears on a direct synchronous failure, so this Recipe no longer serves one.

## Sources

- Documentation: https://plaid.com/docs/api
- Machine-readable description: https://raw.githubusercontent.com/plaid/plaid-openapi/master/2020-09-14.yml, last checked 2026-08-30
  `cauldron drift plaid` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve plaid     # run it
cauldron verify plaid -v # check every claim
```
