# weaviate

Emulates the weaviate API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A batch import that half failed is a 200. Weaviate's own description says so: the request was processed, and each object's outcome is reported individually in an array inside the body, `result.status` of SUCCESS or FAILED per element. A client that checks the status code and moves on has imported some of its documents -- not none, which would be obvious, and not all, which is what it believes -- and the only symptom later is a search that quietly returns less than it should.

## Sources

- Documentation: https://weaviate.io/developers/weaviate/api/rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve weaviate     # run it
cauldron verify weaviate -v # check every claim
```
