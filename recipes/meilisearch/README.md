# meilisearch

Emulates the meilisearch API (v1.11), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://www.meilisearch.com/docs/reference/api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve meilisearch     # run it
cauldron verify meilisearch -v # check every claim
```
