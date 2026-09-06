# Algolia

Emulates the Algolia API (1), for local development and tests.

**17 conformance cases, 2 checked against the live API.**

Two were struck live against latency-dsn.algolia.net on 2026-09-05, and both corrected this file. The casing was wrong -- Algolia sends "Invalid Application-ID or API key", and this had it as "Application-Id" throughout. More than casing: this had also claimed that omitting the application-id header gets a distinct, differently-worded failure from sending a bad key. Live, with a correct-looking key and no Application-Id header at all, Algolia sends the identical generic message -- it does not say which credential was the problem, and this file's belief that it did was never checked before now.

## What this Recipe found

Algolia's credentials are two separate headers, and neither is Authorization -- sending a bearer token gets nowhere. The record identifier is objectID, supplied by the caller rather than minted by Algolia, so code that reads hit.id after a search finds nothing.

The only fidelity gap: Cauldron doesn't implement search. A query returns everything in the index rather than the records actually matching the term, so relevance, ranking and highlighting aren't modelled -- only the response shape and the failure modes are.

## Sources

- Documentation: https://www.algolia.com/doc/rest-api/search
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve algolia     # run it
cauldron verify algolia -v # check every claim
```
