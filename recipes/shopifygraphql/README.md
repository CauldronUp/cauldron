# Shopify GraphQL

Emulates the Shopify GraphQL API (graphql), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

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
