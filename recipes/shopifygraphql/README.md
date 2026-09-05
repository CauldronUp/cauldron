# Shopify GraphQL

Emulates the Shopify GraphQL API (graphql), for local development and tests.

**11 conformance cases, 1 checked against the live API.**

Everything past one credential shape still cites documentation rather than an observation, because a real, running store needs an account. That one shape was checked directly against a shop subdomain Shopify has never issued, on 2026-09-05.

## What this Recipe found

Checked live: a POST with no Authorization credential at all answers exactly the shape-two failure this file already had, byte for byte, `{"errors":"[API] Invalid API key or access token (unrecognized login or wrong password)"}`. But a POST carrying a syntactically fine, fictitious token answers something else entirely, `{"errors":"Not Found"}` -- the store's existence is resolved once any credential is present, and only its total absence reaches the authentication check this file had documented for both. Whether a genuinely wrong token against a real shop reaches the 401 shape instead is not settled by a nonexistent-store probe, so that case stays documentation-sourced rather than verified.

GraphQL answers 200 when it refuses, so `if (response.ok)` is true for a request that did nothing. Failures arrive in three different shapes on this one API: a throttled request is 200 with `errors` as an array of objects; a bad token is 401 with `errors` as a bare string, so `errors[0].message` reads a stray character rather than throwing; and a mutation refused on business grounds is 200 with no top-level `errors` at all -- the failure is nested under the mutation's own name in `userErrors`, present and empty even on success.

Identifiers are also URLs (`gid://shopify/Product/1234567890`) rather than the numeric ids the REST API used, and a connection is `edges` and `node` rather than an array, so `data.products.map(...)` finds nothing.

## Sources

- Documentation: https://shopify.dev/docs/api/admin-graphql
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shopifygraphql     # run it
cauldron verify shopifygraphql -v # check every claim
```
