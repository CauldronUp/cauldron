# meilisearch

Emulates the meilisearch API (v1.11), for local development and tests.

**28 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Adding a document to Meilisearch answers 202 with a task number, not the document, not an id, not even a count -- just a number and the word `enqueued`. The document is not searchable yet and might never be, because the task can still fail after that 202. A test that writes a document and immediately searches for it finds nothing, and reaching for a sleep just turns a race into a slower race; the only way to know a write landed is to poll `/tasks/{uid}` until it leaves `enqueued`.

Two more surprises sit right next to it: `estimatedTotalHits` is an estimate that changes between requests for the same query, so a paginator computing a page count from it produces a count that was only ever true once, and a search against a nonexistent index answers 404 while a search that simply matches nothing answers 200 with an empty array -- both look identical to code that only checks whether `hits` is empty. The primary key is also inferred from a collection's first document if never set explicitly, and inferring it wrong is permanent -- the whole index has to be deleted to fix it.

Auth failures were read from Meilisearch's own Rust source rather than a live call, since no reachable instance exists without standing up a project: a missing Authorization header and a wrong key answer different status codes, 401 versus 403, not just different sentences.

## Sources

- Documentation: https://www.meilisearch.com/docs/reference/api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve meilisearch     # run it
cauldron verify meilisearch -v # check every claim
```
