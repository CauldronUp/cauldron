# Algolia

Emulates the Algolia API (1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Algolia's credentials are two separate headers, and neither is Authorization -- sending a bearer token gets nowhere, and omitting the application-id header is a distinct, differently-worded failure from sending a bad key. The record identifier is objectID, supplied by the caller rather than minted by Algolia, so code that reads hit.id after a search finds nothing.

The only fidelity gap: Cauldron doesn't implement search. A query returns everything in the index rather than the records actually matching the term, so relevance, ranking and highlighting aren't modelled -- only the response shape and the failure modes are.

## Sources

- Documentation: https://www.algolia.com/doc/rest-api/search
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve algolia     # run it
cauldron verify algolia -v # check every claim
```
